package master

import (
	"errors"
	"go.etcd.io/etcd/clientv3"
	pb "cluster-balance/api"
	"sync/atomic"
	"time"
)

type request struct {
	req   interface{}       //request
	reply chan interface{}  //
}

type addResourceRequest pb.ResourceSpec
type removeResourceRequest pb.ResourceSpec
type getResourceRequest pb.ResourceSpec

type Master struct {
	leader      string      //leader hosts
	title       int32       //1 : leader, 0: follower
	masterID    string      //it's the the only criterion of master
	config      *Config
	etcdConn    *clientv3.Client

	nodes       map[string]pb.NodeSpec
	resources   map[string]pb.ResourceSpec

	scheduler   Scheduler
	requests    chan request
	shutdown    chan bool
}

func NewMaster(config *Config, masterID string)(*Master, error){
	if  err := config.Validate(); err != nil{
		return nil, err
	}

	cfg := clientv3.Config{
		Endpoints:   config.EetcdHosts,
		DialTimeout: 5 * time.Second,
	}
	conn, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Master{
		masterID:   masterID,
		etcdConn:   conn,
		config:     config,
		nodes:      make(map[string]pb.NodeSpec),
		resources:  make(map[string]pb.ResourceSpec),
		requests:   make(chan request, 16),
		shutdown:   make(chan bool),
	}, nil
}

func (m *Master)isLeader() bool {
	return atomic.LoadInt32(&m.title) != 0
}

func (m *Master)addResource(spec pb.ResourceSpec) (*pb.NodeSpec, error){
	ns := &pb.NodeSpec{}


	return nil, nil
}

func (m *Master)dispatchCommand(req request){
	switch v := req.req.(type) {
	//just leader own the permission to add, remove Resource
	case addResourceRequest:
		if m.isLeader(){
			req.reply <- errors.New("not leader")
		}else{
			m.addResource(pb.ResourceSpec(v))
		}
	default:
		req.reply <- errors.New("unknown command")
	}

}

func (m *Master)startElect() *Election{
	return nil
}

func (m *Master)Run(){

MainLoop:
	for  {
		select {
		case req, ok := <-m.requests:
			if ok {
				m.dispatchCommand(req)
			}else{
				break MainLoop;
			}
		}
	}
}