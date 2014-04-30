package paxoswraprpc

import (
	"github.com/itsmadpig/rpc/paxosrpc"
)

type RemotePaxosWrapper interface {
	//put methods here
	Prepare(*paxosrpc.PrepareArgs, *paxosrpc.PrepareReply) error
	Accept(*paxosrpc.AcceptArgs, *paxosrpc.AcceptReply) error
	Commit(*paxosrpc.CommitArgs, *paxosrpc.CommitReply) error
	GetLogs(*paxosrpc.GetArgs, *paxosrpc.GetReply) error
}

type PaxosWrapper struct {
	// Embed all methods into the struct.
	RemotePaxosWrapper
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemotePaxosWrapper) RemotePaxosWrapper {
	return &PaxosWrapper{s}
}
