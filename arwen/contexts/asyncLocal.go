package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (context *asyncContext) executeAsyncLocalCalls() error {
	localCalls := make([]*arwen.AsyncCall, 0)

	for _, group := range context.asyncCallGroups {
		for _, call := range group.AsyncCalls {
			if call.IsLocal() {
				localCalls = append(localCalls, call)
			}
		}
	}

	for _, call := range localCalls {
		err := context.executeAsyncLocalCall(call)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO split this method into smaller ones
func (context *asyncContext) executeAsyncLocalCall(asyncCall *arwen.AsyncCall) error {
	if asyncCall.ExecutionMode == arwen.ESDTTransferOnCallBack {
		context.executeESDTTransferOnCallback(asyncCall)
		context.completeChild(asyncCall.CallID, 0)
		return nil
	}

	destinationCallInput, err := context.createContractCallInput(asyncCall)
	if err != nil {
		logAsync.Trace("executeAsyncLocalCall failed", "error", err)
		return err
	}

	logAsync.Trace("executeAsyncLocalCall",
		"caller", destinationCallInput.CallerAddr,
		"dest", destinationCallInput.RecipientAddr,
		"func", destinationCallInput.Function,
		"args", destinationCallInput.Arguments,
		"gasProvided", destinationCallInput.GasProvided,
		"gasLocked", destinationCallInput.GasLocked)

	// Briefly restore the AsyncCall GasLimit, after it was consumed in its
	// entirety by addAsyncCall(); this is required, because ExecuteOnDestContext()
	// must also consume the GasLimit in its entirety, before starting execution,
	// but will restore any GasRemaining to the current instance.
	metering := context.host.Metering()
	metering.RestoreGas(asyncCall.GetGasLimit())

	vmOutput, isComplete, err := context.host.ExecuteOnDestContext(destinationCallInput)
	if vmOutput == nil {
		return arwen.ErrNilDestinationCallVMOutput
	}

	logAsync.Trace("executeAsyncLocalCall",
		"retCode", vmOutput.ReturnCode,
		"message", vmOutput.ReturnMessage,
		"data", vmOutput.ReturnData,
		"gasRemaining", vmOutput.GasRemaining,
		"error", err)

	asyncCall.UpdateStatus(vmOutput.ReturnCode)

	if isComplete {
		if asyncCall.HasCallback() {
			// Restore gas locked while still on the caller instance; otherwise, the
			// locked gas will appear to have been used twice by the caller instance.
			isCallbackComplete, callbackVMOutput := context.ExecuteSyncCallbackAndFinishOutput(asyncCall, vmOutput, destinationCallInput, 0, err)
			if callbackVMOutput == nil {
				return arwen.ErrAsyncNoOutputFromCallback
			}

			if isCallbackComplete {
				callbackGasRemaining := callbackVMOutput.GasRemaining
				callbackVMOutput.GasRemaining = 0
				return context.completeChild(asyncCall.CallID, callbackGasRemaining)
			}
		} else {
			return context.completeChild(asyncCall.CallID, 0)
		}
	}

	return nil
}

// TODO rename to executeLocalCallbackAndFinishOutput
func (context *asyncContext) ExecuteSyncCallbackAndFinishOutput(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	destinationCallInput *vmcommon.ContractCallInput,
	gasAccumulated uint64,
	err error) (bool, *vmcommon.VMOutput) {
	callbackVMOutput, isComplete, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, gasAccumulated, err)
	context.finishAsyncLocalCallbackExecution(callbackVMOutput, callbackErr, vmOutput.ReturnCode)
	return isComplete, callbackVMOutput
}

// TODO rename to executeLocalCallback
func (context *asyncContext) executeSyncCallback(
	asyncCall *arwen.AsyncCall,
	destinationVMOutput *vmcommon.VMOutput,
	gasAccumulated uint64,
	destinationErr error,
) (*vmcommon.VMOutput, bool, error) {
	callbackInput, err := context.createCallbackInput(asyncCall, destinationVMOutput, gasAccumulated, destinationErr)
	if err != nil {
		logAsync.Trace("executeSyncCallback", "error", err)
		return nil, true, err
	}

	logAsync.Trace("executeSyncCallback",
		"caller", callbackInput.CallerAddr,
		"dest", callbackInput.RecipientAddr,
		"func", callbackInput.Function,
		"args", callbackInput.Arguments,
		"gasProvided", callbackInput.GasProvided,
		"gasLocked", callbackInput.GasLocked)

	context.host.Metering().RestoreGas(asyncCall.GasLocked)
	callbackVMOutput, isComplete, callbackErr := context.host.ExecuteOnDestContext(callbackInput)
	if callbackVMOutput != nil {
		logAsync.Trace("async call: sync callback call",
			"retCode", callbackVMOutput.ReturnCode,
			"message", callbackVMOutput.ReturnMessage,
			"data", callbackVMOutput.ReturnData,
			"gasRemaining", callbackVMOutput.GasRemaining,
			"error", callbackErr)
	}

	return callbackVMOutput, isComplete, callbackErr
}

func (context *asyncContext) executeESDTTransferOnCallback(asyncCall *arwen.AsyncCall) {
	context.host.Output().PrependFinish(asyncCall.Data)

	// The contract has already paid the gas for GasLimit and
	// GasLocked, as if the call were destined to another contract. Both
	// GasLimit and GasLocked are restored in the case of
	// ESDTTransferOnCallBack because:
	// * GasLocked isn't needed, since no callback will be called
	// * GasLimit cannot be paid here, because it's the *destination*
	// contract that ends up paying the gas for the ESDTTransfer
	context.host.Metering().RestoreGas(asyncCall.GasLimit)
	context.host.Metering().RestoreGas(asyncCall.GasLocked)
	asyncCall.UpdateStatus(vmcommon.Ok)
}

// executeSyncHalfOfBuiltinFunction will synchronously call the requested
// built-in function. This is required for all cross-shard calls to built-in
// functions, because they will handle cross-shard calls themselves, by
// generating entries in vmOutput.OutputAccounts, and they need to be executed
// synchronously to do that. As a consequence, it is not necessary to call
// sendAsyncCallCrossShard(). The vmOutput produced by the built-in function,
// containing the cross-shard call, has ALREADY been merged into the main
// output by the inner call to host.ExecuteOnDestContext(). Moreover, the
// status of the AsyncCall is not updated here - it will be updated by
// PostprocessCrossShardCallback(), when the cross-shard call returns.
func (context *asyncContext) executeSyncHalfOfBuiltinFunction(asyncCall *arwen.AsyncCall) error {
	destinationCallInput, err := context.createContractCallInput(asyncCall)
	if err != nil {
		return err
	}

	// Briefly restore the AsyncCall GasLimit, after it was consumed in its
	// entirety by addAsyncCall(); this is required, because ExecuteOnDestContext()
	// must also consume the GasLimit in its entirety, before starting execution,
	// but will restore any GasRemaining to the current instance.
	metering := context.host.Metering()
	metering.RestoreGas(asyncCall.GetGasLimit())

	vmOutput, _, err := context.host.ExecuteOnDestContext(destinationCallInput)
	if err != nil {
		return err
	}

	// If the in-shard half of the built-in function call has failed, go no
	// further and execute the error callback of this AsyncCall.
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.Reject()
		callbackVMOutput, _, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, 0, err)
		context.finishAsyncLocalCallbackExecution(callbackVMOutput, callbackErr, 0)
	}

	// The gas that remains after executing the in-shard half of the built-in
	// function is provided to the cross-shard half.
	asyncCall.GasLimit = vmOutput.GasRemaining

	return nil
}

