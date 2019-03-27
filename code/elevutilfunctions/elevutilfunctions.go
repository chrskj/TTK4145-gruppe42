package elevutilfunctions

import (
	"math"

	. "../elevio"
	. "../util"
)

func CalculateCostFunction(elevator Elev, order ChannelPacket) float64 {
	var cost float64
	//if order.Dir != elevator.Dir {
	//	cost = cost + 2
	//}
	switch elevator.State {
	case Idle:
		return math.Abs(float64(order.Floor - elevator.Floor))
	case Running: //Checks if the elevator is on it's way towards the potential new order
		if (elevator.Dir == 2 && (order.Floor-elevator.Floor > 0)) ||
			(elevator.Dir == 0 && (order.Floor-elevator.Floor) < 0) {
			return math.Abs(float64(order.Floor-elevator.Floor)) + 0.5*
				float64(QueueFuncCountOrders(elevator))
		} else if (elevator.Dir == 0 && (order.Floor-elevator.Floor) > 0) ||
			(elevator.Dir == 2 && (order.Floor-elevator.Floor) < 0) {
			return float64(2*NumFloors-elevator.Floor-order.Floor-2) + 0.5*
				float64(QueueFuncCountOrders(elevator))
		}
	case DoorOpen:
		if (elevator.Dir == 2 && (order.Floor-elevator.Floor > 0)) ||
			(elevator.Dir == 0 && (order.Floor-elevator.Floor) < 0) {
			return math.Abs(float64(order.Floor-elevator.Floor)) + 0.5*
				float64(QueueFuncCountOrders(elevator))
		} else if (elevator.Dir == 0 && (order.Floor-elevator.Floor) > 0) ||
			(elevator.Dir == 2 && (order.Floor-elevator.Floor) < 0) {
			return float64(2*NumFloors-elevator.Floor-order.Floor-2) + 0.5*
				float64(QueueFuncCountOrders(elevator))
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

func ClearOrders(floor int, elevator *Elev) {
	elevator.OrdersQueue[floor][ButtonCab] = false //erases orders to current floor from queue
	elevator.OrdersQueue[floor][ButtonUp] = false
	elevator.OrdersQueue[floor][ButtonDown] = false
	SetButtonLamp(BT_HallDown, floor, false)
	SetButtonLamp(BT_HallUp, floor, false)
	SetButtonLamp(BT_Cab, floor, false)
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
