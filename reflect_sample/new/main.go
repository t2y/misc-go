package main

import (
	"fmt"
	"reflect"
)

func newString() {
	var s string
	rs := reflect.New(reflect.TypeOf(s)).Elem()
	rs.SetString("test")
	fmt.Println("new and set:", rs)
}

func main() {
	newString()
}
