package main

import (
	"fmt"
	"log"
	"time"
)

func slowProcess() error {
	for i := 0; i < 3; i++ {
		log.Println("doing something...", i)
		time.Sleep(time.Duration(1 * time.Second))
	}
	log.Println("something is done")
	return nil
}

func handle() {
	resultCh := make(chan error, 1)
	go func() { resultCh <- slowProcess() }()

	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("slowProcess Timedout.")
	case err := <-resultCh:
		fmt.Println("Result:", err)
	}
}

func main() {
	handle()
}
