// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type queue struct {
	q []str
}

func (q *queue) enqueue(value str) {
	q.q = append(q.q, value)
}

func (q *queue) dequeue() (str, bool) {
	if len(q.q) == 0 {
		return str{}, false
	}

	result := q.q[0]

	q.q = q.q[1:]

	return result, true
}
