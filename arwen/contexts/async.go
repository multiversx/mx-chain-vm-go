package contexts

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.AsyncContext = (*asyncContext)(nil)

var logAsync = logger.GetOrCreate("arwen/async")

type asyncContext struct {
	host        arwen.VMHost
	stateStack  []*asyncContext
	marshalizer *marshal.GogoProtoMarshalizer

	callerAddr         []byte
	callback           string
	callbackData       []byte
	gasAccumulated     uint64
	returnData         []byte
	asyncCallGroups    []*arwen.AsyncCallGroup
	callArgsParser     arwen.CallArgsParser
	esdtTransferParser vmcommon.ESDTTransferParser

	callsCounter      uint64 // incremented and decremented during run
	totalCallsCounter uint64 // used for callid generation

	childResults           *vmcommon.VMOutput
	contextCallbackEnabled bool

	address                      []byte
	callID                       []byte
	callType                     vm.CallType
	callerCallID                 []byte
	callbackAsyncInitiatorCallID []byte

	asyncStorageDataPrefix []byte
	callbackParentCall     *arwen.AsyncCall
}

// NewAsyncContext creates a new asyncContext.
func NewAsyncContext(
	host arwen.VMHost,
	callArgsParser arwen.CallArgsParser,
	esdtTransferParser vmcommon.ESDTTransferParser,
	marshalizer *marshal.GogoProtoMarshalizer,
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

	storage := host.Storage()

	context := &asyncContext{
		host:                   host,
		stateStack:             nil,
		callerAddr:             nil,
		callback:               "",
		callbackData:           nil,
		gasAccumulated:         0,
		returnData:             nil,
		asyncCallGroups:        make([]*arwen.AsyncCallGroup, 0),
		callArgsParser:         callArgsParser,
		esdtTransferParser:     esdtTransferParser,
		contextCallbackEnabled: false,
		asyncStorageDataPrefix: storage.GetVmProtectedPrefix(arwen.AsyncDataPrefix),
		callbackParentCall:     nil,
	}

	return context, nil
}

// InitState initializes the internal state of the AsyncContext.
func (context *asyncContext) InitState() {
	context.address = nil
	context.callID = nil
	context.callerCallID = nil
	context.callerAddr = make([]byte, 0)
	context.gasAccumulated = 0
	context.returnData = make([]byte, 0)
	context.asyncCallGroups = make([]*arwen.AsyncCallGroup, 0)
	context.callback = ""
	context.callbackData = make([]byte, 0)
	context.callbackAsyncInitiatorCallID = nil
	context.callsCounter = 0
	context.totalCallsCounter = 0
	context.childResults = nil
	context.callbackParentCall = nil
}

// InitStateFromInput initializes the internal state of the AsyncContext with
// information provided by a ContractCallInput.
func (context *asyncContext) InitStateFromInput(input *vmcommon.ContractCallInput) error {
	context.InitState()
	context.callerAddr = input.CallerAddr
	context.callType = input.CallType

	runtime := context.host.Runtime()
	context.address = runtime.GetContextAddress()

	emptyStack := len(context.stateStack) == 0
	if emptyStack && !context.isCallAsync() {
		context.callID = input.CurrentTxHash
		context.callerCallID = nil
	} else {
		if input.AsyncArguments == nil {
			return vmcommon.ErrAsyncParams
		}
		context.callID = input.AsyncArguments.CallID
		context.callerCallID = input.AsyncArguments.CallerCallID
	}

	if input.CallType == vm.AsynchronousCallBack {
		context.callbackAsyncInitiatorCallID = input.AsyncArguments.CallbackAsyncInitiatorCallID
		context.gasAccumulated = input.AsyncArguments.GasAccumulated
	}

	if logAsync.GetLevel() == logger.LogTrace {
		logAsync.Trace("Calling", "function", input.Function)
		logAsync.Trace("", "address", string(context.address))
		logAsync.Trace("", "callID", context.callID)
		logAsync.Trace("", "input.GasProvided", input.GasProvided)
		logAsync.Trace("", "input.GasLocked", input.GasLocked)
		logAsync.Trace("", "callerAddr", string(context.callerAddr))
		logAsync.Trace("", "callerCallID", context.callerCallID)
		logAsync.Trace("", "callbackAsyncInitiatorCallID", context.callbackAsyncInitiatorCallID)
		logAsync.Trace("", "gasAccumulated", context.gasAccumulated)
	}

	return nil
}

