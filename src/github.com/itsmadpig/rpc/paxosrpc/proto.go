package paxosrpc

const (
	REJECT = iota //0
	OK     = iota //1
)

type PrepareArgs struct {
	ProposalNumber int
}

type PrepareReply struct {
	Value              string
	HighestAcceptedNum int
	Status             int
}

type AcceptArgs struct {
	Value          string
	ProposalNumber int
}

type AcceptReply struct {
	HighestSeen int
	Status      int
}

type CommitArgs struct {
	Value string
}

type CommitReply struct {
	Value string
}
