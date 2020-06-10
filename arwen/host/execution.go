package host

import (
	"bytes"
	"encoding/json"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, blockchain, _, output, runtime, storage := host.GetContexts()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)

	output.AddTxValueToAccount(address, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
	}

	vmOutput, err := host.performCodeDeploy(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}
	return vmOutput
}

func (host *vmHost) performCodeDeploy(input arwen.CodeDeployInput) (*vmcommon.VMOutput, error) {
	log.Trace("performCodeDeploy", "address", input.ContractAddress, "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	_, _, metering, output, runtime, _ := host.GetContexts()

	err := metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return nil, err
	}

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(input.ContractCode, vmInput.GasProvided)
	if err != nil {
		log.Debug("performCodeDeploy/StartWasmerInstance", "err", err)
		return nil, arwen.ErrContractInvalid
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		log.Debug("performCodeDeploy/VerifyContractCode", "err", err)
		return nil, arwen.ErrContractInvalid
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextID(idContext)

	err = host.callInitFunction()
	if err != nil {
		return nil, err
	}

	output.DeployCode(input)
	vmOutput := output.GetVMOutput()
	return vmOutput, nil
}

func (host *vmHost) doRunSmartContractUpgrade(input *vmcommon.ContractCallInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, _, _, output, runtime, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	code, codeMetadata, err := runtime.GetCodeUpgradeFromArgs()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrInvalidUpgradeArguments)
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
	}

	vmOutput, err := host.performCodeDeploy(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}
	return vmOutput
}

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput) {
	host.InitState()
	defer host.Clean()

	_, blockchain, metering, output, runtime, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	contract, err := blockchain.GetCode(runtime.GetSCAddress())
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractNotFound)
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
	}

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(contract, vmInput.GasProvided)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractInvalid)
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextID(idContext)

	err = host.callSCMethod()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	metering.UnlockGasIfAsyncStep()

	vmOutput = output.GetVMOutput()
	return
}

func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, asyncInfo *vmcommon.AsyncContextInfo, err error) {
	log.Trace("ExecuteOnDestContext", "function", input.Function)

	bigInt, _, _, output, runtime, storage := host.GetContexts()

	bigInt.PushState()
	bigInt.InitState()

	output.PushState()
	output.CensorVMOutput()

	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	storage.PushState()
	storage.SetAddress(host.Runtime().GetSCAddress())

	defer func() {
		vmOutput = host.finishExecuteOnDestContext(err)
	}()

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err = output.Transfer(input.RecipientAddr, input.CallerAddr, 0, input.CallValue, nil)
	if err != nil {
		return
	}

	err = host.execute(input)
	if err != nil {
		return
	}

	asyncInfo = runtime.GetAsyncContextInfo()
	_, err = host.processAsyncInfo(asyncInfo)
	return
}

func (host *vmHost) finishExecuteOnDestContext(executeErr error) *vmcommon.VMOutput {
	bigInt, _, _, output, runtime, storage := host.GetContexts()

	if executeErr != nil {
		// Execution failed: restore contexts as if the execution didn't happen,
		// but first create a vmOutput to capture the error.
		vmOutput := output.CreateVMOutputInCaseOfError(executeErr)

		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()
		storage.PopSetActiveState()

		return vmOutput
	}

	// Extract the VMOutput produced by the execution in isolation, before
	// restoring the contexts. This needs to be done before popping any state
	// stacks.
	vmOutput := host.Output().GetVMOutput()

	// Execution successful: restore the previous context states, except Output,
	// which will merge the current state (VMOutput) with the initial state.
	bigInt.PopSetActiveState()
	output.PopMergeActiveState()
	runtime.PopSetActiveState()
	storage.PopSetActiveState()

	return vmOutput
}

func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (asyncInfo *vmcommon.AsyncContextInfo, err error) {
	log.Trace("ExecuteOnSameContext", "function", input.Function)

	bigInt, _, _, output, runtime, _ := host.GetContexts()

	// Back up the states of the contexts (except Storage, which isn't affected
	// by ExecuteOnSameContext())
	bigInt.PushState()
	output.PushState()
	runtime.PushState()

	runtime.InitStateFromContractCallInput(input)

	defer func() {
		host.finishExecuteOnSameContext(err)
	}()

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err = output.Transfer(input.RecipientAddr, input.CallerAddr, 0, input.CallValue, nil)
	if err != nil {
		return
	}

	err = host.execute(input)
	if err != nil {
		return
	}

	asyncInfo = runtime.GetAsyncContextInfo()

	return
}

