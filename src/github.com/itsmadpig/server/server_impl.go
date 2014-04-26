package server

import (
	//"errors"
	"fmt"
	"github.com/itsmadpig/rpc/serverrpc"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type pacmanServer struct {
	masterServerHostPort string           //master's hostport
	selfNode             *serverrpc.Node  //node of itself
	nodes                []serverrpc.Node // map of all nodes
	masterConn           *rpc.Client      //connection to master
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
	pacmanServer.selfNode = new(serverrpc.Node)
	pacmanServer.nodes = make([]serverrpc.Node, serverrpc.InitCliNum)
	pacmanServer.masterServerHostPort = masterServerHostPort
	if masterServerHostPort == "" {
		//if it's the master client
		pacmanServer.nodes[0] = serverrpc.Node{HostPort: net.JoinHostPort("localhost", strconv.Itoa(port)), NodeID: nodeID}

		listener, err := net.Listen("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)))
		if err != nil {
			return nil, err
		}
		// Wrap the tribServer before registering it for RPC.
		err = rpc.RegisterName("PacmanServer", serverrpc.Wrap(pacmanServer))
		if err != nil {
			return nil, err
		}

		rpc.HandleHTTP()
		go http.Serve(listener, nil)

		return pacmanServer, nil
	} else {
		//if it's the slave client
		fmt.Println("slave server")
		conn, err := rpc.DialHTTP("tcp", masterServerHostPort)
		if err != nil {
			return nil, err
		}
		pacmanServer.masterConn = conn
		args := &serverrpc.RegisterArgs{ServerInfo: serverrpc.Node{HostPort: net.JoinHostPort("localhost", strconv.Itoa(port)), NodeID: nodeID}}
		var reply serverrpc.RegisterReply
		err = pacmanServer.masterConn.Call("PacmanServer.RegisterServer", args, &reply)
		if err != nil {
			return nil, err
		}
		fmt.Println(reply.Status)
		for reply.Status != serverrpc.OK {
			//if master server is still busy
			fmt.Println("retrying to connect")
			time.Sleep(1000 * time.Millisecond)
			err = pacmanServer.masterConn.Call("PacmanServer.RegisterServer", args, &reply)
			if err != nil {
				return nil, err
			}
		}

		//do paxos?
		pacmanServer.nodes = reply.Servers

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

		return pacmanServer, nil

	}

}

func (cl *pacmanServer) RegisterServer(args *serverrpc.RegisterArgs, reply *serverrpc.RegisterReply) error {
	fmt.Println("HI")
	return nil
}
