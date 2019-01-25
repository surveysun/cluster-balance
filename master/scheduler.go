package master

import (
	pb "cluster-balance/api"
)

// Scheduler interface provides strategy to assign resource to a node.
// Must not return nil on success.
type Scheduler interface {
	Assign(view map[string]pb.NodeSpec, spec *pb.ResourceSpec) (*pb.NodeSpec, error)
	Rebalance(view map[string]pb.NodeSpec, resources map[string]pb.ResourceSpec) (map[string]string, error)
}

const rebalanceHighWatermark = 0.8

// QuotaBasedScheduler schedules resource to nodes with corresponding labels
// and makes sure node quotas are not exceeded.
type QuotaBasedScheduler struct {
}

func (s* QuotaBasedScheduler)Assign(view map[string]pb.NodeSpec, spec *pb.ResourceSpec) (*pb.NodeSpec, error){
	return nil, nil
}

func (s* QuotaBasedScheduler)Rebalance(view map[string]pb.NodeSpec, resources map[string]pb.ResourceSpec) (map[string]string, error){
	return nil, nil
}