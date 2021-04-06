package dex

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

func (pfe *fuzzDexExecutor) unstake(amountMax int, statistics *eventsStatistics, rand *rand.Rand) error {
	stakersLen := len(pfe.stakers)
	if stakersLen == 0 {
		return nil
	}

	nonce := rand.Intn(stakersLen)
	user := pfe.stakers[nonce].user
	amount := pfe.stakers[nonce].value
	if pfe.stakers[nonce].value == 0 {
		return nil
	}

	unstakeAmount := int64(amountMax)
	if int64(amountMax) > amount {
		unstakeAmount = amount
	} else {
		unstakeAmount = int64(amountMax)
	}
	pfe.stakers[nonce] = StakeInfo{
		value: amount - unstakeAmount,
		user: user,
	}

	wegldBefore, err := pfe.getTokens([]byte(user), pfe.wegldTokenId)
	if err != nil {
		return nil
	}


	reward, err := pfe.querySingleResult(pfe.ownerAddress, pfe.stakingAddress,
		"calculateRewardsForGivenPosition", fmt.Sprintf(`"%d", "%d"`, nonce, unstakeAmount))
	if err != nil {
		return err
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "stake",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "unstake",
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
		pfe.stakingAddress,
		"STAKING-abcdef",
		unstakeAmount,
		nonce,
	))
	if output == nil {
		return errors.New("NULL Output")
	}

	success := output.ReturnMessage == ""
	if success {
		statistics.unstakeHits += 1

		wegldAfter, err := pfe.getTokens([]byte(user), pfe.wegldTokenId)
		if err != nil {
			return nil
		}

		if wegldAfter.Cmp(wegldBefore) == 1 {
			statistics.unstakeWithRewards += 1
		} else if wegldAfter.Cmp(wegldBefore) == -1 {
			return errors.New("LOST wegld while unstake")
		}

		if wegldAfter.Cmp(big.NewInt(0).Add(wegldBefore, big.NewInt(0).SetBytes(reward[0]))) != 0 {
			return errors.New("BAD reward received")
		}

	} else {
		statistics.unstakeMisses += 1
		pfe.log("unstake")
		pfe.log("could not unstake because %s", output.ReturnMessage)

		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}