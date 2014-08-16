package nserv_test

import (
	"fmt"
	"gopkg.in/kornel661/nserv.v0"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"testing"
	"time"
)

func newServer() *nserv.Server {
	srv := &nserv.Server{}
	srv.Addr = addr
	return srv
}

func TestServerStop1(t *testing.T) {
	srv := newServer()
	finish := make(chan struct{})
	go func() {
		if !srv.Stop() {
			t.Error("We didn't stop the server.")
		}
		close(finish)
	}()
	if err := srv.ListenAndServe(); err != nil {
		t.Error(err)
	}
	<-finish
	srv.Wait()
	if srv.Stop() {
		t.Error("We did stop the server.")
	}
}

// just to test possibly different goroutine ordering
func TestServerStop2(t *testing.T) {
	srv := newServer()
	finish := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			t.Error(err)
		}
		close(finish)
	}()
	if !srv.Stop() {
		t.Error("We didn't stop the server.")
	}
	<-finish
	if srv.Stop() {
		t.Error("We did stop the server.")
	}
}

func TestThrottling(t *testing.T) {
	srv := newServer()
	srv.ReadTimeout = 1 * time.Second
	srv.WriteTimeout = 1 * time.Second
	max := 10
	finish := make(chan struct{})
	counter := make(chan struct{}, max)
	for i := 0; i < max; i++ {
		counter <- struct{}{}
	}
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", html.EscapeString(r.URL.Path))
	})
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
		t.Error("Server exited ungracefully.")
		<-clientDone
	}
}
