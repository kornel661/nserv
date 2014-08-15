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

func TestServerStop(t *testing.T) {
	srv := newServer()
	go srv.Stop()
	srv.ListenAndServe()
}
