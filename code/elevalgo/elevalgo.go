package elevalgo

import (
	"fmt"
	"time"

	"../elevio"
	"../util"
	"../utilfunc"
	wdog "../watchdog"
)

func InitElev(elevPort string) {
	ipString := "localhost:" + elevPort
	elevio.Init(ipString, util.NumFloors)
	for i := 0; i < util.NumFloors; i++ { //Turn of all the lights in case they are still on
		elevio.SetButtonLamp(elevio.BT_Cab, i, false)
		elevio.SetButtonLamp(elevio.BT_HallDown, i, false)
		elevio.SetButtonLamp(elevio.BT_HallUp, i, false)
		fmt.Printf(" %d ", i)
	}
}

func ElevStateMachine(ElevAlgoToOrders, ComToElevAlgo, ElevAlgoToCom,
	OrdersToElevAlgo chan util.ChannelPacket, elevPort string, elevID int) {
	InitElev(elevPort)
	elevio.SetMotorDirection(elevio.MD_Up)
	elevator := util.Elev{
		State:       util.Idle,
		Dir:         util.DirUp,
		OrdersQueue: [util.NumFloors][util.NumOrderTypes]bool{},
	}
	var elevatorPtr *util.Elev = &elevator
	//Start watchdogs
	engineWatchDog := wdog.New(3 * time.Second)
	engineWatchDog.Reset()
	engineWatchDog.Stop()
	var engineFlag bool = false
	//Start timers
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()

	//Initialize channels
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)

	//Start polling
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)

	var ElevGoDirection = func(elevator *util.Elev) string {
		if elevator.Dir == util.DirDown {
			elevio.SetMotorDirection(elevio.MD_Down)
			engineWatchDog.Reset()
			elevator.State = util.Running
			return "Doing next order in queue, going down"
		} else if elevator.Dir == util.DirUp {
			elevio.SetMotorDirection(elevio.MD_Up)
			engineWatchDog.Reset()
			elevator.State = util.Running
			return "Doing next order in queue, going up"
		} else if elevator.Dir == util.DirStop {
			elevator.State = util.Idle
			return "No orders in queue"
		} else {
			return "elevator.Dir out of bounds"
		}
	}
	var IdleCheck = func() string {
		if elevator.State == util.Idle {
			elevator.Dir = utilfunc.QueueFuncChooseDirection(elevator)
			return ElevGoDirection(&elevator)
		} else {
			return "Elevator not idle, continuing on queue"
		}
	}

	for {
		utilfunc.PrintElevState(elevator)
		select {
		case a := <-OrdersToElevAlgo:
			switch a.PacketType {
			case "cabOrder":
				fmt.Printf("Recieved %s from Orders\n", a.PacketType)
				elevator.OrdersQueue[a.Floor][util.ButtonCab] = true
				if a.Floor == elevator.Floor {
					go func() { drv_floors <- int(a.Floor) }()
				} else {
					elevio.SetButtonLamp(elevio.BT_Cab, int(a.Floor), true)
					IdleCheck()
				}
			case "newOrder": //if newOrder is from orders, do the order
				fmt.Printf("Got new order from Orders, printing packet\n")
				fmt.Println(a)
				if a.Floor == elevator.Floor {
					go func() { drv_floors <- int(a.Floor) }()
				}
				utilfunc.SetOrder(a.Direction, int(a.Floor), elevatorPtr)
				fmt.Printf("%s\n", IdleCheck())
			}
		case a := <-ComToElevAlgo:
			fmt.Printf("Entering ComToElevAlgo\n")
			switch a.PacketType {
			case "requestCostFunc":
				fmt.Printf("Entering ComToElevAlgo\n Responding cost function \n")
				go func(ElevAlgoToCom chan util.ChannelPacket) {
					ElevAlgoToCom <- utilfunc.CreateCostPacket(a, elevatorPtr, engineFlag)
				}(ElevAlgoToCom)
			case "newOrder": //if newOrder is from comm, only switch on the light
				elevio.SetButtonLamp(utilfunc.DirBoolToButtonType(a.Direction), int(a.Floor), true)
			case "orderComplete":
				elevio.SetButtonLamp(elevio.BT_HallDown, int(a.Floor), false)
				elevio.SetButtonLamp(elevio.BT_HallUp, int(a.Floor), false)
			}

		case a := <-drv_buttons:
			fmt.Printf("Entering drv_buttons\n")
			//This order will go straight to orders, unless its a cab call!
			NewOrder := util.ChannelPacket{
				PacketType: "buttonPress",
				Floor:      int64(a.Floor),
				Timestamp:  uint64(time.Now().UnixNano()),
			}
			if a.Floor == int(elevator.Floor) {
				if elevator.State == util.Idle || elevator.State == util.DoorOpen {
					go func() { drv_floors <- a.Floor }()
				} else {
					if a.Button == elevio.BT_Cab {
						elevator.OrdersQueue[a.Floor][util.ButtonCab] = true
						elevio.SetButtonLamp(a.Button, a.Floor, true)
						ElevAlgoToOrders <- util.ChannelPacket{
							PacketType: "newOrder",
							Floor:      int64(a.Floor),
							Elevator:   0,
							Timestamp:  uint64(time.Now().UnixNano()),
						}
						fmt.Println(IdleCheck())
					} else {
						utilfunc.SetOrder(utilfunc.DirButtonTypeToBool(a.Button), a.Floor, &elevator)
						NewOrder.Direction = utilfunc.DirButtonTypeToBool(a.Button)
						ElevAlgoToOrders <- NewOrder
					}
				}

			} else {
				if a.Button == elevio.BT_Cab {
					elevator.OrdersQueue[a.Floor][util.ButtonCab] = true
					elevio.SetButtonLamp(a.Button, a.Floor, true)
					ElevAlgoToOrders <- util.ChannelPacket{
						PacketType: "newOrder",
						Floor:      int64(a.Floor),
						Elevator:   0,
						Timestamp:  uint64(time.Now().UnixNano()),
					}
					fmt.Println(IdleCheck())
				} else {
					NewOrder.Direction = utilfunc.DirButtonTypeToBool(a.Button)
					ElevAlgoToOrders <- NewOrder
				}
			}
		case a := <-drv_floors:
			fmt.Printf("Entering drv_floors\n")
			engineFlag = false
			engineWatchDog.Reset()
			elevio.SetFloorIndicator(a)
			fmt.Printf("We are on floor nr. %+v\n", a)
			elevator.Floor = int64(a)
			if utilfunc.QueueFuncShouldStop(elevator) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				engineWatchDog.Stop()
				utilfunc.ClearOrders(a, &elevator)
				packet := util.ChannelPacket{
					Elevator:   elevID,
					PacketType: "orderComplete",
					Floor:      elevator.Floor,
					Timestamp:  uint64(time.Now().UnixNano()),
				}
				ElevAlgoToCom <- packet          //Notifying that order is complete
				doorTimer.Reset(3 * time.Second) //begin 3 seconds of waiting for people to enter and leave car
				elevio.SetDoorOpenLamp(true)
				elevator.State = util.DoorOpen

			}

		case <-engineWatchDog.TimeOverChannel():
			fmt.Printf("Engine has timed out. Entering emergency stop mode .\n")
			engineFlag = true
			//notify the system of engine failure
			packet := util.ChannelPacket{
				PacketType: "engineTimeOut",
				Floor:      elevator.Floor,
				Direction:  utilfunc.DirIntToBool(elevator.Dir),
				Timestamp:  uint64(time.Now().UnixNano()),
			}
			ElevAlgoToOrders <- packet
		case <-doorTimer.C:
			fmt.Printf("Entering doorTimer\n")
			elevio.SetDoorOpenLamp(false)
			elevator.Dir = utilfunc.QueueFuncChooseDirection(elevator)
			if elevator.Dir == util.DirDown {
				elevio.SetMotorDirection(elevio.MD_Down)
				engineWatchDog.Reset()
				elevator.State = util.Running
			} else if elevator.Dir == util.DirUp {
				elevio.SetMotorDirection(elevio.MD_Up)
				engineWatchDog.Reset()
				elevator.State = util.Running
			} else {
				elevator.State = util.Idle
				engineWatchDog.Stop()
			}
		}
	}
}
