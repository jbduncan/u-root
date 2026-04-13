// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestGraph(t *testing.T) {
	testAllNodes(t, fixtureGraph())
	testNodeCount(t, fixtureGraph())
	testSuccessors(t, fixtureGraph())
	testInDegree(t, fixtureGraph())
	testRemoveNode(t, fixtureGraph())
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
	// ...where edges are pointed downwards
	g := newGraph()
	putEdge(g, "a", "d")
	putEdge(g, "a", "e")
	putEdge(g, "b", "e")
	putEdge(g, "b", "f")
	putEdge(g, "b", "i")
	putEdge(g, "b", "i")
	putEdge(g, "e", "h")
	putEdge(g, "e", "i")
	putEdge(g, "f", "i")
	putEdge(g, "c", "g")
	addNode(g, "j")
	return g
}

func testAllNodes(t *testing.T, g *graph) {
	got := slices.Collect(g.nodes())
	want := slice("a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	if diff := orderInsensitiveDiffByValue(got, want); diff != "" {
		t.Fatalf(
			"allNodes mismatch (-actual +expected):\n%s",
			diff)
	}
}

func testNodeCount(t *testing.T, g *graph) {
	if got, want := g.nodeCount(), 10; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}

	addNode(g, "k")
	if got, want := g.nodeCount(), 11; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}
}

func testSuccessors(t *testing.T, g *graph) {
	if diff := orderInsensitiveDiffByValue(successors(g, "a"), slice("d", "e")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"a\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiffByValue(successors(g, "b"), slice("e", "f", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiffByValue(successors(g, "e"), slice("h", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiffByValue(successors(g, "f"), slice("i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"f\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiffByValue(successors(g, "c"), slice("g")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"c\") +expected):\n%s",
			diff)
	}
	if got := successors(g, "h"); len(got) > 0 {
		t.Errorf(`g.successors("h"): want empty, got %v`, got)
	}
	if got := successors(g, "i"); len(got) > 0 {
		t.Errorf(`g.successors("i"): want empty, got %v`, got)
	}
	if got := successors(g, "j"); len(got) > 0 {
		t.Errorf(`g.successors("j"): want empty, got %v`, got)
	}
	caughtPanic := catchPanic(func() { successors(g, "absent") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "node is not in graph") {
		t.Errorf(
			`g.successors("absent"): want panic with message "node is not in graph", got %#v`,
			caughtPanic)
	}
}

func testInDegree(t *testing.T, g *graph) {
	if inDegree(g, "a") != 0 {
		t.Errorf(`g.inDegree("a"): want 0, got %d`, inDegree(g, "a"))
	}
	if inDegree(g, "d") != 1 {
		t.Errorf(`g.inDegree("d"): want 1, got %d`, inDegree(g, "d"))
	}
	if inDegree(g, "e") != 2 {
		t.Errorf(`g.inDegree("e"): want 2, got %d`, inDegree(g, "e"))
	}
	if inDegree(g, "i") != 3 {
		t.Errorf(`g.inDegree("i"): want 3, got %d`, inDegree(g, "e"))
	}
	if inDegree(g, "absent-node") != 0 {
		t.Errorf(
			`g.inDegree("absent-node"): want 0, got %d`,
			inDegree(g, "absent-node"))
	}
}

func testRemoveNode(t *testing.T, g *graph) {
	caughtPanic := catchPanic(func() { removeNode(g, "absent-node") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "node is not in graph") {
		t.Errorf(
			`g.removeNode("absent-node"): want panic with message "node is not in graph", got %#v`,
			caughtPanic)
	}

	removeNode(g, "j")
	if diff := orderInsensitiveDiffByValue(
		slices.Collect(g.nodes()),
		slice("a", "b", "c", "d", "e", "f", "g", "h", "i"),
	); diff != "" {
		t.Fatalf("g.removeNode(\"j\"): nodes mismatch (-got +want):\n%s", diff)
	}

	removeNode(g, "c")
	if diff := orderInsensitiveDiffByValue(
		slices.Collect(g.nodes()),
		slice("a", "b", "d", "e", "f", "g", "h", "i"),
	); diff != "" {
		t.Errorf("g.removeNode(\"c\"): nodes mismatch (-got +want):\n%s", diff)
	}
}

func testRemoveEdge(t *testing.T, g *graph) {
	caughtPanic := catchPanic(func() { removeEdge(g, "absent-source-node", "a") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "source node is not in graph") {
		t.Errorf(
			`g.removeEdge("absent-source-node", "a"): want panic with message "source node is not in graph", got %#v`,
			caughtPanic)
	}
	testSuccessors(t, g) // test that there were no changes

	caughtPanic = catchPanic(func() { removeEdge(g, "a", "absent-target-node") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "target node is not in graph") {
		t.Errorf(
			`g.removeEdge("absent-target-node", "a"): want panic with message "target node is not in graph", got %#v`,
			caughtPanic)
	}
	testSuccessors(t, g) // test that there were no changes

	removeEdge(g, "b", "e")
	if diff := orderInsensitiveDiffByValue(successors(g, "b"), slice("f", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiffByValue(successors(g, "e"), slice("h", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if inDegree(g, "e") != 1 {
		t.Errorf(
			`g.removeEdge("b", "e"): want g.inDegree("e") to be 1, got %d`,
			inDegree(g, "e"))
	}
}

func slice(values ...string) []str {
	result := make([]str, 0, len(values))
	for _, value := range values {
		result = append(result, strOf(value))
	}
	return result
}

func putEdge(g *graph, source string, target string) {
	g.putEdge(strOf(source), strOf(target))
}

func addNode(g *graph, node string) {
	g.addNode(strOf(node))
}

func successors(g *graph, node string) []str {
	return slices.Collect(g.successors(strOf(node)))
}

func inDegree(g *graph, node string) int {
	return g.inDegree(strOf(node))
}

func removeNode(g *graph, node string) {
	g.removeNode(strOf(node))
}

func removeEdge(g *graph, source string, target string) {
	g.removeEdge(strOf(source), strOf(target))
}

func catchPanic(f func()) (caughtPanic error) {
	defer func() {
		if e := recover(); e != nil {
			caughtPanic = fmt.Errorf("%v", e)
		}
	}()

	f()
	return
}
