package contexts

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.AsyncContext = (*asyncContext)(nil)

var logAsync = logger.GetOrCreate("arwen/async")

type asyncContext struct {
	host       arwen.VMHost
	stateStack []*asyncContext

	callerAddr         []byte
	callback           string
	callbackData       []byte
	gasPrice           uint64
	gasAccumulated     uint64
	returnData         []byte
	asyncCallGroups    []*arwen.AsyncCallGroup
	callArgsParser     arwen.CallArgsParser
	esdtTransferParser vmcommon.ESDTTransferParser

	groupCallbacksEnabled  bool
	contextCallbackEnabled bool
}

type serializableAsyncContext struct {
	CallerAddr      []byte
	Callback        string
	CallbackData    []byte
	GasPrice        uint64
	GasAccumulated  uint64
	ReturnData      []byte
	AsyncCallGroups []*arwen.AsyncCallGroup
}

// NewAsyncContext creates a new asyncContext.
func NewAsyncContext(
	host arwen.VMHost,
	callArgsParser arwen.CallArgsParser,
	esdtTransferParser vmcommon.ESDTTransferParser,
) (*asyncContext, error) {
	if check.IfNil(host) {
		return nil, arwen.ErrNilVMHost
	}
	if check.IfNil(callArgsParser) {
		return nil, arwen.ErrNilCallArgsParser
	}
	if check.IfNil(esdtTransferParser) {
		return nil, arwen.ErrNilESDTTransferParser
	}

	context := &asyncContext{
		host:                   host,
		stateStack:             nil,
		callerAddr:             nil,
		callback:               "",
		callbackData:           nil,
		gasPrice:               0,
		gasAccumulated:         0,
		returnData:             nil,
		asyncCallGroups:        make([]*arwen.AsyncCallGroup, 0),
		callArgsParser:         callArgsParser,
		esdtTransferParser:     esdtTransferParser,
		groupCallbacksEnabled:  false,
		contextCallbackEnabled: false,
	}

	return context, nil
}

// InitState initializes the internal state of the AsyncContext.
func (context *asyncContext) InitState() {
	context.callerAddr = make([]byte, 0)
	context.gasPrice = 0
	context.gasAccumulated = 0
	context.returnData = make([]byte, 0)
	context.asyncCallGroups = make([]*arwen.AsyncCallGroup, 0)
	context.callback = ""
}

// InitStateFromInput initializes the internal state of the AsyncContext with
// information provided by a ContractCallInput.
func (context *asyncContext) InitStateFromInput(input *vmcommon.ContractCallInput) {
	context.InitState()
	context.callerAddr = input.CallerAddr
	context.gasPrice = input.GasPrice
	context.gasAccumulated = 0
}

