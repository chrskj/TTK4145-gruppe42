// -spawne en backup (phoenix)
//
// Holde styr på alle andres ordre-lagre til fil, holde alle orders oppdatert
//
// -be om cost function
// -lagre til fil, holde alle orders oppdatert
// -tildel ordre

//sende ut alle orders den har til komm

//ta imot alle andres ordre fra komm, og vurdere hva som er nye ordre,
//og hva som er ferdige ordre (og da slette)

//Hvor lagres alle de andre heisene sine ordre?
//Når leses det fra fil for å hente gamle ordre?
package orders

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	. "../util"
)

var thisElevator int
var costChan chan ChannelPacket
var data []ChannelPacket
var localOrders [2][]ChannelPacket
var comparing bool = false
var ordersRecieved = false

func InitOrders(OrdersToCom, ComToOrders, ElevAlgoToOrders,
	OrdersToElevAlgo chan ChannelPacket, elevID int) {

	thisElevator = elevID
	readFile()
	data = localOrders[0]
	go orderRoutine(OrdersToCom, ComToOrders, ElevAlgoToOrders, OrdersToElevAlgo)
	time.Sleep(1 * time.Second)
	for _, val := range localOrders[0] {
		if val.Elevator != -1 {
			val.PacketType = "newOrder"
			OrdersToElevAlgo <- val
		}
	}
	for _, val := range localOrders[1] {
		val.PacketType = "cabOrder"
		OrdersToElevAlgo <- val
	}
	OrdersToCom <- ChannelPacket{
		PacketType: "getOrderList",
		Elevator:   thisElevator,
	}
}

func orderRoutine(OrdersToCom, ComToOrders, ElevAlgoToOrders,
	OrdersToElevAlgo chan ChannelPacket) {

	costChan := make(chan ChannelPacket)
	for {
		select {
		case temp := <-ComToOrders:
			switch temp.PacketType {
			case "cost":
				if comparing {
					fmt.Println("before where I think it stops")
					costChan <- temp
					fmt.Println("after where I think it stops")
				}
			case "orderComplete":
				removeOrder(temp)
			case "newOrder":
				addOrder(temp)
			case "getOrderList":
				packet := ChannelPacket{
					PacketType: "orderList",
					OrderList:  data,
				}
				OrdersToCom <- packet
			case "orderList":
				if !ordersRecieved {
					data = temp.OrderList
					var locOrdersTemp []ChannelPacket
					for _, val := range data {
						if val.Elevator == thisElevator {
							locOrdersTemp = append(locOrdersTemp, val)
						}
					}
					localOrders[0] = locOrdersTemp
					ordersRecieved = true
				}
			case "elevLost":
				fmt.Printf("Recieved %s from comm. Redistributing orders", temp.PacketType)
				for len(localOrders[1]) > 0 {
					localOrders[1][0].Elevator = -1
					if costCompare(localOrders[1][0], OrdersToCom, costChan) { //order was assigned successfully
						fmt.Printf("Order was assigned to elevator x\n")
						removeOrder(localOrders[1][0])
					}
				}
			}
		case temp := <-ElevAlgoToOrders:
			switch temp.PacketType {
			case "newOrder":
				addOrder(temp)
			case "buttonPress":
				fmt.Println("Orders recieved " + temp.PacketType +
					" from elevAlgo")
				newOrder := ChannelPacket{
					Elevator:  -1, //Skal det ikke være heisens ID her?
					Floor:     temp.Floor,
					Direction: temp.Direction,
					Timestamp: uint64(time.Now().UnixNano()),
				}
				//check if order already exists
				for _, value := range data {
					if value.Floor == newOrder.Floor &&
						value.Direction == newOrder.Direction {
						newOrder.Timestamp = 0
						break
					}
				}
				//if not: start the cost compare
				if newOrder.Timestamp > 0 {
					go costCompare(newOrder, OrdersToCom, costChan)
				}
			case "engineTimeOut":
				fmt.Println("Motor has stopped. Redistributing orders")
				for len(localOrders[1]) > 0 {
					if costCompare(localOrders[1][0], OrdersToCom, costChan) {
						//order was assigned successfully
						fmt.Printf("Order was assigned to elevator %d\n")
						removeOrder(localOrders[1][0])
					}
				}
			}
		}
	}
}

