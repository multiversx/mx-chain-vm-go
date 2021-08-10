package math

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

type SeedRandReader struct {
	rand *rand.Rand
}

func NewSeedRandReader(seed []byte) (*SeedRandReader, error) {
	if len(seed) == 0 {
		return nil, ErrSeedLengthIsZero
	}

	seedHash := sha256.Sum256(seed)
	seedNumber := binary.BigEndian.Uint64(seedHash[:])

	source := rand.NewSource(int64(seedNumber))
	rand := rand.New(source)

	return &SeedRandReader{
		rand: rand,
	}, nil
}

// Read will read upto len(p) bytes. It will rotate the existing byte buffer (seed) until it will fill up the provided
// p buffer
func (srr *SeedRandReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, ErrSeedLengthIsZero
	}

	return srr.rand.Read(p)
}
