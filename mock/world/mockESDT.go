package worldmock

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
)

func (bf *BuiltinFunctionsWrapper) MakeTokenKey(tokenName []byte) []byte {
	keyPrefix := []byte(core.ElrondProtectedKeyPrefix + core.ESDTKeyIdentifier)

	tokenKey := append(keyPrefix, tokenName...)
	return tokenKey
}

func (bf *BuiltinFunctionsWrapper) PerformDirectESDTTransfer(
	sender []byte,
	receiver []byte,
	token []byte,
	value *big.Int,
	callType vmcommon.CallType,
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

	esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, token, value.Bytes())
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

func (bf *BuiltinFunctionsWrapper) GetTokenBalance(address []byte, tokenKey []byte) (*big.Int, error) {
	tokenData, err := bf.GetTokenData(address, tokenKey)
	if err != nil {
		return nil, err
	}

	return tokenData.Value, nil
}

func (bf *BuiltinFunctionsWrapper) SetTokenBalance(address []byte, tokenKey []byte, balance *big.Int) error {
	tokenData, err := bf.GetTokenData(address, tokenKey)
	if err != nil {
		return err
	}

	if balance.Sign() < 0 {
		return data.ErrNegativeValue
	}

	tokenData.Value = balance
	return bf.SetTokenData(address, tokenKey, tokenData)
}

func (bf *BuiltinFunctionsWrapper) GetTokenData(address []byte, tokenKey []byte) (*esdt.ESDigitalToken, error) {
	account := bf.World.AcctMap.GetAccount(address)
	esdtData := &esdt.ESDigitalToken{Value: big.NewInt(0), Type: uint32(core.Fungible)}

	marshaledData, err := account.DataTrieTracker().RetrieveValue(tokenKey)
	if err != nil || len(marshaledData) == 0 {
		return esdtData, nil
	}

	err = bf.Marshalizer.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, err
	}

	return esdtData, nil
}

func (bf *BuiltinFunctionsWrapper) SetTokenData(address []byte, tokenKey []byte, tokenData *esdt.ESDigitalToken) error {
	account := bf.World.AcctMap.GetAccount(address)

	marshaledData, err := bf.Marshalizer.Marshal(tokenData)
	if err != nil {
		return err
	}

	return account.DataTrieTracker().SaveKeyValue(tokenKey, marshaledData)
}
