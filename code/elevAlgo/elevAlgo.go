//spawne phoenix backup
//fikse den watchdogen
package elevAlgo

import (
	"fmt"
	"time"

	. "../QueueFunctions"
	. "../elevio"
	. "../util"
	w "../watchdog"
)

<<<<<<< HEAD
func InitElev(elevPort string){
	ipString := "localhost:" + elevPort
	Init(ipString, NumFloors)
	for i:= 0; i<NumFloors;i++{ //Turn of all the lights in case they are still on
			SetButtonLamp(BT_Cab, i, false)
			SetButtonLamp(BT_HallDown, i, false)
			SetButtonLamp(BT_HallUp, i, false)
			fmt.Printf(" %d ", i)
			}
}

func ElevStateMachine(OrdersToElevAlgo, ElevAlgoToOrders, ComToElevAlgo,
	ElevAlgoToCom chan ChannelPacket, elevPort string) {
	InitElev(elevPort)
	SetMotorDirection(MD_Up)
=======
func ElevStateMachine(OrdersToElevAlgo, ElevAlgoToOrders, ComToElevAlgo,
	ElevAlgoToCom chan ChannelPacket) {
	Init("localhost:15657", NumFloors)

	var d MotorDirection = MD_Up
	SetMotorDirection(d)
>>>>>>> windows_struggles

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
			fmt.Printf("Enteng OrdersToElevAlgo\n")
			if a.Direction {
				elevator.OrdersQueue[a.Floor][ButtonUp] = true
<<<<<<< HEAD
				SetButtonLamp(BT_HallUp, int(a.Floor), true)
=======
>>>>>>> windows_struggles
			} else {
				elevator.OrdersQueue[a.Floor][ButtonDown] = true
				SetButtonLamp(BT_HallDown, int(a.Floor), true)
			}

		case a := <-ComToElevAlgo:
			fmt.Printf("Entering ComToElevAlgo\n")
			packet := ChannelPacket{
				PacketType: "cost",
				Cost: CalculateCostFunction(elevator, ChannelPacket{
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
				elevator.OrdersQueue[a.Floor][ButtonCab] = true
<<<<<<< HEAD
				SetButtonLamp(a.Button, a.Floor, true)
				if elevator.State == Idle {
					elevator.Dir = QueueFuncChooseDirection(elevator)
					if elevator.Dir == DirDown {
						SetMotorDirection(MD_Down)
						engineWatchDog.Reset()
						elevator.State = Running
					} else if elevator.Dir == DirUp {
						SetMotorDirection(MD_Up)
						engineWatchDog.Reset()
						elevator.State = Running
=======
				if elevator.State == Idle {
					elevator.Dir = QueueFuncChooseDirection(elevator)
					elevator.State = Running
					if elevator.Dir == DirDown {
						SetMotorDirection(MD_Down)
					} else if elevator.Dir == DirUp {
						SetMotorDirection(MD_Up)
>>>>>>> windows_struggles
					} else {
						fmt.Printf("Dafuq?")
					}
				}
<<<<<<< HEAD
=======

>>>>>>> windows_struggles
			}

		case a := <-drv_floors:
			fmt.Printf("Entering drv_floors\n")
			engineWatchDog.Reset()
			if aTemp != a {
				SetFloorIndicator(a)
				fmt.Printf("We are on floor nr. %+v\n", a)
				elevator.Floor = int64(a)
				//elevAlgoToOrders <- a //Sends the current floor to orders
				if QueueFuncShouldStop(elevator) {
					SetMotorDirection(MD_Stop)
					engineWatchDog.Stop()
					elevator.Dir = DirStop
					elevator.OrdersQueue[a][ButtonCab] = false    //erases cab order from queue
					elevator.OrdersQueue[a][elevator.Dir] = false //erases order in correct direction
					SetButtonLamp(BT_Cab, a, false)               //Turn of button lamp in cab
					if elevator.Dir == DirDown { //Turn of button lamp in the correct direction
						SetButtonLamp(BT_HallDown, a, false)
					} else if elevator.Dir == DirUp {
						SetButtonLamp(BT_HallUp, a, false)
					} else {

					}
					packet := ChannelPacket{
						PacketType: "OrderComplete",
						Floor: elevator.Floor,
						Direction: DirIntToBool(elevator.Dir),
						Timestamp: uint64(time.Now().UnixNano()),
					}
					ElevAlgoToCom <- packet //Notifying that order is complete
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
			elevator.State = Running
		case a := <-drv_stop:
			elevator.State = EmergencyStop
			fmt.Printf("Entered shit-hit-the-fan-mode \n", a)
			SetMotorDirection(MD_Stop)
			//si ifra om emergency stop

			//Lage en pakke her!
			packet := ChannelPacket{
				PacketType: "EmergencyStop",
				Floor: elevator.Floor,
				Direction: DirIntToBool(elevator.Dir),
				Timestamp: uint64(time.Now().UnixNano()),
			}
			ElevAlgoToOrders <- packet

		case <-engineWatchDog.TimeOverChannel():
			fmt.Printf("Engine has timed out. Entering emergency stop mode .\n")
			drv_stop <- true

		case <-doorTimer.C:
			fmt.Printf("Entering doorTimer\n")
			SetDoorOpenLamp(false)
			elevator.Dir = QueueFuncChooseDirection(elevator)
			if elevator.Dir == DirDown {
				SetMotorDirection(MD_Down)
				//engineWatchDog.Reset()
				elevator.State = Running
			} else if elevator.Dir == DirUp {
				SetMotorDirection(MD_Up)
				//engineWatchDog.Reset()
				elevator.State = Running
			} else { //elevator.Dir == DirStop
				elevator.State = Idle
				engineWatchDog.Stop()
			}
		}
	}
}
