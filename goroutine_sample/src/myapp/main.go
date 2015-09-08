package main

import (
	"fmt"
	"time"
)

// http://guzalexander.com/2013/12/06/golang-channels-tutorial.html

func printGoroutine() {
	go fmt.Println("printGoroutine goroutine")
	fmt.Println("printGoroutine function message")
}

func synchronousChannel() {
	done := make(chan bool)
	go func() {
		fmt.Println("synchronousChannel goroutine")
		done <- true
	}()

	fmt.Println("synchronousChannel function message")
	<-done
}

func asynchronousChannel() {
	message := make(chan string, 2) // buffered
	count := 3

	go func() {
		for i := 1; i <= count; i++ {
			fmt.Println("send message:", i)
			message <- fmt.Sprintf("message %d", i)
			fmt.Println("sent message:", i)
		}
	}()

	time.Sleep(time.Second * 3)

	for i := 1; i <= count; i++ {
		fmt.Println(<-message)
	}
}

// http://stackoverflow.com/questions/15715605/multiple-goroutines-listening-on-one-channel

func multipleGoroutine() {
	c := make(chan string)

	for i := 1; i <= 5; i++ {
		// not ensure the order of running goroutine
		go func(i int, co chan<- string) {
			for j := 1; j <= 5; j++ {
				co <- fmt.Sprintf("hi from %d.%d", i, j)
			}
		}(i, c)
	}

	for i := 1; i <= 25; i++ {
		fmt.Println(<-c)
	}
}

func main() {
	fmt.Println("Main start")

	// e.g. 1
	printGoroutine()
	synchronousChannel()
	asynchronousChannel()

	// e.g. 2
	multipleGoroutine()

	fmt.Println("Main end")
}
