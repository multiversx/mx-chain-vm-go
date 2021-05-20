package dex

import (
	"fmt"
)

func (pfe *fuzzDexExecutor) checkPairViews(user string, swapPair SwapPair, stats *eventsStatistics) error {

	outputAmountInA, errAmountInA := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

	outputAmountOutA, errAmountOutA := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

	outputEquivalentOutA, errEquivalentA := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

	outputAmountInB, errAmountInB := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.secondToken, 1000))

	outputAmountOutB, errAmountOutB := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.secondToken, 1000))

	outputEquivalentOutB, errEquivalentB := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.secondToken, 1000))

	if errAmountInA != nil || errAmountInB != nil || errAmountOutA != nil || errAmountOutB != nil ||
		errEquivalentA != nil || errEquivalentB != nil {
		pfe.log("some query returned errors")
		stats.queryPairsMisses += 1
	} else {
		stats.queryPairsHits += 1
	}

	Use(outputAmountInA, outputAmountInB, outputAmountOutA, outputAmountOutB, outputEquivalentOutA, outputEquivalentOutB)

	return nil
}
