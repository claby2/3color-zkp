package main

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/claby2/3color-zkp/graph"
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
	g := graph.New()

	g.AddVertex(0)
	g.AddVertex(1)
	g.AddVertex(2)

	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)

	coloring := map[int]string{
		0: "red",
		1: "green",
		2: "blue",
	}

	prover, err := newProver(g, coloring)
	if err != nil {
		panic(err)
	}

	verifier := newVerifier(g, prover.hashes())

	a, b := verifier.randomEdge()

	oc1, oc2 := prover.openCommitments(a, b)
	if verifier.verify(a, b, oc1, oc2) {
		println("✅ Edge is valid.")
	} else {
		println("❌ Edge is invalid.")
	}
}
