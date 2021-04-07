package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) stake(user string, tokenA string, tokenB string, amount int, statistics *eventsStatistics) error {

	if tokenA == tokenB {
		return nil
	}

	pairAddressRaw, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
		"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
	if err != nil {
		return err
	}

	lpTokenRaw, err := pfe.querySingleResult(pfe.ownerAddress, pairAddressRaw[0],
		"get_lp_token_identifier", "")
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
		pfe.stakingAddress,
		lpTokenHex,
		amount,
	))
	if output == nil {
		return errors.New("NULL Output")
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		statistics.stakeHits += 1

		pfe.currentNftNonce += 1
		nonce := pfe.currentNftNonce
		bigint, errGet := pfe.getTokensWithNonce([]byte(user), "STAKING-abcdef", nonce)
		if errGet != nil {
			return errGet
		}
		pfe.stakers[nonce] = StakeInfo{
			user: user,
			value: bigint.Int64(),
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