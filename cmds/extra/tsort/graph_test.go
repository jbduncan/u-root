// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"
)

func TestGraph(t *testing.T) {
	testNodes(t, fixtureGraph())
	testSuccessors(t, fixtureGraph())
	testInDegree(t, fixtureGraph())
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

func testNodes(t *testing.T, g *graph) {
	if diff := orderInsensitiveIterDiff(
		g.nodes(),
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	); diff != "" {
		t.Errorf(
			"values mismatch (-g.nodes() +expected):\n%s",
			diff)
	}
}

func testSuccessors(t *testing.T, g *graph) {
	if diff := orderInsensitiveIterDiff(g.successors("a"), "d", "e"); diff != "" {
		t.Errorf(
			"values mismatch (-g.successors(\"a\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveIterDiff(g.successors("b"), "e", "f", "i"); diff != "" {
		t.Errorf(
			"values mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveIterDiff(g.successors("e"), "h", "i"); diff != "" {
		t.Errorf(
			"values mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveIterDiff(g.successors("f"), "i"); diff != "" {
		t.Errorf(
			"values mismatch (-g.successors(\"f\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveIterDiff(g.successors("c"), "g"); diff != "" {
		t.Errorf(
			"values mismatch (-g.successors(\"c\") +expected):\n%s",
			diff)
	}
	caughtPanic := catchPanic(func() { g.successors("absent") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "node is not in graph") {
		t.Errorf(
			`g.successors("absent"): want panic with message "node is not in graph", got %#v`,
			caughtPanic)
	}
}

func testInDegree(t *testing.T, g *graph) {
	if g.inDegree("a") != 0 {
		t.Errorf(`g.inDegree("a"): want 0, got %d`, g.inDegree("a"))
	}
	if g.inDegree("d") != 1 {
		t.Errorf(`g.inDegree("d"): want 1, got %d`, g.inDegree("d"))
	}
	if g.inDegree("e") != 2 {
		t.Errorf(`g.inDegree("e"): want 2, got %d`, g.inDegree("e"))
	}
	if g.inDegree("i") != 3 {
		t.Errorf(`g.inDegree("i"): want 3, got %d`, g.inDegree("e"))
	}
	if g.inDegree("absent-node") != 0 {
		t.Errorf(
			`g.inDegree("absent-node"): want 0, got %d`,
			g.inDegree("absent-node"))
	}
}

func testRemoveEdge(t *testing.T, g *graph) {
	caughtPanic := catchPanic(func() { g.removeEdge("absent-source-node", "a") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "source node is not in graph") {
		t.Errorf(
			`g.removeEdge("absent-source-node", "a"): want panic with message "source node is not in graph", got %#v`,
			caughtPanic)
	}
	testSuccessors(t, g) // test that there were no changes

	caughtPanic = catchPanic(func() { g.removeEdge("a", "absent-target-node") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "target node is not in graph") {
		t.Errorf(
			`g.removeEdge("absent-target-node", "a"): want panic with message "target node is not in graph", got %#v`,
			caughtPanic)
	}
	testSuccessors(t, g) // test that there were no changes

	g.removeEdge("b", "e")
	if diff := orderInsensitiveIterDiff(g.successors("b"), "f", "i"); diff != "" {
		t.Errorf(
			"values mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveIterDiff(g.successors("e"), "h", "i"); diff != "" {
		t.Errorf(
			"values mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if g.inDegree("e") != 1 {
		t.Errorf(
			`g.removeEdge("b", "e"): want g.inDegree("e") to be 1, got %d`,
			g.inDegree("e"))
	}
}
