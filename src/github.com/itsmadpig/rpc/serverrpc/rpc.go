package serverrpc

type RemoteServer interface {
	//put methods here
	Temp(*TempArgs, *TempReply) error
}

type PacmanServer struct {
	// Embed all methods into the struct.
	RemoteServer
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemoteServer) RemoteServer {
	return &PacmanServer{s}
}
