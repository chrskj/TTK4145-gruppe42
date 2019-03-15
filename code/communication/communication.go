package communication

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

import (
	//"flag"
	"fmt"
	//"os"
	"time"
	//"math/rand"
    //"strconv"
    "github.com/chrskj/TTK4145-gruppe44/code/network/bcast"
    "github.com/chrskj/TTK4145-gruppe44/code/network/peers"
    //. "github.com/chrskj/TTK4145-gruppe44/code/util"
)

type MessageStruct struct {
	Message string
	Iter    int
}

func SendHeartbeat(id string) {
	peerTxEnable := make(chan bool)
	go peers.Transmitter(16569, id, peerTxEnable)
}

func ReceiveHeartbeat() {
	peerUpdateCh := make(chan peers.PeerUpdate)
	go peers.Receiver(16569, peerUpdateCh)
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
		}
    }
}

func SendMessage(id string) {
	sendMessage := make(chan MessageStruct)
	go bcast.Transmitter(16570, sendMessage)

    helloMsg := MessageStruct{"Hello from " + id, 0}
    for {
        helloMsg.Iter++
        sendMessage <- helloMsg
        time.Sleep(1 * time.Second)
    }
}

func ReceiveMessage() {
    receiveMessage := make(chan MessageStruct)
    go bcast.Receiver(16570, receiveMessage)
    for {
        select {
        case a := <-receiveMessage:
            fmt.Printf("Received: %v\n", a)
        }
    }
}

/*
func ListenForModules(fromElevAlgo, fromOrder, toElevAlgo, toOrder) {
    for {
        select {
        case <-fromElevAlgo:
            fmt.Printf("elevAlgo")
        case <-fromOrders:
            fmt.Printf("orders")
        default:
            fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
        }
    }
}
*/

func ListenForModules(fromElevAlgo chan int) {
    for {
        select {
        case temp := <-fromElevAlgo:
            fmt.Println(temp)
        default:
            fmt.Println("    .")
			time.Sleep(1000 * time.Millisecond)
        }
    }
}
