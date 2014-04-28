package server

import (
	//"errors"
	"fmt"
	"github.com/itsmadpig/paxos"
	"github.com/itsmadpig/rpc/loadbalancerrpc"
	"github.com/itsmadpig/rpc/serverrpc"
	"githum.com/itsmadpig/rpc/paxosrpc"
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
	paxos                paxos.Paxos
}

//two types of disconnection
//explicit disconnect (user stops game)
//implicit disconnect (user crashes)

//first do static paxos, then implement dynamic paxos (quorum dynamic)

//static paxos: wait till all connected, then only runs paxos if master server disconnects
//dynamic paxos: running paxos even before all nodes connect.

//if masterServerHostPort isn't empty, then it's a masterclient,
//else it's a slaveclient to start with
func NewServer(masterServerHostPort string, port int, nodeID uint32) (PacmanServer, error) {
	pacmanServer := new(pacmanServer)
	pacmanServer.selfNode = new(loadbalancerrpc.Node)
	pacmanServer.nodes = make([]loadbalancerrpc.Node, loadbalancerrpc.InitCliNum)
	pacmanServer.masterServerHostPort = masterServerHostPort
	//if it's the slave client
	conn, err := rpc.DialHTTP("tcp", masterServerHostPort)
	if err != nil {
		return nil, err
	}
	pacmanServer.masterConn = conn
	args := &loadbalancerrpc.RegisterArgs{ServerInfo: loadbalancerrpc.Node{HostPort: net.JoinHostPort("localhost", strconv.Itoa(port)), NodeID: nodeID}}
	var reply loadbalancerrpc.RegisterReply
	err = pacmanServer.masterConn.Call("LoadBalancer.RegisterServer", args, &reply)
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
	err = rpc.RegisterName("PacmanServer", serverrpc.Wrap(pacmanServer))
	if err != nil {
		return nil, err
	}
	rpc.HandleHTTP()
	go http.Serve(listener, nil)

	pacmanServer.paxos, err = paxos.NewPaxos(net.JoinHostPort("localhost", strconv.Itoa(port)), 0, hostPorts)
	if err != nil {
		return nil, err
	}
	return pacmanServer, nil

}

func (cl *pacmanServer) Temp(args *serverrpc.TempArgs, reply *serverrpc.TempReply) error {
	fmt.Println("HI")
	return nil
}

func (cl *pacmanServer) MakeMove(args *serverrpc.MoveArgs, reply *serverrpc.MoveReply) error {
	direction := args.Direction
	fmt.Println("server got : ", direction)
	err := cl.paxos.RequestValue(direction)
	pack := new(serverrpc.MoveReply)
	pack.Direction = direction
	*reply = *pack
	return nil

}
