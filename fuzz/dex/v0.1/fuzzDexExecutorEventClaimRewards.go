package dex

import (
	"errors"
	"fmt"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) claimRewards(amountMax int, statistics *eventsStatistics, rand *rand.Rand) error {
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

	claimAmount := int64(amountMax)
	if int64(amountMax) > amount {
		claimAmount = amount
	} else {
		claimAmount = int64(amountMax)
	}

	farm := pfe.farmers[nonce].farm
	pfe.farmers[nonce] = FarmerInfo{
		value:   amount - claimAmount,
		user:    user,
		farm: 	 farm,
	}

	mexBefore, err := pfe.getTokens(user, pfe.mexTokenId)
	if err != nil {
		return nil
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "claimRewards",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "claimRewards",
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
		claimAmount,
		nonce,
	))
	if output == nil {
		return errors.New("NULL Output")
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		statistics.claimRewardsHits += 1

		mexAfter, err := pfe.getTokens(user, pfe.mexTokenId)
		if err != nil {
			return nil
		}

		if mexAfter.Cmp(mexBefore) == 1 {
			statistics.claimRewardsWithRewards += 1
		} else if mexAfter.Cmp(mexBefore) == -1 {
			return errors.New("LOST mex while claimRewards")
		}

		pfe.currentFarmTokenNonce[farm.address] += 1
		nonce := pfe.currentFarmTokenNonce[farm.address]
		bigint, errGet := pfe.getTokensWithNonce(user, farm.farmToken, nonce)
		if errGet != nil {
			return errGet
		}
		pfe.farmers[nonce] = FarmerInfo{
			user:    user,
			value:   bigint.Int64(),
			farm: 	 farm,
		}
	} else {
		statistics.claimRewardsMisses += 1
		pfe.log("could not claimRewards because %s", output.ReturnMessage)

		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
