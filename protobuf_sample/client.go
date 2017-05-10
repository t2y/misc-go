package main

import (
	"log"
	"net"
)

func sendData(data []byte) {
	conn, err := net.Dial(tcp, uri)
	if err != nil {
		log.Fatal(err)
	}

	n, err := conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sent data: %d bytes", n)
}