func (host *vmHost) finishExecuteOnSameContext(executeErr error) {
	bigInt, _, _, output, runtime, _ := host.GetContexts()

	if executeErr != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return
	}

	// Execution successful: discard the backups made at the beginning and
	// resume from the new state.
	bigInt.PopDiscard()
	output.PopDiscard()
	runtime.PopSetActiveState()
}

func (host *vmHost) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

func (host *vmHost) isBuiltinFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	_, ok := host.protocolBuiltinFunctions[functionName]
	return ok
}

func (host *vmHost) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	log.Trace("CreateNewContract", "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	_, blockchain, metering, output, runtime, _ := host.GetContexts()

	// Use all gas initially. In case of successful deployment, the unused gas
	// will be restored.
	initialGasProvided := input.GasProvided
	metering.UseGas(initialGasProvided)

	if runtime.ReadOnly() {
		return nil, arwen.ErrInvalidCallOnReadOnlyMode
	}

	runtime.PushState()

	runtime.SetVMInput(&input.VMInput)
	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		runtime.PopSetActiveState()
		return nil, err
	}

	err = output.Transfer(address, input.CallerAddr, 0, input.CallValue, nil)
	if err != nil {
		runtime.PopSetActiveState()
		return nil, err
	}

	blockchain.IncreaseNonce(input.CallerAddr)
	runtime.SetSCAddress(address)

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
	}

	err = metering.DeductInitialGasForIndirectDeployment(codeDeployInput)
	if err != nil {
		runtime.PopSetActiveState()
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	gasForDeployment := runtime.GetVMInput().GasProvided
	err = runtime.StartWasmerInstance(input.ContractCode, gasForDeployment)
	if err != nil {
		runtime.PopInstance()
		runtime.PopSetActiveState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		runtime.PopInstance()
		runtime.PopSetActiveState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	runtime.SetInstanceContextID(idContext)

	err = host.callInitFunction()
	if err != nil {
		runtime.PopInstance()
		runtime.PopSetActiveState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	output.DeployCode(codeDeployInput)

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	runtime.PopSetActiveState()
	arwen.RemoveHostContext(idContext)

	metering.RestoreGas(gasToRestoreToCaller)
	return address, nil
}

// TODO: Add support for indirect smart contract upgrades.
func (host *vmHost) execute(input *vmcommon.ContractCallInput) error {
	if host.isBuiltinFunctionBeingCalled() {
		return host.callBuiltinFunction(input)
	}

	// Use all gas initially, on the Wasmer instance of the caller
	// (runtime.PushInstance() is called later). In case of successful execution,
	// the unused gas will be restored.
	_, _, metering, output, runtime, _ := host.GetContexts()
	initialGasProvided := input.GasProvided
	metering.UseGas(initialGasProvided)

	if host.isInitFunctionBeingCalled() {
		return arwen.ErrInitFuncCalledInRun
	}

	contract, err := host.Blockchain().GetCode(runtime.GetSCAddress())
	if err != nil {
		return err
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	gasForExecution := runtime.GetVMInput().GasProvided
	err = runtime.StartWasmerInstance(contract, gasForExecution)
	if err != nil {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return err
	}

	runtime.SetInstanceContextID(idContext)

	err = host.callSCMethodIndirect()
	if err != nil {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return err
	}

	if output.ReturnCode() != vmcommon.Ok {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return arwen.ErrReturnCodeNotOk
	}

	metering.UnlockGasIfAsyncStep()

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	metering.RestoreGas(gasToRestoreToCaller)
	arwen.RemoveHostContext(idContext)

	return nil
}

func (host *vmHost) callSCMethodIndirect() error {
	function, err := host.Runtime().GetFunctionToCall()
	if err != nil {
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	if err != nil {
		return err
	}

	return err
}

func (host *vmHost) callBuiltinFunction(input *vmcommon.ContractCallInput) error {
	_, _, metering, output, _, _ := host.GetContexts()

	vmOutput, err := host.blockChainHook.ProcessBuiltInFunction(input)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return err
	}

	gasConsumed := input.GasProvided - vmOutput.GasRemaining
	if vmOutput.GasRemaining < input.GasProvided {
		metering.UseGas(gasConsumed)
	}

	output.AddToActiveState(vmOutput)
	return nil
}

func (host *vmHost) EthereumCallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmHost) callInitFunction() error {
	runtime := host.Runtime()
	init := runtime.GetInitFunction()
	if init == nil {
		return nil
	}

	_, err := init()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	return err
}

func (host *vmHost) callSCMethod() error {
	runtime := host.Runtime()

	err := host.verifyAllowedFunctionCall()
	if err != nil {
		return err
	}

	// If this is callback we should figure out what
	function, err := runtime.GetFunctionToCall()
	if err != nil {
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}
	if err != nil {
		return err
	}

	switch runtime.GetVMInput().CallType {
	case vmcommon.AsynchronousCall:
		pendingMap, paiErr := host.processAsyncInfo(runtime.GetAsyncContextInfo())
		if paiErr != nil {
			return err
		}
		if len(pendingMap.AsyncContextMap) == 0 {
			err = host.sendCallbackToDestination()
		}
		break
	case vmcommon.AsynchronousCallBack:
		err = host.processCallbackStack()
		break
	default:
		_, err = host.processAsyncInfo(runtime.GetAsyncContextInfo())
	}

	return err
}

func (host *vmHost) verifyAllowedFunctionCall() error {
	runtime := host.Runtime()
	functionName := runtime.Function()

	isInit := functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
	if isInit {
		return arwen.ErrInitFuncCalledInRun
	}

	isCallBack := functionName == arwen.CallBackFunctionName
	isInAsyncCallBack := runtime.GetVMInput().CallType == vmcommon.AsynchronousCallBack
	if isCallBack && !isInAsyncCallBack {
		return arwen.ErrCallBackFuncCalledInRun
	}

	return nil
}

// The first four bytes is the method selector. The rest of the input data are method arguments in chunks of 32 bytes.
// The method selector is the kecccak256 hash of the method signature.
func (host *vmHost) createETHCallInput() []byte {
	newInput := make([]byte, 0)

	function := host.Runtime().Function()
	if len(function) > 0 {
		hashOfFunction, err := host.cryptoHook.Keccak256([]byte(function))
		if err != nil {
			return nil
		}

		newInput = append(newInput, hashOfFunction[0:4]...)
	}

	for _, arg := range host.Runtime().Arguments() {
		paddedArg := make([]byte, arwen.ArgumentLenEth)
		copy(paddedArg[arwen.ArgumentLenEth-len(arg):], arg)
		newInput = append(newInput, paddedArg...)
	}

	return newInput
}

/**
 * processAsyncInfo takes a list of async calls and for each of them, if the code exists and can be processed on this
 *  host it will. For all others, a vm output account is generated for an actual async call
 *
 * returns a list of pending calls (the ones that should be processed on other hosts)
 */
func (host *vmHost) processAsyncInfo(asyncInfo *vmcommon.AsyncContextInfo) (*vmcommon.AsyncContextInfo, error) {
	pendingMap := host.getExternalAsyncCalls(asyncInfo)
	if len(pendingMap.AsyncContextMap) > 0 {
		err := host.savePendingAsyncCalls(pendingMap)
		if err != nil {
			return nil, err
		}
	}

	host.setupAsyncCallsGasByPercentages(asyncInfo)

	for _, asyncContext := range asyncInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if !host.canExecuteSynchronouslyOnDest(asyncCall.Destination) {
				err := host.sendAsyncCallToDestination(asyncCall)
				if err != nil {
					return nil, err
				}
				continue
			}

			err := host.processAsyncCall(asyncCall)
			if err != nil {
				return nil, err
			}
		}
	}

	return pendingMap, nil
}

