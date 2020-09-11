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
			"from": "''%s",
			"to": "''%s",
			"value": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.faucetAddress),
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

func (pfe *fuzzDelegationExecutor) withdrawInactiveStake(delegatorIndex int, amount *big.Int) error {
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
		string(pfe.delegatorAddress(delegatorIndex)),
		string(pfe.delegationContractAddress),
		amount,
	))
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		// keep track of stake withdrawn
		pfe.totalStakeWithdrawn.Add(pfe.totalStakeWithdrawn, amount)

		pfe.log("withdraw inactive stake, delegator: %d, amount: %d", delegatorIndex, amount)

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
			string(pfe.delegatorAddress(delegatorIndex)),
			string(pfe.withdrawTargetAddress),
			amount,
		))
		if err != nil {
			return err
		}
	} else {
		pfe.log("withdraw inactive stake, delegator: %d, amount: %d, fail, %s", delegatorIndex, amount, output.ReturnMessage)
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) unStake(delegatorIndex int, stake *big.Int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "unStake",
			"arguments": ["%d"],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegatorIndex)),
		string(pfe.delegationContractAddress),
		stake,
	))
	if err != nil {
		return err
	}
	if output.ReturnCode == vmi.Ok {
		pfe.log("unStake, delegator: %d", delegatorIndex)
	} else {
		pfe.log("unStake, delegator: %d, fail, %s", delegatorIndex, output.ReturnMessage)
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) withdrawAllInactiveStake(delegatorIndex int) error {
	inactiveStake, err := pfe.getUserInactiveStake(delegatorIndex)
	if err != nil {
		return err
	}

	return pfe.withdrawInactiveStake(delegatorIndex, inactiveStake)
}

func (pfe *fuzzDelegationExecutor) getUserInactiveStake(delegatorIndex int) (*big.Int, error) {
	return pfe.delegatorQuery("getUserInactiveStake", delegatorIndex)
}

func (pfe *fuzzDelegationExecutor) unBond(delegatorIndex int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "unBond",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegatorIndex)),
		string(pfe.delegationContractAddress),
	))
	if err != nil {
		return err
	}

	if output.ReturnCode == vmi.Ok {
		pfe.log("unBond, delegator: %d", delegatorIndex)
	} else {
		pfe.log("unBond, delegator: %d, fail, %s", delegatorIndex, output.ReturnMessage)
	}

	return nil
}
