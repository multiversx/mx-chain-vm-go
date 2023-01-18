package host

import (
	"github.com/multiversx/wasm-vm/arwen"
)

func (host *vmHost) handleBreakpointIfAny(executionErr error) error {
	if executionErr == nil {
		return nil
	}

	runtime := host.Runtime()
	breakpointValue := runtime.GetRuntimeBreakpointValue()
	log.Trace("handleBreakpointIfAny", "value", breakpointValue)
	if breakpointValue != arwen.BreakpointNone {
		err := host.handleBreakpoint(breakpointValue)
		runtime.AddError(err, runtime.FunctionName())
		return err
	}

	log.Trace("wasmer execution error", "err", executionErr)
	runtime.AddError(executionErr, runtime.FunctionName())
	return arwen.ErrExecutionFailed
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
	if breakpointValue == arwen.BreakpointMemoryLimit {
		return arwen.ErrMemoryLimit
	}

	return arwen.ErrUnhandledRuntimeBreakpoint
}
