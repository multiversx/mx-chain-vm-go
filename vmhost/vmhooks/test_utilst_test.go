package vmhooks

import (
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
)

func createTestVMHooks() (*VMHooksImpl, *mockery.MockVMHost, *mockery.MockRuntimeContext, *mockery.MockMeteringContext, *mockery.MockOutputContext, *mockery.MockStorageContext) {
	host := &mockery.MockVMHost{}
	runtime := &mockery.MockRuntimeContext{}
	metering := &mockery.MockMeteringContext{}
	output := &mockery.MockOutputContext{}
	storage := &mockery.MockStorageContext{}
	instance := &mockery.MockInstance{}

	host.On("Runtime").Return(runtime)
	host.On("Metering").Return(metering)
	host.On("Output").Return(output)
	host.On("Storage").Return(storage)
	runtime.On("FailExecution", mock.Anything).Return()
	runtime.On("GetInstance").Return(instance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return(nil, nil)
	instance.On("MemStore", mock.Anything, mock.Anything).Return(nil)

	gasSchedule, _ := config.CreateGasConfig(config.MakeGasMapForTests())
	metering.On("GasSchedule").Return(gasSchedule)
	metering.On("StartGasTracing", mock.Anything)
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	metering.On("GasLeft").Return(uint64(100))

	hooks := NewVMHooksImpl(host)
	return hooks, host, runtime, metering, output, storage
}
