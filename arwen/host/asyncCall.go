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

	asyncCallInfo := runtime.GetAsyncCallInfo()

	execMode, err := host.determineAsyncCallExecutionMode(asyncCallInfo)
	if err != nil {
		return err
	}

	if execMode == arwen.AsyncUnknown {
		return host.sendAsyncCallToDestination(asyncCallInfo)
	}

	// Cross-shard calls for built-in functions must be executed in both the
	// sender and destination shards.
	if execMode == arwen.AsyncBuiltinFunc {
		_, err := host.executeSyncDestinationCall(asyncCallInfo)
		return err
	}

	// Start calling the destination SC, synchronously.
	destinationVMOutput, destinationErr := host.executeSyncDestinationCall(asyncCallInfo)

	callbackVMOutput, callBackErr := host.executeSyncCallbackCall(asyncCallInfo, destinationVMOutput, destinationErr)

	err = host.processCallbackVMOutput(callbackVMOutput, callBackErr)
	if err != nil {
		return err
	}

	return nil
}

func (host *vmHost) determineAsyncCallExecutionMode(asyncCallInfo *arwen.AsyncCallInfo) (arwen.AsyncCallExecutionMode, error) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()

	// If ArgParser cannot read the Data field, then this is neither a SC call,
	// nor a built-in function call.
	argParser := parsers.NewCallArgsParser()
	functionName, _, err := argParser.ParseData(string(asyncCallInfo.Data))
	if err != nil {
		return arwen.AsyncUnknown, err
	}

	code, err := blockchain.GetCode(asyncCallInfo.Destination)
	if len(code) > 0 && err == nil {
		return arwen.SyncCall, nil
	}

	shardOfSC := blockchain.GetShardOfAddress(runtime.GetSCAddress())
	shardOfDest := blockchain.GetShardOfAddress(asyncCallInfo.Destination)
	sameShard := shardOfSC == shardOfDest

	if host.IsBuiltinFunctionName(functionName) {
		if sameShard {
			return arwen.SyncCall, nil
		}
		return arwen.AsyncBuiltinFunc, nil
	}

	return arwen.AsyncUnknown, nil
}

func (host *vmHost) executeSyncDestinationCall(asyncCallInfo *arwen.AsyncCallInfo) (*vmcommon.VMOutput, error) {
	destinationCallInput, err := host.createDestinationContractCallInput(asyncCallInfo)
	if err != nil {
		return nil, err
	}

	destinationVMOutput, _, err := host.ExecuteOnDestContext(destinationCallInput)
	return destinationVMOutput, err
}

func (host *vmHost) executeSyncCallbackCall(
	asyncCallInfo *arwen.AsyncCallInfo,
	destinationVMOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.VMOutput, error) {
	callbackCallInput, err := host.createCallbackContractCallInput(
		destinationVMOutput,
		asyncCallInfo.Destination,
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
	// TODO replace this function in promise-related code below.
	blockchain := host.Blockchain()
	calledSCCode, err := blockchain.GetCode(destination)

	return len(calledSCCode) > 0 && err == nil
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

func (host *vmHost) sendCallbackToCurrentCaller() error {
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
		currentCall.CallValue,
		retData,
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

	destinationAccount, _ := output.GetOutputAccount(callerAddress)
	destinationAccount.CallType = vmcommon.AsynchronousCallBack

	err := output.Transfer(
		callerAddress,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		currentCall.CallValue,
		returnData,
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

	argParser := parsers.NewCallArgsParser()
	function, arguments, err := argParser.ParseData(string(asyncCallInfo.GetData()))
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
 * processAsyncInfo takes a list of async calls and for each of them, if the code exists and can be processed on this
 *  host it will. For all others, a vm output account is generated for an actual async call.
 *  Given the fact that the generated async calls that remain pending will be saved on storage, the processing is
 *  done in two steps in order to correctly use all remaining gas. We first split the gas as specified by the developer,
 *  then we save the storage, then we split again the gas to calls that leave this shard.
 *
 * returns a list of pending calls (the ones that should be processed on other hosts)
 */
func (host *vmHost) processAsyncInfo(asyncInfo *arwen.AsyncContextInfo) (*arwen.AsyncContextInfo, error) {
	if len(asyncInfo.AsyncContextMap) == 0 {
		return asyncInfo, nil
	}

	err := host.setupAsyncCallsGas(asyncInfo)
	if err != nil {
		return nil, err
	}

	for _, asyncContext := range asyncInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if !host.canExecuteSynchronously(asyncCall.Destination, asyncCall.Data) {
				continue
			}

			procErr := host.processAsyncCall(asyncCall)
			if procErr != nil {
				return nil, procErr
			}
		}
	}

	pendingMapInfo := host.getPendingAsyncCalls(asyncInfo)
	if len(pendingMapInfo.AsyncContextMap) == 0 {
		return pendingMapInfo, nil
	}

	err = host.savePendingAsyncCalls(pendingMapInfo)
	if err != nil {
		return nil, err
	}

	err = host.setupAsyncCallsGas(pendingMapInfo)
	if err != nil {
		return nil, err
	}

	for _, asyncContext := range pendingMapInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if !host.canExecuteSynchronously(asyncCall.Destination, asyncCall.Data) {
				sendErr := host.sendAsyncCallToDestination(asyncCall)
				if sendErr != nil {
					return nil, sendErr
				}
			}
		}
	}

	return pendingMapInfo, nil
}

