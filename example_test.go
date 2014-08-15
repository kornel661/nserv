package nserv_test

import (
	"fmt"
	"gopkg.in/kornel661/nserv.v0"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func ExampleServer() {
	// initialize database, etc.
	defer func() {
		// do clean-up (close DB connections, etc.)
	}()
	// set-up server:
	srv := &nserv.Server
	srv.Addr = "localhost:12345"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	// catch signals:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		<-signals
		log.Println("Caught signal. Shutting down gracefully.")
		srv.Stop()
	}()
	// start serving:
	log.Printf("Serving at http://%s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
