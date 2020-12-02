package host

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

var defaultVMType = []byte{0xF, 0xF}
var errAccountNotFound = errors.New("account not found")

var userAddress = []byte("userAccount.....................")
var parentAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fparentSC..............")
var childAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fchildSC...............")

var CustomGasSchedule = config.GasScheduleMap(nil)

// DefaultTestArwenForDeployment creates an Arwen vmHost configured for testing deployments
func DefaultTestArwenForDeployment(t *testing.T, _ uint64, newAddress []byte) *vmHost {
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &mock.AccountMock{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host, _ := DefaultTestArwen(t, stubBlockchainHook)
	return host
}

func DefaultTestArwenForCall(tb testing.TB, code []byte, balance *big.Int) (*vmHost, *mock.BlockchainHookStub) {
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &mock.AccountMock{
				Code:    code,
				Balance: balance,
			}, nil
		}
		return nil, errAccountNotFound
	}

	host, _ := DefaultTestArwen(tb, stubBlockchainHook)
	return host, stubBlockchainHook
}

// DefaultTestArwenForTwoSCs creates an Arwen vmHost configured for testing calls between 2 SmartContracts
func DefaultTestArwenForTwoSCs(t *testing.T, parentCode []byte, childCode []byte, parentSCBalance *big.Int) (*vmHost, *mock.BlockchainHookStub) {
	stubBlockchainHook := &mock.BlockchainHookStub{}

	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &mock.AccountMock{
				Code:    parentCode,
				Balance: parentSCBalance,
			}, nil
		}
		if bytes.Equal(scAddress, childAddress) {
			return &mock.AccountMock{
				Code: childCode,
			}, nil
		}

		return nil, errAccountNotFound
	}

	host, _ := DefaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

func DefaultTestArwen(tb testing.TB, blockchain vmcommon.BlockchainHook) (*vmHost, error) {
	gasSchedule := CustomGasSchedule
	if gasSchedule == nil {
		gasSchedule = config.MakeGasMapForTests()
	}

	host, err := NewArwenVM(blockchain, &arwen.VMHostParameters{
		VMType:                   defaultVMType,
		BlockGasLimit:            uint64(1000),
		GasSchedule:              gasSchedule,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
		UseWarmInstance:          false,
		DynGasLockEnableEpoch:    0,
	})
	require.Nil(tb, err)
	require.NotNil(tb, host)
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
			CallerAddr:  userAddress,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: parentAddress,
		Function:      "function",
	}
}

func LoadGasScheduleConfig(filepath string) (config.GasScheduleMap, error) {
	gasScheduleConfig, err := arwen.LoadTomlFileToMap(filepath)
	if err != nil {
		return nil, err
	}

	flattenedGasSchedule := make(config.GasScheduleMap)
	for libType, costs := range gasScheduleConfig {
		flattenedGasSchedule[libType] = make(map[string]uint64)
		costsMap := costs.(map[string]interface{})
		for operationName, cost := range costsMap {
			flattenedGasSchedule[libType][operationName] = uint64(cost.(int64))
		}
	}

	return flattenedGasSchedule, nil
}
