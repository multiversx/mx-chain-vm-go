package host

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var defaultVmType = []byte{0xF, 0xF}
var ErrCodeNotFound = errors.New("code not found")
var firstAddress = []byte("firstSC.........................")
var secondAddress = []byte("secondSC........................")

func GetSCCode(fileName string) []byte {
	code, _ := ioutil.ReadFile(filepath.Clean(fileName))

	return code
}

func GetTestSCCode(scName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/" + scName + ".wasm"
	return GetSCCode(pathToSC)
}

func DefaultTestArwenForDeployment(t *testing.T, ownerNonce uint64, newAddress []byte) *vmHost {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetNonceCalled = func(address []byte) (uint64, error) {
		return 24, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	return host
}

func DefaultTestArwenForCall(t *testing.T, code []byte) (*vmHost, *mock.BlockchainHookStub) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetCodeCalled = func(address []byte) ([]byte, error) {
		return code, nil
	}
	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

func DefaultTestArwenForTwoSCs(t *testing.T, firstCode []byte, secondCode []byte) (*vmHost, *mock.BlockchainHookStub) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetCodeCalled = func(scAddress []byte) ([]byte, error) {
		if bytes.Equal(scAddress, firstAddress) {
			return firstCode, nil
		}
		if bytes.Equal(scAddress, secondAddress) {
			return secondCode, nil
		}
		return nil, ErrCodeNotFound
	}
	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

func DefaultTestArwen(t *testing.T, blockchain vmcommon.BlockchainHook, crypto vmcommon.CryptoHook) (*vmHost, error) {
	host, err := NewArwenVM(blockchain, crypto, defaultVmType, uint64(1000), config.MakeGasMap(1))
	require.Nil(t, err)
	require.NotNil(t, host)
	return host, err
}

func DefaultTestContractCreateInput() *vmcommon.ContractCreateInput {
	return &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments: [][]byte{
				[]byte("argument 1"),
				[]byte("argument 2"),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		ContractCode: []byte("contract"),
	}
}

func DefaultTestContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  []byte("caller"),
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: []byte("smartcontract"),
		Function:      "function",
	}
}

func MakeVMOutput() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		ReturnData:      make([][]byte, 0),
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
	}
}

func AddFinishData(vmOutput *vmcommon.VMOutput, data []byte) {
	vmOutput.ReturnData = append(vmOutput.ReturnData, data)
}

func AddNewOutputAccount(vmOutput *vmcommon.VMOutput, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	account := &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(balanceDelta),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		Code:           nil,
		Data:           data,
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

func SetStorageUpdate(account *vmcommon.OutputAccount, key []byte, data []byte) {
	keyString := string(key)
	update, exists := account.StorageUpdates[keyString]
	if !exists {
		update = &vmcommon.StorageUpdate{}
		account.StorageUpdates[keyString] = update
	}
	update.Offset = key
	update.Data = data
}

func SetStorageUpdateStrings(account *vmcommon.OutputAccount, key string, data string) {
	SetStorageUpdate(account, []byte(key), []byte(data))
}
