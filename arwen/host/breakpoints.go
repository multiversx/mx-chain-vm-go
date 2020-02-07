package host

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

func (host *vmHost) handleBreakpoint(
	breakpointValue arwen.BreakpointValue,
	result wasmer.Value,
) error {
	if breakpointValue == arwen.BreakpointAsyncCall {
		return host.handleAsyncCallBreakpoint(result)
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
