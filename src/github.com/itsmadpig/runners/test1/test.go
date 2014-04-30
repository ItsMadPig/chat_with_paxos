package main

import (
	"flag"
	"fmt"
	"github.com/itsmadpig/client"
	"log"
	//"time"
)

const defaultPort = 9009
const masterHostPort = "localhost:8009"

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()
	fmt.Println("Starting Test 1. Test 1 checks 1 client 1 loadbalancer 3 servers")

	// If nodeID is 0, then assign a random 32-bit integer instead.

	fmt.Println("starting client")
	client1, err := client.NewPacClient(masterHostPort, 2002, "Karan")
	if err != nil {
		fmt.Println("Failed to create client", err)

	}

	fmt.Println("Writing Msgs ..")
	// Run the storage server forever.
	err = client1.MakeMove("hi")
	if err != nil {
		fmt.Println("Failed to make move", err)
	}
	err = client1.MakeMove("how are you?")
	if err != nil {
		fmt.Println("Failed to make move", err)
	}
	err = client1.MakeMove("this is a test")
	if err != nil {
		fmt.Println("Failed to make move", err)
	}
	err = client1.MakeMove("Now to see if all messages are stored.")
	if err != nil {
		fmt.Println("Failed to make move", err)
	}

	logs := client1.GetLogs()
	for index, value := range logs {
		if index == 0 {
			if value != "Karan:hi" {
				fmt.Println("0: failed")
			}
		}
		if index == 1 {
			if value != "Karan:how are you?" {
				fmt.Println("1: failed")
			}
		}
		if index == 2 {
			if value != "Karan:this is a test" {
				fmt.Println("2: failed")
			}
		}
		if index == 3 {
			if value != "Karan:Now to see if all messages are stored." {
				fmt.Println("3: failed")
			}
		}

	}
	fmt.Println("All 4 tests have passed")

}
