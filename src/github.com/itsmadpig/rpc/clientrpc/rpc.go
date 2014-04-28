package clientrpc

type RemotePacClient interface {
	//put methods here
}

type PacClient struct {
	// Embed all methods into the struct.
	RemotePacClient
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemotePacClient) RemotePacClient {
	return &PacClient{s}
}
