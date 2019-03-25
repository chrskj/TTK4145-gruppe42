package util

import (
	"fmt"
)

const (
	NumElevators = 3
	NumFloors     = 4
	NumOrderTypes = 3
)

type currentFloor int //stor bokstav
type FSM_state int

const ( //stor bokstav
	Initialize    = 0
	Idle          = 1
	Running       = 2
	DoorOpen      = 3
	EmergencyStop = 4
)

type Direction int //for elevator IO use, not orders

const (
	DirDown Direction = iota - 1
	DirStop
	DirUp
)

type Elev struct {
	State       FSM_state
	Dir         Direction
	Floor       int64
	OrdersQueue [NumFloors][NumOrderTypes]bool
}

type ChannelPacket struct {
	PacketType string
	Elevator   int
	Floor      int64
	Direction  bool
	Timestamp  uint64
	Cost       float64
	OrderList   []ChannelPacket
}

type button int

const ( //stor bokstav
	ButtonDown button = 0
	ButtonCab         = 1
	ButtonUp          = 2
)

func ElevatorPrinter(elev Elev) {
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
	switch elev.Dir {
	case -1:
		fmt.Printf("Going down...\t")
	case 0:
		fmt.Printf("Standing still...\t")
	case 1:
		fmt.Printf("Going up...\t")
	}
	fmt.Printf("On floor %d\n", elev.Floor)
	fmt.Printf("%t\n", elev.OrdersQueue)
}

func DirBoolToInt(direction bool) Direction{
	if direction{
		return DirUp
	}else{
		return DirDown
	}
}

func DirIntToBool(direction Direction)bool{
	if direction == DirDown{
		return false
	}else if direction == DirUp{
		return true
	}else{
		fmt.Printf("Error: DirStop cannot be converted to bool\n")
		return false
	}
}