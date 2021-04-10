package dex

import (
	"errors"
	"fmt"
)

func (pfe *fuzzDexExecutor) checkPairViews(user string, tokenA string, tokenB string, stats *eventsStatistics) error {
	pairAddressRaw, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
		"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
	if err != nil {
		return err
	}

	pairHexStr := "0x"
	for i := 0; i < len(pairAddressRaw[0]); i++ {
		toAppend := fmt.Sprintf("%02x", pairAddressRaw[0][i])
		pairHexStr += toAppend
	}

	if pairHexStr == "0x0000000000000000000000000000000000000000000000000000000000000000" && tokenA != tokenB {
		return errors.New("NULL pair for different tokens")
	}

	if tokenA == tokenB {
		return nil
	}

	outputAmountInA, errAmountInA := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	outputAmountOutA, errAmountOutA := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	outputEquivalentOutA, errEquivalentA := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	outputAmountInB, errAmountInB := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", tokenB, 1000))

	outputAmountOutB, errAmountOutB := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", tokenB, 1000))

	outputEquivalentOutB, errEquivalentB := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenB, 1000))

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
