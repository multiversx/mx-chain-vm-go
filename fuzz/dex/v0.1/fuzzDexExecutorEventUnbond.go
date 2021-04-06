package dex

import (
	"errors"
	"fmt"
)

func (pfe *fuzzDexExecutor) unbond(user string, tokenA string, tokenB string, statistics *eventsStatistics) error {

	if tokenA == tokenB {
		return nil
	}

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

	rawLpResponse, err := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"get_lp_token_identifier", "")
	if err != nil {
		return err
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "remove-liq",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "unbond",
			"arguments": [
				"str:%s"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		string(pfe.stakingAddress),
		string(rawLpResponse[0]),
	))
	if output == nil {
		return errors.New("NULL Output")
	}

	success := output.ReturnMessage == ""
	if success {
		statistics.unbondHits += 1
	} else {
		statistics.unbondMisses += 1
		pfe.log("unbond %s -> %s", tokenA, tokenB)
		pfe.log("could not unbond because %s", output.ReturnMessage)

		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}