package elevutilfunctions

import (
	. "../elevio"
	. "../util"
)

//stor bokstav på ALT !!

func CalculateCostFunction(elevator Elev, order ChannelPacket) float64 {
	var cost float64
	//if order.Direction != elevator.Dir {
	//	cost = cost + 2
	//}
	switch elevator.Dir {
	case DirDown:
		if order.Floor > elevator.Floor {
			cost = cost + 2
		}
	case DirUp:
		if order.Floor < elevator.Floor {
			cost = cost + 2
		}
	}
	return float64(QueueFuncCountOrders(elevator)) + cost
}

/*func OpenDoor(elevator *Elev, doorTimer **time.Timer) {
	doorTimer..Reset(3 * time.Second) //begin 3 seconds of waiting for people to enter and leave car
	SetDoorOpenLamp(true)
	elevator.State = DoorOpen
}*/

func SetOrder(direction bool, floor int, elevator *Elev) {
	if direction {
		elevator.OrdersQueue[floor][ButtonUp] = true
		SetButtonLamp(BT_HallUp, floor, true)
	} else {
		elevator.OrdersQueue[floor][ButtonDown] = true
		SetButtonLamp(BT_HallDown, floor, true)
	}
}

func CreateCostPacket(order ChannelPacket, elevator *Elev) ChannelPacket {
	packet := ChannelPacket{
		PacketType: "cost",
		Cost: CalculateCostFunction(*elevator, ChannelPacket{
			Elevator:  order.Elevator,
			Floor:     order.Floor,
			Direction: order.Direction}),
	}
	return packet
}

func QueueFuncCountOrders(elevator Elev) int {
	var sum int
	for i := 0; i < NumFloors; i++ {
		for j := 0; j < NumOrderTypes; j++ {
			if elevator.OrdersQueue[i][j] {
				sum = sum + 1
			}

		}
	}
	return sum
}

func QueueFuncOrdersAboveInQueue(elevator Elev) bool {
	for i := elevator.Floor + 1; i < NumFloors; i++ {
		for j := 0; j < NumOrderTypes; j++ {
			if elevator.OrdersQueue[i][j] {
				return true
			}
		}
	}
	return false
}

func QueueFuncOrdersBelowInQueue(elevator Elev) bool {
	for i := int64(0); i < elevator.Floor; i++ {
		for j := 0; j < NumOrderTypes; j++ {
			if elevator.OrdersQueue[i][j] {
				return true
			}
		}
	}
	return false
}

func QueueFuncChooseDirection(elevator Elev) ElevDir {
	switch elevator.Dir {
	case DirUp:
		if QueueFuncOrdersAboveInQueue(elevator) {
			return DirUp
		} else if QueueFuncOrdersBelowInQueue(elevator) {
			return DirDown
		} else {
			return DirStop
		}
	case DirDown:
		if QueueFuncOrdersBelowInQueue(elevator) {
			return DirDown
		} else if QueueFuncOrdersAboveInQueue(elevator) {
			return DirUp
		} else {
			return DirStop
		}
	case DirStop:
		if QueueFuncOrdersBelowInQueue(elevator) {
			return DirDown
		} else if QueueFuncOrdersAboveInQueue(elevator) {
			return DirUp
		} else {
			return DirStop
		}
	}
	return DirStop
}

func QueueFuncShouldStop(elevator Elev) bool {
	switch elevator.Dir {
	case DirDown:
		return (elevator.OrdersQueue[elevator.Floor][ButtonCab] ||
			elevator.OrdersQueue[elevator.Floor][ButtonDown] ||
			!QueueFuncOrdersBelowInQueue(elevator))
	case DirUp:
		return (elevator.OrdersQueue[elevator.Floor][ButtonCab] ||
			elevator.OrdersQueue[elevator.Floor][ButtonUp] ||
			!QueueFuncOrdersAboveInQueue(elevator))
	default:
		return true
	}
}
