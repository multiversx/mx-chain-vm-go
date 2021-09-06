package contexts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
)

var _ arwen.AsyncContext = (*asyncContext)(nil)

var logAsync = logger.GetOrCreate("arwen/async")

type asyncContext struct {
	host       arwen.VMHost
	stateStack []*asyncContext

	address  []byte
	callID   []byte
	callType vm.CallType

	callerAddr                 []byte
	callerCallID               []byte
	callAsyncIdentifierAsBytes []byte

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

	callsCounter      uint64 // incremented and decremented during run
	totalCallsCounter uint64 // used for callid generation
	childResults      *vmcommon.VMOutput
}

type serializableAsyncContext struct {
	Address  []byte
	CallID   []byte
	CallType vm.CallType

	CallerAddr                       []byte
	CallerCallID                     []byte
	CallerCallAsyncIdentifierAsBytes []byte

	Callback        string
	CallbackData    []byte
	GasPrice        uint64
	GasAccumulated  uint64
	ReturnData      []byte
	AsyncCallGroups []*arwen.AsyncCallGroup

	ChildResults      *vmcommon.VMOutput
	CallsCounter      uint64 // incremented and decremented during run
	TotalCallsCounter uint64 // used for callid generation
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
	context.callAsyncIdentifierAsBytes = nil
}

// InitStateFromInput initializes the internal state of the AsyncContext with
// information provided by a ContractCallInput.
func (context *asyncContext) InitStateFromInput(input *vmcommon.ContractCallInput) {
	context.InitState()
	context.callerAddr = input.CallerAddr
	context.gasPrice = input.GasPrice
	context.gasAccumulated = 0

	runtime := context.host.Runtime()
	context.address = runtime.GetSCAddress()

	// TODO matei-p change to debug logging
	fmt.Println("Calling function ", input.Function)
	if len(context.stateStack) == 0 && input.CallType != vm.AsynchronousCall && input.CallType != vm.AsynchronousCallBack {
		context.callID = input.CurrentTxHash
		context.callerCallID = nil
	} else {
		context.callID = runtime.GetAndEliminateFirstArgumentFromList()
		context.callerCallID = runtime.GetAndEliminateFirstArgumentFromList()
	}
	context.callType = input.CallType
	context.callsCounter = 0
	context.totalCallsCounter = 0

	// TODO matei-p change to debug logging
	fmt.Println("\taddress", string(context.address))
	fmt.Println("\tcallID", DebugCallIDAsString(context.callID))
	if input.CallType == vm.AsynchronousCall || input.CallType == vm.AsynchronousCallBack {
		context.callAsyncIdentifierAsBytes = runtime.GetAndEliminateFirstArgumentFromList()
		// TODO matei-p change to debug logging
		asynCallIdentifier, _ := context.GetCallerAsyncCallIdentifier()
		asynCallIdentifierAsString, _ := json.Marshal(asynCallIdentifier)
		fmt.Println("\tcallerAddr", string(context.callerAddr))
		fmt.Println("\tcallerCallID", DebugCallIDAsString(context.callerCallID))
		fmt.Println("\tcallerCallAsyncIdentifier", string(asynCallIdentifierAsString))
	}
	// TODO matei-p change to debug logging
	fmt.Println("\tinput.GasProvided", input.GasProvided)
	if input.GasLocked != 0 {
		fmt.Println("\tinput.GasLocked", input.GasLocked)
	}
}

