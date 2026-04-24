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
		nodeToSuccessors: make(map[string]set),
	}
}

type graph struct {
	nodeToSuccessors map[string]set
}

func (g *graph) addNode(node string) {
	_ = g.addNodeInternal(node)
}

func (g *graph) addNodeInternal(node string) set {
	data, ok := g.nodeToSuccessors[node]
	if !ok {
		data = makeSet()
		g.nodeToSuccessors[node] = data
	}

	return data
}

func (g *graph) putEdge(source, target string) {
	sourceData := g.addNodeInternal(source)
	_ = g.addNodeInternal(target)

	sourceData.add(target)
}

func (g *graph) nodeCount() int {
	return len(g.nodeToSuccessors)
}

func (g *graph) nodes() iter.Seq[string] {
	return maps.Keys(g.nodeToSuccessors)
}

func (g *graph) successors(node string) iter.Seq[string] {
	return g.nodeToSuccessors[node].all()
}

func (g *graph) removeEdge(source, target string) {
	g.nodeToSuccessors[source].remove(target)
}
