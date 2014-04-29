package paxosrpc

const (
	REJECT      = iota //0
	OK          = iota //1
	OldInstance = iota //2
	Prepareres  = iota //3
)

type PrepareArgs struct {
	ProposalNumber int
	Round          int
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
