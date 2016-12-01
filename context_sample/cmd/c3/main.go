package main

import (
	"fmt"
	"log"
	"time"
)

func slowProcess(done <-chan struct{}) error {
	for i := 0; i < 10; i++ {
		log.Println("doing something...", i)
		select {
		case <-time.After(1 * time.Second):
		case <-done:
			log.Println("slowProcess done.", i)
			return nil
		}
	}
	log.Println("something is done")
	return nil
}

func handle() {
	done := make(chan struct{}, 1)
	resultCh := make(chan error, 1)
	go func() { resultCh <- slowProcess(done) }()

	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		close(done)
		fmt.Println("slowProcess Timedout.")
	case err := <-resultCh:
		fmt.Println("Result:", err)
	}
}

func main() {
	handle()
}
