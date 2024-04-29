package party

import (
	"crypto/sha256"
	"encoding/binary"
)

const Lambda = 32

// commit creates a commitment for the given color and randomness.
func commit(color string, r int64) [32]byte {
	rBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(rBytes, uint64(r))
	hash := sha256.Sum256([]byte(color + string(rBytes)))
	return hash
}

// OpenCommitment represents an open commitment to a color and a random value.
type OpenCommitment struct {
	Color string `json:"color"`
	R     int64  `json:"r"`
}
