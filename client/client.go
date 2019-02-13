package client

import (
	"cluster-balance/comm"
	pb "cluster-balance/api"
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
	"encoding/json"
)

type Worker interface {
	HanderDoPut(key, value []byte) error
	HanderDoDelete(key, value []byte) error
}

type Client struct {
	clientID    string
	config      *comm.Config
	etcd        *comm.EtcdHander
	node		*pb.NodeSpec
	woker       Worker

	shutdown    chan bool
	close       chan struct{}
	createTime  int64
	workpath    string
}

func NewClient(config *comm.Config, node *pb.NodeSpec, woker Worker) (*Client, error){
	if err := config.Validate(); err != nil {
		return nil, err
	}

	etcd, err := comm.NewEtcdHander(config.EetcdHosts)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientID:   node.Id,
		config:     config,
		etcd:       etcd,
		woker:      woker,
		workpath:   config.GetNodesPath() + "/" + node.Id,
		node:		node,
	}, nil
}

func (c *Client)mangerExist() bool {
	kvs, err := c.etcd.GetAll(c.config.GetMasterPath())
	if len(kvs) > 0 && err == nil {
		return true
	}

	return false
}

func (c *Client)register() error {
	if c.node == nil {
		return errors.New("node config is error")
	}

	if !c.mangerExist() {
		return errors.New("this is no manger run")
	}

	data, err := proto.Marshal(c.node)
	if err != nil{
		return errors.New("Marshal node info failed")
	}

	err = c.etcd.Put( c.config.GetNodesPath() + "/" + c.clientID, string(data))
	if err != nil {
		return ErrRegister
	}

	return nil
}

func (c *Client)heartBeat(status comm.NodeStatus) error {
	heartBeat := comm.HeartBeat{
		Node:   c.clientID,
		Status: status,
		UpdateTime: time.Now().UTC().UnixNano()/1000000,      //毫秒
		CreateTime: c.createTime,
	}

	str, err := json.Marshal(&heartBeat)
	if err != nil{
		return err
	}

	err = c.etcd.Put(c.config.GetClusterPath(), string(str))
	return err
}

func (c *Client)Run() error {
	tickC := time.NewTimer(time.Second * 10)

	err := c.register()
	if err != nil {
		return err
	}

MainLoop:
	for{
		select {
		case <-c.close:
			break MainLoop;
		case <-tickC.C:
			tickC.Reset(time.Second * 10)

			c.heartBeat(comm.STATUS_RUNING)
		default:
			c.etcd.Watch(c.workpath, c.woker.HanderDoPut, c.woker.HanderDoDelete)
		}
	}

	c.shutdown <- true
	return nil
}

func (c *Client)Stop(){
	c.close <- struct{}{}
	<-c.shutdown
}