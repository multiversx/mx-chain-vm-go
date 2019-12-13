package context

import (
	"math/big"

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
	senderVMOutput := host.createVMOutput(convertedResult.Bytes())

	dest := host.asyncCallDest

	// If SC code is not found, it means this is either a cross-shard call or a wrong call.
	calledSCCode, err := host.GetCode(dest)
	if err != nil || len(calledSCCode) == 0 {
		return senderVMOutput, nil
	}

	// Start calling the destination SC, synchronously.
	sender := host.GetSCAddress()
	valueBytes := host.asyncCallValueBytes
	data := host.asyncCallData
	gasLimit := host.asyncCallGasLimit

	err = host.argParser.ParseData(string(data))
	if err != nil {
		return createVMOutputInCaseOfBreakpointError(err), nil
	}

	function, err := host.argParser.GetFunction()
	if err != nil {
		return createVMOutputInCaseOfBreakpointError(err), nil
	}

	arguments, err := host.argParser.GetArguments()
	if err != nil {
		return createVMOutputInCaseOfBreakpointError(err), nil
	}

	destinationCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0).SetBytes(valueBytes),
			GasPrice:    0,
			GasProvided: host.BoundGasLimit(int64(gasLimit)),
		},
		RecipientAddr: dest,
		Function:      function,
	}

	destinationVMOutput, err := host.executeOnNewContextAndGetVMOutput(destinationCallInput)

	if err != nil {
		return createVMOutputInCaseOfBreakpointError(err), nil
	}

	// Return to the sender SC, calling its callback() method.
	callbackCallInput := &vmcommon.ContractCallInput{}

	var callbackVMOutput *vmcommon.VMOutput

	mergedVMOutput := mergeVMOutputs(senderVMOutput, destinationVMOutput, callbackVMOutput)

	return mergedVMOutput, nil
}

func (host *vmContext) executeOnNewContextAndGetVMOutput(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	currVmInput := host.vmInput
	currScAddress := host.scAddress
	currCallFunction := host.callFunction

	currContext := host.copyToNewContext()

	defer func() {
		// Restore the original host context
		host.copyFromContext(currContext)
		host.vmInput = currVmInput
		host.scAddress = currScAddress
		host.callFunction = currCallFunction
	}()

	host.vmInput = input.VMInput
	host.scAddress = input.RecipientAddr
	host.callFunction = input.Function

	host.initInternalValues()

	var vmOutput *vmcommon.VMOutput
	err := host.execute(input)
	if err != nil {
		vmOutput = host.createVMOutput(make([]byte, 0))
	}

	return vmOutput, err
}

func mergeVMOutputs(sender *vmcommon.VMOutput, destination *vmcommon.VMOutput, callback *vmcommon.VMOutput) *vmcommon.VMOutput {
	return nil
}

func createVMOutputInCaseOfBreakpointError(err error) *vmcommon.VMOutput {
	return nil
}
