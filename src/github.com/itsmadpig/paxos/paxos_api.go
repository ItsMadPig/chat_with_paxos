package paxos

import "github.com/itsmadpig/rpc"

type Paxos interface {
	Prepare(*paxosrpc.PrepareArgs, *paxosrpc.PrepareReply) error
	Accept(*paxosrpc.AcceptArgs, *paxosrpc.AcceptReply) error
	RequestValue(*paxosrpc.RequestArgs, *paxosrpc.RequestReply) error
}
