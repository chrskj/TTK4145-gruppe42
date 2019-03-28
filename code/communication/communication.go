package communication

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

import (
	"fmt"
	"strconv"

	"../network/bcast"
	"../network/peers"
	. "../util"
)

func InitCom(toElevAlgo, toOrders, fromElevAlgo, fromOrders chan ChannelPacket,
	elevID int) {

	sendMessage := make(chan ChannelPacket)
	go bcast.Transmitter(16570, sendMessage)

	receiveMessage := make(chan ChannelPacket)
	go bcast.Receiver(16570, receiveMessage)

	peerTxEnable := make(chan bool)
	go peers.Transmitter(16569, strconv.Itoa(elevID), peerTxEnable)

	peerUpdateCh := make(chan peers.PeerUpdate)
	go peers.Receiver(16569, peerUpdateCh)

	for {
		select {
		case msg := <-fromElevAlgo:
			//fmt.Printf("Comm Recieved packet of type %s from ElevAlgo\n", msg.PacketType)
			msg.Elevator = elevID
			sendMessage <- msg
		case msg := <-fromOrders:
			//fmt.Printf("Comm Recieved packet of type %s from Orders\n", msg.PacketType)
			switch msg.PacketType {
			case "requestCostFunc":
				sendMessage <- msg
			case "getOrderList":
				sendMessage <- msg
			case "newOrder":
				//fmt.Printf("newOrder.Elevator = %d\n", msg.Elevator)
				if msg.Elevator == elevID {
					toElevAlgo <- msg
					sendMessage <- msg
				} else {
					sendMessage <- msg
				}
			case "orderList":
				sendMessage <- msg
			}
		case msg := <-receiveMessage:
			//fmt.Printf("Comm Recieved packet of type %s from broadcast\n", msg.PacketType)
			switch msg.PacketType {
			case "newOrder":
				if msg.Elevator == elevID {
					toOrders <- msg
				} else if msg.Elevator != elevID {
					toElevAlgo <- msg
				}
			case "orderList":
				if msg.Elevator == elevID {
					toOrders <- msg
				}
			case "getOrderList":
				toOrders <- msg
			case "cost":
				toOrders <- msg
			case "orderComplete":
				toOrders <- msg
				toElevAlgo <- msg
			case "requestCostFunc":
				toElevAlgo <- msg
			}
		case msg := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", msg.Peers)
			fmt.Printf("  New:      %q\n", msg.New)
			fmt.Printf("  Lost:     %q\n", msg.Lost)
			if len(msg.Lost) > 0 {
				idLost, _ := strconv.Atoi(msg.Lost[0])
				toOrders <- ChannelPacket{
					PacketType: "elevLost",
					Elevator:   idLost,
				}
			}
		default:
		}
	}
}
