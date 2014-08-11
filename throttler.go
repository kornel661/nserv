package nserv

// throttler is run in a separate goroutine. It listens on srv.setMaxThrottle
// and adds or removes tokens from the srv.throttle channnel.
// Negative value on srv.setMaxThrottle channel signals exit.
func (srv *Server) throttler() {
	var (
		instMax   = 0 // instantenous max == number of throttling tokens at large
		targetMax = 0 // target instantenous max, we want to make instMax = targetMax
	)
	// removes a token from the jar
	decrease := func() {
		select {
		case <-srv.throttle:
			instMax--
		case targetMax = <-srv.setMaxThrottle:
		}
	}
	// adds a token to the jar
	increase := func() {
		select {
		case srv.throttle <- token{}:
			instMax++
		case targetMax = <-srv.setMaxThrottle:
		}
	}
	// listens for a new instMax
	idle := func() {
		select {
		case targetMax = <-srv.setMaxThrottle:
		}
	}
	// loop until signaled to exit (targetMax < 0)
	for targetMax >= 0 {
		switch {
		case instMax < targetMax:
			increase()
		case instMax == targetMax:
			idle()
		case instMax > targetMax:
			decrease()
		}
	}
	// server is shutting down, switch off keep-alive connections
	srv.serv.SetKeepAlivesEnabled(false)
	// reclaim all tokens (i.e., wait for all connections to finish)
	for i := 0; i < instMax; i++ {
		<-srv.throttle
	}
	instMax = 0
	// signal we're finished
	srv.finished <- token{}
}

// throttlerStop signals the throttler goroutine to start shutdown procedure:
// to wait for all requests to finish and signal on srv.finished at the end.
func (srv *Server) throttlerStop() {
	select {
	case srv.setMaxThrottle <- -1: // signal to stop
	default: // throttler's been already signalled to stop
	}
}
