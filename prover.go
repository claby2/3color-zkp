package main

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

// Prover represents a prover in the graph coloring protocol.
type Prover struct {
	graph       graph.Graph
	coloring    map[int]string
	commitments map[int]Commitment
}

// newProver creates a new prover instance for the given graph.
func newProver(g graph.Graph, coloring map[int]string) (Prover, error) {
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

	p.commitments = make(map[int]Commitment)
	for _, v := range vertices {
		permColor := permutation[coloring[v]]

		r := rand.Int63n(int64(math.Pow(2, Lambda)))
		hash := commit(permColor, r)
		c := Commitment{
			color: permColor,
			hash:  hash, r: r,
		}
		p.commitments[v] = c
	}
	return p, nil
}

func (p *Prover) hashes() map[int][32]byte {
	hashes := make(map[int][32]byte)
	for vertex, commitment := range p.commitments {
		hashes[vertex] = commitment.hash
	}
	return hashes
}

// Commitment represents a commitment to a color and a random value.
type Commitment struct {
	color string
	hash  [32]byte
	r     int64
}

// OpenCommitment represents an open commitment to a color and a random value.
type OpenCommitment struct {
	color string
	r     int64
}

// fromCommitment creates an open commitment from a commitment.
func fromCommitment(c Commitment) OpenCommitment {
	return OpenCommitment{
		color: c.color,
		r:     c.r,
	}
}

// openCommitments returns the open commitments of the vertices connected by the given edge.
func (p *Prover) openCommitments(a, b int) (oc1 OpenCommitment, oc2 OpenCommitment) {
	oc1 = fromCommitment(p.commitments[a])
	oc2 = fromCommitment(p.commitments[b])
	return oc1, oc2
}
