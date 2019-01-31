package master

type Election interface {
	IsLeader() bool

}