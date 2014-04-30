package loadbalancerrpc

const (
	INIT  = iota //0
	RETRY = iota //1
)
const (
	NotReady = iota //0
	OK       = iota //1
	MOSTFAIL = iota // 2
	INVALID  = iota
)
const (
	InitCliNum = 3
)

type RouteArgs struct {
	Attempt  int
	HostPort string
}

type RouteReply struct {
	HostPort string
	Status   int
}

type Node struct {
	HostPort string // The host:port address of the storage server node.
	NodeID   int    // The ID identifying this storage server node.
}

type RegisterArgs struct {
	ServerInfo Node
}
type RegisterReply struct {
	Status     int
	Servers    []Node
	Recovering bool
}
