package nserv_test

import (
	"gopkg.in/kornel661/nserv.v0"
	"testing"
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
