package queue

import "./elevio"
import "fmt"



//første indeks: etasjenr, andre indeks:
// 0 = ned, 1 = opp, 2 = cab. True or false
var ordersQueue bool[numFloors][numOrderTypes] 

func QcountOrders(){
	var sum int
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numOrderTypes; j++ {
			sum = sum + ordersQueue[i][j]
		}
	}
	return sum
}

func QshouldStop(Elev elevator){
	if ordersQueue[floor][2]{
		return 1
	}
	select {
	case dir = 0:
		return QiterateDownwardCallDown(3) // = floor?

	case dir = 1:
		return Qite
	}
}