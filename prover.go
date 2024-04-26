package main

import (
	"math"
	"math/rand"

	"github.com/dominikbraun/graph"
)

// Prover represents a prover in the graph coloring protocol.
type Prover[K comparable, T any] struct {
	graph       graph.Graph[K, T]
	commitments map[K]Commitment
}

// newProver creates a new prover instance for the given graph.
func newProver[K comparable, T any](g graph.Graph[K, T]) (Prover[K, T], error) {
	p := Prover[K, T]{graph: g}

	vertices, err := p.vertices()
	if err != nil {
		return p, err
	}

	colors, err := p.colors(vertices)
	if err != nil {
		return p, err
	}
	perm, err := p.colorPermutation(colors)

	p.commitments = make(map[K]Commitment)
	for _, v := range vertices {
		_, properties, err := p.graph.VertexWithProperties(v)
		if err != nil {
			return Prover[K, T]{}, err
		}
		color := properties.Attributes["color"]

		permColor := perm[color]
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

func (p *Prover[K, T]) hashes() map[K][32]byte {
	hashes := make(map[K][32]byte)
	for vertex, commitment := range p.commitments {
		hashes[vertex] = commitment.hash
	}
	return hashes
}

// vertices returns all vertices of the graph that are connected by an edge.
// Note: this function does not return isolated vertices, but this is not
// relevant for the protocol.
func (p *Prover[K, T]) vertices() ([]K, error) {
	edges, err := p.graph.Edges()
	vertexSet := make(map[K]struct{})
	vertices := make([]K, 0, len(vertexSet))
	if err != nil {
		return vertices, err
	}

	for _, edge := range edges {
		vertexSet[edge.Source] = struct{}{}
		vertexSet[edge.Target] = struct{}{}
	}
	for vertex := range vertexSet {
		vertices = append(vertices, vertex)
	}
	return vertices, nil
}

// colors returns all colors of the vertices in the given slice.
func (p *Prover[K, T]) colors(vertices []K) ([]string, error) {
	colorSet := make(map[string]struct{})
	for _, v := range vertices {
		_, properties, err := p.graph.VertexWithProperties(v)
		if err != nil {
			return []string{}, err
		}
		colorSet[properties.Attributes["color"]] = struct{}{}
	}
	colors := make([]string, 0, len(colorSet))
	for color := range colorSet {
		colors = append(colors, color)
	}
	return colors, nil
}

// colorPermutation returns a random permutation of the given colors.
// The permutation is represented as a map where the key is the original
// color and the value is the new color.
func (p *Prover[K, T]) colorPermutation(colors []string) (map[string]string, error) {
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
	return perm, nil
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
func (p *Prover[K, T]) openCommitments(e graph.Edge[K]) (oc1 OpenCommitment, oc2 OpenCommitment) {
	oc1 = fromCommitment(p.commitments[e.Source])
	oc2 = fromCommitment(p.commitments[e.Target])
	return oc1, oc2
}
