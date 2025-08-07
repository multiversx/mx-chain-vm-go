package vmhooks

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_ManagedSCAddress(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetContextAddress").Return([]byte("sc-address"))
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedSCAddress(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("sc-address"))
}

func TestVMHooksImpl_ManagedOwnerAddress(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	blockchain.On("GetOwnerAddress").Return([]byte("owner-address"), nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedOwnerAddress(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("owner-address"))
}

func TestVMHooksImpl_ManagedCaller(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
		},
	})
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedCaller(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("caller"))
}

func TestVMHooksImpl_ManagedGetOriginalCallerAddr(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			OriginalCallerAddr: []byte("original-caller"),
		},
	})
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetOriginalCallerAddr(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("original-caller"))
}

func TestVMHooksImpl_ManagedGetRelayerAddr(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			RelayerAddr: []byte("relayer"),
		},
	})
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetRelayerAddr(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("relayer"))
}

func TestVMHooksImpl_ManagedSignalError(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("error"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("SignalUserError", "error").Return()

	hooks.ManagedSignalError(0)
	runtime.AssertCalled(t, "SignalUserError", "error")
}

func TestVMHooksImpl_ManagedWriteLog(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("topic")}, uint64(1), nil)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("GetContextAddress").Return([]byte("address"))
	output.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return()

	hooks.ManagedWriteLog(0, 0)
}

func TestVMHooksImpl_ManagedGetOriginalTxHash(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetOriginalTxHash").Return([]byte("tx-hash"))
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetOriginalTxHash(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("tx-hash"))
}

func TestVMHooksImpl_ManagedGetStateRootHash(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	blockchain.On("GetStateRootHash").Return([]byte("state-root-hash"))
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetStateRootHash(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("state-root-hash"))
}

func TestVMHooksImpl_ManagedGetBlockRandomSeed(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	blockchain.On("CurrentRandomSeed").Return([]byte("random-seed"))
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetBlockRandomSeed(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("random-seed"))
}

func TestVMHooksImpl_ManagedGetPrevBlockRandomSeed(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	blockchain.On("LastRandomSeed").Return([]byte("random-seed"))
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetPrevBlockRandomSeed(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("random-seed"))
}

func TestVMHooksImpl_ManagedGetReturnData(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	output.On("ReturnData").Return([][]byte{[]byte("data")})
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetReturnData(0, 0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("data"))
}

func TestVMHooksImpl_ManagedGetMultiESDTCallValue(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			ESDTTransfers: []*vmcommon.ESDTTransfer{},
		},
	})
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)

	hooks.ManagedGetMultiESDTCallValue(0)
}

func TestVMHooksImpl_ManagedGetAllTransfersCallValue(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			ESDTTransfers: []*vmcommon.ESDTTransfer{},
			CallValue:     big.NewInt(0),
		},
	})
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)

	hooks.ManagedGetAllTransfersCallValue(0)
}

func TestVMHooksImpl_ManagedGetBackTransfers(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBackTransfers").Return(nil, big.NewInt(0))
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("GetBigIntOrCreate", mock.Anything).Return(big.NewInt(0))

	hooks.ManagedGetBackTransfers(0, 0)
}