func (context *asyncContext) finishAsyncLocalCallbackExecution(
	vmOutput *vmcommon.VMOutput,
	err error,
	destinationReturnCode vmcommon.ReturnCode) {
	// output := context.host.Output()
	// if err == nil {
	// 	if setReturnCode {
	// 		output.SetReturnCode(destinationReturnCode)
	// 	}
	// 	return
	// }

	runtime := context.host.Runtime()

	runtime.GetVMInput().GasProvided = 0

	// if vmOutput == nil {
	// 	vmOutput = output.CreateVMOutputInCaseOfError(err)
	// }

	// if setReturnCode {
	// 	if vmOutput.ReturnCode != vmcommon.Ok {
	// 		output.SetReturnCode(vmOutput.ReturnCode)
	// 	} else {
	// 		output.SetReturnCode(destinationReturnCode)
	// 	}
	// }

	// output.SetReturnMessage(vmOutput.ReturnMessage)
	// output.Finish([]byte(vmOutput.ReturnCode.String()))
	// output.Finish(runtime.GetCurrentTxHash())
}

func (context *asyncContext) createContractCallInput(asyncCall *arwen.AsyncCall) (*vmcommon.ContractCallInput, error) {
	host := context.host
	runtime := host.Runtime()
	sender := runtime.GetContextAddress()

	function, arguments, err := context.callArgsParser.ParseData(string(asyncCall.GetData()))
	if err != nil {
		return nil, err
	}

	gasLimit := asyncCall.GetGasLimit()
	gasToUse := host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0).SetBytes(asyncCall.GetValue()),
			CallType:       vm.AsynchronousCall,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			GasLocked:      asyncCall.GetGasLocked(),
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: asyncCall.GetDestination(),
		Function:      function,
	}
	context.SetAsyncArgumentsForCall(contractCallInput)
	asyncCall.CallID = contractCallInput.AsyncArguments.CallID

	return contractCallInput, nil
}

