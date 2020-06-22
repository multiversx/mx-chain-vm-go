package delegation

import (
	"errors"
	"fmt"
	"math/big"
)

func (pfe *fuzzDelegationExecutor) addRewards(amount *big.Int) error {
	// keep track of added rewards
	pfe.totalRewards.Add(pfe.totalRewards, amount)

	// simulate system rewards to delegation contract
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "validatorReward",
		"txId": "%d",
		"tx": {
			"to": "''%s",
			"value": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegationContractAddress),
		amount,
	))
	pfe.log("reward: %d", amount)
	return err
}

func (pfe *fuzzDelegationExecutor) getClaimableRewards(delegIndex int) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "getClaimableRewards",
			"arguments": [
				"''%s"
			],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "",
			"logs": "*",
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegIndex)),
		string(pfe.delegationContractAddress),
		string(pfe.delegatorAddress(delegIndex)),
	))
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(output.ReturnData[0]), nil
}

func (pfe *fuzzDelegationExecutor) claimRewards(delegIndex int) error {
	claimableRewardsBefore, err := pfe.getClaimableRewards(delegIndex)
	if err != nil {
		return err
	}
	if claimableRewardsBefore.Sign() == 0 {
		pfe.log("no rewards, delegator: %d", delegIndex)
		return nil
	}
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "claimRewards",
			"arguments": [],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": "*",
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegIndex)),
		string(pfe.delegationContractAddress),
	))
	if err != nil {
		return err
	}
	claimableRewardsAfter, err := pfe.getClaimableRewards(delegIndex)
	if err != nil {
		return err
	}
	if claimableRewardsAfter.Sign() != 0 {
		return errors.New("getClaimableRewards should always yield 0 after claimRewards")
	}
	pfe.log("claim rewards, delegator: %d, amount: %d", delegIndex, claimableRewardsBefore)
	return nil
}

func (pfe *fuzzDelegationExecutor) computeAllRewards() error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "computeAllRewards",
			"arguments": [],
			"gasLimit": "1,000,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": "*",
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
	))
	pfe.log("computeAllRewards")
	return err
}
