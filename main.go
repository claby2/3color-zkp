package main

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/dominikbraun/graph"
)

// Lambda is the security parameter of the protocol.
const Lambda = 128

// commit creates a commitment for the given color and randomness.
func commit(color string, r int64) [32]byte {
	rBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(rBytes, uint64(r))
	hash := sha256.Sum256([]byte(color + string(rBytes)))
	return hash
}

func main() {
	g := graph.New(graph.IntHash)
	g.AddVertex(0, graph.VertexAttribute("color", "red"))
	g.AddVertex(1, graph.VertexAttribute("color", "blue"))
	g.AddVertex(2, graph.VertexAttribute("color", "green"))

	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	prover, err := newProver(g)
	if err != nil {
		panic(err)
	}

	verifier := newVerifier(g, prover.hashes())

	edge, err := verifier.randomEdge()
	if err != nil {
		panic(err)
	}

	oc1, oc2 := prover.openCommitments(edge)
	if verifier.verify(edge, oc1, oc2) {
		println("✅ Edge is valid.")
	} else {
		println("❌ Edge is invalid.")
	}
}
