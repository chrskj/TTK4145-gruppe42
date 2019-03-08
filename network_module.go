// Network module for go
package main

import (
	"./network/bcast"
	//"./network/localip"
	//"./network/peers"
	//"flag"
	"fmt"
	//"os"
	"time"	
)

type MessageStruct struct {
	Message string
}

func main() {
	transmitMessage := make(chan MessageStruct)
	receiveMessage := make(chan MessageStruct)

	go bcast.Transmitter(16569, transmitMessage)
	go bcast.Receiver(16569, receiveMessage)

	go func() {
		helloMsg := MessageStruct{"Hello from me"}
		for {
			receiveMessage <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case a := <-receiveMessage:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}

//func UDP_init()


//func sendMessage(adresse, beskjed)


//func recieveMessage(addresse, beskjed)