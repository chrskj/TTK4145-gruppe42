//spawne phoenix backup
//execute elevator
//calculate cost function
//comment
package elevAlgo

import (
	"fmt"
	"time"

	. "../QueueFunctions"
	. "../elevio"
	. "../util"
	w "../watchdog"
)

func ElevStateMachine(OrdersToElevAlgo, ElevAlgoToOrders, ComToElevAlgo,
		ElevAlgoToCom chan ChannelPacket) {	
	Init("localhost:15657", NumFloors)


    var d MotorDirection = MD_Up
	SetMotorDirection(d)
	
	elevator := Elev{
		State:       Idle,
		Dir:         DirStop,
		OrdersQueue: [NumFloors][NumOrderTypes]bool{},
	}
	var aTemp int
	//Start watchdogs
	engineWatchDog := w.New(time.Second)
	engineWatchDog.Reset()
	engineWatchDog.Stop()

	//Start timers
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()

	//Initialize channels
	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	//Start polling
	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)
	go PollStopButton(drv_stop)
	//elevFSM.FSMinit()

	for {
		select {

		case a := <-OrdersToElevAlgo: //recieves a new ordre from orders
			if a.Direction {
				elevator.OrdersQueue[a.Floor][ButtonUp] = true
			} else{
				elevator.OrdersQueue[a.Floor][ButtonDown] = true
			}

		case a := <-ComToElevAlgo:
			packet := ChannelPacket{
				PacketType: "cost",
				Cost: CalculateCostFunction(elevator, Order{
					Elevator : a.Elevator,
					Floor: a.Floor,
					Direction: a.Direction}),
			}
			ElevAlgoToCom <- packet

		case a := <-drv_buttons:
			//This will go straight to orders, unless its a cab call!
			NewOrder := ChannelPacket{
				Floor: int64(a.Floor),
			}
			if a.Button == BT_HallUp {
				NewOrder.Direction = true
				ElevAlgoToOrders <- NewOrder
			} else if a.Button == BT_HallDown {
				NewOrder.Direction = false
				ElevAlgoToOrders <- NewOrder
			} else {
				elevator.OrdersQueue[a.Floor][ButtonCab] = true
			}

		case a := <-drv_floors:
			if aTemp != a {
				fmt.Printf("We are on floor nr. %+v\n", a)
				elevator.Floor = int64(a)
				//elevAlgoToOrders <- a //Sends the current floor to orders
				if QueueFuncShouldStop(elevator) {
					SetMotorDirection(MD_Stop)
					elevator.OrdersQueue[a][ButtonCab] = false    //erases cab order from queue
					elevator.OrdersQueue[a][elevator.Dir] = false //erases order in correct direction
					//notify orders that its done!
					doorTimer.Reset(3 * time.Second) //begin 3 seconds of waiting for people to enter and leave car
					SetDoorOpenLamp(true)
				}
			}
			aTemp = a

		//If someone is trying to get into the elevator when doors are closing, the elevator will wait 3 more seconds
		case <-drv_obstr:
			fmt.Printf("Obstruction in door! Someone is trying to get in! \n")
			SetMotorDirection(MD_Stop)
			elevator.State = DoorOpen
			doorTimer.Reset(3 * time.Second)

		case a := <-drv_stop:
			elevator.State = EmergencyStop
			fmt.Printf("Entered shit-hit-the-fan-mode \n", a)
			SetMotorDirection(MD_Stop)
			//si ifra om emergency stop

		case <-engineWatchDog.TimeOverChannel():
			fmt.Printf("Engine has timed out. Entering emergency stop mode .\n")
			drv_stop <- true

		case <-doorTimer.C:
			SetDoorOpenLamp(false)
			elevator.Dir = QueueFuncChooseDirection(elevator)
			//SetMotorDirection(elevator.Dir)
			if elevator.Dir == DirDown {
				SetMotorDirection(MD_Down)
				elevator.State = Running
			} else if elevator.Dir == DirUp {
				SetMotorDirection(MD_Up)
				elevator.State = Running
			} else { //elevator.Dir == DirStop
				elevator.State = Idle
				engineWatchDog.Stop()
			}

		}

	}
}
