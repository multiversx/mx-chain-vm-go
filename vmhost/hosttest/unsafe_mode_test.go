package hostCoretest

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/assert"
)

type unsafeVMHostMock struct {
	vmhost.VMHost
	unsafeMode          bool
	failExecutionCalled bool
	breakpointValue     vmhost.BreakpointValue
}

func (h *unsafeVMHostMock) SetUnsafeMode(unsafeMode bool) {
	h.unsafeMode = unsafeMode
}

func (h *unsafeVMHostMock) IsUnsafeMode() bool {
	return h.unsafeMode
}

func (h *unsafeVMHostMock) FailExecution(err error) {
	h.failExecutionCalled = true
	h.breakpointValue = vmhost.BreakpointExecutionFailed
}

func (h *unsafeVMHostMock) FailExecutionConditionally(err error) {
	if !h.unsafeMode {
		h.FailExecution(err)
	}
}

func (h *unsafeVMHostMock) Runtime() vmhost.RuntimeContext {
	return &unsafeRuntimeContextMock{
		VMHost: h,
	}
}

type unsafeRuntimeContextMock struct {
	vmhost.RuntimeContext
	VMHost vmhost.VMHost
}

func (r *unsafeRuntimeContextMock) FailExecution(err error) {
	r.VMHost.Runtime().FailExecution(err)
}

func (r *unsafeRuntimeContextMock) FailExecutionConditionally(err error) {
	r.VMHost.Runtime().FailExecutionConditionally(err)
}

func (r *unsafeRuntimeContextMock) GetRuntimeBreakpointValue() vmhost.BreakpointValue {
	return r.VMHost.(*unsafeVMHostMock).breakpointValue
}

func TestUnsafeMode(t *testing.T) {
	t.Run("safe mode", func(t *testing.T) {
		host := &unsafeVMHostMock{}
		host.SetUnsafeMode(false)

		host.FailExecutionConditionally(errors.New("test error"))

		assert.True(t, host.failExecutionCalled)
		assert.Equal(t, vmhost.BreakpointExecutionFailed, host.breakpointValue)
	})

	t.Run("unsafe mode", func(t *testing.T) {
		host := &unsafeVMHostMock{}
		host.SetUnsafeMode(true)

		host.FailExecutionConditionally(errors.New("test error"))

		assert.False(t, host.failExecutionCalled)
		assert.Equal(t, vmhost.BreakpointNone, host.breakpointValue)
	})
}
