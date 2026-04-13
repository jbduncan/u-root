// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"iter"
	"maps"
)

type set[V comparable] map[V]struct{}

func makeSet[V comparable]() set[V] {
	return make(set[V])
}

func (s set[V]) add(value V) {
	s[value] = struct{}{}
}

func (s set[V]) has(value V) bool {
	_, ok := s[value]
	return ok
}

func (s set[V]) remove(value V) {
	delete(s, value)
}

func (s set[V]) all() iter.Seq[V] {
	return maps.Keys(s)
}