func costCompare(newOrder ChannelPacket, OrdersToCom,
	costChan chan ChannelPacket) (returnVar bool) {
	comparing = true
	OrdersToCom <- ChannelPacket{
		PacketType: "requestCostFunc",
		Elevator:   thisElevator,
		Floor:      newOrder.Floor,
	}
	//costTicker := time.NewTicker(10 * time.Millisecond)
	tttimer := time.NewTimer(5 * time.Second)
	timein := true
	go func() {
		<-tttimer.C
		timein = false
	}()
	//var ticks uint
	var costs []ChannelPacket
	for recievedOrders := 0; recievedOrders < NumElevators && timein; {
		select {
		case temp := <-costChan:
			fmt.Println("Orders recieved cost")
			unique := true
			if len(costs) > 0 {
				for _, val := range costs {
					if val.Elevator == temp.Elevator {
						unique = false
					}
				}
			}
			if unique {
				costs = append(costs, temp)
				recievedOrders++
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	if !timein {
		fmt.Println("timed out on cost compare")
	}
	max := 9999.0
	fmt.Printf("-+-+-+-+--+-+-++-+-+--+-+-+-++--+-+-+-++--+-+-+-++\n")
	for _, val := range costs {
		fmt.Printf("The cost function of elevator %d is %f\n",
			val.Elevator, val.Cost)
		if val.Cost < max {
			max = val.Cost
			newOrder.Elevator = val.Elevator
		}
	}
	fmt.Printf("-+-+-+-+--+-+-++-+-+--+-+-+-++--+-+-+-++--+-+-+-++\n")
	if newOrder.Elevator != -1 {
		addOrder(newOrder)
		temp := newOrder
		temp.PacketType = "newOrder"
		OrdersToCom <- temp
		returnVar = true
	} else {
		returnVar = false
	}
	comparing = false
	return
}

func readFile() {
	file, err := os.Open(fmt.Sprintf("orders%d.csv", thisElevator))
	if err != nil {
		file, err := os.Create(fmt.Sprintf("orders%d.csv", thisElevator))
		checkError("Cannot create file", err)
		writer := csv.NewWriter(file)
		writer.Flush()
		file.Close()
	} else {
		defer file.Close()

		reader := csv.NewReader(file)
		fmt.Println("before read")
		for {
			input, error := reader.Read()
			if error == io.EOF {
				break
			}
			FloorTemp, _ := strconv.ParseInt(input[0], 10, 64)
			DirectionTemp, _ := strconv.ParseBool(input[1])
			TimestampTemp, _ := strconv.ParseUint(input[2], 10, 64)
			localOrders[0] = append(localOrders[0], ChannelPacket{
				Elevator:  thisElevator,
				Floor:     FloorTemp,
				Direction: DirectionTemp,
				Timestamp: TimestampTemp,
			})
			if len(input) > 3 {
				FloorTemp, _ := strconv.ParseInt(input[3], 10, 64)
				DirectionTemp, _ := strconv.ParseBool(input[4])
				TimestampTemp, _ := strconv.ParseUint(input[5], 10, 64)
				localOrders[1] = append(localOrders[1], ChannelPacket{
					Elevator:  0,
					Floor:     FloorTemp,
					Direction: DirectionTemp,
					Timestamp: TimestampTemp,
				})
			}
		}
	}
}

func writeToFile() {
	fmt.Println("before write")
	if len(localOrders) > 0 {
		file, err := os.Create(fmt.Sprintf("orders%d.csv", thisElevator))
		checkError("Cannot create file", err)
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		fmt.Printf("len0 = %d\n", len(localOrders[0]))
		fmt.Printf("len1 = %d\n", len(localOrders[0]))
		length := len(localOrders[0])
		if len(localOrders[1]) > length {
			length = len(localOrders[1])
		}
		for j := 0; j < length; j++ {
			var valueStr = []string{}
			if j < len(localOrders[0]) {
				valueStr = append(valueStr, []string{strconv.FormatInt(localOrders[0][j].Floor, 10), strconv.FormatBool(localOrders[0][j].Direction)}...)
				valueStr = append(valueStr, strconv.FormatUint(localOrders[0][j].Timestamp, 10))
				fmt.Println(valueStr[0])
			} else {
				valueStr = append(valueStr, []string{"-1", "", ""}...)
			}
			if j < len(localOrders[1]) {
				valueStr = append(valueStr, []string{strconv.FormatInt(localOrders[1][j].Floor, 10), "0"}...)
				valueStr = append(valueStr, strconv.FormatUint(localOrders[1][j].Timestamp, 10))
				fmt.Println(valueStr[0])
			}
			err = writer.Write(valueStr)
			checkError("Cannot write to file", err)
		}

	}
}

func addOrder(newOrder ChannelPacket) {
	fmt.Printf("Adding some orders to this party!\n")
	if newOrder.Elevator != 0 {
		data = append(data, newOrder)
	}
	if newOrder.Elevator == thisElevator {
		localOrders[0] = append(localOrders[0], newOrder)
		writeToFile()
	}
	if newOrder.Elevator == 0 {
		unique := true
		for _, val := range localOrders[1] {
			if newOrder.Elevator == val.Elevator && newOrder.Floor == val.Floor {
				unique = false
			}
		}
		if unique {
			localOrders[1] = append(localOrders[1], newOrder)
			writeToFile()
		}
	}
}

func removeOrder(toRemove ChannelPacket) {
	for index, value := range data { //checks all normal orders
		if value.Floor == toRemove.Floor {
			if len(data) == 1 {
				data = []ChannelPacket{}
			} else if index > 0 { //index-1 >= 0
				data = append(data[:index-1], data[index+1:]...)
			} else {
				data = data[index+1:]
			}
		}
	}
	for i, val := range localOrders {
		for index, value := range val {
			if value.Floor == toRemove.Floor {
				if index > 0 { //index-1 >= 0
					localOrders[i] = append(localOrders[i][:index-1], localOrders[i][index+1:]...)
				} else {
					localOrders[i] = localOrders[i][index+1:]
				}
			}
		}
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
