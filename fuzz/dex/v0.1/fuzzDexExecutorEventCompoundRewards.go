package dex

import (
	"errors"
	"fmt"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) compoundRewards(r *rand.Rand, statistics *eventsStatistics) error {
	amountMax := r.Intn(pfe.claimRewardsMaxValue) + 1

	stakersLen := len(pfe.farmers)
	if stakersLen == 0 {
		return nil
	}

	nonce := rand.Intn(stakersLen) + 1
	user := pfe.farmers[nonce].user
	amount := pfe.farmers[nonce].value
	rpsBefore := []byte(pfe.farmers[nonce].rps)
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
		value: amount - claimAmount,
		user:  user,
		farm:  farm,
		rps:   string(rpsBefore),
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "claimRewards",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "compoundRewards",
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
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		statistics.compoundRewardsHits += 1

		rpsAfter, err := pfe.querySingleResult(user, farm.address, "getRewardPerShare", "")
		if err != nil {
			return err
		}

		pfe.currentFarmTokenNonce[farm.address] += 1
		nonce := pfe.currentFarmTokenNonce[farm.address]
		bigint, errGet := pfe.getTokensWithNonce(user, farm.farmToken, nonce)
		if errGet != nil {
			return errGet
		}

		pfe.farmers[nonce] = FarmerInfo{
			user:  user,
			value: bigint.Int64(),
			farm:  farm,
			rps:   string(rpsAfter[0]),
		}
	} else {
		statistics.compoundRewardsMisses += 1

		expectedErrors := map[string]bool{
			"Farming token differ from reward token": true,
			"Farming token amount is zero":           true,
		}

		_, expected := expectedErrors[output.ReturnMessage]
		if !expected {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}
