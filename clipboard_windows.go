package clipboard

// #include <stdlib.h>
//
// char *get();
// int set(const char *);
import "C"

import (
	"errors"
	"github.com/tjgq/ticker"
	"time"
	"unsafe"
)

const (
	pollInterval = time.Second
)

var tick = ticker.New(pollInterval)

var last string

func init() {
	go loop()
}

func loop() {
	for {
		select {
		case <-tick.C:
			poll()
		case r := <-ctl:
			handle(r)
		}
	}
}

func handle(r request) {
	switch r.t {
	case watchRequest:
		tick.Start()
		r.c <- reply{}
	case unwatchRequest:
		tick.Stop()
		r.c <- reply{}
	case getRequest:
		s := C.get()
		if s == nil {
			r.c <- reply{"", errors.New("clipboard read error")}
		} else {
			r.c <- reply{C.GoString(s), nil}
			C.free(unsafe.Pointer(s))
		}
	case setRequest:
		s := C.CString(r.s)
		if C.set(s) == 0 {
			r.c <- reply{"", errors.New("clipboard write error")}
		} else {
			r.c <- reply{"", nil}
		}
		C.free(unsafe.Pointer(s))
	}
}

func poll() {
	cstr := C.get()
	if cstr != nil {
		s := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		if s != last {
			last = s
			broadcast(s)
		}
	}
}
