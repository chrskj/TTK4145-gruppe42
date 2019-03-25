package main

import (
    "fmt"
    "flag"
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

    go orders.InitOrders(OrdersToCom, ComToOrders,OrdersToElevAlgo,
            ElevAlgoToOrders)

    var elevPort string
    flag.StringVar(&elevPort, "port", "15657", "Port of elevator to connect to")
    flag.Parse()

    go elevAlgo.ElevStateMachine(OrdersToElevAlgo, ElevAlgoToOrders,
            ComToElevAlgo, ElevAlgoToCom, elevPort)

    go communication.InitCom(ComToElevAlgo, ComToOrders, ElevAlgoToCom,
            OrdersToCom)

    for{}

	//done
}
