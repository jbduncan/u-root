// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func newGraph() *graph {
	return &graph{
		nodeToData: make(map[string]*nodeData),
	}
}

type graph struct {
	nodeToData map[string]*nodeData
}

type nodeData struct {
	inDegree   int
	successors set
}

func (g *graph) addNode(node string) {
	g.addNodeInternal(node)
}

func (g *graph) putEdge(source, target string) {
	sourceData := g.addNodeInternal(source)
	targetData := g.addNodeInternal(target)

	successors := sourceData.successors
	if !successors.has(target) {
		successors.add(target)
		targetData.inDegree++
	}
}

func (g *graph) addNodeInternal(node string) *nodeData {
	data, ok := g.nodeToData[node]
	if !ok {
		data = &nodeData{
			inDegree:   0,
			successors: makeSet(),
		}
		g.nodeToData[node] = data
	}
	return data
}

func (g *graph) successors(node string) set {
	data, ok := g.nodeToData[node]
	if !ok {
		panic("node is not in graph")
	}

	return data.successors
}

func (g *graph) removeEdge(source, target string) {
	if _, ok := g.nodeToData[source]; !ok {
		panic("source node is not in graph")
	}
	if _, ok := g.nodeToData[target]; !ok {
		panic("target node is not in graph")
	}

	g.nodeToData[source].successors.remove(target)
	g.nodeToData[target].inDegree--
}

func (g *graph) inDegree(node string) int {
	data, ok := g.nodeToData[node]
	if !ok {
		return 0
	}
	return data.inDegree
}
