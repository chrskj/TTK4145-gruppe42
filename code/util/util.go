package util

const (
	NumElevators  = 3
	NumFloors     = 4
	NumOrderTypes = 3
)

type currentFloor int //stor bokstav
type FSM_state int

const ( //stor bokstav
	Initialize = iota
	Idle
	Running
	DoorOpen
	EmergencyStop
)

type ElevDir int //for elevator IO use, not orders

const (
	DirDown ElevDir = iota
	DirStop
	DirUp
)

type Elev struct {
	State       FSM_state
	Dir         ElevDir
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
	OrderList  []ChannelPacket
}

type Button int

const ( //stor bokstav
	ButtonDown Button = iota
	ButtonCab
	ButtonUp
)
