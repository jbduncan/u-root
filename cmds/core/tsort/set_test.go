// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"slices"
	"testing"
)

func TestIntSet(t *testing.T) {
	s := makeSet()

	if len(s) != 0 {
		t.Errorf(`set %#v: want len of 0, got %d`, s, len(s))
	}

	s.add(2)
	s.add(4)
	s.add(3)
	s.add(3)
	s.add(5)
	s.add(1)

	if diff := orderInsensitiveDiff(
		slices.Collect(s.all()),
		[]int{1, 2, 3, 4, 5},
	); diff != "" {
		t.Errorf("set iterator mismatch (-got +want):\n%s", diff)
	}

	s.remove(1)
	s.remove(5)
	s.remove(3)

	if diff := orderInsensitiveDiff(
		slices.Collect(s.all()),
		[]int{2, 4},
	); diff != "" {
		t.Errorf("set iterator mismatch (-got +want):\n%s", diff)
	}
}
