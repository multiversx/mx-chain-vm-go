//nolint:all
package delegation

import (
	"fmt"
	"math/big"

	vmi "github.com/multiversx/mx-chain-vm-common-go"
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

func (pfe *fuzzDelegationExecutor) stakeGenesis(delegIndex int, amount *big.Int) error {
	// keep track of stake added
	pfe.totalStakeAdded.Add(pfe.totalStakeAdded, amount)

	// actual staking
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "stakeGenesis",
			"arguments": [
				"%d"
			],
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
	pfe.log("stakeGenesis, delegator: %d, amount: %d", delegIndex, amount)
	return err
}

func (pfe *fuzzDelegationExecutor) activateGenesis() error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "activateGenesis",
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
	pfe.log("activateGenesis")
	return err
}

func (pfe *fuzzDelegationExecutor) stakeAllAvailable(delegIndex int) error {
	pfe.log("stakeAllAvailable, called by delegator: %d", delegIndex)
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "stakeAllAvailable",
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
	))
	return err
}

func (pfe *fuzzDelegationExecutor) withdrawInactiveStake(delegIndex int, amount *big.Int) error {
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
		// keep track of stake withdrawn
		pfe.totalStakeWithdrawn.Add(pfe.totalStakeWithdrawn, amount)

		pfe.log("withdraw inactive stake, delegator: %d, amount: %d", delegIndex, amount)

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
			string(pfe.withdrawTargetAddress),
			amount,
		))
		if err != nil {
			return err
		}
	} else {
		pfe.log("withdraw inactive stake, delegator: %d, amount: %d, fail, %s", delegIndex, amount, output.ReturnMessage)
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) getUserInactiveStake(delegIndex int) (*big.Int, error) {
	return pfe.delegatorQuery("getUserInactiveStake", delegIndex)
}

func (pfe *fuzzDelegationExecutor) withdrawAllInactiveStake(delegIndex int) error {
	inactiveStake, err := pfe.getUserInactiveStake(delegIndex)
	if err != nil {
		return err
	}
	return pfe.withdrawInactiveStake(delegIndex, inactiveStake)
}

func (pfe *fuzzDelegationExecutor) getUserActiveStake(delegIndex int) (*big.Int, error) {
	return pfe.delegatorQuery("getUserActiveStake", delegIndex)
}

func (pfe *fuzzDelegationExecutor) announceUnStake(delegIndex int, amount *big.Int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "announceUnStake",
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
		pfe.log("announceUnStake, delegator: %d, amount: %d", delegIndex, amount)
	} else {
		pfe.log("announceUnStake, delegator: %d, amount: %d, fail, %s", delegIndex, amount, output.ReturnMessage)
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) announceUnStakeAll(delegIndex int) error {
	userActiveStake, err := pfe.getUserActiveStake(delegIndex)
	if err != nil {
		return err
	}
	return pfe.announceUnStake(delegIndex, userActiveStake)
}

func (pfe *fuzzDelegationExecutor) purchaseStake(sellerIndex, buyerIndex int, amount *big.Int) error {
	// get the value from the big sack
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
		string(pfe.delegatorAddress(buyerIndex)),
		amount,
	))
	if err != nil {
		return err
	}

	// the purchase itself
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "%d",
			"function": "purchaseStake",
			"arguments": [
				"''%s"
			],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(buyerIndex)),
		string(pfe.delegationContractAddress),
		amount,
		string(pfe.delegatorAddress(sellerIndex)),
	))
	if err != nil {
		return err
	}
	if output.ReturnCode == vmi.Ok {
		pfe.log("purchaseStake, seller: %d, buyer: %d, amount: %d", sellerIndex, buyerIndex, amount)

		// forward received sum
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
			string(pfe.delegatorAddress(sellerIndex)),
			string(pfe.stakePurchaseForwardAddress),
			amount,
		))
		if err != nil {
			return err
		}
	} else {
		pfe.log("purchaseStake, seller: %d, buyer: %d, amount: %d, fail, %s", sellerIndex, buyerIndex, amount, output.ReturnMessage)

		// return the value
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
			string(pfe.delegatorAddress(buyerIndex)),
			string(pfe.faucetAddress),
			amount,
		))
		if err != nil {
			return err
		}
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) unStake(delegIndex int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "unStake",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegIndex)),
		string(pfe.delegationContractAddress),
	))
	if err != nil {
		return err
	}
	if output.ReturnCode == vmi.Ok {
		pfe.log("unStake, delegator: %d", delegIndex)
	} else {
		pfe.log("unStake, delegator: %d, fail, %s", delegIndex, output.ReturnMessage)
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) unBondAllAvailable() error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "unBondAllAvailable",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
	))
	if err != nil {
		return err
	}
	if output.ReturnCode == vmi.Ok {
		pfe.log("unBondAllAvailable")
	} else {
		pfe.log("unBondAllAvailable, fail, %s", output.ReturnMessage)
	}

	return nil
}
