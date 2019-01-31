package master

import (
	pb "cluster-balance/api"
	"cluster-balance/comm"
	"errors"
	"github.com/golang/protobuf/proto"
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
	status    map[string]string

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

func (m *Master) deleteNodeSpec(nod string, spec string) bool {
	soures := m.nodes[nod].GetStatus()
	if soures == nil || soures.Quotas == nil  {
		return false
	}
	delete(soures.Quotas, spec)
	return true
}

func (m *Master) addResource(spec pb.ResourceSpec) (*pb.NodeSpec, error) {
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

	err = m.etcd.Put( m.config.GetNodesPath() + "/" + spec.AssignedNode + "/" + spec.GetId() , string(data))
	if err != nil{
		return nil, err
	}
	// resource/task_id use to backups
	err = m.etcd.Put( m.config.GetResourcesPath() + "/" + spec.GetId() , string(data))
	if err != nil{
		return nil, err
	}

	newns := proto.Clone(ns).(*pb.NodeSpec)
	return newns, nil
}

func (m *Master) deleteResource(spec pb.ResourceSpec) (*pb.NodeSpec, error) {
	if m.isLeader() {
		if nodeId, ok := m.findNodeBySpec(spec); ok {
			if nodeId != ""{
				_, err := m.etcd.Delete( m.config.GetNodesPath() + "/" + nodeId + "/" + spec.GetId())
				if err != nil{
					return nil, err
				}

				//delete node task info
				m.deleteNodeSpec(nodeId, spec.GetId())
			}

			_, err := m.etcd.Delete( m.config.GetResourcesPath() + "/" + spec.GetId())
			if err != nil{
				return nil, err
			}

			//delete resoures map task info

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
}
