package main

import (
	"flag"
	"fmt"
	"time"
)

type callReq struct {
	resp     chan error
	timeout  chan struct{}
	timer    *time.Timer
	duration time.Duration
}

var durationSecond = flag.Uint("duration", 0, "")
var reset = flag.Bool("reset", false, "")

func main() {
	flag.Parse()
	fmt.Println("start")
	fmt.Println(fmt.Sprintf("close: %v", *durationSecond))

	c := callReq{}
	c.timeout = make(chan struct{})
	defer close(c.timeout)
	c.duration = time.Duration(*durationSecond) * time.Second

	fmt.Println(c.duration)
	if c.duration > 0 {
		fmt.Println("duration > 0")
		c.timer = time.NewTimer(c.duration)
		fmt.Println("create NewTimer")
		time.Sleep(1 * time.Second)
		fmt.Println("slept, then timer block")
		<-c.timer.C
		fmt.Println("come back from timer")
	}

	if *reset {
		c.timer.Reset(c.duration)
	}

	time.Sleep(2 * time.Second)

	if c.timer != nil {
		fmt.Println("has timer")
		if !c.timer.Stop() {
			// timer is already expired, or stopped
			select {
			case <-c.timer.C:
				fmt.Println("drain if timer has already expired or been stopped")
			default:
				fmt.Println("maybe not call this sentence")
			}
		}

		fmt.Println("reset timer using c.duration")
		c.timer.Reset(c.duration)
	}

	fmt.Println("end")
}