// PushState creates a deep clone of the internal state and pushes it onto the
// internal state stack.
func (context *asyncContext) PushState() {
	// fmt.Println("---> PUSH " + string(context.address) + " " + DebugCallIDAsString(context.callID))
	newState := &asyncContext{
		address:                    context.address,
		callerAddr:                 context.callerAddr,
		callerCallID:               context.callerCallID,
		callType:                   context.callType,
		callAsyncIdentifierAsBytes: context.callAsyncIdentifierAsBytes,
		callback:                   context.callback,
		callbackData:               context.callbackData,
		gasPrice:                   context.gasPrice,
		gasAccumulated:             context.gasAccumulated,
		returnData:                 context.returnData,
		asyncCallGroups:            context.asyncCallGroups, // TODO matei-p use cloneCallGroups()?
		callID:                     context.callID,
		callsCounter:               context.callsCounter,
		totalCallsCounter:          context.totalCallsCounter,
		childResults:               context.childResults,
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
	// fmt.Println("---> POP " + string(prevState.address) + " " + DebugCallIDAsString(prevState.callID))
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.address = prevState.address
	context.callID = prevState.callID

	// if context.Load() != arwen.ErrNoStoredAsyncContextFound {
	// 	return
	// }

	context.callerAddr = prevState.callerAddr
	context.callerCallID = prevState.callerCallID
	context.callType = prevState.callType
	context.callAsyncIdentifierAsBytes = prevState.callAsyncIdentifierAsBytes
	context.callback = prevState.callback
	context.callbackData = prevState.callbackData
	context.gasPrice = prevState.gasPrice
	context.gasAccumulated = prevState.gasAccumulated
	context.returnData = prevState.returnData
	context.asyncCallGroups = prevState.asyncCallGroups
	context.childResults = prevState.childResults
	context.callsCounter = prevState.callsCounter
	context.totalCallsCounter = prevState.totalCallsCounter
}

func (context *asyncContext) Clone() arwen.AsyncContext {
	return &asyncContext{
		address:                    context.address,
		callerAddr:                 context.callerAddr,
		callerCallID:               context.callerCallID,
		callType:                   context.callType,
		callAsyncIdentifierAsBytes: context.callAsyncIdentifierAsBytes,
		callback:                   context.callback,
		callbackData:               context.callbackData,
		gasPrice:                   context.gasPrice,
		gasAccumulated:             context.gasAccumulated,
		returnData:                 context.returnData,
		asyncCallGroups:            context.cloneCallGroups(),
		callID:                     context.callID,
		callsCounter:               context.callsCounter,
		totalCallsCounter:          context.totalCallsCounter,
		childResults:               context.childResults,
	}
}

// PopMergeActiveState is a no-operation for the AsyncContext.
func (context *asyncContext) PopMergeActiveState() {
}

// ClearStateStack deletes all the states stored on the internal state stack.
func (context *asyncContext) ClearStateStack() {
	context.stateStack = make([]*asyncContext, 0)
}

// GetAddress returns the address of the context.
func (context *asyncContext) GetAddress() []byte {
	return context.address
}

// GetCallerAddress returns the address of the original caller.
func (context *asyncContext) GetCallerAddress() []byte {
	return context.callerAddr
}

// GetCallerCallID returns the callID of the original caller.
func (context *asyncContext) GetCallerCallID() []byte {
	return context.callerCallID
}

func (context *asyncContext) GetCallerAsyncCallIdentifier() (*arwen.AsyncCallIdentifier, error) {
	asyncCallIdentifier, err := arwen.ReadAsyncCallIdentifierFromBytes(context.callAsyncIdentifierAsBytes)
	if err != nil {
		return nil, err
	}
	return asyncCallIdentifier, nil
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

// UpdateCurrentAsyncCallStatus detects the AsyncCall returning as callback,
// extracts the ReturnCode from data provided by the destination call, and updates
// the status of the AsyncCall with its value.
func (context *asyncContext) UpdateCurrentAsyncCallStatus(address []byte, callID []byte, asyncCallIdentifier *arwen.AsyncCallIdentifier, vmInput *vmcommon.VMInput) (*arwen.AsyncCall, error) {
	deserializedContext, err := context.readSerializedAsyncContextFromStackOrStore(address, callID)
	if err != nil {
		return nil, err
	}

	if vmInput.CallType != vm.AsynchronousCallBack {
		return nil, nil
	}

	if len(vmInput.Arguments) == 0 {
		return nil, arwen.ErrCannotInterpretCallbackArgs
	}

	call, err := deserializedContext.GetCallByAsyncIdentifier(asyncCallIdentifier)
	if err != nil {
		return nil, err
	}

	// The first argument of the callback is the return code of the destination call
	destReturnCode := big.NewInt(0).SetBytes(vmInput.Arguments[0]).Uint64()
	call.UpdateStatus(vmcommon.ReturnCode(destReturnCode))

	return call, nil
}

func (context *asyncContext) GetCallByAsyncIdentifier(asyncCallIdentifier *arwen.AsyncCallIdentifier) (*arwen.AsyncCall, error) {
	return getCallByAsyncIdentifier(context.asyncCallGroups, asyncCallIdentifier)
}

func (context *serializableAsyncContext) GetCallByAsyncIdentifier(asyncCallIdentifier *arwen.AsyncCallIdentifier) (*arwen.AsyncCall, error) {
	return getCallByAsyncIdentifier(context.AsyncCallGroups, asyncCallIdentifier)
}

func getCallByAsyncIdentifier(groups []*arwen.AsyncCallGroup, asyncCallIdentifier *arwen.AsyncCallIdentifier) (*arwen.AsyncCall, error) {
	groupIndex, groupFound := findGroupByIDInAsyncCallGroups(groups, asyncCallIdentifier.GroupIdentifier)
	if !groupFound {
		return nil, arwen.ErrAsyncCallNotFound
	}
	callsInGroup := groups[groupIndex].AsyncCalls
	// we can't directly use the IndexInGroup, some calls will be deleted
	for _, callInGroup := range callsInGroup {
		if callInGroup.Identifier.IndexInGroup == asyncCallIdentifier.IndexInGroup {
			return callInGroup, nil
		}
	}
	return nil, arwen.ErrAsyncCallNotFound
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

	call.Source = context.host.Runtime().GetSCAddress()

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
	// TODO matei-p remove for logging
	// fmt.Println("result produced by function", context.host.Runtime().Function(), " call id", context.callID)
	// retData := context.host.Output().GetVMOutput().ReturnData
	// for d := 0; d < len(retData); d++ {
	// 	data := retData[d]
	// 	if d != len(retData)-1 {
	// 		fmt.Println("\t", data)
	// 	} else {
	// 		fmt.Println("\t" + string(data))
	// 	}
	// }
	// end debug

	context.childResults = context.host.Output().GetVMOutput()

	if context.HasPendingCallGroups() {
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

		// logAsync.Trace("async.Execute() complete locals")
		//
		// This call to closeCompletedAsyncCall() is necessary to remove the
		// AsyncCall that has been just before async.Execute() was called, within
		// host.callSCMethod(). This happens when a cross-shard callback returns and
		// finalizes an AsyncCall.
		// context.closeCompletedAsyncCalls()
		// if context.groupCallbacksEnabled {
		// 	context.executeCompletedGroupCallbacks()
		// }
		// context.deleteCompletedGroups()

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
	} else {
		logAsync.Trace("no async calls")
		if context.IsComplete() && context.callType == vm.DirectCall {
			context, _ = context.LoadFromStackOrStore(context.callerAddr, context.callerCallID)
			context.NotifyChildIsComplete(nil, false)
		}
	}

	// save context and all it's stack parents to store
	// (if exists on stack are either sync or local async calls)
	if !context.IsComplete() {
		context.Save()
		for _, stackContext := range context.stateStack {
			stackContext.host = context.host
			stackContext.Save()
		}
	}

	// TODO matei-p change to debug logging
	fmt.Println("GasLeft ->", context.host.Metering().GasLeft())

	return nil
}

func (context *asyncContext) LoadFromStackOrStore(address []byte, callID []byte) (*asyncContext, bool) {
	stackContext := context.getContextFromStack(address, callID)
	if stackContext != nil {
		return stackContext, false
	}
	context.LoadSpecifiedContext(address, callID)
	return context, true
}

func (context *asyncContext) ReadSerializedFromStackOrStore(address []byte, callID []byte) (*serializableAsyncContext, bool) {
	stackContext := context.getContextFromStack(address, callID)
	if stackContext != nil {
		return stackContext.toSerializable(), false
	}

	serializedContext, _ := NewSerializedAsyncContextFromStore(context.host.Storage(), address, callID)
	return serializedContext, true
}

func (context *asyncContext) removeAsyncCallIfCompleted(asyncCallIdentifier *arwen.AsyncCallIdentifier, returnCode vmcommon.ReturnCode) error {
	asyncCall, err := context.GetCallByAsyncIdentifier(asyncCallIdentifier)
	if err != nil {
		return err
	}
	// The vmOutput instance returned by host.ExecuteOnDestContext() is never nil,
	// by design. Using it without checking for err is safe here.
	asyncCall.UpdateStatus(returnCode)

	context.closeCompletedAsyncCalls()
	if context.groupCallbacksEnabled {
		context.executeCompletedGroupCallbacks()
	}
	context.deleteCompletedGroups()

	return nil
}

func (context *asyncContext) AreAllChildrenComplete() (bool, error) {
	return context.IsStoredContextComplete(context.address, context.callID)
}

func (context *asyncContext) IsStoredContextComplete(address []byte, callID []byte) (bool, error) {
	// fmt.Println("check completed childdren for address", string(address), "callID", DebugPartialArrayToString(callID))
	serializedAsync, err := context.readSerializedAsyncContextFromStackOrStore(address, callID)
	if err != nil && err != arwen.ErrNoStoredAsyncContextFound {
		return false, err
	}
	if serializedAsync == nil {
		return true, nil
	}
	return serializedAsync.IsComplete(), nil
}

func (context *serializableAsyncContext) IsComplete() bool {
	return context.CallsCounter == 0 && len(context.AsyncCallGroups) == 0
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
func (context *asyncContext) PostprocessCrossShardCallback(callID []byte, asyncCallIdentifier *arwen.AsyncCallIdentifier) error {
	// TODO matei-p uncomment or delete if function is no longer used

	// currentGroupID := asyncCallIdentifier.GroupIdentifier
	// asyncCallIndex := asyncCallIdentifier.IndexInGroup

	// currentCallGroup, ok := context.GetCallGroup(currentGroupID)
	// if !ok {
	// 	return arwen.ErrCallBackFuncNotExpected
	// }

	// currentCallGroup.DeleteAsyncCall(asyncCallIndex)
	// if currentCallGroup.HasPendingCalls() {
	// 	return nil
	// }

	// if context.groupCallbacksEnabled {
	// 	// The current group expects no more callbacks, so its own callback can be
	// 	// executed now.
	// 	context.executeCallGroupCallback(currentCallGroup)
	// }

	// context.deleteCallGroupByID(currentGroupID)
	// // Are we still waiting for callbacks to return?
	// if context.HasPendingCallGroups() {
	// 	return nil
	// }

	// // There are no more callbacks to return from other shards. The context can
	// // be deleted from storage.
	// err := context.Delete()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (context *asyncContext) NotifyChildIsComplete(asyncCallIdentifier *arwen.AsyncCallIdentifier, isCrossShardCallChain bool) error {
	// TODO matei-p remove for logging
	fmt.Println("NofityChildIsComplete")
	fmt.Println("\taddress", string(context.address))
	fmt.Println("\tcallID", DebugCallIDAsString(context.callID))
	fmt.Println("\tcallerAddr", string(context.callerAddr))
	fmt.Println("\tcallerCallID", DebugCallIDAsString(context.callerCallID))

	context.DecrementCallsCounter()

	if asyncCallIdentifier != nil {
		currentGroupID := asyncCallIdentifier.GroupIdentifier
		asyncCallIndex := asyncCallIdentifier.IndexInGroup

		currentCallGroup, ok := context.GetCallGroup(currentGroupID)
		if !ok {
			return arwen.ErrCallBackFuncNotExpected
		}

		currentCallGroup.DeleteAsyncCall(asyncCallIndex)

		if context.groupCallbacksEnabled {
			// The current group expects no more callbacks, so its own callback can be
			// executed now.
			context.executeCallGroupCallback(currentCallGroup)
		}

		if currentCallGroup.IsComplete() {
			context.deleteCallGroupByID(currentGroupID)
		}
	}

	context.Save()

	if context.IsComplete() {
		// There are no more callbacks to return from other shards. The context can
		// be deleted from storage.
		err := context.Delete()
		if err != nil {
			return err
		}

		if context.callerCallID == nil {
			// first call, stop notification chain
			return nil
		}

		if isCrossShardCallChain {
			currentAsyncCallIdentifier, _ := context.GetCallerAsyncCallIdentifier()
			// callback for an completed async call is called here only if notifications
			// started with a cross shard call, otherwise it will be called in regular
			// async call code, after the local async call is completed
			if context.callType == vm.AsynchronousCall /*&& isCrossShardCallChain*/ {
				vmOutput := context.childResults
				context.LoadSpecifiedContext(context.callerAddr, context.callerCallID)
				// TODO matei-p remove
				// context, _ := context.LoadFromStackOrStore(context.callerAddr, context.callerCallID)
				isComplete, _ := context.CallCallbackForCompleteAsyncCrossShardCall(currentAsyncCallIdentifier, vmOutput)
				if isComplete {
					// check if we reached the first call and need to end the notification chain
					if asyncCallIdentifier != nil {
						context.NotifyChildIsComplete(currentAsyncCallIdentifier, isCrossShardCallChain)
					}
				}
			} else if context.callType == vm.AsynchronousCallBack {
				// TODO matei-p remove
				// var parentContext *serializableAsyncContext
				// stackContext := context.getContextFromStack(context.address, context.callerCallID)
				// if stackContext != nil {
				// 	parentContext = stackContext.toSerializable()
				// } else {
				// 	parentContext, _ = context.ReadSerializedFromStackOrStore(context.address, context.callerCallID)
				// }
				parentContext, _ := context.ReadSerializedFromStackOrStore(context.address, context.callerCallID)
				asyncCallInParent, _ := parentContext.GetCallByAsyncIdentifier(currentAsyncCallIdentifier)
				if asyncCallInParent != nil && asyncCallInParent.ExecutionMode == arwen.AsyncUnknown {
					// TODO matei-p remove
					// context, _ := context.LoadFromStackOrStore(context.address, context.callerCallID)
					context.LoadSpecifiedContext(context.address, context.callerCallID)
					context.NotifyChildIsComplete(currentAsyncCallIdentifier, isCrossShardCallChain)
				}
			} else if context.callType == vm.DirectCall {
				// TODO matei-p remove
				// context, _ := context.LoadFromStackOrStore(context.callerAddr, context.callerCallID)
				context.LoadSpecifiedContext(context.callerAddr, context.callerCallID)
				context.NotifyChildIsComplete(nil, isCrossShardCallChain)
			}
		}
	}

	return nil
}

func (context *asyncContext) CallCallbackForCompleteAsyncCrossShardCall(asyncCallIdentifier *arwen.AsyncCallIdentifier, vmOutput *vmcommon.VMOutput) (bool, error) {
	asyncCallInParent, _ := context.GetCallByAsyncIdentifier(asyncCallIdentifier)
	return context.callCallback(asyncCallInParent, vmOutput, nil)
}

func (context *asyncContext) callCallback(asyncCall *arwen.AsyncCall, vmOutput *vmcommon.VMOutput, err error) (bool, error) {
	sender := asyncCall.Destination
	destination := asyncCall.Source

	sameShard := context.host.AreInSameShard(sender, destination)
	if !sameShard {
		data := context.GetEncodedDataForAsyncCallbackTransfer(asyncCall, vmOutput)
		return false, sendCrossShardCallback(context.host, sender, destination, data)
	}

	isComplete := context.executeSyncCallbackAndAccumulateGas(asyncCall, vmOutput, err)
	return isComplete, nil
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
	// it's possible that the counter is 0, but further async calls will follow so
	// the context is not finished yet
	return context.callsCounter == 0 && len(context.asyncCallGroups) == 0
}

func (context *asyncContext) Save() error {
	return context.saveUsingStorage(context.host.Storage())
}

// Save serializes and saves the AsyncContext to the storage of the contract, under a protected key.
func (context *asyncContext) saveUsingStorage(storage arwen.StorageContext) error {
	address := context.address
	callID := context.callID
	// fmt.Println("save address", string(address), "callID", DebugPartialArrayToString(callID), "callsCounter", context.callsCounter)

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, callID)
	data, err := json.Marshal(context.toSerializable())
	if err != nil {
		return err
	}

	_, err = storage.SetProtectedStorageToAddress(address, storageKey, data)
	if err != nil {
		return err
	}

	return nil
}

func (context *serializableAsyncContext) HasPendingCallGroups() bool {
	return len(context.AsyncCallGroups) > 0
}

// func (context *asyncContext) Load() error {
// 	return context.LoadSpecifiedContext(context.address, context.callID)
// }

// Load restores the internal state of the AsyncContext from the storage of the contract.
func (context *asyncContext) LoadSpecifiedContext(address []byte, callID []byte) error {
	// fmt.Println("loaded address", string(address), "callID", DebugPartialArrayToString(callID), "callsCounter", context.callsCounter)
	loadedContext, err := NewSerializedAsyncContextFromStore(context.host.Storage(), address, callID)
	if err != nil {
		return err
	}

	context.address = loadedContext.Address
	context.callID = loadedContext.CallID
	context.callerAddr = loadedContext.CallerAddr
	context.callerCallID = loadedContext.CallerCallID
	context.callAsyncIdentifierAsBytes = loadedContext.CallerCallAsyncIdentifierAsBytes
	context.callType = loadedContext.CallType
	context.returnData = loadedContext.ReturnData
	context.asyncCallGroups = loadedContext.AsyncCallGroups
	context.callsCounter = loadedContext.CallsCounter
	context.totalCallsCounter = loadedContext.TotalCallsCounter
	context.childResults = loadedContext.ChildResults

	return nil
}

func (context *asyncContext) readSerializedAsyncContextFromStackOrStore(address []byte, callID []byte) (*asyncContext, error) {
	loadedContext := context.getContextFromStack(address, callID)
	if loadedContext == nil {
		storage := context.host.Storage()
		deserializedContext, err := NewSerializedAsyncContextFromStore(storage, address, callID)
		if err != nil {
			return nil, err
		}
		loadedContext = fromSerializable(deserializedContext)

	}

	loadedContext.host = context.host
	loadedContext.callArgsParser = context.callArgsParser
	loadedContext.esdtTransferParser = context.esdtTransferParser

	return loadedContext, nil
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

// NewSerializedAsyncContextFromStore -
func NewSerializedAsyncContextFromStore(storage arwen.StorageContext, address []byte, callID []byte) (*serializableAsyncContext, error) {
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, callID)
	data := storage.GetStorageFromAddressNoChecks(address, storageKey)
	if len(data) == 0 {
		return nil, arwen.ErrNoStoredAsyncContextFound
	}

	deserializedContext, err := deserializeAsyncContext(data)
	if err != nil {
		return nil, err
	}
	return deserializedContext, nil
}

// Delete deletes the persisted state of the AsyncContext from the contract storage.
func (context *asyncContext) Delete() error {
	storage := context.host.Storage()
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, context.callID)
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

	sameShard := context.host.AreInSameShard(runtime.GetSCAddress(), destination)
	if context.host.IsBuiltinFunctionName(functionName) {
		if sameShard {
			vmInput := runtime.GetVMInput()
			isESDTTransfer, _, _ := context.isESDTTransferOnReturnDataFromFunctionAndArgs(
				runtime.GetSCAddress(),
				destination,
				functionName,
				args)
			isAsyncCall := vmInput.CallType == vm.AsynchronousCall
			isReturningCall := bytes.Equal(vmInput.CallerAddr, destination)

			if isESDTTransfer && isAsyncCall && isReturningCall {
				return arwen.ESDTTransferOnCallBack, nil
			}

			return arwen.AsyncBuiltinFuncIntraShard, nil
		}

		return arwen.AsyncBuiltinFuncCrossShard, nil
	}

	code, err := blockchain.GetCode(destination)
	if len(code) > 0 && err == nil {
		return arwen.SyncExecution, nil
	}

	return arwen.AsyncUnknown, nil
}

func (context *asyncContext) sendAsyncCallCrossShard(asyncCall *arwen.AsyncCall) error {
	host := context.host
	runtime := host.Runtime()
	output := host.Output()

	function, arguments, err := context.callArgsParser.ParseData(string(asyncCall.GetData()))
	if err != nil {
		return err
	}

	callData := txDataBuilder.NewBuilder()
	callData.Func(function)
	callData.Bytes(context.GenerateNewCallID())
	// next two are for the callback to identify this async call in callers async context
	callData.Bytes(context.GetCallID())
	callData.Bytes(asyncCall.Identifier.ToBytes())
	for _, argument := range arguments {
		callData.Bytes(argument)
	}

	err = output.Transfer(
		asyncCall.GetDestination(),
		runtime.GetSCAddress(),
		asyncCall.GetGasLimit(),
		asyncCall.GetGasLocked(),
		big.NewInt(0).SetBytes(asyncCall.GetValue()),
		callData.ToBytes(),
		vm.AsynchronousCall,
	)
	if err != nil {
		return err
	}

	return nil
}

// executeAsyncContextCallback will either execute a sync call (in-shard) to
// the original caller by invoking its callback directly, or will dispatch a
// cross-shard callback to it.
func (context *asyncContext) executeContextCallback() error {
	if !context.HasCallback() {
		// TODO decide whether context.gasAccumulated should be restored here to
		// mark it as available for VMOutput.GasRemaining
		return nil
	}

	callbackCallInput := context.createContextCallbackInput()
	callbackVMOutput, callBackErr, _ := context.host.ExecuteOnDestContext(callbackCallInput)
	context.finishAsyncLocalExecution(callbackVMOutput, callBackErr)

	return nil
}

// TODO compare with host.sendAsyncCallbackToCaller()
func sendCrossShardCallback(host arwen.VMHost, sender []byte, destination []byte, data []byte) error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	gasLeft := metering.GasLeft()

	err := output.Transfer(
		destination,
		sender,
		gasLeft, /// TODO matei-p de discutat cu camil
		0,
		currentCall.CallValue,
		data,
		vm.AsynchronousCallBack,
	)
	metering.UseGas(gasLeft)
	if err != nil {
		runtime.FailExecution(err)
		return err
	}

	log.Trace(
		"sendAsyncCallbackToCaller",
		"caller", currentCall.CallerAddr,
		"data", data,
		"gas", gasLeft)

	return nil

	/*
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

		log.Trace(
			"sendContextCallbackToOriginalCaller",
			"caller", context.callerAddr,
			"data", context.returnData,
			"gas", context.gasAccumulated)

		return nil
	*/
}

