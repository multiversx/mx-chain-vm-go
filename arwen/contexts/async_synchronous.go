package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (context *asyncContext) executeSynchronousCalls() error {
	for groupIndex, group := range context.asyncCallGroups {
		for _, call := range group.AsyncCalls {
			if call.ExecutionMode != arwen.SyncExecution {
				continue
			}

			err := context.executeSyncCall(call)
			if err != nil {
				return err
			}
		}

		group.DeleteCompletedAsyncCalls()

		// If all the AsyncCalls in the AsyncCallGroup were executed synchronously,
		// then the AsyncCallGroup can have its callback executed.
		if group.IsComplete() {
			context.executeCallGroupCallback(group)
			context.deleteCallGroup(groupIndex)
		}
	}
	return nil
}

func (context *asyncContext) executeSyncCall(asyncCall *arwen.AsyncCall) error {
	// Briefly restore the AsyncCall GasLimit, after it was consumed in its
	// entirety by addAsyncCall(); this is required, because ExecuteOnDestContext()
	// must also consume the GasLimit in its entirety, before starting execution,
	// but will restore any GasRemaining to the current instance.
	metering := context.host.Metering()
	metering.RestoreGas(asyncCall.GetGasLimit())

	destinationCallInput, err := context.createContractCallInput(asyncCall)
	if err != nil {
		return err
	}

	vmOutput, err := context.host.ExecuteOnDestContext(destinationCallInput)

	// The vmOutput instance returned by host.ExecuteOnDestContext() is never nil,
	// by design. Using it without checking for err is safe here.
	asyncCall.UpdateStatus(vmOutput.ReturnCode)

	callbackVMOutput, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, err)
	context.finishSyncExecution(callbackVMOutput, callbackErr)

	// TODO accumulate remaining gas from the callback into the AsyncContext,
	// after fixing the bug caught by TestExecution_ExecuteOnDestContext_GasRemaining().

	return nil
}

func (context *asyncContext) executeSyncCallback(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	err error,
) (*vmcommon.VMOutput, error) {

	callbackInput, err := context.createCallbackInput(asyncCall, vmOutput, err)
	if err != nil {
		return nil, err
	}

	return context.host.ExecuteOnDestContext(callbackInput)
}

// executeCallGroupCallback synchronously executes the designated callback of
// the AsyncCallGroup, as it was set with SetGroupCallback().
//
// Gas for the execution has been already paid for when SetGroupCallback() was
// set. The remaining gas is refunded to context.callerAddr, which initiated
// the call and paid for the gas in the first place.
func (context *asyncContext) executeCallGroupCallback(group *arwen.AsyncCallGroup) {
	if !group.HasCallback() {
		return
	}

	input := context.createGroupCallbackInput(group)
	vmOutput, err := context.host.ExecuteOnDestContext(input)
	context.finishSyncExecution(vmOutput, err)
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

	vmOutput, err := context.host.ExecuteOnDestContext(destinationCallInput)
	if err != nil {
		return err
	}

	// If the synchronous half of the built-in function call has failed, go no
	// further and execute the error callback of this AsyncCall.
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.UpdateStatus(vmOutput.ReturnCode)
		callbackVMOutput, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, err)
		context.finishSyncExecution(callbackVMOutput, callbackErr)
	}

	return nil
}

// executeSyncContextCallback will execute the callback of the original caller
// synchronously, already assuming the original caller is in the same shard
func (context *asyncContext) executeSyncContextCallback() {
	callbackCallInput := context.createContextCallbackInput()
	callbackVMOutput, callBackErr := context.host.ExecuteOnDestContext(callbackCallInput)
	context.finishSyncExecution(callbackVMOutput, callBackErr)
}

// TODO return values are never used by code that calls finishSyncExecution
func (context *asyncContext) finishSyncExecution(vmOutput *vmcommon.VMOutput, err error) {
	if err == nil {
		return
	}

	runtime := context.host.Runtime()
	output := context.host.Output()

	runtime.GetVMInput().GasProvided = 0

	if vmOutput == nil {
		vmOutput = output.CreateVMOutputInCaseOfError(err)
	}

	output.SetReturnCode(vmOutput.ReturnCode)
	output.SetReturnMessage(vmOutput.ReturnMessage)
	output.Finish([]byte(vmOutput.ReturnCode.String()))
	output.Finish(runtime.GetCurrentTxHash())
}

func (context *asyncContext) createContractCallInput(asyncCall arwen.AsyncCallHandler) (*vmcommon.ContractCallInput, error) {
	host := context.host
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()

	function, arguments, err := host.CallArgsParser().ParseData(string(asyncCall.GetData()))
	if err != nil {
		return nil, err
	}

	gasLimit := asyncCall.GetGasLimit()
	gasToUse := host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0).SetBytes(asyncCall.GetValue()),
			CallType:       vmcommon.AsynchronousCall,
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

	return contractCallInput, nil
}

func (context *asyncContext) createCallbackInput(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.ContractCallInput, error) {
	metering := context.host.Metering()
	runtime := context.host.Runtime()

	// always provide return code as the first argument to callback function
	retCodeBytes := big.NewInt(int64(vmOutput.ReturnCode)).Bytes()
	if len(retCodeBytes) == 0 {
		retCodeBytes = []byte{0}
	}
	arguments := [][]byte{retCodeBytes}

	if destinationErr == nil {
		// when execution went Ok, callBack arguments are:
		// [0, result1, result2, ....]
		arguments = append(arguments, vmOutput.ReturnData...)
	} else {
		// when execution returned error, callBack arguments are:
		// [error code, error message]
		arguments = append(arguments, []byte(vmOutput.ReturnMessage))
	}

	callbackFunction := asyncCall.GetCallbackName()

	gasLimit := vmOutput.GasRemaining + asyncCall.GetGasLocked()
	dataLength := computeDataLengthFromArguments(callbackFunction, arguments)

	gasToUse := metering.GasSchedule().ElrondAPICost.AsyncCallStep
	gas := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(dataLength)
	gasToUse = math.AddUint64(gasToUse, gas)
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	// Return to the sender SC, calling its specified callback method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     asyncCall.Destination,
			Arguments:      arguments,
			CallValue:      big.NewInt(0),
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      callbackFunction,
	}

	return contractCallInput, nil
}

func (context *asyncContext) createGroupCallbackInput(group *arwen.AsyncCallGroup) *vmcommon.ContractCallInput {
	runtime := context.host.Runtime()
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     context.callerAddr,
			Arguments:      [][]byte{group.CallbackData},
			CallValue:      big.NewInt(0),
			GasPrice:       context.gasPrice,
			GasProvided:    group.GasLocked,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      group.Callback,
	}

	return input
}

func (context *asyncContext) createContextCallbackInput() *vmcommon.ContractCallInput {
	host := context.host
	runtime := host.Runtime()

	_, arguments, err := host.CallArgsParser().ParseData(string(context.returnData))
	if err != nil {
		arguments = [][]byte{context.returnData}
	}

	// TODO ensure a new value for VMInput.CurrentTxHash
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     runtime.GetSCAddress(),
			Arguments:      arguments,
			CallValue:      runtime.GetVMInput().CallValue,
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    context.gasRemaining,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: context.callerAddr,

		// TODO Function is not actually necessary, because the original caller will decide the appropriate callback function
		Function: arwen.CallbackFunctionName,
	}
	return input
}
