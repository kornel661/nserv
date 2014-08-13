package nserv

import (
	"gopkg.in/kornel661/http.v0"
	"log"
	"net"
	stdtp "net/http"
)

func (srv *Server) wrapConnState(fun func(net.Conn, stdtp.ConnState)) func(net.Conn, http.ConnState) {
	log.Println("returning function")
	return func(conn net.Conn, newState http.ConnState) {
		log.Printf("connection: ConnState, state: %v\n", newState)
		if fun != nil {
			fun(conn, stdtp.ConnState(newState))
		}
		// put the token back to the jar when finished
		switch newState {
		case http.StateClosed, http.StateHijacked:
			log.Println("connection: returning token")
			// return token to the jar when connection closes
			srv.throttle <- token{}
		}
	}
}
