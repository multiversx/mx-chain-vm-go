package worldmock

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-scenario-go/esdtconvert"
	mj "github.com/multiversx/mx-chain-scenario-go/model"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// GetTokenBalance returns the ESDT balance of an account for the given token
// key (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) GetTokenBalance(address []byte, tokenIdentifier []byte, nonce uint64) (*big.Int, error) {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return nil, fmt.Errorf("account does not exist: %s", hex.EncodeToString(address))
	}
	return esdtconvert.GetTokenBalance(tokenIdentifier, nonce, account.Storage)
}

// GetTokenData gets the ESDT information related to a token from the storage of an account
// (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) GetTokenData(address []byte, tokenIdentifier []byte, nonce uint64) (*esdt.ESDigitalToken, error) {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return nil, fmt.Errorf("account does not exist: %s", hex.EncodeToString(address))
	}
	systemAccStorage := make(map[string][]byte)
	systemAcc := bf.World.AcctMap.GetAccount(vmcommon.SystemAccountAddress)
	if systemAcc != nil {
		systemAccStorage = systemAcc.Storage
	}
	return account.GetTokenData(tokenIdentifier, nonce, systemAccStorage)
}

// SetTokenData sets the ESDT information related to a token from the storage of an account
// (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) SetTokenData(address []byte, tokenIdentifier []byte, nonce uint64, tokenData *esdt.ESDigitalToken) error {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return fmt.Errorf("account does not exist: %s", hex.EncodeToString(address))
	}
	return account.SetTokenData(tokenIdentifier, nonce, tokenData)
}

// PerformDirectESDTTransfer calls the real ESDTTransfer function immediately;
// only works for in-shard transfers for now, but it will be expanded to
// cross-shard.
// TODO rewrite to simulate what the SCProcessor does when executing a tx with
// data "ESDTTransfer@token@value@contractfunc@contractargs..."
// TODO this function duplicates code from host.ExecuteESDTTransfer(), must refactor
func (bf *BuiltinFunctionsWrapper) PerformDirectESDTTransfer(
	sender []byte,
	receiver []byte,
	token []byte,
	nonce uint64,
	value *big.Int,
	callType vm.CallType,
	gasLimit uint64,
	gasPrice uint64,
) (uint64, error) {
	esdtTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    callType,
			GasPrice:    gasPrice,
			GasProvided: gasLimit,
			GasLocked:   0,
		},
		RecipientAddr:     receiver,
		Function:          core.BuiltInFunctionESDTTransfer,
		AllowInitFunction: false,
	}

	if nonce > 0 {
		esdtTransferInput.Function = core.BuiltInFunctionESDTNFTTransfer
		esdtTransferInput.RecipientAddr = esdtTransferInput.CallerAddr
		nonceAsBytes := big.NewInt(0).SetUint64(nonce).Bytes()
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, token, nonceAsBytes, value.Bytes(), receiver)
	} else {
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, token, value.Bytes())
	}

	vmOutput, err := bf.ProcessBuiltInFunction(esdtTransferInput)
	if err != nil {
		return 0, err
	}

	if vmOutput.ReturnCode != vmcommon.Ok {
		return 0, fmt.Errorf(
			"ESDTtransfer failed: retcode = %d, msg = %s",
			vmOutput.ReturnCode,
			vmOutput.ReturnMessage)
	}

	return vmOutput.GasRemaining, nil
}

// PerformDirectMultiESDTTransfer -
func (bf *BuiltinFunctionsWrapper) PerformDirectMultiESDTTransfer(
	sender []byte,
	receiver []byte,
	esdtTransfers []*mj.ESDTTxData,
	callType vm.CallType,
	gasLimit uint64,
	gasPrice uint64,
) (uint64, error) {
	nrTransfers := len(esdtTransfers)
	nrTransfersAsBytes := big.NewInt(0).SetUint64(uint64(nrTransfers)).Bytes()

	multiTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    callType,
			GasPrice:    gasPrice,
			GasProvided: gasLimit,
			GasLocked:   0,
		},
		RecipientAddr:     sender,
		Function:          core.BuiltInFunctionMultiESDTNFTTransfer,
		AllowInitFunction: false,
	}
	multiTransferInput.Arguments = append(multiTransferInput.Arguments, receiver, nrTransfersAsBytes)

	for i := 0; i < nrTransfers; i++ {
		token := esdtTransfers[i].TokenIdentifier.Value
		nonceAsBytes := big.NewInt(0).SetUint64(esdtTransfers[i].Nonce.Value).Bytes()
		value := esdtTransfers[i].Value.Value

		multiTransferInput.Arguments = append(multiTransferInput.Arguments, token, nonceAsBytes, value.Bytes())
	}

	vmOutput, err := bf.ProcessBuiltInFunction(multiTransferInput)
	if err != nil {
		return 0, err
	}

	if vmOutput.ReturnCode != vmcommon.Ok {
		return 0, fmt.Errorf(
			"MultiESDTtransfer failed: retcode = %d, msg = %s",
			vmOutput.ReturnCode,
			vmOutput.ReturnMessage)
	}

	return vmOutput.GasRemaining, nil
}
