import "./elevAlgo"
import "fmt"

type Order struct {
	elevator  int
	toFloor   int64
	direction bool
	timestamp uint64
}

type ChannelPacket struct {
	packetType string
	elevator   int
	toFloor    int64
	direction  bool //True = opp, false = ned
	timestamp  uint64
	cost       float64
	dataJson   []byte
}

//Kanal orders -> komm (orders)
ordersToCom := make(chan struct ChannelPacket)

//Kanal komm -> orders (orders)
comToOrders := make(chan struct ChannelPacket)


//Kanal orders -> heisalgo (ønsket floor, direction)
ordersToElevAlgo := make(chan ChannelPacket)
//Kanal heisalgo -> orders (current floor)
elevAlgoToOrders := make(chan ChannelPacket)

//Kanal komm -> heisalgo (request om cost function)
comToElevAlgo := make(chan orderStruct)
//Kanal heisalgo -> komm (cost function)
costFuncToCom := make(chan float)
//Kanal heisalgo -> komm (ny ordre )
newOrderToCom := make(chan orderStruct)

go orders.initialize(ordersToCom, comToOrders,ordersToElevAlgo,elevAlgoToOrders)
go elevAlgo.main(ordersToElevAlgo,elevAlgoToOrders, comToElevAlgo,costFuncToCom,newOrderToCom)
go com(comToElevAlgo,elevAlgoToCom)

//done

