package nserv

import (
	"gopkg.in/kornel661/limitnet.v1"
	"net"
	"net/http"
	"strings"
	"sync"
)

// Server with graceful exit and throttling.
//
// Server is an extension of http.Server from the standard library (its API is
// a superset of that of http.Server).
type Server struct {
	http.Server                                     // standard net.Server functionality
	InitialMaxConns int                             // initial limit on simultaneous connections
	tlist           chan limitnet.ThrottledListener // list for Close()
	twlist          chan limitnet.ThrottledListener // list for Wait()
	initOnce        sync.Once                       // for initialization
}

// initialize initializes the server.
func (srv *Server) initialize() {
	srv.initOnce.Do(func() {
		srv.tlist = make(chan limitnet.ThrottledListener, 1)
		srv.twlist = make(chan limitnet.ThrottledListener, 1)
	})
}

// Serve accepts incoming connections on the Listener listn (wrapped with
// ThrottledListener from the gopkg.in/kornel661/limitnet.v1 package), creating
// a new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
// Don't close listn. Rather use srv.Stop() method to exit gracefully.
// Serve returns on unrecoverable errors and when the server is explicitly
// stopped by srv.Stop(). By the time Serve returns the listener listn is closed.
func (srv *Server) Serve(listn net.Listener) error {
	srv.initialize()
	l := limitnet.NewThrottledListener(listn)
	l.MaxConns(srv.InitialMaxConns)
	srv.tlist <- l
	err := srv.Server.Serve(l)
	stopped := !srv.Stop()
	if strings.Contains(err.Error(), "use of closed network connection") && stopped {
		err = nil // server's been stopped by the user (most probably)
	}
	srv.Wait()
	return err
}

// Wait returns only when the server is closed and all connections terminated.
func (srv *Server) Wait() {
	srv.initialize()
	if tl, ok := <-srv.twlist; ok {
		tl.Wait()
		close(srv.twlist)
	}
}

// Stop gracefully stops a running server. Returns false if server had already
// been stopped before.
func (srv *Server) Stop() bool {
	srv.SetKeepAlivesEnabled(false) // do it early (as if it matters)
	srv.initialize()
	if tl, ok := <-srv.tlist; ok {
		tl.Close()
		close(srv.tlist)
		srv.twlist <- tl
		return true
	}
	return false
}

// MaxConns sets new throttling limit (max number of simultaneous connections),
// returns number of free slots for incoming connections. For n < 0 doesn't change
// the limit. See limitnet.ThrottledListener for more detailed description.
//
// Won't return until srv.Serve is called.
func (srv *Server) MaxConns(n int) (free int) {
	srv.initialize()
	if tl, ok := <-srv.tlist; ok {
		free = tl.MaxConns(n)
		srv.tlist <- tl
		return
	}
	return 0
}
