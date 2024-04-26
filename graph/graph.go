package graph

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
