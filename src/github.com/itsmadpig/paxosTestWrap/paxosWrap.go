package paxosWrap

import (
	"github.com/itsmadpig/paxos"
	"github.com/itsmadpig/rpc/paxosrpc"
	"github.com/itsmadpig/rpc/paxoswraprpc"
	"net/rpc"
	"strconv"
	"strings"
	"time"
)

type paxosWrap struct {
	paxos    paxos.Paxos
	flags    []string
	timeout  time.Duration
	rTimeout time.Duration
	pTimeout time.Duration
	aTimeout time.Duration
	gTimeout time.Duration
	cTimeout time.Duration
	rDelay   bool
	pDelay   bool
	aDelay   bool
	gDelay   bool
	cDelay   bool
	sleep    bool
	die      bool
}

func NewPaxosWrap(myHostPort string, ID int, serverHostPorts []string, flags []string) (PaxosWrap, error) {
	Wrapper := new(paxosWrap)
	paxos, err := paxos.NewPaxos(myHostPort, ID, serverHostPorts, true)
	Wrapper.timeout = 0
	Wrapper.die = false
	Wrapper.sleep = false
	if err != nil {
		return nil, err
	}
	Wrapper.paxos = paxos
	err = rpc.RegisterName("Paxos", paxoswraprpc.Wrap(Wrapper))
	if err != nil {
		return nil, err
	}
	for _, value := range flags {
		split := strings.Split(value, ":")
		if split[0] == "s" {
			Wrapper.sleep = true
			timer, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, err
			}
			Wrapper.timeout = time.Duration(timer)
		} else if split[0] == "d" {
			Wrapper.die = true
		} else if split[0] == "p" {
			Wrapper.pDelay = true
			timer, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, err
			}
			Wrapper.pTimeout = time.Duration(timer)
		} else if split[0] == "a" {
			Wrapper.aDelay = true
			timer, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, err
			}
			Wrapper.aTimeout = time.Duration(timer)
		} else if split[0] == "g" {
			Wrapper.gDelay = true
			timer, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, err
			}
			Wrapper.gTimeout = time.Duration(timer)
		} else if split[0] == "c" {
			Wrapper.cDelay = true
			timer, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, err
			}
			Wrapper.cTimeout = time.Duration(timer)
		} else if split[0] == "r" {
			Wrapper.rDelay = true
			timer, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, err
			}
			Wrapper.rTimeout = time.Duration(timer)
		}
	}

	return Wrapper, nil

}

func (pax *paxosWrap) Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	//takes in number and checks if number is higher than highestSeenProposal
	//if so highestSeenProposal = n. returns acceptedProposal number.
	if pax.sleep {
		time.Sleep(time.Second * pax.timeout)
	}
	if pax.die {
		return nil
	}
	if pax.pDelay {
		time.Sleep(time.Second * pax.pTimeout)
	}
	return pax.paxos.Prepare(args, reply)
}

func (pax *paxosWrap) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error { //returns the highestSeenProposal
	//takes in a value and an int. Checks if the int is equal to or greater than highestSeenProposal
	//sets value if it is, and returns the min proposal = n.
	if pax.sleep {
		time.Sleep(time.Second * pax.timeout)
	}
	if pax.die {
		return nil
	}
	if pax.aDelay {
		time.Sleep(time.Second * pax.aTimeout)
	}
	return pax.paxos.Accept(args, reply)
}

func (pax *paxosWrap) Commit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	if pax.sleep {
		time.Sleep(time.Second * pax.timeout)
	}
	if pax.die {
		return nil
	}
	if pax.cDelay {
		time.Sleep(time.Second * pax.cTimeout)
	}
	return pax.paxos.Commit(args, reply)
}

func (pax *paxosWrap) GetLogs(args *paxosrpc.GetArgs, reply *paxosrpc.GetReply) error {
	if pax.sleep {
		time.Sleep(time.Second * pax.timeout)
	}
	if pax.die {
		return nil
	}
	if pax.gDelay {
		time.Sleep(time.Second * pax.gTimeout)
	}
	return pax.paxos.GetLogs(args, reply)
}

func (pax *paxosWrap) RequestValue(reqValue string) error {
	//takes in a string, and acts as a proposer for the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if highestSeenProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.
	if pax.sleep {
		time.Sleep(time.Second * pax.timeout)
	}
	if pax.die {
		return nil
	}
	if pax.rDelay {
		time.Sleep(time.Second * pax.rTimeout)
	}
	return pax.paxos.RequestValue(reqValue)

}
