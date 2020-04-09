package host

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	// TODO also determine whether caller and callee are in the same Shard (based
	// on address?), by account addresses - this would make the empty SC code an
	// unrecoverable error, so returning nil here will not be appropriate anymore.
	if !host.canExecuteSynchronously() {
		return host.sendAsyncCallToDestination()
	}

	// Start calling the destination SC, synchronously.
	destinationCallInput, err := host.createDestinationContractCallInput()
	if err != nil {
		return err
	}

	destinationVMOutput, destinationErr := host.ExecuteOnDestContext(destinationCallInput)

	callbackCallInput, err := host.createCallbackContractCallInput(destinationVMOutput, destinationErr)
	if err != nil {
		return err
	}

	callbackVMOutput, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	err = host.processCallbackVMOutput(callbackVMOutput, callBackErr)
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

func (host *vmHost) sendAsyncCallToDestination() error {
	runtime := host.Runtime()
	output := host.Output()

	asyncCallInfo := runtime.GetAsyncCallInfo()
	destination := asyncCallInfo.Destination
	destinationAccount, _ := output.GetOutputAccount(destination)
	destinationAccount.CallType = vmcommon.AsynchronousCall

	err := output.Transfer(
		destination,
		runtime.GetSCAddress(),
		asyncCallInfo.GasLimit,
		big.NewInt(0).SetBytes(asyncCallInfo.ValueBytes),
		asyncCallInfo.Data,
	)
	if err != nil {
		metering := host.Metering()
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
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

	arguments, err := argParser.GetFunctionArguments()
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
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0).SetBytes(asyncCallInfo.ValueBytes),
			CallType:       vmcommon.AsynchronousCall,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
		},
		RecipientAddr: asyncCallInfo.Destination,
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmHost) createCallbackContractCallInput(destinationVMOutput *vmcommon.VMOutput, destinationErr error) (*vmcommon.ContractCallInput, error) {
	metering := host.Metering()
	runtime := host.Runtime()

	arguments := destinationVMOutput.ReturnData
	gasLimit := destinationVMOutput.GasRemaining
	function := "callBack"

	if destinationErr != nil {
		arguments = [][]byte{
			[]byte(destinationVMOutput.ReturnCode.String()),
			[]byte(runtime.GetCurrentTxHash()),
		}
	}

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
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0),
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
		},
		RecipientAddr: dest,
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmHost) processCallbackVMOutput(callbackVMOutput *vmcommon.VMOutput, callBackErr error) error {
	if callBackErr == nil {
		return nil
	}

	runtime := host.Runtime()
	output := host.Output()

	runtime.GetVMInput().GasProvided = 0
	output.Finish([]byte(callbackVMOutput.ReturnCode.String()))
	output.Finish([]byte(runtime.GetCurrentTxHash()))

	return nil
}

func (host *vmHost) computeDataLengthFromArguments(function string, arguments [][]byte) int {
	// Calculate what length would the Data field have, were it of the
	// form "callback@arg1@arg4...

	// TODO this needs tests, especially for the case when the arguments slice
	// contains an empty []byte
	numSeparators := len(arguments)
	dataLength := len(function) + numSeparators
	for _, element := range arguments {
		dataLength += len(element)
	}

	return dataLength
}
