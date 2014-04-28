package loadbalancer

import (
	"errors"
	"fmt"
	"github.com/itsmadpig/rpc/loadbalancerrpc"
	"net"
	"net/http"
	"net/rpc"
)

type loadBalancer struct {
	hostPort          string
	serverHostPorts   []string
	numCurrentNodes   int
	nodes             []loadbalancerrpc.Node
	nodesClientNum    []int
	nodesFailed       []bool
	nodesFailedNumber int

	numOKs int
}

func NewLoadBalancer(hostPort string) (LoadBalancer, error) {
	loadBalancer := new(loadBalancer)
	loadBalancer.hostPort = hostPort
	loadBalancer.numCurrentNodes = 0
	loadBalancer.nodes = make([]loadbalancerrpc.Node, loadbalancerrpc.InitCliNum)
	loadBalancer.nodesClientNum = make([]int, loadbalancerrpc.InitCliNum)
	loadBalancer.nodesFailed = make([]bool, loadbalancerrpc.InitCliNum)
	loadBalancer.nodesFailedNumber = 0
	loadBalancer.numOKs = 0

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
	fmt.Println(lb.nodesFailedNumber)
	pack := new(loadbalancerrpc.RouteReply)
	if lb.nodesFailedNumber > (loadbalancerrpc.InitCliNum / 2) {
		//if most servers already failed
		pack.Status = loadbalancerrpc.MOSTFAIL
		*reply = *pack
		return nil
	}
	if (loadbalancerrpc.InitCliNum == lb.numCurrentNodes) && (lb.numOKs == loadbalancerrpc.InitCliNum) {
		//if most servers not failed and servers ready
		pack.Status = loadbalancerrpc.OK
	} else {
		pack.Status = loadbalancerrpc.NotReady
		*reply = *pack
		return nil
	}

	lowestNum := -1
	lowestIdx := -1
	i := 0
	if args.Attempt == loadbalancerrpc.INIT {
		//first time connecting
		for i = 0; i < loadbalancerrpc.InitCliNum; i++ {
			if (lb.nodesFailed[i] == false) && (lowestNum == -1 || lb.nodesClientNum[i] < lowestNum) {
				//if server hasn't failed and (first lowest or next is smaller)
				lowestNum = lb.nodesClientNum[i]
				lowestIdx = i
			}
		}
		if lowestIdx == -1 {
			return errors.New("shouldn't reach here")
		}
		lb.nodesClientNum[lowestIdx] += 1
		pack.HostPort = lb.nodes[lowestIdx].HostPort
		*reply = *pack
		return nil
	} else if args.Attempt == loadbalancerrpc.RETRY {
		//second time connecting
		for i = 0; i < loadbalancerrpc.InitCliNum; i++ {
			if lb.nodes[i].HostPort == args.HostPort {
				if lb.nodesFailed[i] == false {
					lb.nodesFailed[i] = true
					lb.nodesFailedNumber++
					if lb.nodesFailedNumber > (loadbalancerrpc.InitCliNum / 2) {
						//if most servers already failed
						pack.Status = loadbalancerrpc.MOSTFAIL
						*reply = *pack
						return nil
					}
				}
			} else {
				if (lb.nodesFailed[i] == false) && (lowestNum == -1 || lb.nodesClientNum[i] < lowestNum) {
					//if server hasn't failed and (first lowest or next is smaller)
					lowestNum = lb.nodesClientNum[i]
					lowestIdx = i
				}
			}

		}
		if lowestIdx == -1 {
			return errors.New("shouldn't reach here")
		}
		lb.nodesClientNum[lowestIdx] += 1
		pack.HostPort = lb.nodes[lowestIdx].HostPort
		*reply = *pack
		return nil

	} else {
		pack.Status = loadbalancerrpc.MOSTFAIL
		*reply = *pack
		return errors.New("routeToServer attempt number invalid")
	}

}

func (lb *loadBalancer) RegisterServer(args *loadbalancerrpc.RegisterArgs, reply *loadbalancerrpc.RegisterReply) error {
	fmt.Println("called registerServer")
	pack := new(loadbalancerrpc.RegisterReply)
	add := true
	added := false

	for i := 0; i < lb.numCurrentNodes; i++ {

		if (lb.nodes[i]).NodeID == args.ServerInfo.NodeID {
			add = false
			break
		}
	}
	if add {
		//fmt.Println("adding :", args.ServerInfo.NodeID)
		newNodeList := make([]loadbalancerrpc.Node, loadbalancerrpc.InitCliNum)

		for i := 0; i < lb.numCurrentNodes; i++ {
			if (lb.nodes[i]).NodeID <= args.ServerInfo.NodeID {
				//		////fmt.Println("Here1")
				newNodeList[i] = lb.nodes[i]
			} else {
				//		////fmt.Println("Here2")
				added = true
				newNodeList[i] = args.ServerInfo
				for j := i; j < lb.numCurrentNodes; j++ {
					//		////fmt.Println(j + 1)
					//		////fmt.Println(newNodeList)
					newNodeList[j+1] = lb.nodes[j]
				}

				break

			}
		}
		if !added {
			//		////fmt.Println("newModeList=", newNodeList)
			//		////fmt.Println("numCurrentNodes=", lb.numCurrentNodes)
			newNodeList[lb.numCurrentNodes] = args.ServerInfo
		}
		lb.nodes = newNodeList
		lb.numCurrentNodes += 1
	}

	if lb.numCurrentNodes != loadbalancerrpc.InitCliNum {
		//if not ready
		pack.Status = loadbalancerrpc.NotReady
		pack.Servers = nil
		*reply = *pack
		//fmt.Println("not ready", lb.nodes)
		return errors.New("not ready yet")
	} else {
		//if ready
		pack.Status = loadbalancerrpc.OK
		pack.Servers = lb.nodes
		*reply = *pack
		lb.numOKs += 1
		//fmt.Println("ready pack.servers = ", lb.nodes)
		return nil
	}
}
