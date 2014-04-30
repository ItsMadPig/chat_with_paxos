Write up:

Four components,
1. Client - chat window in command line for users to communicate with other clients
2. Load Balancer - Routes Clients to servers that are least occupied
3. Server - Servers each maintain a replicated log and communicate with each other servers through paxos
4. Paxos object - object imbedded in servers that allow the servers to 


Everything assumes that the Load balancer does not fail and is always running (Load Balancer Router)
First, load balancer waits for all servers to connect. Once all are connected, service to client begins
If another server with different ID tries to connect to the server after it already initialized, don't allow.

All servers keep a log. The log is a map of key:round to value: string that was commited
Every commited round is unique.

1. failure of Nodes / Node recovery
	Let there be 2f+1 servers

	If more than f servers failed (>=f+1), then everything stops working, have to reset
	else if just a few servers crashes, we will be able to:
		a) Automatic client reroute
			  If server crashes, all clients connected to that server will be routed by Load Balancer to another working server.
		b) Automatic server update all previous logs through paxos
			  If that crashed server reconnects back, it will be updated with all previous logs automatically.
				Since the server just started up, its highestRound will be 0,
				The server does this by multiple iterations of RequestValue that proposes a Nop (empty string) and its currentRound
				when other servers see an old round that was already commited, they will send back the already commited value, so
				and it will continue updating its logs until its highestRound is not commited (not stored in the log).
		Resume regular functionality


2. Network temporarily disabled causes package delay or package Dropped
	a) Proposer delays: 
						If it's incoming proposes/accepts delays for more than willing to wait , then it got rejected
						
						During that time when its acceptor receives propose or accept, then check

	b) Acceptor delays: Temporarily disable connection of server for more than the proposer is willing to wait, if proposer got majority and already moved on (at accept or commit phase)
	   other servers that are already ahead in round number will ignore it. That server will get updated on its log if it proposes something later.

3. State 
	All the logs are replicated through paxos for correctness