// GetEncodedDataForAsyncCallbackTransfer -
func (context *asyncContext) GetEncodedDataForAsyncCallbackTransfer(asyncCall *arwen.AsyncCall, vmOutput *vmcommon.VMOutput) []byte {
	transferData := txDataBuilder.NewBuilder()
	/*
		callbackFunction is not used by arwen, used by testing frameworks
		arwen uses callerCallAsyncIdentifierAsBytes to get the AsyncCall with function name, gas provided etc.
	*/
	// async := host.Async()
	// transferData.Func(callbackFunction)
	// transferData.Bytes(async.GenerateNewCallbackID())
	// transferData.Bytes(callerCallID)
	// transferData.Bytes(callerCallAsyncIdentifierAsBytes)
	transferData.Func(asyncCall.SuccessCallback)
	transferData.Bytes(context.GenerateNewCallbackID())
	transferData.Bytes(context.callID)
	transferData.Bytes(asyncCall.Identifier.ToBytes())

	retCode := vmOutput.ReturnCode

	// TODO matei-p should include return code / data from all children
	transferData.Int64(int64(retCode))
	if retCode == vmcommon.Ok {
		for _, data := range vmOutput.ReturnData {
			transferData.Bytes(data)
		}
	} else {
		transferData.Str(vmOutput.ReturnMessage)
	}
	return transferData.ToBytes()
}

