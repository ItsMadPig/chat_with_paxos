package paxos

import "github.com/itsmadpig/rpc/paxosrpc"

type Paxos interface {
	Prepare(*paxosrpc.PrepareArgs, *paxosrpc.PrepareReply) error
	Accept(*paxosrpc.AcceptArgs, *paxosrpc.AcceptReply) error
	Commit(*paxosrpc.CommitArgs, *paxosrpc.CommitReply) error
}
