package client

import (
	"errors"
	"fmt"
	"github.com/itsmadpig/rpc/loadbalancerrpc"
	"github.com/itsmadpig/rpc/serverrpc"
	"net/rpc"
	//"strconv"
	"time"
)

type pacClient struct {
	serverConn     *rpc.Client
	loadHostPort   string
	loadBalancer   *rpc.Client
	serverHostPort string
	ID             string
	logs           map[int]string
}

func NewPacClient(loadHostPort string, port int, ID string) (PacClient, error) {
	pac := new(pacClient)
	pac.loadHostPort = loadHostPort
	pac.ID = ID
	pac.logs = make(map[int]string)
	cli, err := rpc.DialHTTP("tcp", loadHostPort)
	if err != nil {
		return nil, err
	}

	pac.loadBalancer = cli

	args := &loadbalancerrpc.RouteArgs{Attempt: loadbalancerrpc.INIT, HostPort: ""}
	var reply loadbalancerrpc.RouteReply
	cli.Call("LoadBalancer.RouteToServer", args, &reply)

	for reply.Status == loadbalancerrpc.NotReady {
		fmt.Println("retrying to connect")
		time.Sleep(1000 * time.Millisecond)
		err = cli.Call("LoadBalancer.RouteToServer", args, &reply)

	}
	if reply.Status == loadbalancerrpc.MOSTFAIL {
		return nil, err
	}

	//connect to server
	cli2, err := rpc.DialHTTP("tcp", reply.HostPort)
	pac.serverConn = cli2
	pac.serverHostPort = reply.HostPort
	if err != nil {
		err1 := pac.ReconnectToLB()
		/*for err != nil {
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
		}*/
		if err1 != nil {
			fmt.Println("SERVER ALL FAILED")
			return nil, errors.New("reconnect fail, most servers dead")
		}
	}
	pac.GetLogs()
	go pac.RefreshTimer()
	return pac, nil
}

//if fail connection, do RouteToServer with failed HostPort
//if all fail, stop client\

func (pc *pacClient) RefreshTimer() {
	for {
		select {
		case _ = <-time.After(time.Second):
			pc.GetLogs()
		}
	}
}

func (pc *pacClient) ReconnectToLB() error {
	fmt.Println("reconnect called")
	args := &loadbalancerrpc.RouteArgs{Attempt: loadbalancerrpc.RETRY, HostPort: pc.serverHostPort}
	reply := new(loadbalancerrpc.RouteReply)
	pc.loadBalancer.Call("LoadBalancer.RouteToServer", args, &reply)

	if reply.Status == loadbalancerrpc.MOSTFAIL {
		fmt.Println("SERVER ALL FAILED")
		return errors.New("reconnect fail, most servers dead")
	}

	serverConn, err := rpc.DialHTTP("tcp", reply.HostPort)

	for err != nil {
		fmt.Println("trying to get new server")
		args := &loadbalancerrpc.RouteArgs{Attempt: loadbalancerrpc.RETRY, HostPort: reply.HostPort}
		pc.loadBalancer.Call("LoadBalancer.RouteToServer", args, &reply)

		if reply.Status == loadbalancerrpc.MOSTFAIL {
			fmt.Println("SERVER ALL FAILED")
			return errors.New("reconnect fail, most servers dead")
		}

		serverConn, err = rpc.DialHTTP("tcp", reply.HostPort)
	}
	pc.serverHostPort = reply.HostPort
	pc.serverConn = serverConn
	return nil
}

func (pc *pacClient) isNewMessage(newMap map[int]string) bool {
	for index, value := range newMap {
		thisVal, ok := pc.logs[index]
		if !ok {
			pc.logs = newMap
			return true
		}
		if value != thisVal {
			pc.logs = newMap
			return true
		}
	}
	return false
}

func (pc *pacClient) GetLogs() map[int]string {

	reply := new(serverrpc.GetReply)
	args := new(serverrpc.GetArgs)
	pc.serverConn.Call("PacmanServer.GetLogs", args, &reply)
	if pc.isNewMessage(reply.Logs) == false {
		return pc.logs
	}
	fmt.Println("**************************SCREEN*************************")
	for _, value := range pc.logs {
		fmt.Println(value)
	}
	fmt.Println("***********************END OF SCREEN ********************")
	fmt.Print("Message : ")
	return pc.logs
}

func (pc *pacClient) MakeMove(direction string) error {
	fmt.Println(pc.ID, ":", direction)
	args := new(serverrpc.MoveArgs)
	reply := new(serverrpc.MoveReply)
	args.Direction = pc.ID + ":" + direction
	err := pc.serverConn.Call("PacmanServer.MakeMove", args, &reply)
	if err != nil {
		err = pc.ReconnectToLB()
		if err != nil {
			fmt.Println("all servers failed.. closing down..")
			return errors.New("No servers available")
		}
		pc.MakeMove(direction)
	}
	return nil
}
