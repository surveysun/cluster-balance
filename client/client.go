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
	HanderDoDelte(key, value []byte) error
}

type Client struct {
	clientID    string
	config      *comm.Config
	etcd        *comm.EtcdHander
	woker       Worker

	shutdown    chan bool
	close       chan struct{}
	createTime  int64
	workpath    string
} 

func NewClient(config *comm.Config, clientId string, woker Worker) (*Client, error){
	if err := config.Validate(); err != nil {
		return nil, err
	}

	etcd, err := comm.NewEtcdHander(config.EetcdHosts)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientID:   clientId,
		config:     config,
		etcd:       etcd,
		woker:      woker,
		workpath:   config.GetNodesPath() + "/" + clientId,
	}, nil
}

func (c *Client)Register(labels map[string]string, quotas map[string]uint64) error {
	if labels == nil{
		return errors.New("the Node labels is nil")
	}

	if quotas == nil {
		return errors.New("the Node quotas is nil")
	}

	node := pb.NodeSpec{
		ClusterId:  c.config.ClusterID,
		Id:         c.clientID,
		Labels:     labels,
		Quotas:     quotas,
		Version:    1,
	}

	data, err := proto.Marshal(&node)
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

func (c *Client)Run(){
	tickC := time.NewTimer(time.Second * 10)
MainLoop:
	for{
		select {
		case <-c.close:
			break MainLoop;
		case <-tickC.C:
			tickC.Reset(time.Second * 10)

			c.heartBeat(comm.STATUS_RUNING)
		default:
			c.etcd.Watch(c.workpath, c.woker.HanderDoPut, c.woker.HanderDoDelte)
		}
	}

	c.shutdown <- true
}

func (c *Client)Stop(){
	c.close <- struct{}{}
	<-c.shutdown
}