package contexts

import (
	"encoding/json"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	extramath "github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ arwen.AsyncContext = (*asyncContext)(nil)

type asyncContext struct {
	host       arwen.VMHost
	stateStack []*asyncContext

	callerAddr      []byte
	gasPrice        uint64
	returnData      []byte
	asyncCallGroups []*arwen.AsyncCallGroup
}

// NewAsyncContext creates a new asyncContext
func NewAsyncContext(host arwen.VMHost) *asyncContext {
	return &asyncContext{
		host:            host,
		stateStack:      nil,
		callerAddr:      nil,
		returnData:      nil,
		asyncCallGroups: make([]*arwen.AsyncCallGroup, 0),
	}
}

func (context *asyncContext) InitState() {
	context.callerAddr = make([]byte, 0)
	context.gasPrice = 0
	context.returnData = make([]byte, 0)
	context.asyncCallGroups = make([]*arwen.AsyncCallGroup, 0)
}

func (context *asyncContext) InitStateFromInput(input *vmcommon.ContractCallInput) {
	context.InitState()
	context.callerAddr = input.CallerAddr
	context.gasPrice = input.GasPrice
}

func (context *asyncContext) PushState() {
	newState := &asyncContext{
		callerAddr:      context.callerAddr,
		gasPrice:        context.gasPrice,
		returnData:      context.returnData,
		asyncCallGroups: context.cloneCallGroups(),
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *asyncContext) cloneCallGroups() []*arwen.AsyncCallGroup {
	groupCount := len(context.asyncCallGroups)
	clonedGroups := make([]*arwen.AsyncCallGroup, groupCount)

	for i := 0; i < groupCount; i++ {
		clonedGroups[i] = context.asyncCallGroups[i].Clone()
	}

	return clonedGroups
}

func (context *asyncContext) PopDiscard() {
}

func (context *asyncContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.callerAddr = prevState.callerAddr
	context.gasPrice = prevState.gasPrice
	context.returnData = prevState.returnData
	context.asyncCallGroups = prevState.asyncCallGroups
}

func (context *asyncContext) PopMergeActiveState() {
}

func (context *asyncContext) ClearStateStack() {
	context.stateStack = make([]*asyncContext, 0)
}

func (context *asyncContext) GetCallerAddress() []byte {
	return context.callerAddr
}

func (context *asyncContext) GetGasPrice() uint64 {
	return context.gasPrice
}

func (context *asyncContext) GetReturnData() []byte {
	return context.returnData
}

func (context *asyncContext) SetReturnData(data []byte) {
	context.returnData = data
}

func (context *asyncContext) GetCallGroup(groupID string) (*arwen.AsyncCallGroup, bool) {
	index, ok := context.findGroupByID(groupID)
	if ok {
		return context.asyncCallGroups[index], true
	}
	return nil, false
}

func (context *asyncContext) addCallGroup(group *arwen.AsyncCallGroup) error {
	_, exists := context.findGroupByID(group.Identifier)
	if exists {
		return arwen.ErrAsyncCallGroupExistsAlready
	}

	context.asyncCallGroups = append(context.asyncCallGroups, group)
	return nil
}

func (context *asyncContext) SetGroupCallback(groupID string, callbackName string, data []byte, gas uint64) error {
	group, exists := context.GetCallGroup(groupID)
	if !exists {
		return arwen.ErrAsyncCallGroupDoesNotExist
	}

	err := context.host.Runtime().ValidateCallbackName(callbackName)
	if err != nil {
		return err
	}

	metering := context.host.Metering()
	gasToLock := metering.ComputeGasLockedForAsync() + gas
	err = metering.UseGasBounded(gasToLock)
	if err != nil {
		return err
	}

	group.Callback = callbackName
	group.GasLocked = gasToLock
	group.CallbackData = data

	return nil
}

func (context *asyncContext) HasPendingCallGroups() bool {
	return len(context.asyncCallGroups) > 0
}

func (context *asyncContext) IsComplete() bool {
	return len(context.asyncCallGroups) == 0
}

func (context *asyncContext) DeleteCallGroupByID(groupID string) {
	index, ok := context.findGroupByID(groupID)
	if !ok {
		return
	}

	context.DeleteCallGroup(index)
}

func (context *asyncContext) DeleteCallGroup(index int) {
	groups := context.asyncCallGroups
	if len(groups) == 0 {
		return
	}

	last := len(groups) - 1
	if index < 0 || index > last {
		return
	}

	groups[index] = groups[last]
	groups = groups[:last]
	context.asyncCallGroups = groups
}

func (context *asyncContext) AddCall(groupID string, call *arwen.AsyncCall) error {
	if context.host.IsBuiltinFunctionName(call.SuccessCallback) {
		return arwen.ErrCannotUseBuiltinAsCallback
	}
	if context.host.IsBuiltinFunctionName(call.ErrorCallback) {
		return arwen.ErrCannotUseBuiltinAsCallback
	}

	group, ok := context.GetCallGroup(groupID)
	if !ok {
		group = arwen.NewAsyncCallGroup(groupID)
		context.addCallGroup(group)
	}

	// TODO lock gas for callback

	execMode, err := context.determineExecutionMode(call.Destination, call.Data)
	if err != nil {
		return err
	}

	call.ExecutionMode = execMode
	call.Status = arwen.AsyncCallPending
	group.AddAsyncCall(call)

	return nil
}

func (context *asyncContext) isValidCallbackName(callback string) bool {
	if callback == arwen.InitFunctionName {
		return false
	}
	if context.host.IsBuiltinFunctionName(callback) {
		return false
	}

	return true
}

func (context *asyncContext) FindCall(destination []byte) (string, int, error) {
	for _, group := range context.asyncCallGroups {
		callIndex, ok := group.FindByDestination(destination)
		if ok {
			return group.Identifier, callIndex, nil
		}
	}

	return "", -1, arwen.ErrAsyncCallNotFound
}

func (context *asyncContext) UpdateCurrentCallStatus() (*arwen.AsyncCall, error) {
	vmInput := context.host.Runtime().GetVMInput()
	if vmInput.CallType != vmcommon.AsynchronousCallBack {
		return nil, nil
	}

	if len(vmInput.Arguments) == 0 {
		return nil, arwen.ErrCannotInterpretCallbackArgs
	}

	// The first argument of the callback is the return code of the destination call
	destReturnCode := big.NewInt(0).SetBytes(vmInput.Arguments[0]).Uint64()

	groupID, index, err := context.FindCall(vmInput.CallerAddr)
	if err != nil {
		return nil, err
	}

	group, _ := context.GetCallGroup(groupID)
	call := group.AsyncCalls[index]
	call.UpdateStatus(vmcommon.ReturnCode(destReturnCode))

	return call, nil
}

// PrepareLegacyAsyncCall builds an AsyncCall struct from its arguments, sets it as
// the default async call and informs Wasmer to stop contract execution with BreakpointAsyncCall
func (context *asyncContext) PrepareLegacyAsyncCall(address []byte, data []byte, value []byte) error {
	legacyGroupID := arwen.LegacyAsyncCallGroupID

	_, exists := context.GetCallGroup(legacyGroupID)
	if exists {
		return arwen.ErrOnlyOneLegacyAsyncCallAllowed
	}

	gasToLock, err := context.prepareGasForAsyncCall()
	if err != nil {
		return err
	}

	metering := context.host.Metering()
	gas := metering.GasLeft()

	err = context.AddCall(legacyGroupID, &arwen.AsyncCall{
		Status:          arwen.AsyncCallPending,
		Destination:     address,
		Data:            data,
		ValueBytes:      value,
		SuccessCallback: arwen.CallbackFunctionName,
		ErrorCallback:   arwen.CallbackFunctionName,
		ProvidedGas:     gas,
		GasLocked:       gasToLock,
	})
	if err != nil {
		return err
	}

	context.host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointAsyncCall)

	return nil
}

// Execute is the entry-point of the async calling mechanism; it is called by
// host.ExecuteOnDestContext() and host.callSCMethod(). When Execute()
// finishes, there should be no remaining AsyncCalls that can be executed
// synchronously, and all AsyncCalls that require asynchronous execution must
// already have corresponding entries in vmOutput.OutputAccounts, to be
// dispatched across shards.
//
// Execute() does NOT handle the callbacks of cross-shard AsyncCalls. See
// PostprocessCrossShardCallback() for that.
//
// Note that Execute() is mutually recursive with host.ExecuteOnDestContext(),
// because synchronous AsyncCalls are executed with
// host.ExecuteOnDestContext(), which, in turn, calls asyncContext.Execute() to
// resolve AsyncCalls generated by the AsyncCalls, and so on.

// Moreover, host.ExecuteOnDestContext() will push the state stack of the
// AsyncContext and work with a clean state before calling Execute(), making
// Execute() and host.ExecuteOnDestContext() mutually reentrant.
func (context *asyncContext) Execute() error {
	if context.IsComplete() {
		return nil
	}

	// Step 1: execute all AsyncCalls that can be executed synchronously
	// (includes smart contracts and built-in functions in the same shard)
	err := context.setupAsyncCallsGas()
	if err != nil {
		return err
	}

	err = context.executeSynchronousCalls()
	if err != nil {
		return err
	}

	// Step 2: redistribute unspent gas; then, in one combined step, do the
	// following:
	// * locally execute built-in functions with cross-shard
	//   destinations, whereby the cross-shard OutputAccount entries are generated
	// * call host.sendAsyncCallCrossShard() for each pending AsyncCall, to
	//   generate the corresponding cross-shard OutputAccount entries
	err = context.setupAsyncCallsGas()
	if err != nil {
		return err
	}

	for _, group := range context.asyncCallGroups {
		for _, call := range group.AsyncCalls {
			err = context.executeAsyncCall(call)
			if err != nil {
				return err
			}
		}
	}

	context.DeleteCallGroupByID(arwen.LegacyAsyncCallGroupID)

	err = context.Save()
	if err != nil {
		return err
	}

	return nil
}

// executeAsyncCall executes asynchronous cross-shard calls to built-in
// functions and contracts
func (context *asyncContext) executeAsyncCall(asyncCall *arwen.AsyncCall) error {
	if asyncCall.ExecutionMode == arwen.AsyncUnknown {
		return context.sendAsyncCallCrossShard(asyncCall)
	}

	if asyncCall.ExecutionMode == arwen.AsyncBuiltinFunc {
		// Built-in functions will handle cross-shard calls themselves, by
		// generating entries in vmOutput.OutputAccounts, but they need to be
		// executed synchronously to do that. It is not necessary to call
		// sendAsyncCallCrossShard(). The vmOutput produced by the built-in
		// function, containing the cross-shard call, has ALREADY been merged into
		// the main output by the inner call to host.ExecuteOnDestContext().  The
		// status of the AsyncCall is not updated here - it will be updated by
		// PostprocessCrossShardCallback(), when the cross-shard call returns.
		destinationCallInput, err := context.createSyncCallInput(asyncCall)
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

	return nil
}

func (context *asyncContext) prepareGasForAsyncCall() (uint64, error) {
	metering := context.host.Metering()
	err := metering.UseGasForAsyncStep()
	if err != nil {
		return 0, err
	}

	var shouldLockGas bool

	if !context.host.IsDynamicGasLockingEnabled() {
		// Legacy mode: static gas locking, always enabled
		shouldLockGas = true
	} else {
		// Dynamic mode: lock only if callBack() exists
		shouldLockGas = context.host.Runtime().HasCallbackMethod()
	}

	gasToLock := uint64(0)
	if shouldLockGas {
		gasToLock = metering.ComputeGasLockedForAsync()
		err = metering.UseGasBounded(gasToLock)
		if err != nil {
			return 0, err
		}
	}

	return gasToLock, nil
}

/**
 * postprocessCrossShardCallback() is called by host.callSCMethod() after it
 * has locally executed the callback of a returning cross-shard AsyncCall,
 * which means that the AsyncContext corresponding to the original transaction
 * must be loaded from storage, and then the corresponding AsyncCall must be
 * deleted from the current AsyncContext.
 */
func (context *asyncContext) PostprocessCrossShardCallback() error {
	runtime := context.host.Runtime()
	if runtime.Function() == arwen.CallbackFunctionName {
		// Legacy callbacks do not require postprocessing.
		return nil
	}

	// TODO FindAsyncCallByDestination() only returns the first matched AsyncCall
	// by destination, but there could be multiple matches in an AsyncContext.
	vmInput := runtime.GetVMInput()
	currentGroupID, asyncCallIndex, err := context.FindCall(vmInput.CallerAddr)
	if err != nil {
		return err
	}

	currentCallGroup, ok := context.GetCallGroup(currentGroupID)
	if !ok {
		return arwen.ErrCallBackFuncNotExpected
	}

	currentCallGroup.DeleteAsyncCall(asyncCallIndex)
	if currentCallGroup.HasPendingCalls() {
		return nil
	}

	// The current group expects no more callbacks, so its own callback can be
	// executed now.
	context.executeCallGroupCallback(currentCallGroup)
	context.DeleteCallGroupByID(currentGroupID)
	// Are we still waiting for callbacks to return?
	if context.HasPendingCallGroups() {
		return nil
	}

	// There are no more callbacks to return from other shards. The context can
	// be deleted from storage.
	err = context.Delete()
	if err != nil {
		return err
	}

	return context.executeContextCallback()
}

func (context *asyncContext) Save() error {
	if len(context.asyncCallGroups) == 0 {
		return nil
	}

	storage := context.host.Storage()
	runtime := context.host.Runtime()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetPrevTxHash())
	data, err := context.serialize()
	if err != nil {
		return err
	}

	_, err = storage.SetStorage(storageKey, data)
	if err != nil {
		return err
	}

	return nil
}

func (context *asyncContext) Load() error {
	runtime := context.host.Runtime()
	storage := context.host.Storage()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetPrevTxHash())
	data := storage.GetStorage(storageKey)
	if len(data) == 0 {
		return arwen.ErrNoStoredAsyncContextFound
	}

	loadedContext, err := context.deserialize(data)
	if err != nil {
		return err
	}

	context.callerAddr = loadedContext.callerAddr
	context.returnData = loadedContext.returnData
	context.asyncCallGroups = loadedContext.asyncCallGroups

	return nil
}

