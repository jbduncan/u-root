// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"slices"
)

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
	successors []string
}

func (g *graph) addNode(node string) {
	if _, ok := g.nodeToData[node]; !ok {
		g.nodeToData[node] = &nodeData{
			inDegree:   0,
			successors: make([]string, 0),
		}
	}
}

func (g *graph) putEdge(source, target string) {
	g.addNode(source)
	g.addNode(target)

	g.nodeToData[source].successors = append(g.nodeToData[source].successors, target)
	g.nodeToData[target].inDegree++
}

func (g *graph) successors(node string) []string {
	n, ok := g.nodeToData[node]
	if !ok {
		panic("node is not in graph")
	}

	return n.successors
}

func (g *graph) removeEdge(source, target string) {
	if _, ok := g.nodeToData[source]; !ok {
		panic("source node is not in graph")
	}
	if _, ok := g.nodeToData[target]; !ok {
		panic("target node is not in graph")
	}

	index := slices.Index(g.nodeToData[source].successors, target)
	g.nodeToData[source].successors =
		slices.Delete(g.nodeToData[source].successors, index, index+1)
	g.nodeToData[target].inDegree--
}

func (g *graph) inDegree(node string) int {
	n, ok := g.nodeToData[node]
	if !ok {
		return 0
	}
	return n.inDegree
}
