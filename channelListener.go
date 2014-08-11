package nserv

import (
	"net"
	"time"
)

// chanListener returns a channel of accepted connections from the listener l.
// Close the listener to stop. Returned channel is closed when started goroutine
// exits.
func (srv *Server) chanListener(l net.Listener) <-chan net.Conn {
	conns := make(chan net.Conn, 0)
	go func() {
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
					srv.serv.Logf("http: Accept error: %v; retrying in %v", e, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				// filter out accept of closed listener errors
				if opErr, ok := e.(*net.OpError); ok && opErr.Op == "accept" && opErr.Err.Error() == "use of closed network connection" {
					e = nil
				} else if e.Error() == "use of closed network connection" {
					e = nil
				}
				srv.serverError <- e
				// signal server to stop
				srv.Stop()
				break
			}
			conns <- rw
			tempDelay = 0
		}
		close(conns)
	}()
	return conns
}
