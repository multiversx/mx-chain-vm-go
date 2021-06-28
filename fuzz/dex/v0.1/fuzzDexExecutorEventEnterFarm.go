package dex

import (
	"errors"
	"fmt"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func (pfe *fuzzDexExecutor) enterFarm(r *rand.Rand, statistics *eventsStatistics) error {
	user := pfe.userAddress(r.Intn(pfe.numUsers) + 1)
	amount := r.Intn(pfe.enterFarmMaxValue) + 1
	farm := pfe.farms[r.Intn(len(pfe.farms))]

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
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
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
			user:  user,
			value: bigint.Int64(),
			farm:  farm,
			rps:   string(rps[0]),
		}
	} else {
		statistics.enterFarmMisses += 1

		pfe.log("stake %s", farm.farmingToken)
		pfe.log("could enter farm because %s", output.ReturnMessage)

		return errors.New(output.ReturnMessage)
	}

	return nil
}
