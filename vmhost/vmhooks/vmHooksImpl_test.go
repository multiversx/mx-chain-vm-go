package vmhooks

import (
	"errors"
	"github.com/multiversx/mx-chain-core-go/core"
	mock2 "github.com/multiversx/mx-chain-vm-common-go/mock"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewVMHooksImpl(t *testing.T) {
	t.Parallel()
	host := &mockery.MockVMHost{}
	hooks := NewVMHooksImpl(host)
	require.NotNil(t, hooks)
	require.Equal(t, host, hooks.host)
}

func TestVMHooksImpl_MemLoad(t *testing.T) {
	t.Parallel()
	vmHooks := createTestVMHooksClear()
	instance := vmHooks.instance
	runtime := vmHooks.runtime
	hooks := vmHooks.hooks

	runtime.On("GetInstance").Return(instance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return([]byte("data"), nil)

	data, err := hooks.MemLoad(0, 0)
	require.Nil(t, err)
	require.Equal(t, []byte("data"), data)
}

func TestVMHooksImpl_MemLoadMultiple(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, _, _, _ := createTestVMHooks()
	instance := &mockery.MockInstance{}
	runtime.On("GetInstance").Return(instance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return([]byte("data"), nil)

	data, err := hooks.MemLoadMultiple(0, []int32{4, 4})
	require.Nil(t, err)
	require.Len(t, data, 2)
}

func TestVMHooksImpl_MemStore(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, _, _, _ := createTestVMHooks()
	instance := &mockery.MockInstance{}
	runtime.On("GetInstance").Return(instance)
	instance.On("MemStore", mock.Anything, mock.Anything).Return(nil)

	err := hooks.MemStore(0, []byte("data"))
	require.Nil(t, err)
}

func TestVMHooksImpl_Getters(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, output, storage := createTestVMHooks()

	require.Equal(t, host, hooks.GetVMHost())
	require.Equal(t, host.Blockchain(), hooks.GetBlockchainContext())
	require.Equal(t, runtime, hooks.GetRuntimeContext())
	require.Equal(t, host.ManagedTypes(), hooks.GetManagedTypesContext())
	require.Equal(t, output, hooks.GetOutputContext())
	require.Equal(t, metering, hooks.GetMeteringContext())
	require.Equal(t, storage, hooks.GetStorageContext())
	require.Equal(t, host.EnableEpochsHandler(), hooks.GetEnableEpochsHandler())
}

func TestVMHooksImpl_FailExecution(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("GasLeft").Return(uint64(1000))
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	runtime.On("FailExecution", mock.Anything).Return()

	hooks.FailExecution(errors.New("test error"))
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestFailExecution(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("GasLeft").Return(uint64(1000))
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	runtime.On("FailExecution", mock.Anything).Return()

	FailExecution(hooks.GetVMHost(), errors.New("test error"))
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestVMHooksImpl_FailExecutionConditionally_SafeMode(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("GasLeft").Return(uint64(1000))
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	runtime.On("FailExecution", mock.Anything).Return()

	hooks.FailExecutionConditionally(errors.New("test error"))
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestVMHooksImpl_FailExecutionConditionally_NotActive(t *testing.T) {
	t.Parallel()
	vmHooksMockery := &mockeryStruct{}

	vmHooksMockery.host = &mockery.MockVMHost{}
	vmHooksMockery.runtime = &mockery.MockRuntimeContext{}
	vmHooksMockery.metering = &mockery.MockMeteringContext{}

	vmHooksMockery.host.On("Metering").Return(vmHooksMockery.metering)
	vmHooksMockery.host.On("Runtime").Return(vmHooksMockery.runtime)
	vmHooksMockery.runtime.On("FailExecution", mock.Anything).Return()
	vmHooksMockery.host.On("EnableEpochsHandler").Return(&mock2.EnableEpochsHandlerStub{
		IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
			return false
		},
	})
	vmHooksMockery.metering.On("GasLeft").Return(uint64(1000))
	vmHooksMockery.metering.On("UseGasBounded", mock.Anything).Return(nil)
	vmHooksMockery.hooks = NewVMHooksImpl(vmHooksMockery.host)

	vmHooksMockery.runtime.On("IsUnsafeMode").Return(true)
	vmHooksMockery.hooks.FailExecutionConditionally(errors.New("test error"))
	vmHooksMockery.runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestVMHooksImpl_FailExecutionConditionally_UnSafeMode(t *testing.T) {
	t.Parallel()

	vmHooksMockery := &mockeryStruct{}
	vmHooksMockery.host = &mockery.MockVMHost{}
	vmHooksMockery.runtime = &mockery.MockRuntimeContext{}
	vmHooksMockery.metering = &mockery.MockMeteringContext{}

	vmHooksMockery.host.On("Metering").Return(vmHooksMockery.metering)
	vmHooksMockery.host.On("Runtime").Return(vmHooksMockery.runtime)
	vmHooksMockery.runtime.On("FailExecution", mock.Anything).Return()
	vmHooksMockery.host.On("EnableEpochsHandler").Return(&mock2.EnableEpochsHandlerStub{
		IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
			return true
		},
	})
	vmHooksMockery.metering.On("GasLeft").Return(uint64(1000))
	vmHooksMockery.metering.On("UseGasBounded", mock.Anything).Return(nil)
	vmHooksMockery.hooks = NewVMHooksImpl(vmHooksMockery.host)

	vmHooksMockery.runtime.On("IsUnsafeMode").Return(true)
	vmHooksMockery.runtime.On("AddError", mock.Anything, mock.Anything)
	vmHooksMockery.hooks.FailExecutionConditionally(errors.New("test error"))
	vmHooksMockery.runtime.AssertNotCalled(t, "FailExecution", mock.Anything)
	vmHooksMockery.runtime.AssertCalled(t, "AddError", mock.Anything, mock.Anything)
}
