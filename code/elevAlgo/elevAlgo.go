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
	ElevAlgoToCom chan ChannelPacket, elevPort string) {
	ipString := "localhost:" + elevPort
	Init(ipString, NumFloors)

	var d MotorDirection = MD_Up
	SetMotorDirection(d)

	elevator := Elev{
		State:       Idle,
		Dir:         DirUp,
		OrdersQueue: [NumFloors][NumOrderTypes]bool{},
	}
	var aTemp int = -1
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
		ElevatorPrinter(elevator)
		select {
		case a := <-OrdersToElevAlgo: //recieves a new ordre from orders
			fmt.Printf("Entering OrdersToElevAlgo\n")
			if a.Direction {
				elevator.OrdersQueue[a.Floor][ButtonUp] = true
				SetButtonLamp(BT_HallUp, int(a.Floor), true)
			} else {
				elevator.OrdersQueue[a.Floor][ButtonDown] = true
				SetButtonLamp(BT_HallDown, int(a.Floor), true)
			}

		case a := <-ComToElevAlgo:
			fmt.Printf("Entering ComToElevAlgo\n")
			packet := ChannelPacket{
				PacketType: "cost",
				Cost: CalculateCostFunction(elevator, Order{
					Elevator:  a.Elevator,
					Floor:     a.Floor,
					Direction: a.Direction}),
			}
			ElevAlgoToCom <- packet

		case a := <-drv_buttons:
			fmt.Printf("Entering drv_buttons\n")
			//This will go straight to orders, unless its a cab call!
			NewOrder := ChannelPacket{
				PacketType: "buttonPress",
				Floor:      int64(a.Floor),
			}
			fmt.Printf("%d\n", a.Button)
			if a.Button == BT_HallUp {
				NewOrder.Direction = true
				ElevAlgoToOrders <- NewOrder
			} else if a.Button == BT_HallDown {
				NewOrder.Direction = false
				ElevAlgoToOrders <- NewOrder
			} else {
				fmt.Printf("Why the hell did I end up here?")
				elevator.OrdersQueue[a.Floor][ButtonCab] = true
				SetButtonLamp(a.Button, a.Floor, true)
				if elevator.State == Idle {
					elevator.Dir = QueueFuncChooseDirection(elevator)
					elevator.State = Running
					if elevator.Dir == DirDown {
						SetMotorDirection(MD_Down)
					} else if elevator.Dir == DirUp {
						SetMotorDirection(MD_Up)
					} else {
						fmt.Printf("Dafuq?")
					}
				}

			}

		case a := <-drv_floors:
			fmt.Printf("Entering drv_floors\n")
			if aTemp != a {
				SetFloorIndicator(a)
				fmt.Printf("We are on floor nr. %+v\n", a)
				elevator.Floor = int64(a)
				//elevAlgoToOrders <- a //Sends the current floor to orders
				if QueueFuncShouldStop(elevator) {
					SetMotorDirection(MD_Stop)
					elevator.Dir = DirStop
					elevator.OrdersQueue[a][ButtonCab] = false    //erases cab order from queue
					elevator.OrdersQueue[a][elevator.Dir] = false //erases order in correct direction
					SetButtonLamp(BT_Cab, a, false)               //Turn of button lamp in cab
					//SetButtonLamp(elevator.Dir, a, false)         //Turn of button lamp in the correct direction
					if elevator.Dir == DirDown {
						SetButtonLamp(BT_HallDown, a, false)
					} else if elevator.Dir == DirUp {
						SetButtonLamp(BT_HallUp, a, false)
					} else {

					}

					//notify orders that its done!
					doorTimer.Reset(3 * time.Second) //begin 3 seconds of waiting for people to enter and leave car
					SetDoorOpenLamp(true)
					elevator.State = DoorOpen
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
			fmt.Printf("Entering doorTimer\n")
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