/**
 * processAsyncCall executes an async call and processes the callback if no extra calls are pending
 */
func (host *vmHost) processAsyncCall(asyncCall *arwen.AsyncGeneratedCall) error {
	input, _ := host.createDestinationContractCallInput(asyncCall)
	output, asyncMap, executionError := host.ExecuteOnDestContext(input)

	pendingMap := host.getPendingAsyncCalls(asyncMap)
	if len(pendingMap.AsyncContextMap) == 0 {
		return host.callbackAsync(asyncCall, output, executionError)
	}

	return executionError
}

/**
 * callbackAsync will execute a callback from an async call that was ran on this host and set it's status to resolved or rejected
 */
func (host *vmHost) callbackAsync(asyncCall *arwen.AsyncGeneratedCall, vmOutput *vmcommon.VMOutput, executionError error) error {
	asyncCall.Status = arwen.AsyncCallResolved
	callbackFunction := asyncCall.SuccessCallback
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.Status = arwen.AsyncCallRejected
		callbackFunction = asyncCall.ErrorCallback
	}

	callbackCallInput, err := host.createCallbackContractCallInput(
		vmOutput,
		asyncCall.Destination,
		callbackFunction,
		executionError,
	)
	if err != nil {
		return err
	}

	// Callback omits for now any async call - TODO: take into consideration async calls generated from callbacks
	callbackVMOutput, _, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	err = host.processCallbackVMOutput(callbackVMOutput, callBackErr)
	if err != nil {
		return err
	}

	return nil
}

/**
 * savePendingAsyncCalls takes a list of pending async calls and save them to storage so the info will be available on callback
 */
func (host *vmHost) savePendingAsyncCalls(pendingAsyncMap *arwen.AsyncContextInfo) error {
	if len(pendingAsyncMap.AsyncContextMap) == 0 {
		return nil
	}

	storage := host.Storage()
	runtime := host.Runtime()

	asyncCallStorageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	data, err := json.Marshal(pendingAsyncMap)
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
 * saveCrossShardCalls goes through the list of async calls and saves the ones that are cross shard
 */
func (host *vmHost) saveCrossShardCalls(asyncInfo *arwen.AsyncContextInfo) error {
	crossMap := &arwen.AsyncContextInfo{
		CallerAddr:      asyncInfo.CallerAddr,
		ReturnData:      asyncInfo.ReturnData,
		AsyncContextMap: make(map[string]*arwen.AsyncContext),
	}

	for contextIdentifier, asyncContext := range asyncInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if !host.canExecuteSynchronously(asyncCall.Destination, asyncCall.Data) {
				_, ok := crossMap.AsyncContextMap[contextIdentifier]
				if !ok {
					crossMap.AsyncContextMap[contextIdentifier] = &arwen.AsyncContext{
						Callback:   asyncContext.Callback,
						AsyncCalls: make([]*arwen.AsyncGeneratedCall, 0),
					}
				}
				crossMap.AsyncContextMap[contextIdentifier].AsyncCalls = append(
					crossMap.AsyncContextMap[contextIdentifier].AsyncCalls,
					asyncCall,
				)
			}
		}
	}

	return host.savePendingAsyncCalls(crossMap)
}

/**
 * getPendingAsyncCalls returns only pending async calls from a list that can also contain resolved/rejected entries
 */
func (host *vmHost) getPendingAsyncCalls(asyncInfo *arwen.AsyncContextInfo) *arwen.AsyncContextInfo {
	pendingMap := &arwen.AsyncContextInfo{
		CallerAddr:      asyncInfo.CallerAddr,
		ReturnData:      asyncInfo.ReturnData,
		AsyncContextMap: make(map[string]*arwen.AsyncContext),
	}

	for contextIdentifier, asyncContext := range asyncInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if asyncCall.Status != arwen.AsyncCallPending {
				continue
			}

			_, ok := pendingMap.AsyncContextMap[contextIdentifier]
			if !ok {
				pendingMap.AsyncContextMap[contextIdentifier] = &arwen.AsyncContext{
					Callback:   asyncContext.Callback,
					AsyncCalls: make([]*arwen.AsyncGeneratedCall, 0),
				}
			}
			pendingMap.AsyncContextMap[contextIdentifier].AsyncCalls = append(
				pendingMap.AsyncContextMap[contextIdentifier].AsyncCalls,
				asyncCall,
			)
		}
	}

	return pendingMap
}

