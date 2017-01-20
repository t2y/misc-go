package main

import (
	"flag"
	"fmt"
	"time"
)

// async channel sender sample
// http://guzalexander.com/2013/12/06/golang-channels-tutorial.html

var sleepSecond = flag.Uint("sleep", 0, "")

func main() {
	flag.Parse()
	fmt.Println("start")
	fmt.Println(fmt.Sprintf("sleep: %v", *sleepSecond))

	c1 := make(chan struct{})
	c2 := make(chan int, 1)

	go func() {
		// async channel sender
		fmt.Println("sender goroutine start")

		if *sleepSecond > 0 {
			time.Sleep(time.Duration(*sleepSecond) * time.Second)
		}

		select {
		case <-c1:
			fmt.Println("receive empty struct")
		case c2 <- 10:
			fmt.Println("send integer")
		default:
			fmt.Println("default")
		}
		fmt.Println("sender goroutine end")
		return
	}()

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- struct{}{}
		fmt.Println("send empty struct")
	}()

	go func() {
		fmt.Println("receiver goroutine start")
		select {
		case i := <-c2:
			fmt.Println("got integer from channel", i)
		}
		fmt.Println("receiver goroutine end")
	}()

	time.Sleep(2 * time.Second)
	// close(c1)
	close(c2)

	time.Sleep(2 * time.Second)
	fmt.Println("end")
}
