package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) removeLiquidity(user string, tokenA string, tokenB string, amount int, amountAmin int,
	amountBmin int, statistics *eventsStatistics) error {

	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return err
	}

	if tokenA == tokenB {
		return nil
	}

	err, lpTokenStr, _ := pfe.getLpTokenIdentifier(pairHexStr)
	if err != nil {
		return err
	}

	rawEquivalent, errEquivalent := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	tokenABefore, err := pfe.getTokens([]byte(user), tokenA)
	if err != nil {
		return nil
	}
	tokenBBefore, err := pfe.getTokens([]byte(user), tokenB)
	if err != nil {
		return nil
	}
	tokenLpBefore, err := pfe.getTokens([]byte(user), lpTokenStr)
	if err != nil {
		return err
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "remove-liq",
		"tx": {
			"from": "''%s",
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
		pairHexStr,
		lpTokenStr,
		amount,
		amountAmin,
		amountBmin,
	))
	if output == nil {
		return errors.New("NULL Output")
	}

	tokenAAfter, err := pfe.getTokens([]byte(user), tokenA)
	if err != nil {
		return nil
	}
	tokenBAfter, err := pfe.getTokens([]byte(user), tokenB)
	if err != nil {
		return nil
	}
	tokenLpAfter, err := pfe.getTokens([]byte(user), lpTokenStr)
	if err != nil {
		return err
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		statistics.removeLiquidityHits += 1

		rawOutput, erro := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

		if tokenABefore.Cmp(tokenAAfter) > -1 ||
			tokenBBefore.Cmp(tokenBAfter) > -1 ||
			tokenLpBefore.Cmp(tokenLpAfter) < 1 {
			return errors.New("FAILED remove liquidity balances on success")
		}
		if errEquivalent == nil && erro == nil {
			statistics.removeLiquidityPriceChecks += 1
			if !equalMatrix(rawEquivalent, rawOutput) {
				return errors.New("PRICE CHANGED after success remove")
			}
		}
	} else {
		pfe.log("remove liquidity %s -> %s", tokenA, tokenB)
		pfe.log("could not remove because %s", output.ReturnMessage)
		statistics.removeLiquidityMisses += 1

		if tokenABefore.Cmp(tokenAAfter) != 0 ||
			tokenBBefore.Cmp(tokenBAfter) != 0 ||
			tokenLpBefore.Cmp(tokenLpAfter) != 0 {
			return errors.New("FAILED remove liquidity balances on success")
		}
		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
