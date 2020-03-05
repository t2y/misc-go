package main

import (
	"log"
)

func check(v interface{}) {
	if c, ok := (v).(string); ok {
		log.Println("converted: %+v", c)
	} else {
		log.Println("cannot cast")
	}
}

func main() {
	s := "test"
	check(s)

	var v interface{} = nil
	check(v)
}
