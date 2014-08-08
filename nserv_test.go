package nserv

/*
Tests internals of nserv.
*/

import (
	"runtime"
	"time"
)

const (
	sleepInt = 1 * time.Millisecond
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
