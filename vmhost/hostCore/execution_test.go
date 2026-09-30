package hostCore

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	hostmock "github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

// TestAddress1 is a test address
var TestAddress1 = []byte("addr1")

// StorageSimpleSC is a simple smart contract that stores a value.
var StorageSimpleSC = []byte(`
(module
  (func (export "store") (param i32)
    (set_local 0 (get_local 0))
  )
  (func (export "load") (result i32)
    (i32.const 42)
  )
)
`)

func TestDoRunSmartContractCreate(t *testing.T) {
	t.Parallel()

	h, err := vmhost.NewVMHost(
		&contextmock.BlockchainHookStub{},
		&vmhost.VMHostParameters{
			ESDTTransferParser:   &parsers.ESDTTransferParser{},
			BuiltInFuncContainer: &builtInFunctions.BuiltInFunctionContainer{},
			EpochNotifier:        &hostmock.EpochNotifierStub{},
			EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
				IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
					return true
				},
			},
			Hasher: &mock.HasherMock{},
			VMType: "mock",
		},
	)
	require.NoError(t, err)

	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: TestAddress1,
			CallValue:  big.NewInt(0),
		},
		ContractCode: StorageSimpleSC,
	}

	vmOutput := h.(*vmHost).doRunSmartContractCreate(input)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
}

func TestDoRunSmartContractCreate_NilInput(t *testing.T) {
	t.Parallel()

	h, err := vmhost.NewVMHost(
		&contextmock.BlockchainHookStub{},
		&vmhost.VMHostParameters{
			ESDTTransferParser:   &parsers.ESDTTransferParser{},
			BuiltInFuncContainer: &builtInFunctions.BuiltInFunctionContainer{},
			EpochNotifier:        &hostmock.EpochNotifierStub{},
			EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
				IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
					return true
				},
			},
			Hasher: &mock.HasherMock{},
			VMType: "mock",
		},
	)
	require.NoError(t, err)

	vmOutput := h.(*vmHost).doRunSmartContractCreate(nil)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ErrFatal, vmOutput.ReturnCode)
}

func TestDoRunSmartContractCall_NilInput(t *testing.T) {
	t.Parallel()

	h, err := vmhost.NewVMHost(
		&contextmock.BlockchainHookStub{},
		&vmhost.VMHostParameters{
			ESDTTransferParser:   &parsers.ESDTTransferParser{},
			BuiltInFuncContainer: &builtInFunctions.BuiltInFunctionContainer{},
			EpochNotifier:        &hostmock.EpochNotifierStub{},
			EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
				IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
					return true
				},
			},
			Hasher: &mock.HasherMock{},
			VMType: "mock",
		},
	)
	require.NoError(t, err)

	vmOutput := h.(*vmHost).doRunSmartContractCall(nil)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ErrFatal, vmOutput.ReturnCode)
}

func TestDoRunSmartContractCall(t *testing.T) {
	t.Parallel()

	h, err := vmhost.NewVMHost(
		&contextmock.BlockchainHookStub{},
		&vmhost.VMHostParameters{
			ESDTTransferParser:   &parsers.ESDTTransferParser{},
			BuiltInFuncContainer: &builtInFunctions.BuiltInFunctionContainer{},
			EpochNotifier:        &hostmock.EpochNotifierStub{},
			EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
				IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
					return true
				},
			},
			Hasher: &mock.HasherMock{},
			VMType: "mock",
		},
	)
	require.NoError(t, err)

	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: TestAddress1,
			CallValue:  big.NewInt(0),
		},
		Function: "load",
	}

	vmOutput := h.(*vmHost).doRunSmartContractCall(input)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
}
