package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) addLiquidity(user string, swapPair SwapPair, amountA int, amountB int,
	amountAmin int, amountBmin int, statistics *eventsStatistics) error {

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
		"txId": "accept-esdt-payment",
		"tx": {
			"from": "%s",
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
		swapPair.address,
		swapPair.firstToken,
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
			"from": "%s",
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
		swapPair.address,
		swapPair.secondToken,
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
			"from": "%s",
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
		swapPair.address,
		amountA,
		amountB,
		amountAmin,
		amountBmin,
	))
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		statistics.addLiquidityHits += 1

		rawEquivalentAfter, errAfter := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))
		if errAfter != nil {
			return errAfter
		}

		if errEquivalent == nil && !(len(rawEquivalent) == 1 && len(rawEquivalent[0]) == 0) {
			statistics.addLiquidityPriceChecks += 1
			if !equalMatrix(rawEquivalentAfter, rawEquivalent) {
				return errors.New("PRICE CHANGED after add liquidity")
			}
		}

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

		if tokenABefore.Cmp(tokenAAfter) < 1 ||
			tokenBBefore.Cmp(tokenBAfter) < 1 ||
			tokenLpBefore.Cmp(tokenLpAfter) > -1 {
			return errors.New("FAILED add liquidity balances on success")
		}
	} else {
		statistics.addLiquidityMisses += 1

		pfe.log("add liquidity %s -> %s", swapPair.firstToken, swapPair.secondToken)
		pfe.log("could not add because %s", output.ReturnMessage)

		expectedErrors := map[string]bool{
			"Insufficient second token computed amount": true,
			"Optimal amount greater than desired amount": true,
			"Insufficient first token computed amount": true,
			"First tokens needs to be greater than minimum liquidity": true,
			"Insufficient liquidity minted": true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
