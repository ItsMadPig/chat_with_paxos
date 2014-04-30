package paxosrpc

const (
	REJECT      = 1 //1
	OK          = 2 //2
	OldInstance = 3 //3
	Prepareres  = 4 //4
)

type PrepareArgs struct {
	ProposalNumber int
	Round          int
	HostPort       string
}

type PrepareReply struct {
	Value              string
	HighestAcceptedNum int
	Status             int
}

type AcceptArgs struct {
	Value          string
	ProposalNumber int
	Round          int
}

type AcceptReply struct {
	HighestSeen int
	Status      int
	Value       string
}

type CommitArgs struct {
	Value string
	Round int
}

type CommitReply struct {
	Value string
}

type GetArgs struct {
	ID string
}

type GetReply struct {
	Logs map[int]string
	ID   string
}