func (context *asyncContext) Serialize() ([]byte, error) {
	serializableContext := context.toSerializable()
	return json.Marshal(serializableContext)
}

func deserializeAsyncContext(data []byte) (*serializableAsyncContext, error) {
	deserializedContext := &serializableAsyncContext{}
	err := json.Unmarshal(data, deserializedContext)
	if err != nil {
		return nil, err
	}
	return deserializedContext, nil
}

func (context *asyncContext) toSerializable() *serializableAsyncContext {
	return &serializableAsyncContext{
		Address:                          context.address,
		CallID:                           context.callID,
		CallerAddr:                       context.callerAddr,
		CallerCallID:                     context.callerCallID,
		CallType:                         context.callType,
		CallerCallAsyncIdentifierAsBytes: context.callAsyncIdentifierAsBytes,
		Callback:                         context.callback,
		CallbackData:                     context.callbackData,
		GasPrice:                         context.gasPrice,
		GasAccumulated:                   context.gasAccumulated,
		ReturnData:                       context.returnData,
		AsyncCallGroups:                  context.asyncCallGroups,
		CallsCounter:                     context.callsCounter,
		TotalCallsCounter:                context.totalCallsCounter,
		ChildResults:                     context.childResults,
	}
}

func fromSerializable(serializedContext *serializableAsyncContext) *asyncContext {
	return &asyncContext{
		host:                       nil,
		stateStack:                 nil,
		address:                    serializedContext.Address,
		callID:                     serializedContext.CallID,
		callsCounter:               serializedContext.CallsCounter,
		totalCallsCounter:          serializedContext.TotalCallsCounter,
		callerAddr:                 serializedContext.CallerAddr,
		callerCallID:               serializedContext.CallerCallID,
		callType:                   serializedContext.CallType,
		callAsyncIdentifierAsBytes: serializedContext.CallerCallAsyncIdentifierAsBytes,
		callback:                   serializedContext.Callback,
		callbackData:               serializedContext.CallbackData,
		gasPrice:                   serializedContext.GasPrice,
		gasAccumulated:             serializedContext.GasAccumulated,
		returnData:                 serializedContext.ReturnData,
		asyncCallGroups:            serializedContext.AsyncCallGroups,
		childResults:               serializedContext.ChildResults,
	}
}

