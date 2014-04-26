package server

import "github.com/itsmadpig/rpc/serverrpc"

type PacmanServer interface {
	//all the mothods
	RegisterServer(*serverrpc.RegisterArgs, *serverrpc.RegisterReply) error
}
