package main

import (
    . "fmt"
    "runtime"
    "time"
)

var i = 0

func incrementing() {
    for j := 0; j < 1000000; j++ {
        i += 1
    }
}

func decrementing() {
    for j := 0; j < 1000000; j++ {
        i -= 1
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    go incrementing()
    go decrementing()
    
    // We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
    // We'll come back to using channels in Exercise 2. For now: Sleep.
    time.Sleep(100*time.Millisecond)
    Println("The magic number is:", i)
}