func (context *asyncContext) findGroupByID(groupID string) (int, bool) {
	return findGroupByIDInAsyncCallGroups(context.asyncCallGroups, groupID)
}

func findGroupByIDInAsyncCallGroups(asyncCallGroups []*arwen.AsyncCallGroup, groupID string) (int, bool) {
	for index, group := range asyncCallGroups {
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

func (context *asyncContext) GetCallID() []byte {
	return context.callID
}

func (context *asyncContext) GenerateNewCallID() []byte {
	return context.generateNewCallID(false)
}

func (context *asyncContext) GenerateNewCallbackID() []byte {
	return context.generateNewCallID(true)
}

func (context *asyncContext) generateNewCallID(isCallback bool) []byte {
	if !isCallback {
		context.callsCounter++
	}
	context.totalCallsCounter++
	newCallID := append(context.callID, big.NewInt(int64(context.totalCallsCounter)).Bytes()...)
	newCallID, _ = context.host.Crypto().Sha256(newCallID)
	// TODO matei-p only for debug purposes
	newCallID = append([]byte("_"), newCallID...)
	newCallID = append([]byte(strconv.Itoa(int(context.totalCallsCounter))), newCallID...)
	newCallID = append([]byte("_"), newCallID...)
	newCallID = append([]byte(context.host.Runtime().Function()), newCallID...)
	return newCallID
}

func (context *asyncContext) DecrementCallsCounter() {
	// fmt.Println("---> decrement " + string(context.address) + " " + DebugCallIDAsString(context.callID))
	context.callsCounter--
}

// DebugCallIDAsString - just for debug purposes
func DebugCallIDAsString(arr []byte) string {
	if len(arr) > 3 {
		return "[" + string(arr)[:5] + "...]"
	}
	return fmt.Sprint(arr)
}
