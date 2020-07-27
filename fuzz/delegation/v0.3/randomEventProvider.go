package delegation

import (
	"math/rand"
	"time"
)

type randomEventProvider struct {
	randProvider  *rand.Rand
	currentRand   float32
	cumulatedProb float32
}

func newRandomEventProvider() *randomEventProvider {
	re := &randomEventProvider{
		randProvider: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	re.reset()
	return re
}

func (re *randomEventProvider) reset() {
	re.currentRand = re.randProvider.Float32()
	re.cumulatedProb = 0
}

func (re *randomEventProvider) withProbability(p float32) bool {
	re.cumulatedProb += p
	if re.cumulatedProb > 1 {
		panic("probabilities exceed 1")
	}
	return re.currentRand < re.cumulatedProb
}