/**
 * processAsyncCall executes an async call and processes the callback if no extra calls are pending
 */
func (host *vmHost) processAsyncCall(asyncCall *vmcommon.AsyncCall) error {
	input, _ := host.createDestinationContractCallInput(asyncCall)
	output, asyncMap, executionError := host.ExecuteOnDestContext(input)

	if executionError != nil {
		return executionError
	}

	pendingMap := host.getPendingAsyncCalls(asyncMap)
	if len(pendingMap.AsyncContextMap) == 0 {
		return host.callbackAsync(asyncCall, output, executionError)
	}

	return nil
}

/**
 * callbackAsync will execute a callback from an async call that was ran on this host and set it's status to resolved or rejected
 */
func (host *vmHost) callbackAsync(asyncCall *vmcommon.AsyncCall, vmOutput *vmcommon.VMOutput, executionError error) error {
	asyncCall.Status = vmcommon.AsyncCallResolved
	callbackFunction := asyncCall.SuccessCallback
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.Status = vmcommon.AsyncCallRejected
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
func (host *vmHost) savePendingAsyncCalls(pendingAsyncMap *vmcommon.AsyncContextInfo) error {
	storage := host.Storage()
	runtime := host.Runtime()

	asyncCallStorageKey, err := arwen.CustomStorageKey(host.Crypto(), arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	if err != nil {
		return err
	}

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
 * getPendingAsyncCalls returns only pending async calls from a list that can also contain resolved/rejected entries
 */
func (host *vmHost) getPendingAsyncCalls(asyncInfo *vmcommon.AsyncContextInfo) *vmcommon.AsyncContextInfo {
	pendingMap := &vmcommon.AsyncContextInfo{
		AsyncInitiator: vmcommon.AsyncInitiator{
			CallerAddr: asyncInfo.CallerAddr,
			ReturnData: asyncInfo.ReturnData,
		},
		AsyncContextMap: make(map[string]*vmcommon.AsyncContext),
	}

	for contextIdentifier, asyncContext := range asyncInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if asyncCall.Status != vmcommon.AsyncCallPending {
				continue
			}
			if pendingMap.AsyncContextMap[contextIdentifier] == nil {
				pendingMap.AsyncContextMap[contextIdentifier] = &vmcommon.AsyncContext{
					Callback: asyncContext.Callback,
					AsyncCalls: make([]*vmcommon.AsyncCall, 0),
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
 * getExternalAsyncCalls returns async calls for which we don't have the code to execute. This should be saved
 *  before we execute the rest of the calls so we can correctly split the gas after saving them in storage
 */
func (host *vmHost) getExternalAsyncCalls(asyncInfo *vmcommon.AsyncContextInfo) *vmcommon.AsyncContextInfo {
	pendingMap := &vmcommon.AsyncContextInfo{
		AsyncInitiator: vmcommon.AsyncInitiator{
			CallerAddr: asyncInfo.CallerAddr,
			ReturnData: asyncInfo.ReturnData,
		},
		AsyncContextMap: make(map[string]*vmcommon.AsyncContext),
	}

	for contextIdentifier, asyncContext := range asyncInfo.AsyncContextMap {
		for _, asyncCall := range asyncContext.AsyncCalls {
			if !host.canExecuteSynchronouslyOnDest(asyncCall.Destination) {
				continue
			}
			if pendingMap.AsyncContextMap[contextIdentifier] == nil {
				pendingMap.AsyncContextMap[contextIdentifier] = &vmcommon.AsyncContext{
					Callback: asyncContext.Callback,
					AsyncCalls: make([]*vmcommon.AsyncCall, 0),
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

	storageKey, err := arwen.CustomStorageKey(host.Crypto(), arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	if err != nil {
		return err
	}

	buff := storage.GetStorage(storageKey)

	asyncInfo := &vmcommon.AsyncContextInfo{}
	err = json.Unmarshal(buff, &asyncInfo)
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

	if len(asyncInfo.AsyncContextMap) == 0 {
		err = host.sendStorageCallbackToDestination(asyncInfo.AsyncInitiator)
		if err != nil {
			return err
		}

		// Delete storage, we are no longer expecting any callback
		_, err = storage.SetStorage(storageKey, nil)
		if err != nil {
			return err
		}

	}

	return nil
}

/**
 * setupAsyncCallsGasByPercentages takes the percentage of gas set up by the SC developer for each call
 *  from the gas left after the original SC call execution. If there is extra gas after divisions it
 *  is added to the last async call. There is no check here for the total of percentages to be less
 *  than 100, that check is done while the async call is added to the list
 */
func (host *vmHost) setupAsyncCallsGasByPercentages(asyncInfo *vmcommon.AsyncContextInfo) {
	gasLeft := host.Metering().GasLeft()
	gasAdded := uint64(0)

	var lastContextIdentifier string
	var lastAsyncCallIndex int
	for identifier, asyncContext := range asyncInfo.AsyncContextMap {
		lastContextIdentifier = identifier
		for index, asyncCall := range asyncContext.AsyncCalls {
			lastAsyncCallIndex = index
			gasLimit := gasLeft*(uint64(asyncCall.GasPercentage)/100)
			asyncInfo.AsyncContextMap[identifier].AsyncCalls[index].GasLimit = gasLimit
			gasAdded += gasLimit
		}
	}
	if len(lastContextIdentifier) > 0 && gasAdded < gasLeft {
		asyncInfo.AsyncContextMap[lastContextIdentifier].AsyncCalls[lastAsyncCallIndex].GasLimit += gasLeft - gasAdded
	}
}
