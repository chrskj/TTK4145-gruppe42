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

type Direction int

const (
	DirDown Direction = iota - 1
	DirStop
	DirUp
)

type Elev struct {
	State       FSM_state
	Dir         Direction
	Floor       int
	OrdersQueue [NumFloors][NumOrderTypes]bool //int fordi lettere å utføre matematiske operasjoner senere.
}

type Order struct {
	Dir   Direction //0 er ned og 1 er opp
	Floor int
}

type button int

const ( //stor bokstav
	ButtonDown button = 0
	ButtonCab         = 1
	ButtonUp          = 2
)
