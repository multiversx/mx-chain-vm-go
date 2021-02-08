package worldmock

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// BuiltInFunctionESDTTransfer is the key for the elrond standard digital token transfer built-in function
const BuiltInFunctionESDTTransfer = "ESDTTransfer"

func getBuiltinFunctionNames() vmcommon.FunctionNames {
	builtinFunctionNames := make(vmcommon.FunctionNames)
	var empty struct{}
	builtinFunctionNames[BuiltInFunctionESDTTransfer] = empty
	return builtinFunctionNames
}

func (b *MockWorld) processBuiltInFunction(input *vmcommon.ContractCallInput) (*vmi.VMOutput, error) {
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
		return true, fmt.Errorf("Tx recipient (address: %s) does not exist", hex.EncodeToString(to))
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

func (b *MockWorld) runESDTTransferCall(input *vmcommon.ContractCallInput) (*vmi.VMOutput, error) {
	if len(input.Arguments) != 2 {
		return nil, errors.New("ESDTTransfer expects 2 arguments")
	}
	tokenName := string(input.Arguments[0])
	amount := big.NewInt(0).SetBytes(input.Arguments[1])

	enoughFunds, err := b.StartTransferESDT(input.CallerAddr, input.RecipientAddr, tokenName, amount)
	if err != nil {
		return nil, err
	}

	if !enoughFunds {
		return nil, ErrInsufficientFunds
	}

	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    input.GasProvided - input.GasLocked,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}, nil
}
