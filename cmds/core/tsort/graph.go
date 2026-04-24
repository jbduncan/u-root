// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"iter"
	"maps"
)

func newGraph() *graph {
	return &graph{
		nodeToData: make(map[str]nodeData),
	}
}

type graph struct {
	nodeToData map[str]nodeData
}

type nodeData struct {
	inDegree   int
	successors set[str]
}

func (g *graph) addNode(node str) {
	g.addNodeInternal(node)
}

func (g *graph) addNodeInternal(node str) nodeData {
	var data nodeData
	var ok bool
	if data, ok = g.nodeToData[node]; !ok {
		data = nodeData{
			inDegree:   0,
			successors: makeSet[str](),
		}
		g.nodeToData[node] = data
	}
	return data
}

func (g *graph) putEdge(source, target str) {
	sourceData := g.addNodeInternal(source)
	targetData := g.addNodeInternal(target)

	if !sourceData.successors.has(target) {
		sourceData.successors.add(target)
		g.nodeToData[target] = nodeData{
			inDegree:   targetData.inDegree + 1,
			successors: targetData.successors,
		}
	}
}

func (g *graph) nodeCount() int {
	return len(g.nodeToData)
}

func (g *graph) nodes() iter.Seq[str] {
	return maps.Keys(g.nodeToData)
}

func (g *graph) inDegree(node str) int {
	data, ok := g.nodeToData[node]
	if !ok {
		return 0
	}
	return data.inDegree
}

func (g *graph) successors(node str) iter.Seq[str] {
	data, ok := g.nodeToData[node]
	if !ok {
		panic("node is not in graph")
	}

	return data.successors.all()
}

func (g *graph) removeNode(node str) {
	if _, ok := g.nodeToData[node]; !ok {
		panic("node is not in graph")
	}

	// In a general-purpose graph type, removing a node would also require
	// updating all predecessors' successor lists. This would be O(N) since
	// we don't track predecessors, so every node would need to be looped over
	// to figure out which ones are the predecessors.
	//
	// However, this optimization is safe for the variant of Kahn's algorithm
	// used in tsort.go because we follow a strict dequeue-process-delete
	// pattern (see decreaseInDegree for more information). We never query the
	// predecessors of the remaining nodes (only successors, which are
	// tracked). So even though the graph state becomes incomplete, it suffices
	// for the algorithm's needs.
	delete(g.nodeToData, node)
}

func (g *graph) removeEdge(source, target str) {
	sourceData, ok := g.nodeToData[source]
	if !ok {
		panic("source node is not in graph")
	}
	targetData, ok := g.nodeToData[target]
	if !ok {
		panic("target node is not in graph")
	}

	sourceData.successors.remove(target)
	g.nodeToData[target] = nodeData{
		inDegree:   targetData.inDegree - 1,
		successors: targetData.successors,
	}
}

func (g *graph) decreaseInDegree(node str) int {
	// Optimization for the variant of Kahn's algorithm used in tsort.go.
	// Unlike the standard algorithm (which removes edges), this variant
	// decreases in-degrees and immediately removes processed nodes. These
	// design choices enable early cycle detection: when the queue empties but
	// nodes remain, a cycle exists.
	//
	// The decreaseInDegree optimization is safe because of these choices:
	// 1. We dequeue a root node (in-degree 0)
	// 2. Decrement in-degrees of its successors  <- we are here
	// 3. Immediately remove the dequeued node via removeNode()
	//
	// Since we delete dequeued nodes before the next iteration, we never
	// query predecessors of remaining nodes. In-degree values alone suffice
	// to identify new roots. This avoids the O(N) cost of maintaining
	// predecessor links (which this graph implementation doesn't track).
	//
	// WARNING: This optimization depends on the immediate dequeue-process-delete
	// pattern. Code that breaks this pattern (e.g., deferring node removal or
	// querying predecessors) will cause algorithm failure.

	data, ok := g.nodeToData[node]
	if !ok {
		panic("target node is not in graph")
	}

	newInDegree := data.inDegree - 1
	g.nodeToData[node] = nodeData{
		inDegree:   newInDegree,
		successors: data.successors,
	}

	return newInDegree
}
