package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) stake(user string, tokenA string, tokenB string, amount int, statistics *eventsStatistics) error {
	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return err
	}

	if tokenA == tokenB {
		return nil
	}

	err, lpTokenStr, lpTokenHexStr := pfe.getLpTokenIdentifier(pairHexStr)
	if err != nil {
		return err
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "stake",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "stake",
			"esdt": {
				"tokenIdentifier": "%s",
				"value": "%d"
			},
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pfe.wegldStakingAddress,
		lpTokenHexStr,
		amount,
	))
	if output == nil {
		return errors.New("NULL Output")
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		statistics.stakeHits += 1

		pfe.currentStakeTokenNonce += 1
		nonce := pfe.currentStakeTokenNonce
		bigint, errGet := pfe.getTokensWithNonce([]byte(user), "STAKING-abcdef", nonce)
		if errGet != nil {
			return errGet
		}
		pfe.stakers[nonce] = StakeInfo{
			user: user,
			value: bigint.Int64(),
			lpToken: lpTokenStr,
		}
	} else {
		statistics.stakeMisses += 1
		pfe.log("stake %s -> %s", tokenA, tokenB)
		pfe.log("could stake add because %s", output.ReturnMessage)

		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}