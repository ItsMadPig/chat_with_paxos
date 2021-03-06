package serverrpc

const (
	NotReady = iota // NotReady = 0
	OK       = iota // OK = 1
)

type Node struct {
	HostPort string // The host:port address of the storage server node.
	NodeID   int    // The ID identifying this storage server node.
}

type TempArgs struct {
	Attempt int
}

type TempReply struct {
	HostPort string
	Status   int
}

type MoveArgs struct {
	Direction string
}

type MoveReply struct {
	Direction string
}

type GetArgs struct {
	ID string
}

type GetReply struct {
	Logs map[int]string
	ID   string
}
