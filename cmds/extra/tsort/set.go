// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"iter"
	"maps"
)

type set struct {
	m map[string]struct{}
}

func makeSet() set {
	return set{
		m: make(map[string]struct{}),
	}
}

func (s set) add(value string) {
	s.m[value] = struct{}{}
}

func (s set) has(value string) bool {
	_, ok := s.m[value]
	return ok
}

func (s set) remove(value string) {
	if !s.has(value) {
		panic("set is empty")
	}

	delete(s.m, value)
}

func (s set) all() iter.Seq[string] {
	return maps.Keys(s.m)
}
