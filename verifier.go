package main

import (
	"math/rand"

	"github.com/dominikbraun/graph"
)

// Verifier represents a verifier in the graph coloring protocol.
type Verifier[K comparable, T any] struct {
	graph  graph.Graph[K, T]
	hashes map[K][32]byte
}

// newVerifier creates a new verifier instance for the given graph.
func newVerifier[K comparable, T any](g graph.Graph[K, T], h map[K][32]byte) Verifier[K, T] {
	return Verifier[K, T]{graph: g, hashes: h}
}

// randomEdge returns a random edge from the graph.
func (v *Verifier[K, T]) randomEdge() (graph.Edge[K], error) {
	edges, err := v.graph.Edges()
	if err != nil {
		return graph.Edge[K]{}, err
	}
	return edges[rand.Intn(len(edges))], nil
}

// verify checks whether the given edge and its open commitments are valid.
func (v *Verifier[K, T]) verify(edge graph.Edge[K], oc1 OpenCommitment, oc2 OpenCommitment) bool {
	return v.hashes[edge.Source] == commit(oc1.color, oc1.r) &&
		v.hashes[edge.Target] == commit(oc2.color, oc2.r) &&
		oc1.color != oc2.color
}
