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

type Order struct {
    Dir Direction
    Floor int
}
