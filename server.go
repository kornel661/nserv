package nserv

import (
	"gopkg.in/kornel661/http.v0"
	"net"
	"time"
)

// tokens used for communication via channels
type token struct{}

// Server with graceful exit and throttling.
// It's an extension of http.Server from the standard library (its API is
// a superset of that of http.Server).
type Server struct {
	// standard http.Server
	*http.Server
	// max number of simultaneous
	throttleMax int
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
//
// Based on the standard library, see:
// http://golang.org/src/pkg/net/http/server.go?s=50405:50451#L1684
func (s *Server) Serve(l net.Listener) error {
	srv := s.Server
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.Logf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c, err := srv.NewConn(rw)
		if err != nil {
			continue
		}
		c.SetState(c.Getrwc(), http.StateNew) // before Serve can return
		go c.Serve()
	}
}
