package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"math/big"
)

func (pfe *fuzzDexExecutor) swapFixedInput(user string, tokenA string, amountA int, tokenB string,
	amountB int, statistics *eventsStatistics) error {
	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return err
	}

	if tokenA == tokenB {
		return nil
	}

	tokenABefore, err := pfe.getTokens([]byte(user), tokenA)
	if err != nil {
		return nil
	}
	tokenBBefore, err := pfe.getTokens([]byte(user), tokenB)
	if err != nil {
		return nil
	}

	amoutOutRaw, amoutOutErr := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, amountA))

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "swap-fixed-input",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "swapTokensFixedInput",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
				"str:%s",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pairHexStr,
		tokenA,
		amountA,
		tokenB,
		amountB,
	))
	if output == nil {
		return errors.New("NULL output")
	}

	tokenAAfter, err := pfe.getTokens([]byte(user), tokenA)
	if err != nil {
		return nil
	}
	tokenBAfter, err := pfe.getTokens([]byte(user), tokenB)
	if err != nil {
		return nil
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		statistics.swapFixedInputHits += 1

		if amoutOutErr == nil {
			//Check if tokens send vs received are correct
			if tokenAAfter.Cmp(big.NewInt(0).Sub(tokenABefore, big.NewInt(int64(amountA)))) != 0 {
				return errors.New("SWAP fixed input wrong amount A")
			}
			if tokenBAfter.Cmp(big.NewInt(0).Add(tokenBBefore, big.NewInt(0).SetBytes(amoutOutRaw[0]))) != 0 {
				return errors.New("SWAP fixed input wrong amount B")
			}
		}
	} else {
		statistics.swapFixedInputMisses += 1
		pfe.log("swapFixedInput %s -> %s", tokenA, tokenB)
		pfe.log("could not swap because %s", output.ReturnMessage)

		if tokenAAfter.Cmp(tokenABefore) != 0 {
			return errors.New("SWAP wrong amount A")
		}
		if tokenBAfter.Cmp(tokenBBefore) != 0 {
			return errors.New("SWAP wrong amount B")
		}
		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
		if output.ReturnMessage == "K invariant failed" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}

func (pfe *fuzzDexExecutor) swapFixedOutput(user string, tokenA string, amountA int, tokenB string,
	amountB int, statistics *eventsStatistics) error {

	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return err
	}

	if tokenA == tokenB {
		return nil
	}

	tokenABefore, err := pfe.getTokens([]byte(user), tokenA)
	if err != nil {
		return nil
	}
	tokenBBefore, err := pfe.getTokens([]byte(user), tokenB)
	if err != nil {
		return nil
	}

	amountInRaw, amountInErr := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", tokenB, amountB))

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "swap-fixed-input",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "swapTokensFixedOutput",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
				"str:%s",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pairHexStr,
		tokenA,
		amountA,
		tokenB,
		amountB,
	))
	if output == nil {
		return errors.New("NULL output")
	}

	tokenAAfter, err := pfe.getTokens([]byte(user), tokenA)
	if err != nil {
		return nil
	}
	tokenBAfter, err := pfe.getTokens([]byte(user), tokenB)
	if err != nil {
		return nil
	}

	success := output.ReturnMessage == ""
	if success {
		statistics.swapFixedOutputHits += 1

		if amountInErr == nil {
			//Check if tokens send vs received are correct
			if tokenAAfter.Cmp(big.NewInt(0).Sub(tokenABefore, big.NewInt(0).SetBytes(amountInRaw[0]))) != 0 {
				return errors.New("SWAP fixed output wrong amount A")
			}
			if tokenBAfter.Cmp(big.NewInt(0).Add(tokenBBefore, big.NewInt(int64(amountB)))) != 0 {
				return errors.New("SWAP fixed output wrong amount B")
			}
		}
	} else {
		statistics.swapFixedOutputMisses += 1
		pfe.log("swapFixedOutput %s -> %s", tokenA, tokenB)
		pfe.log("could not swap because %s", output.ReturnMessage)

		if tokenAAfter.Cmp(tokenABefore) != 0 {
			return errors.New("SWAP wrong amount A")
		}
		if tokenBAfter.Cmp(tokenBBefore) != 0 {
			return errors.New("SWAP wrong amount B")
		}
		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
		if output.ReturnMessage == "K invariant failed" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}