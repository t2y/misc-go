package main

import (
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"

	"example"
)

const (
	tcp  = "tcp"
	host = "127.0.0.1"
	port = 2110
)

var uri string

func init() {
	uri = fmt.Sprintf("%s:%d", host, port)
}

func dumpData(data *example.Test) {
	log.Printf("dump data\n\tlabel: %s\n\ttype: %d\n\treps: %v\n\tgroup: %v\n",
		data.GetLabel(),
		data.GetType(),
		data.GetReps(),
		data.GetOptionalgroup(),
	)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("read error: ", err)
	}
	log.Printf("received data: %d bytes", n)

	data := new(example.Test)
	err = proto.Unmarshal(buf[0:n], data)
	if err != nil {
		log.Fatal("unmarshal error: ", err)
	}

	if data == nil {
		log.Println("data is null")
		return
	}

	dumpData(data)
}

func runServer() {
	listener, err := net.Listen(tcp, uri)
	if err != nil {
		log.Fatal("listener error: ", err)
	}

	for {
		log.Println("can accept ...")
		if conn, err := listener.Accept(); err == nil {
			handleConnection(conn)
		} else {
			continue
		}
	}
}