// PushState creates a deep clone of the internal state and pushes it onto the
// internal state stack.
func (context *asyncContext) PushState() {
	newState := &asyncContext{
		callID:          context.callID,
		callerCallID:    context.callerCallID,
		callerAddr:      context.callerAddr,
		callback:        context.callback,
		callbackData:    context.callbackData,
		gasAccumulated:  context.gasAccumulated,
		returnData:      context.returnData,
		asyncCallGroups: context.asyncCallGroups, // TODO matei-p use cloneCallGroups()?

		callType:                     context.callType,
		callbackAsyncInitiatorCallID: context.callbackAsyncInitiatorCallID,
		address:                      context.address,
		callsCounter:                 context.callsCounter,
		totalCallsCounter:            context.totalCallsCounter,
		childResults:                 context.childResults,
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

	context.address = prevState.address
	context.callID = prevState.callID

	context.callerAddr = prevState.callerAddr
	context.callerCallID = prevState.callerCallID
	context.callType = prevState.callType
	context.callbackAsyncInitiatorCallID = prevState.callbackAsyncInitiatorCallID
	context.callback = prevState.callback
	context.callbackData = prevState.callbackData
	context.returnData = prevState.returnData
	context.asyncCallGroups = prevState.asyncCallGroups
	context.childResults = prevState.childResults
	context.callsCounter = prevState.callsCounter
	context.totalCallsCounter = prevState.totalCallsCounter
	context.gasAccumulated = math.AddUint64(context.gasAccumulated, prevState.gasAccumulated)
}

// Clone creates a clone of the given context
func (context *asyncContext) Clone() arwen.AsyncContext {
	return &asyncContext{
		address:                      context.address,
		callerAddr:                   context.callerAddr,
		callerCallID:                 context.callerCallID,
		callType:                     context.callType,
		callbackAsyncInitiatorCallID: context.callbackAsyncInitiatorCallID,
		callback:                     context.callback,
		callbackData:                 context.callbackData,
		gasAccumulated:               context.gasAccumulated,
		returnData:                   context.returnData,
		asyncCallGroups:              context.cloneCallGroups(),
		callID:                       context.callID,
		callsCounter:                 context.callsCounter,
		totalCallsCounter:            context.totalCallsCounter,
		childResults:                 context.childResults,
		host:                         context.host,
		marshalizer:                  context.marshalizer,
		callArgsParser:               context.callArgsParser,
		esdtTransferParser:           context.esdtTransferParser,
		stateStack:                   context.stateStack,
	}
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

// GetCallerCallID returns the callID of the original caller.
func (context *asyncContext) GetCallerCallID() []byte {
	return context.callerCallID
}

// GetCallbackAsyncInitiatorCallID returns the callID of the original caller.
func (context *asyncContext) GetCallbackAsyncInitiatorCallID() []byte {
	return context.callbackAsyncInitiatorCallID
}

// GetCallID is a getter for the async call's callID
func (context *asyncContext) GetCallID() []byte {
	return context.callID
}

// SetCallID is only used in integration tests.
func (context *asyncContext) SetCallID(callID []byte) {
	context.callID = callID
}

// SetCallIDForCallInGroup is only used in integration tests.
func (context *asyncContext) SetCallIDForCallInGroup(groupIndex int, callIndex int, callID []byte) {
	context.asyncCallGroups[groupIndex].AsyncCalls[callIndex].CallID = callID
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
	gasToLock := math.AddUint64(gas, metering.ComputeExtraGasLockedForAsync())
	err = metering.UseGasBounded(gasToLock)
	if err != nil {
		return err
	}

	context.gasAccumulated = gasToLock
	context.callback = callbackName
	context.callbackData = data

	return nil
}

// SetAsyncArgumentsForCall sets standard async context arguments
func (context *asyncContext) SetAsyncArgumentsForCall(input *vmcommon.ContractCallInput) {
	newCallID := context.generateNewCallID()
	context.incrementCallsCounter()
	input.VMInput.AsyncArguments = &vmcommon.AsyncArguments{
		CallID:       newCallID,
		CallerCallID: context.GetCallID(),
	}
}

// SetAsyncArgumentsForCall sets standard async context arguments
func (context *asyncContext) SetAsyncArgumentsForCallback(
	input *vmcommon.ContractCallInput,
	asyncCall *arwen.AsyncCall,
	gasAccumulated uint64) {
	newCallID := context.generateNewCallID()
	input.VMInput.AsyncArguments = &vmcommon.AsyncArguments{
		CallID:                       newCallID,
		CallerCallID:                 asyncCall.CallID,
		CallbackAsyncInitiatorCallID: context.callID,
		GasAccumulated:               gasAccumulated,
	}
}

type asyncCallLocation struct {
	asyncCall  *arwen.AsyncCall
	groupIndex int
	callIndex  int
	err        error
}

func (callInfo *asyncCallLocation) GetAsyncCall() *arwen.AsyncCall {
	return callInfo.asyncCall
}

func (callInfo *asyncCallLocation) GetGroupIndex() int {
	return callInfo.groupIndex
}

func (callInfo *asyncCallLocation) GetCallIndex() int {
	return callInfo.callIndex
}

func (callInfo *asyncCallLocation) GetError() error {
	return callInfo.err
}

// GetAsyncCallByCallID gets from the context the call with the given callID
func (context *asyncContext) GetAsyncCallByCallID(callID []byte) arwen.AsyncCallLocation {
	for groupIndex, group := range context.asyncCallGroups {
		for callIndex, callInGroup := range group.AsyncCalls {
			if bytes.Equal(callInGroup.CallID, callID) {
				return &asyncCallLocation{
					asyncCall:  callInGroup,
					groupIndex: groupIndex,
					callIndex:  callIndex,
					err:        nil,
				}
			}
		}
	}

	return &asyncCallLocation{
		asyncCall:  nil,
		groupIndex: -1,
		callIndex:  -1,
		err:        arwen.ErrAsyncCallNotFound,
	}
}

func (context *asyncContext) generateNewCallID() []byte {
	context.totalCallsCounter++
	return GenerateNewCallID(context.host.Crypto(), context.callID, big.NewInt(int64(context.totalCallsCounter)).Bytes())
}

func (context *asyncContext) incrementCallsCounter() {
	context.callsCounter++
}

func (context *asyncContext) decrementCallsCounter() {
	context.callsCounter--
}

// SetResults fills the child result of the async context
func (context *asyncContext) SetResults(vmOutput *vmcommon.VMOutput) {
	if context.host.Runtime().GetVMInput().CallType == vm.AsynchronousCall {
		context.childResults = vmOutput
	}
}

// GetGasAccumulated is a getter for gas accumulated
func (context *asyncContext) GetGasAccumulated() uint64 {
	return context.gasAccumulated
}

// IsCrossShard returns true if the current async call is cross shard
func (context *asyncContext) IsCrossShard() bool {
	return len(context.stateStack) == 0 && (context.callType == vm.AsynchronousCall || context.callType == vm.AsynchronousCallBack)
}

func (context *asyncContext) IsFirstCall() bool {
	return context.callerCallID == nil
}

func (context *asyncContext) HasCallback() bool {
	return context.callback != ""
}

// HasPendingCallGroups returns true if the AsyncContext still contains AsyncCallGroup.
func (context *asyncContext) HasPendingCallGroups() bool {
	return len(context.asyncCallGroups) > 0
}

// IsComplete returns true if the calls counter is 0 and if there are no more
// AsyncCallGroups contained in the AsyncContext.
func (context *asyncContext) IsComplete() bool {
	return context.callsCounter == 0 && len(context.asyncCallGroups) == 0
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
		call.GasLocked = math.AddUint64(call.GasLocked, metering.ComputeExtraGasLockedForAsync())
	}

	err := metering.UseGasForAsyncStep()
	if err != nil {
		return err
	}

	call.CallID = nil
	err = context.addAsyncCall(groupID, call)
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
	metering := context.host.Metering()
	logAsync.Trace("RegisterLegacyAsyncCall", "gas left", metering.GasLeft())
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

	gasLimit, err := context.computeGasLimitForLegacyAsyncCall(gasToLock)
	if err != nil {
		return err
	}

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
	metering := context.host.Metering()

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
	if execMode == arwen.ESDTTransferOnCallBack {
		context.incrementCallsCounter()
		call.CallID = context.generateNewCallID()
	}

	if context.isMultiLevelAsync(call) {
		return arwen.ErrAsyncNoMultiLevel
	}

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
	metering := context.host.Metering()
	gasLeft := metering.GasLeft()

	if context.HasPendingCallGroups() {
		logAsync.Trace("async.Execute() begin", "gas left", gasLeft, "gas acc", context.gasAccumulated)
		logAsync.Trace("async.Execute() execute locals")

		// Step 1: execute all AsyncCalls that can be executed synchronously
		// (includes smart contracts and built-in functions in the same shard)
		err := context.executeAsyncLocalCalls()
		if err != nil {
			return err
		}

		logAsync.Trace("async.Execute() execute remote")
		// Step 2: in one combined step, do the following:
		// * locally execute built-in functions with cross-shard
		//   destinations, whereby the cross-shard OutputAccount entries are generated
		// * call host.sendAsyncCallCrossShard() for each pending AsyncCall, to
		//   generate the corresponding cross-shard OutputAccount entries
		// Note that all async calls below this point are pending by definition.
		for _, group := range context.asyncCallGroups {
			for _, call := range group.AsyncCalls {
				if call.Status != arwen.AsyncCallPending {
					continue
				}
				err = context.executeAsyncCall(call)
				if err != nil {
					return err
				}
			}
		}

		context.deleteCallGroupByID(arwen.LegacyAsyncCallGroupID)
	}

	return nil
}

// UpdateCurrentAsyncCallStatus detects the AsyncCall returning as callback,
// extracts the ReturnCode from data provided by the destination call, and updates
// the status of the AsyncCall with its value.
func (context *asyncContext) UpdateCurrentAsyncCallStatus(
	address []byte,
	callID []byte,
	vmInput *vmcommon.VMInput,
) (*arwen.AsyncCall, bool, error) {
	if vmInput.CallType != vm.AsynchronousCallBack {
		return nil, false, nil
	}

	if len(vmInput.Arguments) == 0 {
		return nil, false, arwen.ErrCannotInterpretCallbackArgs
	}

	loadedContext, err := readAsyncContextFromStorage(
		context.host.Storage(),
		address,
		context.callbackAsyncInitiatorCallID,
		context.marshalizer)
	if err != nil {
		if err == arwen.ErrNoStoredAsyncContextFound {
			return getLegacyCallback(address, vmInput), true, nil
		} else {
			return nil, false, err
		}
	}

	asyncCallInfo := loadedContext.GetAsyncCallByCallID(callID)
	call := asyncCallInfo.GetAsyncCall()
	err = asyncCallInfo.GetError()
	if err != nil {
		if err == arwen.ErrAsyncCallNotFound {
			return getLegacyCallback(address, vmInput), true, nil
		} else {
			return nil, false, err
		}
	}

	// The first argument of the callback is the return code of the destination call
	destReturnCode := big.NewInt(0).SetBytes(vmInput.Arguments[0]).Uint64()
	call.UpdateStatus(vmcommon.ReturnCode(destReturnCode))

	return call, false, nil
}

func getLegacyCallback(address []byte, vmInput *vmcommon.VMInput) *arwen.AsyncCall {
	var valueBytes []byte = nil
	if vmInput.CallValue != nil {
		valueBytes = vmInput.CallValue.Bytes()
	}
	return &arwen.AsyncCall{
		Status:          arwen.AsyncCallResolved,
		Destination:     address,
		ValueBytes:      valueBytes,
		SuccessCallback: arwen.CallbackFunctionName,
		ErrorCallback:   arwen.CallbackFunctionName,
		GasLimit:        vmInput.GasProvided,
		GasLocked:       vmInput.GasLocked,
	}
}

func (context *asyncContext) isMultiLevelAsync(call *arwen.AsyncCall) bool {
	// ESDTTransferOnCallback must be allowed as an exception, even if it appears
	// to be a 2-level async call.
	return context.isCallAsyncOnStack() && call.ExecutionMode != arwen.ESDTTransferOnCallBack
}

func (context *asyncContext) isCallAsyncOnStack() bool {
	if context.isCallAsync() {
		return true
	}

	for index := len(context.stateStack) - 1; index >= 0; index-- {
		stackContext := context.stateStack[index]
		if stackContext.isCallAsync() {
			return true
		}
	}
	return false
}

func (context *asyncContext) isCallAsync() bool {
	return IsCallAsync(context.callType)
}

// IsCallAsync checks if the call is an async or callback async
func IsCallAsync(callType vm.CallType) bool {
	return callType == vm.AsynchronousCall || callType == vm.AsynchronousCallBack
}

// IsCallback checks if the call is a callback async
func IsCallback(callType vm.CallType) bool {
	return callType == vm.AsynchronousCallBack
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
		gasToLock = metering.ComputeExtraGasLockedForAsync()
	}

	return gasToLock, nil
}

