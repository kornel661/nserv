package nserv

// This file is, for the most part, copied from the standard net.Server
// implementation.

import (
	"crypto/tls"
	"net"
	"time"
)

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type TCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts the next incoming call and returns the new
// connection. KeepAlivePeriod is set properly.
func (ln *TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// Some default values set by the ListenAndServe method family. Feel free to modify
// these variables, e.g.,
//     nserv.DefaultMaxConns = 100
//     srv.ListenAndServe()
// Or, alternatively (arguably a better solution if you're writing a 'serious'
// server), use srv.Listen directly or reimplement ListenAndServe.
var (
	DefaultReadTimeout  = 60 * time.Second // default ReadTimeout set by the ListenAndServe methods
	DefaultWriteTimeout = 60 * time.Second // default WriteTimeout set by the ListenAndServe methods
	DefaultMaxConns     = 1000             // default MaxConns set by the ListenAndServe methods
)

// saneDefaults sets some 'sane' default timeouts and limits.
func (srv *Server) saneDefaults() {
	if srv.ReadTimeout == 0 {
		srv.ReadTimeout = DefaultReadTimeout
	}
	if srv.WriteTimeout == 0 {
		srv.WriteTimeout = DefaultWriteTimeout
	}
	if srv.InitialMaxConns == 0 {
		srv.InitialMaxConns = DefaultMaxConns
	}
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
// srv.Addr is blank, ":http" is used.
//
// If either of ReadTimeout, WriteTimeout or MaxConns is 0, it's going to be set
// to a 'sane' default value, see the corresponding Default... variable.
func (srv *Server) ListenAndServe() error {
	srv.saneDefaults()
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(&TCPKeepAliveListener{ln.(*net.TCPListener)})
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
//
// Filenames containing a certificate and matching private key for
// the server must be provided. If the certificate is signed by a
// certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
//
// If either of ReadTimeout, WriteTimeout or MaxConns is 0, it's going to be set
// to a 'sane' default value, see the corresponding Default... variable.
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	srv.saneDefaults()
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(&TCPKeepAliveListener{ln.(*net.TCPListener)}, config)
	return srv.Serve(tlsListener)
}
