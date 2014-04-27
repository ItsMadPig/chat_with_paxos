package paxosrpc

const (
	REJECT = iota //0
	OK     = iota //1
)

type PrepareArgs struct {
	ProposalNumber int
}

type PrepareReply struct {
	Value            string
	ProposalAwaiting int
	Status           int
}

type AcceptArgs struct {
	Value          string
	ProposalNumber int
}

type AcceptReply struct {
	ProposalCommited int
}

type RequestArgs struct {
	Value string
}
