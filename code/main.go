package main

import ( 
	"./elevAlgo"
	"fmt"
	"os"
	"./util"

)

type orderStruct struct {
	int floor
	int direction //0 er ned og 1 er opp
} 


//Kanal orders -> komm (orders)
//ordersToCom := make(chan struct med noe)

//Kanal komm -> orders (orders)
//comToOrders := make(chan struct med noe)


//Kanal orders -> heisalgo (ønsket floor, direction)
OrdersToElevAlgo := make(chan Order)
//Kanal heisalgo -> orders (current floor)
ElevAlgoToOrders := make(chan Order)

//Kanal komm -> heisalgo (request om cost function)
ComToElevAlgo := make(chan orderStruct)
//Kanal heisalgo -> komm (cost function)
CostFuncToCom := make(chan int)
//Kanal heisalgo -> komm (ny ordre )
NewOrderToCom := make(chan Order)

//go orders(ordersToCom, comToOrders)
go elevAlgo.ElevStateMachine(ordersToElevAlgo,elevAlgoToOrders, comToElevAlgo,costFuncToCom,newOrderToCom)
//go com(comToElevAlgo,elevAlgoToCom)

//done

