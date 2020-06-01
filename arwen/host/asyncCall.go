package host

import (
	"encoding/hex"
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
		return host.sendAsyncCallToDestination(runtime.GetAsyncCallInfo())
	}

	// Start calling the destination SC, synchronously.
	destinationCallInput, err := host.createDestinationContractCallInput(runtime.GetAsyncCallInfo())
	if err != nil {
		return err
	}

	// TODO: If this generates async calls before jumping out here, we should execute those calls?
	destinationVMOutput, _, destinationErr := host.ExecuteOnDestContext(destinationCallInput)

	callbackCallInput, err := host.createCallbackContractCallInput(
		destinationVMOutput,
		runtime.GetAsyncCallInfo().Destination,
		arwen.CallbackDefault,
		destinationErr,
	)
	if err != nil {
		return err
	}

	callbackVMOutput, _, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	err = host.processCallbackVMOutput(callbackVMOutput, callBackErr)
	if err != nil {
		return err
	}

	return nil
}

func (host *vmHost) canExecuteSynchronouslyOnDest(dest []byte) bool {
	blockchain := host.Blockchain()
	calledSCCode, err := blockchain.GetCode(dest)

	return len(calledSCCode) != 0 && err == nil
}

func (host *vmHost) canExecuteSynchronously() bool {
	// TODO replace with a blockchain hook that verifies if the caller and callee
	// are in the same Shard.
	runtime := host.Runtime()
	asyncCallInfo := runtime.GetAsyncCallInfo()
	dest := asyncCallInfo.Destination

	return host.canExecuteSynchronouslyOnDest(dest)
}

func (host *vmHost) sendAsyncCallToDestination(asyncCallInfo arwen.AsyncCallInfoHandler) error {
	runtime := host.Runtime()
	output := host.Output()

	destination := asyncCallInfo.GetDestination()
	destinationAccount, _ := output.GetOutputAccount(destination)
	destinationAccount.CallType = vmcommon.AsynchronousCall

	err := output.Transfer(
		destination,
		runtime.GetSCAddress(),
		asyncCallInfo.GetGasLimit(),
		big.NewInt(0).SetBytes(asyncCallInfo.GetValueBytes()),
		asyncCallInfo.GetData(),
	)
	if err != nil {
		metering := host.Metering()
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}


func (host *vmHost) sendCallbackToDestination() error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	destination := currentCall.CallerAddr
	destinationAccount, _ := output.GetOutputAccount(destination)
	destinationAccount.CallType = vmcommon.AsynchronousCallBack

	retData := []byte("@" + hex.EncodeToString([]byte(output.ReturnCode().String())))
	for _, data := range output.ReturnData() {
		retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
	}

	err := output.Transfer(
		destination,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		big.NewInt(0).SetUint64(metering.GasLeft()*currentCall.GasPrice),
		retData,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

func (host *vmHost) sendStorageCallbackToDestination(asyncCallInitiator vmcommon.AsyncInitiator) error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	destination := asyncCallInitiator.CallerAddr
	destinationAccount, _ := output.GetOutputAccount(destination)
	destinationAccount.CallType = vmcommon.AsynchronousCallBack

	err := output.Transfer(
		destination,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		big.NewInt(0).SetUint64(metering.GasLeft()*currentCall.GasPrice),
		asyncCallInitiator.ReturnData,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

func (host *vmHost) createDestinationContractCallInput(asyncCallInfo arwen.AsyncCallInfoHandler) (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()

	argParser := runtime.ArgParser()
	err := argParser.ParseData(string(asyncCallInfo.GetData()))
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

	gasLimit := asyncCallInfo.GetGasLimit()
	gasToUse := host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0).SetBytes(asyncCallInfo.GetValueBytes()),
			CallType:       vmcommon.AsynchronousCall,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
		},
		RecipientAddr: asyncCallInfo.GetDestination(),
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmHost) createCallbackContractCallInput(
	destinationVMOutput *vmcommon.VMOutput,
	callbackInitiator []byte,
	callbackFunction string,
	destinationErr error,
) (*vmcommon.ContractCallInput, error) {
	metering := host.Metering()
	runtime := host.Runtime()

	// always provide return code as the first argument to callback function
	arguments := [][]byte{
		big.NewInt(int64(destinationVMOutput.ReturnCode)).Bytes(),
	}
	if destinationErr == nil {
		// when execution went Ok, callBack arguments are:
		// [0, result1, result2, ....]
		arguments = append(arguments, destinationVMOutput.ReturnData...)
	} else {
		// when execution returned error, callBack arguments are:
		// [error code, error message]
		arguments = append(arguments, []byte(destinationVMOutput.ReturnMessage))
	}

	gasLimit := destinationVMOutput.GasRemaining
	dataLength := host.computeDataLengthFromArguments(callbackFunction, arguments)

	gasToUse := metering.GasSchedule().ElrondAPICost.AsyncCallStep
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(dataLength)
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	// Return to the sender SC, calling its callback() method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     callbackInitiator,
			Arguments:      arguments,
			CallValue:      big.NewInt(0),
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      callbackFunction,
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
