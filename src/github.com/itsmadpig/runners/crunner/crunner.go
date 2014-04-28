package main

import (
	"flag"
	"github.com/itsmadpig/client"
	"log"
)

const defaultPort = 9009

var (
	port1 = flag.Int("port1", defaultPort, "port number to listen on")
	port2 = flag.Int("port2", defaultPort, "port number to listen on")
	port3 = flag.Int("port3", defaultPort, "port number to listen on")
	port4 = flag.Int("port4", defaultPort, "port number to listen on")

	masterHostPort = flag.String("master", "", "master storage server host port (if non-empty then this storage server is a slave)")
	nodeID         = flag.Uint("id", 0, "a 32-bit unsigned node ID to use for consistent hashing")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()
	*masterHostPort = "localhost:8009"

	// Create and start the StorageServer.
	client1, err := client.NewPacClient(*masterHostPort, *port1)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}
	// Run the storage server forever.

	client2, err := client.NewPacClient(*masterHostPort, *port2)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}
	client3, err := client.NewPacClient(*masterHostPort, *port3)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}
	client4, err := client.NewPacClient(*masterHostPort, *port4)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}
	err = client1.MakeMove("up")
	err = client2.MakeMove("down")
	err = client3.MakeMove("left")
	err = client4.MakeMove("right")

}
