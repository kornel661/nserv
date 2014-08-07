package nserv

/* This file is basically copy-paste from the standard library.
   FIXME: If you have a better idea how to imitate inheritance please let me
          know.
   If implementation in the standard library changes this file should be updated.
*/

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// A conn represents the server side of an HTTP connection.
// See:
// http://golang.org/src/pkg/net/http/server.go#L106
type conn struct {
	remoteAddr string               // network address of remote side
	server     *http.Server         // the Server on which the connection arrived
	rwc        net.Conn             // i/o connection
	sr         liveSwitchReader     // where the LimitReader reads from; usually the rwc
	lr         *io.LimitedReader    // io.LimitReader(sr)
	buf        *bufio.ReadWriter    // buffered(lr,rwc), reading from bufio->limitReader->sr->rwc
	tlsState   *tls.ConnectionState // or nil when not using TLS

	mu           sync.Mutex // guards the following
	clientGone   bool       // if client has disconnected mid-request
	closeNotifyc chan bool  // made lazily
	hijackedv    bool       // connection has been hijacked by handler
}

func (c *conn) hijacked() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hijackedv
}

func (c *conn) hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.hijackedv {
		return nil, nil, ErrHijacked
	}
	if c.closeNotifyc != nil {
		return nil, nil, errors.New("http: Hijack is incompatible with use of CloseNotifier")
	}
	c.hijackedv = true
	rwc = c.rwc
	buf = c.buf
	c.rwc = nil
	c.buf = nil
	c.setState(rwc, StateHijacked)
	return
}

func (c *conn) closeNotify() <-chan bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc == nil {
		c.closeNotifyc = make(chan bool, 1)
		if c.hijackedv {
			// to obey the function signature, even though
			// it'll never receive a value.
			return c.closeNotifyc
		}
		pr, pw := io.Pipe()

		readSource := c.sr.r
		c.sr.Lock()
		c.sr.r = pr
		c.sr.Unlock()
		go func() {
			_, err := io.Copy(pw, readSource)
			if err == nil {
				err = io.EOF
			}
			pw.CloseWithError(err)
			c.noteClientGone()
		}()
	}
	return c.closeNotifyc
}

func (c *conn) noteClientGone() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc != nil && !c.clientGone {
		c.closeNotifyc <- true
	}
	c.clientGone = true
}

// Create new connection from rwc.
// Replacement of
// func (srv *Server) newConn(rwc net.Conn) (c *conn, err error)
// See:
// http://golang.org/src/pkg/net/http/server.go#L423
func httpServer_newConn(srv *http.Server, rwc net.Conn) (c *conn, err error) {
	c = new(conn)
	c.remoteAddr = rwc.RemoteAddr().String()
	c.server = srv
	c.rwc = rwc
	c.sr = liveSwitchReader{r: c.rwc}
	c.lr = io.LimitReader(&c.sr, noLimit).(*io.LimitedReader)
	br := newBufioReader(c.lr)
	bw := newBufioWriterSize(c.rwc, 4<<10)
	c.buf = bufio.NewReadWriter(br, bw)
	return c, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// Basically copied from standard library, see:
// http://golang.org/src/pkg/net/http/server.go#L1938
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// Replacement for
// func (s *Server) logf(format string, args ...interface{})
// see:
// http://golang.org/src/pkg/net/http/server.go#L1741
func httpServer_logf(s *http.Server, format string, args ...interface{}) {
	if s.ErrorLog != nil {
		s.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
// srv.Addr is blank, ":http" is used.
//
// Basically copied from standard library, see:
// http://golang.org/src/pkg/net/http/server.go?s=49983:50024#L1669
func (s *Server) ListenAndServe() error {
	srv := s.Server
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
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
// Basically copied from standard library, see:
// http://golang.org/src/pkg/net/http/server.go?s=54154:54222#L1813
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	srv := s.Server
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

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	return s.Serve(tlsListener)
}
