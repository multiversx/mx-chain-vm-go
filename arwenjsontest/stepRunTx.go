package arwenjsontest

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func (ae *ArwenTestExecutor) executeTx(tx *ij.Transaction) (*vmi.VMOutput, error) {
	beforeErr := ae.world.UpdateWorldStateBefore(
		tx.From.Value,
		tx.GasLimit.Value,
		tx.GasPrice.Value)
	if beforeErr != nil {
		return nil, beforeErr
	}

	var output *vmi.VMOutput

	sender := ae.world.AcctMap.GetAccount(tx.From.Value)
	if sender.Balance.Cmp(tx.Value.Value) < 0 {
		// out of funds is handled by the protocol, so it needs to be mocked here
		output = outOfFundsResult()
	} else {
		switch tx.Type {
		case ij.ScDeploy:
			var err error
			output, err = ae.scCreate(tx)
			if err != nil {
				return nil, err
			}
		case ij.ScCall:
			var err error
			output, err = ae.scCall(tx)
			if err != nil {
				return nil, err
			}
		case ij.Transfer:
			return nil, errors.New("simple transfers not yet implemented")
		default:
			return nil, errors.New("unknown transaction type")
		}

	}

	if output.ReturnCode == vmi.Ok {
		err := ae.updateStateAfterTx(tx, output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

func outOfFundsResult() *vmi.VMOutput {
	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.OutOfFunds,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}
}

func (ae *ArwenTestExecutor) scCreate(tx *ij.Transaction) (*vmi.VMOutput, error) {
	input := &vmi.ContractCreateInput{
		ContractCode: tx.Code.Value,
		VMInput: vmi.VMInput{
			CallerAddr:  tx.From.Value,
			Arguments:   ij.JSONBytesValues(tx.Arguments),
			CallValue:   tx.Value.Value,
			GasPrice:    tx.GasPrice.Value,
			GasProvided: tx.GasLimit.Value,
		},
	}

	return ae.vm.RunSmartContractCreate(input)
}

func (ae *ArwenTestExecutor) scCall(tx *ij.Transaction) (*vmi.VMOutput, error) {
	recipient := ae.world.AcctMap.GetAccount(tx.To.Value)
	if recipient == nil {
		return nil, fmt.Errorf("Tx recipient (address: %s) does not exist", hex.EncodeToString(tx.To.Value))
	}
	if len(recipient.Code) == 0 {
		return nil, fmt.Errorf("Tx recipient (address: %s) is not a smart contract", hex.EncodeToString(tx.To.Value))
	}
	input := &vmi.ContractCallInput{
		RecipientAddr: tx.To.Value,
		Function:      tx.Function,
		VMInput: vmi.VMInput{
			CallerAddr:  tx.From.Value,
			Arguments:   ij.JSONBytesValues(tx.Arguments),
			CallValue:   tx.Value.Value,
			GasPrice:    tx.GasPrice.Value,
			GasProvided: tx.GasLimit.Value,
		},
	}

	return ae.vm.RunSmartContractCall(input)
}

func (ae *ArwenTestExecutor) updateStateAfterTx(
	tx *ij.Transaction,
	output *vmi.VMOutput) error {

	// subtract call value from sender (this is not reflected in the delta)
	_ = ae.world.UpdateBalanceWithDelta(tx.From.Value, big.NewInt(0).Neg(tx.Value.Value))

	accountsSlice := make([]*vmi.OutputAccount, len(output.OutputAccounts))
	i := 0
	for _, account := range output.OutputAccounts {
		accountsSlice[i] = account
		i++
	}

	// update accounts based on deltas
	updErr := ae.world.UpdateAccounts(accountsSlice, output.DeletedAccounts)
	if updErr != nil {
		return updErr
	}

	// sum of all balance deltas should equal call value (unless we got an error)
	sumOfBalanceDeltas := big.NewInt(0)
	for _, oa := range output.OutputAccounts {
		sumOfBalanceDeltas = sumOfBalanceDeltas.Add(sumOfBalanceDeltas, oa.BalanceDelta)
	}
	if sumOfBalanceDeltas.Cmp(tx.Value.Value) != 0 {
		return fmt.Errorf("sum of balance deltas should equal call value. Sum of balance deltas: %d (0x%x). Call value: %d (0x%x)",
			sumOfBalanceDeltas, sumOfBalanceDeltas, tx.Value.Value, tx.Value.Value)
	}

	return nil
}
