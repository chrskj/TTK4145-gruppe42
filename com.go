// Network module for go

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

package main

import (
	"./network/bcast"
	//"./network/localip"
	//"./network/peers"
	//"flag"
	"fmt"
	//"os"
	"time"	
	"math/rand"	
    //"strconv"
)

type MessageStruct struct {
	Message string
}

func main() {
	fmt.Println("Started")

    // Send Heartbeat
    rand.Seed(time.Now().UnixNano())
    ranInt := rand.Intn(20)
	transmitMessage := make(chan MessageStruct)
	go bcast.Transmitter(16569, transmitMessage)

	go func() {
        response := fmt.Sprintf("Heartbeat from %d", ranInt)
		helloMsg := MessageStruct{response}
		for {
			transmitMessage <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

    // Listen to heartbeat
	receiveMessage := make(chan MessageStruct)
	go bcast.Receiver(16569, receiveMessage)

	for {
		select {
		case a := <-receiveMessage:
			fmt.Printf("Received: %v\n", a)
		}
	}
}

//func UDP_init()


//func sendMessage(adresse, beskjed)


//func recieveMessage(addresse, beskjed)
