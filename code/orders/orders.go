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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"../util"
)

var numElevators int = 3
var thisElevator int
var costChan chan ChannelPacket
var data []Order
var tstamp uint64 = 1
var ordersToCom chan ChannelPacket
var comToOrders chan ChannelPacket
var ordersToElevAlgo chan ChannelPacket
var elevAlgoToOrders chan ChannelPacket


func main() {
	data = readFile()

	writeToFile()
}

func Initialize(ordersToCom chan ChannelPacket, comToOrders chan ChannelPacket, ordersToElevAlgo chan ChannelPacket, elevAlgoToOrders chan ChannelPacket) {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			<-ticker.C
			tstamp++
		}
	}()
	//try to get data from others

	go orderRoutine(ordersToCom, comToOrders, ordersToElevAlgo, elevAlgoToOrders)
}

func orderRoutine(ordersToCom chan ChannelPacket, comToOrders chan ChannelPacket, ordersToElevAlgo chan ChannelPacket, elevAlgoToOrders chan ChannelPacket) {
	var costChan chan ChannelPacket
	for {
		select {
		case temp := <-comToOrders:
			switch temp.PacketType {
			case "elevID":
				thisElevator = temp.Elevator
			case "compareCost":
				costChan <- temp
			case "orderComplete":
				removeOrder(Order{
					Elevator:  temp.Elevator,
					Floor:   temp.Floor,
					Direction: temp.Direction,
					Timestamp: temp.Timestamp,
				})
			case "addOrder":
				addOrder(Order{
					Elevator:  temp.Elevator,
					Floor:   temp.Floor,
					Direction: temp.Direction,
					Timestamp: temp.Timestamp,
				})
			case "getOrderList":
				packet := ChannelPacket{
					PacketType: "orderList",
					DataJson:   getOrderJson(),
				}
				ordersToCom <- packet
			case "orderJson":
				json.Unmarshal(temp.DataJson, data)
			}
		case temp := <-elevAlgoToOrders:
			switch temp.PacketType {
			case "buttonPress":
				newOrder := Order{
					Elevator:  -1,
					Floor:   temp.Floor,
					Direction: temp.Direction,
					Timestamp: tstamp,
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
					go costCompare(newOrder, ordersToCom)
				}
			}
		}
	}
}

func costCompare(newOrder Order, ordersToCom chan ChannelPacket) {
	ordersToCom <- ChannelPacket{
		PacketType: "requestCostFunc",
		Elevator:   thisElevator,
	}
	costTicker := time.NewTicker(10 * time.Millisecond)
	var ticks uint = 0
	var costs []ChannelPacket
	for recievedOrders := 0; recievedOrders < numElevators && ticks < 200; {
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
		if val.cost < max {
			max = val.cost
			newOrder.Elevator = val.Elevator
		}
	}
	if newOrder.Elevator != -1 {
		data = addOrder(newOrder)
		temp := ChannelPacket{
			PacketType: "newOrder"
			Elevator: newOrder.Elevator
			Floor:	newOrder.Floor
			Direction:	newOrder.Direction
			Timestamp: newOrder.Timestamp
		}
		if temp.Elevator == thisElevator {
			ordersToElevAlgo <- temp
		} else {
			ordersToCom <- temp
		}
	} else {
		//error, no costs received
	}
}

func readFile() []Order {
	var data []Order
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
		for i := 0; i < numElevators; i++ {
			FloorTemp, _ := strconv.ParseInt(input[0+3*i], 10, 64)
			DirectionTemp, _ := strconv.ParseBool(input[1+3*i])
			tstampTemp, _ := strconv.ParseUint(input[2+3*i], 10, 64)
			ElevatorTemp := i + 1
			data = append(data, Order{
				Elevator:  ElevatorTemp,
				Floor:   FloorTemp,
				Direction: DirectionTemp,
				Timestamp: tstampTemp,
			})
		}
		FloorTemp, _ := strconv.ParseInt(input[3*numElevators], 10, 64)
		tstampTemp, _ := strconv.ParseUint(input[3*numElevators+1], 10, 64)
		data = append(data, Order{
			Elevator:  0,
			Floor:   FloorTemp,
			Timestamp: tstampTemp,
		})
	}
	return data
}

func getOrderJson() []byte {
	var temp []Order
	for i := 0; i < len(data)/4; i++ {
		for j := 0; j < 3; j++ {
			temp = append(temp, data[i])
		}
	}
	valueJson, _ := json.Marshal(temp)
	return valueJson
}

func writeToFile() {
	fmt.Println("before write")
	file, err := os.Create("orders.csv")
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	//var writeData []string
	for i := 0; i < (len(data) / (numElevators + 1)); i++ {
		values := data[((numElevators + 1) * i):((numElevators+1)*i + (numElevators + 1))]
		var value []string
		for j := 0; j < 3; j++ {
			value = append(value, strconv.FormatInt(values[j].Floor, 10))
			value = append(value, strconv.FormatBool(values[j].Direction))
			value = append(value, strconv.FormatUint(values[j].Timestamp, 10))
		}
		value = append(value, strconv.FormatInt(values[numElevators].Floor, 10))
		value = append(value, strconv.FormatUint(values[numElevators].Timestamp, 10))
		var valueStr []string
		for j := 0; j < 3*numElevators+1; j++ {
			valueStr = append(valueStr, value[j]) // + ","
		}
		valueStr = append(valueStr, value[3*numElevators])
		valueStr = append(valueStr, value[3*numElevators+1])
		valueStr = valueStr[:len(valueStr)-1]
		//writeData = append(writeData, valueStr)
		err = writer.Write(valueStr)
		checkError("Cannot write to file", err)
	}
	/*
		err = writer.Write(writeData)
		checkError("Cannot write to file", err)
	*/
}

func addOrder(newOrder Order) []Order {
	blankOrder := Order{
		Elevator:  0,
		Floor:   0,
		Direction: false,
		Timestamp: 0,
	}
	for index, value := range data {
		if value.Timestamp == 0 && value.Elevator == newOrder.Elevator {
			data[index] = newOrder
			newOrder.Timestamp = 0
		}
	}
	if newOrder.Timestamp != 0 {
		data = append(data, newOrder)
		for i := 0; i < numElevators; i++ {
			data = append(data, blankOrder)
		}
	}
	return data
}

func removeOrder(toRemove Order) []Order {
	blankOrder := Order{
		Elevator:  0,
		Floor:   0,
		Direction: false,
		Timestamp: 0,
	}
	for index, value := range data {
		if value == toRemove {
			data[index] = blankOrder
		}
	}
	return data
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
        TESTLINE
	}
}
