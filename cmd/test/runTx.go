package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func executeTx(tx *ij.Transaction,
	world *worldhook.BlockchainHookMock,
	vm vmi.VMExecutionHandler) (*vmi.VMOutput, error) {

	beforeErr := world.UpdateWorldStateBefore(tx.From, tx.GasLimit, tx.GasPrice)
	if beforeErr != nil {
		return nil, beforeErr
	}

	arguments := make([][]byte, len(tx.Arguments))
	for i, arg := range tx.Arguments {
		arguments[i] = append(arguments[i], arg...)
	}
	var output *vmi.VMOutput

	sender := world.AcctMap.GetAccount(tx.From)
	if sender.Balance.Cmp(tx.Value) < 0 {
		// out of funds is handled by the protocol, so needs to be mocked here
		output = &vmcommon.VMOutput{
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
	} else if tx.IsCreate {
		// SC create
		input := &vmi.ContractCreateInput{
			ContractCode: tx.Code,
			VMInput: vmi.VMInput{
				CallerAddr:  tx.From,
				Arguments:   arguments,
				CallValue:   tx.Value,
				GasPrice:    tx.GasPrice,
				GasProvided: tx.GasLimit,
			},
		}

		var err error
		output, err = vm.RunSmartContractCreate(input)
		if err != nil {
			return nil, err
		}
	} else {
		// SC call
		recipient := world.AcctMap.GetAccount(tx.To)
		if recipient == nil {
			return nil, fmt.Errorf("Tx recipient (address: %s) does not exist", hex.EncodeToString(tx.To))
		}
		if len(recipient.Code) == 0 {
			return nil, fmt.Errorf("Tx recipient (address: %s) is not a smart contract", hex.EncodeToString(tx.To))
		}
		input := &vmi.ContractCallInput{
			RecipientAddr: tx.To,
			Function:      tx.Function,
			VMInput: vmi.VMInput{
				CallerAddr:  tx.From,
				Arguments:   arguments,
				CallValue:   tx.Value,
				GasPrice:    tx.GasPrice,
				GasProvided: tx.GasLimit,
			},
		}

		var err error
		output, err = vm.RunSmartContractCall(input)
		if err != nil {
			return nil, err
		}
	}

	if output.ReturnCode == vmi.Ok {
		// subtract call value from sender (this is not reflected in the delta)
		_ = world.UpdateBalanceWithDelta(tx.From, big.NewInt(0).Neg(tx.Value))

		accountsSlice := make([]*vmi.OutputAccount, len(output.OutputAccounts))
		i := 0
		for _, account := range output.OutputAccounts {
			accountsSlice[i] = account
			i++
		}

		// update accounts based on deltas
		updErr := world.UpdateAccounts(accountsSlice, output.DeletedAccounts)
		if updErr != nil {
			return nil, updErr
		}

		// sum of all balance deltas should equal call value (unless we got an error)
		sumOfBalanceDeltas := big.NewInt(0)
		for _, oa := range output.OutputAccounts {
			sumOfBalanceDeltas = sumOfBalanceDeltas.Add(sumOfBalanceDeltas, oa.BalanceDelta)
		}
		if sumOfBalanceDeltas.Cmp(tx.Value) != 0 {
			return nil, fmt.Errorf("sum of balance deltas should equal call value. Sum of balance deltas: %d (0x%x). Call value: %d (0x%x)",
				sumOfBalanceDeltas, sumOfBalanceDeltas, tx.Value, tx.Value)
		}
	}

	return output, nil
}
