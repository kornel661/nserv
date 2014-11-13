package nserv

import (
	"errors"
	"fmt"
	"gopkg.in/kornel661/limitnet.v0"
)

// InitializeZeroDowntime sets up the commandline flags used by this package for
// supporting zero-downtime restarts. See also: limitnet.InitializeZeroDowntime().
// You need to execute flag.Parse() after InitializeZeroDowntime() for it to work.
// See "gopkg.in/kornel661/limitnet.v0/ZeroDowntime-example" how to use this
// feature.
func InitializeZeroDowntime() {
	limitnet.InitializeZeroDowntime()
}

// CanResumeServe tells if it seems possible to resume serving.
func CanResumeServe() bool {
	return limitnet.CanRetrieveListeners()
}

// ResumeServe tries to resume serving. Returns true if resuption was successful.
// Typically you execute InitializeZeroDowntime() and flag.Parse() first.
func (srv *Server) ResumeServe() (ok bool, err error) {
	listeners, err := limitnet.RetrieveListeners()
	if err != nil {
		return false, err
	}
	if len(listeners) != 1 {
		return false, fmt.Errorf("Inherited %d listeners instead of 1.", len(listeners))
	}
	return true, srv.Serve(listeners[0])
}

// ZeroDowntimeRestart shuts down the server and launches binary named the same
// as currently executing program with command line arguments args. The newly
// executed program inherits the file descriptor the srv server used.
func (srv *Server) ZeroDowntimeRestart(args ...string) error {
	err := srv.OperateOnListener(func(l limitnet.ThrottledListener) error {
		// prepare the command to be executed
		cmd, err := limitnet.PrepareCmd("", args, nil, l)
		if err != nil {
			return err
		}
		// start the command, return error
		return cmd.Start()
	})
	if err == nil {
		srv.Stop()
	}
	return err
}

// OperateOnListener applies function fun to the server's listener. It ensures
// the server is running during execution of fun.
func (srv *Server) OperateOnListener(fun func(limitnet.ThrottledListener) error) error {
	srv.initialize()
	// take the listener
	l, ok := <-srv.tlist
	if !ok {
		return errors.New("Server not running.")
	}
	defer func() {
		//replace the listener
		srv.tlist <- l
	}()
	return fun(l)
}
