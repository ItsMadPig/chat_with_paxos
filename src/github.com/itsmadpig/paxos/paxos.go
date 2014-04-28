//Implements Paxos package

package paxos

import (
	"fmt"
	"github.com/itsmadpig/rpc/paxosrpc"
	"net/rpc"
	"time"
)

type paxos struct {
	ID                  int
	round               int
	highestSeenProposal int
	proposalNum         int
	//
	acceptedProposal int
	value            string

	//
	myHostPort      string
	serverHostPorts []string
	logs            map[int]string
	paxosServers    []*rpc.Client
}

func NewPaxos(myHostPort string, ID int, serverHostPorts []string) (Paxos, error) {
	paxos := new(paxos)
	paxos.ID = ID
	paxos.round = 0
	paxos.highestSeenProposal = 0
	paxos.proposalNum = pax.ID
	paxos.value = ""
	paxos.myHostPort = myHostPort
	paxos.serverHostPorts = serverHostPorts
	paxos.logs = make(map[int]string)
	err = rpc.RegisterName("Paxos", paxosrpc.Wrap(paxos))
	if err != nil {
		return nil, err
	}

	//dial all other paxos and create a list of them to call.
	err = paxos.DialAllServers()
	if err != nil {
		return nil, err
	}

	return paxos, nil
}

func (pax *paxos) commitLog(r int, s string) error {
	//takes in r round number and s string
	pax.logs[r] = s
}
func (pax *paxos) dialAllServers() error {
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
	//takes in number and checks if number is higher than highestSeenProposal
	//if so highestSeenProposal = n. returns acceptedProposal number.

	number = args.ProposalNumber
	pack := new(paxosrpc.PrepareReply)
	if number >= pax.highestSeenProposal { /////////////////////////check
		pax.highestSeenProposal = number
		pack.Value = pax.value
		pack.HighestAcceptedNum = pax.acceptedProposal
		pack.Status = paxosrpc.OK
	} else {

		pack.Status = paxosrpc.REJECT
	}
	*reply = *pack
	return nil
}

func (pax *paxos) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) { //returns the highestSeenProposal
	//takes in a value and an int. Checks if the int is equal to or greater than highestSeenProposal
	//sets value if it is, and returns the min proposal = n.
	number := args.ProposalNumber
	pack := new(paxosrpc.AcceptReply)
	pack.HighestSeen = pax.highestSeenProposal
	if number >= pax.highestSeenProposal {
		pax.highestSeenProposal = number
		pax.value = args.Value
		pax.acceptedProposal = number
		pack.Status = paxosrpc.OK
	} else {
		pack.Status = paxosrpc.REJECT
	}

	//pax.round += (len(pax.paxosServers) + 1) // do we need this
	*reply = *pack
}

func (pax *paxos) Commit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) {

}

func (pax *paxos) RequestValue(args *paxosrpc.RequestArgs) error {
	//takes in a string, and acts as a proposer for the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if highestSeenProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.

	proposalNum := pax.proposalNum + 10 //
	pax.proposalNum = max(pax.highestSeenProposal, pax.proposalNum)
	majority := ((len(pax.paxosServers) + 1) / 2) + 1

	propArgument := new(paxosrpc.PrepareArgs)
	propArgument.ProposalNumer = proposalNum

	propReply := make([]*paxosrpc.PrepareReply, len(pax.paxosServers)+1)
	propChan := make(chan *rpc.Call, len(pax.paxosServers))

	i := 0
	for cli := range pax.paxosServers {
		propReply[i] = new(paxosrpc.PrepareReply)
		cli.Go("Paxos.Prepare", propArgument, &propReply[i], propChan) //it blocks, if doesn't return for some time, false
		i++
	}

	//fix this
	count := 0
	for count < i {
		select {
		case _ = <-propChan.Done:
			count++
			if count >= i {
				break
			}
		case _ = <-time.After(2 * time.Second): //how does this work?
			break
		}
	}

	/*
		propChanList := make([]*rpc.Call, len(pax.paxosServers))
		for cli := range pax.paxosServers {
			propChanList[i] = new(rpc.Call)
			propChanList[i] = cli.Go("Paxos.Prepare", propArgument, &propReply[i], nil) //it blocks, if doesn't return for some time, false
			i++
		}

		count := 0
		for count < i {
			select {
			case _ = <-propChanList[count].Done:
				count++
				if count >= i {
					break
				}
			case _ = <-time.After(3 * time.Second): //how does this work?
				break
			}
		}
	*/

	pax.Prepare(argument, &propReply[i]) //not sure if stored in reply
	count++

	//quick check if not majority, restart
	if count < majority {
		time.Sleep(3 * time.Second)
		return pax.RequestValue(args)
	}

	//check what the highest proposal number and highest value is
	propAccepted := 0
	tempHighest := 0
	tempValue := ""
	for rep := range propReply {
		//can rep be null?
		if rep.Status == paxosrpc.OK { //check this algo
			propAccepted++
			if rep.HighestAcceptedNum > tempHighest { /////////////////// && empty string?
				tempValue = rep.Value
				tempHighest = rep.HighestAcceptedNum
			}
		}
	}

	if tempValue == "" {
		value := args.Value
	} else {
		value = tempValue
	}

	if propAccepted < majority {
		time.Sleep(3 * time.Second)
		return pax.RequestValue(args)

	} else {
		///////////////////////////////////////////
		////////////////////////////accept phase
		///////////////////////////////////////
		//if majority accepted proposal, send accept to all nodes
		k := 0
		acceptArgument := new(paxosrpc.AcceptArgs)
		acceptArgument.Value = value
		acceptArgument.ProposalNumber = proposalNum
		acceptReply := make([]*paxosrpc.AcceptReply, len(pax.paxosServers)+1)
		acceptChan := make(chan *rpc.Call, len(pax.paxosServers))

		for cli := range pax.paxosServers {
			acceptReply[k] = new(paxosrpc.AcceptReply)
			cli.GO("Paxos.Accept", acceptArgument, &acceptReply[k], acceptChan) //it blocks, if doesn't return for some time, false
			k++
		}

		acceptCount := 0
		for acceptCount < k {
			select {
			case _ = <-acceptChan.Done:
				acceptCount++
				if acceptCount >= k {
					break
				}
			case _ = <-time.After(3 * time.Second):
				break
			}
		}
		pax.Accept()
		acceptCount++

		acceptAccepted := 0
		for rep := range acceptReply {
			//can rep be null?
			if rep.HighestSeen > proposalNum {
				//increment proposalNum because trying to get other servers to know the rep.highestseen (done at start)
				time.Sleep(3 * time.Second)
				return pax.RequestValue(args)
			}
		}
		if acceptCount >= majority {

		} else {
			time.Sleep(3 * time.Second)
			return pax.RequestValue(args)
		}

	}
}
