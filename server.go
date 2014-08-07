package nserv

// tokens used for communication via channels
type token struct{}

type Server struct {
	// max number of simultaneous
	throttleMax int
}
