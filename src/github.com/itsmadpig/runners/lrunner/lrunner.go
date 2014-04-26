package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/itsmadpig/loadbalancer"
)

const defaultHostPort = "localhost:8009"

var (
	hostPort = flag.String("hostPort", defaultHostPort, "port number to listen on")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()
	if *hostPort == "" {
		// If masterHostPort string is empty, then this storage server is the master.
		*hostPort = defaultHostPort
	}

	// Create and start the StorageServer.
	_, err := loadbalancer.NewLoadBalancer(*hostPort)
	if err != nil {
		log.Fatalln("Failed to create storage server:", err)
	}
	// Run the storage server forever.
	fmt.Println("hostPort=", hostPort)
	select {}
}
