package communication

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

import (
	"fmt"
	"os"
	"time"
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

    peerTxEnable := make(chan bool)
	go peers.Transmitter(16569, strconv.Itoa(id), peerTxEnable)

	peerUpdateCh := make(chan peers.PeerUpdate)
	go peers.Receiver(16569, peerUpdateCh)

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
		case temp := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", temp.Peers)
			fmt.Printf("  New:      %q\n", temp.New)
			fmt.Printf("  Lost:     %q\n", temp.Lost)
        default:
			fmt.Println("    .")
			time.Sleep(time.Second)
		}
	}
}
