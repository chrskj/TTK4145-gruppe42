package main

import (
	"fmt"
    "time"
    w "github.com/chrskj/TTK4145-gruppe44/code/watchdog"
)

func main() {
	tick := time.Tick(100 * time.Millisecond)
    wdog := w.New(time.Second)
    wdog.Reset()
    for {
        select {
		case <-tick:
			fmt.Println("tick.")
        case <-wdog.TimeOverChannel():
            fmt.Println("Time over")
            return
        }
    }
}