package server

import "github.com/itsmadpig/rpc/serverrpc"

type PacmanServer interface {
	//all the mothods
	MakeMove(*serverrpc.MoveArgs, *serverrpc.MoveReply) error
	GetLogs(*serverrpc.GetArgs, *serverrpc.GetReply) error
}
