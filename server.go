package nserv

import (
	"net/http"
)

// tokens used for communication via channels
type token struct{}

// Server with graceful exit and throttling
type Server struct {
	*http.Server
	// max number of simultaneous
	throttleMax int
}
