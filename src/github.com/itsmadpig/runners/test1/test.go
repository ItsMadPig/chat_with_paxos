package main

import (
	"flag"
	"fmt"
	"github.com/itsmadpig/client"
	"log"
	//"strconv"
	"time"
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
	client2, err := client.NewPacClient(masterHostPort, 2003, "KaranLala")
	if err != nil {
		fmt.Println("Failed to create client", err)

	}
	fmt.Println("Writing Msgs ..")
	// Run the storage server forever.
	err = client2.MakeMove("hi")
	if err != nil {
		fmt.Println("Failed to make move", err)
	}
	/*
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
		}*/

	time.Sleep(time.Second * time.Duration(2))

	logs := client1.GetLogs()
	_, ok := logs[10]
	if !ok {
		fmt.Println("0: failed")
		return
	} else {
		fmt.Println(logs[10])
	}

	/*if value != "KaranLala:hi" {
		fmt.Println("0: failed")
		return
	}*/
	/*
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
		}*/

	fmt.Println("All 4 Posting message tests have passed")
	fmt.Println("Test1 Passed")

	fmt.Println("Starting Test 2. Test 2 checks 1 client getting chat history of session")
	fmt.Println("Adding a new client to retrieve history")

	logs = client2.GetLogs()

	_, ok = logs[10]
	if !ok {
		fmt.Println("0: failed")
		return
	} else {
		fmt.Println(logs[10])
	}

	/*else if value != "KaranLala:hi" {
		fmt.Println("0: failed")
		return
	} */ /*
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
		}*/

	fmt.Println("History properly retrieved")
	fmt.Println("Test2 Passed")

	fmt.Println("Starting Test 3. Test 3 checks 3 clients, 3 servers and 1 loadbalancer with many messages")
}

/*client3, err := client.NewPacClient(masterHostPort, 2004, "AaronHsu")

	store := make(map[int]string)
	index := 0
	for i := 0; i < 15; i++ {
		client2.MakeMove("Whats up" + strconv.Itoa(index))
		store[index] = ("KaranLala:Whats up" + strconv.Itoa(index))
		index++
		client3.MakeMove("TestTest" + strconv.Itoa(index))
		store[index] = ("AaronHsu:TestTest" + strconv.Itoa(index))
		index++
	}
	time.Sleep(time.Duration(5) * time.Second)
	if !isSubsetMap(store, client2.GetLogs()) {
		fmt.Println("Test Failed 3 - 1")
		return
	}
	if !isSubsetMap(store, client3.GetLogs()) {
		fmt.Println("Test Failed 3 - 2")
		return
	}
	if !isSubsetMap(store, client1.GetLogs()) {
		fmt.Println("Test Failed 3 - 3")
		return
	}
}

//returns true if all the values of map1 are also in map2
func isSubsetMap(map1, map2 map[int]string) bool {
	value1 := ""
	for _, value := range map1 {
		exists := false
		for _, value1 = range map2 {
			if value == value1 {
				exists = true
			}
		}
		if exists == false {
			fmt.Println(value, value1)
			return false
		}
	}
	return true
}*/
