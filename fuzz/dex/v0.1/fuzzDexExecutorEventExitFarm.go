package dex

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (pfe *fuzzDexExecutor) exitFarm(amountMax int, statistics *eventsStatistics, rand *rand.Rand) error {
	stakersLen := len(pfe.farmers)
	if stakersLen == 0 {
		return nil
	}

	nonce := rand.Intn(stakersLen)
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
	lpToken := pfe.farmers[nonce].lpToken
	pfe.farmers[nonce] = FarmerInfo{
		value:   amount - unstakeAmount,
		user:    user,
		lpToken: lpToken,
	}

	wegldBefore, err := pfe.getTokens(user, pfe.wegldTokenId)
	if err != nil {
		return nil
	}

	reward, err := pfe.querySingleResult(pfe.ownerAddress, pfe.wegldFarmingAddress,
		"calculateRewardsForGivenPosition", fmt.Sprintf(`"%d", "%d"`, nonce, unstakeAmount))
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
		pfe.wegldFarmingAddress,
		"FARM-abcdef",
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

		if wegldAfter.Cmp(big.NewInt(0).Add(wegldBefore, big.NewInt(0).SetBytes(reward[0]))) != 0 {
			return errors.New("BAD reward received")
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
