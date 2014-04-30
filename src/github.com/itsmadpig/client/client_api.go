package client

//import "github.com/itsmadpig/rpc/clientrpc"

type PacClient interface {
	MakeMove(string) error
	GetLogs() map[int]string
}
