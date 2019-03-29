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

func ElevStateMachine(ElevAlgoToOrders, ComToElevAlgo, ElevAlgoToCom,
	OrdersToElevAlgo chan ChannelPacket, elevPort string, elevID int) {
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
	var engineFlag bool = false
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
	var ElevGoDirection = func(elevator *Elev) string {
		if elevator.Dir == DirDown {
			SetMotorDirection(MD_Down)
			engineWatchDog.Reset()
			elevator.State = Running
			return "Doing next order in queue, going down"
		} else if elevator.Dir == DirUp {
			SetMotorDirection(MD_Up)
			engineWatchDog.Reset()
			elevator.State = Running
			return "Doing next order in queue, going up"
		} else if elevator.Dir == DirStop {
			elevator.State = Idle
			return "No orders in queue"
		} else {
			return "elevator.Dir out of bounds"
		}
	}
	var IdleCheck = func() string {
		if elevator.State == Idle {
			elevator.Dir = QueueFuncChooseDirection(elevator)
			return ElevGoDirection(&elevator)
		} else {
			return "Elevator not idle, continuing on queue"
		}
	}

	for {
		ElevatorPrinter(elevator)
		select {
		case a := <-OrdersToElevAlgo:
			switch a.PacketType {
			case "cabOrder":
				fmt.Printf("Recieved %s from Orders\n", a.PacketType)
				elevator.OrdersQueue[a.Floor][ButtonCab] = true
				if a.Floor == elevator.Floor {
					go func() { drv_floors <- int(a.Floor) }()
				} else {
					SetButtonLamp(BT_Cab, int(a.Floor), true)
					IdleCheck()
				}
			case "newOrder": //if newOrder is from orders, do the order
				fmt.Printf("Got new order from Orders, printing packet\n")
				fmt.Println(a)
				if a.Floor == elevator.Floor {
					go func() { drv_floors <- int(a.Floor) }()
				}
				SetOrder(a.Direction, int(a.Floor), elevatorPtr)
				fmt.Printf("%s\n", IdleCheck())
			}
		case a := <-ComToElevAlgo:
			fmt.Printf("Entering ComToElevAlgo\n")
			switch a.PacketType {
			case "requestCostFunc":
				fmt.Printf("Entering ComToElevAlgo\n Responding cost function \n")
				go func(ElevAlgoToCom chan ChannelPacket) {
					ElevAlgoToCom <- CreateCostPacket(a, elevatorPtr, engineFlag)
				}(ElevAlgoToCom) //ble tidligere stuck her, bør kanskje endre
			case "newOrder": //if newOrder is from comm, only switch on the light
				SetButtonLamp(DirBoolToButtonType(a.Direction), int(a.Floor), true)
			case "orderComplete":
				SetButtonLamp(BT_HallDown, int(a.Floor), false)
				SetButtonLamp(BT_HallUp, int(a.Floor), false)
			}

		case a := <-drv_buttons:
			fmt.Printf("Entering drv_buttons\n")
			//This will go straight to orders, unless its a cab call!
			NewOrder := ChannelPacket{
				PacketType: "buttonPress",
				Floor:      int64(a.Floor),
				Timestamp:  uint64(time.Now().UnixNano()),
			}
			if a.Floor == int(elevator.Floor) {
				if elevator.State == Idle || elevator.State == DoorOpen {
					go func() { drv_floors <- a.Floor }()
				} else {
					if a.Button == BT_Cab {
						elevator.OrdersQueue[a.Floor][ButtonCab] = true
						SetButtonLamp(a.Button, a.Floor, true)
						ElevAlgoToOrders <- ChannelPacket{
							PacketType: "newOrder",
							Floor:      int64(a.Floor),
							Elevator:   0,
							Timestamp:  uint64(time.Now().UnixNano()),
						}
						fmt.Println(IdleCheck())
					} else {
						SetOrder(DirButtonTypeToBool(a.Button), a.Floor, &elevator)
						NewOrder.Direction = DirButtonTypeToBool(a.Button)
						ElevAlgoToOrders <- NewOrder
					}
				}

			} else {
				if a.Button == BT_Cab {
					elevator.OrdersQueue[a.Floor][ButtonCab] = true
					SetButtonLamp(a.Button, a.Floor, true)
					ElevAlgoToOrders <- ChannelPacket{
						PacketType: "newOrder",
						Floor:      int64(a.Floor),
						Elevator:   0,
						Timestamp:  uint64(time.Now().UnixNano()),
					}
					fmt.Println(IdleCheck())
				} else {
					NewOrder.Direction = DirButtonTypeToBool(a.Button)
					ElevAlgoToOrders <- NewOrder
				}
			}
		case a := <-drv_floors:
			fmt.Printf("Entering drv_floors\n")
			engineFlag = false
			engineWatchDog.Reset()
			SetFloorIndicator(a)
			fmt.Printf("We are on floor nr. %+v\n", a)
			elevator.Floor = int64(a)
			//elevAlgoToOrders <- a //Sends the current floor to orders
			if QueueFuncShouldStop(elevator) {
				SetMotorDirection(MD_Stop)
				engineWatchDog.Stop()
				ClearOrders(a, &elevator)
				packet := ChannelPacket{
					Elevator:   elevID,
					PacketType: "orderComplete",
					Floor:      elevator.Floor,
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
			engineFlag = true
			//SI IFRa om at motoren har stanset
			packet := ChannelPacket{
				PacketType: "engineTimeOut",
				Floor:      elevator.Floor,
				Direction:  DirIntToBool(elevator.Dir),
				Timestamp:  uint64(time.Now().UnixNano()),
			}
			ElevAlgoToOrders <- packet
		case <-doorTimer.C:
			fmt.Printf("Entering doorTimer\n")
			SetDoorOpenLamp(false)
			elevator.Dir = QueueFuncChooseDirection(elevator)
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
