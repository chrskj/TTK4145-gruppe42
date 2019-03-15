package communication

//sende heartbeats
//phoenix backup
// Hele tiden oppdatere de andre heisene på sine egne orders
// - hele tiden sende ned til orders en komplett ordreliste
// - ta imot orders sin ordrelsite og sende ut

import (
    "github.com/chrskj/TTK4145-gruppe44/code/network/bcast"
	//"./network/localip"
	//"./network/peers"
	//"flag"
	"fmt"
	//"os"
	"time"
	"math/rand"
    //"strconv"
)

type MessageStruct struct {
	Message string
}

func SendHeartbeat() {
    rand.Seed(time.Now().UnixNano())
    ranInt := rand.Intn(20)
    transmitHeartbeat := make(chan MessageStruct)
    go bcast.Transmitter(16569, transmitHeartbeat)
    response := fmt.Sprintf("Heartbeat from %d", ranInt)
    helloMsg := MessageStruct{response}
    for {
        transmitHeartbeat <- helloMsg
        time.Sleep(1 * time.Second)
    }
}

func ListenHeartbeat() {
    receiveHeartbeat := make(chan MessageStruct)
    go bcast.Receiver(16569, receiveHeartbeat)
    for {
        select {
        case a := <-receiveHeartbeat:
            fmt.Printf("Received: %v\n", a)
        }
    }
}

func SendMessage() {
    transmitMessage := make(chan MessageStruct)
    go bcast.Transmitter(16570, transmitMessage)
    response := fmt.Sprintf("Heartbeat from %d", ranInt)
    helloMsg := MessageStruct{response}
    for {
        transmitMessage<- helloMsg
        time.Sleep(1 * time.Second)
    }
}

func ListenMessage(addresse, beskjed) {
    receiveMessage := make(chan MessageStruct)
    go bcast.Receiver(16570, receiveMessage)
    for {
        select {
        case a := <-receiveMessage:
            fmt.Printf("Received: %v\n", a)
        }
    }
}






