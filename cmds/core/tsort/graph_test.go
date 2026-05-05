// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"slices"
	"testing"
)

func TestGraph(t *testing.T) {
	testAllNodes(t, graphFixture())
	testNodeCount(t, graphFixture())
	testSuccessors(t, graphFixture())
	testRemoveEdge(t, graphFixture())
}

func graphFixture() *graph {
	//    1     2      3   4
	//   / \   /|\     |
	//  /   \ / | \    |
	// 5     6  |  7   8
	//       |\ | /
	//       | \|/
	//       9  10
	// ...where edges are pointing downwards
	g := newGraph()
	g.putEdge(1, 5)
	g.putEdge(1, 6)
	g.putEdge(2, 6)
	g.putEdge(2, 7)
	g.putEdge(2, 10)
	g.putEdge(3, 8)
	g.putEdge(6, 9)
	g.putEdge(6, 10)
	g.putEdge(7, 10)
	g.addNode(4)
	return g
}

func testAllNodes(t *testing.T, g *graph) {
	got := slices.Collect(g.nodes())
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if diff := orderInsensitiveDiff(got, want); diff != "" {
		t.Fatalf(
			"allNodes mismatch (-actual +expected):\n%s",
			diff)
	}
}

func testNodeCount(t *testing.T, g *graph) {
	if got, want := g.nodeCount(), 10; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}

	g.addNode(11)
	if got, want := g.nodeCount(), 11; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}

	g.addNode(11)
	if got, want := g.nodeCount(), 11; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}
}

func testSuccessors(t *testing.T, g *graph) {
	if diff := orderInsensitiveDiff(slices.Collect(g.successors(1)), []int{5, 6}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"a\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors(2)), []int{6, 7, 10}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors(3)), []int{8}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors(6)), []int{9, 10}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"f\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors(7)), []int{10}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"c\") +expected):\n%s",
			diff)
	}
	if got := slices.Collect(g.successors(5)); len(got) > 0 {
		t.Errorf(`g.successors(8): want empty, got %v`, got)
	}
	if got := slices.Collect(g.successors(9)); len(got) > 0 {
		t.Errorf(`g.successors(9): want empty, got %v`, got)
	}
	if got := slices.Collect(g.successors(10)); len(got) > 0 {
		t.Errorf(`g.successors(10): want empty, got %v`, got)
	}
}

func testRemoveEdge(t *testing.T, g *graph) {
	g.removeEdge(42, 1)
	testSuccessors(t, g) // test that there were no changes

	g.removeEdge(1, 42)
	testSuccessors(t, g) // test that there were no changes

	g.removeEdge(2, 6)
	if diff := orderInsensitiveDiff(slices.Collect(g.successors(2)), []int{7, 10}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
}
