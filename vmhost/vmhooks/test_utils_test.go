package vmhooks

import (
	mock2 "github.com/multiversx/mx-chain-vm-common-go/mock"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
)

type mockeryStruct struct {
	hooks       *VMHooksImpl
	host        *mockery.MockVMHost
	runtime     *mockery.MockRuntimeContext
	metering    *mockery.MockMeteringContext
	output      *mockery.MockOutputContext
	storage     *mockery.MockStorageContext
	blockchain  *mockery.MockBlockchainContext
	managedType *mockery.MockManagedTypesContext
	async       *mockery.MockAsyncContext
	instance    *mockery.MockInstance
}

func createTestVMHooks() (*VMHooksImpl, *mockery.MockVMHost, *mockery.MockRuntimeContext, *mockery.MockMeteringContext, *mockery.MockOutputContext, *mockery.MockStorageContext) {
	hooks, host, runtime, metering, output, storage, _, _ := createTestVMHooksFull()
	return hooks, host, runtime, metering, output, storage
}

func createTestVMHooksFull() (*VMHooksImpl, *mockery.MockVMHost, *mockery.MockRuntimeContext, *mockery.MockMeteringContext, *mockery.MockOutputContext, *mockery.MockStorageContext, *mockery.MockBlockchainContext, *mockery.MockManagedTypesContext) {
	vmHooks := createTestVMHooksClear()

	vmHooks.host.On("Runtime").Return(vmHooks.runtime)
	vmHooks.host.On("Metering").Return(vmHooks.metering)
	vmHooks.host.On("Output").Return(vmHooks.output)
	vmHooks.host.On("Storage").Return(vmHooks.storage)
	vmHooks.host.On("Blockchain").Return(vmHooks.blockchain)
	vmHooks.host.On("ManagedTypes").Return(vmHooks.managedType)

	vmHooks.host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	vmHooks.host.On("IsBuiltinFunctionCall", mock.Anything).Return(false)
	vmHooks.runtime.On("FailExecution", mock.Anything).Return()
	vmHooks.runtime.On("IsUnsafeMode").Return(false)

	baseInstanceSetup(vmHooks.runtime, vmHooks.instance)

	baseMeteringSetup(vmHooks.metering)

	vmHooks.blockchain.On("LastRandomSeed").Return([]byte("rand"))
	vmHooks.blockchain.On("RoundTime").Return(uint64(6000))
	vmHooks.blockchain.On("LastRound").Return(uint64(6000))
	vmHooks.blockchain.On("EpochStartBlockRound").Return(uint64(6000))
	vmHooks.blockchain.On("EpochStartBlockTimeStampMs").Return(uint64(12345000))
	vmHooks.blockchain.On("GetCode", mock.Anything).Return([]byte("code"), nil)

	return vmHooks.hooks, vmHooks.host, vmHooks.runtime, vmHooks.metering, vmHooks.output, vmHooks.storage, vmHooks.blockchain, vmHooks.managedType
}

func createTestVMHooksClear() *mockeryStruct {
	vmHooksMockery := &mockeryStruct{}

	vmHooksMockery.host = &mockery.MockVMHost{}
	vmHooksMockery.runtime = &mockery.MockRuntimeContext{}
	vmHooksMockery.metering = &mockery.MockMeteringContext{}
	vmHooksMockery.output = &mockery.MockOutputContext{}
	vmHooksMockery.storage = &mockery.MockStorageContext{}
	vmHooksMockery.instance = &mockery.MockInstance{}
	vmHooksMockery.blockchain = &mockery.MockBlockchainContext{}
	vmHooksMockery.managedType = &mockery.MockManagedTypesContext{}
	vmHooksMockery.async = &mockery.MockAsyncContext{}

	vmHooksMockery.host.On("Runtime").Return(vmHooksMockery.runtime)
	vmHooksMockery.host.On("Metering").Return(vmHooksMockery.metering)
	vmHooksMockery.host.On("Output").Return(vmHooksMockery.output)
	vmHooksMockery.host.On("Storage").Return(vmHooksMockery.storage)
	vmHooksMockery.host.On("Blockchain").Return(vmHooksMockery.blockchain)
	vmHooksMockery.host.On("ManagedTypes").Return(vmHooksMockery.managedType)
	vmHooksMockery.host.On("Async").Return(vmHooksMockery.async)
	vmHooksMockery.host.On("EnableEpochsHandler").Return(&mock2.EnableEpochsHandlerStub{})

	vmHooksMockery.hooks = NewVMHooksImpl(vmHooksMockery.host)
	return vmHooksMockery
}

func baseInstanceSetup(runtime *mockery.MockRuntimeContext, instance *mockery.MockInstance) {
	runtime.On("GetInstance").Return(instance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return(nil, nil)
	instance.On("MemStore", mock.Anything, mock.Anything).Return(nil)
}

func baseMeteringSetup(metering *mockery.MockMeteringContext) {
	gasSchedule, _ := config.CreateGasConfig(config.MakeGasMapForTests())
	metering.On("GasSchedule").Return(gasSchedule)
	metering.On("StartGasTracing", mock.Anything)
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	metering.On("GasLeft").Return(uint64(100))
	metering.On("BoundGasLimit", mock.Anything).Return(uint64(100))
}

func createTestVMHooksWithSetMetering() (*VMHooksImpl, *mockery.MockVMHost, *mockery.MockRuntimeContext, *mockery.MockMeteringContext, *mockery.MockOutputContext, *mockery.MockStorageContext, *mockery.MockBlockchainContext, *mockery.MockManagedTypesContext, *mockery.MockAsyncContext, *mockery.MockInstance) {
	vmHooks := createTestVMHooksClear()
	baseMeteringSetup(vmHooks.metering)
	baseInstanceSetup(vmHooks.runtime, vmHooks.instance)
	return vmHooks.hooks, vmHooks.host, vmHooks.runtime, vmHooks.metering, vmHooks.output, vmHooks.storage, vmHooks.blockchain, vmHooks.managedType, vmHooks.async, vmHooks.instance
}

func createHooksWithBaseSetup() *mockeryStruct {
	vmHooks := createTestVMHooksClear()
	baseMeteringSetup(vmHooks.metering)
	baseInstanceSetup(vmHooks.runtime, vmHooks.instance)
	vmHooks.host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	vmHooks.host.On("IsBuiltinFunctionCall", mock.Anything).Return(false)
	vmHooks.runtime.On("FailExecution", mock.Anything).Return()
	vmHooks.runtime.On("SetRuntimeBreakpointValue", mock.Anything).Return()

	return vmHooks
}
