package host

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
)

func (host *vmHost) handleBreakpointIfAny(executionErr error) error {
	if executionErr == nil {
		return nil
	}

	runtime := host.Runtime()
	breakpointValue := runtime.GetRuntimeBreakpointValue()
	if breakpointValue != arwen.BreakpointNone {
		if breakpointValue != arwen.BreakpointAsyncCall {
			runtime.AddError(executionErr, runtime.Function())
			return host.handleBreakpoint(breakpointValue)
		}
		return host.handleBreakpoint(breakpointValue)
	}

	return executionErr
}

func (host *vmHost) handleBreakpoint(breakpointValue arwen.BreakpointValue) error {
	if breakpointValue == arwen.BreakpointAsyncCall {
		return host.handleAsyncCallBreakpoint()
	}
	if breakpointValue == arwen.BreakpointExecutionFailed {
		return arwen.ErrExecutionFailed
	}
	if breakpointValue == arwen.BreakpointSignalError {
		return arwen.ErrSignalError
	}
	if breakpointValue == arwen.BreakpointOutOfGas {
		return arwen.ErrNotEnoughGas
	}

	return arwen.ErrUnhandledRuntimeBreakpoint
}
