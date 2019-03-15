//spawne phoenix backup
//execute elevator
//calculate cost function
//
package elevAlgo

import "./elevio"
import "fmt"
import "./elevFSM"
import w "./watchdog"
import utils "./elevUtilFuncs"

var numFloors int = 4
var numOrderTypes int = 3
var currentFloor int


type order struct{
	int floor
	int direction //0 er ned og 1 er opp
}

type button int

const (
	buttonDown button = 0
	buttonCab = 1
	buttonUp = 2
)

type FSM_state int

const (
	initialize    FSM_state = 0
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

	//Start watchdogs
	engineWatchDog:= w.New(time.Second)
	engineWatchDog.Reset()
	engineWatchDog.Stop()

	
	//Start timers
	doorTimer := time.NewTimer(3* time.Second)
	doorTimer.Stop()

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

		case a := <-ordersToElevAlgo: //recieves a new ordre from orders
			elevator.Queue[a.floor][a.direction] = 1

		case a := <-comToElevAlgo:
			costFuncToCom <- calculateCostFunc(a,elevator)

		case a := <-drv_buttons:
			//This will go straight to orders, unless its a cab call!
			NewOrder := order{
				floor: a.Floor
				direction: 0
			}
			if a.Button == BT_HallUp{
				NewOrder.direction = 1
				elevAlgoToOrders <-NewOrder
			}else if a.Button == BT_HallDown{
				NewOrder.direction = 0
				elevAlgoToOrders <-NewOrder
			}else {
				elevator.Queue[a.floor][buttonCab] = 1
			}


		case a := <-drv_floors:
			var a_temp int
			if a_temp != a {
			fmt.Printf("We are on floor nr. %+v\n", a)
			elevator.Floor = a
			elevAlgoToOrders <- a //Sends the current floor to orders
			if utils.utilShouldStop(elevator){
				elevio.SetMotorDirection(elevio.MD_Stop)
				ordersQueue[a][2] = 0 //erases order from queue
				doorTimedOut.Reset(3 * time.Second)//begin 3 seconds of waiting for people to enter and leave car
				elevio.SetDoorOpenLamp(1)
			}}
			a_temp = a
			

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
			fmt.Printf("Entered shit-hit-the-fan-mode \n", a)
			elevio.SetMotorDirection(elevio.MD_Stop)
			//si ifra om emergency stop


		
		case <-engineWatchDog.TimeOverChannel():
			fmt.Printf("Engine has timed out. Entering emergency stop mode .\n")
			drv_stop <- 1
		}

		case <-doorTimedOut.C:
			elevio.SetDoorOpenLamp(0)
			elevator.Dir = chooseDirection(elevator)
			elevio.SetMotorDirection(elevator.Dir)
			if elevator.Dir == DirStop {
				elevator.State = idle
				engineWatchDog.Stop()
			}else {
				elevator.State = running
			}


	}
}
