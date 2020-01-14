package host

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

func (host *vmHost) handleBreakpoint(
	breakpointValue arwen.BreakpointValue,
	result wasmer.Value,
	err error,
) error {

	if breakpointValue == arwen.BreakpointAsyncCall {
		// return host.handleAsyncCallBreakpoint(result, err)
	}

	if breakpointValue == arwen.BreakpointSignalError {
		return nil
	}

	if breakpointValue == arwen.BreakpointSignalExit {
		return nil
	}

	return arwen.ErrUnhandledRuntimeBreakpoint
}
