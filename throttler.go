package nserv

import "log"

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
	// listens for a new instMax (when instMax == targetMax)
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

	log.Println("throttler: exitting...")
	// reclaim all tokens (i.e., wait for all connections to finish)
	for i := 0; i < instMax; i++ {
		<-srv.throttle
	}
	instMax = 0
	log.Println("throttler: signal finish...")
	// signal we're finished
	srv.finished <- token{}
}
