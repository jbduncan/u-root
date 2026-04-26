// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// tsort writes to standard output a totally ordered list of items consistent
// with a partial ordering of items contained in the input. The standard input
// will be used if no file is specified.
//
// The input is a sequence of pairs of items, separated by <blank> characters.
// Pairs of different items (e.g., "a b") indicate ordering. Pairs of identical
// items (e.g., "c c") indicate presence, but not ordering.
//
// Synopsis:
//
//	tsort [FILE]
//
// Example:
//
//	tsort <<EOF
//	a b c c d e
//	g g
//	f g e f
//	h h
//	EOF
//
// produces an output like:
//
//	a
//	b
//	c
//	d
//	e
//	f
//	g
//	h
//
// which is one valid total ordering, but this is not guaranteed, it could
// equally be:
//
//	h
//	a
//	c
//	d
//	b
//	e
//	f
//	g
//
// or any other ordering where the following holds true:
//
//	- a is before b
//	- d is before e
//	- f is before g
//	- e is before f
//	- c is anywhere
//	- h is anywhere

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

var (
	errNonFatal     = errors.New("non-fatal")
	errOddDataCount = errors.New("odd data count")
)

func run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) error {
	var err error
	in := io.NopCloser(stdin)
	if len(args) >= 1 {
		in, err = os.Open(args[0])
		if err != nil {
			return err
		}
	}
	defer in.Close()

	var buf strings.Builder
	if _, err = io.Copy(&buf, in); err != nil {
		return err
	}

	g := newGraph()
	if err = parseInto(buf.String(), g); err != nil {
		return err
	}

	topologicalOrdering(
		g,
		func(node string) {
			fmt.Fprintln(stdout, node)
		},
		func(cycle []string) {
			fmt.Fprintln(stderr, "tsort: cycle in data")
			for _, node := range cycle {
				fmt.Fprintf(stderr, "tsort: %v\n", node)
			}
			err = errNonFatal
		})
	return err
}

func parseInto(buf string, g *graph) error {
	fields := strings.Fields(buf)
	if len(fields)%2 == 1 {
		return errOddDataCount
	}

	for i := 0; i < len(fields); i += 2 {
		a, b := fields[i], fields[i+1]
		if a == b {
			g.addNode(a)
		} else {
			g.putEdge(a, b)
		}
	}

	return nil
}

func topologicalOrdering(
	g *graph,
	f func(node string),
	cycles func(cycle []string),
) {
	// A topological ordering algorithm based on the depth-first search
	// algorithm by Cormen et al. It returns an ordering even for graphs with
	// cycles by breaking cycles when they are found.

	type visitState int
	const (
		notVisited visitState = iota
		partiallyVisited
		fullyVisited
	)

	type stackFrame struct {
		node       string
		succs      []string
		idx        int
		finalizing bool
	}

	var path []string
	var stack []stackFrame
	result := make([]string, 0, g.nodeCount())
	nodeToVisitState := make(map[string]visitState, g.nodeCount())

	topologicalOrderingStartingFrom := func(node string) {
		path = append(path, node)
		stack = append(stack,
			stackFrame{
				node:       node,
				finalizing: true,
			},
			stackFrame{
				node:       node,
				succs:      slices.Collect(g.successors(node)),
				idx:        0,
				finalizing: false,
			},
		)
		nodeToVisitState[node] = partiallyVisited

		for len(stack) > 0 {
			frame := stack[len(stack)-1]

			if frame.finalizing {
				path = path[:len(path)-1]
				stack = stack[:len(stack)-1]
				result = append(result, frame.node)
				nodeToVisitState[frame.node] = fullyVisited
				continue
			}

			if frame.idx == len(frame.succs) {
				stack = stack[:len(stack)-1]
				continue
			}

			succ := frame.succs[frame.idx]
			switch nodeToVisitState[succ] {
			case notVisited:
				path = append(path, succ)
				stack = append(
					stack,
					stackFrame{
						node:       succ,
						finalizing: true,
					},
					stackFrame{
						node:       succ,
						succs:      slices.Collect(g.successors(succ)),
						idx:        0,
						finalizing: false,
					},
				)
				nodeToVisitState[succ] = partiallyVisited
			case partiallyVisited:
				// Cycle detected; report it, break it and
				// continue as if the cycle never existed.
				idx := slices.Index(path, succ)
				cycle := path[idx:]
				cycles(cycle)
				g.removeEdge(frame.node, succ)
				stack[len(stack)-1].idx++
			case fullyVisited:
				stack[len(stack)-1].idx++
			}
		}
	}

	for node := range g.nodes() {
		if nodeToVisitState[node] != fullyVisited {
			topologicalOrderingStartingFrom(node)
		}
	}

	for _, node := range slices.Backward(result) {
		f(node)
	}
}

func main() {
	err := run(os.Stdin, os.Stdout, os.Stderr, os.Args[1:]...)
	if errors.Is(err, errNonFatal) {
		// All non-fatal warnings have been printed already, so just exit.
		os.Exit(1)
	}
	if err != nil {
		log.Fatalf("tsort: %v", err)
	}
}
