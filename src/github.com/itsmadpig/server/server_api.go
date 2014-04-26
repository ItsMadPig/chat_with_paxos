package server

import "github.com/itsmadpig/rpc"

type PacmanServer interface {
	//all the mothods
	RegisterServer(*serverrpc.RegisterArgs, *serverrpc.RegisterReply) error
}
