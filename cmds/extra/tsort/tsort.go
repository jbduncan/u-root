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

var errNonFatal = errors.New("non-fatal")
var errOddDataCount = errors.New("odd data count")

func run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) error {
	in := io.NopCloser(stdin)
	if len(args) >= 1 {
		var err error
		if in, err = os.Open(args[0]); err != nil {
			return err
		}
	}
	defer in.Close()

	var buf strings.Builder
	if _, err := io.Copy(&buf, in); err != nil {
		return err
	}

	g := newGraph()

	if err := parseInto(buf.String(), g); err != nil {
		return err
	}

	var cycleFound bool
	for nodes, cycle := range topologicalOrdering(g) {
		for _, node := range nodes {
			fmt.Fprintf(stdout, "%v\n", node)
		}
		if cycle != nil {
			fmt.Fprintf(stderr, "tsort: %v\n", "cycle in data")
			for _, node := range cycle {
				fmt.Fprintf(stderr, "tsort: %v\n", node)
			}
			cycleFound = true
		}
	}
	if cycleFound {
		return errNonFatal
	}
	return nil
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

type nodes []string
type cycle []string

func topologicalOrdering(g *graph) iter.Seq2[nodes, cycle] {
	return func(yield func(nodes, cycle) bool) {
		for {
			// Kahn's algorithm
			var result nodes
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
				yield(result, nil)
				return
			}

			// Break a cycle and try Kahn's algorithm again
			for next := range nonRoots.allUnique() {
				cycle := cycleStartingAt(g, next)
				if len(cycle) == 0 {
					continue
				}

				g.removeEdge(cycle[len(cycle)-1], cycle[0])
				if !yield(nil, cycle) {
					return
				}
				break
			}
		}
	}
}

func rootsOf(g *graph) queue {
	result := queue{}
	for node := range g.nodes() {
		if g.inDegree(node) == 0 {
			result.enqueue(node)
		}
	}
	return result
}

func nonRootsOf(g *graph) multiset {
	result := newMultiset()
	for node := range g.nodes() {
		if g.inDegree(node) > 0 {
			result.add(node, g.inDegree(node))
		}
	}
	return result
}

func cycleStartingAt(g *graph, node string) cycle {
	s := makeStack()
	s.push(node)
	inStack := makeSet()
	inStack.add(node)

	var cycle cycle
	var dfs func() bool
	dfs = func() bool {
		for succ := range g.successors(s.peek()) {
			if inStack.has(succ) {
				// cycle found
				cycle = append(cycle, s.pop())
				for cycle[len(cycle)-1] != succ {
					cycle = append(cycle, s.pop())
				}
				slices.Reverse(cycle)
				return true
			}

			s.push(succ)
			inStack.add(succ)
			if dfs() {
				return true
			}
		}

		inStack.remove(s.pop())
		return false
	}
	dfs()
	return cycle
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
