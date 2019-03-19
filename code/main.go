package main

import ( 
	"./elevAlgo"
	"fmt"
	"os"


)

type orderStruct struct {
	int floor
	int direction //0 er ned og 1 er opp
} 


//Kanal orders -> komm (orders)
ordersToCom := make(chan struct med noe)

//Kanal komm -> orders (orders)
comToOrders := make(chan struct med noe)


//Kanal orders -> heisalgo (ønsket floor, direction)
ordersToElevAlgo := make(chan orderStruct)
//Kanal heisalgo -> orders (current floor)
elevAlgoToOrders := make(chan int)

//Kanal komm -> heisalgo (request om cost function)
comToElevAlgo := make(chan orderStruct)
//Kanal heisalgo -> komm (cost function)
costFuncToCom := make(chan float)
//Kanal heisalgo -> komm (ny ordre )
newOrderToCom := make(chan orderStruct)

go orders(ordersToCom, comToOrders)
go elevAlgo.main(ordersToElevAlgo,elevAlgoToOrders, comToElevAlgo,costFuncToCom,newOrderToCom)
go com(comToElevAlgo,elevAlgoToCom)

//done

