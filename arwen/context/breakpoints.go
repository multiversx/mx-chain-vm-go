package context

import (
	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

func (host *vmContext) handleBreakpoint(
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

	return ErrUnhandledRuntimeBreakpoint
}
