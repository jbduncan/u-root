package main

import "testing"

func TestIDAssigner(t *testing.T) {
	i := newIDAssigner()

	if got, want := i.assignID("foo"), 0; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := i.valueFor(0), "foo"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := i.assignID("bar"), 1; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := i.valueFor(1), "bar"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := i.assignID("foo"), 0; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
