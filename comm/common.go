package comm

type NodeStatus int

const (
	STATUS_PREPARE NodeStatus = iota
	STATUS_RUNING
	STATUS_OFF
)
//secend
const HeartbeatTime = 10

type HeartBeat struct {
	Node       string     `json:"node"`
	Status     NodeStatus `json:"status"`
	UpdateTime int64      `json:"update_time"`
	CreateTime int64      `json:"create_time"`
}

type MasterInfo  HeartBeat