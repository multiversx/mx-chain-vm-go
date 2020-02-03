package host

import (
	"math/big"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) handleAsyncCallBreakpoint(result wasmer.Value, argError error) error {
	runtime := host.Runtime()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	// TODO also determine whether caller and callee are in the same Shard, by
	// account addresses - this makes the empty SC code an error
	syncCall, err := host.canExecuteSynchronously()
	if err != nil {
		return err
	}
	if !syncCall {
		return argError
	}

	// Start calling the destination SC, synchronously.
	destinationCallInput, err := host.createDestinationContractCallInput()
	if err != nil {
		return err
	}

	destinationVMOutput, err := host.ExecuteOnDestContext(destinationCallInput)
	// TODO pass error to the SC callback (append as argument)
	// TODO consume remaining gas
	if err != nil {
		return err
	}

	callbackCallInput, err := host.createCallbackContractCallInput(destinationVMOutput)
	if err != nil {
		return err
	}

	_, err = host.ExecuteOnDestContext(callbackCallInput)
	if err != nil {
		return err
	}

	return nil
}

func (host *vmHost) canExecuteSynchronously() (bool, error) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()
	asyncCallInfo := runtime.GetAsyncCallInfo()
	dest := asyncCallInfo.Destination
	calledSCCode, err := blockchain.GetCode(dest)

	return len(calledSCCode) != 0, err
}

func (host *vmHost) createDestinationContractCallInput() (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()
	asyncCallInfo := runtime.GetAsyncCallInfo()

	argParser := runtime.ArgParser()
	err := argParser.ParseData(string(asyncCallInfo.Data))
	if err != nil {
		return nil, err
	}

	function, err := argParser.GetFunction()
	if err != nil {
		return nil, err
	}

	arguments, err := argParser.GetArguments()
	if err != nil {
		return nil, err
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0).SetBytes(asyncCallInfo.ValueBytes),
			GasPrice:    0,
			GasProvided: asyncCallInfo.GasLimit,
		},
		RecipientAddr: asyncCallInfo.Destination,
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmHost) createCallbackContractCallInput(destinationVMOutput *vmcommon.VMOutput) (*vmcommon.ContractCallInput, error) {
	metering := host.Metering()
	runtime := host.Runtime()

	arguments := destinationVMOutput.ReturnData
	gasLimit := destinationVMOutput.GasRemaining
	function := "callBack"

	// Calculate what length would the Data field have, were it of the
	// form "callback@arg1@arg4...
	dataLength := len(function) + 1
	for i, element := range arguments {
		if len(element) == 0 {
			continue
		}
		if i != 0 && dataLength > 0 {
			dataLength += 1
		}
		dataLength += len(element)
	}

	// TODO define gas cost specific to async calls - asyncCall should cost more
	// than ExecuteOnDestContext, especially cross-shard async calls
	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(dataLength)
	gasLimit -= gasToUse

	sender := runtime.GetAsyncCallInfo().Destination
	dest := runtime.GetSCAddress()

	// Return to the sender SC, calling its callback() method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0),
			GasPrice:    0,
			GasProvided: gasLimit,
		},
		RecipientAddr: dest,
		Function:      function,
	}

	return contractCallInput, nil
}
