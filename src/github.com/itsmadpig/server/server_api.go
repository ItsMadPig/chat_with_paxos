package server

import "github.com/itsmadpig/rpc/serverrpc"

type PacmanServer interface {
	//all the mothods
	Temp(*serverrpc.TempArgs, *serverrpc.TempReply) error
	MakeMove(*serverrpc.MoveArgs, *serverrpc.MoveReply) error
}
