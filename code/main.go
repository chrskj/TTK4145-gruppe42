package main

import (
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
}
