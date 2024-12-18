package main

import (
	"iter"
	"slices"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func orderInsensitiveIterDiff(iter iter.Seq[string], values ...string) string {
	return orderInsensitiveDiff(slices.Collect(iter), values)
}

func orderInsensitiveDiff(a []string, b []string) string {
	return cmp.Diff(
		a, b, cmpopts.SortSlices(func(x, y string) bool { return x < y }))
}
