package client

import "github.com/itsmadpig/rpc"

type ClientServer interface {
	//all the mothods
	RegisterClient(*clientrpc.RegisterArgs, *clientrpc.RegisterReply) error
}
