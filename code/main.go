package main

import (
	"flag"
	"fmt"

	"./communication"
	"./elevAlgo"
	"./orders"
	. "./util"
)

func main() {
	fmt.Println("Started")

	OrdersToCom := make(chan ChannelPacket)
	ComToOrders := make(chan ChannelPacket)

	OrdersToElevAlgo := make(chan ChannelPacket)
	ElevAlgoToOrders := make(chan ChannelPacket)

	ComToElevAlgo := make(chan ChannelPacket)
	ElevAlgoToCom := make(chan ChannelPacket)

	var elevPort string
	flag.StringVar(&elevPort, "port", "15657", "Port of elevator to connect to")
	var elevID int
	flag.IntVar(&elevID, "id", -1, "Elevator ID")
	flag.Parse()

	go communication.InitCom(ComToElevAlgo, ComToOrders, ElevAlgoToCom,
		OrdersToCom, elevID)

	go elevAlgo.ElevStateMachine(ElevAlgoToOrders, ComToElevAlgo, ElevAlgoToCom,
		OrdersToElevAlgo, elevPort, elevID)

	go orders.InitOrders(OrdersToCom, ComToOrders, ElevAlgoToOrders,
		OrdersToElevAlgo, elevID)

	for {
	}
}
