package communication

import (
	"fmt"
	"strconv"
	"time"

	"../network/bcast"
	"../network/peers"
	. "../util"
)

func RedundantBcast(msg ChannelPacket, sendMessage chan ChannelPacket) {
	for tries := 0; tries < 5; tries++ {
		sendMessage <- msg
		time.Sleep(100 * time.Millisecond)
	}
}

func InitCom(ComToElevAlgo, ComToOrders, ElevAlgoToCom,
	OrdersToCom chan ChannelPacket, elevID int) {

	sendMessage := make(chan ChannelPacket)
	go bcast.Transmitter(16570, sendMessage)

	receiveMessage := make(chan ChannelPacket)
	go bcast.Receiver(16570, receiveMessage)

	peerTxEnable := make(chan bool)
	go peers.Transmitter(16569, strconv.Itoa(elevID), peerTxEnable)

	peerUpdateCh := make(chan peers.PeerUpdate)
	go peers.Receiver(16569, peerUpdateCh)

	lastRecieved := []uint64{0, 0, 0, 0, 0, 0}

	for {
		select {
		case msg := <-ElevAlgoToCom:
			go RedundantBcast(msg, sendMessage)
		case msg := <-OrdersToCom:
			go RedundantBcast(msg, sendMessage)
		case msg := <-receiveMessage:
			switch msg.PacketType {
			case "newOrder":
				if lastRecieved[0] != msg.Timestamp {
					lastRecieved[0] = msg.Timestamp
					ComToOrders <- msg
					if msg.Elevator != elevID {
						ComToElevAlgo <- msg
					}
				}
			case "orderList":
				if lastRecieved[1] != msg.Timestamp {
					lastRecieved[1] = msg.Timestamp
					if msg.Elevator == elevID {
						ComToOrders <- msg
					}
				}
			case "getOrderList":
				if lastRecieved[2] != msg.Timestamp {
					lastRecieved[2] = msg.Timestamp
					ComToOrders <- msg

				}
			case "cost":
				if lastRecieved[3] != msg.Timestamp {
					lastRecieved[3] = msg.Timestamp
					ComToOrders <- msg
				}

			case "orderComplete":
				if lastRecieved[4] != msg.Timestamp {
					lastRecieved[4] = msg.Timestamp
					ComToOrders <- msg
					ComToElevAlgo <- msg
				}

			case "requestCostFunc":
				if lastRecieved[5] != msg.Timestamp {
					lastRecieved[5] = msg.Timestamp
					ComToElevAlgo <- msg
				}
			}
		case msg := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", msg.Peers)
			fmt.Printf("  New:      %q\n", msg.New)
			fmt.Printf("  Lost:     %q\n", msg.Lost)
			if len(msg.Lost) > 0 {
				idLost, _ := strconv.Atoi(msg.Lost[0])
				ComToOrders <- ChannelPacket{
					PacketType: "elevLost",
					Elevator:   idLost,
				}
			}
		default:
		}
	}
}
