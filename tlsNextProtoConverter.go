package nserv

import (
	"crypto/tls"
	"gopkg.in/kornel661/http.v0"
	stdtp "net/http"
	"unsafe"
)

func init() {
	// mild sanity check
	if (unsafe.Sizeof(http.Server{}) != unsafe.Sizeof(stdtp.Server{})) {
		panic("You're using unsupported version of Go.")
	}
}

// tlsNPC conterts TLSNextProto values between packages
func tlsNPC(std map[string]func(*stdtp.Server, *tls.Conn, stdtp.Handler)) map[string]func(*http.Server, *tls.Conn, http.Handler) {
	m := make(map[string]func(*http.Server, *tls.Conn, http.Handler), len(std))
	for str, fun := range std {
		newFun := func(s *http.Server, c *tls.Conn, h http.Handler) {
			fun((*stdtp.Server)(unsafe.Pointer(s)), c, handlerRConv{h})
		}
		m[str] = newFun
	}
	return m
}
