package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func slowProcess(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		log.Println("doing something...", i)
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			log.Println("slowProcess done.", i)
			return ctx.Err()
		}
	}
	log.Println("something is done")
	return nil
}

func handle() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resultCh := make(chan error, 1)
	go func() { resultCh <- slowProcess(ctx) }()

	err := <-resultCh
	fmt.Println("Result:", err)
}

func main() {
	handle()
}
