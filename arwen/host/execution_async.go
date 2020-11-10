package host

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	asyncCall := runtime.GetDefaultAsyncCall()
	execMode, err := host.determineAsyncCallExecutionMode(asyncCall)
	if err != nil {
		return err
	}

	if execMode == arwen.AsyncUnknown {
		return host.sendAsyncCallToDestination(asyncCall)
	}

	// Cross-shard calls for built-in functions must be executed in both the
	// sender and destination shards.
	if execMode == arwen.AsyncBuiltinFunc {
		_, err := host.executeSyncDestinationCall(asyncCall)
		return err
	}

	// Start calling the destination SC, synchronously.
	destinationVMOutput, destinationErr := host.executeSyncDestinationCall(asyncCall)

	callbackVMOutput, callBackErr := host.executeSyncCallbackCall(asyncCall, destinationVMOutput, destinationErr)

	err = host.processCallbackVMOutput(callbackVMOutput, callBackErr)
	if err != nil {
		return err
	}

	return nil
}

func (host *vmHost) determineAsyncCallExecutionMode(asyncCall *arwen.AsyncCall) (arwen.AsyncCallExecutionMode, error) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()

	// If ArgParser cannot read the Data field, then this is neither a SC call,
	// nor a built-in function call.
	argParser := parsers.NewCallArgsParser()
	functionName, _, err := argParser.ParseData(string(asyncCall.Data))
	if err != nil {
		return arwen.AsyncUnknown, err
	}

	code, err := blockchain.GetCode(asyncCall.Destination)
	if len(code) > 0 && err == nil {
		return arwen.SyncExecution, nil
	}

	shardOfSC := blockchain.GetShardOfAddress(runtime.GetSCAddress())
	shardOfDest := blockchain.GetShardOfAddress(asyncCall.Destination)
	sameShard := shardOfSC == shardOfDest

	if host.IsBuiltinFunctionName(functionName) {
		if sameShard {
			return arwen.SyncExecution, nil
		}
		return arwen.AsyncBuiltinFunc, nil
	}

	return arwen.AsyncUnknown, nil
}

func (host *vmHost) executeSyncDestinationCall(asyncCall *arwen.AsyncCall) (*vmcommon.VMOutput, error) {
	destinationCallInput, err := host.createDestinationContractCallInput(asyncCall)
	if err != nil {
		return nil, err
	}

	destinationVMOutput, _, err := host.ExecuteOnDestContext(destinationCallInput)
	return destinationVMOutput, err
}

func (host *vmHost) executeSyncCallbackCall(
	asyncCall *arwen.AsyncCall,
	destinationVMOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.VMOutput, error) {

	asyncCall.Status = arwen.AsyncCallResolved
	callbackFunction := asyncCall.SuccessCallback
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.Status = arwen.AsyncCallRejected
		callbackFunction = asyncCall.ErrorCallback
	}

	callbackCallInput, err := host.createCallbackContractCallInput(
		destinationVMOutput,
		asyncCall.Destination,
		arwen.CallbackDefault,
		destinationErr,
	)
	if err != nil {
		return nil, err
	}

	callbackVMOutput, _, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	return callbackVMOutput, callBackErr
}

func (host *vmHost) canExecuteSynchronously(destination []byte, _ []byte) bool {
	panic("remove this function")
}

