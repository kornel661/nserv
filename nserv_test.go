package nserv_test

import (
	"gopkg.in/kornel661/nserv.v0"
	"net"
	"net/http"
	"runtime"
	"testing"
	"time"
)

const (
	// time to wait (let other goroutines do their job first)
	delay = 40 * time.Millisecond
	// address to listen to
	addr = "localhost:1234"
)

// serverTest tests server srv, checks throttling and graceful exit
// LAS is a ListenAndServe-type function applied to srv
// getFunc should get addr/path
func serverTest(t *testing.T, srv *nserv.Server, handler http.HandlerFunc,
	LAS func(srv *nserv.Server) error, getFunc func(t *testing.T, path string)) {

	defer func() { // reset the default muxer
		http.DefaultServeMux = http.NewServeMux()
	}()
	srv.ReadTimeout = 2 * time.Second
	srv.WriteTimeout = 2 * time.Second
	max := 10
	m := 6 // number of conns = m * max
	if testing.Short() {
		m = 2
	}
	// finish is closed when server finishes
	finish := make(chan struct{})
	// counter to check if throttling limit is obeyed
	counter := make(chan struct{}, max)
	for i := 0; i < max; i++ {
		counter <- struct{}{}
	}
	// count active connections via ConnState
	srv.ConnState = func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			// give other connections chance to return their tokens
			runtime.Gosched()
			time.Sleep(delay / 2)
			select {
			case <-counter:
				// took a token
				if len(counter) == 0 {
					t.Log("server counter: got 0 -- good")
				}
			default:
				// no tokens in the counter, the limit's been exceeded
				t.Error("Possible error. Exceeded(?) limit of simultaneous connections.")
				// take the token we missed
				<-counter
			}
		case http.StateClosed, http.StateHijacked:
			// return the token (before or after? the connection is closed)
			// FIXME: needs to be before, otherwise might get false positives
			//        errors of exceeding the limit
			counter <- struct{}{}
		}
	}
	// setup handler and start the server
	http.HandleFunc("/", handler)
	go func() {
		if err := LAS(srv); err != nil {
			t.Error(err)
		}
		close(finish)
	}()
	runtime.Gosched()
	time.Sleep(delay)

	// test throttling
	t.Log("Testing throttling...")
	srv.MaxConns(max)
	path := "/test"
	// getFinished: when a "get" function finishes it puts a token here
	getFinished := make(chan struct{}, m*max)
	get := func() {
		getFunc(t, path)
		getFinished <- struct{}{}
	}
	for i := 0; i < m*max; i++ {
		go get()
	}
	for i := 0; i < m*max; i++ {
		<-getFinished
	}

	// test graceful shutdown
	t.Log("Testing graceful shutdown...")
	clientDone := make(chan struct{})
	// client
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		srv.Stop()            // signal server to stop
		runtime.Gosched()     // give the server a chance to exit ungracefully
		time.Sleep(delay * 6) // give the server a chance to exit ungracefully
		close(clientDone)     // signal we're about to exit
	}()
	// let's see who exits first
	select {
	case <-clientDone:
		// client exited first, test passed
		t.Log("Client exited first, the server is graceful.")
		<-finish
	case <-finish:
		t.Error("Error. Server (most probably) exited ungracefully.")
		<-clientDone
	}

	runtime.Gosched()
	time.Sleep(delay)
	if n := len(counter); n != cap(counter) {
		t.Errorf("Error. The number of tokens in the counter is: %d.\n", n)
	}
}
