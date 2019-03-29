package communication

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"

	"../network/bcast"
	"../network/conn"
	. "../util"
	w "../watchdog"
)

func InitCom(toElevAlgo, toOrders, fromElevAlgo, fromOrders chan ChannelPacket,
	elevID int) {

	sendMessage := make(chan ChannelPacket)
	go bcast.Transmitter(16570, sendMessage)

	receiveMessage := make(chan ChannelPacket)
	go bcast.Receiver(16570, receiveMessage)

	go SendHeartbeat(16569, strconv.Itoa(elevID))

	go ReceiveHeartbeat(16569, toOrders)

	/*
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
	*/

	//handShakeChan := make(chan ChannelPacket)

	lastRecieved := []uint64{0, 0, 0, 0, 0, 0}

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
				if lastRecieved[0] != msg.Timestamp {
					lastRecieved[0] = msg.Timestamp
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
				if lastRecieved[1] != msg.Timestamp {
					lastRecieved[1] = msg.Timestamp
					//start
					if msg.Elevator == elevID {
						toOrders <- msg
					}
					//end
				}
			case "getOrderList":
				if lastRecieved[2] != msg.Timestamp {
					lastRecieved[2] = msg.Timestamp
					//start
					toOrders <- msg
					//end
				}
			case "cost":
				if lastRecieved[3] != msg.Timestamp {
					lastRecieved[3] = msg.Timestamp
					//start
					toOrders <- msg
					//end
				}

			case "orderComplete":
				if lastRecieved[4] != msg.Timestamp {
					lastRecieved[4] = msg.Timestamp
					//start
					toOrders <- msg
					toElevAlgo <- msg
					//end
				}

			case "requestCostFunc":
				if lastRecieved[5] != msg.Timestamp {
					lastRecieved[5] = msg.Timestamp
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
	for tries := 0; tries < 10; tries++ {
		sendMessage <- msg
		time.Sleep(100 * time.Millisecond)
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

func SendHeartbeat(port int, id string) {
	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
	for {
		time.Sleep(60 * time.Millisecond)
		conn.WriteTo([]byte(id), addr)
	}
}

func ReceiveHeartbeat(port int, toOrders chan ChannelPacket) {
	var beatList []string
	var dogList []*w.Watchdog
	var buf [1048576]byte
	conn := conn.DialBroadcastUDP(port)
	for {
		n, _, _ := conn.ReadFrom(buf[0:])
		id := string(buf[:n])

		unique := true
		var index int

		for i, val := range beatList {
			if id == val {
				unique = false
				index = i
			}
		}
		if unique {
			beatList = append(beatList, id)
			dogList = append(dogList, w.New(3*time.Second))
		} else {
			dogList[index].Reset()
		}

		cases := make([]reflect.SelectCase, len(dogList))
		for i, _ := range dogList {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(dogList[i].TimeOverChannel())}
		}
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})

		chosen, _, _ := reflect.Select(cases)

		if chosen != len(dogList) {
			lostID, _ := strconv.Atoi(beatList[chosen])
			toOrders <- ChannelPacket{
				PacketType: "elevLost",
				Elevator:   lostID,
			}
		}
		//fmt.Println(id)
	}
}
