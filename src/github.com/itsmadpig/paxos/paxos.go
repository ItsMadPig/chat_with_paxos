//Implements Paxos package

package paxos

import (
	"errors"
	"fmt"
	"github.com/itsmadpig/rpc/paxosrpc"
	"math"
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
	thisPaxos := new(paxos)
	thisPaxos.ID = ID
	thisPaxos.round = 0
	thisPaxos.highestSeenProposal = 0
	thisPaxos.proposalNum = ID
	thisPaxos.value = ""
	thisPaxos.myHostPort = myHostPort
	thisPaxos.serverHostPorts = serverHostPorts
	thisPaxos.logs = make(map[int]string)

	err := rpc.RegisterName("Paxos", paxosrpc.Wrap(thisPaxos))
	if err != nil {
		return nil, err
	}
	//dial all other paxos and create a list of them to call.
	err = thisPaxos.DialAllServers()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		err = thisPaxos.DialAllServers()
		if err == nil {
			break
		}
		if i == 4 {
			return nil, errors.New("dial all servers error")
		}
	}

	return thisPaxos, nil
}

func (pax *paxos) commitLog(r int, s string) error {
	//takes in r round number and s string
	pax.logs[r] = s
	return nil
}

func (pax *paxos) DialAllServers() error {
	pax.paxosServers = make([]*rpc.Client, len(pax.serverHostPorts)-1)
	i := 0
	for _, server := range pax.serverHostPorts {
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
	fmt.Println("prepare Called")
	number := args.ProposalNumber
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

func (pax *paxos) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error { //returns the highestSeenProposal
	//takes in a value and an int. Checks if the int is equal to or greater than highestSeenProposal
	//sets value if it is, and returns the min proposal = n.
	fmt.Println("Accept Called")
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
	return nil
}

func (pax *paxos) Commit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	fmt.Println("Commit Called")
	if pax.value != args.Value {
		fmt.Println("wrong commit")
		return nil
	}
	pax.logs[pax.acceptedProposal] = pax.value
	pax.proposalNum = pax.proposalNum + 10
	fmt.Println(args.Value)
	return nil
}

func (pax *paxos) GetLogs(args *paxosrpc.GetArgs, reply *paxosrpc.GetReply) error {

	pack := new(paxosrpc.GetReply)
	pack.Logs = pax.logs
	*reply = *pack
	return nil
}

func (pax *paxos) RequestValue(direction string) error {
	//takes in a string, and acts as a proposer for the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if highestSeenProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.
	fmt.Println("request Value called")
	proposalNum := pax.proposalNum + 10 //
	fmt.Println("proposal Num : ", proposalNum)
	pax.proposalNum = int(math.Max(float64(pax.highestSeenProposal), float64(pax.proposalNum)))
	majority := ((len(pax.paxosServers) + 1) / 2) + 1

	propArgument := new(paxosrpc.PrepareArgs)
	propArgument.ProposalNumber = proposalNum

	propReply := make([]*paxosrpc.PrepareReply, len(pax.paxosServers)+1)
	propChan := make(chan *rpc.Call, len(pax.paxosServers))

	i := 0
	for _, cli := range pax.paxosServers {
		propReply[i] = new(paxosrpc.PrepareReply)
		cli.Go("Paxos.Prepare", propArgument, propReply[i], propChan) //it blocks, if doesn't return for some time, false
		i++
	}

	//fix this
	count := 0
	for count < i {
		select {
		case _ = <-propChan:
			count++
			if count >= i {
				break
			}
		case _ = <-time.After(2 * time.Second): //how does this work?
			break
		}
	}

	propReply[i] = new(paxosrpc.PrepareReply)
	pax.Prepare(propArgument, propReply[i]) //not sure if stored in reply
	count++

	//quick check if not majority, restart
	if count < majority {
		time.Sleep(3 * time.Second)
		return pax.RequestValue(direction)
	}

	//check what the highest proposal number and highest value is
	propAccepted := 0
	tempHighest := 0
	tempValue := ""
	for _, rep := range propReply {
		//can rep be null?
		if rep != nil {
			if rep.Status == paxosrpc.OK { //check this algo
				propAccepted++
				if rep.HighestAcceptedNum > tempHighest { /////////////////// && empty string?
					tempValue = rep.Value
					tempHighest = rep.HighestAcceptedNum
				}
			}
		}

	}
	value := ""
	if tempValue != "" && tempHighest >= proposalNum {
		value = tempValue
	} else {
		value = direction
	}

	if propAccepted < majority {
		time.Sleep(3 * time.Second)
		return pax.RequestValue(direction)

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

		for _, cli := range pax.paxosServers {
			acceptReply[k] = new(paxosrpc.AcceptReply)
			cli.Go("Paxos.Accept", acceptArgument, &acceptReply[k], acceptChan) //it blocks, if doesn't return for some time, false
			k++
		}
		acceptReply[k] = new(paxosrpc.AcceptReply)
		pax.Accept(acceptArgument, acceptReply[k])
		acceptCount := 0
		for acceptCount < k {
			select {
			case _ = <-acceptChan:
				acceptCount++
				if acceptCount >= k {
					break
				}
			case _ = <-time.After(3 * time.Second):
				break
			}
		}
		//
		acceptAccepted := 0
		for _, rep := range acceptReply {
			//can rep be null?
			if rep != nil {
				if rep.Status == paxosrpc.OK {
					acceptAccepted++
				} else {
					//increment proposalNum because trying to get other servers to know the rep.highestseen (done at start)
					time.Sleep(3 * time.Second)
					return pax.RequestValue(direction)
				}
			}

		}
		if acceptAccepted >= majority {
			commitArg := new(paxosrpc.CommitArgs)
			commitArg.Value = value
			for _, cli := range pax.paxosServers {
				acceptReply[k] = new(paxosrpc.AcceptReply)
				cli.Go("Paxos.Commit", commitArg, nil, nil) //it blocks, if doesn't return for some time, false
			}
			pax.Commit(commitArg, nil)
		} else {
			time.Sleep(3 * time.Second)
			return pax.RequestValue(direction)
		}

	}
	return nil
}
