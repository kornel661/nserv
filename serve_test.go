package nserv_test

import (
	"gopkg.in/kornel661/nserv.v0"
	"net/http"
	"testing"
	"time"
)

const (
	deadlockDelay = 1 * time.Second
	deadlockTest  = time.Second / 2
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
	case <-srv.Stop():
		t.Error("srv.Stop() returned before the server shut down.")
	case <-time.After(deadlockTest): // OK
	}
}

func TestServerStartStop1(t *testing.T) {
	srv := nserv.New(opts, 10)
	go srv.Stop()
	t.Log("starting server, it should terminate almost instantaneously")
	if err := srv.ListenAndServe(); err != nil {
		t.Error(err)
	}
}

func TestServerStartStop2(t *testing.T) {
	srv := nserv.New(opts, 10)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			t.Error(err)
		}
	}()
	select {
	case <-srv.Stop(): // shouldn't deadlock
	case <-time.After(deadlockDelay):
		t.Error("deadlock")
	}
}
