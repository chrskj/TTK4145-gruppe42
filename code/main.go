package main

import (
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
}
