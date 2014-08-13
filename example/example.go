package main

import (
	"fmt"
	"gopkg.in/kornel661/nserv.v0"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	// initialize database, etc.
	// initialize()
	// do clean-up
	defer func() {
		// cleanup()
	}()
	// set-up server
	srv := nserv.New(nil, 100)
	srv.Addr = "localhost:12345"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	// catch signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		<-signals
		log.Println("Caught signal. Shutting down gracefully.")
		srv.Stop()
	}()
	// start serving
	log.Println("Serving at http://localhost:12345/")
	srv.ListenAndServe()
	//http.ListenAndServe("localhost:12345", nil)
}
