package client

import (
	//"errors"
	"github.com/itsmadpig/rpc/loadbalancerrpc"
	"github.com/itsmadpig/rpc/serverrpc"
	//"net"
	"net/rpc"
	//"strings"
	//"strconv"
	"errors"
	"fmt"
	"time"
)

type pacClient struct {
	client         *rpc.Client
	serverHostPort string
	loadBalancer   *rpc.Client
}

func NewPacClient(serverHostPort string, port int) (PacClient, error) {
	cli, err := rpc.DialHTTP("tcp", serverHostPort)
	if err != nil {
		return nil, err
	}

	args := &loadbalancerrpc.RouteArgs{Attempt: 0, HostPort: ""}
	var reply loadbalancerrpc.RouteReply
	cli.Call("LoadBalancer.RouteToServer", args, &reply)
	for reply.Status != loadbalancerrpc.OK {
		if reply.Status == loadbalancerrpc.MOSTFAIL {
			return nil, errors.New("most servers failed")
		}
		fmt.Println("retrying to connect")
		time.Sleep(1000 * time.Millisecond)
		cli.Call("LoadBalancer.RouteToServer", args, &reply)
	}
	cli2, err := rpc.DialHTTP("tcp", reply.HostPort)
	if err != nil {
		fmt.Println("Server failed to respond")
		for err != nil {
			fmt.Println("trying to get new server")
			args.HostPort = reply.HostPort
			args.Attempt = loadbalancerrpc.RETRY
			cli.Call("LoadBalancer.RouteToServer", args, &reply)
			time.Sleep(time.Second)
			if reply.Status != loadbalancerrpc.OK {
				if reply.Status == loadbalancerrpc.MOSTFAIL {
					return nil, errors.New("most servers failed")
				}
			}
			cli2, err = rpc.DialHTTP("tcp", reply.HostPort)
		}
	}
	fmt.Println("server Hostport=", reply.HostPort)
	pac := new(pacClient)
	pac.client = cli2
	pac.serverHostPort = reply.HostPort
	pac.loadBalancer = cli
	return pac, nil
}

//if fail connection, do RouteToServer with failed HostPort
//if all fail, stop client

func (pc *pacClient) ReconnectToLB() error {
	fmt.Println("reconnect called")
	args := &loadbalancerrpc.RouteArgs{Attempt: loadbalancerrpc.RETRY, HostPort: pc.serverHostPort}
	var reply loadbalancerrpc.RouteReply
	pc.loadBalancer.Call("LoadBalancer.RouteToServer", args, &reply)
	for reply.Status != loadbalancerrpc.OK {
		if reply.Status == loadbalancerrpc.MOSTFAIL {
			return errors.New("most servers failed")
		}
		fmt.Println("retrying to connect")
		time.Sleep(1000 * time.Millisecond)
		pc.loadBalancer.Call("LoadBalancer.RouteToServer", args, &reply)
	}
	cli2, err := rpc.DialHTTP("tcp", reply.HostPort)
	if err != nil {
		fmt.Println("Server failed to respond")
		for err != nil {
			fmt.Println("trying to get new server")
			args := &loadbalancerrpc.RouteArgs{Attempt: loadbalancerrpc.RETRY, HostPort: reply.HostPort}
			pc.loadBalancer.Call("LoadBalancer.RouteToServer", args, &reply)
			time.Sleep(time.Second)
			if reply.Status != loadbalancerrpc.OK {
				if reply.Status == loadbalancerrpc.MOSTFAIL {
					return errors.New("most servers failed")
				}
			}
			cli2, err = rpc.DialHTTP("tcp", reply.HostPort)
		}
	}
	pc.serverHostPort = reply.HostPort
	pc.client = cli2
	return nil
}

func (pc *pacClient) MakeMove(direction string) error {
	fmt.Print("Trying to go ")
	fmt.Println(direction)
	args := new(serverrpc.MoveArgs)
	reply := new(serverrpc.MoveReply)
	args.Direction = direction
	err := pc.client.Call("PacmanServer.MakeMove", args, &reply)
	if err != nil {
		fmt.Println("server not responding.. trying to find new server")
		err = pc.ReconnectToLB()
		if err != nil {
			fmt.Println("all servers failed.. closing down..")
			return errors.New("No servers available")
		}
		pc.MakeMove(direction)
	}
	return nil
}