func (context *asyncContext) Delete() error {
	runtime := context.host.Runtime()
	storage := context.host.Storage()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetPrevTxHash())
	_, err := storage.SetStorage(storageKey, nil)
	return err
}

func (context *asyncContext) determineExecutionMode(destination []byte, data []byte) (arwen.AsyncCallExecutionMode, error) {
	runtime := context.host.Runtime()
	blockchain := context.host.Blockchain()

	// If ArgParser cannot read the Data field, then this is neither a SC call,
	// nor a built-in function call.
	functionName, _, err := context.host.CallArgsParser().ParseData(string(data))
	if err != nil {
		return arwen.AsyncUnknown, err
	}

	shardOfSC := blockchain.GetShardOfAddress(runtime.GetSCAddress())
	shardOfDest := blockchain.GetShardOfAddress(destination)
	sameShard := shardOfSC == shardOfDest

	if sameShard {
		return arwen.SyncExecution, nil
	}

	if context.host.IsBuiltinFunctionName(functionName) {
		return arwen.AsyncBuiltinFunc, nil
	}

	return arwen.AsyncUnknown, nil
}

// TODO decide whether this function is necessary at all, because unspent gas should
// be accummulated into the AsyncContext, then refunded to the original caller.
// Redistribution of gas among calls in the same group should not be necessary.
func (context *asyncContext) setupAsyncCallsGas() error {
	gasLeft := context.host.Metering().GasLeft()
	gasNeeded := uint64(0)
	callsWithZeroGas := uint64(0)

	for _, group := range context.asyncCallGroups {
		for _, asyncCall := range group.AsyncCalls {
			var err error
			gasNeeded, err = extramath.AddUint64(gasNeeded, asyncCall.ProvidedGas)
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

			asyncCall.GasLimit = asyncCall.ProvidedGas
		}
	}

	if callsWithZeroGas == 0 {
		return nil
	}

	if gasLeft <= gasNeeded {
		return arwen.ErrNotEnoughGas
	}

	gasShare := (gasLeft - gasNeeded) / callsWithZeroGas
	for _, group := range context.asyncCallGroups {
		for _, asyncCall := range group.AsyncCalls {
			if asyncCall.ProvidedGas == 0 {
				asyncCall.GasLimit = gasShare
			}
		}
	}

	return nil
}

