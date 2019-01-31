package client

import "cluster-balance/comm"

type Client struct {
	config   *comm.Config
	etcd     *comm.EtcdHander
} 
