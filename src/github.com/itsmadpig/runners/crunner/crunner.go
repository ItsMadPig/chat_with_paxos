package main

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	"github.com/itsmadpig/client"
	"log"
	"math"
	"math/big"
	"math/rand"
)

const defaultPort = 9009

var (
	port           = flag.Int("port", defaultPort, "port number to listen on")
	port2          = flag.Int("port2", defaultPort, "port number to listen on")
	masterHostPort = flag.String("master", "", "master storage server host port (if non-empty then this storage server is a slave)")
	nodeID         = flag.Uint("id", 0, "a 32-bit unsigned node ID to use for consistent hashing")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()
	if *masterHostPort == "" && *port == 0 {
		// If masterHostPort string is empty, then this storage server is the master.
		*port = defaultPort
	}
	*masterHostPort = "localhost:8009"

	// If nodeID is 0, then assign a random 32-bit integer instead.
	randID := uint32(*nodeID)
	if randID == 0 {
		randint, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		rand.Seed(randint.Int64())
		randID = rand.Uint32()
	}

	// Create and start the StorageServer.
	client1, err := client.NewPacClient(*masterHostPort, *port)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}
	// Run the storage server forever.
	fmt.Println("masterHostPort=", masterHostPort)
	fmt.Println("port=", port)

	client2, err := client.NewPacClient(*masterHostPort, *port2)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}
	err = client1.MakeMove("up")
	err = client2.MakeMove("down")

}
