package main

import (
	"./elevAlgo"
	. "./util"
)

func main() {

	//Kanal orders -> komm (orders)
	//ordersToCom := make(chan struct med noe)

	//Kanal komm -> orders (orders)
	//comToOrders := make(chan struct med noe)

	//Kanal orders -> heisalgo (ønsket floor, direction)
	OrdersToElevAlgo := make(chan Order)
	//Kanal heisalgo -> orders (current floor)
	ElevAlgoToOrders := make(chan Order)

	//Kanal komm -> heisalgo (request om cost function)
	ComToElevAlgo := make(chan Order)
	//Kanal heisalgo -> komm (cost function)
	CostFuncToCom := make(chan int)
	//Kanal heisalgo -> komm (ny ordre )
	NewOrderToCom := make(chan Order)

	//go orders(ordersToCom, comToOrders)
	go elevAlgo.ElevStateMachine(OrdersToElevAlgo, ElevAlgoToOrders, ComToElevAlgo, CostFuncToCom, NewOrderToCom)
	//go com(comToElevAlgo,elevAlgoToCom)

	//done
}
