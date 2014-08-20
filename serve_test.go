package nserv_test

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/kornel661/nserv.v0"
	"html"
	"io/ioutil"
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

func getFunc(t *testing.T, path string) {
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
}

func getTLSFunc(t *testing.T, path string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if resp, err := client.Get("https://" + addr + path); err != nil {
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
}

func TestThrottling(t *testing.T) {
	srv := newServer()
	serverTest(t, srv, handler, (*nserv.Server).ListenAndServe, getFunc)
}

func TestThrottlingTLS(t *testing.T) {
	srv := newServer()
	LAS := func(srv *nserv.Server) error {
		return srv.ListenAndServeTLS("_test.crt", "_test.key")
	}
	serverTest(t, srv, handler, LAS, getTLSFunc)
}
