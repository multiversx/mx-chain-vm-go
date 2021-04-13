package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) stake(user string, tokenA string, tokenB string, amount int, statistics *eventsStatistics) error {
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

	lpTokenRaw, err := pfe.querySingleResult(pfe.ownerAddress, pairAddressRaw[0],
		"getLpTokenIdentifier", "")
	if err != nil {
		return err
	}

	lpTokenHex := "0x"
	for i := 0; i < len(lpTokenRaw[0]); i++ {
		toAppend := fmt.Sprintf("%02x", lpTokenRaw[0][i])
		lpTokenHex += toAppend
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
		lpTokenHex,
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
			lpToken: string(lpTokenRaw[0]),
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