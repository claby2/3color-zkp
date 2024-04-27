default: prover verifier

prover:
	go build -o prover ./cmd/prover

verifier:
	go build -o verifier ./cmd/verifier

clean:
	go clean
	rm -f prover verifier

.PHONY: default prover verifier clean
