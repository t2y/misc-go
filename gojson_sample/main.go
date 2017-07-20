package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/gojson"
)

func main() {
	f, err := os.Open("example1.json")
	if err != nil {
		log.Fatalf("reading file: %s", err)
	}
	defer f.Close()

	tagList := []string{"json"}
	subStruct := true
	if output, err := gojson.Generate(
		f, gojson.ParseJson, "example", "mypkg", tagList, subStruct,
	); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(output))
	}
}