// PushState creates a deep clone of the internal state and pushes it onto the
// internal state stack.
func (context *asyncContext) PushState() {
	newState := &asyncContext{
		callerAddr:      context.callerAddr,
		callback:        context.callback,
		callbackData:    context.callbackData,
		gasPrice:        context.gasPrice,
		gasAccumulated:  context.gasAccumulated,
		returnData:      context.returnData,
		asyncCallGroups: context.asyncCallGroups,
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

// PopDiscard is a no-operation for the AsyncContext.
func (context *asyncContext) PopDiscard() {
}

// PopSetActiveState pops the state found at the top of the internal state
// stack and sets it as the 'active' state of the AsyncContext.
func (context *asyncContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.callerAddr = prevState.callerAddr
	context.callback = prevState.callback
	context.callbackData = prevState.callbackData
	context.gasPrice = prevState.gasPrice
	context.gasAccumulated = prevState.gasAccumulated
	context.returnData = prevState.returnData
	context.asyncCallGroups = prevState.asyncCallGroups
}

// PopMergeActiveState is a no-operation for the AsyncContext.
func (context *asyncContext) PopMergeActiveState() {
}

// ClearStateStack deletes all the states stored on the internal state stack.
func (context *asyncContext) ClearStateStack() {
	context.stateStack = make([]*asyncContext, 0)
}

// GetCallerAddress returns the address of the original caller.
func (context *asyncContext) GetCallerAddress() []byte {
	return context.callerAddr
}

// GetGasPrice retrieves the gas price set by the original caller.
func (context *asyncContext) GetGasPrice() uint64 {
	return context.gasPrice
}

// GetReturnData returns the data to be sent back to the original caller.
func (context *asyncContext) GetReturnData() []byte {
	return context.returnData
}

// SetReturnData sets the data to be sent back to the original caller.
func (context *asyncContext) SetReturnData(data []byte) {
	context.returnData = data
}

// GetCallGroup retrieves an AsyncCallGroup by its ID.
func (context *asyncContext) GetCallGroup(groupID string) (*arwen.AsyncCallGroup, bool) {
	index, ok := context.findGroupByID(groupID)
	if ok {
		return context.asyncCallGroups[index], true
	}
	return nil, false
}

// AddCallGroup adds the provided AsyncCallGroup to the AsyncContext, if it does not exist already.
func (context *asyncContext) AddCallGroup(group *arwen.AsyncCallGroup) error {
	_, exists := context.findGroupByID(group.Identifier)
	if exists {
		return arwen.ErrAsyncCallGroupExistsAlready
	}

	context.asyncCallGroups = append(context.asyncCallGroups, group)
	return nil
}

// SetGroupCallback registers the name of the callback method to be called upon the completion of the specified AsyncCallGroup.
func (context *asyncContext) SetGroupCallback(groupID string, callbackName string, data []byte, gas uint64) error {
	if !context.groupCallbacksEnabled {
		return arwen.ErrGroupCallbacksDisabled
	}

	group, exists := context.GetCallGroup(groupID)
	if !exists {
		return arwen.ErrAsyncCallGroupDoesNotExist
	}

	if group.IsComplete() {
		return arwen.ErrAsyncCallGroupAlreadyComplete
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

// SetContextCallback registers the name of the callback method to be called upon the completion of all the groups
func (context *asyncContext) SetContextCallback(callbackName string, data []byte, gas uint64) error {
	if !context.contextCallbackEnabled {
		return arwen.ErrContextCallbackDisabled
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

	context.gasAccumulated = gasToLock
	context.callback = callbackName
	context.callbackData = data

	return nil
}

func (context *asyncContext) deleteCallGroupByID(groupID string) {
	index, ok := context.findGroupByID(groupID)
	if !ok {
		return
	}

	context.deleteCallGroup(index)
}

func (context *asyncContext) deleteCallGroup(index int) {
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

func (context *asyncContext) isValidCallbackName(callback string) bool {
	if callback == arwen.InitFunctionName {
		return false
	}
	if context.host.IsBuiltinFunctionName(callback) {
		return false
	}

	err := context.host.Runtime().ValidateCallbackName(callback)
	if err != nil {
		return false
	}

	return true
}

func (context *asyncContext) callExists(destination []byte) bool {
	_, _, err := context.findCall(destination)
	if err != nil {
		return false
	}
	return true
}

func (context *asyncContext) findCall(destination []byte) (string, int, error) {
	for _, group := range context.asyncCallGroups {
		callIndex, ok := group.FindByDestination(destination)
		if ok {
			return group.Identifier, callIndex, nil
		}
	}

	return "", -1, arwen.ErrAsyncCallNotFound
}

// UpdateCurrentCallStatus detects the AsyncCall returning as callback,
// extracts the ReturnCode from data provided by the destination call, and updates
// the status of the AsyncCall with its value.
func (context *asyncContext) UpdateCurrentCallStatus() (*arwen.AsyncCall, error) {
	vmInput := context.host.Runtime().GetVMInput()
	if vmInput.CallType != vm.AsynchronousCallBack {
		return nil, nil
	}

	if len(vmInput.Arguments) == 0 {
		return nil, arwen.ErrCannotInterpretCallbackArgs
	}

	call, err := context.getCurrentCall()
	if err != nil {
		return nil, err
	}

	// The first argument of the callback is the return code of the destination call
	destReturnCode := big.NewInt(0).SetBytes(vmInput.Arguments[0]).Uint64()
	call.UpdateStatus(vmcommon.ReturnCode(destReturnCode))

	return call, nil
}

func (context *asyncContext) getCurrentCall() (*arwen.AsyncCall, error) {
	vmInput := context.host.Runtime().GetVMInput()
	groupID, index, err := context.findCall(vmInput.CallerAddr)
	if err != nil {
		return nil, err
	}

	group, _ := context.GetCallGroup(groupID)
	call := group.AsyncCalls[index]
	return call, nil
}

// RegisterAsyncCall validates the provided AsyncCall adds it to the specified
// group (adding the AsyncCall consumes its gas entirely).
func (context *asyncContext) RegisterAsyncCall(groupID string, call *arwen.AsyncCall) error {
	runtime := context.host.Runtime()
	metering := context.host.Metering()

	// Lock gas only if a callback is defined (either for success or for error).
	shouldLockGas := false
	if call.SuccessCallback != "" {
		err := runtime.ValidateCallbackName(call.SuccessCallback)
		if err != nil {
			return err
		}
		shouldLockGas = true
	}
	if call.ErrorCallback != "" {
		err := runtime.ValidateCallbackName(call.ErrorCallback)
		if err != nil {
			return err
		}
		shouldLockGas = true
	}

	if shouldLockGas {
		call.GasLocked = metering.ComputeGasLockedForAsync()
	}

	err := context.addAsyncCall(groupID, call)
	if err != nil {
		return err
	}

	return nil
}

// RegisterLegacyAsyncCall builds a legacy AsyncCall from provided arguments,
// computes the gas to lock depending on legacy configuration (non-dynamic gas
// locking), then adds the AsyncCall to the predefined legacy
// call group and informs Wasmer to stop contract execution with
// BreakpointAsyncCall (adding the AsyncCall consumes its gas entirely).
func (context *asyncContext) RegisterLegacyAsyncCall(address []byte, data []byte, value []byte) error {
	if !context.canRegisterLegacyAsyncCall() {
		return arwen.ErrLegacyAsyncCallInvalid
	}

	legacyGroupID := arwen.LegacyAsyncCallGroupID
	_, exists := context.GetCallGroup(legacyGroupID)
	if exists {
		return arwen.ErrOnlyOneLegacyAsyncCallAllowed
	}

	gasToLock, err := context.computeGasLockForLegacyAsyncCall()
	if err != nil {
		return err
	}

	metering := context.host.Metering()
	gasLimit := math.SubUint64(metering.GasLeft(), gasToLock)

	callbackFunction := ""
	if context.host.Runtime().HasFunction(arwen.CallbackFunctionName) {
		callbackFunction = arwen.CallbackFunctionName
	}

	err = context.addAsyncCall(legacyGroupID, &arwen.AsyncCall{
		Status:          arwen.AsyncCallPending,
		Destination:     address,
		Data:            data,
		ValueBytes:      value,
		SuccessCallback: callbackFunction,
		ErrorCallback:   callbackFunction,
		GasLimit:        gasLimit,
		GasLocked:       gasToLock,
	})
	if err != nil {
		return err
	}

	context.host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointAsyncCall)

	return nil
}

func (context *asyncContext) canRegisterLegacyAsyncCall() bool {
	vmInput := context.host.Runtime().GetVMInput()
	noGroups := len(context.asyncCallGroups) == 0
	notInCallback := vmInput.CallType != vm.AsynchronousCallBack

	return noGroups && notInCallback
}

// addAsyncCall adds the provided AsyncCall to the specified AsyncCallGroup
func (context *asyncContext) addAsyncCall(groupID string, call *arwen.AsyncCall) error {
	runtime := context.host.Runtime()
	metering := context.host.Metering()

	// TODO discuss
	// TODO add exception for the first callback instance of the same address,
	// which must be allowed to modify the AsyncContext
	scOccurrences := runtime.CountSameContractInstancesOnStack(runtime.GetSCAddress())
	callType := runtime.GetVMInput().CallType
	modifiableAsyncContext := (scOccurrences == 0) || (callType == vm.AsynchronousCallBack)
	if !modifiableAsyncContext {
		return arwen.ErrAsyncContextUnmodifiableUnlessFirstSCOrFirstCallback
	}

	err := metering.UseGasBounded(call.GasLocked)
	if err != nil {
		return err
	}
	err = metering.UseGasBounded(call.GasLimit)
	if err != nil {
		return err
	}
	execMode, err := context.determineExecutionMode(call.Destination, call.Data)
	if err != nil {
		return err
	}

	call.ExecutionMode = execMode
	group, ok := context.GetCallGroup(groupID)
	if !ok {
		group = arwen.NewAsyncCallGroup(groupID)
		err := context.AddCallGroup(group)
		if err != nil {
			return err
		}
	}
	group.AddAsyncCall(call)

	logAsync.Trace(
		"added async call",
		"group", groupID,
		"dest", string(call.Destination),
		"mode", call.ExecutionMode,
		"gas limit", call.GasLimit,
		"gas locked", call.GasLocked,
	)

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
//
// Moreover, host.ExecuteOnDestContext() will push the state stack of the
// AsyncContext and work with a clean state before calling Execute(), making
// Execute() and host.ExecuteOnDestContext() mutually reentrant.
func (context *asyncContext) Execute() error {
	if context.IsComplete() {
		logAsync.Trace("no async calls")
		return nil
	}

	metering := context.host.Metering()
	gasLeft := metering.GasLeft()
	context.accumulateGas(gasLeft)
	logAsync.Trace("async.Execute() begin", "gas left", gasLeft, "gas acc", context.gasAccumulated)
	logAsync.Trace("async.Execute() execute locals")

	// Step 1: execute all AsyncCalls that can be executed synchronously
	// (includes smart contracts and built-in functions in the same shard)
	err := context.executeAsyncLocalCalls()
	if err != nil {
		return err
	}

	logAsync.Trace("async.Execute() complete locals")

	// This call to closeCompletedAsyncCall() is necessary to remove the
	// AsyncCall that has been just before async.Execute() was called, within
	// host.callSCMethod(). This happens when a cross-shard callback returns and
	// finalizes an AsyncCall.
	context.closeCompletedAsyncCalls()
	if context.groupCallbacksEnabled {
		context.executeCompletedGroupCallbacks()
	}
	context.deleteCompletedGroups()

	logAsync.Trace("async.Execute() execute remote")
	// Step 2: in one combined step, do the following:
	// * locally execute built-in functions with cross-shard
	//   destinations, whereby the cross-shard OutputAccount entries are generated
	// * call host.sendAsyncCallCrossShard() for each pending AsyncCall, to
	//   generate the corresponding cross-shard OutputAccount entries
	// Note that all async calls below this point are pending by definition.
	for _, group := range context.asyncCallGroups {
		for _, call := range group.AsyncCalls {
			err = context.executeAsyncCall(call)
			if err != nil {
				return err
			}
		}
	}

	context.deleteCallGroupByID(arwen.LegacyAsyncCallGroupID)
	if !context.HasPendingCallGroups() && context.contextCallbackEnabled {
		logAsync.Trace("async.Execute() execute context callback")
		context.executeContextCallback()
	}

	logAsync.Trace("async.Execute() save")
	err = context.Save()
	if err != nil {
		return err
	}

	return nil
}

func (context *asyncContext) executeAsyncCall(asyncCall *arwen.AsyncCall) error {
	// Cross-shard calls to built-in functions have two halves: an intra-shard
	// half, followed by sending the call across shards.
	if asyncCall.ExecutionMode == arwen.AsyncBuiltinFuncCrossShard {
		err := context.executeSyncHalfOfBuiltinFunction(asyncCall)
		if err != nil || asyncCall.Status == arwen.AsyncCallRejected {
			return err
		}

		return nil
	}

	return context.sendAsyncCallCrossShard(asyncCall)
}

func (context *asyncContext) computeGasLockForLegacyAsyncCall() (uint64, error) {
	metering := context.host.Metering()
	err := metering.UseGasForAsyncStep()
	if err != nil {
		return 0, err
	}

	gasToLock := uint64(0)
	if context.host.Runtime().HasFunction(arwen.CallbackFunctionName) {
		gasToLock = metering.ComputeGasLockedForAsync()
	}

	return gasToLock, nil
}

// PostprocessCrossShardCallback is called by host.callSCMethod() after it
// has locally executed the callback of a returning cross-shard AsyncCall,
// which means that the AsyncContext corresponding to the original transaction
// must be loaded from storage, and then the corresponding AsyncCall must be
// deleted from the current AsyncContext.
func (context *asyncContext) PostprocessCrossShardCallback() error {
	runtime := context.host.Runtime()
	if runtime.Function() == arwen.CallbackFunctionName {
		// Legacy callbacks do not require postprocessing.
		return nil
	}

	// TODO FindAsyncCallByDestination() only returns the first matched AsyncCall
	// by destination, but there could be multiple matches in an AsyncContext.
	vmInput := runtime.GetVMInput()
	currentGroupID, asyncCallIndex, err := context.findCall(vmInput.CallerAddr)
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

	if context.groupCallbacksEnabled {
		// The current group expects no more callbacks, so its own callback can be
		// executed now.
		context.executeCallGroupCallback(currentCallGroup)
	}
	context.deleteCallGroupByID(currentGroupID)
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

	if context.contextCallbackEnabled {
		return context.executeContextCallback()
	}

	return nil
}

func (context *asyncContext) HasCallback() bool {
	return context.callback != ""
}

// HasPendingCallGroups returns true if the AsyncContext still contains AsyncCallGroup.
func (context *asyncContext) HasPendingCallGroups() bool {
	return len(context.asyncCallGroups) > 0
}

// IsComplete returns true if there are no more AsyncCallGroups contained in the AsyncContext.
func (context *asyncContext) IsComplete() bool {
	return len(context.asyncCallGroups) == 0
}

// Save serializes and saves the AsyncContext to the storage of the contract, under a protected key.
func (context *asyncContext) Save() error {
	if len(context.asyncCallGroups) == 0 {
		return nil
	}

	storage := context.host.Storage()
	runtime := context.host.Runtime()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetPrevTxHash())
	data, err := context.Serialize()
	if err != nil {
		return err
	}

	_, err = storage.SetProtectedStorage(storageKey, data)
	if err != nil {
		return err
	}

	return nil
}

// Load restores the internal state of the AsyncContext from the storage of the contract.
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

// Delete deletes the persisted state of the AsyncContext from the contract storage.
func (context *asyncContext) Delete() error {
	runtime := context.host.Runtime()
	storage := context.host.Storage()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetPrevTxHash())
	_, err := storage.SetProtectedStorage(storageKey, nil)
	return err
}

func (context *asyncContext) determineExecutionMode(destination []byte, data []byte) (arwen.AsyncCallExecutionMode, error) {
	runtime := context.host.Runtime()
	blockchain := context.host.Blockchain()

	// If ArgParser cannot read the Data field, then this is neither a SC call,
	// nor a built-in function call.
	functionName, args, err := context.callArgsParser.ParseData(string(data))
	if err != nil {
		return arwen.AsyncUnknown, err
	}

	actualDestination := context.determineDestinationForAsyncCall(destination, data)
	sameShard := context.host.AreInSameShard(runtime.GetSCAddress(), actualDestination)
	if context.host.IsBuiltinFunctionName(functionName) {
		if sameShard {
			vmInput := runtime.GetVMInput()
			isESDTTransfer, _, _ := context.isESDTTransferOnReturnDataFromFunctionAndArgs(
				runtime.GetSCAddress(),
				actualDestination,
				functionName,
				args)
			isAsyncCall := vmInput.CallType == vm.AsynchronousCall
			isReturningCall := bytes.Equal(vmInput.CallerAddr, actualDestination)

			if isESDTTransfer && isAsyncCall && isReturningCall {
				return arwen.ESDTTransferOnCallBack, nil
			}

			return arwen.AsyncBuiltinFuncIntraShard, nil
		}

		return arwen.AsyncBuiltinFuncCrossShard, nil
	}

	code, err := blockchain.GetCode(actualDestination)
	if len(code) > 0 && err == nil {
		return arwen.SyncExecution, nil
	}

	return arwen.AsyncUnknown, nil
}

func (context *asyncContext) determineDestinationForAsyncCall(destination []byte, data []byte) []byte {
	if !bytes.Equal(context.host.Runtime().GetSCAddress(), destination) {
		return destination
	}

	argsParser := context.callArgsParser
	functionName, args, err := argsParser.ParseData(string(data))
	if !context.host.IsBuiltinFunctionName(functionName) {
		return destination
	}

	parsedTransfer, err := context.esdtTransferParser.ParseESDTTransfers(destination, destination, functionName, args)
	if err != nil {
		return destination
	}

	return parsedTransfer.RcvAddr
}

func (context *asyncContext) sendAsyncCallCrossShard(asyncCall *arwen.AsyncCall) error {
	host := context.host
	runtime := host.Runtime()
	output := host.Output()

	err := output.Transfer(
		asyncCall.GetDestination(),
		runtime.GetSCAddress(),
		asyncCall.GetGasLimit(),
		asyncCall.GetGasLocked(),
		big.NewInt(0).SetBytes(asyncCall.GetValue()),
		asyncCall.GetData(),
		vm.AsynchronousCall,
	)
	if err != nil {
		return err
	}

	return nil
}

// executeContextCallback will either execute a sync call (in-shard) to
// the original caller by invoking its callback directly, or will dispatch a
// cross-shard callback to it.
func (context *asyncContext) executeContextCallback() error {
	if !context.HasCallback() {
		// TODO decide whether context.gasAccumulated should be restored here to
		// mark it as available for VMOutput.GasRemaining
		return nil
	}

	sameShard := context.host.AreInSameShard(context.host.Runtime().GetSCAddress(), context.callerAddr)
	if !sameShard {
		return context.sendContextCallbackToOriginalCaller()
	}

	// The caller is in the same shard, execute its callback
	context.executeSyncContextCallback()

	return nil
}

// TODO compare with host.sendAsyncCallbackToCaller()
func (context *asyncContext) sendContextCallbackToOriginalCaller() error {
	host := context.host
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	err := output.Transfer(
		context.callerAddr,
		runtime.GetSCAddress(),
		context.gasAccumulated,
		0,
		currentCall.CallValue,
		context.returnData,
		vm.AsynchronousCallBack,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	logAsync.Trace(
		"sendContextCallbackToOriginalCaller",
		"caller", context.callerAddr,
		"data", context.returnData,
		"gas", context.gasAccumulated)

	return nil
}

func (context *asyncContext) Serialize() ([]byte, error) {
	serializableContext := context.toSerializable()
	return json.Marshal(serializableContext)
}

func (context *asyncContext) deserialize(data []byte) (*asyncContext, error) {
	deserializedContext := &serializableAsyncContext{}
	err := json.Unmarshal(data, deserializedContext)
	if err != nil {
		return nil, err
	}

	return context.fromSerializable(deserializedContext), nil
}

func (context *asyncContext) toSerializable() *serializableAsyncContext {
	return &serializableAsyncContext{
		CallerAddr:      context.callerAddr,
		Callback:        context.callback,
		CallbackData:    context.callbackData,
		GasPrice:        context.gasPrice,
		GasAccumulated:  context.gasAccumulated,
		ReturnData:      context.returnData,
		AsyncCallGroups: context.asyncCallGroups,
	}
}

func (context *asyncContext) fromSerializable(serializedContext *serializableAsyncContext) *asyncContext {
	return &asyncContext{
		host:            nil,
		stateStack:      nil,
		callerAddr:      serializedContext.CallerAddr,
		callback:        serializedContext.Callback,
		callbackData:    serializedContext.CallbackData,
		gasPrice:        serializedContext.GasPrice,
		gasAccumulated:  serializedContext.GasAccumulated,
		returnData:      serializedContext.ReturnData,
		asyncCallGroups: serializedContext.AsyncCallGroups,
	}
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
	// form "callback@arg1hex@arg2hex...

	separator := uint64(1)
	hexSize := uint64(2)
	dataLength := uint64(len(function))
	for _, argument := range arguments {
		dataLength = math.AddUint64(dataLength, separator)
		encodedArgumentLength := math.MulUint64(uint64(len(argument)), hexSize)
		dataLength = math.AddUint64(dataLength, encodedArgumentLength)
	}

	return int(dataLength)
}

func (context *asyncContext) accumulateGas(gas uint64) {
	context.gasAccumulated = math.AddUint64(context.gasAccumulated, gas)
	logAsync.Trace("async gas accumulated", "gas", context.gasAccumulated)
}

// deleteCompletedGroups removes all completed AsyncGroups
func (context *asyncContext) deleteCompletedGroups() {
	remainingAsyncGroups := make([]*arwen.AsyncCallGroup, 0)
	for _, group := range context.asyncCallGroups {
		if !group.IsComplete() {
			remainingAsyncGroups = append(remainingAsyncGroups, group)
		} else {
			logAsync.Trace("deleted group", "group", group.Identifier)
		}
	}

	context.asyncCallGroups = remainingAsyncGroups
}

func (context *asyncContext) closeCompletedAsyncCalls() {
	for _, group := range context.asyncCallGroups {
		group.DeleteCompletedAsyncCalls()
	}
}

func (context *asyncContext) getCallByIndex(groupIndex int, callIndex int) (*arwen.AsyncCall, error) {
	if groupIndex > len(context.asyncCallGroups)-1 {
		return nil, arwen.ErrAsyncCallGroupDoesNotExist
	}
	if callIndex > len(context.asyncCallGroups[groupIndex].AsyncCalls)-1 {
		return nil, arwen.ErrAsyncCallNotFound
	}
	return context.asyncCallGroups[groupIndex].AsyncCalls[callIndex], nil
}