func (context *asyncContext) computeGasLimitForLegacyAsyncCall(gasToLock uint64) (uint64, error) {
	gasLimit := math.SubUint64(context.host.Metering().GasLeft(), gasToLock)
	return gasLimit, nil
}

// DeleteAsyncCallAndCleanGroup deletes the specified async call and the group if this is the last call
func (context *asyncContext) DeleteAsyncCallAndCleanGroup(callID []byte) error {
	asyncCallInfo := context.GetAsyncCallByCallID(callID)
	groupIndex := asyncCallInfo.GetGroupIndex()
	callIndex := asyncCallInfo.GetCallIndex()
	err := asyncCallInfo.GetError()
	if err != nil {
		return err
	}

	currentCallGroup := context.asyncCallGroups[groupIndex]
	currentCallGroup.DeleteAsyncCall(callIndex)

	if currentCallGroup.IsComplete() {
		context.deleteCallGroup(groupIndex)
	}

	return nil
}

func (context *asyncContext) callCallback(callID []byte, vmOutput *vmcommon.VMOutput, err error) (bool, *vmcommon.VMOutput, error) {
	sender := context.address
	destination := context.callerAddr

	sameShard := context.host.AreInSameShard(sender, destination)
	if !sameShard {
		err = context.SendCrossShardCallback(vmOutput.ReturnCode, vmOutput.ReturnData, vmOutput.ReturnMessage)
		return false, nil, err
	}

	gasAccumulated := context.gasAccumulated
	loadedContext, _ := context.LoadParentContextFromStackOrStorage()
	asyncCallInfo := loadedContext.GetAsyncCallByCallID(callID)
	asyncCall := asyncCallInfo.GetAsyncCall()
	errLoad := asyncCallInfo.GetError()
	if errLoad != nil {
		return false, nil, errLoad
	}

	context.host.Metering().DisableRestoreGas()
	isComplete, callbackVMOutput := loadedContext.ExecuteSyncCallbackAndFinishOutput(asyncCall, vmOutput, nil, gasAccumulated, err)
	context.host.Metering().EnableRestoreGas()
	return isComplete, callbackVMOutput, nil
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

func (context *asyncContext) getContextFromStack(address []byte, callID []byte) *asyncContext {
	var loadedContext *asyncContext
	for _, stackContext := range context.stateStack {
		if bytes.Equal(stackContext.address, address) && bytes.Equal(stackContext.callID, callID) {
			loadedContext = stackContext
			loadedContext.host = context.host
			break
		}
	}
	return loadedContext
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
	sameShard := context.host.AreInSameShard(runtime.GetContextAddress(), actualDestination)
	if context.host.IsBuiltinFunctionName(functionName) {
		if sameShard {
			vmInput := runtime.GetVMInput()
			isESDTTransfer, _, _ := context.isESDTTransferOnReturnDataFromFunctionAndArgs(
				runtime.GetContextAddress(),
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
	if !bytes.Equal(context.host.Runtime().GetContextAddress(), destination) {
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

func (context *asyncContext) findGroupByID(groupID string) (int, bool) {
	for index, group := range context.asyncCallGroups {
		if group.Identifier == groupID {
			return index, true
		}
	}
	return -1, false
}

// computeDataLengthFromArguments salculates what length would the Data field have, were it of the
// form "callback@arg1hex@arg2hex..."
func computeDataLengthFromArguments(function string, arguments [][]byte) int {
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

// HasLegacyGroup checks if the a legacy async group was created
func (context *asyncContext) HasLegacyGroup() bool {
	_, hasLegacyGroup := context.GetCallGroup(arwen.LegacyAsyncCallGroupID)
	return hasLegacyGroup
}

// SetCallbackParentCall sets the async call that triggered the callback (used for callback closure)
func (context *asyncContext) SetCallbackParentCall(asyncCall *arwen.AsyncCall) {
	context.callbackParentCall = asyncCall
}

// GetCallbackClosure gets the async call callback closure
func (context *asyncContext) GetCallbackClosure() ([]byte, error) {
	if context.callbackParentCall == nil {
		stackContext := context.Clone()
		stackContext, err := stackContext.LoadParentContextFromStackOrStorage()
		if err != nil {
			return nil, arwen.ErrAsyncNoCallbackForClosure
		}
		context.callbackParentCall = stackContext.
			GetAsyncCallByCallID(context.callerCallID).
			GetAsyncCall()
	}
	if context.callbackParentCall == nil {
		return nil, arwen.ErrAsyncNoCallbackForClosure
	}
	return context.callbackParentCall.CallbackClosure, nil
}

// DebugCallIDAsString - just for debug purposes
func DebugCallIDAsString(arr []byte) string {
	if len(arr) > 3 {
		return "[" + string(arr)[:5] + "...]"
	}
	return fmt.Sprint(arr)
}
