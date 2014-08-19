package nserv_test

/*
Tests internals of nserv.
*/

import (
	"gopkg.in/kornel661/nserv.v0"
	"io/ioutil"
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

func init() {
	// enable parallelism (if your hardware supports it)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func serverTest(t *testing.T, srv *nserv.Server, handler http.HandlerFunc) {
	//srv := newServer()
	srv.ReadTimeout = 1 * time.Second
	srv.WriteTimeout = 1 * time.Second
	max := 10
	// finish is closed when server finishes
	finish := make(chan struct{})
	// counter to chack if throttling limit is obeyed
	counter := make(chan struct{}, max)
	for i := 0; i < max; i++ {
		counter <- struct{}{}
	}
	// count via ConnState
	srv.ConnState = func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			select {
			case <-counter:
				if len(counter) == 0 {
					t.Log("server counter: got 0 -- good")
				}
			default: // no tokens in the counter, the limit's been exceeded
				t.Error("Exceeded limit of simultaneous connections.")
			}
			runtime.Gosched() // give other connections chance to connect
		case http.StateClosed:
			counter <- struct{}{}
		}
	}
	http.HandleFunc("/", handler)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			t.Error(err)
		}
		close(finish)
	}()
	// test throttling
	t.Log("Testing throttling...")
	srv.MaxConns(max)
	path := "/test"
	// getFinished: when a "get" function finishes it puts a token here
	getFinished := make(chan struct{}, 10*max)
	get := func() {
		if resp, err := http.Get("http://" + addr + path); err != nil {
			t.Error(err)
		} else {
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Error(err)
			} else {
				if string(body) != path {
					t.Errorf("Got message `%s`.", body)
				}
			}
			resp.Body.Close()
		}
		getFinished <- struct{}{}
	}
	for i := 0; i < 10*max; i++ {
		go get()
	}
	for i := 0; i < 10*max; i++ {
		<-getFinished
	}
	// test gracefull shutdown
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
		srv.Stop()             // signal server to stop
		runtime.Gosched()      // give the server a chance to exit ungracefully
		time.Sleep(delay * 10) // give the server a chance to exit ungracefully
		close(clientDone)      // signal we're about to exit
	}()
	// let's see who exits first
	select {
	case <-clientDone:
		// client exited first, test passed
		<-finish
	case <-finish:
		t.Error("Server (most probably) exited ungracefully.")
		<-clientDone
	}
}
