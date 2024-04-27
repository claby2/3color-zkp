package party

import (
	"math/rand"

	"github.com/claby2/3color-zkp/graph"
)

// Verifier represents a verifier in the graph coloring protocol.
type Verifier struct {
	graph  graph.Graph
	hashes map[int][32]byte
}

// NewVerifier creates a new verifier instance for the given graph.
func NewVerifier(g graph.Graph, h map[int][32]byte) Verifier {
	return Verifier{graph: g, hashes: h}
}

// RandomEdge returns a random edge from the graph.
func (v *Verifier) RandomEdge() (int, int) {
	edges := v.graph.Edges()
	edge := edges[rand.Intn(len(edges))]
	return edge[0], edge[1]
}

func acceptableColor(color string) bool {
	return color == "red" || color == "green" || color == "blue"
}

// Verify checks whether the given edge and its open commitments are valid.
func (v *Verifier) Verify(a, b int, oc1 OpenCommitment, oc2 OpenCommitment) bool {
	return v.hashes[a] == commit(oc1.Color, oc1.R) &&
		v.hashes[b] == commit(oc2.Color, oc2.R) &&
		oc1.Color != oc2.Color &&
		acceptableColor(oc1.Color) &&
		acceptableColor(oc2.Color)
}
