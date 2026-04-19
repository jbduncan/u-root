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
		nodeToData: make(map[string]nodeData),
	}
}

type graph struct {
	nodeToData map[string]nodeData
}

type nodeData struct {
	inDegree   int
	successors set
}

func (g *graph) addNode(node string) {
	g.addNodeInternal(node)
}

func (g *graph) addNodeInternal(node string) nodeData {
	var data nodeData
	var ok bool
	if data, ok = g.nodeToData[node]; !ok {
		data = nodeData{
			inDegree:   0,
			successors: makeSet(),
		}
		g.nodeToData[node] = data
	}
	return data
}

func (g *graph) putEdge(source, target string) {
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

func (g *graph) nodes() iter.Seq[string] {
	return maps.Keys(g.nodeToData)
}

func (g *graph) inDegree(node string) int {
	data, ok := g.nodeToData[node]
	if !ok {
		return 0
	}
	return data.inDegree
}

func (g *graph) successors(node string) iter.Seq[string] {
	data, ok := g.nodeToData[node]
	if !ok {
		panic("node is not in graph")
	}

	return data.successors.all()
}

func (g *graph) removeNode(node string) {
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

func (g *graph) removeEdge(source, target string) {
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
