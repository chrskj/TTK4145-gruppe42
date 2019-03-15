package util

const (
    numFloors       = 4
    numOrderTypes   = 3
)

type currentFloor int
type FSM_state int

const (
	initialize      = 0
	idle            = 1
	running         = 2
	doorOpen        = 3
	emergencyStop   = 4
)

type Direction int

const (
	DirDown Direction = iota - 1
	DirStop
	DirUp
)

type Elev struct {
	State FSM_state
	Dir Direction
	Floor int
	ordersQueue [numFloors][numOrderTypes]bool
}

type order struct {
	floor int
	direction int //0 er ned og 1 er opp
}

type ChannelPacket struct{
	packetType string
	elevatorID int
	toFloor int64
	direction int64
	timestamp uint64
	cost float64
}
