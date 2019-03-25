package communication

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

import (
	//"flag"
	"fmt"
	"os"
	"time"

	//"math/rand"
	"strconv"

	"../network/bcast"
	"../network/peers"
	. "../util"
)

func InitCom(toElevAlgo, toOrders, fromElevAlgo, fromOrders chan ChannelPacket) {
	id := os.Getpid()

	sendMessage := make(chan ChannelPacket)
	go bcast.Transmitter(16570, sendMessage)

	receiveMessage := make(chan ChannelPacket)
	go bcast.Receiver(16570, receiveMessage)

	go SendHeartbeat(strconv.Itoa(id))
	go ReceiveHeartbeat()

	idPacket := ChannelPacket{
		PacketType: "elevID",
		Elevator:   id,
	}

	toOrders <-idPacket

	for {
		select {
		case temp := <-fromElevAlgo:
			fmt.Printf("Received packet of type %s from elevAlgo:\n", temp.PacketType)
			// Skal begge meldinger sendes over nettet? (cost & ordersComplete)
            sendMessage <-temp
			toOrders <-temp
		case temp := <-fromOrders:
			fmt.Println(temp.PacketType)
			switch temp.PacketType {
			case "requestCostFunction":
                sendMessage <-temp
				toElevAlgo <-temp
			case "getOrderList":
				// Hva må gjøres her?
                sendMessage <-temp
			case "newOrder":
				// Hva må gjøres her?
                sendMessage <-temp
			case "orderList":
				// Hva må gjøres her?
                sendMessage <-temp
			}
		case temp := <-receiveMessage:
			fmt.Printf("Received packet of type %s from the net:\n", temp.PacketType)
			switch temp.PacketType {
			case "newOrder":
				toOrders <-temp
			case "orderList":
				toOrders <-temp
			case "getOrderList":
				toOrders <-temp
			case "cost":
				toOrders <-temp
			case "orderComplete":
				toOrders <-temp
			case "requestCostFunction":
				toElevAlgo <-temp
			}
		default:
			fmt.Println("    .")
			time.Sleep(time.Second)
		}
	}
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
