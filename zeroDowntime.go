package nserv

import (
	"errors"
	"fmt"
	"gopkg.in/kornel661/limitnet.v0"
	"os"
)

// InitializeZeroDowntime sets up the command-line flags used by this package for
// supporting zero-downtime restarts. See also: limitnet.InitializeZeroDowntime().
// You need to execute flag.Parse() after InitializeZeroDowntime() for it to work.
// See "gopkg.in/kornel661/limitnet.v0/ZeroDowntime-example" how to use this
// feature.
func InitializeZeroDowntime() {
	limitnet.InitializeZeroDowntime()
}

// CanResume tells if it seems possible to resume serving.
func CanResume() bool {
	return limitnet.CanRetrieveListeners()
}

// ResumeAndServe tries to resume serving.
// Typically you execute InitializeZeroDowntime() and flag.Parse() first.
//
// First, limitnet.RetrieveListeners() is called to retrieve a listener. Next,
// if either of srv.ReadTimeout, srv.WriteTimeout or srv.MaxConns is 0, it's
// going to be set to a 'sane' default value, see the corresponding Default...
// variables. Finally, srv.Serve method is invoked with the retrieved listener
// as its argument.
func (srv *Server) ResumeAndServe() error {
	listeners, err := limitnet.RetrieveListeners()
	if err != nil {
		return err
	}
	if len(listeners) != 1 {
		for _, l := range listeners {
			l.Close()
		}
		return fmt.Errorf("Inherited %d listeners instead of 1.", len(listeners))
	}
	srv.saneDefaults()
	return srv.Serve(listeners[0])
}

// ZeroDowntimeRestart shuts down the server and launches binary named the same
// as currently executing program with command line arguments args. The newly
// executed program inherits the file descriptor the srv server used.
//
// Error behavior similar to Server.OperateOnListener or due to command
// execution error.
func (srv *Server) ZeroDowntimeRestart(args ...string) error {
	err := srv.OperateOnListener(func(l limitnet.ThrottledListener) error {
		// prepare the command to be executed
		cmd, err := limitnet.PrepareCmd("", args, nil, l)
		if err != nil {
			return err
		}
		// start the command, return error
		err = cmd.Start()
		cmd.ExtraFiles[0].Close() // close unused file
		return err
	})
	if err == nil {
		srv.Stop()
	}
	return err
}

// CopyListenerFD returns DUP of the file descriptor associated with the listener.
// If the server isn't running the behavior is as in Server.OperateOnListener.
func (srv *Server) CopyListenerFD() (fd *os.File, err error) {
	srv.OperateOnListener(func(l limitnet.ThrottledListener) error {
		fd, err = limitnet.CopyFD(l)
		return err
	})
	return
}

// OperateOnListener applies function fun to the server's listener. It ensures
// the server is running during execution of fun (returns an error if stopped or
// hangs if it hasn't been started).
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
