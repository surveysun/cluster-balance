package client

import (
	pb "cluster-balance/api"
	"cluster-balance/comm"
	"cluster-balance/utils"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

type Worker interface {
	HanderDoPut(key, value []byte) error
	HanderDoDelete(key, value []byte) error
}

type Client struct {
	clientID string
	config   *comm.Config
	etcd     *comm.EtcdHander
	node     *pb.NodeSpec
	woker    Worker

	shutdown   chan bool
	close      chan struct{}
	createTime int64
	workpath   string
}

func NewClient(config *comm.Config, node *pb.NodeSpec, woker Worker) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	etcd, err := comm.NewEtcdHander(config.EetcdHosts)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientID: node.Id,
		config:   config,
		etcd:     etcd,
		woker:    woker,
		close:	  make(chan struct{}),
		shutdown: make(chan bool),
		workpath: config.GetNodesPath() + "/" + node.Id,
		node:     node,
	}, nil
}

func (c *Client) mangerExist() bool {
	kvs, err := c.etcd.GetAll(c.config.GetMasterPath())
	if len(kvs) > 0 && err == nil {
		for _, kv := range kvs {
			var master comm.MasterInfo
			err := json.Unmarshal([]byte(kv.Value), &master)
			if err != nil{
				log.Warn("Unmarshal failed, err:",err)
				continue
			}
			if master.Status == comm.STATUS_RUNING && utils.CheckTime(master.UpdateTime, comm.HeartbeatTime){
				return true
			}
		}
	}

	return false
}

func (c *Client) register() error {
	if c.node == nil {
		return errors.New("node config is error")
	}

	if !c.mangerExist() {
		return errors.New("this is no manger run")
	}

	data, err := proto.Marshal(c.node)
	if err != nil {
		return errors.New("Marshal node info failed")
	}

	err = c.etcd.Put(c.config.GetNodesPath()+"/"+c.clientID, string(data))
	if err != nil {
		return ErrRegister
	}

	return nil
}

func (c *Client) heartBeat(status comm.NodeStatus) error {
	if status == comm.STATUS_PREPARE {
		c.createTime = time.Now().UTC().UnixNano() / 1000000 //毫秒
	}

	heartBeat := comm.HeartBeat{
		Node:       c.clientID,
		Status:     status,
		UpdateTime: time.Now().UTC().UnixNano() / 1000000, //毫秒
		CreateTime: c.createTime,
	}

	str, err := json.Marshal(&heartBeat)
	if err != nil {
		return err
	}

	err = c.etcd.Put(c.config.GetClusterPath()+"/"+c.clientID, string(str))
	return err
}

func (c *Client) Run() error {
	tickC := time.NewTimer(time.Second * comm.HeartbeatTime)

	err := c.register()
	if err != nil {
		return err
	}
	watchChan := c.etcd.Watch(c.workpath)
MainLoop:
	for {
		select {
		case <-c.close:
			break MainLoop
		case <-tickC.C:
			tickC.Reset(time.Second * 10)
			c.heartBeat(comm.STATUS_RUNING)
		case resp := <- watchChan:
			for _, ev := range resp.Events {
				if ev.Type == comm.EventTypePut {
					c.woker.HanderDoPut(ev.Kv.Key, ev.Kv.Value)
				} else if ev.Type == comm.EventTypeDelete {
					c.woker.HanderDoDelete(ev.Kv.Key, ev.Kv.Value)
				} else {
					log.Info("error event.Type:", ev.Type, " key:", string(ev.Kv.Key))
				}
			}
		default:
		}
	}

	c.etcd.CloseWatchAll()
	c.shutdown <- true
	return nil
}

func (c *Client) Stop() {
	c.close <- struct{}{}
	<-c.shutdown
}
