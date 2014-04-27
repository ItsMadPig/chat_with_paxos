//Implements Paxos package

package paxos

import (
	"fmt"
	"github.com/itsmadpig/rpc/paxosrpc"
	"net/rpc"
)

type paxos struct {
	ID                   int
	sequence             int
	minProposal          int
	value                string
	masterServerHostPort string
	myHostPort           string
	serverHostPorts      []string
	paxosServers         []*rpc.Client
}

func NewPaxos(masterServerHostPort, myHostPort string, ID int, serverHostPorts []string) (Paxos, error) {
	paxos := new(paxos)
	paxos.minProposal = 0
	paxos.value = ""
	paxos.ID = ID
	paxos.masterServerHostPort = masterServerHostPort
	paxos.myHostPort = myHostPort
	paxos.serverHostPorts = serverHostPorts
	err = rpc.RegisterName("Paxos", paxosrpc.Wrap(paxos))
	if err != nil {
		return nil, err
	}

	//dial all other paxos and create a list of them to call.
	err = DialAllServers()
	if err != nil {
		return nil, err
	}

	return paxos, nil
}

func (pax *paxos) DialAllServers() error {
	pax.paxosServers = make([]*rpc.Client, len(pax.serverHostPorts)-1)
	i := 0
	for server := range pax.serverHostPorts {
		if server != pax.myHostPort {
			cli, err := rpc.DialHTTP("tcp", server)
			if err != nil {
				return err
			}
			pax.paxosServers[i] = cli
			i++
		}
	}
	return nil

}

func (pax *paxos) Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) {
	//takes in number and checks if number is higher than minProposal
	//if so minProposal = n. returns acceptedProposal number.

	number = args.ProposalNumber
	if number > pax.minProposal {
		pax.minProposal = number
	}
	var pack = &paxosrpc.PrepareReply{Value: pax.value, ProposalAwaiting: pax.minProposal}
	*reply = *pack

}

func (pax *paxos) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) { //returns the minProposal
	//takes in a value and an int. Checks if the int is equal to or greater than minProposal
	//sets value if it is, and returns the min proposal = n.
	number = args.ProposalNumber
	if number >= pax.minProposal {
		pax.minProposal = number
		pax.value = args.Value
	}
	pack = &paxosrpc.AcceptReply{ProposalCommited: pax.minProposal}
	*reply = *pack

}

func (pax *paxos) RequestValue(args *paxosrpc.RequestArgs) {
	//takes in a string, and acts as a proposer for the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if minProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.
	value = args.Value
	proposalNum = pax.sequence + pax.ID
	majority = (len(pax.paxosServers) / 2) + 1
	args := &paxosrpc.PrepareArgs{ProposalNumber: proposalNum}
	reply := make([]*paxosrpc.PrepareReply)
	i := 0
	for cli := range pax.paxosServers {
		cli.Call("Paxos.Prepare", args, &reply[i])
		i++
	}
	accepted := 0
	totalReplied := 0
	for rep := range reply {
		if rep.ProposalAwaiting == proposalNum { //check this algo
			accepted++
		} else if rep.ProposalAwaiting >= ProposalNum
	}

	//for list of acceptors, call Prepare
	if accepted >= majority {
		//send commits
	}

}
