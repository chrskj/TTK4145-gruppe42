package main

import (
  "fmt"
	"./orders"
	"./communication"
	"./elevAlgo"
	. "./util"
)

//comment

func main(){
  fmt.Println("Started")

  //Kanal orders -> komm (orders)
  OrdersToCom := make(chan ChannelPacket)

  //Kanal komm -> orders (orders)
  ComToOrders := make(chan ChannelPacket)


  //Kanal orders -> heisalgo (ønsket floor, direction)
  OrdersToElevAlgo := make(chan ChannelPacket)
  //Kanal heisalgo -> orders (current floor)
  ElevAlgoToOrders := make(chan ChannelPacket)

  //Kanal komm -> heisalgo (request om cost function)
  ComToElevAlgo := make(chan ChannelPacket)
  //Kanal heisalgo -> komm (cost function)
  ElevAlgoToCom := make(chan ChannelPacket)

  go orders.InitOrders(OrdersToCom, ComToOrders,OrdersToElevAlgo,ElevAlgoToOrders)
  go elevAlgo.ElevStateMachine(OrdersToElevAlgo, ElevAlgoToOrders,
      ComToElevAlgo, ElevAlgoToCom)
  go communication.InitCom(ComToElevAlgo, ComToOrders, ElevAlgoToCom, OrdersToCom)

  for{}

	//done
}
