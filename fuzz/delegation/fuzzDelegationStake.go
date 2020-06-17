package delegation

import (
	"fmt"
	"math/big"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func (pfe *fuzzDelegationExecutor) stake(delegIndex int, amount *big.Int) error {
	// keep track of stake added
	pfe.totalStakeAdded.Add(pfe.totalStakeAdded, amount)

	// get the stake from the big sack
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "transfer",
		"txId": "%d",
		"tx": {
			"from": "''endless_sack_of_erd___________s1",
			"to": "''%s",
			"value": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegIndex)),
		amount,
	))
	if err != nil {
		return err
	}

	// actual staking
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "%d",
			"function": "stake",
			"arguments": [],
			"gasLimit": "100,000,000",
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
		amount,
	))
	pfe.log("stake, delegator: %d, amount: %d", delegIndex, amount)
	return err
}

func (pfe *fuzzDelegationExecutor) withdrawInactiveStake(delegIndex int, amount *big.Int) error {
	// keep track of stake withdrawn
	pfe.totalStakeWithdrawn.Add(pfe.totalStakeWithdrawn, amount)

	// actual withdraw
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "withdrawInactiveStake",
			"arguments": [
				"%d"
			],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegIndex)),
		string(pfe.delegationContractAddress),
		amount,
	))
	if err != nil {
		return err
	}
	if output.ReturnCode == vmi.Ok {
		pfe.log("unstake, delegator: %d, amount: %d", delegIndex, amount)
	} else {
		pfe.log("unstake, delegator: %d, amount: %d, fail, %s", delegIndex, amount, output.ReturnMessage)
	}

	// move withdrawn stake to a special account
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "transfer",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegIndex)),
		pfe.withdrawTargetAddress,
		amount,
	))
	return err
}

func (pfe *fuzzDelegationExecutor) getUserInactiveStake(delegIndex int) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "getUserInactiveStake",
			"arguments": [
				"''%s"
			],
			"gasLimit": "100,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "",
			"logs": [],
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

	result := big.NewInt(0).SetBytes(output.ReturnData[0])
	pfe.log("getUserInactiveStake -> %d", result)
	return result, nil
}

func (pfe *fuzzDelegationExecutor) withdrawAllInactiveStake(delegIndex int) error {
	inactiveStake, err := pfe.getUserInactiveStake(delegIndex)
	if err != nil {
		return err
	}
	return pfe.withdrawInactiveStake(delegIndex, inactiveStake)
}