func (context *asyncContext) sendAsyncCallCrossShard(asyncCall arwen.AsyncCallHandler) error {
	host := context.host
	runtime := host.Runtime()
	output := host.Output()

	err := output.Transfer(
		asyncCall.GetDestination(),
		runtime.GetSCAddress(),
		asyncCall.GetGasLimit(),
		asyncCall.GetGasLocked(),
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

// executeAsyncContextCallback will either execute a sync call (in-shard) to
// the original caller by invoking its callback directly, or will dispatch a
// cross-shard callback to it.
func (context *asyncContext) executeContextCallback() error {
	execMode, err := context.determineExecutionMode(context.callerAddr, context.returnData)
	if err != nil {
		return err
	}

	if execMode != arwen.SyncExecution {
		return context.sendContextCallbackToOriginalCaller()
	}

	// The caller is in the same shard, execute its callback
	callbackCallInput := context.createSyncContextCallbackInput()

	callbackVMOutput, callBackErr := context.host.ExecuteOnDestContext(callbackCallInput)
	context.finishSyncExecution(callbackVMOutput, callBackErr)

	return nil
}

func (context *asyncContext) sendContextCallbackToOriginalCaller() error {
	host := context.host
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	err := output.Transfer(
		context.callerAddr,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		0,
		currentCall.CallValue,
		context.returnData,
		vmcommon.AsynchronousCallBack,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

func (context *asyncContext) serialize() ([]byte, error) {
	serializableContext := &asyncContext{
		host:            nil,
		stateStack:      nil,
		callerAddr:      context.callerAddr,
		returnData:      context.returnData,
		asyncCallGroups: context.asyncCallGroups,
	}
	return json.Marshal(serializableContext)
}

func (context *asyncContext) deserialize(data []byte) (*asyncContext, error) {
	deserializedContext := &asyncContext{}
	err := json.Unmarshal(data, deserializedContext)
	if err != nil {
		return nil, err
	}

	return deserializedContext, nil
}

func (context *asyncContext) findGroupByID(groupID string) (int, bool) {
	for index, group := range context.asyncCallGroups {
		if group.Identifier == groupID {
			return index, true
		}
	}
	return -1, false
}

func computeDataLengthFromArguments(function string, arguments [][]byte) int {
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
