package paxosrpc

type PrepareArgs struct {
	ProposalNumber int
}

type PrepareReply struct {
	Value            string
	ProposalAwaiting int
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
