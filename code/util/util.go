package util

const (
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

type Order struct {
	Elevator  int
	Floor     int64
	Direction bool //True = opp, False = ned
	Timestamp uint64
}

type ChannelPacket struct {
	PacketType string
	Elevator   int
	Floor    int64
	Direction  bool
	Timestamp  uint64
	Cost       float64
	DataJson   []byte
}

type button int

const ( //stor bokstav
	ButtonDown button = 0
	ButtonCab         = 1
	ButtonUp          = 2
)
