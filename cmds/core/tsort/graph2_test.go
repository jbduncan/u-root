// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"slices"
	"testing"
)

// TODO: rename
func TestGraph2(t *testing.T) {
	// TODO: rename
	testValueFor2(t, graphFixture2())
	testNodeIDs2(t, graphFixture2())
	testNodeCount2(t, graphFixture2())
	testSuccessorIDs2(t, graphFixture2())
	testRemoveEdgeBetweenNodeIDs2(t, graphFixture2())
}

func graphFixture2() *graph2 {
	//    a     b      c   j
	//   / \   /|\     |
	//  /   \ / | \    |
	// d     e  |  f   g
	//       |\ | /
	//       | \|/
	//       h  i
	// ...where edges are pointed downwards
	g := newGraph2()
	g.putEdge("a", "d") // node IDs 0 and 1
	g.putEdge("a", "e") // node IDs 0 and 2
	g.putEdge("b", "e") // node IDs 3 and 2
	g.putEdge("b", "f") // node IDs 3 and 4
	g.putEdge("b", "i") // node IDs 3 and 5
	g.putEdge("b", "i") // node IDs 3 and 5
	g.putEdge("e", "h") // node IDs 2 and 6
	g.putEdge("e", "i") // node IDs 2 and 5
	g.putEdge("f", "i") // node IDs 4 and 5
	g.putEdge("c", "g") // node IDs 7 and 8
	g.addNode("j")      // node ID 9
	return g
}

func testValueFor2(t *testing.T, g *graph2) {
	for _, tt := range []struct {
		id   int32
		node string
	}{
		{id: 0, node: "a"},
		{id: 1, node: "d"},
		{id: 2, node: "e"},
		{id: 3, node: "b"},
		{id: 4, node: "f"},
		{id: 5, node: "i"},
		{id: 6, node: "h"},
		{id: 7, node: "c"},
		{id: 8, node: "g"},
		{id: 9, node: "j"},
	} {
		t.Run(
			fmt.Sprintf("g.valueFor(%d) == %s", tt.id, tt.node),
			func(t *testing.T) {
				if got, want := g.valueFor(tt.id), tt.node; got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			},
		)
	}
}

func testNodeIDs2(t *testing.T, g *graph2) {
	t.Run("g.nodeIDs()", func(t *testing.T) {
		got := slices.Collect(g.nodeIDs())
		want := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		if diff := orderInsensitiveDiff(got, want); diff != "" {
			t.Fatalf(
				"mismatch (-actual +expected):\n%s",
				diff)
		}
	})
}

func testNodeCount2(t *testing.T, g *graph2) {
	t.Run("g.nodeCount()", func(t *testing.T) {
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
	})
}

func testSuccessorIDs2(t *testing.T, g *graph2) {
	for _, tt := range []struct {
		id         int32
		successors []int32
	}{
		{id: 0, successors: []int32{1, 2}},
		{id: 1, successors: []int32{}},
		{id: 2, successors: []int32{5, 6}},
		{id: 3, successors: []int32{2, 4, 5}},
		{id: 4, successors: []int32{5}},
		{id: 5, successors: []int32{}},
		{id: 6, successors: []int32{}},
		{id: 7, successors: []int32{8}},
		{id: 8, successors: []int32{}},
		{id: 9, successors: []int32{}},
	} {
		t.Run(
			fmt.Sprintf("g.successorIDs(%d) == %v", tt.id, tt.successors),
			func(t *testing.T) {
				if diff := orderInsensitiveDiff(slices.Collect(g.successorIDs(tt.id)), tt.successors); diff != "" {
					t.Errorf("mismatch (-g.successorIDs(%d) +expected):\n%s", tt.id, diff)
				}
			},
		)
	}
}

func testRemoveEdgeBetweenNodeIDs2(t *testing.T, g *graph2) {
	t.Run("g.removeEdgeBetweenNodeIDs(3, 2)", func(t *testing.T) {
		g.removeEdgeBetweenNodeIDs(3, 2)
		if diff := orderInsensitiveDiff(slices.Collect(g.successorIDs(3)), []int32{4, 5}); diff != "" {
			t.Errorf("mismatch (-g.successorIDs(3) +expected):\n%s", diff)
		}
	})

}
