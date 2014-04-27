package loadbalancer

import (
	"fmt"
	"github.com/itsmadpig/rpc/loadbalancerrpc"
	"net"
	"net/http"
	"net/rpc"
)

type loadBalancer struct {
	hostPort        string
	serverHostPorts []string
}

func NewLoadBalancer(hostPort string) (LoadBalancer, error) {
	loadBalancer := new(loadBalancer)
	loadBalancer.hostPort = hostPort
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return nil, err
	}
	// Wrap the tribServer before registering it for RPC.
	err = rpc.RegisterName("LoadBalancer", loadbalancerrpc.Wrap(loadBalancer))
	if err != nil {
		return nil, err
	}

	rpc.HandleHTTP()
	go http.Serve(listener, nil)

	return loadBalancer, nil
}

func (lb *loadBalancer) RouteToServer(args *loadbalancerrpc.RouteArgs, reply *loadbalancerrpc.RouteReply) error {
	if args.Attempt == loadbalancerrpc.INIT {
		//first time connecting

	} else {
		//second time connecting
	}
	pack := new(loadbalancerrpc.RouteReply)
	pack.Status = loadbalancerrpc.OK
	pack.HostPort = "localhost:9009"
	*reply = *pack
	return nil

}

func (cl *loadBalancer) RegisterServer(args *loadbalancerrpc.RegisterArgs, reply *loadbalancerrpc.RegisterReply) error {
	fmt.Println("HI")
	return nil
}
