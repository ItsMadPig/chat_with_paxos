package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/itsmadpig/client"
	"log"
	"os"
)

const defaultPort = 9009
const defaultID = 1

var (
	port = flag.Int("port1", defaultPort, "port number to listen on")

	masterHostPort = flag.String("master", "", "master storage server host port (if non-empty then this storage server is a slave)")
	ID             = flag.Int("id", defaultID, "a unique ID for client")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()
	*masterHostPort = "localhost:8009"

	// Create and start the StorageServer.
	client, err := client.NewPacClient(*masterHostPort, *port, *ID)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		//blocks

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("please retype string")
		}
		message := string([]byte(input))
		err = client.MakeMove(message)
		if err != nil {
			fmt.Println("Server is dead")
			return
		}

	}

}
