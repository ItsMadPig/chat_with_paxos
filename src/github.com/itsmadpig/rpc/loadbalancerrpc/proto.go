package loadbalancerrpc

const (
	INIT  = iota //0
	RETRY = iota //1
)
const (
	OK    = iota //0
	NOTOK = iota //1
)

type RouteArgs struct {
	Attempt int
}

type RouteReply struct {
	HostPort string
	Status   int
}
