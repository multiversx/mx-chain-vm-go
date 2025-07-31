package hostCore

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createTestHostWithMocks() (*vmHost, *context.MockRuntimeContext, *context.MockAsyncContext) {
	runtimeMock := &context.MockRuntimeContext{}
	asyncMock := &context.MockAsyncContext{}
	host := &vmHost{
		runtimeContext: runtimeMock,
		asyncContext:   asyncMock,
	}
	return host, runtimeMock, asyncMock
}

func TestVmHost_handleBreakpoint(t *testing.T) {
	t.Parallel()

	t.Run("async call", func(t *testing.T) {
		t.Parallel()
		h, runtimeMock, asyncMock := createTestHostWithMocks()
		runtimeMock.On("SetRuntimeBreakpointValue", vmhost.BreakpointNone)
		asyncMock.On("GetCallGroup", vmhost.LegacyAsyncCallGroupID).Return(nil, false)
		err := h.handleBreakpoint(vmhost.BreakpointAsyncCall)
		require.Equal(t, vmhost.ErrLegacyAsyncCallNotFound, err)
	})

	errTests := []struct {
		breakpoint vmhost.BreakpointValue
		expected   error
	}{
		{vmhost.BreakpointExecutionFailed, vmhost.ErrExecutionFailed},
		{vmhost.BreakpointSignalError, vmhost.ErrSignalError},
		{vmhost.BreakpointOutOfGas, vmhost.ErrNotEnoughGas},
		{vmhost.BreakpointMemoryLimit, vmhost.ErrMemoryLimit},
		{vmhost.BreakpointValue(99), vmhost.ErrUnhandledRuntimeBreakpoint},
	}

	for _, tt := range errTests {
		tt := tt
		t.Run(tt.expected.Error(), func(t *testing.T) {
			t.Parallel()
			h, _, _ := createTestHostWithMocks()
			err := h.handleBreakpoint(tt.breakpoint)
			require.Equal(t, tt.expected, err)
		})
	}
}

func TestVmHost_handleBreakpointIfAny(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		t.Parallel()
		h, _, _ := createTestHostWithMocks()
		err := h.handleBreakpointIfAny(nil)
		require.Nil(t, err)
	})

	t.Run("with error, no breakpoint", func(t *testing.T) {
		t.Parallel()
		h, runtimeMock, _ := createTestHostWithMocks()
		runtimeMock.On("GetRuntimeBreakpointValue").Return(vmhost.BreakpointNone)
		runtimeMock.On("AddError", mock.Anything, mock.Anything).Return()
		runtimeMock.On("FunctionName").Return("testFunc")
		execErr := errors.New("exec error")
		err := h.handleBreakpointIfAny(execErr)
		require.Equal(t, vmhost.ErrExecutionFailed, err)
	})

	t.Run("with error and breakpoint", func(t *testing.T) {
		t.Parallel()
		h, runtimeMock, asyncMock := createTestHostWithMocks()
		runtimeMock.On("GetRuntimeBreakpointValue").Return(vmhost.BreakpointOutOfGas)
		runtimeMock.On("AddError", mock.Anything, mock.Anything).Return()
		runtimeMock.On("FunctionName").Return("testFunc")
		asyncMock.On("GetCallGroup", vmhost.LegacyAsyncCallGroupID).Return(nil, false)
		execErr := errors.New("exec error")
		err := h.handleBreakpointIfAny(execErr)
		require.Equal(t, vmhost.ErrNotEnoughGas, err)
	})
}

func TestVmHost_handleAsyncCallBreakpoint(t *testing.T) {
	t.Parallel()

	t.Run("group not found", func(t *testing.T) {
		t.Parallel()
		h, runtimeMock, asyncMock := createTestHostWithMocks()
		runtimeMock.On("SetRuntimeBreakpointValue", vmhost.BreakpointNone).Return()
		asyncMock.On("GetCallGroup", vmhost.LegacyAsyncCallGroupID).Return(nil, false)
		err := h.handleAsyncCallBreakpoint()
		require.Equal(t, vmhost.ErrLegacyAsyncCallNotFound, err)
	})

	t.Run("group is complete", func(t *testing.T) {
		t.Parallel()
		h, runtimeMock, asyncMock := createTestHostWithMocks()
		group := &vmhost.AsyncCallGroup{}
		runtimeMock.On("SetRuntimeBreakpointValue", vmhost.BreakpointNone).Return()
		asyncMock.On("GetCallGroup", vmhost.LegacyAsyncCallGroupID).Return(group, true)
		err := h.handleAsyncCallBreakpoint()
		require.Equal(t, vmhost.ErrLegacyAsyncCallInvalid, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		h, runtimeMock, asyncMock := createTestHostWithMocks()
		group := &vmhost.AsyncCallGroup{AsyncCalls: []*vmhost.AsyncCall{{}}}
		runtimeMock.On("SetRuntimeBreakpointValue", vmhost.BreakpointNone).Return()
		asyncMock.On("GetCallGroup", vmhost.LegacyAsyncCallGroupID).Return(group, true)
		err := h.handleAsyncCallBreakpoint()
		require.Nil(t, err)
	})
}
