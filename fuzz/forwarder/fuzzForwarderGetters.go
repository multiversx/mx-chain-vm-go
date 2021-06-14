package fuzzForwarder

import (
	"fmt"
	"math/rand"
)

func (pfe *fuzzExecutor) nextTxIndex() int {
	pfe.txIndex++
	return pfe.txIndex
}

func (pfe *fuzzExecutor) forwarderAddress(index int) string {
	return fmt.Sprintf("sc:forwarder-%02d", index)
}

func (pfe *fuzzExecutor) fungibleTokenName(index int) string {
	return fmt.Sprintf("FUNG-%02d", index)
}

func (pfe *fuzzExecutor) semiFungibleTokenName(index int) string {
	return fmt.Sprintf("SFT-%02d", index)
}

func (pfe *fuzzExecutor) randomTokenNameAndNonce(r *rand.Rand) (string, int) {
	index := r.Intn(1 + pfe.data.numFungibleTokens + pfe.data.numSemiFungibleTokens)
	if index == 0 {
		return "EGLD", 0
	}
	if index <= pfe.data.numFungibleTokens {
		return pfe.fungibleTokenName(index), 0
	}
	return pfe.semiFungibleTokenName(index - pfe.data.numFungibleTokens), 1
}
