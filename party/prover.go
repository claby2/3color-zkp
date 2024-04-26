package party

import (
	"math"
	"math/rand"

	"github.com/claby2/3color-zkp/graph"
)

// colorPermutation returns a random permutation of the given colors.
// The permutation is represented as a map where the key is the original
// color and the value is the new color.
func colorPermutation(colors []string) map[string]string {
	shuffled := make([]string, len(colors))
	copy(shuffled, colors)
	for i := range shuffled {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	perm := make(map[string]string)
	for i, color := range colors {
		perm[color] = shuffled[i]
	}
	return perm
}

// commitment represents a commitment to a color and a random value.
type commitment struct {
	color string
	hash  [32]byte
	r     int64
}

// OpenCommitment represents an open commitment to a color and a random value.
func (c commitment) Open() OpenCommitment {
	return OpenCommitment{Color: c.color, R: c.r}
}

// Prover represents a prover in the graph coloring protocol.
type Prover struct {
	graph       graph.Graph
	coloring    map[int]string
	commitments map[int]commitment
}

// NewProver creates a new prover instance for the given graph.
func NewProver(g graph.Graph, coloring map[int]string) (Prover, error) {
	p := Prover{graph: g}

	vertices := p.graph.Vertices()

	colorSet := make(map[string]struct{})
	for _, v := range vertices {
		colorSet[coloring[v]] = struct{}{}
	}

	colors := make([]string, 0, len(colorSet))
	for color := range colorSet {
		colors = append(colors, color)
	}

	permutation := colorPermutation(colors)

	p.commitments = make(map[int]commitment)
	for _, v := range vertices {
		permColor := permutation[coloring[v]]

		r := rand.Int63n(int64(math.Pow(2, Lambda)))
		hash := commit(permColor, r)
		c := commitment{
			color: permColor,
			hash:  hash, r: r,
		}
		p.commitments[v] = c
	}
	return p, nil
}

// Hashes returns the hashes of the commitments.
func (p *Prover) Hashes() map[int][32]byte {
	hashes := make(map[int][32]byte)
	for vertex, commitment := range p.commitments {
		hashes[vertex] = commitment.hash
	}
	return hashes
}

// Open returns the open commitments of the vertices connected by the given edge.
func (p *Prover) Open(a, b int) (oc1 OpenCommitment, oc2 OpenCommitment) {
	oc1 = p.commitments[a].Open()
	oc2 = p.commitments[b].Open()
	return oc1, oc2
}
