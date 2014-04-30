package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/itsmadpig/server"
)

const defaultPort = 9009

var (
	port           = flag.Int("port", defaultPort, "port number to listen on")
	masterHostPort = flag.String("master", "", "master storage server host port (if non-empty then this storage server is a slave)")
	nodeID         = flag.Int("id", 0, "a unique integer")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()
	if *masterHostPort == "" {
		// If masterHostPort string is empty, then this storage server is the master.
		*masterHostPort = "localhost:8009"
	}
	if *port == 0 {
		*port = defaultPort
	}

	// If nodeID is 0, then assign a random 32-bit integer instead.

	// Create and start the StorageServer.
	flags := []string{""}
	_, err := server.NewServer(*masterHostPort, *port, *nodeID, true, flags)
	if err != nil {
		log.Fatalln("Failed to create storage server:", err)
	}
	// Run the storage server forever.
	fmt.Println("masterHostPort=", masterHostPort)
	fmt.Println("port=", port)
	fmt.Println("nodeID=", nodeID)
	select {}
}
