package host

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) handleAsyncCallBreakpoint(result wasmer.Value) error {
	runtime := host.Runtime()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	// TODO also determine whether caller and callee are in the same Shard (based
	// on address?), by account addresses - this would make the empty SC code an
	// unrecoverable error, so returning nil here will not be appropriate anymore.
	if !host.canExecuteSynchronously() {
		host.setAsyncCallToDestination()
		return nil
	}

	// Start calling the destination SC, synchronously.
	destinationCallInput, err := host.createDestinationContractCallInput()
	if err != nil {
		return err
	}

	destinationVMOutput, err := host.ExecuteOnDestContext(destinationCallInput)
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

func (host *vmHost) canExecuteSynchronously() bool {
	// TODO replace with a blockchain hook that verifies if the caller and callee
	// are in the same Shard.
	runtime := host.Runtime()
	blockchain := host.Blockchain()
	asyncCallInfo := runtime.GetAsyncCallInfo()
	dest := asyncCallInfo.Destination
	calledSCCode, err := blockchain.GetCode(dest)

	return len(calledSCCode) != 0 && err == nil
}

func (host *vmHost) setAsyncCallToDestination() {
	runtime := host.Runtime()
	output := host.Output()
	destination := runtime.GetAsyncCallInfo().Destination
	destinationAccount, _ := output.GetOutputAccount(destination)
	destinationAccount.CallType = vmcommon.AsynchronousCall
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

	gasLimit := asyncCallInfo.GasLimit
	gasToUse := host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0).SetBytes(asyncCallInfo.ValueBytes),
			CallType:    vmcommon.AsynchronousCall,
			GasPrice:    runtime.GetVMInput().GasPrice,
			GasProvided: gasLimit,
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

	dataLength := host.computeDataLengthFromArguments(function, arguments)

	gasToUse := metering.GasSchedule().ElrondAPICost.AsyncCallStep
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(dataLength)
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	sender := runtime.GetAsyncCallInfo().Destination
	dest := runtime.GetSCAddress()

	// Return to the sender SC, calling its callback() method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.AsynchronousCallBack,
			GasPrice:    runtime.GetVMInput().GasPrice,
			GasProvided: gasLimit,
		},
		RecipientAddr: dest,
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmHost) computeDataLengthFromArguments(function string, arguments [][]byte) int {
	// Calculate what length would the Data field have, were it of the
	// form "callback@arg1@arg4...

	// TODO change this after allowing empty entries in VMCommon.ReturnData
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

	return dataLength
}
