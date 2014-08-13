package nserv_test

import (
	"fmt"
	"gopkg.in/kornel661/nserv.v0"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	deadlockDelay = time.Second / 2
	deadlockTest  = time.Second / 4
	waitDelay     = 30 * time.Millisecond
	addr          = "localhost:1234"
)

var (
	opts = &http.Server{Addr: addr}
)

func TestDoubleInitialize(t *testing.T) {
	srv := nserv.New(opts, 0)
	defer func() {
		if err := recover(); err == nil {
			t.Error("Second initialization didn't panic.")
		}
	}()
	srv.Initialize(10, 5)
}

func TestInitializeNegative(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Initialization with negative throttling limit didn't panic.")
		}
	}()
	nserv.New(nil, -1)
}

func TestServerStartStop0(t *testing.T) {
	srv := nserv.New(nil, 10)
	select {
	case <-srv.StopChan():
		t.Error("srv.Stop() returned before the server shut down.")
	case <-time.After(deadlockTest): // OK
		t.Log("Waited deadlockTest seconds and nothing. Good.")
	}
}

func TestServerStartStop1(t *testing.T) {
	srv := nserv.New(opts, 10)
	go srv.Stop()
	t.Log("starting server, it should terminate almost instantaneously, without reporting any errors")
	if err := srv.ListenAndServe(); err != nil {
		t.Fatal(err)
	}
	t.Log("Waiting for server to shutdown.")
	srv.StopWait()
	srv.StopWait()
	srv.StopWait()
	srv.StopWait()
	srv.StopWait()
}

func TestServerStartStop2(t *testing.T) {
	srv := nserv.New(opts, 10)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			t.Error(err)
		}
	}()
	select {
	case <-srv.StopChan(): // shouldn't deadlock
	case <-time.After(deadlockDelay):
		t.Error("Waited deadlockDelay seconds. Deadlock?")
	}
}

func TestServingBasic(t *testing.T) {
	hw := "Hello World!"
	log.Println("New server...")
	srv := nserv.New(nil, 10)
	srv.Addr = addr
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, hw)
	})
	log.Println("Starting server...")
	finished := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Error listening: %s\n", err)
			t.Fatal(err)
		}
		finished <- struct{}{}
	}()
	var (
		resp *http.Response
		body []byte
		err  error
	)
	time.Sleep(waitDelay)
	log.Println("Getting a response...")
	if resp, err = http.Get("http://" + addr); err != nil {
		t.Fatal(err)
	}
	log.Println("Reading the response...")
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
	if string(body) != hw {
		t.Errorf("Got body `%s`, should be `%s`", body, hw)
	}
	log.Println("Stopping server.")
	srv.Stop()
	log.Println("Waiting for the server to stop.")
	<-finished
}
