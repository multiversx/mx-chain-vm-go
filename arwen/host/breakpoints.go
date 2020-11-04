package host

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
)

func (host *vmHost) handleBreakpointIfAny(executionErr error) error {
	if executionErr == nil {
		return nil
	}

	runtime := host.Runtime()
	breakpointValue := runtime.GetRuntimeBreakpointValue()
	log.Info("handleBreakpointIfAny", "value", breakpointValue)
	if breakpointValue != arwen.BreakpointNone {
		executionErr = host.handleBreakpoint(breakpointValue)
	}

	return executionErr
}

func (host *vmHost) handleBreakpoint(breakpointValue arwen.BreakpointValue) error {
	log.Info("handleBreakPoint", "value", breakpointValue)
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
		log.Info("OUT OF GAS breakpoint")
		return arwen.ErrNotEnoughGas
	}

	return arwen.ErrUnhandledRuntimeBreakpoint
}
