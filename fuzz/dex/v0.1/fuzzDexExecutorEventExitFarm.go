package dex

import (
	"errors"
	"fmt"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) exitFarm(amountMax int, statistics *eventsStatistics, rand *rand.Rand) error {
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

	unstakeAmount := int64(amountMax)
	if int64(amountMax) > amount {
		unstakeAmount = amount
	} else {
		unstakeAmount = int64(amountMax)
	}
	farm := pfe.farmers[nonce].farm
	pfe.farmers[nonce] = FarmerInfo{
		value:   amount - unstakeAmount,
		user:    user,
		farm: 	 farm,
	}

	wegldBefore, err := pfe.getTokens(user, pfe.wegldTokenId)
	if err != nil {
		return nil
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
	if output == nil {
		return errors.New("NULL Output")
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		statistics.exitFarmHits += 1

		wegldAfter, err := pfe.getTokens(user, pfe.wegldTokenId)
		if err != nil {
			return nil
		}

		if wegldAfter.Cmp(wegldBefore) == 1 {
			statistics.exitFarmWithRewards += 1
		} else if wegldAfter.Cmp(wegldBefore) == -1 {
			return errors.New("LOST wegld while unstake")
		}

	} else {
		statistics.exitFarmMisses += 1
		pfe.log("exitFarm")
		pfe.log("could not exitFarm because %s", output.ReturnMessage)

		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
