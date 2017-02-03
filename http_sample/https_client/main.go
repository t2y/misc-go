package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var server = flag.String("server", "localhost", "")
var port = flag.String("port", "4443", "")
var path = flag.String("path", "/", "")

var insecure = flag.Bool("insecure", false, "")

func main() {
	flag.Parse()

	c := new(tls.Config)
	if *insecure {
		c.InsecureSkipVerify = true
	}

	tr := &http.Transport{TLSClientConfig: c}
	client := &http.Client{Transport: tr}
	res, err := client.Get("https://" + *server + ":" + *port + *path)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	log.Println(res)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", b)
}
