package QueueFunctions

import (
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

func QueueFuncChooseDirection(elevator Elev) ElevDir{
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
