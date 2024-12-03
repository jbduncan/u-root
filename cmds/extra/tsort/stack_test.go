package main

import (
	"testing"
)

func TestStack(t *testing.T) {
	s := makeStack()

	s.push("a")
	if s.peek() != "a" {
		t.Errorf(`stack %#v: want peeked element to be "a", but was not`, s)
	}
	s.push("b")
	if s.peek() != "b" {
		t.Errorf(`stack %#v: want peeked element to be "b", but was not`, s)
	}

	if s.pop() != "b" {
		t.Errorf(`stack %#v: want popped element to be "b", but was not`, s)
	}
	if s.pop() != "a" {
		t.Errorf(`stack %#v: want popped element to be "a", but was not`, s)
	}

	caughtPanic := catchPanic(func() { s.peek() })
	if caughtPanic == nil {
		t.Fatalf(`stack %#v: want peek to panic, got no panic`, s)
	}
	if caughtPanic.Error() != "stack is empty" {
		t.Fatalf(
			`stack %#v: want peek to panic with "stack is empty", got %q`,
			s, caughtPanic)
	}

	caughtPanic = catchPanic(func() { s.pop() })
	if caughtPanic == nil {
		t.Fatalf(`stack %#v: want pop to panic, got no panic`, s)
	}
	if caughtPanic.Error() != "stack is empty" {
		t.Fatalf(
			`stack %#v: want pop to panic with "stack is empty", got %q`,
			s, caughtPanic)
	}
}
