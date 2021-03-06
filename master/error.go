package master

import "errors"

var(
	ErrNoNode                   = errors.New("no available nodes")
	ErrNodeResourceExhausted    = errors.New("node resource exhausted")
	ErrNotLeader                = errors.New("not the leader")
	ErrNotFindResource          = errors.New("not find the resource")
	ErrResourceExits            = errors.New("resources already exists")
)