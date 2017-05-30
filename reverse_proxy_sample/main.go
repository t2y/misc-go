package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

func runInsideOfLocalhost() {
	log.Println("---> start runInsideOfLocalhost")

	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	defer backendServer.Close()

	log.Println("backend server:", backendServer.URL)

	rpURL, err := url.Parse(backendServer.URL)
	if err != nil {
		log.Fatal(err)
	}
	frontendProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(rpURL))
	defer frontendProxy.Close()

	log.Println("frontend proxy:", frontendProxy.URL)

	resp, err := http.Get(frontendProxy.URL)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", b)
	log.Println("<--- end runInsideOfLocalhost")
}

var localServer = "localhost:5000"
var backendServer string

func init() {
	flag.StringVar(&backendServer, "backendServer", "", "set target server. e.g.) http://hostname:8000")
}

func runOutsideOfLocalhost() {
	log.Println("---> start runOutsideOfLocalhost")
	flag.Parse()
	if backendServer == "" {
		log.Println("no backend server")
		return
	}

	log.Println("backend server:", backendServer)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(backendServer)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("proxy request:")
		log.Println("header:", r.Header)
		log.Println("host:", r.Host)
		log.Println("uri:", r.RequestURI)

		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(w, r)
	})

	log.Println(".... serve as reverse proxy:", localServer)
	err := http.ListenAndServe(localServer, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println("<--- end runOutsideOfLocalhost")
}

func main() {
	runInsideOfLocalhost()
	runOutsideOfLocalhost()
}
