package nserv

// tokens used for communication via channels
type token struct{}

// takeToken takes one token from the srv.throttle jar and returns true
// or accepts srv.finish token and returns false
func (srv *Server) takeToken() bool {
	select {
	case <-srv.throttle:
		return true
	case <-srv.finish:
		return false
	}
}

// replaceToken puts a token back in srv.throttle
func (srv *Server) replaceToken() {
	srv.throttle <- token{}
}

// waits for all requests to finish processing
func (srv *Server) waitShutdown() {
	// wait for all requests to finish
	<-srv.finished
	// replace token
	srv.finished <- token{}
}
