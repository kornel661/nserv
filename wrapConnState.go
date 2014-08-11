package nserv

import (
	"gopkg.in/kornel661/http.v0"
	"net"
	stdtp "net/http"
)

func (srv *Server) wrapConnState(fun func(net.Conn, stdtp.ConnState)) func(net.Conn, http.ConnState) {
	return func(conn net.Conn, newState http.ConnState) {
		fun(conn, stdtp.ConnState(newState))
		// put the token back to the jar when finished
		switch newState {
		case http.StateClosed, http.StateHijacked:
			// return token to the jar when connection closes
			srv.throttle <- token{}
		}
	}
}
