package main

import (
	//"./network/localip"
	//"./network/peers"
	//"flag"
	"fmt"
	"os"
	//"time"	
	//"math/rand"	
    //"strconv"
    com "github.com/chrskj/TTK4145-gruppe44/code/communication"
    //. "github.com/chrskj/TTK4145-gruppe44/code/util"
)

func main() {
    fmt.Println("Started")

    id := fmt.Sprintf("%d", os.Getpid())

    //elevAlgoToOrders := make(chan int)
    //ordersToElevAlgo := make(chan int)

    //comToElevAlgo := make(chan int)
    //elevAlgoToCom := make(chan int)

    //ordersToCom := make(chan struct med noe)
    //comToOrders := make(chan struct med noe)

    go com.SendHeartbeat(id)
    go com.ReceiveHeartbeat()
    //go com.SendMessage(id)
    go com.ReceiveMessage()

    //go com.ListenForModules(elevAlgoToCom, ordersToCom, comToElevAlgo,
    //    comToOrders)
    //go com.ListenForModules(elevAlgoToCom)

    /*
    //test-func
    go func() {
        for i := 0; i < 10; i++ {
			elevAlgoToCom<-i
        }
    }()
    */

    for{}
}
