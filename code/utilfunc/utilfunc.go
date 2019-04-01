package utilfunc

import (
	"fmt"
	"math"
	"time"

	"../elevio"
	"../util"
)

func CalculateCostFunction(elevator util.Elev, order util.ChannelPacket,
	engineFlag bool) float64 {

	//If engine failure, high cost on malfunctioning elevator
	if engineFlag {
		return 9999.0
	}

	//High cost if elevator is heading in the wrong direction, or already has many orders
	switch elevator.State {
	case util.Idle:
		return math.Abs(float64(order.Floor - elevator.Floor))
	case util.Running: //Checks if the elevator is on it's way towards the potential new order
		if (elevator.Dir == 2 && (order.Floor-elevator.Floor > 0)) ||
			(elevator.Dir == 0 && (order.Floor-elevator.Floor) < 0) {
			return math.Abs(float64(order.Floor-elevator.Floor)) + 0.5*
				float64(QueueFuncCountOrders(elevator))
		} else if (elevator.Dir == 0 && (order.Floor-elevator.Floor) > 0) ||
			(elevator.Dir == 2 && (order.Floor-elevator.Floor) < 0) {
			return float64(2*util.NumFloors-elevator.Floor-order.Floor-2) + 0.5*
				float64(QueueFuncCountOrders(elevator))
		}
	case util.DoorOpen:
		if (elevator.Dir == 2 && (order.Floor-elevator.Floor > 0)) ||
			(elevator.Dir == 0 && (order.Floor-elevator.Floor) < 0) {
			return math.Abs(float64(order.Floor-elevator.Floor)) + 0.5*
				float64(QueueFuncCountOrders(elevator))
		} else if (elevator.Dir == 0 && (order.Floor-elevator.Floor) > 0) ||
			(elevator.Dir == 2 && (order.Floor-elevator.Floor) < 0) {
			return float64(2*util.NumFloors-elevator.Floor-order.Floor-2) + 0.5*
				float64(QueueFuncCountOrders(elevator))
		}
	}
	return float64(QueueFuncCountOrders(elevator))
}

func SetOrder(direction bool, floor int, elevator *util.Elev) {
	if direction {
		elevator.OrdersQueue[floor][util.ButtonUp] = true
		elevio.SetButtonLamp(elevio.BT_HallUp, floor, true)
	} else {
		elevator.OrdersQueue[floor][util.ButtonDown] = true
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, true)
	}
}

func ClearOrders(floor int, elevator *util.Elev) {
	elevator.OrdersQueue[floor][util.ButtonCab] = false //erases orders to current floor from queue
	elevator.OrdersQueue[floor][util.ButtonUp] = false
	elevator.OrdersQueue[floor][util.ButtonDown] = false
	elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
	elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}

func CreateCostPacket(order util.ChannelPacket, elevator *util.Elev,
	engineFlag bool) util.ChannelPacket {
	packet := util.ChannelPacket{
		PacketType: "cost",
		Cost: CalculateCostFunction(*elevator, util.ChannelPacket{
			Elevator:  order.Elevator,
			Floor:     order.Floor,
			Direction: order.Direction}, engineFlag),
		Timestamp: uint64(time.Now().UnixNano()),
	}
	return packet
}

func QueueFuncCountOrders(elevator util.Elev) int {
	var sum int
	for i := 0; i < util.NumFloors; i++ {
		for j := 0; j < util.NumOrderTypes; j++ {
			if elevator.OrdersQueue[i][j] {
				sum = sum + 1
			}
		}
	}
	return sum
}

func QueueFuncOrdersAboveInQueue(elevator util.Elev) bool {
	for i := elevator.Floor + 1; i < util.NumFloors; i++ {
		for j := 0; j < util.NumOrderTypes; j++ {
			if elevator.OrdersQueue[i][j] {
				return true
			}
		}
	}
	return false
}

func QueueFuncOrdersBelowInQueue(elevator util.Elev) bool {
	for i := int64(0); i < elevator.Floor; i++ {
		for j := 0; j < util.NumOrderTypes; j++ {
			if elevator.OrdersQueue[i][j] {
				return true
			}
		}
	}
	return false
}

func QueueFuncChooseDirection(elevator util.Elev) util.ElevDir {
	switch elevator.Dir {
	case util.DirUp:
		if QueueFuncOrdersAboveInQueue(elevator) {
			return util.DirUp
		} else if QueueFuncOrdersBelowInQueue(elevator) {
			return util.DirDown
		} else {
			return util.DirStop
		}
	case util.DirDown:
		if QueueFuncOrdersBelowInQueue(elevator) {
			return util.DirDown
		} else if QueueFuncOrdersAboveInQueue(elevator) {
			return util.DirUp
		} else {
			return util.DirStop
		}
	case util.DirStop:
		if QueueFuncOrdersBelowInQueue(elevator) {
			return util.DirDown
		} else if QueueFuncOrdersAboveInQueue(elevator) {
			return util.DirUp
		} else {
			return util.DirStop
		}
	}
	return util.DirStop
}

func QueueFuncShouldStop(elevator util.Elev) bool {
	switch elevator.Dir {
	case util.DirDown:
		return (elevator.OrdersQueue[elevator.Floor][util.ButtonCab] ||
			elevator.OrdersQueue[elevator.Floor][util.ButtonDown] ||
			!QueueFuncOrdersBelowInQueue(elevator) || (elevator.Floor == 0))
	case util.DirUp:
		return (elevator.OrdersQueue[elevator.Floor][util.ButtonCab] ||
			elevator.OrdersQueue[elevator.Floor][util.ButtonUp] ||
			!QueueFuncOrdersAboveInQueue(elevator) || (elevator.Floor == 3))
	default:
		return true
	}
}

func DirBoolToInt(direction bool) util.ElevDir {
	if direction {
		return util.DirUp
	} else {
		return util.DirDown
	}
}
func DirIntToBool(direction util.ElevDir) bool {
	if direction == util.DirDown {
		return false
	} else if direction == util.DirUp {
		return true
	} else {
		fmt.Printf("Error: DirStop cannot be converted to bool\n")
		return false
	}
}

func DirBoolToButtonType(direction bool) elevio.ButtonType {
	if direction {
		return elevio.BT_HallUp
	} else {
		return elevio.BT_HallDown
	}
}

func DirButtonTypeToBool(direction elevio.ButtonType) bool {
	if direction == elevio.BT_HallUp {
		return true
	} else {
		return false
	}
}

func PrintElevState(elev util.Elev) {
	fmt.Printf("State: ")
	switch elev.State {
	case 0:
		fmt.Printf("Initialize\t")
	case 1:
		fmt.Printf("Idle\t")
	case 2:
		fmt.Printf("Running\t")
	case 3:
		fmt.Printf("DoorOpen\t")
	case 4:
		fmt.Printf("EmergencyStop\t")
	}
	fmt.Printf("| Current Direction: ")
	switch elev.Dir {
	case -1:
		fmt.Printf("Going down...\t")
	case 0:
		fmt.Printf("Standing still...\t")
	case 1:
		fmt.Printf("Going up...\t")
	}
	fmt.Printf("| Floor: %d\n", elev.Floor)
	//fmt.Printf("%t\n", elev.OrdersQueue)
	for i := 0; i < len(elev.OrdersQueue); i++ {
		for j := 0; j < len(elev.OrdersQueue[0]); j++ {
			fmt.Printf("%t\t", elev.OrdersQueue[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("==========================================================\n")
}
