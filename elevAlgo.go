//Keep a queue
//spawne phoenix backup
//execute elevator
//calculate cost function
//
package elevAlgo

import "./elevio"
import "fmt"

//Variables that need initializing:
// -orders
// -current floor
var numFloors int= 4
var numOrderTypes int = 3
var currentFloor int
var ordersQueue bool[numFloors][numOrderTypes] //første indeks: etasjenr, andre indeks: 0 = ned, 1 = opp, 2 = cab. True or false 

func countOrders(){
	var sum int
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numOrderTypes; j++ {
			sum = sum + ordersQueue[i][j]
		}
	}
}

func calculateCostFunc(orderStruct order){
	return countOrders()
}

func FSM(){
	
}

func main(ordersToElevAlgo,elevAlgoToOrders, comToElevAlgo,costFuncToCom,newOrderToCom) {

	numFloors := 4
	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

a = 3

	for {
		select {

		case a := <-ordersToElevAlgo: //recieves a new ordre from orders
			ordersQueue[a.floor][a.direction] = 1

		case a := <-comToElevAlgo:
			costFuncToCom <- calculateCostFunc(a)

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
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}
