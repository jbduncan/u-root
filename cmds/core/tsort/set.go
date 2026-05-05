// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"iter"
	"maps"
)

type set map[int]struct{}

func makeSet() set {
	return make(set)
}

func (s set) add(value int) {
	s[value] = struct{}{}
}

func (s set) remove(value int) {
	delete(s, value)
}

func (s set) all() iter.Seq[int] {
	return maps.Keys(s)
}
