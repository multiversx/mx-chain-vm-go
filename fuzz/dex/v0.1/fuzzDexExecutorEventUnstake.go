package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
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
	lpToken := pfe.stakers[nonce].lpToken
	pfe.stakers[nonce] = StakeInfo{
		value: amount - unstakeAmount,
		user: user,
		lpToken: lpToken,
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

	success := output.ReturnCode == vmi.Ok
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

		pfe.currentUnstakeTokenNonce += 1
		nonce = pfe.currentUnstakeTokenNonce
		bigint, errGet := pfe.getTokensWithNonce([]byte(user), "UNSTAK-abcdef", nonce)
		if errGet != nil {
			return errGet
		}
		pfe.unstakers[nonce] = UnstakeInfo{
			user: user,
			value: bigint.Int64(),
			lpToken: lpToken,
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