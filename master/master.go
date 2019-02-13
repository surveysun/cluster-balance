package master

import (
	pb "cluster-balance/api"
	"cluster-balance/comm"
	"errors"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

type request struct {
	req   interface{}      //request
	reply chan interface{} //
}

type addResourceRequest pb.ResourceSpec
type removeResourceRequest pb.ResourceSpec
type getResourceRequest pb.ResourceSpec


type Master struct {
	leader   string //leader hosts
	title    int32  //1 : leader, 0: follower
	masterID string //it's the the only criterion of master
	config   *comm.Config
	etcd     *comm.EtcdHander

	nodes     map[string]pb.NodeSpec
	resources map[string]pb.ResourceSpec
	status    map[string]comm.HeartBeat

	scheduler Scheduler
	requests  chan request
	shutdown  chan bool
}

func NewMaster(config *comm.Config, masterID string) (*Master, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	etcd, err := comm.NewEtcdHander(config.EetcdHosts)
	if err != nil {
		return nil, err
	}

	return &Master{
		masterID:  masterID,
		etcd:      etcd,
		config:    config,
		nodes:     make(map[string]pb.NodeSpec),
		resources: make(map[string]pb.ResourceSpec),
		requests:  make(chan request, 16),
		shutdown:  make(chan bool),
	}, nil
}

func (m *Master) isLeader() bool {
	return atomic.LoadInt32(&m.title) != 0
}

// logger returns a new log entry with manger id label.
// log.Entry is not thread-safe.
func (m *Master) logger() *log.Entry {
	return log.WithField("managerID", m.masterID)
}

func (m *Master) findNodeBySpec(spec pb.ResourceSpec) ( string, bool) {
	for k, v := range m.resources {
		if k == spec.GetId() {
			if _, ok := m.nodes[v.AssignedNode]; ok {
				return v.AssignedNode, true
			}
			return "", true
		}
	}
	return "", false
}

func (m *Master) deleteMapNodeSpec(nod string, spec string) bool {
	soures := m.nodes[nod].GetStatus()
	if soures == nil || soures.Quotas == nil  {
		return false
	}

	if _, ok := soures.Quotas[spec]; ok {
		delete(soures.Quotas, spec)
		return true
	}

	return false
}

func (m *Master) deleteMapSpec(spec string) bool {
	if  _, ok := m.resources[spec]; ok {
		delete(m.resources, spec)
		return true
	}
	return false
}

func (m *Master) exitsResource(spec string) bool{
	if _, ok := m.resources[spec]; ok {
		return true
	}
	return  false
}

func (m *Master) addResource(spec pb.ResourceSpec) (*pb.NodeSpec, error) {
	if m.exitsResource(spec.GetId()) {
		return nil, ErrResourceExits
	}

	ns := &pb.NodeSpec{}
	if m.isLeader() {
		ns, err:= m.scheduler.Assign( m.nodes, &spec)
		if err != nil{
			return nil, err
		}

		if ns == nil{
			panic("scheduler return nil")
		}

		spec.AssignedNode = ns.GetId()
	}else{
		return nil, ErrNotLeader
	}

	spec.CreationTime = time.Now().UnixNano()
	data, err := proto.Marshal(&spec)
	if err != nil {
		return nil, err
	}

	// resource/task_id use to backups
	err = m.etcd.Put( m.config.GetResourcesPath() + "/" + spec.GetId() , string(data))
	if err != nil{
		return nil, err
	}
	m.resources[spec.GetId()] = spec
	err = m.etcd.Put( m.config.GetNodesPath() + "/" + spec.AssignedNode + "/" + spec.GetId() , string(data))
	if err != nil{
		return nil, err
	}

	ns.GetStatus().Quotas[spec.GetId()] = 1

	m.logger().Info("assign ", spec.GetId(), " ---->", ns.GetId())

	newns := proto.Clone(ns).(*pb.NodeSpec)
	return newns, nil
}

func (m *Master) deleteResource(spec pb.ResourceSpec) (*pb.NodeSpec, error) {
	if m.isLeader() {
		if nodeId, ok := m.findNodeBySpec(spec); ok {
			//first delete backups Resources
			_, err := m.etcd.Delete( m.config.GetResourcesPath() + "/" + spec.GetId())
			if err != nil{
				return nil, err
			}
			//delete resoures map task info
			m.deleteMapSpec(spec.GetId())

			if nodeId != ""{
				_, err := m.etcd.Delete( m.config.GetNodesPath() + "/" + nodeId + "/" + spec.GetId())
				if err != nil{
					return nil, err
				}

				//delete node task info
				m.deleteMapNodeSpec(nodeId, spec.GetId())
				return nil, nil
			}

			node, _ := m.nodes[nodeId]
			return &node, nil
		}
		return nil, ErrNotFindResource
	}else {
		return nil, ErrNotLeader
	}
}

func (m *Master) dispatchCommand(req request) {
	switch v := req.req.(type) {
	//just leader own the permission to add, remove Resource
	case addResourceRequest:
		if !m.isLeader() {
			req.reply <- ErrNotLeader
		} else {
			m.addResource(pb.ResourceSpec(v))
		}
	default:
		req.reply <- errors.New("unknown command")
	}

}

func (m *Master) startElect() *Election {
	return nil
}

func (m *Master) Run() {

MainLoop:
	for {
		select {
		case req, ok := <-m.requests:
			if ok {
				m.dispatchCommand(req)
			} else {
				break MainLoop
			}
		}
	}

	//elector Resign

	m.etcd.CloseWatchAll()
	m.logger().Info("manager exiting")
	atomic.StoreInt32(&m.title, 0)


	m.shutdown <- true
}

func (m *Master) AddResource(spec *pb.ResourceSpec) (*pb.NodeSpec, error){
	if spec.GetId() == "" {
		return nil, errors.New("invalid resource id")
	}
	//must be set by worker
	if spec.GetStatus() != nil {
		return nil, errors.New("resource status must be nil")
	}

	reply := make(chan interface{}, 1)
	m.requests <- request{req: addResourceRequest(*spec), reply: reply}
	r := <-reply
	if err, ok := r.(error); ok {
		return nil, err
	}
	return r.(*pb.NodeSpec), nil
}

func (m *Master) RemoveResource(spec *pb.ResourceSpec) error {
	reply := make(chan interface{}, 1)
	m.requests <- request{req: removeResourceRequest(*spec), reply: reply}
	r := <-reply
	if r == nil {
		return nil
	}
	return r.(error)
}

func (m *Master) Close(){
	close(m.requests)
	<-m.shutdown
}