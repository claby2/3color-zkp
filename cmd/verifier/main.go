// Description: Verifier for the 3-color ZKP protocol.
package main

import (
	"flag"
	"net"
	"os"

	"github.com/claby2/3color-zkp/party"
)

func main() {
	port := flag.String("port", "", "port to listen on")
	repetitions := flag.Int("repetitions", 1, "number of repetitions")
	flag.Parse()

	if *port == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}

	graph := party.ProverToVerifierGraph{}
	err = party.Receive(&graph, conn)
	if err != nil {
		panic(err)
	}

	err = party.Send(party.VerifierToProverRepetitions{Repetitions: *repetitions}, conn)
	if err != nil {
		panic(err)
	}

	failed := 0
	for range *repetitions {
		hashes := party.ProverToVerifierHashes{}
		err = party.Receive(&hashes, conn)
		if err != nil {
			panic(err)
		}

		verifier := party.NewVerifier(graph.Graph, hashes.Hashes)

		a, b := verifier.RandomEdge()
		err = party.Send(party.VerifierToProverEdge{A: a, B: b}, conn)
		if err != nil {
			panic(err)
		}

		commitments := party.ProverToVerifierCommitment{}
		err = party.Receive(&commitments, conn)
		if err != nil {
			panic(err)
		}

		valid := verifier.Verify(a, b, commitments.A, commitments.B)

		if !valid {
			failed++
		}
	}

	if failed == 0 {
		print("✅ ")
		println("All", *repetitions, "repetitions passed.")
	} else {
		print("❌ ")
		println("Summary:", failed, "failed out of", *repetitions)
	}
}
