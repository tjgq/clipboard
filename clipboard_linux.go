package clipboard

import (
	"github.com/conformal/gotk3/gdk"
	"github.com/conformal/gotk3/glib"
	"github.com/conformal/gotk3/gtk"
	"runtime"
)

func init() {
	go loop()
}

func loop() {

	runtime.LockOSThread()

	gtk.Init(nil)

	cb, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	if err != nil {
		panic("clipboard: unable to get clipboard")
	}

	h, err := cb.Connect("owner-change", poll)
	if err != nil {
		panic("clipboard: unable to connect signal handler")
	}

	cb.HandlerBlock(h)

	go func() {
		for {
			select {
			case r := <-ctl:
				glib.IdleAdd(handle, r, cb, h)
			}
		}
	}()

	gtk.Main()
}

func handle(r request, cb *gtk.Clipboard, h glib.SignalHandle) bool {
	switch r.t {
	case watchRequest:
		cb.HandlerUnblock(h)
		r.c <- reply{}
	case unwatchRequest:
		cb.HandlerBlock(h)
		r.c <- reply{}
	case getRequest:
		s, err := cb.WaitForText()
		r.c <- reply{s, err}
	case setRequest:
		cb.SetText(r.s)
		r.c <- reply{"", nil}
	}
	return false
}

func poll(cb *gtk.Clipboard, _ *gdk.Event) {
	s, err := cb.WaitForText()
	if err == nil {
		broadcast(s)
	}
}
