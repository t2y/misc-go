package main

import (
	"flag"
	"fmt"
	"time"
)

type callReq struct {
	resp    chan error
	timeout chan struct{}
	timer   *time.Timer
}

var closeSecond = flag.Uint("close", 1, "")

func main() {
	flag.Parse()
	fmt.Println("start")
	fmt.Println(fmt.Sprintf("close: %v", *closeSecond))

	c := callReq{}
	c.timeout = make(chan struct{})

	go func() {
		time.Sleep(3 * time.Second)
		c.timeout <- struct{}{}
		fmt.Println("send empty data to channel")
	}()

	go func() {
		time.Sleep(time.Duration(*closeSecond) * time.Second)
		close(c.timeout)
		fmt.Println("close channel")
	}()

	fmt.Println("block main")
	<-c.timeout
	fmt.Println("come back main")

	time.Sleep(1 * time.Second)
	fmt.Println("end")
}
