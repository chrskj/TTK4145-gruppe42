package queueFunctions



func calculateCostFunction(elevator Elev, order Order){
	var cost int = 0
	orderDir := order.Dir
	orderFloor := order.Floor
	if order.Dir != elevator.Dir{
		cost = cost + 2
	}
	switch elevator.Dir {
	case DirDown:
		if order.Floor > elev.Floor:
		cost = cost + 2
	case DirUp:
		if order.Floor < elev.Floor:
		cost = cost + 2
	}
	return utilCountOrders(elev) + cost
}

func utilCountOrders(elevator Elev) {
	var sum int
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numOrderTypes; j++ {
			sum = sum + ordersQueue[i][j]
		}
	}
	return sum
}

func utilShouldStop(elevator Elev) bool {
	switch elevator.Dir {
	case DirDown:
		return elevator.ordersQueue[elevator.Floor][buttonCab] ||
			elevator.ordersQueue[elevator.Floor][buttonDown] ||
			!ordersBelowInQueue(elevator)
	case DirUp:
		return elevator.ordersQueue[elevator.Floor][buttonCab] ||
			elevator.ordersQueue[elevator.Floor][buttonUp] ||
			!ordersAboveInQueue(elevator)

	default:

	}
	return false
}

func ordersAboveInQueue(elevator Elev) bool {
	for i := elevator.Floor; i < numFloors; i++ {
		for j := 0; j < numOrderTypes; j++ {
			if elevator.ordersQueue[i][j] {
				return true
			}

		}
	}
	return false
}

func ordersBelowInQueue(elevator Elev) bool {
	for i := 0; i < elevator.Floor; i++ {
		for j := 0; j < numOrderTypes; j++ {
			if elevator.ordersQueue[i][j] {
				return true
			}
		}
	}
	return false
}

func chooseDirection(elevator Elev) Direction {
	switch elevator.Dir {
	case DirStop:
		if ordersAboveInQueue(elevator) {
			return DirUp
		} else if ordersBelowInQueue(elevator) {
			return DirDown
		} else {
			return DirStop
		}

	case DirDown:
		if ordersBelowInQueue(elevator) {
			return DirDown
		} else if ordersAboveInQueue(elevator) {
			DirUp
		} else {
			DirStop
		}

	case DirUp:
		if ordersAboveInQueue(elevator) {
			return DirUp
		} else if ordersBelowInQueue(elevator) {
			DirDown
		} else {
			DirStop
		}
	}
	return DirStop
}
