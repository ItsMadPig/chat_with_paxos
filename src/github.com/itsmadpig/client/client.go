package client

import (
	//"errors"
	"github.com/itsmadpig/rpc/loadbalancerrpc"
	//"github.com/itsmadpig/rpc/serverrpc"
	//"net"
	"net/rpc"
	//"strings"
	//"strconv"
	"errors"
	"fmt"
	"time"
)

type pacClient struct {
	client *rpc.Client
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
			return nil, errors.New("most servers failed alread")
		}
		fmt.Println("retrying to connect")
		time.Sleep(1000 * time.Millisecond)
		cli.Call("LoadBalancer.RouteToServer", args, &reply)
	}
	cli2, err := rpc.DialHTTP("tcp", reply.HostPort)
	if err != nil {
		return nil, err
	}
	fmt.Println("server Hostport=", reply.HostPort)
	pac := new(pacClient)
	pac.client = cli2
	return pac, nil
}

//if fail connection, do RouteToServer with failed HostPort
//if all fail, stop client
