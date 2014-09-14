package clipboard

import (
	"testing"
	"time"
)

const wait = time.Second

func TestGetSet(t *testing.T) {
	want := "test"

	if err := Set(want); err != nil {
		t.Fatalf("can't set clipboard: %v", err)
	}

	got, err := Get()
	if err != nil {
		t.Fatalf("can't get clipboard: %v", err)
	}

	if want != got {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestNotifyUnnotify(t *testing.T) {

	var testStrings = []string{
		"foo",
		"bar",
		"baz",
	}

	c := make(chan string, len(testStrings))

	Notify(c)

	for i := 0; i < len(testStrings); i++ {
		Set(testStrings[i])
		time.Sleep(2 * wait)
	}

	if len(c) != 3 {
		panic("notification lost")
	}

	for j := 0; j < len(testStrings); j++ {
		want := testStrings[j]
		got := <-c
		if want != got {
			t.Fatalf("want %s got %s", want, got)
		}
	}

	Unnotify(c)
	Set("another")
	time.Sleep(2 * wait)
	if len(c) > 0 {
		t.Fatal("notification after Unnotify")
	}
}
