package main

import "unique"

type str unique.Handle[string]

func strOf(value string) str {
	return str(unique.Make(value))
}

func (s str) String() string {
	return unique.Handle[string](s).Value()
}

func (s str) Equal(u str) bool {
	return s == u
}
