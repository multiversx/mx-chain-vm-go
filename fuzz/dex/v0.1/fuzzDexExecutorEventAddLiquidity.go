package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) addLiquidity(user string, tokenA string, tokenB string, amountA int,
	amountB int , amountAmin int, amountBmin int, statistics *eventsStatistics) error {

	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return err
	}

	if tokenA == tokenB {
		return nil
	}

	rawEquivalent, errEquivalent := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	err, lpTokenStr, _ := pfe.getLpTokenIdentifier(pairHexStr)
	if err != nil {
		return err
	}

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
		"txId": "accept-esdt-payment",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "acceptEsdtPayment",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "0",
			"message": "",
			"gas": "*",
			"refund": "*"
		}
	}`,
		user,
		pairHexStr,
		tokenA,
		amountA,
	))
	if err != nil {
		return err
	}

	output, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "accept-esdt-payment",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "acceptEsdtPayment",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "0",
			"message": "",
			"gas": "*",
			"refund": "*"
		}
	}`,
		user,
		pairHexStr,
		tokenB,
		amountB,
	))
	if err != nil {
		return err
	}

	output, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "add-liquidity",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "addLiquidity",
			"arguments": [
				"%d",
				"%d",
				"%d",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pairHexStr,
		amountA,
		amountB,
		amountAmin,
		amountBmin,
	))
	if output == nil {
		return errors.New("NULL output")
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		// Add liquidity is good
		statistics.addLiquidityHits += 1

		// Get New price
		rawEquivalentAfter, errAfter := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))
		if errAfter != nil {
			return errAfter
		}

		// New and old prices should be the same
		if errEquivalent == nil {
			statistics.addLiquidityPriceChecks += 1
			if  !equalMatrix(rawEquivalentAfter, rawEquivalent) {
				return errors.New("PRICE CHANGED after add liquidity")
			}
		}
	} else {
		statistics.addLiquidityMisses += 1
		pfe.log("add liquidity %s -> %s", tokenA, tokenB)
		pfe.log("could not add because %s", output.ReturnMessage)

		//In case we get these errors but values are !=0, its an error
		if (output.ReturnMessage == "PAIR: INSSUFICIENT TOKEN A FUNDS SENT" ||
			output.ReturnMessage == "PAIR: INSSUFICIENT TOKEN B FUNDS SENT" ||
			output.ReturnMessage == "PAIR: NO AVAILABLE TOKEN A FUNDS" ||
			output.ReturnMessage == "PAIR: NO AVAILABLE TOKEN B FUNDS") &&
			(amountA > 0 && amountB > 0) {
			return errors.New(output.ReturnMessage)
		}

		if output.ReturnMessage == "Pair: FIRST TOKENS NEEDS TO BE GRATER THAN MINIMUM LIQUIDITY: 1000 * 1000e-18" &&
			amountA > 1000 && amountB > 1000 {
			return errors.New(output.ReturnMessage)
		}

		//No way we should receive this
		if output.ReturnMessage == "K invariant failed" {
			return errors.New(output.ReturnMessage)
		}

		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}

		// Other errors are fine
	}

	output, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "reclaim-temporary-funds",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "reclaimTemporaryFunds",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "0",
			"message": "",
			"gas": "*",
			"refund": "*"
		}
	}`,
		user,
		pairHexStr,
	))
	if err != nil {
		return err
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

	if success {
		if tokenABefore.Cmp(tokenAAfter) < 1 ||
			tokenBBefore.Cmp(tokenBAfter) < 1 ||
			tokenLpBefore.Cmp(tokenLpAfter) > -1 {
			return errors.New("FAILED add liquidity balances on success")
		}
	} else {
		if tokenABefore.Cmp(tokenAAfter) != 0 ||
			tokenBBefore.Cmp(tokenBAfter) != 0 ||
			tokenLpBefore.Cmp(tokenLpAfter) != 0 {
			return errors.New("FAILED add liquidity balances on fail")
		}
	}

	return nil
}
