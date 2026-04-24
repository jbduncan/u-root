// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"slices"
	"testing"
)

func TestGraph(t *testing.T) {
	testAllNodes(t, fixtureGraph())
	testNodeCount(t, fixtureGraph())
	testSuccessors(t, fixtureGraph())
	testRemoveEdge(t, fixtureGraph())
}

func fixtureGraph() *graph {
	//    a     b      c   j
	//   / \   /|\     |
	//  /   \ / | \    |
	// d     e  |  f   g
	//       |\ | /
	//       | \|/
	//       h  i
	// ...where edges are pointing downwards
	g := newGraph()
	g.putEdge("a", "d")
	g.putEdge("a", "e")
	g.putEdge("b", "e")
	g.putEdge("b", "f")
	g.putEdge("b", "i")
	g.putEdge("b", "i")
	g.putEdge("e", "h")
	g.putEdge("e", "i")
	g.putEdge("f", "i")
	g.putEdge("c", "g")
	g.addNode("j")
	return g
}

func testAllNodes(t *testing.T, g *graph) {
	got := slices.Collect(g.nodes())
	want := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
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

	g.addNode("k")
	if got, want := g.nodeCount(), 11; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}

	g.addNode("k")
	if got, want := g.nodeCount(), 11; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}
}

func testSuccessors(t *testing.T, g *graph) {
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("a")), []string{"d", "e"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"a\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("b")), []string{"e", "f", "i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("e")), []string{"h", "i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("f")), []string{"i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"f\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("c")), []string{"g"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"c\") +expected):\n%s",
			diff)
	}
	if got := slices.Collect(g.successors("h")); len(got) > 0 {
		t.Errorf(`g.successors("h"): want empty, got %v`, got)
	}
	if got := slices.Collect(g.successors("i")); len(got) > 0 {
		t.Errorf(`g.successors("i"): want empty, got %v`, got)
	}
	if got := slices.Collect(g.successors("j")); len(got) > 0 {
		t.Errorf(`g.successors("j"): want empty, got %v`, got)
	}
}

func testRemoveEdge(t *testing.T, g *graph) {
	g.removeEdge("absent-source-node", "a")
	testSuccessors(t, g) // test that there were no changes

	g.removeEdge("a", "absent-target-node")
	testSuccessors(t, g) // test that there were no changes

	g.removeEdge("b", "e")
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("b")), []string{"f", "i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("e")), []string{"h", "i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
}
