// Copyright 2012-2024 the u-root Authors. All rights reserved
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
	"iter"
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
			fmt.Fprintf(stdout, "%v\n", node)
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
	var i int
	var odd bool

	next := func() (string, bool) {
		if i == len(fields) {
			return "", false
		}
		odd = !odd
		result := fields[i]
		i++
		return result, true
	}

	for {
		a, ok := next()
		if !ok {
			break
		}

		b, ok := next()
		if !ok {
			break
		}

		if a == b {
			g.addNode(a)
		} else {
			g.putEdge(a, b)
		}
	}

	if odd {
		return errOddDataCount
	}

	return nil
}

func topologicalOrdering(
	g *graph,
	f func(node string),
	cycles func(cycle []string),
) {
	for {
		// Kahn's algorithm
		var result []string
		roots := rootsOf(g)
		nonRoots := nonRootsOf(g)
		for !roots.isEmpty() {
			next := roots.dequeue()
			result = append(result, next)
			for succ := range g.successors(next) {
				nonRoots.removeOne(succ)
				if !nonRoots.has(succ) {
					roots.enqueue(succ)
				}
			}
		}
		if nonRoots.isEmpty() {
			// No cycles left
			for _, value := range result {
				f(value)
			}
			break
		}

		// Break a cycle and try Kahn's algorithm again
		cycle := cycleStartingAtAnyOfNodes(g, nonRoots.allUnique())
		g.removeEdge(cycle[len(cycle)-1], cycle[0])
		cycles(cycle)
	}
}

func rootsOf(g *graph) queue {
	result := queue{}
	for node := range g.allNodes() {
		if g.inDegree(node) == 0 {
			result.enqueue(node)
		}
	}
	return result
}

func nonRootsOf(g *graph) multiset {
	result := newMultiset()
	for node := range g.allNodes() {
		if g.inDegree(node) > 0 {
			result.add(node, g.inDegree(node))
		}
	}
	return result
}

type nodeState int

const (
	notVisited nodeState = iota
	partiallyVisited
	fullyVisited
)

func (n nodeState) String() string {
	switch n {
	case notVisited:
		return "notVisited"
	case partiallyVisited:
		return "partiallyVisited"
	case fullyVisited:
		return "fullyVisited"
	}
	panic(fmt.Sprintf("unknown nodeState: %d", n))
}

func cycleStartingAtAnyOfNodes(g *graph, nodes iter.Seq[string]) []string {
	var path []string
	nodeToState := make(map[string]nodeState)

	var findCycle func(node string) []string
	findCycle = func(node string) []string {
		nodeToState[node] = partiallyVisited
		path = append(path, node)

		for succ := range g.successors(node) {
			if nodeToState[succ] == partiallyVisited {
				// cycle found
				start := slices.Index(path, succ)
				cycle := slices.Clone(path[start:])
				return cycle
			}
			if nodeToState[succ] == notVisited {
				if cycle := findCycle(succ); len(cycle) != 0 {
					return cycle
				}
			}
		}

		path = path[:len(path)-1]
		nodeToState[node] = fullyVisited
		return nil
	}

	for node := range nodes {
		if nodeToState[node] != notVisited {
			continue
		}
		if cycle := findCycle(node); len(cycle) != 0 {
			return cycle
		}
	}
	return nil
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
