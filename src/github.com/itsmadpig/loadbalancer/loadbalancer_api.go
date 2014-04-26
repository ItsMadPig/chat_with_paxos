package loadbalancer

import "github.com/itsmadpig/rpc/loadbalancerrpc"

type LoadBalancer interface {
	RouteToServer(*loadbalancerrpc.RouteArgs, *loadbalancerrpc.RouteReply) error
}
