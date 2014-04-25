package client

import (
	//"errors"
	"fmt"
	"github.com/itsmadpig/rpc"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type clientServer struct {
	masterClientHostPort string           //master's hostport
	selfNode             *clientrpc.Node  //node of itself
	nodes                []clientrpc.Node // map of all nodes
	masterConn           *rpc.Client      //connection to master
}

//treat each node as both a client and server

//two types of disconnection
//explicit disconnect (user stops game)
//implicit disconnect (user crashes)

//first do static paxos, then implement dynamic paxos (quorum dynamic)

//static paxos: wait till all connected, then only runs paxos if master client disconnects
//dynamic paxos: running paxos even before all nodes connect.

//if masterClientHostPort isn't empty, then it's a masterclient,
//else it's a slaveclient to start with
func NewClient(masterClientHostPort string, port int, nodeID uint32) (ClientServer, error) {
	clientServer := new(clientServer)
	clientServer.selfNode = new(clientrpc.Node)
	clientServer.nodes = make([]clientrpc.Node, clientrpc.InitCliNum)
	clientServer.masterClientHostPort = masterClientHostPort
	if masterClientHostPort == "" {
		//if it's the master client
		clientServer.nodes[0] = clientrpc.Node{HostPort: net.JoinHostPort("localhost", strconv.Itoa(port)), NodeID: nodeID}

		listener, err := net.Listen("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)))
		if err != nil {
			return nil, err
		}
		// Wrap the tribServer before registering it for RPC.
		err = rpc.RegisterName("ClientServer", clientrpc.Wrap(clientServer))
		if err != nil {
			return nil, err
		}

		rpc.HandleHTTP()
		go http.Serve(listener, nil)

		return clientServer, nil
	} else {
		//if it's the slave client
		fmt.Println("slave client")
		conn, err := rpc.DialHTTP("tcp", masterClientHostPort)
		if err != nil {
			return nil, err
		}
		clientServer.masterConn = conn
		args := &clientrpc.RegisterArgs{ServerInfo: clientrpc.Node{HostPort: net.JoinHostPort("localhost", strconv.Itoa(port)), NodeID: nodeID}}
		var reply clientrpc.RegisterReply
		err = clientServer.masterConn.Call("ClientServer.RegisterClient", args, &reply)
		if err != nil {
			fmt.Println("HEREHEREHERE")
			return nil, err
		}
		fmt.Println(reply.Status)
		for reply.Status != clientrpc.OK {
			//if master server is still busy
			fmt.Println("retrying to connect")
			time.Sleep(1000 * time.Millisecond)
			err = clientServer.masterConn.Call("ClientServer.RegisterClient", args, &reply)
			if err != nil {
				fmt.Println("HEREHERE")
				return nil, err
			}
		}

		//do paxos?
		clientServer.nodes = reply.Servers

		listener, err := net.Listen("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)))
		if err != nil {
			return nil, err
		}
		err = rpc.RegisterName("ClientServer", clientrpc.Wrap(clientServer))
		if err != nil {
			return nil, err
		}
		rpc.HandleHTTP()
		go http.Serve(listener, nil)

		return clientServer, nil

	}

}

func (cl *clientServer) RegisterClient(args *clientrpc.RegisterArgs, reply *clientrpc.RegisterReply) error {
	fmt.Println("HI")
	return nil
}
