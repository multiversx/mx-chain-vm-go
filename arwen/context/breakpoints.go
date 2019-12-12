package context

import (
	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

func (host *vmContext) reachedBreakpoint(err error) bool {
	return err != nil && host.GetRuntimeBreakpointValue() != arwen.BreakpointNone
}

func (host *vmContext) handleBreakpoint(result wasmer.Value, err error) (*vmcommon.VMOutput, error) {
	breakpointValue := host.GetRuntimeBreakpointValue()

	if breakpointValue == arwen.BreakpointAsyncCall {
		return host.handleAsyncCallBreakpoint(result, err)
	}

	return nil, ErrUnhandledRuntimeBreakpoint
}

func (host *vmContext) handleAsyncCallBreakpoint(result wasmer.Value, err error) (*vmcommon.VMOutput, error) {
	host.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	convertedResult := arwen.ConvertReturnValue(result)
	callerSCVmOutput := host.createVMOutput(convertedResult.Bytes())

	sender := host.GetSCAddress()
	dest := host.asyncCallDest

	calledSCCode, err := host.GetCode(dest)
	if err != nil || len(calledSCCode) == 0 {
		return callerSCVmOutput, nil
	}

	return nil, nil
}
