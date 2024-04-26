package graph

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidGraph is returned when the input graph is invalid.
var (
	ErrInvalidGraph    = errors.New("invalid graph")
	ErrInvalidColoring = errors.New("invalid coloring")
)

// Parse a graph from a string, producing a graph and a coloring.
func Parse(s string) (Graph, map[int]string, error) {
	lines := strings.Split(s, "\n")
	g := New()
	coloring := make(map[int]string)

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			v1, err := strconv.Atoi(parts[0])
			if err != nil {
				return g, coloring, err
			}
			g.AddVertex(v1)

			v2, err := strconv.Atoi(parts[1])
			if err != nil {
				coloring[v1] = parts[1]
			} else {
				g.AddEdge(v1, v2)
			}
		} else {
			return g, coloring, ErrInvalidGraph
		}
	}

	if err := g.Verify(); err != nil {
		return g, coloring, err
	}

	return g, coloring, nil
}

// Graph represents a graph with vertices and edges.
type Graph struct {
	vertices map[int]struct{}
	edges    map[int]map[int]struct{}
}

// New creates a new graph instance.
func New() Graph {
	return Graph{
		vertices: map[int]struct{}{},
		edges:    map[int]map[int]struct{}{},
	}
}

// Verify ensures edges are bidirectional and all vertices mentioned in edges
// are present in the graph.
func (g *Graph) Verify() error {
	for v1, v2s := range g.edges {
		if _, ok := g.vertices[v1]; !ok {
			return ErrInvalidGraph
		}

		for v2 := range v2s {
			if _, ok := g.edges[v2]; !ok {
				return ErrInvalidGraph
			}
			if _, ok := g.edges[v2][v1]; !ok {
				return ErrInvalidGraph
			}

			if _, ok := g.vertices[v2]; !ok {
				return ErrInvalidGraph
			}
		}
	}
	return nil
}

// AddVertex adds a vertex to the graph.
func (g *Graph) AddVertex(v int) {
	g.vertices[v] = struct{}{}
}

// AddEdge adds an edge between two vertices to the graph.
func (g *Graph) AddEdge(v1, v2 int) {
	if _, ok := g.edges[v1]; !ok {
		g.edges[v1] = map[int]struct{}{}
	}

	g.edges[v1][v2] = struct{}{}

	if _, ok := g.edges[v2]; !ok {
		g.edges[v2] = map[int]struct{}{}
	}

	g.edges[v2][v1] = struct{}{}
}

// HasVertex checks whether the graph contains a vertex.
func (g *Graph) HasVertex(v int) bool {
	_, ok := g.vertices[v]
	return ok
}

// HasEdge checks whether the graph contains an edge between two vertices.
func (g *Graph) HasEdge(v1, v2 int) bool {
	if _, ok := g.edges[v1]; !ok {
		return false
	}
	_, ok := g.edges[v1][v2]
	return ok
}

// Vertices returns all vertices of the graph.
func (g *Graph) Vertices() []int {
	vertices := make([]int, 0, len(g.vertices))
	for v := range g.vertices {
		vertices = append(vertices, v)
	}
	return vertices
}

// Edges returns all edges of the graph.
func (g *Graph) Edges() [][2]int {
	edges := make([][2]int, 0, len(g.edges))
	for v1, v2s := range g.edges {
		for v2 := range v2s {
			edges = append(edges, [2]int{v1, v2})
		}
	}
	return edges
}
