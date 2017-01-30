package main

import (
	"log"
	"net/http"
)

// reference:
// http://qiita.com/ryurock/items/f55db5944397619735bf

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello\n"))
}

func main() {
	http.HandleFunc("/hello", hello)

	err := http.ListenAndServeTLS(":4443", "ssl/development/myself.crt", "ssl/development/myself.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
