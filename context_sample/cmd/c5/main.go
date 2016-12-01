package main

import (
	"context"
	"fmt"
)

func handle() {
	key1 := "key1"
	value1 := "myValue1"
	ctx := context.WithValue(context.Background(), key1, value1)

	key2 := "key2"
	value2 := 2
	ctx = context.WithValue(ctx, key2, value2)

	fmt.Println("Result1:", ctx.Value(key1).(string))
	fmt.Println("Result2:", ctx.Value(key2).(int))
}

func main() {
	handle()
}
