package main

import (
	"fmt"
)

func catchPanic(f func()) (caughtPanic error) {
	defer func() {
		if e := recover(); e != nil {
			caughtPanic = fmt.Errorf("%v", e)
		}
	}()

	f()
	return
}
