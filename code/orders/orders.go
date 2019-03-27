// -spawne en backup (phoenix)
//
// Holde styr på alle andres ordre-lagre til fil, holde alle orders oppdatert
//
// -be om cost function
// -lagre til fil, holde alle orders oppdatert
// -tildel ordre

//sende ut alle orders den har til komm

//ta imot alle andres ordre fra komm, og vurdere hva som er nye ordre, og hva som er ferdige ordre (og da slette)

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

func InitOrders(OrdersToCom, ComToOrders, OrdersToElevAlgo,
	ElevAlgoToOrders chan ChannelPacket) {
	//try to get data from others

	go orderRoutine(OrdersToCom, ComToOrders, OrdersToElevAlgo, ElevAlgoToOrders)
}

func orderRoutine(OrdersToCom chan ChannelPacket, ComToOrders chan ChannelPacket, OrdersToElevAlgo chan ChannelPacket, ElevAlgoToOrders chan ChannelPacket) {
	costChan := make(chan ChannelPacket)
	for {
		select {
		case temp := <-ComToOrders:
			switch temp.PacketType {
			case "elevID":
				fmt.Println("Recieved Elevator ID")
				thisElevator = temp.Elevator
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
				data = temp.OrderList
			}
		case temp := <-ElevAlgoToOrders:
			switch temp.PacketType {
			case "buttonPress":
				fmt.Println("Orders recieved " + temp.PacketType + " from elevAlgo")
				newOrder := ChannelPacket{
					Elevator:  -1, //Skal det ikke være heisens ID her?
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
					go costCompare(newOrder, OrdersToElevAlgo, OrdersToCom, costChan)
				}
			}
		}
	}
}

func costCompare(newOrder ChannelPacket, OrdersToElevAlgo, OrdersToCom, costChan chan ChannelPacket) {
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
<<<<<<< HEAD
		fmt.Printf("The cost function of elevator %d is %f\n",
			val.Elevator, val.Cost)
=======
		fmt.Printf("The cost function of elevator %d is %d", val.Cost, val.Elevator)
>>>>>>> 661181a00b56b70de672f24c6f1ecfb89b731c1b
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
	} else {
		//error, no costs received
	}
	comparing = false
}

func readFile() {
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

func writeToFile() {
	fmt.Println("before write")
	if len(localOrders) > 0 {

		file, err := os.Create("orders.csv")
		checkError("Cannot create file", err)
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		length := len(localOrders[0])
		if len(localOrders[1]) > length {
			length = len(localOrders[1])
		}
		var valueStr []string
		for j := 0; j < length; j++ {
			if j < len(localOrders[0]) {
				valueStr = append(valueStr, strconv.FormatInt(localOrders[0][j].Floor, 10)+","+strconv.FormatBool(localOrders[0][j].Direction)+",")
				valueStr[j] = valueStr[j] + strconv.FormatUint(localOrders[0][j].Timestamp, 10)
			} else {
				valueStr = append(valueStr, "0,false,0")
			}
			if j < len(localOrders[1]) {
				valueStr[j] = valueStr[j] + "," + strconv.FormatInt(localOrders[1][j].Floor, 10) + "," + strconv.FormatBool(localOrders[1][j].Direction) + ","
				valueStr[j] = valueStr[j] + strconv.FormatUint(localOrders[1][j].Timestamp, 10)
			}
		}
		err = writer.Write(valueStr)
		checkError("Cannot write to file", err)
	}
}

func addOrder(newOrder ChannelPacket) {
	data = append(data, newOrder)
	if newOrder.Elevator == thisElevator {
		localOrders[0] = append(localOrders[0], newOrder)
	} else if newOrder.Elevator == 0 {
		localOrders[1] = append(localOrders[1], newOrder)
	}
}

func removeOrder(toRemove ChannelPacket) []ChannelPacket {
	for index, value := range data {
		if value.Timestamp == toRemove.Timestamp {
			data = append(data[:index-1], data[index+1:]...)
		}
	}
	if toRemove.Elevator == thisElevator || toRemove.Elevator == 0 {
		for index, value := range data {
			if value.Timestamp == toRemove.Timestamp {
				data = append(data[:index-1], data[index+1:]...)
			}
		}
	}
	return data
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
