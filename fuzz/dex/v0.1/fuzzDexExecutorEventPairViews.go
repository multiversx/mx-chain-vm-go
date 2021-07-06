package dex

import (
	"fmt"
	"math/rand"
)

func (pfe *fuzzDexExecutor) checkPairViews(r *rand.Rand, stats *eventsStatistics) error {
	swapPair := pfe.swaps[r.Intn(len(pfe.swaps))]

	_, errAmountInA := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

	_, errAmountOutA := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

	_, errEquivalentA := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

	_, errAmountInB := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.secondToken, 1000))

	_, errAmountOutB := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.secondToken, 1000))

	_, errEquivalentB := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.secondToken, 1000))

	if errAmountInA != nil || errAmountInB != nil ||
		errAmountOutA != nil || errAmountOutB != nil ||
		errEquivalentA != nil || errEquivalentB != nil {
		stats.queryPairsMisses += 1

		pfe.log("some queries returned errors")
	} else {
		stats.queryPairsHits += 1
	}

	return nil
}
