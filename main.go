package main

import (
	"crypto/sha256"
	"encoding/binary"
	"os"
	"strconv"

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
	if len(os.Args) < 2 || len(os.Args) > 3 {
		println("Usage: ./3color-zkp <graph> [repetitions]")
		os.Exit(1)
	}

	graphFile := os.Args[1]
	graphData, err := os.ReadFile(graphFile)
	if err != nil {
		panic(err)
	}

	repetitions := 1
	if len(os.Args) == 3 {
		repetitions, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
	}

	g, coloring, err := graph.Parse(string(graphData))
	if err != nil {
		panic(err)
	}

	prover, err := newProver(g, coloring)
	if err != nil {
		panic(err)
	}

	verifier := newVerifier(g, prover.hashes())

	failed := 0
	for i := 0; i < repetitions; i++ {
		a, b := verifier.randomEdge()

		oc1, oc2 := prover.openCommitments(a, b)
		if !verifier.verify(a, b, oc1, oc2) {
			failed++
		}
	}
	if failed == 0 {
		print("✅ ")
	} else {
		print("❌ ")
	}
	println("Summary:", failed, "failed out of", repetitions)
}
