package main

import (
	"math/rand"

	"github.com/claby2/3color-zkp/graph"
)

// Verifier represents a verifier in the graph coloring protocol.
type Verifier struct {
	graph  graph.Graph
	hashes map[int][32]byte
}

// newVerifier creates a new verifier instance for the given graph.
func newVerifier(g graph.Graph, h map[int][32]byte) Verifier {
	return Verifier{graph: g, hashes: h}
}

// randomEdge returns a random edge from the graph.
func (v *Verifier) randomEdge() (int, int) {
	edges := v.graph.Edges()
	edge := edges[rand.Intn(len(edges))]
	return edge[0], edge[1]
}

// verify checks whether the given edge and its open commitments are valid.
func (v *Verifier) verify(a, b int, oc1 OpenCommitment, oc2 OpenCommitment) bool {
	return v.hashes[a] == commit(oc1.color, oc1.r) &&
		v.hashes[b] == commit(oc2.color, oc2.r) &&
		oc1.color != oc2.color
}
