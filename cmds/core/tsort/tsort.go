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
	"bufio"
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
	out := bufio.NewWriter(stdout)
	defer out.Flush()

	var err error
	in := io.NopCloser(stdin)
	if len(args) >= 1 {
		if in, err = os.Open(args[0]); err != nil {
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
		func(node nodeID) {
			// More optimal than fmt.Fprintln
			_, _ = out.WriteString(g.valueFor(node))
			_ = out.WriteByte('\n')
		},
		func(cycle []nodeID) {
			_, _ = fmt.Fprintln(stderr, "tsort: cycle in data")
			for _, node := range cycle {
				_, _ = fmt.Fprintln(stderr, "tsort: "+g.valueFor(node))
			}
			err = errNonFatal
		})
	return err
}

func parseInto(buf string, g *graph) error {
	var prev string
	var odd bool
	for field := range strings.FieldsSeq(buf) {
		if odd {
			if prev == field {
				g.addNode(prev)
			} else {
				g.putEdge(prev, field)
			}
		} else {
			prev = field
		}

		odd = !odd
	}

	if odd {
		return errOddDataCount
	}
	return nil
}

func topologicalOrdering(
	g *graph,
	f func(node nodeID),
	cycles func(cycle []nodeID),
) {
	// A topological ordering algorithm based on the depth-first search
	// algorithm in "Introduction to Algorithms" by Cormen et al.
	//
	// Unlike normal topological ordering, it returns an ordering even for
	// cyclic graphs by breaking cycles when they are found.

	type visitState int8
	const (
		notVisited visitState = iota
		partiallyVisited
		fullyVisited
	)

	var path []nodeID
	result := make([]nodeID, 0, g.nodeCount())
	nodeToVisitState := make([]visitState, g.nodeCount())

	var doTopologicalOrdering func(node nodeID)
	doTopologicalOrdering = func(node nodeID) {
		nodeToVisitState[node] = partiallyVisited
		path = append(path, node)

		for succ := range g.successorIDs(node) {
			switch nodeToVisitState[succ] {
			case notVisited:
				doTopologicalOrdering(succ)
			case partiallyVisited:
				// Cycle detected; report it, break it and
				// continue as if the cycle never existed.
				idx := slices.Index(path, succ)
				cycle := path[idx:]
				cycles(cycle)
				g.removeEdge(node, succ)
			case fullyVisited:
				continue
			}
		}

		path = path[:len(path)-1]
		nodeToVisitState[node] = fullyVisited

		result = append(result, node)
	}

	for node := range g.nodeIDs() {
		if nodeToVisitState[node] != fullyVisited {
			doTopologicalOrdering(node)
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
