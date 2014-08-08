package nserv

import (
	"runtime"
	"testing"
	"time"
)

func TestThrottler(t *testing.T) {
	max := 100
	srv := New(nil, max)
	time.Sleep(sleepInt) // wait for tokens
	if l := len(srv.throttle); l != max {
		t.Errorf("Number of throtte tokens is %d instead of %d.", l, max)
	}

	// test bouds checks
	if srv.SetThrottle(max+1) == nil {
		t.Error("Set too high throttle limit.")
	}
	if srv.SetThrottle(-1) == nil {
		t.Error("Set too low throttle limit.")
	}
	if l := len(srv.throttle); l != max {
		t.Errorf("Number of throtte tokens is %d instead of %d.", l, max)
	}

	// test setting
	if srv.SetThrottle(max/2) != nil {
		t.Error("Error setting throttle max/2.")
	}
	time.Sleep(sleepInt) // wait for tokens
	if l := len(srv.throttle); l != max/2 {
		t.Errorf("Number of throtte tokens is %d instead of %d.", l, max/2)
	}

	for i := 1; i <= 10; i++ {
		if srv.SetThrottle(max/i) != nil {
			t.Errorf("loop1: Error setting throttle max/%d.", i)
		}
	}
	for i := 10; i >= 1; i-- {
		if srv.SetThrottle(max/i) != nil {
			t.Errorf("loop2: Error setting throttle max/%d.", i)
		}
	}
	for i := 1; i <= 10; i++ {
		if srv.SetThrottle(max/i) != nil {
			t.Errorf("loop3: Error setting throttle max/%d.", i)
		}
		runtime.Gosched()
	}
	runtime.Gosched()
	time.Sleep(sleepInt) // wait for tokens
	if l := len(srv.throttle); l != max/10 {
		t.Errorf("Number of throtte tokens is %d instead of %d.", l, max/10)
	}

	srv.setMaxThrottle <- -1 // signal the end
	<-srv.finished           // wait for throttler to finish
	if len(srv.throttle) > 0 {
		t.Error("tokens left in srv.throttle after finish.")
	}
}
