package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) removeLiquidity(user string, swapPair SwapPair, amount int, amountAmin int,
	amountBmin int, statistics *eventsStatistics) error {

	rawEquivalent, errEquivalent := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

	tokenABefore, err := pfe.getTokens(user, swapPair.firstToken)
	if err != nil {
		return err
	}
	tokenBBefore, err := pfe.getTokens(user, swapPair.secondToken)
	if err != nil {
		return err
	}
	tokenLpBefore, err := pfe.getTokens(user, swapPair.lpToken)
	if err != nil {
		return err
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "remove-liq",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "removeLiquidity",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
				"%d",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		swapPair.address,
		swapPair.lpToken,
		amount,
		amountAmin,
		amountBmin,
	))

	if err == nil && output.ReturnCode == vmi.Ok {
		statistics.removeLiquidityHits += 1

		tokenAAfter, err := pfe.getTokens(user, swapPair.firstToken)
		if err != nil {
			return err
		}
		tokenBAfter, err := pfe.getTokens(user, swapPair.secondToken)
		if err != nil {
			return err
		}
		tokenLpAfter, err := pfe.getTokens(user, swapPair.lpToken)
		if err != nil {
			return err
		}

		rawOutput, erro := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))

		if tokenABefore.Cmp(tokenAAfter) > -1 ||
			tokenBBefore.Cmp(tokenBAfter) > -1 ||
			tokenLpBefore.Cmp(tokenLpAfter) < 1 {
			return errors.New("FAILED remove liquidity balances on success")
		}
		if errEquivalent == nil && erro == nil && !(len(rawEquivalent) == 1 && len(rawEquivalent[0]) == 0) {
			statistics.removeLiquidityPriceChecks += 1
			if !equalMatrix(rawEquivalent, rawOutput) {
				return errors.New("PRICE CHANGED after success remove")
			}
		}
	} else {
		statistics.removeLiquidityMisses += 1

		if output == nil {
			return errors.New("output is nil")
		}

		pfe.log("remove liquidity %s -> %s", swapPair.firstToken, swapPair.secondToken)
		pfe.log("could not remove because %s", output.ReturnMessage)

		expectedErrors := map[string]bool{
			"Not enough LP token supply": true,
			"Insufficient liquidity burned": true,
			"Not enough reserve": true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
