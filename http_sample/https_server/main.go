package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

// reference:
// http://qiita.com/ryurock/items/f55db5944397619735bf

var MemProfileRate int = 1

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("hello called")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("hello\n"))
}

func main() {
	http.HandleFunc("/hello", hello)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err := http.ListenAndServeTLS(":4443", "ssl/development/myself.crt", "ssl/development/myself.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
