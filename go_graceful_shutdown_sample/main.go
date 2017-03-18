package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var port = flag.String("port", "8000", "")
var shutdownTimeout = flag.Duration("shutdownTimeout", time.Second*3, "")
var mustClose = flag.Bool("mustClose", false, "")

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("hello called")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("hello\n"))
	time.Sleep(5 * time.Second)
	log.Println("write again")
	w.Write([]byte("hello again\n"))
	time.Sleep(5 * time.Second)
	log.Println("write bye")
	w.Write([]byte("bye\n"))
}

type Server struct {
	Port            string
	ShutdownTimeout *time.Duration
	MustClose       bool
}

func (s *Server) Serve(sigChan chan os.Signal) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	httpServer := &http.Server{
		Addr:    ":" + s.Port,
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("ListenAndServe returns an error", err)
			if err != http.ErrServerClosed {
				log.Fatalln("HTTPServer closed with error:", err)
			}
		}
		log.Println("ListenAndServe goroutine completed.")
	}()

	log.Println("Starting http server...")
	log.Printf("SIGNAL %d received, then shutting down...\n", <-sigChan)

	ctx, cancel := context.WithTimeout(context.Background(), *s.ShutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("Failed to gracefully shutdown HTTPServer:", err)
		if s.MustClose {
			httpServer.Close()
			log.Println("Server closed immediately")
		}
	}
	log.Println("HTTPServer shutdown.")
}

func NewServer(
	port string, shutdownTimeout *time.Duration, mustClose bool,
) (s *Server) {
	s = &Server{
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
		MustClose:       mustClose,
	}
	log.Println("Server:", s)
	return
}

func main() {
	flag.Parse()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt) // catch SIGTERM and SIGINT

	log.Println("My process id is", os.Getpid())
	s := NewServer(*port, shutdownTimeout, *mustClose)
	s.Serve(sigCh)

	log.Printf("Server shutdown, but waiting 5 second to exit process ...")
	time.Sleep(5 * time.Second)
}
