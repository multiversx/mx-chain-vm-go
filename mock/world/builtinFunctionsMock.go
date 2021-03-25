package worldmock

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// BuiltInFunctionESDTTransfer is the key for the elrond standard digital token transfer built-in function
const BuiltInFunctionESDTTransfer = "ESDTTransfer"

func getBuiltinFunctionNames() vmcommon.FunctionNames {
	builtinFunctionNames := make(vmcommon.FunctionNames)
	var empty struct{}
	builtinFunctionNames[BuiltInFunctionESDTTransfer] = empty
	return builtinFunctionNames
}

func (b *MockWorld) processBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if input.Function == BuiltInFunctionESDTTransfer {
		output, err := b.runESDTTransferCall(input)
		return output, err
	}

	return nil, fmt.Errorf("function %s is not provided by the mock as a builtin function", input.Function)
}

// StartTransferESDT updates ESDT balance deltas, without changing the actual balances.
// The deltas need to be committed afterwards.
func (b *MockWorld) StartTransferESDT(from, to []byte, tokenName string, amount *big.Int) (bool, error) {
	sender := b.AcctMap.GetAccount(from)
	senderESDT := sender.ESDTData[tokenName]
	if senderESDT == nil {
		return false, nil
	}

	balancePlusDelta := big.NewInt(0).Add(senderESDT.Balance, senderESDT.BalanceDelta)
	if amount.Cmp(balancePlusDelta) > 0 {
		return false, nil
	}

	senderESDT.BalanceDelta = senderESDT.BalanceDelta.Sub(senderESDT.BalanceDelta, amount)

	recipient := b.AcctMap.GetAccount(to)
	if recipient == nil {
		return true, fmt.Errorf("tx recipient (address: %s) does not exist", hex.EncodeToString(to))
	}
	recipientESDT := recipient.ESDTData[tokenName]
	if recipientESDT == nil {
		recipientESDT = &ESDTData{
			Balance:      big.NewInt(0),
			BalanceDelta: big.NewInt(0),
			Frozen:       false,
		}
		recipient.ESDTData[tokenName] = recipientESDT
	}
	recipientESDT.BalanceDelta = recipientESDT.BalanceDelta.Add(recipientESDT.BalanceDelta, amount)

	return true, nil
}

func (b *MockWorld) runESDTTransferCall(vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if len(vmInput.Arguments) < 2 {
		return nil, errors.New("ESDTTransfer expects at least 2 arguments")
	}
	tokenName := string(vmInput.Arguments[0])
	amount := big.NewInt(0).SetBytes(vmInput.Arguments[1])

	enoughFunds, err := b.StartTransferESDT(vmInput.CallerAddr, vmInput.RecipientAddr, tokenName, amount)
	if err != nil {
		return nil, err
	}

	if !enoughFunds {
		return nil, ErrInsufficientFunds
	}

	vmOutput := &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    vmInput.GasProvided,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}

	isSCCallAfter := len(vmInput.Arguments) > 2 && b.IsSmartContract(vmInput.RecipientAddr)

	if isSCCallAfter {
		endpointName := string(vmInput.Arguments[2])
		var callArgs [][]byte
		if len(vmInput.Arguments) > 3 {
			callArgs = vmInput.Arguments[3:]
		}
		addOutputTransferToVMOutput(
			endpointName,
			callArgs,
			vmInput.RecipientAddr,
			vmInput.GasLocked,
			vmOutput)
	}

	return vmOutput, nil
}

func addOutputTransferToVMOutput(
	function string,
	arguments [][]byte,
	recipient []byte,
	gasLocked uint64,
	vmOutput *vmcommon.VMOutput,
) {
	esdtTransferTxData := function
	for _, arg := range arguments {
		esdtTransferTxData += "@" + hex.EncodeToString(arg)
	}
	outTransfer := vmcommon.OutputTransfer{
		Value:     big.NewInt(0),
		GasLimit:  vmOutput.GasRemaining,
		GasLocked: gasLocked,
		Data:      []byte(esdtTransferTxData),
		CallType:  vmcommon.AsynchronousCall,
	}
	vmOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
	vmOutput.OutputAccounts[string(recipient)] = &vmcommon.OutputAccount{
		Address:         recipient,
		OutputTransfers: []vmcommon.OutputTransfer{outTransfer},
	}
}