/**
 * processCallbackStack is triggered when a callback was received from another host through a transaction.
 *  It will return an error if we receive a callback and we don't have it's associated data in the storage.
 *  If the associated callback was found in the pending set, it will be removed - It should not be executed
 *   again since it was executed in the callSCMethod step
 */
func (host *vmHost) processCallbackStack() error {
	runtime := host.Runtime()
	storage := host.Storage()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	buff := storage.GetStorageUnmetered(storageKey)
	if len(buff) == 0 {
		return nil
	}

	asyncInfo := &arwen.AsyncContextInfo{}
	err := json.Unmarshal(buff, &asyncInfo)
	if err != nil {
		return err
	}

	vmInput := runtime.GetVMInput()
	var asyncCallPosition int
	var currentContextIdentifier string
	for contextIdentifier, asyncContext := range asyncInfo.AsyncContextMap {
		for position, asyncCall := range asyncContext.AsyncCalls {
			if bytes.Equal(vmInput.CallerAddr, asyncCall.Destination) {
				asyncCallPosition = position
				currentContextIdentifier = contextIdentifier
				break
			}
		}

		if len(currentContextIdentifier) > 0 {
			break
		}
	}

	if len(currentContextIdentifier) == 0 {
		return arwen.ErrCallBackFuncNotExpected
	}

	// Remove current async call from the pending list
	currentContextCalls := asyncInfo.AsyncContextMap[currentContextIdentifier].AsyncCalls
	currentContextCalls[asyncCallPosition] = currentContextCalls[len(currentContextCalls)-1]
	currentContextCalls[len(currentContextCalls)-1] = nil
	currentContextCalls = currentContextCalls[:len(currentContextCalls)-1]

	if len(currentContextCalls) == 0 {
		// call OUR callback for resolving a full context
		delete(asyncInfo.AsyncContextMap, currentContextIdentifier)
	}

	// If we are still waiting for callbacks we return
	if len(asyncInfo.AsyncContextMap) > 0 {
		return nil
	}

	_, err = storage.SetStorage(storageKey, nil)
	if err != nil {
		return err
	}

	// Now figure out if we can execute the callback here or different shard
	if !host.canExecuteSynchronously(asyncInfo.CallerAddr, asyncInfo.ReturnData) {
		err = host.sendStorageCallbackToDestination(asyncInfo.CallerAddr, asyncInfo.ReturnData)
		if err != nil {
			return err
		}

		return nil
	}

	// The caller is in the same shard, execute it's callback
	callbackCallInput, err := host.createCallbackContractCallInput(
		host.Output().GetVMOutput(),
		asyncInfo.CallerAddr,
		arwen.CallbackDefault,
		nil,
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

/**
 * setupAsyncCallsGas sets the gasLimit for each async call with the amount of gas provided by the
 *  SC developer. The remaining gas is split between the async calls where the developer
 *  did not specify any gas amount
 */
func (host *vmHost) setupAsyncCallsGas(asyncInfo *arwen.AsyncContextInfo) error {
	gasLeft := host.Metering().GasLeft()
	gasNeeded := uint64(0)
	callsWithZeroGas := uint64(0)

	for identifier, asyncContext := range asyncInfo.AsyncContextMap {
		for index, asyncCall := range asyncContext.AsyncCalls {
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

			asyncInfo.AsyncContextMap[identifier].AsyncCalls[index].GasLimit = asyncCall.ProvidedGas
		}
	}

	if callsWithZeroGas == 0 {
		return nil
	}

	if gasLeft <= gasNeeded {
		return arwen.ErrNotEnoughGas
	}

	gasShare := (gasLeft - gasNeeded) / callsWithZeroGas
	for identifier, asyncContext := range asyncInfo.AsyncContextMap {
		for index, asyncCall := range asyncContext.AsyncCalls {
			if asyncCall.ProvidedGas == 0 {
				asyncInfo.AsyncContextMap[identifier].AsyncCalls[index].GasLimit = gasShare
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

	asyncInfo, err := host.getCurrentAsyncInfo()
	if err != nil {
		return nil, err
	}

	vmInput := runtime.GetVMInput()

	customCallback := false
	for _, asyncContext := range asyncInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if bytes.Equal(vmInput.CallerAddr, asyncCall.Destination) {
				customCallback = true
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

func (host *vmHost) getCurrentAsyncInfo() (*arwen.AsyncContextInfo, error) {
	runtime := host.Runtime()
	storage := host.Storage()

	asyncInfo := &arwen.AsyncContextInfo{}
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	buff := storage.GetStorageUnmetered(storageKey)
	if len(buff) == 0 {
		return asyncInfo, nil
	}

	err := json.Unmarshal(buff, &asyncInfo)
	if err != nil {
		return nil, err
	}

	return asyncInfo, nil
}
