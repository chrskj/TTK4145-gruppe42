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

	handShakeChan := make(chan ChannelPacket)
	lastRecieved := ChannelPacket{ //Dr.Frankenstein's FrankenPacket
		orderList: []ChannelPacket{ChannelPacket{PacketType: "newOrder",Floor: 0,},
			ChannelPacket{PacketType: "orderList", Floor: 1,},
			ChannelPacket{PacketType: "getOrderList",Floor: 2,},
			ChannelPacket{PacketType: "cost",Floor: 3,},
			ChannelPacket{PacketType: "orderComplete",Floor: 4,},
			ChannelPacket{PacketType: "requestCostFunc",Floor: 5,}}
	}

	for {
		select {
		case msg := <-fromElevAlgo:
			msg.Elevator = elevID
			SendImportantMsg(msg, sendMessage, handShakeChan)
		case msg := <-fromOrders:
			SendImportantMsg(msg, sendMessage, handShakeChan)
		case msg := <-receiveMessage:
			//fmt.Printf("Comm Recieved packet of type %s from broadcast\n", msg.PacketType)
			switch msg.PacketType {
			case "newOrder":
				if(lastRecieved.OrderList[0].TimeStamp != msg.TimeStamp){
					lastRecieved.OrderList[0].TimeStamp = msg.TimeStamp
					//start
					toOrders <- msg
					if msg.Elevator != elevID {
						toElevAlgo <- msg
					}
					//end
				}
				msg.PacketType = "handShake"
				msg.Elevator = elevID
				sendMessage <- msg
			case "orderList":
				if(lastRecieved.OrderList[1].TimeStamp != msg.TimeStamp){
					lastRecieved.OrderList[1].TimeStamp = msg.TimeStamp
					//start
					if msg.Elevator == elevID {
						toOrders <- msg
					}
					//end
				}
				msg.PacketType = "handShake"
				msg.Elevator = elevID
				sendMessage <- msg
			case "getOrderList":
				if(lastRecieved.OrderList[2].TimeStamp != msg.TimeStamp){
					lastRecieved.OrderList[2].TimeStamp = msg.TimeStamp
					//start
					toOrders <- msg
					//end
				}
				msg.PacketType = "handShake"
				msg.Elevator = elevID
				sendMessage <- msg
			case "cost":
				if(lastRecieved.OrderList[3].TimeStamp != msg.TimeStamp){
					lastRecieved.OrderList[3].TimeStamp = msg.TimeStamp
					//start
					toOrders <- msg
					//end
				}
				msg.PacketType = "handShake"
				msg.Elevator = elevID
				sendMessage <- msg
			case "orderComplete":
				if(lastRecieved.OrderList[4].TimeStamp != msg.TimeStamp){
					lastRecieved.OrderList[4].TimeStamp = msg.TimeStamp
					//start
					toOrders <- msg
					toElevAlgo <- msg
					//end
				}
				msg.PacketType = "handShake"
				msg.Elevator = elevID
				sendMessage <- msg
			case "requestCostFunc":
				if(lastRecieved.OrderList[5].TimeStamp != msg.TimeStamp){
					lastRecieved.OrderList[5].TimeStamp = msg.TimeStamp
					//start
					toElevAlgo <- msg
					//end
				}
				msg.PacketType = "handShake"
				msg.Elevator = elevID
				sendMessage <- msg
			case "handShake":
				handShakeChan <- msg
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

func SendImportantMsg(msg ChannelPacket, sendMessage, handShakeChan chan ChannelPacket){
	recievedHandShakes := []int{}
	sendMessage <- msg
	for tries := 0; tries<10 && len(recievedHandShakes)<NumElevators{
		select{
		case temp <- handShakeChan:
			if temp.TimeStamp == msg.TimeStamp{
				unique := true
				for _, val in recievedHandShakes{
					if temp.Elevator == val.Elevator{unique = false}
				}
				if unique {
					recievedHandShakes = append(recievedHandShakes, temp)
				}
			} else {
				handShakeChan <- temp
				tries++
				time.Sleep(300*time.Millisecond)
			}
		default:
			sendMessage <- msg
			tries++
			time.Sleep(300*time.Millisecond)
		}
	}
}