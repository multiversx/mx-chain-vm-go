package math

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

type seedRandReader struct {
	rand *rand.Rand
}

// NewSeedRandReader creates and returns a new SeedRandReader
func NewSeedRandReader(seed []byte) *seedRandReader {
	seedHash := sha256.Sum256(seed)
	seedNumber := binary.BigEndian.Uint64(seedHash[:])

	source := rand.NewSource(int64(seedNumber))
	randomizer := rand.New(source)

	return &seedRandReader{
		rand: randomizer,
	}
}

// Read generates len(p) random bytes and writes them into p. It always returns len(p) and a nil error.
func (srr *seedRandReader) Read(p []byte) (n int, err error) {
	return srr.rand.Read(p)
}

// IsInterfaceNil returns true if there is no value under the interface
func (srr *seedRandReader) IsInterfaceNil() bool {
	return srr == nil
}
