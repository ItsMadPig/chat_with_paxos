package server

import (
	//"errors"
	"fmt"
	"github.com/itsmadpig/paxos"
	"github.com/itsmadpig/paxosTestWrap"
	"github.com/itsmadpig/rpc/loadbalancerrpc"
	"github.com/itsmadpig/rpc/paxosrpc"
	//"github.com/itsmadpig/rpc/paxoswraprpc"
	"github.com/itsmadpig/rpc/serverrpc"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type pacmanServer struct {
	masterServerHostPort string                 //master's hostport
	selfNode             *loadbalancerrpc.Node  //node of itself
	nodes                []loadbalancerrpc.Node // map of all nodes
	masterConn           *rpc.Client            //connection to master
	paxos                paxosWrap.PaxosWrap    //paxosTestWrap.PaxosWrap - for testing //paxos.Paxos - for not testing
	ID                   int
}

//two types of disconnection
//explicit disconnect (user stops game)
//implicit disconnect (user crashes)

//first do static paxos, then implement dynamic paxos (quorum dynamic)

//static paxos: wait till all connected, then only runs paxos if master server disconnects
//dynamic paxos: running paxos even before all nodes connect.

//if masterServerHostPort isn't empty, then it's a masterclient,
//else it's a slaveclient to start with
func NewServer(masterServerHostPort string, port int, nodeID int, test bool, flags []string) (PacmanServer, error) {
	pacmanServer := new(pacmanServer)
	pacmanServer.selfNode = new(loadbalancerrpc.Node)
	pacmanServer.nodes = make([]loadbalancerrpc.Node, loadbalancerrpc.InitCliNum)
	pacmanServer.masterServerHostPort = masterServerHostPort
	pacmanServer.ID = nodeID
	//if it's the slave client
	conn, err := rpc.DialHTTP("tcp", masterServerHostPort)
	if err != nil {
		return nil, err
	}
	pacmanServer.masterConn = conn
	args := &loadbalancerrpc.RegisterArgs{ServerInfo: loadbalancerrpc.Node{HostPort: net.JoinHostPort("localhost", strconv.Itoa(port)), NodeID: nodeID}}
	var reply loadbalancerrpc.RegisterReply

	err = pacmanServer.masterConn.Call("LoadBalancer.RegisterServer", args, &reply)
	fmt.Println("trying to connect , Host port:", port, " ready :", reply.Status)
	if err != nil {
		fmt.Println(err)
	}
	for reply.Status != loadbalancerrpc.OK {
		//if master server is still busy
		fmt.Println(reply.Status)
		fmt.Println("retrying to connect")
		time.Sleep(1000 * time.Millisecond)
		err = pacmanServer.masterConn.Call("LoadBalancer.RegisterServer", args, &reply)
	}

	pacmanServer.nodes = reply.Servers
	hostPorts := make([]string, len(pacmanServer.nodes))
	i := 0
	for _, node := range pacmanServer.nodes {
		hostPorts[i] = node.HostPort
		i++
	}

	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	if test {
		err = rpc.RegisterName("PacmanServer", serverrpc.Wrap(pacmanServer))
		if err != nil {
			return nil, err
		}
	} else {
		err = rpc.RegisterName("PacmanServer", serverrpc.Wrap(pacmanServer))
		if err != nil {
			return nil, err
		}
	}
	rpc.HandleHTTP()
	go http.Serve(listener, nil)
	if test {
		fmt.Println("Test Mode : On")
		pacmanServer.paxos, err = paxosWrap.NewPaxosWrap(net.JoinHostPort("localhost", strconv.Itoa(port)), nodeID, hostPorts, flags)
	} else {
		pacmanServer.paxos, err = paxos.NewPaxos(net.JoinHostPort("localhost", strconv.Itoa(port)), nodeID, hostPorts, false)
	}

	if err != nil {
		return nil, err
	}
	return pacmanServer, nil

}

func (cl *pacmanServer) GetLogs(args *serverrpc.GetArgs, reply *serverrpc.GetReply) error {
	thisReply := new(paxosrpc.GetReply)
	arg := new(paxosrpc.GetArgs)
	arg.ID = args.ID
	cl.paxos.GetLogs(arg, thisReply)
	pack := new(serverrpc.GetReply)
	pack.Logs = thisReply.Logs
	*reply = *pack
	return nil
}

func (cl *pacmanServer) MakeMove(args *serverrpc.MoveArgs, reply *serverrpc.MoveReply) error {
	direction := args.Direction
	fmt.Println("server got : ", direction)
	err := cl.paxos.RequestValue(direction)
	if err != nil {
		fmt.Println("request error")
	}
	pack := new(serverrpc.MoveReply)
	pack.Direction = direction
	*reply = *pack
	return nil

}
