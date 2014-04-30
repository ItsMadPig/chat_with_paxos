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
	acceptedProposal int
	value            string
	//
	highestSeenProposal int
	proposalNum         int
	highestRound        int
	currentRound        int
	//
	logs   map[int]string
	stable bool
	//
	ID              int
	myHostPort      string
	serverHostPorts []string
	paxosServers    []*rpc.Client
}

func NewPaxos(myHostPort string, ID int, serverHostPorts []string, test bool) (Paxos, error) {
	thisPaxos := new(paxos)
	thisPaxos.ID = ID
	thisPaxos.currentRound = 0
	thisPaxos.highestSeenProposal = 0
	thisPaxos.proposalNum = 0
	thisPaxos.value = ""
	thisPaxos.myHostPort = myHostPort
	thisPaxos.serverHostPorts = serverHostPorts
	thisPaxos.logs = make(map[int]string)

	if !test {
		fmt.Println("Testing Mode : False")
		err := rpc.RegisterName("Paxos", paxosrpc.Wrap(thisPaxos))
		if err != nil {
			return nil, err
		}

	}
	//dial all other paxos and create a list of them to call.
	err := thisPaxos.DialAllServers()
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
	pack := new(paxosrpc.PrepareReply)
	fmt.Println("pax.highestRound = ", pax.highestRound)
	fmt.Println("received round = ", args.Round)
	fmt.Println("pax.highestSeenProposal = ", pax.highestSeenProposal)
	fmt.Println("received ProposalNumber = ", args.ProposalNumber)

	for i, hP := range pax.serverHostPorts {
		if hP == args.HostPort {
			cli, err := rpc.DialHTTP("tcp", hP)
			if err != nil {
				return err
			}
			pax.paxosServers[i] = cli
		}
	}

	if args.Round <= pax.highestRound {
		pack.HighestAcceptedNum = args.Round
		pack.Value = pax.logs[args.Round]
		pack.Status = paxosrpc.OldInstance

	} else if args.ProposalNumber > pax.highestSeenProposal {
		pax.highestSeenProposal = args.ProposalNumber
		pack.Value = pax.value
		pack.HighestAcceptedNum = pax.acceptedProposal
		pack.Status = paxosrpc.Prepareres
		fmt.Println("Prepareres")
	} else {
		pack.HighestAcceptedNum = pax.highestSeenProposal
		pack.Value = pax.value
		pack.Status = paxosrpc.REJECT
		fmt.Println("Reject")
	}

	*reply = *pack
	return nil
}

func (pax *paxos) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error { //returns the highestSeenProposal
	//takes in a value and an int. Checks if the int is equal to or greater than highestSeenProposal
	//sets value if it is, and returns the min proposal = n.
	fmt.Println("Accept Called")
	pack := new(paxosrpc.AcceptReply)
	pack.HighestSeen = pax.highestSeenProposal

	if args.Round <= pax.highestRound {
		pack.Status = paxosrpc.OldInstance
		pack.Value = pax.logs[args.Round]
		pack.HighestSeen = args.Round
	} else if args.ProposalNumber >= pax.highestSeenProposal {
		pax.highestSeenProposal = args.ProposalNumber /////////////should this be here? doesn't matter
		pax.value = args.Value
		pax.acceptedProposal = args.ProposalNumber
		//don't log yet?

		pack.Status = paxosrpc.OK
	} else {
		pack.Status = paxosrpc.REJECT
	}

	//pax.round += (len(pax.paxosServers) + 1) // do we need this
	*reply = *pack
	return nil
}

func (pax *paxos) Commit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	if pax.highestRound < args.Round {
		pax.logs[args.Round] = args.Value
		pax.highestRound = args.Round
	}
	//pax.proposalNum = pax.proposalNum + 10
	pax.highestSeenProposal = 0
	pax.acceptedProposal = 0
	pax.currentRound = 0
	pax.value = ""
	pax.proposalNum = 0
	pax.stable = true
	fmt.Println(args.Value)
	return nil
}
func (pax *paxos) GetLogs(args *paxosrpc.GetArgs, reply *paxosrpc.GetReply) error {

	pack := new(paxosrpc.GetReply)
	pack.Logs = pax.logs
	*reply = *pack
	return nil
}

