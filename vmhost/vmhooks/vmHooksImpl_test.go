package vmhooks

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	vmcommonmock "github.com/multiversx/mx-chain-vm-common-go/mock"
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

func TestVMHooksImpl_FailExecutionConditionally(t *testing.T) {
	t.Parallel()

	errTest := errors.New("test error")

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		mockRuntime := &mockRuntimeContext{
			failExecutionCalled: func(err error) {
				wasCalled = true
			},
			addErrorCalled: func(err error, otherInfo ...string) {
				wasCalled = true
			},
		}

		host := &mockery.MockVMHost{}
		host.On("Runtime").Return(mockRuntime)
		hooks := NewVMHooksImpl(host)

		hooks.FailExecutionConditionally(nil)

		require.False(t, wasCalled)
	})

	t.Run("unsafe mode OFF", func(t *testing.T) {
		t.Parallel()

		failExecutionCalled := false
		addErrorCalled := false
		mockRuntime := &mockRuntimeContext{
			isUnsafeModeCalled: func() bool {
				return false
			},
			failExecutionCalled: func(err error) {
				failExecutionCalled = true
			},
			addErrorCalled: func(err error, otherInfo ...string) {
				addErrorCalled = true
			},
		}

		host := &mockery.MockVMHost{}
		host.On("Runtime").Return(mockRuntime)
		metering := &mockery.MockMeteringContext{}
		metering.On("GasLeft").Return(uint64(1000))
		metering.On("UseGasBounded", mock.Anything).Return(nil)
		host.On("Metering").Return(metering)
		hooks := NewVMHooksImpl(host)

		hooks.FailExecutionConditionally(errTest)

		require.True(t, failExecutionCalled)
		require.False(t, addErrorCalled)
	})

	t.Run("unsafe mode ON, flag OFF", func(t *testing.T) {
		t.Parallel()

		failExecutionCalled := false
		addErrorCalled := false
		mockRuntime := &mockRuntimeContext{
			isUnsafeModeCalled: func() bool {
				return true
			},
			failExecutionCalled: func(err error) {
				failExecutionCalled = true
			},
			addErrorCalled: func(err error, otherInfo ...string) {
				addErrorCalled = true
			},
		}

		host := &mockery.MockVMHost{}
		host.On("Runtime").Return(mockRuntime)
		metering := &mockery.MockMeteringContext{}
		metering.On("GasLeft").Return(uint64(1000))
		metering.On("UseGasBounded", mock.Anything).Return(nil)
		host.On("Metering").Return(metering)
		enableEpochsHandler := &vmcommonmock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return false
			},
		}
		host.On("EnableEpochsHandler").Return(enableEpochsHandler)
		hooks := NewVMHooksImpl(host)

		hooks.FailExecutionConditionally(errTest)

		require.True(t, failExecutionCalled)
		require.False(t, addErrorCalled)
	})

	t.Run("unsafe mode ON, flag ON", func(t *testing.T) {
		t.Parallel()

		failExecutionCalled := false
		addErrorCalled := false
		mockRuntime := &mockRuntimeContext{
			isUnsafeModeCalled: func() bool {
				return true
			},
			failExecutionCalled: func(err error) {
				failExecutionCalled = true
			},
			addErrorCalled: func(err error, otherInfo ...string) {
				addErrorCalled = true
			},
		}

		host := &mockery.MockVMHost{}
		host.On("Runtime").Return(mockRuntime)
		enableEpochsHandler := &vmcommonmock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		}
		host.On("EnableEpochsHandler").Return(enableEpochsHandler)
		hooks := NewVMHooksImpl(host)

		hooks.FailExecutionConditionally(errTest)

		require.False(t, failExecutionCalled)
		require.True(t, addErrorCalled)
	})
}
