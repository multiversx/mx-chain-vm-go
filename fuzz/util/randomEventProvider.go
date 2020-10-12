package fuzzutil

import (
	"math/rand"
	"time"
)

// RandomEventProvider fuzzing utility.
type RandomEventProvider struct {
	randProvider  *rand.Rand
	currentRand   float32
	cumulatedProb float32
}

// NewRandomEventProvider is a RandomEventProvider constructor.
func NewRandomEventProvider() *RandomEventProvider {
	re := &RandomEventProvider{
		randProvider: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	re.Reset()
	return re
}

// Reset clears the RandomEventProvider.
func (re *RandomEventProvider) Reset() {
	re.currentRand = re.randProvider.Float32()
	re.cumulatedProb = 0
}

// WithProbability randomly provides true, according to cumulated probabilities.
func (re *RandomEventProvider) WithProbability(p float32) bool {
	re.cumulatedProb += p
	if re.cumulatedProb > 1 {
		panic("probabilities exceed 1")
	}
	return re.currentRand < re.cumulatedProb
}
