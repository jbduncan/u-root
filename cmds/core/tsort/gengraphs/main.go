// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/cmds/core/tsort/gen"
)

func main() {
	cases := []struct {
		name string
		g    string
	}{
		{
			name: "small-sparse-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(10, 0.1),
		},
		{
			name: "small-half-total-edges-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(10, 0.5),
		},
		{
			name: "small-edgeless-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(10, 0.0),
		},
		{
			name: "small-tournament-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(10, 1.0),
		},
		{
			name: "medium-sparse-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(100, 0.1),
		},
		{
			name: "medium-half-total-edges-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(100, 0.5),
		},
		{
			name: "medium-edgeless-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(100, 0),
		},
		{
			name: "medium-tournament-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(100, 1.0),
		},
		{
			name: "large-sparse-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(1_000, 0.1),
		},
		{
			name: "large-half-total-edges-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(1_000, 0.5),
		},
		{
			name: "large-edgeless-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(1_000, 0.0),
		},
		{
			name: "large-tournament-directed-acyclic-graph.txt",
			g:    gen.RandomDirectedAcyclicGraph(1_000, 1.0),
		},
		{
			name: "small-cyclic-graph.txt",
			g:    gen.RandomDirectedCyclicGraph(10),
		},
		{
			name: "medium-cyclic-graph.txt",
			g:    gen.RandomDirectedCyclicGraph(50),
		},
		{
			name: "large-cyclic-graph.txt",
			g:    gen.RandomDirectedCyclicGraph(100),
		},
	}
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	for _, c := range cases {
		func() {
			f, err := os.Create(filepath.Join(tempDir, c.name))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if _, err = f.WriteString(c.g); err != nil {
				panic(err)
			}
		}()
	}
	fmt.Println(tempDir)
}
