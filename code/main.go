package main

import (
	//"./network/localip"
	//"./network/peers"
	//"flag"
	//"fmt"
	//"os"
	//"time"	
	//"math/rand"	
    //"strconv"
    com "code/communication"
)

//Starte alle kanaler
//Starte alle gorutines, og passe kanaler som input arguments
/*
//Kanal orders -> komm (orders)
ordersToCom := make(chan struct med noe)

//Kanal komm -> orders (orders)
comToOrders := make(chan struct med noe)


//Kanal orders -> heisalgo (ønsket floor)
ordersToElevAlgo := make(chan int)
//Kanal heisalgo -> orders (current floor)
elevAlgoToOrders := make(chan int)

//Kanal komm -> heisalgo (request om cost function)
comToElevAlgo := make(chan int)
//Kanal heisalgo -> komm (cost function)
elevAlgoToCom := make(chan float)

go orders(ordersToCom, comToOrders)
go elevAlgo(ordersToElevAlgo,elevAlgoToOrders)
go com(comToElevAlgo,elevAlgoToCom)
*/
//done

func main() {
 	fmt.Println("Started")
    go com.sendHeartbeat()
    go com.listenHeartbeat()
}