// TODO function too large; refactor needed
func (context *asyncContext) createCallbackInput(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	gasAccumulated uint64,
	destinationErr error,
) (*vmcommon.ContractCallInput, error) {
	runtime := context.host.Runtime()

	actualCallbackInitiator := context.determineDestinationForAsyncCall(asyncCall.GetDestination(), asyncCall.GetData())

	arguments := context.getArgumentsForCallback(vmOutput, destinationErr)

	esdtFunction := ""
	isESDTOnCallBack := false
	esdtArgs := make([][]byte, 0)
	returnWithError := false
	if destinationErr == nil && vmOutput.ReturnCode == vmcommon.Ok {
		// when execution went Ok, callBack arguments are:
		// [0, result1, result2, ....]
		isESDTOnCallBack, esdtFunction, esdtArgs = context.isESDTTransferOnReturnDataWithNoAdditionalData(
			actualCallbackInitiator,
			runtime.GetContextAddress(),
			vmOutput)
	} else {
		returnWithError = true
	}

	callbackFunction := asyncCall.GetCallbackName()

	dataLength := computeDataLengthFromArguments(callbackFunction, arguments)
	gasLimit, err := context.computeGasLimitForCallback(asyncCall, vmOutput, dataLength)
	if err != nil {
		return nil, err
	}

	// Return to the sender SC, calling its specified callback method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:           actualCallbackInitiator,
			Arguments:            arguments,
			CallValue:            context.computeCallValueFromVMOutput(vmOutput),
			CallType:             vm.AsynchronousCallBack,
			GasPrice:             runtime.GetVMInput().GasPrice,
			GasProvided:          gasLimit,
			GasLocked:            0,
			CurrentTxHash:        runtime.GetCurrentTxHash(),
			OriginalTxHash:       runtime.GetOriginalTxHash(),
			PrevTxHash:           runtime.GetPrevTxHash(),
			ReturnCallAfterError: returnWithError,
		},
		RecipientAddr: context.address,
		Function:      callbackFunction,
	}
	context.SetAsyncArgumentsForCallback(contractCallInput, asyncCall, gasAccumulated)

	if isESDTOnCallBack {
		context.updateContractInputForESDTOnCallback(contractCallInput, esdtFunction, esdtArgs, vmOutput, asyncCall, gasAccumulated)
	}

	return contractCallInput, nil
}

func (context *asyncContext) updateContractInputForESDTOnCallback(
	contractCallInput *vmcommon.ContractCallInput,
	esdtFunction string,
	esdtArgs [][]byte,
	vmOutput *vmcommon.VMOutput,
	asyncCall *arwen.AsyncCall,
	gasAccumulated uint64) {

	oldArgLen := len(contractCallInput.Arguments)
	oldFunction := contractCallInput.Function

	contractCallInput.Function = esdtFunction
	contractCallInput.Arguments = make([][]byte, 0, oldArgLen)
	contractCallInput.Arguments = append(contractCallInput.Arguments, esdtArgs...)
	contractCallInput.Arguments = append(contractCallInput.Arguments, []byte(oldFunction))
	contractCallInput.Arguments = append(contractCallInput.Arguments, ReturnCodeToBytes(vmOutput.ReturnCode))

	if len(vmOutput.ReturnData) > 1 {
		contractCallInput.Arguments = append(contractCallInput.Arguments, vmOutput.ReturnData[1:]...)
	}
	if context.isSameShardNFTTransfer(contractCallInput) {
		contractCallInput.RecipientAddr = contractCallInput.CallerAddr
	}
	context.SetAsyncArgumentsForCallback(contractCallInput, asyncCall, gasAccumulated)

	context.host.Output().DeleteFirstReturnData()
}

func ReturnCodeToBytes(returnCode vmcommon.ReturnCode) []byte {
	if returnCode == vmcommon.Ok {
		return []byte{0}
	}
	return big.NewInt(int64(returnCode)).Bytes()
}

