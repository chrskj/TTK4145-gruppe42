package communication

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

import (
	"fmt"
	"strconv"
	"time"

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

	go func() {
		msg := <-peerUpdateCh
		for {
			msg = <-peerUpdateCh
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
		}
	}()

	//handShakeChan := make(chan ChannelPacket)

	lastRecieved := ChannelPacket{ //Dr.Frankenstein's FrankenPacket
		OrderList: []ChannelPacket{
			ChannelPacket{PacketType: "newOrder", Floor: 0},
			ChannelPacket{PacketType: "orderList", Floor: 1},
			ChannelPacket{PacketType: "getOrderList", Floor: 2},
			ChannelPacket{PacketType: "cost", Floor: 3},
			ChannelPacket{PacketType: "orderComplete", Floor: 4},
			ChannelPacket{PacketType: "requestCostFunc", Floor: 5},
		},
	}

	for {
		select {
		case msg := <-fromElevAlgo:
			msg.Elevator = elevID
			go SendImportantMsg(msg, sendMessage)
		case msg := <-fromOrders:
			go SendImportantMsg(msg, sendMessage)
		case msg := <-receiveMessage:
			//fmt.Printf("Comm Recieved packet of type %s from broadcast\n", msg.PacketType)
			switch msg.PacketType {
			case "newOrder":
				if lastRecieved.OrderList[0].Timestamp != msg.Timestamp {
					lastRecieved.OrderList[0].Timestamp = msg.Timestamp
					//start
					toOrders <- msg
					if msg.Elevator != elevID {
						toElevAlgo <- msg
					}
					//end
				}
				//msg.PacketType = "handShake"
				//msg.Elevator = elevID
				//sendMessage <- msg
			case "orderList":
				if lastRecieved.OrderList[1].Timestamp != msg.Timestamp {
					lastRecieved.OrderList[1].Timestamp = msg.Timestamp
					//start
					if msg.Elevator == elevID {
						toOrders <- msg
					}
					//end
				}
			case "getOrderList":
				if lastRecieved.OrderList[2].Timestamp != msg.Timestamp {
					lastRecieved.OrderList[2].Timestamp = msg.Timestamp
					//start
					toOrders <- msg
					//end
				}
			case "cost":
				if lastRecieved.OrderList[3].Timestamp != msg.Timestamp {
					lastRecieved.OrderList[3].Timestamp = msg.Timestamp
					//start
					toOrders <- msg
					//end
				}

			case "orderComplete":
				if lastRecieved.OrderList[4].Timestamp != msg.Timestamp {
					lastRecieved.OrderList[4].Timestamp = msg.Timestamp
					//start
					toOrders <- msg
					toElevAlgo <- msg
					//end
				}

			case "requestCostFunc":
				if lastRecieved.OrderList[5].Timestamp != msg.Timestamp {
					lastRecieved.OrderList[5].Timestamp = msg.Timestamp
					//start
					toElevAlgo <- msg
					//end
				}
				//case "handShake":
				//	handShakeChan <- msg
			}

		default:
		}
	}
}

func SendImportantMsg(msg ChannelPacket, sendMessage chan ChannelPacket) {
	for tries := 0; tries < 1; tries++ {
		sendMessage <- msg
		time.Sleep(10 * time.Millisecond)
	}
}

/*
func SendImportantMsg(msg ChannelPacket, sendMessage, handShakeChan chan ChannelPacket){
	recievedHandShakes := []int{}
	sendMessage <- msg
	for tries := 0; tries < 10 && len(recievedHandShakes) < NumElevators; {
		select {
		case temp := <-handShakeChan:
			if temp.Timestamp == msg.Timestamp {
				unique := true
				for _, val := range recievedHandShakes {
					if temp.Elevator == val {
						unique = false
					}
				}
				if unique {
					recievedHandShakes = append(recievedHandShakes, temp.Elevator)
				}
			} else {
				handShakeChan <- temp
				tries++
				time.Sleep(300 * time.Millisecond)
			}
		default:
			sendMessage <- msg
			tries++
			time.Sleep(300 * time.Millisecond)
		}
	}
}
*/
