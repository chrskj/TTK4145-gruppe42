//spawne phoenix backup
//execute elevator
//calculate cost function
//
package elevAlgo

import ( 
	"../elevio"
 	"fmt"
	 w "../watchdog"
	utils "../queueFunctions" 
	. "../util"
)

var currentFloor int
var a_temp int




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
		if a.Direction == 1{
			elevator.Queue[a.floor][DirUp] = 1
		}else if a.Direction == 0{
			elevator.Queue[a.floor][DirDown] = 1
		}else{
			fmt.Printf("Something fishy in the orders from Orders, not 0 or 1!")
		}

			

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
			if a_temp != a {
			fmt.Printf("We are on floor nr. %+v\n", a)
			elevator.Floor = a
			elevAlgoToOrders <- a //Sends the current floor to orders
			if utils.utilShouldStop(elevator){
				elevio.SetMotorDirection(elevio.MD_Stop)
				ordersQueue[a][buttonCab] = 0 //erases cab order from queue
				ordersQueue[a][a.direction+1] = 0 //erases order in correct direction
				//notify orders that its done!
				doorTimedOut.Reset(3 * time.Second)//begin 3 seconds of waiting for people to enter and leave car
				elevio.SetDoorOpenLamp(1)
			}}
			a_temp = a
			
		//If someone is trying to get into the elevator when doors are closing, the elevator will wait 3 more seconds
		case a := <-drv_obstr:
			fmt.Printf("Obstruction in door! Someone is trying to get in! \n")
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevState = doorOpen
			doorTimedOut.Reset(3*time.Second) 
			

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
