Write up:

Four components,
1. Client - chat window in command line for users to communicate with other clients
2. Load Balancer - Routes Clients to servers that are least occupied
3. Server - Servers each maintain a replicated log and communicate with each other servers through paxos
4. Paxos object - object imbedded in servers that allow the servers to 

1. failure of Nodes
	

	if client connects later, server pushes all previous logs to client (working)


	if server crashes, all clients connected to that server will be routed by Load Balancer to another working server
	if that crashed server reconnects back, it will be updated with all previous logs automatically, by sending NOPs 
	and resume regular functionality

	if another server with diff ID tries to connect to the server after it already initialized, don't allow.



