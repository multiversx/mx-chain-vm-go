package dex

import (
	"errors"
	"fmt"
	"math/rand"

	vmi "github.com/multiversx/mx-chain-vm-common-go"
)

func (pfe *fuzzDexExecutor) exitFarm(r *rand.Rand, statistics *eventsStatistics) error {
	amountMax := r.Intn(pfe.exitFarmMaxValue) + 1

	stakersLen := len(pfe.farmers)
	if stakersLen == 0 {
		return nil
	}

	nonce := rand.Intn(stakersLen) + 1
	user := pfe.farmers[nonce].user
	amount := pfe.farmers[nonce].value
	if pfe.farmers[nonce].value == 0 {
		return nil
	}

	farm := pfe.farmers[nonce].farm
	var unstakeAmount int64
	if int64(amountMax) > amount {
		unstakeAmount = amount
		delete(pfe.farmers, nonce)
	} else {
		unstakeAmount = int64(amountMax)
		pfe.farmers[nonce] = FarmerInfo{
			value: amount - unstakeAmount,
			user:  user,
			farm:  farm,
		}
	}

	mexBefore, err := pfe.getTokens(user, pfe.mexTokenId)
	if err != nil {
		return err
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "stake",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "exitFarm",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d",
				"nonce": "%d"
			},
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		farm.address,
		farm.farmToken,
		unstakeAmount,
		nonce,
	))
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		statistics.exitFarmHits += 1

		mexAfter, err := pfe.getTokens(user, pfe.mexTokenId)
		if err != nil {
			return err
		}

		if mexAfter.Cmp(mexBefore) == 1 {
			statistics.exitFarmWithRewards += 1
		} else if mexAfter.Cmp(mexBefore) == -1 {
			return errors.New("LOST mex while exiting farm")
		}

	} else {
		statistics.exitFarmMisses += 1

		pfe.log("exitFarm")
		pfe.log("could not exitFarm because %s", output.ReturnMessage)

		expectedErrors := map[string]bool{
			"Exit too early for lock rewards option": true,
			"Farming token amount is zero":           true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
