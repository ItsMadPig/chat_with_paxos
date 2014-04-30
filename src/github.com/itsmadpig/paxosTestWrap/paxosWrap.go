package paxosWrap

import (
	"github.com/itsmadpig/paxos"
	"github.com/itsmadpig/rpc/paxosrpc"
	"github.com/itsmadpig/rpc/paxoswraprpc"
	"net/rpc"
)

type paxosWrap struct {
	paxos paxos.Paxos
	flags []string
}

func NewPaxosWrap(myHostPort string, ID int, serverHostPorts []string, flags []string) (PaxosWrap, error) {
	Wrapper := new(paxosWrap)
	paxos, err := paxos.NewPaxos(myHostPort, ID, serverHostPorts, true)
	if err != nil {
		return nil, err
	}
	Wrapper.paxos = paxos
	err = rpc.RegisterName("Paxos", paxoswraprpc.Wrap(Wrapper))
	if err != nil {
		return nil, err
	}

	return Wrapper, nil

}

func (pax *paxosWrap) Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	//takes in number and checks if number is higher than highestSeenProposal
	//if so highestSeenProposal = n. returns acceptedProposal number.

	return pax.paxos.Prepare(args, reply)
}

func (pax *paxosWrap) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error { //returns the highestSeenProposal
	//takes in a value and an int. Checks if the int is equal to or greater than highestSeenProposal
	//sets value if it is, and returns the min proposal = n.

	return pax.paxos.Accept(args, reply)
}

func (pax *paxosWrap) Commit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {

	return pax.paxos.Commit(args, reply)
}

func (pax *paxosWrap) GetLogs(args *paxosrpc.GetArgs, reply *paxosrpc.GetReply) error {
	return pax.paxos.GetLogs(args, reply)
}

func (pax *paxosWrap) RequestValue(reqValue string) error {
	//takes in a string, and acts as a proposer for the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if highestSeenProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.

	return pax.paxos.RequestValue(reqValue)

}
