package orders

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"../util"
)

var thisElevator int
var costChan chan util.ChannelPacket
var data []util.ChannelPacket
var localOrders [2][]util.ChannelPacket
var comparing bool = false
var ordersRecieved = false

func InitOrders(OrdersToCom, ComToOrders, ElevAlgoToOrders,
	OrdersToElevAlgo chan util.ChannelPacket, elevID int) {

	thisElevator = elevID
	readFile()
	data = localOrders[0]
	go orderRoutine(OrdersToCom, ComToOrders, ElevAlgoToOrders, OrdersToElevAlgo)
	time.Sleep(1 * time.Second)
	for _, val := range localOrders[0] {
		if val.Elevator == 1 {
			val.PacketType = "newOrder"
			OrdersToElevAlgo <- val
		}
	}
	for _, val := range localOrders[1] {
		val.PacketType = "cabOrder"
		OrdersToElevAlgo <- val
	}
	OrdersToCom <- util.ChannelPacket{
		PacketType: "getOrderList",
		Elevator:   thisElevator,
		Timestamp:  uint64(time.Now().UnixNano()),
	}
}

func orderRoutine(OrdersToCom, ComToOrders, ElevAlgoToOrders,
	OrdersToElevAlgo chan util.ChannelPacket) {
	costChan := make(chan util.ChannelPacket)
	var redistributeOrders = func(LostElevator int) bool {
		for _, val := range data {
			if val.Elevator == LostElevator {
				val.Timestamp = uint64(time.Now().UnixNano())
				go costCompare(val, OrdersToCom, OrdersToElevAlgo, costChan)
				time.Sleep(3 * time.Second)
			}
		}
		return true
	}

	for {
		select {
		case temp := <-ComToOrders:
			switch temp.PacketType {
			case "cost":
				if comparing {
					costChan <- temp
				}
			case "orderComplete":
				removeOrder(temp)
			case "newOrder":
				fmt.Println("New order from comm")
				if temp.Elevator == thisElevator {
					OrdersToElevAlgo <- temp
					addOrder(temp)
				} else {
					addOrder(temp)
				}
			case "getOrderList":
				packet := util.ChannelPacket{
					PacketType: "orderList",
					OrderList:  data,
					Timestamp:  uint64(time.Now().UnixNano()),
				}
				OrdersToCom <- packet
			case "orderList":
				if !ordersRecieved {
					data = temp.OrderList
					var locOrdersTemp []util.ChannelPacket
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
				redistributeOrders(temp.Elevator)
			}
		case temp := <-ElevAlgoToOrders:
			switch temp.PacketType {
			case "newOrder":
				fmt.Println("New order from elevAlgo")
				addOrder(temp)
			case "buttonPress":
				fmt.Println("Orders recieved " + temp.PacketType + " from elevAlgo")
				newOrder := util.ChannelPacket{
					Elevator:  -1,         //Gets set to the order with the best cost, if it's still -1 at the end,
					Floor:     temp.Floor, //					  that means that there are no available elevators
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
					go costCompare(newOrder, OrdersToCom, OrdersToElevAlgo, costChan)
				}
			case "engineTimeOut":
				fmt.Println("Motor has stopped. Redistributing orders")
				for _, val := range localOrders[0] {
					val.Timestamp = uint64(time.Now().UnixNano())
					go costCompare(val, OrdersToCom, OrdersToElevAlgo, costChan)
					time.Sleep(3 * time.Second)
				}
			}
		}
	}
}

func costCompare(newOrder util.ChannelPacket, OrdersToCom, OrdersToElevAlgo, costChan chan util.ChannelPacket) {
	comparing = true
	OrdersToCom <- util.ChannelPacket{
		PacketType: "requestCostFunc",
		Elevator:   thisElevator,
		Floor:      newOrder.Floor,
		Timestamp:  uint64(time.Now().UnixNano()),
	}
	tttimer := time.NewTimer(5 * time.Second)
	timein := true
	go func() {
		<-tttimer.C
		timein = false
	}()
	var costs []util.ChannelPacket
	for recievedOrders := 0; recievedOrders < util.NumElevators && timein; {
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
	max := 9998.0
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
		temp := newOrder
		temp.PacketType = "newOrder"
		OrdersToCom <- temp
		if newOrder.Elevator == thisElevator {
			OrdersToElevAlgo <- temp
		}
	}
	comparing = false
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
		for {
			input, error := reader.Read()
			if error == io.EOF {
				break
			}
			FloorTemp, _ := strconv.ParseInt(input[0], 10, 64)
			DirectionTemp, _ := strconv.ParseBool(input[1])
			TimestampTemp, _ := strconv.ParseUint(input[2], 10, 64)
			if FloorTemp != -1 {
				localOrders[0] = append(localOrders[0], util.ChannelPacket{
					Elevator:  thisElevator,
					Floor:     FloorTemp,
					Direction: DirectionTemp,
					Timestamp: TimestampTemp,
				})
			}
			if len(input) > 3 {
				FloorTemp, _ := strconv.ParseInt(input[3], 10, 64)
				DirectionTemp, _ := strconv.ParseBool(input[4])
				TimestampTemp, _ := strconv.ParseUint(input[5], 10, 64)
				localOrders[1] = append(localOrders[1], util.ChannelPacket{
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
	if len(localOrders) > 0 {
		file, err := os.Create(fmt.Sprintf("orders%d.csv", thisElevator))
		checkError("Cannot create file", err)
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		length := len(localOrders[0])
		if len(localOrders[1]) > length {
			length = len(localOrders[1])
		}
		for j := 0; j < length; j++ {
			var valueStr = []string{}
			if j < len(localOrders[0]) {
				valueStr = append(valueStr, []string{strconv.FormatInt(localOrders[0][j].Floor, 10), strconv.FormatBool(localOrders[0][j].Direction)}...)
				valueStr = append(valueStr, strconv.FormatUint(localOrders[0][j].Timestamp, 10))
			} else {
				valueStr = append(valueStr, []string{"-1", "", ""}...)
			}
			if j < len(localOrders[1]) {
				valueStr = append(valueStr, []string{strconv.FormatInt(localOrders[1][j].Floor, 10), "0"}...)
				valueStr = append(valueStr, strconv.FormatUint(localOrders[1][j].Timestamp, 10))
			}
			err = writer.Write(valueStr)
			checkError("Cannot write to file", err)
		}

	}
}

func addOrder(newOrder util.ChannelPacket) {
	fmt.Println("Lets add an order!", newOrder)
	if len(data) > 0 {
		if data[len(data)-1].Timestamp != newOrder.Timestamp ||
			newOrder.Elevator != data[len(data)-1].Elevator {
			if newOrder.Elevator != 0 {
				data = append(data, newOrder)
			}
			if newOrder.Elevator == thisElevator {
				localOrders[0] = append(localOrders[0], newOrder)
				writeToFile()
			}
		}
	} else {
		if newOrder.Elevator != 0 {
			data = append(data, newOrder)
		}
		if newOrder.Elevator == thisElevator {
			localOrders[0] = append(localOrders[0], newOrder)
			writeToFile()
		}
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

func removeOrder(toRemove util.ChannelPacket) {
	if len(data) > 0 {
		for index, value := range data { //checks all normal orders
			if value.Floor == toRemove.Floor {
				if len(data) == 1 { //compares the length of data to the index and executes the correct removal
					data = []util.ChannelPacket{}
				} else if index > 0 && index < len(data)-1 {
					data = append(data[:index], data[index+1:]...)
				} else if index == 0 {
					data = data[index+1:]
				} else {
					data = data[:index]
				}
			}
		}
	}
	if len(localOrders) > 0 {
		for i, val := range localOrders {
			for index, value := range val {
				if value.Floor == toRemove.Floor {
					if index > 0 && index < len(localOrders[i])-1 { //compares the length of localOrders to the index and executes the correct removal
						localOrders[i] = append(localOrders[i][:index], localOrders[i][index+1:]...)
						writeToFile()
					} else if index == 0 {
						localOrders[i] = localOrders[i][index+1:]
						writeToFile()
					} else {
						localOrders[i] = localOrders[i][:index]
						writeToFile()
					}
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
