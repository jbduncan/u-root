// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"iter"
	"slices"
)

func newGraph2() *graph2 {
	return &graph2{
		nodeToID:             make(map[string]int32),
		idToNode:             make([]string, 0),
		nodeIDToSuccessorIDs: make([][]int32, 0),
	}
}

type graph2 struct {
	nodeToID             map[string]int32
	idToNode             []string
	nodeIDToSuccessorIDs [][]int32
}

func (g *graph2) addNode(node string) {
	_ = g.addNodeInternal(node)
}

func (g *graph2) addNodeInternal(node string) int32 {
	if id, ok := g.nodeToID[node]; ok {
		return id
	}

	id := int32(len(g.idToNode))
	g.idToNode = append(g.idToNode, node)
	g.nodeToID[node] = id

	g.nodeIDToSuccessorIDs = append(g.nodeIDToSuccessorIDs, make([]int32, 0))

	return id
}

func (g *graph2) putEdge(source, target string) {
	sourceID := g.addNodeInternal(source)
	targetID := g.addNodeInternal(target)

	succs := g.nodeIDToSuccessorIDs[sourceID]
	if !slices.Contains(succs, targetID) {
		g.nodeIDToSuccessorIDs[sourceID] = append(succs, targetID)
	}
}

func (g *graph2) valueFor(nodeID int32) string {
	return g.idToNode[nodeID]
}

func (g *graph2) nodeCount() int {
	return len(g.nodeIDToSuccessorIDs)
}

func (g *graph2) nodeIDs() iter.Seq[int32] {
	return func(yield func(int32) bool) {
		for id := range len(g.idToNode) {
			if !yield(int32(id)) {
				return
			}
		}
	}
}

func (g *graph2) successorIDs(nodeID int32) iter.Seq[int32] {
	return slices.Values(g.nodeIDToSuccessorIDs[nodeID])
}

func (g *graph2) removeEdgeBetweenNodeIDs(sourceID, targetID int32) {
	succs := g.nodeIDToSuccessorIDs[sourceID]
	idx := slices.Index(succs, targetID)
	g.nodeIDToSuccessorIDs[sourceID] = slices.Delete(succs, idx, idx+1)
}
