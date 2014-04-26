package loadbalancerrpc

type RemoteLoadBalancer interface {
	//put methods here
	RouteToServer(*RouteArgs, *RouteReply) error
}

type LoadBalancer struct {
	// Embed all methods into the struct.
	RemoteLoadBalancer
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemoteLoadBalancer) RemoteLoadBalancer {
	return &LoadBalancer{s}
}
