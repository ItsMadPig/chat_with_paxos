//Implements Paxos package

package paxos

type Paxos struct {
	ID          int
	sequence    int
	minProposal int
	value       string
}

func NewPaxos(master) (Paxos, error) {
	paxos := new(Paxos)
	paxos.minProposal = 0
	value = ""

}

func (pax *Paxos) Prepare(number int) (int, string) {
	//takes in number and checks if number is higher than minProposal
	//if so minProposal = n. returns acceptedProposal number.
	if number > pax.minProposal {
		pax.minProposal = number
	}
	return pax.minProposal, pax.value

}

func (pax *Paxos) Accept(n int, value string) int { //returns the minProposal
	//takes in a value and an int. Checks if the int is equal to or greater than minProposal
	//sets value if it is, and returns the min proposal = n.
	if n >= pax.minProposal {
		pax.minProposal = n
		pax.value = value
	}
	return pax.minProposal

}

func (pax *Paxos) RequestValue(value string) {
	//takes in a string, and initiates the paxos process.
	//send out prepares, wait for majority to reply with same proposal number and empty value
	//if value is not empty, pick the value and proposal number and send commits with it.
	//if value is empty, and mojority replied okay, send out accepts.
	//if minProposal is same as yours, your value is commited and you can return,
	//else start requestValue again.
	proposalNum = pax.sequence + pax.ID
	majority = (numServer / 2) + 1
	//for list of acceptors, call Prepare
	if accepted >= majority {
		//send commits
	}

}
