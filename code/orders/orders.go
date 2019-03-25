// -spawne en backup (phoenix)
//
// Holde styr på alle andres ordre-lagre til fil, holde alle orders oppdatert
//
// -be om cost function
// -lagre til fil, holde alle orders oppdatert
// -tildel ordre

//sende ut alle orders den har til komm

//ta imot alle andres ordre fra komm, og vurdere hva som er nye ordre, og hva som er ferdige ordre (og da slette)

//comment
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

func InitOrders(OrdersToCom, ComToOrders, OrdersToElevAlgo,
	ElevAlgoToOrders chan ChannelPacket) {
	//try to get data from others

	go orderRoutine(OrdersToCom, ComToOrders, OrdersToElevAlgo, ElevAlgoToOrders)
}

func orderRoutine(OrdersToCom chan ChannelPacket, ComToOrders chan ChannelPacket, OrdersToElevAlgo chan ChannelPacket, ElevAlgoToOrders chan ChannelPacket) {
	var costChan chan ChannelPacket
	for {
		select {
		case temp := <-ComToOrders:
			switch temp.PacketType {
			case "elevID":
				thisElevator = temp.Elevator
			case "cost":
				costChan <- temp
			case "orderComplete":
				removeOrder(ChannelPacket{
					Elevator:  temp.Elevator,
					Floor:     temp.Floor,
					Direction: temp.Direction,
					Timestamp: temp.Timestamp,
				})
			case "newOrder":
				addOrder(ChannelPacket{
					Elevator:  temp.Elevator,
					Floor:     temp.Floor,
					Direction: temp.Direction,
					Timestamp: temp.Timestamp,
				})
			case "getOrderList":
				packet := ChannelPacket{
					PacketType: "orderList",
					OrderList:  data,
				}
				OrdersToCom <- packet
			case "orderList":
				data = temp.OrderList
			}
		case temp := <-ElevAlgoToOrders:
			switch temp.PacketType {
			case "buttonPress":
				newOrder := ChannelPacket{
					Elevator:  -1,
					Floor:     temp.Floor,
					Direction: temp.Direction,
					Timestamp: uint64(time.Now().UnixNano()),
				}
				//check if order already exists
				for _, value := range data {
					if value.Floor == newOrder.Floor && value.Direction == newOrder.Direction {
						newOrder.Timestamp = 0
						break
					}
				}
				//if not: start the cost compare
				if newOrder.Timestamp > 0 {
					go costCompare(newOrder, OrdersToElevAlgo, OrdersToCom)
				}
			}
		}
	}
}

func costCompare(newOrder ChannelPacket, OrdersToElevAlgo, OrdersToCom chan ChannelPacket) {
	OrdersToCom <- ChannelPacket{
		PacketType: "requestCostFunc",
		Elevator:   thisElevator,
	}
	costTicker := time.NewTicker(10 * time.Millisecond)
	var ticks uint = 0
	var costs []ChannelPacket
	for recievedOrders := 0; recievedOrders < NumElevators && ticks < 200; {
		select {
		case temp := <-costChan:
			unique := true
			for _, val := range costs {
				if val.Elevator == temp.Elevator {
					unique = false
				}
			}
			if unique {
				costs = append(costs, temp)
				recievedOrders++
			}
		case <-costTicker.C:
			ticks++
		}
	}
	max := 9999.0
	for _, val := range costs {
		if val.Cost < max {
			max = val.Cost
			newOrder.Elevator = val.Elevator
		}
	}
	if newOrder.Elevator != -1 {
		addOrder(newOrder)
		temp := ChannelPacket{
			PacketType: "newOrder",
			Elevator:   newOrder.Elevator,
			Floor:      newOrder.Floor,
			Direction:  newOrder.Direction,
			Timestamp:  newOrder.Timestamp,
		}
		if temp.Elevator == thisElevator {
			OrdersToElevAlgo <- temp
		} else {
			OrdersToCom <- temp
		}
	} else {
		//error, no costs received
	}
}

func readFile() []ChannelPacket {
	var data []ChannelPacket
	file, err := os.Open("orders.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	reader := csv.NewReader(file)
	fmt.Println("before read")
	for {
		input, error := reader.Read()
		if error == io.EOF {
			break
		}
		for i := 0; i < NumElevators; i++ {
			FloorTemp, _ := strconv.ParseInt(input[0+3*i], 10, 64)
			DirectionTemp, _ := strconv.ParseBool(input[1+3*i])
			tstampTemp, _ := strconv.ParseUint(input[2+3*i], 10, 64)
			ElevatorTemp := i + 1
			data = append(data, ChannelPacket{
				Elevator:  ElevatorTemp,
				Floor:     FloorTemp,
				Direction: DirectionTemp,
				Timestamp: tstampTemp,
			})
		}
		FloorTemp, _ := strconv.ParseInt(input[3*NumElevators], 10, 64)
		tstampTemp, _ := strconv.ParseUint(input[3*NumElevators+1], 10, 64)
		data = append(data, ChannelPacket{
			Elevator:  0,
			Floor:     FloorTemp,
			Timestamp: tstampTemp,
		})
	}
	return data
}

func writeToFile() {
	
	} 
	/*
		err = writer.Write(writeData)
		checkError("Cannot write to file", err)
	*/


func addOrder(newOrder ChannelPacket) {
		data = append(data, newOrder)
		for i := 0; i < NumElevators; i++ {
}
}

func removeOrder(toRemove ChannelPacket) []ChannelPacket {

	return data
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}