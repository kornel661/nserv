package nserv_test

import (
	"gopkg.in/kornel661/nserv.v0"
	"net/http"
	"testing"
)

func TestCanResume(t *testing.T) {
	srv := newServer()
	go func() {
		srv.Stop()
	}()
	if nserv.CanResume() {
		if err := srv.ResumeAndServe(); err != nil {
			t.Errorf("srv.ResumeAndServe error: %v", err)
		}
	} else { // can't resume
		if err := srv.ResumeAndServe(); err == nil {
			t.Errorf("srv.ResumeAndServe error: resumed when it wasn't possible")
		}
	}
}

func TestCopyListenerFD(t *testing.T) {
	srv := newServer()
	go func() {
		defer func() {
			srv.Stop()
			// FIXME: can make it reliable some other way?
			resp, _ := http.Get("http://" + srv.Addr + "/")
			// without http.Get srv.Stop above is unreliable
			// (server just waits for a single connection sometimes)
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}

		}()
		fd, err := srv.CopyListenerFD()
		if err != nil {
			t.Errorf("srv.CopyListenerFD error: %v", err)
			return
		}
		if err := fd.Close(); err != nil {
			t.Errorf("fd.Close error: %v", err)
		}
	}()
	if err := srv.ListenAndServe(); err != nil {
		t.Errorf("srv.ListenAndServe error: %v", err)
	}
}
