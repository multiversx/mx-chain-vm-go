package contexts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

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

	callerAddr                   []byte
	callerCallID                 []byte
	callbackAsyncInitiatorCallID []byte

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

type SerializableAsyncContext struct {
	Address  []byte
	CallID   []byte
	CallType vm.CallType

	CallerAddr                   []byte
	CallerCallID                 []byte
	CallbackAsyncInitiatorCallID []byte

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
	context.callbackAsyncInitiatorCallID = nil
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
	fmt.Println("\tcallID", context.callID) // DebugCallIDAsString(
	if input.CallType == vm.AsynchronousCallBack {
		context.callbackAsyncInitiatorCallID = runtime.GetAndEliminateFirstArgumentFromList()
		// TODO matei-p change to debug logging
		fmt.Println("\tcallerAddr", string(context.callerAddr))
		fmt.Println("\tcallerCallID", context.callerCallID)
		fmt.Println("\tcallbackAsyncInitiatorCallID", context.callbackAsyncInitiatorCallID)
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
	newState := &asyncContext{
		address:                      context.address,
		callerAddr:                   context.callerAddr,
		callerCallID:                 context.callerCallID,
		callType:                     context.callType,
		callbackAsyncInitiatorCallID: context.callbackAsyncInitiatorCallID,
		callback:                     context.callback,
		callbackData:                 context.callbackData,
		gasPrice:                     context.gasPrice,
		gasAccumulated:               context.gasAccumulated,
		returnData:                   context.returnData,
		asyncCallGroups:              context.asyncCallGroups, // TODO matei-p use cloneCallGroups()?
		callID:                       context.callID,
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
		address:                      context.address,
		callerAddr:                   context.callerAddr,
		callerCallID:                 context.callerCallID,
		callType:                     context.callType,
		callbackAsyncInitiatorCallID: context.callbackAsyncInitiatorCallID,
		callback:                     context.callback,
		callbackData:                 context.callbackData,
		gasPrice:                     context.gasPrice,
		gasAccumulated:               context.gasAccumulated,
		returnData:                   context.returnData,
		asyncCallGroups:              context.cloneCallGroups(),
		callID:                       context.callID,
		callsCounter:                 context.callsCounter,
		totalCallsCounter:            context.totalCallsCounter,
		childResults:                 context.childResults,
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
func (context *asyncContext) UpdateCurrentAsyncCallStatus(address []byte, callID []byte, asyncCallIdentifier []byte, vmInput *vmcommon.VMInput) (*arwen.AsyncCall, error) {
	deserializedContext, err := NewSerializedAsyncContextFromStore(context.host.Storage(), address, context.callbackAsyncInitiatorCallID)
	if err != nil {
		return nil, err
	}

	if vmInput.CallType != vm.AsynchronousCallBack {
		return nil, nil
	}

	if len(vmInput.Arguments) == 0 {
		return nil, arwen.ErrCannotInterpretCallbackArgs
	}

	call, _, _, err := deserializedContext.GetCallByAsyncIdentifier(asyncCallIdentifier)
	if err != nil {
		return nil, err
	}

	// The first argument of the callback is the return code of the destination call
	destReturnCode := big.NewInt(0).SetBytes(vmInput.Arguments[0]).Uint64()
	call.UpdateStatus(vmcommon.ReturnCode(destReturnCode))

	return call, nil
}

func (context *asyncContext) GetCallByAsyncIdentifier(asyncCallIdentifier []byte) (*arwen.AsyncCall, int, int, error) {
	return getCallByAsyncIdentifier(context.asyncCallGroups, asyncCallIdentifier)
}

func (context *SerializableAsyncContext) GetCallByAsyncIdentifier(asyncCallIdentifier []byte) (*arwen.AsyncCall, int, int, error) {
	return getCallByAsyncIdentifier(context.AsyncCallGroups, asyncCallIdentifier)
}

func getCallByAsyncIdentifier(groups []*arwen.AsyncCallGroup, asyncCallIdentifier []byte) (*arwen.AsyncCall, int, int, error) {
	for groupIndex, group := range groups {
		for call1Index, callInGroup := range group.AsyncCalls {
			if bytes.Equal(callInGroup.Identifier, asyncCallIdentifier) {
				return callInGroup, groupIndex, call1Index, nil
			}
		}
	}

	return nil, -1, -1, arwen.ErrAsyncCallNotFound
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

	//call.Identifier = context.GenerateNewCallID()
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

func (context *asyncContext) removeAsyncCallIfCompleted(asyncCallIdentifier []byte, returnCode vmcommon.ReturnCode) error {
	asyncCall, _, _, err := context.GetCallByAsyncIdentifier(asyncCallIdentifier)
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
	serializedAsync, err := NewSerializedAsyncContextFromStore(context.host.Storage(), address, callID)
	if err != nil && err != arwen.ErrNoStoredAsyncContextFound {
		return false, err
	}
	if serializedAsync == nil {
		return true, nil
	}
	return serializedAsync.IsComplete(), nil
}

func (context *SerializableAsyncContext) IsComplete() bool {
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

func (context *asyncContext) NotifyChildIsComplete(asyncCallIdentifier []byte, isCrossShardCallChain bool) error {
	// TODO matei-p remove for logging
	fmt.Println("NofityChildIsComplete")
	fmt.Println("\taddress", string(context.address))
	fmt.Println("\tcallID", context.callID) // DebugCallIDAsString
	fmt.Println("\tcallerAddr", string(context.callerAddr))
	fmt.Println("\tcallerCallID", context.callerCallID)
	fmt.Println("\tasyncCallIdentifier", asyncCallIdentifier)

	context.DecrementCallsCounter()

	if asyncCallIdentifier != nil {

		_, groupIndex, callIndex, err := context.GetCallByAsyncIdentifier(asyncCallIdentifier)
		if err != nil {
			return err
		}

		currentCallGroup := context.asyncCallGroups[groupIndex]
		currentCallGroup.DeleteAsyncCall(callIndex)

		if context.groupCallbacksEnabled {
			// The current group expects no more callbacks, so its own callback can be
			// executed now.
			context.executeCallGroupCallback(currentCallGroup)
		}

		if currentCallGroup.IsComplete() {
			context.deleteCallGroup(groupIndex)
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

		// if first call, stop notification chain
		if context.IsFirstCall() {
			return nil
		}

		if isCrossShardCallChain {
			currentAsyncCallIdentifier := context.GetCallID()
			if err != nil {
				return err
			}
			// callback for an completed async call is called here only if notifications
			// started with a cross shard call, otherwise it will be called in regular
			// async call code, after the local async call is completed
			if context.callType == vm.AsynchronousCall {
				vmOutput := context.childResults
				// context.LoadParentContext()
				isComplete, _ := context.CallCallback(currentAsyncCallIdentifier, vmOutput, nil)
				if isComplete {
					if !context.IsFirstCall() {
						context.NotifyChildIsComplete(currentAsyncCallIdentifier, isCrossShardCallChain)
					}
				}
			} else if context.callType == vm.AsynchronousCallBack {
				context.LoadParentContext()
				context.NotifyChildIsComplete(currentAsyncCallIdentifier, isCrossShardCallChain)
			} else if context.callType == vm.DirectCall {
				context.LoadParentContext()
				context.NotifyChildIsComplete(nil, isCrossShardCallChain)
			}
		}
	}

	return nil
}

func (context *asyncContext) CallCallback(asyncCallIdentifier []byte, vmOutput *vmcommon.VMOutput, err error) (bool, error) {
	sender := context.address
	destination := context.callerAddr

	sameShard := context.host.AreInSameShard(sender, destination)
	if !sameShard {
		data := context.GetEncodedDataForAsyncCallbackTransfer(vmOutput)
		return false, sendCrossShardCallback(context.host, sender, destination, data)
	}

	context.LoadParentContext()
	asyncCall, _, _, errLoad := context.GetCallByAsyncIdentifier(asyncCallIdentifier)
	if errLoad != nil {
		return false, errLoad
	}
	isComplete := context.executeSyncCallbackAndAccumulateGas(asyncCall, vmOutput, err)
	return isComplete, nil
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

// IsComplete returns true if there are no more AsyncCallGroups contained in the AsyncContext.
func (context *asyncContext) IsComplete() bool {
	// it's possible that the counter is 0, but further async calls will follow so
	// the context is not finished yet
	return context.callsCounter == 0 && len(context.asyncCallGroups) == 0
}

// Save serializes and saves the AsyncContext to the storage of the contract, under a protected key.
func (context *asyncContext) Save() error {
	address := context.address
	callID := context.callID
	storage := context.host.Storage()
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

func (context *SerializableAsyncContext) HasPendingCallGroups() bool {
	return len(context.AsyncCallGroups) > 0
}

func (context *asyncContext) LoadParentContext() error {
	if context.callType != vm.AsynchronousCallBack {
		return context.LoadSpecifiedContext(context.callerAddr, context.callerCallID)
	}
	return context.LoadSpecifiedContext(context.address, context.callbackAsyncInitiatorCallID)
}

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
	context.callbackAsyncInitiatorCallID = loadedContext.CallbackAsyncInitiatorCallID
	context.callType = loadedContext.CallType
	context.returnData = loadedContext.ReturnData
	context.asyncCallGroups = loadedContext.AsyncCallGroups
	context.callsCounter = loadedContext.CallsCounter
	context.totalCallsCounter = loadedContext.TotalCallsCounter
	context.childResults = loadedContext.ChildResults

	return nil
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
func NewSerializedAsyncContextFromStore(storage arwen.StorageContext, address []byte, callID []byte) (*SerializableAsyncContext, error) {
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, callID)
	data := storage.GetStorageFromAddressNoChecks(address, storageKey)
	if len(data) == 0 {
		return nil, arwen.ErrNoStoredAsyncContextFound
	}

	deserializedContext, err := DeserializeAsyncContext(data)
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

	newCallID := context.GenerateNewCallID()
	callData := txDataBuilder.NewBuilder()
	callData.Func(function)
	callData.Bytes(newCallID)
	callData.Bytes(context.GetCallID())

	asyncCall.Identifier = newCallID

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
}

// GetEncodedDataForAsyncCallbackTransfer -
func (context *asyncContext) GetEncodedDataForAsyncCallbackTransfer(vmOutput *vmcommon.VMOutput) []byte {
	transferData := txDataBuilder.NewBuilder()

	transferData.Func("<callback>")
	transferData.Bytes(context.GenerateNewCallbackID())
	transferData.Bytes(context.callID)
	transferData.Bytes(context.callerCallID)

	retCode := vmOutput.ReturnCode

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

func DeserializeAsyncContext(data []byte) (*SerializableAsyncContext, error) {
	deserializedContext := &SerializableAsyncContext{}
	err := json.Unmarshal(data, deserializedContext)
	if err != nil {
		return nil, err
	}
	return deserializedContext, nil
}

func (context *asyncContext) toSerializable() *SerializableAsyncContext {
	return &SerializableAsyncContext{
		Address:                      context.address,
		CallID:                       context.callID,
		CallerAddr:                   context.callerAddr,
		CallerCallID:                 context.callerCallID,
		CallType:                     context.callType,
		CallbackAsyncInitiatorCallID: context.callbackAsyncInitiatorCallID,
		Callback:                     context.callback,
		CallbackData:                 context.callbackData,
		GasPrice:                     context.gasPrice,
		GasAccumulated:               context.gasAccumulated,
		ReturnData:                   context.returnData,
		AsyncCallGroups:              context.asyncCallGroups,
		CallsCounter:                 context.callsCounter,
		TotalCallsCounter:            context.totalCallsCounter,
		ChildResults:                 context.childResults,
	}
}

func fromSerializable(serializedContext *SerializableAsyncContext) *asyncContext {
	return &asyncContext{
		host:                         nil,
		stateStack:                   nil,
		address:                      serializedContext.Address,
		callID:                       serializedContext.CallID,
		callsCounter:                 serializedContext.CallsCounter,
		totalCallsCounter:            serializedContext.TotalCallsCounter,
		callerAddr:                   serializedContext.CallerAddr,
		callerCallID:                 serializedContext.CallerCallID,
		callType:                     serializedContext.CallType,
		callbackAsyncInitiatorCallID: serializedContext.CallbackAsyncInitiatorCallID,
		callback:                     serializedContext.Callback,
		callbackData:                 serializedContext.CallbackData,
		gasPrice:                     serializedContext.GasPrice,
		gasAccumulated:               serializedContext.GasAccumulated,
		returnData:                   serializedContext.ReturnData,
		asyncCallGroups:              serializedContext.AsyncCallGroups,
		childResults:                 serializedContext.ChildResults,
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
	// newCallID = append([]byte("_"), newCallID...)
	// newCallID = append([]byte(strconv.Itoa(int(context.totalCallsCounter))), newCallID...)
	// newCallID = append([]byte("_"), newCallID...)
	// newCallID = append([]byte(context.host.Runtime().Function()), newCallID...)
	return newCallID
}

func (context *asyncContext) DecrementCallsCounter() {
	context.callsCounter--
}

// DebugCallIDAsString - just for debug purposes
func DebugCallIDAsString(arr []byte) string {
	if len(arr) > 3 {
		return "[" + string(arr)[:5] + "...]"
	}
	return fmt.Sprint(arr)
}