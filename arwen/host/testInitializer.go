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
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var defaultVMType = []byte{0xF, 0xF}
var errCodeNotFound = errors.New("code not found")
var firstAddress = []byte("firstSC.........................")
var secondAddress = []byte("secondSC........................")

// GetSCCode retrieves the bytecode of a WASM module from a file
func GetSCCode(fileName string) []byte {
	code, _ := ioutil.ReadFile(filepath.Clean(fileName))

	return code
}

// GetTestSCCode retrieves the bytecode of a WASM testing module
func GetTestSCCode(scName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/" + scName + ".wasm"
	return GetSCCode(pathToSC)
}

// DefaultTestArwenForDeployment creates an Arwen vmHost configured for testing deployments
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

// DefaultTestArwenForCall creates an Arwen vmHost configured for testing SC calls
func DefaultTestArwenForCall(t *testing.T, code []byte) (*vmHost, *mock.BlockchainHookStub) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetCodeCalled = func(address []byte) ([]byte, error) {
		return code, nil
	}
	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

// DefaultTestArwenForTwoSCs creates an Arwen vmHost configured for testing calls between 2 SmartContracts
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
		return nil, errCodeNotFound
	}
	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

// DefaultTestArwen creates an Arwen vmHost configured with the provided BlockchainHook and CryptoHook
func DefaultTestArwen(t *testing.T, blockchain vmcommon.BlockchainHook, crypto vmcommon.CryptoHook) (*vmHost, error) {
	host, err := NewArwenVM(blockchain, crypto, defaultVMType, uint64(1000), config.MakeGasMap(1))
	require.Nil(t, err)
	require.NotNil(t, host)
	return host, err
}

// DefaultTestContractCreateInput creates a vmcommon.ContractCreateInput struct with default values
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

// DefaultTestContractCallInput creates a vmcommon.ContractCallInput struct with default values
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

// MakeVMOutput creates a vmcommon.VMOutput struct with default values
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

// AddFinishData appends the provided []byte to the ReturnData of the given vmOutput
func AddFinishData(vmOutput *vmcommon.VMOutput, data []byte) {
	vmOutput.ReturnData = append(vmOutput.ReturnData, data)
}

// AddNewOutputAccount creates a new vmcommon.OutputAccount from the provided arguments and adds it to OutputAccounts of the provided vmOutput
func AddNewOutputAccount(vmOutput *vmcommon.VMOutput, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	account := &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(balanceDelta),
		Balance:        nil,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		Code:           nil,
		Data:           data,
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

// SetStorageUpdate sets a storage update to the provided vmcommon.OutputAccount
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

// SetStorageUpdateStrings sets a storage update to the provided vmcommon.OutputAccount, from string arguments
func SetStorageUpdateStrings(account *vmcommon.OutputAccount, key string, data string) {
	SetStorageUpdate(account, []byte(key), []byte(data))
}
