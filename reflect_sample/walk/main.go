package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

func parseJson(input io.Reader) (interface{}, error) {
	var result interface{}
	if err := json.NewDecoder(input).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func makeIndent(depth int) (indent string) {
	if depth > 0 {
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		return
	}
	return
}

func printValue(value interface{}, depth int, isSlice bool) {
	if isSlice {
		fmt.Print(makeIndent(depth))
	}

	typ := reflect.TypeOf(value)
	switch typ.Kind() {
	case reflect.Int:
		fmt.Println("value(int):", value)
	case reflect.Float64:
		fmt.Println("value(float64):", value)
	case reflect.Slice:
		fmt.Println("value(slice)")
		v := reflect.Indirect(reflect.ValueOf(value))
		if v.Len() < 1 {
			fmt.Println("value(slice): []")
		} else {
			for i := 0; i < v.Len(); i++ {
				sv := v.Index(i)
				printValue(sv.Interface(), depth+1, true)
			}
		}
	case reflect.Map:
		fmt.Println("{")
		walkJson(value, depth+1, false)
		fmt.Println(makeIndent(depth) + "}")
	case reflect.String:
		fmt.Println("value(string):", value)
	default:
		fmt.Println("value(other):", value)
	}
}

func walkJson(json interface{}, depth int, isSlice bool) {
	for key, value := range json.(map[string]interface{}) {
		fmt.Print(makeIndent(depth))
		fmt.Printf("key: %s, ", key)
		printValue(value, depth, false)
	}
}

func main() {
	f, err := os.Open("example1.json")
	if err != nil {
		log.Fatalf("reading file: %s", err)
	}
	defer f.Close()

	json, err := parseJson(f)
	if err != nil {
		log.Fatalf("parsing json: %s", err)
	}

	fmt.Println(reflect.TypeOf(json))
	walkJson(json, 0, false)
}
