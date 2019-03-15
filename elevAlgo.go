//Keep a queue
//spawne phoenix backup
//execute elevator
//calculate cost function
//
package elevAlgo

import "./elevio"
import "fmt"
import "./elevFSM"

//Variables that need initializing:
// -orders
// -current floor
var numFloors int = 4
var numOrderTypes int = 3
var currentFloor int

type FSM_state int

const (
	init          FSM_state = 0
	idle          FSM_state = 1
	running       FSM_state = 2
	doorOpen      FSM_state = 3
	emergencyStop FSM_state = 4
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

func calculateCostFunc(orderStruct order) {
	return countOrders()
}

func main(ordersToElevAlgo, elevAlgoToOrders, comToElevAlgo, costFuncToCom, newOrderToCom) {
	elevator := Elev{
		State: idle,
		Dir: DirStop,
		Floor: //get floor sensor signal
		Queue: [numFloors][numOrderTypes]bool{},
	}

	//Start timers
	//Doortimer
	//engineFailureTimer
	//

	//Initialize channels
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	//Start polling
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	elevFSM.FSMinit()

	for {
		select {

		//Write an FSM function under each case?
		case a := <-ordersToElevAlgo: //recieves a new ordre from orders
			ordersQueue[a.floor][a.direction] = 1

		case a := <-comToElevAlgo:
			costFuncToCom <- calculateCostFunc(a,elevator)

		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)

			elevio.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			currentFloor = a
			elevAlgoToOrders <- a //Sends the current floor to orders
			ordersQueue[a][2] = 0
			if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevState = doorOpen
			} else {
				elevio.SetMotorDirection(d)
				elevState = running
			}

		case a := <-drv_stop:
			elevState = emergencyStop
			fmt.Printf("%+v\n", a)
			elevFSM.FSMemergencyStop()

		}
	}
}
