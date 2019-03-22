package main

import (
<<<<<<< HEAD
    "fmt"
    //orders "github.com/chrskj/TTK4145-gruppe44/code/communication"
    //elevAlgo "github.com/chrskj/TTK4145-gruppe44/code/communication"
    com "github.com/chrskj/TTK4145-gruppe44/code/communication"
    . "github.com/chrskj/TTK4145-gruppe44/code/util"
)

func main() {
    fmt.Println("Started")

    //elevAlgoToOrders := make(chan ChannelPacket)
    //ordersToElevAlgo := make(chan ChannelPacket)

    comToElevAlgo := make(chan ChannelPacket)
    elevAlgoToCom := make(chan ChannelPacket)

    ordersToCom := make(chan ChannelPacket)
    comToOrders := make(chan ChannelPacket)

    //go orders.Initialize(ordersToCom, comToOrders, ordersToElevAlgo,
    //    elevAlgoToOrders)
    //go elevAlgo.Initialize(elevAlgoToCom, elevAlgoToOrders, comToElevAlgo,
    //    ordersToElevAlgo)
    go com.Initialize(comToElevAlgo, comToOrders, elevAlgoToCom, ordersToCom)

    for{}
=======
	"./elevAlgo"
	. "./util"
)

func main() {

	//Kanal orders -> komm (orders)
  OrdersToCom := make(chan struct ChannelPacket)

  //Kanal komm -> orders (orders)
  ComToOrders := make(chan struct ChannelPacket)


  //Kanal orders -> heisalgo (ønsket floor, direction)
  OrdersToElevAlgo := make(chan ChannelPacket)
  //Kanal heisalgo -> orders (current floor)
  ElevAlgoToOrders := make(chan ChannelPacket)

  //Kanal komm -> heisalgo (request om cost function)
  ComToElevAlgo := make(chan orderStruct)
  //Kanal heisalgo -> komm (cost function)
  ElevAlgoToCom := make(chan orderStruct)

  go orders.initialize(ordersToCom, comToOrders,ordersToElevAlgo,elevAlgoToOrders)
  go elevAlgo.ElevStateMachine(OrdersToElevAlgo, ElevAlgoToOrders, ComToElevAlgo, CostFuncToCom, NewOrderToCom)
  go com(comToElevAlgo,elevAlgoToCom)

	//done
>>>>>>> Noe_spennende
}
