package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/dominikbraun/graph"
	"math"
	"math/rand"
)

const Lambda = 128

type Prover[K comparable, T any] struct {
	graph graph.Graph[K, T]
}

// produces vertices from graph from prover
func (p *Prover[K, T]) vertices() ([]K, error) {
	edges, err := p.graph.Edges()
	vertexSet := make(map[K]struct{}) //dictionary
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

func (p *Prover[K, T]) colors(vertices []K) ([]string, error) {
	colorSet := make(map[string]struct{})
	for _, v := range vertices {
		_, properties, err := p.graph.VertexWithProperties(v)
		if err != nil {
			return []string{}, err
		}
		colorSet[properties.Attributes["color"]] = struct{}{} //Maps to nothing in Go
	}
	colors := make([]string, 0, len(colorSet))
	for color := range colorSet {
		colors = append(colors, color)
	}
	return colors, nil
}

// random shuffle
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

type Commitment struct {
	hash [32]byte
	r    int64
}

func (p *Prover[K, T]) commitments() ([]Commitment, error) {
	vertices, err := p.vertices()
	if err != nil {
		return []Commitment{}, err
	}

	colors, err := p.colors(vertices)
	if err != nil {
		return []Commitment{}, err
	}
	perm, err := p.colorPermutation(colors)

	//generates commitments

	commitments := make([]Commitment, 0, len(vertices))

	for _, v := range vertices {
		_, properties, err := p.graph.VertexWithProperties(v)
		if err != nil {
			return []Commitment{}, err
		}
		color := properties.Attributes["color"]

		permColor := perm[color]
		r := rand.Int63n(int64(math.Pow(2, Lambda)))
		rBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(rBytes, uint64(r))
		hash := sha256.Sum256([]byte(permColor + string(rBytes)))
		c := Commitment{hash: hash, r: r}
		commitments = append(commitments, c)
	}
	return commitments, nil
}

func main() {
	prover := Prover[int, int]{}
	prover.graph = graph.New(graph.IntHash)
	prover.graph.AddVertex(0, graph.VertexAttribute("color", "red"))
	prover.graph.AddVertex(1, graph.VertexAttribute("color", "blue"))
	prover.graph.AddVertex(2, graph.VertexAttribute("color", "green"))

	prover.graph.AddEdge(0, 1)
	prover.graph.AddEdge(1, 2)
	prover.graph.AddEdge(2, 3)

	x := 0
	x = 5

	fmt.Println(x)

}
