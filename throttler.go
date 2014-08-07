package nserv

// throttler is run in a separate goroutine. It listens on srv.setMaxThrottle
// and adds or removes tokens from the srv.throttle channnel.
// Negative value on srv.setMaxThrottle channel signals exit.
func (srv *Server) throttler() {
	var (
		instMax   = 0 // instantenous max == number of throttling tokens at large
		targetMax = 0 // target instantenous max, we want to make instMax = targetMax
	)
	// removes a token
	decrease := func() {
		select {
		case <-srv.throttle:
			instMax--
		case newMax := <-srv.setMaxThrottle:
			targetMax = newMax
		}
	}
	// adds a token
	increase := func() {
		select {
		case srv.throttle <- token{}:
			instMax++
		case newMax := <-srv.setMaxThrottle:
			targetMax = newMax
		}
	}
	// listens for a new instMax
	idle := func() {
		select {
		case newMax := <-srv.setMaxThrottle:
			targetMax = newMax
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
	srv.SetKeepAlivesEnabled(false)
	// reclaim all tokens (i.e., wait for all connections to finish)
	for i := 0; i < instMax; i++ {
		<-srv.throttle
	}
	instMax = 0
	// signal we're finished
	srv.finished <- token{}
}
