Write up:

Four components,
1. Client - chat window in command line for users to communicate with other clients
2. Load Balancer - Routes Clients to servers that are least occupied
3. Server - Servers each maintain a replicated log and communicate with each other servers through paxos
4. Paxos object - object imbedded in servers that allow the servers to 


number of servers is stored in InitCliNum = 3 in loadbalancerrpc/proto.go

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

		/////////////////////////////////////////////////////////////////////////////////
		test for disconnecting server with client on, and client reroutes, paxos works
		test for disconnecting server, with client not on, paxos works
		test for disconnecting server, and paxos works
		test for disconnecting server and reconnecting it back, does NOP and gets previous values
			then connect a client to that server and paxos works


		./lrunner
		./srunner =master="localhost:8009" -port=9010 -id=1
		./srunner =master="localhost:8009" -port=9011 -id=2
		./srunner =master="localhost:8009" -port=9012 -id=3
		./crunner -id=Aaron -port1=2002
		type something 
		type something
		./crunner -id=Lala -port1=2003
		type something
		kill 1st server, then client reroutes to server 3, because least load
		type something, paxos working
		then restart server 1, it loads all previous logs,
		now connect a new client, and new client downloads all logs
		/////////////////////////////////////////////////////////////////////////////////







2. Network temporarily disabled causes package delay or package Dropped
	a) Proposer delays: 
						If it's incoming proposes/accepts delays for more than willing to wait , then it got rejected
						
						During that time when its acceptor receives propose or accept, then check

	b) Acceptor delays: Temporarily disable connection of server for more than the proposer is willing to wait, if proposer got majority and already moved on (at accept or commit phase)
	   other servers that are already ahead in round number will ignore it. That server will get updated on its log if it proposes something later.


	   //////////////////////////////////////////
		test 1
		two clients 3 servers
		server id=1 shouldn't get commited until 30 seconds
		////////////////////////////////////////////
		./lrunner
		./srunner -master="localhost:8009" -port=9010 -id=1 -mode=p:10::a:20::c:30
		./srunner -master="localhost:8009" -port=9011 -id=2
		./srunner -master="localhost:8009" -port=9012 -id=3
		./test1


		/////////////////////////////////////////
		test 2
		two clients 3 servers
		server id=1 should get commited after 10 seconds
		//////////////////////////////////////////
		./lrunner
		./srunner -master="localhost:8009" -port=9010 -id=1 -mode=p:10::a:20::c:10
		./srunner -master="localhost:8009" -port=9011 -id=2
		./srunner -master="localhost:8009" -port=9012 -id=3
		./test1


		/////////////////////////////////////////
		test 3
		two clients 3 servers
		server id
		/////////////////////////////////////////
		./lrunner
		./srunner -master="localhost:8009" -port=9010 -id=1 -mode="d:0"
		./srunner -master="localhost:8009" -port=9012 -id=2
		./srunner -master="localhost:8009" -port=9011 -id=3 -mode="d:0"
		./test1




3. State 
	All the logs are replicated through paxos for correctness










