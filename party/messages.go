package party

import (
	"encoding/binary"
	"encoding/json"
	"net"

	"github.com/claby2/3color-zkp/graph"
)

type ProverToVerifierGraph struct {
	Graph graph.Graph `json:"graph"`
}

type VerifierToProverRepetitions struct {
	Repetitions int `json:"repetitions"`
}

type ProverToVerifierHashes struct {
	Hashes map[int][32]byte `json:"hashes"`
}

type VerifierToProverEdge struct {
	A int `json:"a"`
	B int `json:"b"`
}

type ProverToVerifierCommitment struct {
	A OpenCommitment `json:"A"`
	B OpenCommitment `json:"B"`
}

func Send(v any, conn net.Conn) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(data)))
	_, err = conn.Write(size)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func Receive(v any, conn net.Conn) error {
	size := make([]byte, 4)
	_, err := conn.Read(size)
	if err != nil {
		return err
	}

	dataSize := binary.BigEndian.Uint32(size)
	data := make([]byte, dataSize)
	_, err = conn.Read(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