func (pax *paxos) RequestValue(reqValue string) error {
	//takes in a string, and acts as a proposer for the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if highestSeenProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.

	pax.stable = false
	majority := ((len(pax.paxosServers) + 1) / 2) + 1

	fmt.Println("request Value called")
	if pax.currentRound == 0 {
		pax.currentRound = pax.highestRound + 10 //////////fix, when restart paxos, shouldn't increment
	}
	if (int(math.Max(float64(pax.highestSeenProposal), float64(pax.proposalNum))) % 10) == 0 {
		pax.proposalNum = int(math.Max(float64(pax.highestSeenProposal), float64(pax.proposalNum))) + 10 + pax.ID
	} else {
		pax.proposalNum = int(math.Max(float64(pax.highestSeenProposal), float64(pax.proposalNum))) + 10
	}

	propArgument := new(paxosrpc.PrepareArgs)
	propArgument.ProposalNumber = pax.proposalNum
	propArgument.Round = pax.currentRound

	propReply := make([]*paxosrpc.PrepareReply, len(pax.paxosServers)+1)
	propChan := make(chan *rpc.Call, len(pax.paxosServers))

	i := 0
	for _, cli := range pax.paxosServers {
		propReply[i] = new(paxosrpc.PrepareReply)
		cli.Go("Paxos.Prepare", propArgument, propReply[i], propChan) //it blocks, if doesn't return for some time, false
		i++
	}

	propReply[i] = new(paxosrpc.PrepareReply)
	pax.Prepare(propArgument, propReply[i]) //not sure if stored in reply
	fmt.Println("value", propReply[i].Value)
	fmt.Println("highestacceptnum", propReply[i].HighestAcceptedNum)
	fmt.Println("status", propReply[i].Status)
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

	count++

	//check what the highest proposal number and highest value is
	propAccepted := 0
	tempHighest := 0
	tempValue := ""
	for _, rep := range propReply {
		//can rep be null?
		if rep.Status != 0 {
			if rep.Status == paxosrpc.OldInstance {
				commitArgs := new(paxosrpc.CommitArgs)
				commitArgs.Value = rep.Value
				commitArgs.Round = rep.HighestAcceptedNum
				pax.Commit(commitArgs, nil)
				fmt.Println("Reach here 1")
				return pax.RequestValue(reqValue)
			} else if rep.Status == paxosrpc.REJECT {
				time.Sleep(3 * time.Second)
				fmt.Println("Reach here 2")
				return pax.RequestValue(reqValue)
			} else {
				propAccepted++
				if rep.HighestAcceptedNum > tempHighest {
					tempValue = rep.Value
					tempHighest = rep.HighestAcceptedNum
				}
			}
		}

	}
	value := ""
	if tempValue != "" {
		value = tempValue
	} else {
		value = reqValue
	}

	if propAccepted < majority {
		pax.stable = true
		return errors.New("CRASH, not enough majority responded") //#nodes smaller than majority
	} else {
		///////////////////////////////////////////
		////////////////////////////accept phase
		///////////////////////////////////////
		//if majority accepted proposal, send accept to all nodes
		k := 0
		acceptArgument := new(paxosrpc.AcceptArgs)
		acceptArgument.Value = value
		acceptArgument.ProposalNumber = pax.proposalNum
		acceptArgument.Round = pax.currentRound
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
			case _ = <-time.After(2 * time.Second):
				break
			}
		}

		acceptAccepted := 0
		for _, rep := range acceptReply {
			//can rep be null?
			if rep.Status != 0 {
				if rep.Status == paxosrpc.OK {
					acceptAccepted++
				} else {
					//increment proposalNum because trying to get other servers to know the rep.highestseen (done at start)
					//return pax.RequestValue(reqValue)
				}
			}
		}
		acceptAccepted++
		if acceptAccepted >= majority {
			commitArg := new(paxosrpc.CommitArgs)
			commitArg.Value = value
			commitArg.Round = pax.currentRound
			for _, cli := range pax.paxosServers {
				acceptReply[k] = new(paxosrpc.AcceptReply)
				cli.Go("Paxos.Commit", commitArg, nil, nil) //it blocks, if doesn't return for some time, false
			}
			pax.Commit(commitArg, nil)
		} else {
			pax.highestSeenProposal = 0
			pax.acceptedProposal = 0
			pax.currentRound = 0
			pax.value = ""
			pax.proposalNum = 0
			pax.stable = true
			return errors.New("ABORT")
			//time.Sleep(3 * time.Second)
			//return pax.RequestValue(reqValue)
		}

	}
	return nil
}
