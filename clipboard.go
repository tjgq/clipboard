// Package clipboard gives access to the system's text clipboard.
package clipboard

import (
	"sync"
)

var listeners struct {
	sync.Mutex
	m map[chan<- string]struct{}
}

type requestType int

const (
	watchRequest requestType = iota
	unwatchRequest
	getRequest
	setRequest
)

type reply struct {
	s   string
	err error
}

type request struct {
	t requestType
	s string
	c chan reply
}

// Use a control channel to relay clipboard requests into a handler goroutine.
// This is required because on some systems all clipboard calls must originate
// from the same thread.
var ctl = make(chan request, 1)

// Get returns the current clipboard value.
func Get() (string, error) {
	c := make(chan reply, 1)
	ctl <- request{getRequest, "", c}
	select {
	case r := <-c:
		return r.s, r.err
	}
}

// Set sets the clipboard value.
func Set(s string) error {
	c := make(chan reply, 1)
	ctl <- request{setRequest, s, c}
	select {
	case r := <-c:
		return r.err
	}
}

// Notify causes c to receive future clipboard values.
// If c is already receiving clipboard values, Notify is a no-op.
//
// This package will not block sending to c. It is the caller's responsibility
// to either ensure enough buffer is available.
func Notify(c chan<- string) {
	listeners.Lock()
	defer listeners.Unlock()

	if listeners.m == nil {
		listeners.m = make(map[chan<- string]struct{})
	}

	if len(listeners.m) == 0 {
		defer func() {
			r := make(chan reply, 1)
			ctl <- request{watchRequest, "", r}
			<-r
		}()
	}

	listeners.m[c] = struct{}{}
}

// Unnotify causes c to no longer receive clipboard values.
// If c is not receiving clipboard values, Unnotify is a no-op.
//
// It is guaranteed that c receives no more values after Unnotify returns.
func Unnotify(c chan<- string) {
	listeners.Lock()
	defer listeners.Unlock()

	delete(listeners.m, c)

	if len(listeners.m) == 0 {
		r := make(chan reply, 1)
		ctl <- request{unwatchRequest, "", r}
		<-r
	}
}

func broadcast(s string) {
	listeners.Lock()
	defer listeners.Unlock()
	for c := range listeners.m {
		select {
		case c <- s:
		default:
		}
	}
}
