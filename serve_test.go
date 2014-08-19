package nserv_test

import (
	"fmt"
	"gopkg.in/kornel661/nserv.v0"
	"html"
	"net/http"
	"testing"
)

func TestServerStop(t *testing.T) {
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

func TestServerStopTLS(t *testing.T) {
	srv := newServer()
	finish := make(chan struct{})
	go func() {
		if err := srv.ListenAndServeTLS("_test.crt", "_test.key"); err != nil {
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

func newServer() *nserv.Server {
	srv := &nserv.Server{}
	srv.Addr = addr
	return srv
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", html.EscapeString(r.URL.Path))
}

func TestThrottling(t *testing.T) {
	srv := newServer()
	serverTest(t, srv, handler)
}
