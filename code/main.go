package main

import (
	"flag"
	"fmt"

	"./communication"
	"./elevalgo"
	"./orders"
	"./util"
)

func main() {
	fmt.Println("Started")

	OrdersToCom := make(chan util.ChannelPacket)
	ComToOrders := make(chan util.ChannelPacket)

	OrdersToElevAlgo := make(chan util.ChannelPacket)
	ElevAlgoToOrders := make(chan util.ChannelPacket)

	ComToElevAlgo := make(chan util.ChannelPacket)
	ElevAlgoToCom := make(chan util.ChannelPacket)

	var elevPort string
	flag.StringVar(&elevPort, "port", "15657", "Port of elevator to connect to")
	var elevID int
	flag.IntVar(&elevID, "id", -1, "Elevator ID")
	flag.Parse()

	go communication.InitCom(ComToElevAlgo, ComToOrders, ElevAlgoToCom,
		OrdersToCom, elevID)

	go elevalgo.ElevStateMachine(ElevAlgoToOrders, ComToElevAlgo, ElevAlgoToCom,
		OrdersToElevAlgo, elevPort, elevID)

	go orders.InitOrders(OrdersToCom, ComToOrders, ElevAlgoToOrders,
		OrdersToElevAlgo, elevID)

	for {
	}
}
