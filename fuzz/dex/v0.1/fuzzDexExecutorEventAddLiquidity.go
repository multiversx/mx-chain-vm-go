package dex

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"

	vmi "github.com/multiversx/mx-chain-vm-common-go"
)

func (pfe *fuzzDexExecutor) addLiquidity(r *rand.Rand, statistics *eventsStatistics) error {
	user := pfe.userAddress(r.Intn(pfe.numUsers) + 1)
	swapPair := pfe.swaps[r.Intn(len(pfe.swaps))]

	seed := r.Intn(pfe.addLiquidityMaxValue) + 1
	amountA := seed
	amountB := seed
	amountAmin := seed / 100
	amountBmin := seed / 100

	rawEquivalent, err := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))
	if err != nil {
		return err
	}

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

	_, err = pfe.executeTxStep(fmt.Sprintf(`
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

	_, err = pfe.executeTxStep(fmt.Sprintf(`
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

	output, err := pfe.executeTxStep(fmt.Sprintf(`
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

		rawEquivalentAfter, err := pfe.querySingleResultStringAddr(pfe.ownerAddress, swapPair.address,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", swapPair.firstToken, 1000))
		if err != nil {
			return err
		}

		before := big.NewInt(0).SetBytes(rawEquivalent[0])
		if big.NewInt(0).Cmp(before) != 0 {
			statistics.addLiquidityPriceChecks += 1
			after := big.NewInt(0).SetBytes(rawEquivalentAfter[0])
			difference := big.NewInt(0).Abs(big.NewInt(0).Sub(before, after))

			if big.NewInt(1).Cmp(difference) == -1 {
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
			"Insufficient second token computed amount":               true,
			"Optimal amount greater than desired amount":              true,
			"Insufficient first token computed amount":                true,
			"First tokens needs to be greater than minimum liquidity": true,
			"Insufficient liquidity minted":                           true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
