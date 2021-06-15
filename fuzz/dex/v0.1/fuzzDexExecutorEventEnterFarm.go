package dex

import (
	"errors"
	"fmt"

	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) enterFarm(user string, farm Farm, amount int, statistics *eventsStatistics) error {

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "stake",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "enterFarm",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		farm.address,
		farm.farmingToken,
		amount,
	))

	if err == nil && output.ReturnCode == vmi.Ok {
		statistics.enterFarmHits += 1

		pfe.currentFarmTokenNonce[farm.address] += 1
		nonce := pfe.currentFarmTokenNonce[farm.address]
		bigint, errGet := pfe.getTokensWithNonce(user, farm.farmToken, nonce)
		if errGet != nil {
			return errGet
		}

		rps, err := pfe.querySingleResult(user, farm.address, "getRewardPerShare", "")
		if err != nil {
			return err
		}

		pfe.farmers[nonce] = FarmerInfo{
			user:    user,
			value:   bigint.Int64(),
			farm: 	 farm,
			rps:     string(rps[0]),
		}
	} else {
		statistics.enterFarmMisses += 1

		if output == nil {
			return errors.New("output is nil")
		}

		pfe.log("stake %s", farm.farmingToken)
		pfe.log("could enter farm because %s", output.ReturnMessage)

		return errors.New(output.ReturnMessage)
	}

	return nil
}
