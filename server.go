package nserv

import (
	"errors"
	"gopkg.in/kornel661/http.v0"
	"net"
	"time"
)

// New returns a new nserv.Server initialized with srv and throttleMax.
// Argument throttleMax has the same meaning as for the server initialization,
// serv.Initialize(...)
func New(s *http.Server, throttleMax int) *Server {
	var srv Server
	srv.Server = s
	srv.Initialize(throttleMax, throttleMax)
	return &srv
}

// Server with graceful exit and throttling.
// Server is an extension of http.Server from the standard library (its API is
// a superset of that of http.Server). It should only be created by nserv.New
type Server struct {
	*http.Server            // standard http.Server behaviour (but from slightly modified library)
	throttleMax  int        // max number of simultaneous connections
	throttle     chan token // channel used to throttle the requests (jar of tokens)
	finish       chan token // channel used to signal the server to quit (by srv.Stop())
	finished     chan token // channel used to signal that the server finished processing all requests
	// send a new instantenous limit for number of requests to this channel (<0 to exit)
	setMaxThrottle chan int
}

// Initialize given server. It is an error to initialize a server multiple times.
// It is also an error to use an uninitilized server. In this context error means
// panic (or undefined behaviour if executed in parallel).
//
// Argument throttleMax specifies maximum number of simultaneous requests the
// server will be able to serve. Instantenous maximum can be adjusted during
// runtime between 0 and throttleMax.
// Set throttleMax to 0 to disable throttling.
//
// Initial maximum for throttling is equal to initialMax.
//
// You can tinker with the server options before it is started.
func (srv *Server) Initialize(throttleMax, initialMax int) {
	if srv.finish != nil {
		panic("nserv.Server can be initialized only once.")
	}
	srv.finish = make(chan token, 1)
	srv.finished = make(chan token, 1)
	if srv.Server == nil {
		srv.Server = new(http.Server)
	}
	if throttleMax > 0 {
		srv.setMaxThrottle = make(chan int, 1)
		srv.throttleMax = throttleMax
		srv.throttle = make(chan token, throttleMax)
		// start throttling
		go srv.throttler()
		// set initial max
		srv.SetThrottle(initialMax)
	} else {
		srv.throttleMax = 0
		srv.throttle = nil
	}
}

// SetThrottle sets a new (instantenous) throttling limit.
// Returned error is not nil if n is out of bounds (0 <= n <= srv.throttleMax).
// Can be executed in parallel (from many goroutines) on initialized server that
// hasn't been stopped yet. (It's OK to run it on server that hasn't been
// started yet.)
//
// It may take some time to reach the new maximum (e.g., decreasing the maximum
// won't interrupt any active connections, it will rather wait for the connections
// to end). In other words the throttling limit will be eventually equal to the
// n provided as the argument (unless it's changed in the meantime).
func (srv *Server) SetThrottle(n int) error {
	if srv.throttleMax == 0 {
		return errors.New("nserv.SetThrottle: throttling is disabled")
	}
	if n < 0 {
		return errors.New("nserv.SetThrottle: cannot set negative limit")
	}
	if n > srv.throttleMax {
		return errors.New("nserv.SetThrottle: cannot set too high limit")
	}
	srv.setMaxThrottle <- n
	return nil
}

// Stop stops running server and returns a receiving channel which signals when
// the server stops. In other words:
// srv.Stop() // merely signals the server to finish
// <-srv.Stop() // signals the server to stop and "returns" only when the server stopped
//
// It is "thread-safe" and can be invoked multiple times.
func (srv *Server) Stop() <-chan token {
	ch := make(chan token, 1)
	select {
	case srv.finish <- token{}: // signal to stop
	default: // srv.finish is full, sb has already signalled
	}
	go func() {
		srv.waitShutdown()
		ch <- token{}
	}()
	return ch
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
//
// Based on the standard library, see:
// http://golang.org/src/pkg/net/http/server.go?s=50405:50451#L1684
func (srv *Server) Serve(l net.Listener) error {
	// wait for all connections when shutting down
	defer srv.waitShutdown()
	// the 'actual' server
	serv := srv.Server
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		if !srv.takeToken() {
			// we've been signalled to finish, pass the message to throttler goroutine
			srv.setMaxThrottle <- -1
			return nil
		}
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
				serv.Logf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				// give the token back:
				srv.replaceToken()
				continue
			}
			// give the token back:
			srv.replaceToken()
			// signal finish to throttler goroutine
			srv.setMaxThrottle <- 1
			return e
		}
		tempDelay = 0
		c, err := serv.NewConn(rw)
		if err != nil {
			// give the token back:
			srv.replaceToken()
			continue
		}
		c.SetState(c.Getrwc(), http.StateNew) // before Serve can return
		go c.Serve()
	}
}
