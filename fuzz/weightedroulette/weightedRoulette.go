package weightedroulette

import (
	"math/rand"
)

// Outcome holds a possible outcome of roulette random choice execution.
type Outcome struct {
	Weight int
	Event  func()
}

// RandomChoice executes one of the outcomes randomly, according to the weights.
func RandomChoice(r *rand.Rand, outcomes ...Outcome) {
	weightSum := outputWeightSum(outcomes...)
	randomNum := r.Intn(weightSum)
	cumulative := 0
	for _, outcome := range outcomes {
		cumulative += outcome.Weight
		if randomNum < cumulative {
			outcome.Event()
			return
		}
	}
}

func outputWeightSum(outcomes ...Outcome) int {
	sum := 0
	for _, outcome := range outcomes {
		sum += outcome.Weight
	}
	return sum
}
