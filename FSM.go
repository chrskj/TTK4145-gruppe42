package elevFSM

import "./elevio"

func FSMinit() {
	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	elevio.SetMotorDirection(d)
	if 
}

func FSMidle() {
	if countOrders() {
		
	}
}

func FSMrunning() {

}

func FSMdoorOpen() {

}

func FSMemergencyStop() {
	elevio.SetMotorDirection(elevio.MD_Stop)
	//set stop lamp
	//clear all orders
	//tell the others!
	for f := 0; f < numFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}

}
