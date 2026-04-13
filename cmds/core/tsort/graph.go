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

	// In a general-purpose graph type, the predecessors and successors of
	// the given node would need to be amended too. But this type only tracks
	// in-degrees and successors, so it would take O(N) time to find all of the
	// predecessors and amend them. Therefore, this method "cheats" and only
	// removes the given node, which is good enough for tsort.
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
