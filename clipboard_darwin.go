package clipboard

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
//
// #include <stdlib.h>
//
// long count();
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

var seq C.long

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
		seq = C.count()
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
	newseq := C.count()
	if newseq != seq {
		seq = newseq
		s := C.get()
		if s != nil {
			broadcast(C.GoString(s))
			C.free(unsafe.Pointer(s))
		}
	}
}
