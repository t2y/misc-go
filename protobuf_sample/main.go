package main

import (
	"log"
	"time"

	"example"

	"github.com/golang/protobuf/proto"
)

func createData(
	label string, typ int32, reps []int64, group *example.Test_OptionalGroup,
) ([]byte, error) {
	data := new(example.Test)
	data.Label = proto.String(label)
	data.Type = proto.Int32(typ)
	data.Reps = reps
	data.Optionalgroup = group
	return proto.Marshal(data)
}

func main() {
	go func() {
		var i int32
		for {
			time.Sleep(1 * time.Second)

			data, err := createData(
				"hello world",
				i,
				[]int64{100, int64(100 + i)},
				&example.Test_OptionalGroup{
					RequiredField: proto.String("good bye"),
				},
			)
			if err != nil {
				log.Fatal("marshaling error: ", err)
			}

			sendData(data)
			i++
		}
	}()

	runServer()
}
