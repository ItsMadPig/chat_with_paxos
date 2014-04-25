package clientrpc

type RemoteClientServer interface {
	//put methods here
	RegisterClient(*RegisterArgs, *RegisterReply) error
}

type ClientServer struct {
	// Embed all methods into the struct.
	RemoteClientServer
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemoteClientServer) RemoteClientServer {
	return &ClientServer{s}
}
