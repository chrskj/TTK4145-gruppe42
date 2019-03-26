//spawne phoenix backup

//Hvis i strømmen går, må den rette seg opp hvis den ser en etasje igjen
//Start og stopp funksjon til heis (for mye jobb?)
//Endre lys slik at de også skrus av når heisen er på vei opp men bestillingen er nedover

package elevAlgo

import (
	"fmt"
	"time"

	. "../elevio"
	. "../elevutilfunctions"
	. "../util"
	w "../watchdog"
)

func InitElev(elevPort string) {
	ipString := "localhost:" + elevPort
	Init(ipString, NumFloors)
	for i := 0; i < NumFloors; i++ { //Turn of all the lights in case they are still on
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

	elevator := Elev{
		State:       Idle,
		Dir:         DirUp,
		OrdersQueue: [NumFloors][NumOrderTypes]bool{},
	}
	var elevatorPtr *Elev = &elevator
	//Start watchdogs
	engineWatchDog := w.New(3 * time.Second)
	engineWatchDog.Reset()
	engineWatchDog.Stop()

	//Start timers
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()
	//var doorTimerPtr **time.Timer = &doorTimer

	//Initialize channels
	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)

	//Start polling
	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	//elevFSM.FSMinit()

	for {
		ElevatorPrinter(elevator)
		select {
		case a := <-OrdersToElevAlgo: //recieves a new ordre from orders
			fmt.Printf("Entering OrdersToElevAlgo\n Setting order\n")
			SetOrder(a.Direction, int(a.Floor), elevatorPtr)

		case a := <-ComToElevAlgo:

			fmt.Printf("Entering ComToElevAlgo\n")
			switch a.PacketType {
			case "newOrder":
				fmt.Printf("Got new order from comm\n")
				fmt.Println(a)
				if a.Direction {
					elevator.OrdersQueue[a.Floor][ButtonUp] = true
					SetButtonLamp(BT_HallUp, int(a.Floor), true)
				} else {
					elevator.OrdersQueue[a.Floor][ButtonDown] = true
					SetButtonLamp(BT_HallDown, int(a.Floor), true)
				}
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
					} else {
						fmt.Printf("Dafuq?")
					}
				}
			case "requestCostFunc":
				fmt.Printf("Entering ComToElevAlgo\n Responding cost function \n")
				go func(ElevAlgoToCom chan ChannelPacket) {
					ElevAlgoToCom <- CreateCostPacket(a, elevatorPtr)
				}(ElevAlgoToCom) //ble tidligere stuck her, bør kanskje endre
			}

		case a := <-drv_buttons:
			fmt.Printf("Entering drv_buttons\n")
			//This will go straight to orders, unless its a cab call!
			NewOrder := ChannelPacket{
				PacketType: "buttonPress",
				Floor:      int64(a.Floor),
			}
			if a.Button == BT_HallUp {
				NewOrder.Direction = true
				ElevAlgoToOrders <- NewOrder
			} else if a.Button == BT_HallDown {
				NewOrder.Direction = false
				ElevAlgoToOrders <- NewOrder
			} else {
				elevator.OrdersQueue[a.Floor][ButtonCab] = true
				SetButtonLamp(a.Button, a.Floor, true)

				if a.Floor == int(elevator.Floor) {
					go func() { drv_floors <- a.Floor }()
				}

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
					} else {
						fmt.Println("I prefer just chillin' here for a while if you don't mind")
					}
				}
			}
		case a := <-drv_floors:
			fmt.Printf("Entering drv_floors\n")
			engineWatchDog.Reset()
			SetFloorIndicator(a)
			fmt.Printf("We are on floor nr. %+v\n", a)
			elevator.Floor = int64(a)
			//elevAlgoToOrders <- a //Sends the current floor to orders
			if QueueFuncShouldStop(elevator) {
				SetMotorDirection(MD_Stop)
				engineWatchDog.Stop()
				elevator.OrdersQueue[a][ButtonCab] = false //erases cab order from queue
				SetButtonLamp(BT_Cab, a, false)            //Turn of button lamp in cab

				if elevator.Dir == DirDown { //Turn of button lamp in the correct direction
					SetButtonLamp(BT_HallDown, a, false)
				} else if elevator.Dir == DirUp {
					SetButtonLamp(BT_HallUp, a, false)
				} else {

				}

				packet := ChannelPacket{
					PacketType: "OrderComplete",
					Floor:      elevator.Floor,
					Direction:  DirIntToBool(elevator.Dir),
					Timestamp:  uint64(time.Now().UnixNano()),
				}
				ElevAlgoToCom <- packet //Notifying that order is complete
				//OpenDoor(elevatorPtr, doorTimerPtr) Prosjekt for en annen gang
				doorTimer.Reset(3 * time.Second) //begin 3 seconds of waiting for people to enter and leave car
				SetDoorOpenLamp(true)
				elevator.State = DoorOpen

			}
		//If someone is trying to get into the elevator when doors are closing,
		case <-engineWatchDog.TimeOverChannel():
			fmt.Printf("Engine has timed out. Entering emergency stop mode .\n")
			//drv_stop <- true
		case <-doorTimer.C:
			fmt.Printf("Entering doorTimer\n")
			SetDoorOpenLamp(false)
			elevator.Dir = QueueFuncChooseDirection(elevator)
			fmt.Printf("We need to go %d\n", elevator.Dir)
			if elevator.Dir == DirDown {
				SetMotorDirection(MD_Down)
				engineWatchDog.Reset()
				elevator.State = Running
			} else if elevator.Dir == DirUp {
				SetMotorDirection(MD_Up)
				engineWatchDog.Reset()
				elevator.State = Running
			} else { //elevator.Dir == DirStop
				elevator.State = Idle
				engineWatchDog.Stop()
			}
		}
	}
}