func (host *vmHost) sendAsyncCallToDestination(asyncCall arwen.AsyncCallHandler) error {
	runtime := host.Runtime()
	output := host.Output()

	err := output.Transfer(
		asyncCall.GetDestination(),
		runtime.GetSCAddress(),
		asyncCall.GetGasLimit(),
		big.NewInt(0).SetBytes(asyncCall.GetValueBytes()),
		asyncCall.GetData(),
		vmcommon.AsynchronousCall,
	)
	if err != nil {
		metering := host.Metering()
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

func (host *vmHost) sendCallbackToCurrentCaller() error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	retData := []byte("@" + hex.EncodeToString([]byte(output.ReturnCode().String())))
	for _, data := range output.ReturnData() {
		retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
	}

	err := output.Transfer(
		currentCall.CallerAddr,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		currentCall.CallValue,
		retData,
		vmcommon.AsynchronousCallBack,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

func (host *vmHost) sendStorageCallbackToDestination(callerAddress, returnData []byte) error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	err := output.Transfer(
		callerAddress,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		currentCall.CallValue,
		returnData,
		vmcommon.AsynchronousCallBack,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

func (host *vmHost) createDestinationContractCallInput(asyncCall arwen.AsyncCallHandler) (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()

	argParser := parsers.NewCallArgsParser()
	function, arguments, err := argParser.ParseData(string(asyncCall.GetData()))
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
			CallValue:      big.NewInt(0).SetBytes(asyncCall.GetValueBytes()),
			CallType:       vmcommon.AsynchronousCall,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
		},
		RecipientAddr: asyncCall.GetDestination(),
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
	if destinationErr == nil && destinationVMOutput.ReturnCode == vmcommon.Ok {
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

	// Return to the sender SC, calling its specified callback method.
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

	if callbackVMOutput == nil {
		callbackVMOutput = output.CreateVMOutputInCaseOfError(callBackErr)
	}

	output.SetReturnMessage(callbackVMOutput.ReturnMessage)
	output.Finish([]byte(callbackVMOutput.ReturnCode.String()))
	output.Finish(runtime.GetCurrentTxHash())

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

/**
 * processAsyncContext takes an entire context of async calls and for each of
 * them, if the code exists and can be processed on this host it will. For all
 * others, a vm output account is generated for an actual async call.  Given
 * the fact that the generated async calls that remain pending will be saved on
 * storage, the processing is done in two steps in order to correctly use all
 * remaining gas. We first split the gas as specified by the developer, then we
 * save the storage, then we split again the gas to calls that leave this
 * shard.
 *
 * returns a list of pending calls (the ones that should be processed on other
 * hosts)
 */
func (host *vmHost) processAsyncContext(asyncContext *arwen.AsyncContext) (*arwen.AsyncContext, error) {
	if len(asyncContext.AsyncCallGroups) == 0 {
		return asyncContext, nil
	}

	err := host.setupAsyncCallsGas(asyncContext)
	if err != nil {
		return nil, err
	}

	for _, asyncCallGroup := range asyncContext.AsyncCallGroups {
		for _, asyncCall := range asyncCallGroup.AsyncCalls {
			if !host.canExecuteSynchronouslyOnDest(asyncCall.Destination, asyncCall.Data) {
				continue
			}

			procErr := host.processAsyncCall(asyncCall)
			if procErr != nil {
				return nil, procErr
			}
		}
	}

	pendingAsyncContext := asyncContext.MakeAsyncContextWithPendingCalls()
	if len(pendingAsyncContext.AsyncCallGroups) == 0 {
		return pendingAsyncContext, nil
	}

	err = host.saveAsyncContext(pendingAsyncContext)
	if err != nil {
		return nil, err
	}

	err = host.setupAsyncCallsGas(pendingAsyncContext)
	if err != nil {
		return nil, err
	}

	for _, asyncCallGroup := range pendingAsyncContext.AsyncCallGroups {
		for _, asyncCall := range asyncCallGroup.AsyncCalls {
			if !host.canExecuteSynchronouslyOnDest(asyncCall.Destination, asyncCall.Data) {
				sendErr := host.sendAsyncCallToDestination(asyncCall)
				if sendErr != nil {
					return nil, sendErr
				}
			}
		}
	}

	return pendingAsyncContext, nil
}

/**
 * processAsyncCall executes an async call and processes the callback if no extra calls are pending
 */
func (host *vmHost) processAsyncCall(asyncCall *arwen.AsyncCall) error {
	input, _ := host.createSyncCallInput(asyncCall)
	output, executionError := host.ExecuteOnDestContext(input)

	pendingMap := host.Runtime().GetAsyncContext().MakeAsyncContextWithPendingCalls()
	if len(pendingMap.AsyncCallGroups) == 0 {
		return host.callbackAsync(asyncCall, output, executionError)
	}

	return executionError
}

/**
 * callbackAsync will execute a callback from an async call that was ran on this host and set it's status to resolved or rejected
 */
func (host *vmHost) callbackAsync(asyncCall *arwen.AsyncCall, vmOutput *vmcommon.VMOutput, executionError error) error {
	asyncCall.Status = arwen.AsyncCallResolved
	callbackFunction := asyncCall.SuccessCallback
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.Status = arwen.AsyncCallRejected
		callbackFunction = asyncCall.ErrorCallback
	}

	callbackCallInput, err := host.createSyncCallbackInput(
		vmOutput,
		asyncCall.Destination,
		callbackFunction,
		executionError,
	)
	if err != nil {
		return err
	}

	// Callback omits for now any async call - TODO: take into consideration async calls generated from callbacks
	callbackVMOutput, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	err = host.finishSyncExecution(callbackVMOutput, callBackErr)
	if err != nil {
		return err
	}

	return nil
}

/**
 * saveAsyncContext takes a list of pending async calls and save them to storage so the info will be available on callback
 */
func (host *vmHost) saveAsyncContext(asyncContext *arwen.AsyncContext) error {
	if len(asyncContext.AsyncCallGroups) == 0 {
		return nil
	}

	storage := host.Storage()
	runtime := host.Runtime()

	asyncCallStorageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	data, err := json.Marshal(asyncContext)
	if err != nil {
		return err
	}

	_, err = storage.SetStorage(asyncCallStorageKey, data)
	if err != nil {
		return err
	}

	return nil
}

/**
 * postprocessCrossShardCallback() is called by host.callSCMethod() after it
 * has locally executed the callback of a returning cross-shard AsyncCall,
 * which means that the corresponding AsyncCall must be deleted from the
 * current AsyncContext.

 * Moreover, because individual AsyncCalls are contained by AsyncCallGroups, we
 * must verify whether the containing AsyncCallGroup has any remaining calls
 * pending. If not, the final callback of the containing AsyncCallGroup must be
 * executed as well.
 */
func (host *vmHost) postprocessCrossShardCallback() error {
	runtime := host.Runtime()

	asyncContext, err := host.getCurrentAsyncContext()
	if err != nil {
		return err
	}

	// TODO FindAsyncCallByDestination() only returns the first matched AsyncCall
	// by destination, but there could be multiple matches in an AsyncContext.
	vmInput := runtime.GetVMInput()
	currentGroupID, asyncCallIndex, err := asyncContext.FindAsyncCallByDestination(vmInput.CallerAddr)
	if err != nil {
		return arwen.ErrCallBackFuncNotExpected
	}

	currentCallGroup := asyncContext.AsyncCallGroups[currentGroupID]
	currentCallGroup.DeleteAsyncCall(asyncCallIndex)

	if currentCallGroup.HasPendingCalls() {
		return nil
	}

	// Are we still waiting for callbacks to return?
	if asyncContext.HasPendingCallGroups() {
		return nil
	}

	err = host.deleteCurrentAsyncContext()
	if err != nil {
		return err
	}

	// All callbacks in the current AsyncCallGroup have returned.
	// Now figure out if we can execute the callback here or different shard
	if !host.canExecuteSynchronouslyOnDest(asyncContext.CallerAddr, asyncContext.ReturnData) {
		err = host.sendStorageCallbackToDestination(asyncContext.CallerAddr, asyncContext.ReturnData)
		if err != nil {
			return err
		}

		return nil
	}

	// The caller is in the same shard, execute its callback
	callbackCallInput, err := host.createSyncCallbackInput(
		host.Output().GetVMOutput(),
		asyncContext.CallerAddr,
		arwen.CallbackDefault,
		nil,
	)
	if err != nil {
		return err
	}

	callbackVMOutput, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	err = host.finishSyncExecution(callbackVMOutput, callBackErr)
	if err != nil {
		return err
	}

	return nil
}

/**
 * setupAsyncCallsGas sets the gasLimit for each async call with the amount of gas provided by the
 *  SC developer. The remaining gas is split between the async calls where the developer
 *  did not specify any gas amount
 */
func (host *vmHost) setupAsyncCallsGas(asyncContext *arwen.AsyncContext) error {
	gasLeft := host.Metering().GasLeft()
	gasNeeded := uint64(0)
	callsWithZeroGas := uint64(0)

	for identifier, asyncCallGroup := range asyncContext.AsyncCallGroups {
		for index, asyncCall := range asyncCallGroup.AsyncCalls {
			var err error
			gasNeeded, err = math.AddUint64(gasNeeded, asyncCall.ProvidedGas)
			if err != nil {
				return err
			}

			if gasNeeded > gasLeft {
				return arwen.ErrNotEnoughGas
			}

			if asyncCall.ProvidedGas == 0 {
				callsWithZeroGas++
				continue
			}

			asyncContext.AsyncCallGroups[identifier].AsyncCalls[index].GasLimit = asyncCall.ProvidedGas
		}
	}

	if callsWithZeroGas == 0 {
		return nil
	}

	if gasLeft <= gasNeeded {
		return arwen.ErrNotEnoughGas
	}

	gasShare := (gasLeft - gasNeeded) / callsWithZeroGas
	for identifier, asyncCallGroup := range asyncContext.AsyncCallGroups {
		for index, asyncCall := range asyncCallGroup.AsyncCalls {
			if asyncCall.ProvidedGas == 0 {
				asyncContext.AsyncCallGroups[identifier].AsyncCalls[index].GasLimit = gasShare
			}
		}
	}

	return nil
}

func (host *vmHost) getFunctionByCallType(callType vmcommon.CallType) (wasmer.ExportedFunctionCallback, error) {
	runtime := host.Runtime()

	if callType != vmcommon.AsynchronousCallBack {
		return runtime.GetFunctionToCall()
	}

	asyncContext, err := host.getCurrentAsyncContext()
	if err != nil {
		return nil, err
	}

	vmInput := runtime.GetVMInput()

	customCallback := false
	for _, asyncCallGroup := range asyncContext.AsyncCallGroups {
		for _, asyncCall := range asyncCallGroup.AsyncCalls {
			if bytes.Equal(vmInput.CallerAddr, asyncCall.Destination) {
				customCallback = true
				// TODO why asyncCall.SuccessCallback? why not check for error as well,
				// and set asyncCall.ErrorCallback?
				runtime.SetCustomCallFunction(asyncCall.SuccessCallback)
				break
			}
		}

		if customCallback {
			break
		}
	}

	return runtime.GetFunctionToCall()
}

func (host *vmHost) getCurrentAsyncContext() (*arwen.AsyncContext, error) {
	runtime := host.Runtime()
	storage := host.Storage()

	asyncContext := &arwen.AsyncContext{}
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	buff := storage.GetStorage(storageKey)
	if len(buff) == 0 {
		return asyncContext, nil
	}

	err := json.Unmarshal(buff, &asyncContext)
	if err != nil {
		return nil, err
	}

	return asyncContext, nil
}

func (host *vmHost) deleteCurrentAsyncContext() error {
	runtime := host.Runtime()
	storage := host.Storage()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	_, err := storage.SetStorage(storageKey, nil)
	return err
}
