package nserv

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
)

func init() {
	// enable parallelism (if your hardware supports it)
	runtime.GOMAXPROCS(runtime.NumCPU())
}
