package paxosrpc

type RemotePaxos interface {
	//put methods here
<<<<<<< HEAD
	Prepare(*PrepareArgs, *PrepareReply) error
	Accept(*AcceptArgs, *AcceptReply) error
	Commit(*CommitArgs, *CommitReply) error
=======
	Prepare(*PrepareArgs, *PrepareReply)
	Accept(*AcceptArgs, *AcceptReply)
	Commit(*CommitArgs, *CommitReply)
>>>>>>> e88f4092a443bf2bc7f6ab2b4e07b4f6f305684b
}

type Paxos struct {
	// Embed all methods into the struct.
	RemotePaxos
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemotePaxos) RemotePaxos {
	return &Paxos{s}
}
