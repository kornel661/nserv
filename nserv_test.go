package nserv_test

/*
Tests internals of nserv.
*/

import (
	"runtime"
	"time"
)

const (
	// time to wait (let other goroutines do their job first)
	delay = 40 * time.Millisecond
	// address to listen to
	addr = "localhost:12345"
)

func init() {
	// enable parallelism (if your hardware supports it)
	runtime.GOMAXPROCS(runtime.NumCPU())
}
