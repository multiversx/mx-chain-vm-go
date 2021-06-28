package dex

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func (pfe *fuzzDexExecutor) claimRewards(r *rand.Rand, statistics *eventsStatistics) error {
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
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		statistics.claimRewardsHits += 1

		mexAfter, err := pfe.getTokens(user, pfe.mexTokenId)
		if err != nil {
			return err
		}

		rpsAfter, err := pfe.querySingleResult(user, farm.address, "getRewardPerShare", "")
		if err != nil {
			return err
		}

		if mexAfter.Cmp(mexBefore) == 1 {
			statistics.claimRewardsWithRewards += 1

			rpsDifference := big.NewInt(0).Sub(big.NewInt(0).SetBytes(rpsAfter[0]), big.NewInt(0).SetBytes(rpsBefore))
			shouldEarn := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(claimAmount), rpsDifference), big.NewInt(1000000000000))
			earned := big.NewInt(0).Sub(mexAfter, mexBefore)

			if earned.Cmp(shouldEarn) != 0 {
				return errors.New("reward is not good")
			}

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
			user:  user,
			value: bigint.Int64(),
			farm:  farm,
			rps:   string(rpsAfter[0]),
		}
	} else {
		statistics.claimRewardsMisses += 1

		pfe.log("could not claimRewards because %s", output.ReturnMessage)

		return errors.New(output.ReturnMessage)
	}

	return nil
}