func (context *asyncContext) computeGasLimitForCallback(asyncCall *arwen.AsyncCall, vmOutput *vmcommon.VMOutput, dataLength int) (uint64, error) {
	metering := context.host.Metering()
	gasLimit := math.AddUint64(vmOutput.GasRemaining, asyncCall.GetGasLocked())

	gasToUse := metering.GasSchedule().ElrondAPICost.AsyncCallStep
	copyPerByte := metering.GasSchedule().BaseOperationCost.DataCopyPerByte
	gas := math.MulUint64(copyPerByte, uint64(dataLength))
	gasToUse = math.AddUint64(gasToUse, gas)
	if gasLimit <= gasToUse {
		return 0, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	return gasLimit, nil
}

func (context *asyncContext) getArgumentsForCallback(vmOutput *vmcommon.VMOutput, err error) [][]byte {
	// always provide return code as the first argument to callback function
	arguments := [][]byte{
		ReturnCodeToBytes(vmOutput.ReturnCode),
	}
	if err == nil && vmOutput.ReturnCode == vmcommon.Ok {
		// when execution went Ok, callBack arguments are:
		// [0, result1, result2, ....]
		arguments = append(arguments, vmOutput.ReturnData...)
	} else {
		// when execution returned error, callBack arguments are:
		// [error code, error message]
		arguments = append(arguments, []byte(vmOutput.ReturnMessage))
	}

	return arguments
}

func (context *asyncContext) isSameShardNFTTransfer(contractCallInput *vmcommon.ContractCallInput) bool {
	if !context.host.AreInSameShard(contractCallInput.CallerAddr, contractCallInput.RecipientAddr) {
		return false
	}

	return contractCallInput.Function == core.BuiltInFunctionMultiESDTNFTTransfer ||
		contractCallInput.Function == core.BuiltInFunctionESDTNFTTransfer
}

func (context *asyncContext) createGroupCallbackInput(group *arwen.AsyncCallGroup) *vmcommon.ContractCallInput {
	runtime := context.host.Runtime()

	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallType:       vm.AsynchronousCallBack,
			CallerAddr:     context.callerAddr,
			Arguments:      [][]byte{group.CallbackData},
			CallValue:      big.NewInt(0),
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    group.GasLocked + context.gasAccumulated,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetContextAddress(),
		Function:      group.Callback,
	}

	logAsync.Trace("created group callback input", "group", group.Identifier, "function", input.Function)
	logAsync.Trace("created group callback input gas", "provided", input.GasProvided, "locked", group.GasLocked, "accumulated", context.gasAccumulated)
	return input
}

func (context *asyncContext) createContextCallbackInput() *vmcommon.ContractCallInput {
	host := context.host
	runtime := host.Runtime()

	arguments := [][]byte{context.callbackData}

	// TODO ensure a new value for VMInput.CurrentTxHash
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     context.callerAddr,
			Arguments:      arguments,
			CallValue:      runtime.GetVMInput().CallValue,
			CallType:       vm.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    context.gasAccumulated,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetContextAddress(),
		Function:      context.callback,
	}

	logAsync.Trace("created context callback input", "sc", runtime.GetContextAddress(), "function", input.Function)
	logAsync.Trace("created context callback input gas", "provided", input.GasProvided, "accumulated", context.gasAccumulated)
	return input
}

func (context *asyncContext) isESDTTransferOnReturnDataWithNoAdditionalData(
	sndAddr, dstAddr []byte,
	destinationVMOutput *vmcommon.VMOutput,
) (bool, string, [][]byte) {
	if len(destinationVMOutput.ReturnData) == 0 {
		return false, "", nil
	}

	functionName, args, err := context.callArgsParser.ParseData(string(destinationVMOutput.ReturnData[0]))
	if err != nil {
		return false, "", nil
	}

	return context.isESDTTransferOnReturnDataFromFunctionAndArgs(sndAddr, dstAddr, functionName, args)
}

func (context *asyncContext) isESDTTransferOnReturnDataFromFunctionAndArgs(
	sndAddr, dstAddr []byte,
	functionName string,
	args [][]byte,
) (bool, string, [][]byte) {
	parsedTransfer, err := context.esdtTransferParser.ParseESDTTransfers(sndAddr, dstAddr, functionName, args)
	if err != nil {
		return false, functionName, args
	}

	isNoCallAfter := len(parsedTransfer.CallFunction) == 0
	return isNoCallAfter, functionName, args
}

func (context *asyncContext) computeCallValueFromVMOutput(destinationVMOutput *vmcommon.VMOutput) *big.Int {
	if len(destinationVMOutput.ReturnData) > 0 {
		return big.NewInt(0)
	}

	returnTransfer := big.NewInt(0)
	callBackReceiver := context.host.Runtime().GetContextAddress()
	outAcc, ok := destinationVMOutput.OutputAccounts[string(callBackReceiver)]
	if !ok {
		return returnTransfer
	}

	if len(outAcc.OutputTransfers) == 0 {
		return returnTransfer
	}

	lastOutTransfer := outAcc.OutputTransfers[len(outAcc.OutputTransfers)-1]
	if len(lastOutTransfer.Data) == 0 {
		returnTransfer.Set(lastOutTransfer.Value)
	}

	return returnTransfer
}
