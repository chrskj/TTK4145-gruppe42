package util

const (
	NumElevators  = 3
	NumFloors     = 4
	NumOrderTypes = 3
)

const (
	Initialize = iota
	Idle
	Running
	DoorOpen
	EmergencyStop
)

type ElevDir int

const (
	DirDown ElevDir = iota
	DirStop
	DirUp
)

type FSM_state int

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

const (
	ButtonDown Button = iota
	ButtonCab
	ButtonUp
)
