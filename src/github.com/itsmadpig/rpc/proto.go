package serverrpc

const (
	OK       = iota // OK = 0
	NotReady = iota // NotReady = 1
)
const (
	InitCliNum = 8
)

type Node struct {
	HostPort string // The host:port address of the storage server node.
	NodeID   uint32 // The ID identifying this storage server node.
}

type RegisterArgs struct {
	ServerInfo Node
}
type RegisterReply struct {
	Status  int
	Servers []Node
}
