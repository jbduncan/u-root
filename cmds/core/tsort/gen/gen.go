package gen

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

var (
	rnd = rand.New(rand.NewPCG(1, 1))
)

func SeqFrom0To(size uint) string {
	result := new(strings.Builder)
	for i := uint(0); i < size-1; i++ {
		_, _ = fmt.Fprintln(result, i, i+1)
	}
	return result.String()
}

func RandomDirectedAcyclicGraph(
	nodeCount uint16,
	edgeCountRatio float64,
) string {
	if edgeCountRatio < 0.0 || edgeCountRatio > 1.0 {
		panic(fmt.Sprintf(
			"edgeCountRatio %v must be between 0.0 and 1.0",
			edgeCountRatio,
		))
	}

	totalPossibleEdges := maxEdgesForDirectedAcyclicGraph(nodeCount)
	edgeCount := uint(math.Round(float64(totalPossibleEdges) * edgeCountRatio))

	// filled with `false` by default
	randomEdges := make([]bool, totalPossibleEdges)
	for i := range edgeCount {
		randomEdges[i] = true
	}
	rnd.Shuffle(len(randomEdges), func(i, j int) {
		randomEdges[i], randomEdges[j] = randomEdges[j], randomEdges[i]
	})

	result := new(strings.Builder)
	for i := uint16(0); i < nodeCount; i++ {
		_, _ = fmt.Fprintln(result, i, i)
	}
	index := 0
	for i := uint16(0); i < nodeCount-1; i++ {
		for j := i + 1; j < nodeCount; j++ {
			if randomEdges[index] {
				_, _ = fmt.Fprintln(result, i, j)
			}
			index++
		}
	}

	return result.String()
}

// For any directed acyclic graph, the maximum number of edges is equal to (n * (n - 1) / 2),
// where n is the number of nodes in the graph.
func maxEdgesForDirectedAcyclicGraph(nodeCount uint16) uint {
	return uint(nodeCount) * (uint(nodeCount) - 1) / 2
}

func RandomDirectedCyclicGraph(nodeRange uint) string {
	result := new(strings.Builder)
	// Produces a cyclic graph with a fixed RNG seed and through
	// sheer probability.
	for range 100 * nodeRange {
		x := rnd.UintN(nodeRange + 1)
		y := rnd.UintN(nodeRange + 1)
		_, _ = fmt.Fprintln(result, x, y)
	}
	return result.String()
}