func TestVMHooksImpl_ManagedGetESDTBalance(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(&esdt.ESDigitalToken{Value: big.NewInt(100)}, nil)
	managedType.On("GetBigIntOrCreate", mock.Anything).Return(big.NewInt(0))

	hooks.ManagedGetESDTBalance(0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedGetESDTTokenData(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(&esdt.ESDigitalToken{
		Value:         big.NewInt(100),
		TokenMetaData: &esdt.MetaData{},
	}, nil)
	managedType.On("GetBigIntOrCreate", mock.Anything).Return(big.NewInt(0))
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("WriteManagedVecOfManagedBuffers", mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedGetESDTTokenData(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedGetESDTTokenType(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(&esdt.ESDigitalToken{Type: 1}, nil)
	managedType.On("GetBigIntOrCreate", mock.Anything).Return(big.NewInt(0))

	hooks.ManagedGetESDTTokenType(0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedAsyncCall(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedAsyncCall(0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedCreateAsyncCall(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	async.On("RegisterAsyncCall", mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedCreateAsyncCall(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedGetCallbackClosure(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	async.On("GetCallbackClosure").Return([]byte("closure"), nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedGetCallbackClosure(0)
	managedType.AssertCalled(t, "SetBytes", int32(0), []byte("closure"))
}

func TestVMHooksImpl_ManagedUpgradeFromSourceContract(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	runtime.On("SetRuntimeBreakpointValue", mock.Anything).Return()
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	blockchain.On("GetCode", mock.Anything).Return([]byte("code"), nil)
	async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedUpgradeFromSourceContract(0, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedUpgradeContract(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	runtime.On("SetRuntimeBreakpointValue", mock.Anything).Return()
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedUpgradeContract(0, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedDeleteContract(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	runtime.On("SetRuntimeBreakpointValue", mock.Anything).Return()
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedDeleteContract(0, 0, 0)
}

func TestVMHooksImpl_ManagedDeployFromSourceContract(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.VMInput{})
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	blockchain.On("GetCode", mock.Anything).Return([]byte("code"), nil)
	host.On("CreateNewContract", mock.Anything, mock.Anything).Return([]byte("new-address"), nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedDeployFromSourceContract(0, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedCreateContract(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.VMInput{})
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	host.On("CreateNewContract", mock.Anything, mock.Anything).Return([]byte("new-address"), nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.ManagedCreateContract(0, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedExecuteReadOnly(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.VMInput{})
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	host.On("ExecuteOnDestContext", mock.Anything).Return(&vmcommon.VMOutput{}, true, nil)
	runtime.On("ReadOnly").Return(false)
	runtime.On("SetReadOnly", mock.Anything).Return()
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	async.On("SetAsyncArgumentsForCall", mock.Anything).Return()
	async.On("CompleteChildConditional", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	managedType.On("WriteManagedVecOfManagedBuffers", mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedExecuteReadOnly(0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedExecuteOnSameContext(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.VMInput{})
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	host.On("ExecuteOnSameContext", mock.Anything).Return(nil)
	managedType.On("WriteManagedVecOfManagedBuffers", mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedExecuteOnSameContext(0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedExecuteOnDestContext(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.VMInput{})
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	host.On("ExecuteOnDestContext", mock.Anything).Return(&vmcommon.VMOutput{}, true, nil)
	host.On("CompleteLogEntriesWithCallType", mock.Anything, mock.Anything).Return()
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	async.On("SetAsyncArgumentsForCall", mock.Anything).Return()
	async.On("CompleteChildConditional", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	managedType.On("WriteManagedVecOfManagedBuffers", mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedExecuteOnDestContext(0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedExecuteOnDestContextWithErrorReturn(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.VMInput{})
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	host.On("ExecuteOnDestContext", mock.Anything).Return(&vmcommon.VMOutput{}, true, nil)
	host.On("CompleteLogEntriesWithCallType", mock.Anything, mock.Anything).Return()
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	async.On("SetAsyncArgumentsForCall", mock.Anything).Return()
	async.On("CompleteChildConditional", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	managedType.On("WriteManagedVecOfManagedBuffers", mock.Anything, mock.Anything).Return(nil)

	hooks.ManagedExecuteOnDestContextWithErrorReturn(0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_ManagedMultiTransferESDTNFTExecute(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.VMInput{})
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	host.On("ExecuteOnDestContext", mock.Anything).Return(&vmcommon.VMOutput{}, true, nil)
	host.On("CompleteLogEntriesWithCallType", mock.Anything, mock.Anything).Return()
	async := &mockery.MockAsyncContext{}
	host.On("Async").Return(async)
	async.On("SetAsyncArgumentsForCall", mock.Anything).Return()
	async.On("CompleteChildConditional", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	managedType.On("WriteManagedVecOfManagedBuffers", mock.Anything, mock.Anything).Return(nil)
	output := &mockery.MockOutputContext{}
	host.On("Output").Return(output)
	output.On("TransferESDT", mock.Anything, mock.Anything).Return(uint64(0), nil)
	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	blockchain.On("GetSnapshot").Return(0)

	hooks.ManagedMultiTransferESDTNFTExecute(0, 0, 0, 0, 0)
}
