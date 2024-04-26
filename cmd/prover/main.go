// Description: Prover for the 3-color ZKP protocol.
package main

import (
	"flag"
	"net"
	"os"

	"github.com/claby2/3color-zkp/graph"
	"github.com/claby2/3color-zkp/party"
)

func main() {
	address := flag.String("address", "", "address of the verifier")
	port := flag.String("port", "", "port of the verifier")
	graphPath := flag.String("graph", "", "path to the graph file")
	flag.Parse()

	if *address == "" || *port == "" || *graphPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	graphData, err := os.ReadFile(*graphPath)
	if err != nil {
		panic(err)
	}

	g, coloring, err := graph.Parse(string(graphData))
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", *address+":"+*port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = party.Send(party.ProverToVerifierGraph{Graph: g}, conn)
	if err != nil {
		panic(err)
	}

	repetitions := party.VerifierToProverRepetitions{}
	err = party.Receive(&repetitions, conn)
	if err != nil {
		panic(err)
	}

	for range repetitions.Repetitions {
		prover, err := party.NewProver(g, coloring)
		if err != nil {
			panic(err)
		}

		err = party.Send(party.ProverToVerifierHashes{Hashes: prover.Hashes()}, conn)
		if err != nil {
			panic(err)
		}

		edge := party.VerifierToProverEdge{}
		err = party.Receive(&edge, conn)
		if err != nil {
			panic(err)
		}

		oc1, oc2 := prover.Open(edge.A, edge.B)

		err = party.Send(party.ProverToVerifierCommitment{A: oc1, B: oc2}, conn)
		if err != nil {
			panic(err)
		}
	}
}
