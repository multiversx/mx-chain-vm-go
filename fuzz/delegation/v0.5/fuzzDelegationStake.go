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

func (pfe *fuzzDelegationExecutor) unBond(delegatorIndex int) error {
	deferredPayment, err := pfe.getUserStakeOfType(delegatorIndex, UserDeferredPayment)
	if err != nil {
		return err
	}
	withdrawOnly, err := pfe.getUserStakeOfType(delegatorIndex, UserWithdrawOnly)
	if err != nil {
		return err
	}
	stakeWithdrawn := big.NewInt(0).Add(deferredPayment, withdrawOnly)

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

		if stakeWithdrawn.Cmp(big.NewInt(0)) > 0 {
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
				stakeWithdrawn,
			))
			if err != nil {
				return err
			}
		}

		pfe.totalStakeWithdrawn.Add(pfe.totalStakeWithdrawn, stakeWithdrawn)
	} else {
		pfe.log("unBond, delegator: %d, fail, %s", delegatorIndex, output.ReturnMessage)
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) getUserStakeOfType(delegatorIndex int, fundType string) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "%s",
			"arguments": ["''%s"],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.delegatorAddress(delegatorIndex)),
		string(pfe.delegationContractAddress),
		fundType,
		string(pfe.delegatorAddress(delegatorIndex)),
	))
	if err != nil {
		return nil, err
	}

	if len(output.ReturnData) != 0 {
		result := big.NewInt(0).SetBytes(output.ReturnData[0])
		return result, nil
	}

	return big.NewInt(0), nil
}

func (pfe *fuzzDelegationExecutor) printTotalStakeByType() {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "getTotalStakeByType",
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
		pfe.log("getTotalStakeByType error")
		return
	}

	if len(output.ReturnData) == 5 {
		pfe.log("total funds in contract: "+
			"WithdrawOnly: %d "+
			"Waiting: %d "+
			"Active: %d "+
			"UnStaked %d "+
			"DeferredPayment %d",
			big.NewInt(0).SetBytes(output.ReturnData[0]),
			big.NewInt(0).SetBytes(output.ReturnData[1]),
			big.NewInt(0).SetBytes(output.ReturnData[2]),
			big.NewInt(0).SetBytes(output.ReturnData[3]),
			big.NewInt(0).SetBytes(output.ReturnData[4]),
		)
	}
}
