//Implements Paxos package

package paxos

import (
	"fmt"
	"github.com/itsmadpig/rpc/paxosrpc"
	"net/rpc"
	"time"
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

func (pax *paxos) Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	//takes in number and checks if number is higher than minProposal
	//if so minProposal = n. returns acceptedProposal number.

	number = args.ProposalNumber
	pack := new(paxosrpc.PrepareReply)
	if number >= pax.minProposal { /////////////////////////check
		pax.minProposal = number
		pack.Value = pax.value
		pack.ProposalAwaiting = pax.minProposal
		pack.Status = paxosrpc.OK
	} else {

		pack.Status = paxosrpc.REJECT
	}
	*reply = *pack
	return nil
}

func (pax *paxos) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) { //returns the minProposal
	//takes in a value and an int. Checks if the int is equal to or greater than minProposal
	//sets value if it is, and returns the min proposal = n.
	number = args.ProposalNumber
	if args.Value == pax.value && number >= pax.minProposal {

	}
}

func (pax *paxos) RequestValue(args *paxosrpc.RequestArgs) error {
	//takes in a string, and acts as a proposer for the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if minProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.
	value := args.Value
	proposalNum := pax.sequence + pax.ID
	majority := ((len(pax.paxosServers) + 1) / 2) + 1
	argument := &paxosrpc.PrepareArgs{ProposalNumber: proposalNum}
	reply := make([]*paxosrpc.PrepareReply, len(pax.paxosServers)+1)

	i := 0
	for cli := range pax.paxosServers {
		cli.Call("Paxos.Prepare", argument, &reply[i]) //it blocks, if doesn't return for some time, false
		i++
	}
	pax.Prepare(argument, &reply[i]) //not sure if stored in reply
	accepted := 0
	totalReplied := 0

	for rep := range reply {
		if rep.Status == paxosrpc.OK { //check this algo
			accepted++
			if rep.Value != "" { ///////////////////
				value = rep.Value
			}
		}
	}
	//for list of acceptors, call Prepare
	if accepted >= majority {
		//send commits
		k := 0
		j := 0
		acceptArgument := new(paxosrpc.AcceptArgs)
		acceptArgument.Value = value
		acceptArgument.ProposalNumber = proposalNum
		acceptReply := make([]*paxosrpc.PrepareReply, accepted)

		for cli := range pax.paxosServers {
			if reply[k].Status == paxosrpc.OK {
				cli.Call("Paxos.Accept", argument, &acceptReply[j]) //it blocks, if doesn't return for some time, false
				j++
			}
			k++
		}
	} else {
		time.Sleep(1000 * time.Millisecond)
		pax.RequestValue(args)
	}

}
