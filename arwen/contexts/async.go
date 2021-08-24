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

	address []byte
	callID  []byte

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

	// created for an async call destination
	rootAsyncInfo *rootAsyncContextInfo
}

type rootAsyncContextInfo struct {
	PrevAsyncCallAddress []byte
	PrevAsyncCallID      []byte
	ChildNumberOfCalls   uint32
	ChildResults         []*vmcommon.VMOutput
}

type serializableAsyncContext struct {
	Address []byte
	CallID  []byte

	CallerAddr                       []byte
	CallerCallID                     []byte
	CallerCallAsyncIdentifierAsBytes []byte

	Callback        string
	CallbackData    []byte
	GasPrice        uint64
	GasAccumulated  uint64
	ReturnData      []byte
	AsyncCallGroups []*arwen.AsyncCallGroup
	RootAsyncInfo   *rootAsyncContextInfo
	AddressStack    []*arwen.AddressAndCallID
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
		rootAsyncInfo: &rootAsyncContextInfo{
			ChildNumberOfCalls: 0,
			ChildResults:       nil,
		},
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

	runtime := context.host.Runtime()
	context.address = runtime.GetSCAddress()
	context.callID = runtime.GetCallID()
	if input.CallType == vm.AsynchronousCall {
		context.callerCallID = runtime.GetAndEliminateFirstArgumentFromList()
		context.callAsyncIdentifierAsBytes = runtime.GetAndEliminateFirstArgumentFromList()

		context.rootAsyncInfo.PrevAsyncCallAddress = runtime.GetSCAddress()
		context.rootAsyncInfo.PrevAsyncCallID = runtime.GetCallID()
		context.rootAsyncInfo.ChildNumberOfCalls = 0

		// TODO matei-p remove for logging
		asynCallIdentifier, _ := arwen.ReadAsyncCallIdentifierFromBytes(context.callAsyncIdentifierAsBytes)
		asynCallIdentifierAsString, _ := json.Marshal(asynCallIdentifier)
		fmt.Println("from async input -> callerCallID", context.callerCallID, "callerCallAsyncIdentifier", string(asynCallIdentifierAsString))
		// fmt.Println()
	}
}

// PushState creates a deep clone of the internal state and pushes it onto the
// internal state stack.
func (context *asyncContext) PushState() {
	newState := &asyncContext{
		address:                    context.address,
		callID:                     context.callID,
		callerAddr:                 context.callerAddr,
		callerCallID:               context.callerCallID,
		callAsyncIdentifierAsBytes: context.callAsyncIdentifierAsBytes,
		callback:                   context.callback,
		callbackData:               context.callbackData,
		gasPrice:                   context.gasPrice,
		gasAccumulated:             context.gasAccumulated,
		returnData:                 context.returnData,
		asyncCallGroups:            context.asyncCallGroups, // TODO matei-p use cloneCallGroups()?
		rootAsyncInfo: &rootAsyncContextInfo{
			PrevAsyncCallAddress: context.rootAsyncInfo.PrevAsyncCallAddress,
			PrevAsyncCallID:      context.rootAsyncInfo.PrevAsyncCallID,
			ChildNumberOfCalls:   context.rootAsyncInfo.ChildNumberOfCalls,
			ChildResults:         context.rootAsyncInfo.ChildResults,
		},
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
	context.callAsyncIdentifierAsBytes = prevState.callAsyncIdentifierAsBytes
	context.callback = prevState.callback
	context.callbackData = prevState.callbackData
	context.gasPrice = prevState.gasPrice
	context.gasAccumulated = prevState.gasAccumulated
	context.returnData = prevState.returnData
	context.asyncCallGroups = prevState.asyncCallGroups

	context.rootAsyncInfo = prevState.rootAsyncInfo
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

// func (context *asyncContext) callExists(destination []byte) bool {
// 	_, _, err := context.findCall(destination)
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

// func (context *asyncContext) findCall(destination []byte) (string, int, error) {
// 	return findCallInAsyncCallGroups(context.asyncCallGroups, destination)
// }

// func findCallInAsyncCallGroups(asyncCallGroups []*arwen.AsyncCallGroup, destination []byte) (string, int, error) {
// 	for _, group := range asyncCallGroups {
// 		callIndex, ok := group.FindByDestination(destination)
// 		if ok {
// 			return group.Identifier, callIndex, nil
// 		}
// 	}

// 	return "", -1, arwen.ErrAsyncCallNotFound
// }

// UpdateCurrentAsyncCallStatus detects the AsyncCall returning as callback,
// extracts the ReturnCode from data provided by the destination call, and updates
// the status of the AsyncCall with its value.
func UpdateCurrentAsyncCallStatus(storage arwen.StorageContext, address []byte, callID []byte, asyncCallIdentifier *arwen.AsyncCallIdentifier, vmInput *vmcommon.VMInput) (*arwen.AsyncCall, error) {
	deserializedContext, err := NewSerializedAsyncContextFromStore(storage, address, callID)
	if err != nil {
		return nil, err
	}

	if vmInput.CallType != vm.AsynchronousCallBack {
		return nil, nil
	}

	if len(vmInput.Arguments) == 0 {
		return nil, arwen.ErrCannotInterpretCallbackArgs
	}

	call, err := deserializedContext.getCallByAsyncIdentifier(asyncCallIdentifier)
	if err != nil {
		return nil, err
	}

	// The first argument of the callback is the return code of the destination call
	destReturnCode := big.NewInt(0).SetBytes(vmInput.Arguments[0]).Uint64()
	call.UpdateStatus(vmcommon.ReturnCode(destReturnCode))

	return call, nil
}

func (context *asyncContext) getCallByAsyncIdentifier(asyncCallIdentifier *arwen.AsyncCallIdentifier) (*arwen.AsyncCall, error) {
	return getCallByAsyncIdentifier(context.asyncCallGroups, asyncCallIdentifier)
}

func (context *serializableAsyncContext) getCallByAsyncIdentifier(asyncCallIdentifier *arwen.AsyncCallIdentifier) (*arwen.AsyncCall, error) {
	return getCallByAsyncIdentifier(context.AsyncCallGroups, asyncCallIdentifier)
}

func getCallByAsyncIdentifier(groups []*arwen.AsyncCallGroup, asyncCallIdentifier *arwen.AsyncCallIdentifier) (*arwen.AsyncCall, error) {
	groupIndex, groupFound := findGroupByIDInAsyncCallGroups(groups, asyncCallIdentifier.GroupIdentifier)
	if !groupFound {
		return nil, arwen.ErrAsyncCallNotFound
	}
	callsInGroup := groups[groupIndex].AsyncCalls
	if asyncCallIdentifier.IndexInGroup >= len(callsInGroup) {
		return nil, arwen.ErrAsyncCallNotFound
	}
	return callsInGroup[asyncCallIdentifier.IndexInGroup], nil
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
	if !context.IsComplete() {
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
	} else {
		logAsync.Trace("no async calls")
	}

	runtime := context.host.Runtime()
	areAllChildrenComplete, err := context.areAllChildrenComplete()
	if err != nil {
		return err
	}

	if !context.HasPendingCallGroups() && areAllChildrenComplete {
		if context.contextCallbackEnabled {
			logAsync.Trace("async.Execute() execute context callback")
			context.executeContextCallback()
		}
		if runtime.GetVMInput().CallType == vm.AsynchronousCall {
			context.toSerializable().callCallbackOfOriginalCaller(context.host, context)
		}
	} else {
		logAsync.Trace("async.Execute() save")
		context.rootAsyncInfo.ChildResults = append(context.rootAsyncInfo.ChildResults, context.host.Output().GetVMOutput())
		err := context.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (context *asyncContext) areAllChildrenComplete() (bool, error) {
	runtime := context.host.Runtime()
	return context.areAllChildrenCompleteForStoredContext(runtime.GetSCAddress(), runtime.GetCallID())
}

func (context *asyncContext) areAllChildrenCompleteForStoredContext(address []byte, callID []byte) (bool, error) {
	fmt.Println("check completed childdren for address", string(address), "callID", callID)
	serializedAsync, err := NewSerializedAsyncContextFromStore(context.host.Storage(), address, callID)
	if err != nil && err != arwen.ErrNoStoredAsyncContextFound {
		return false, err
	}
	if serializedAsync == nil || serializedAsync.RootAsyncInfo == nil {
		return true, nil
	}
	return serializedAsync.areAllChildrenComplete(), nil
}

func (context *serializableAsyncContext) areAllChildrenComplete() bool {
	return context.RootAsyncInfo.ChildNumberOfCalls == 0
}

func (context *asyncContext) incrementChildrenNumberForPrevAsyncContext() (*serializableAsyncContext, error) {
	var lastAsyncContext *serializableAsyncContext
	stackCallsList := context.host.Runtime().GetFirstAsyncOrCallbackOnStack()
	for _, stackCall := range stackCallsList {
		storage := context.host.Storage()
		serializedAsync, err := NewSerializedAsyncContextFromStore(storage,
			stackCall.Address,
			stackCall.CallID)
		if err != nil && err != arwen.ErrNoStoredAsyncContextFound {
			return nil, err
		}
		if serializedAsync == nil {
			// async context was never saved so find it on stack
			asyncContext := context
			if stackCall.IndexOnStack != -1 {
				asyncContext = context.stateStack[stackCall.IndexOnStack]
				asyncContext.host = context.host
			}
			serializedAsync = asyncContext.toSerializable()
		}

		serializedAsync.RootAsyncInfo.ChildNumberOfCalls++

		// existing or new, the updated async context will be saved now
		serializedAsync.writeToStorageWithAddress(storage, stackCall.Address, stackCall.CallID)

		lastAsyncContext = serializedAsync
	}
	return lastAsyncContext, nil
}

func (context *serializableAsyncContext) decrementChildrenNumberForPrevAsyncContext(storage arwen.StorageContext) (*serializableAsyncContext, error) {
	var lastAsyncContext *serializableAsyncContext
	for _, stackCall := range context.AddressStack {
		serializedAsync, err := NewSerializedAsyncContextFromStore(storage,
			stackCall.Address,
			stackCall.CallID)
		if err != nil && err != arwen.ErrNoStoredAsyncContextFound {
			return nil, err
		}
		serializedAsync.RootAsyncInfo.ChildNumberOfCalls--

		// existing or new, the updated async context will be saved now
		serializedAsync.writeToStorageWithAddress(storage, stackCall.Address, stackCall.CallID)

		lastAsyncContext = serializedAsync
	}
	return lastAsyncContext, nil
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

	_, err := context.incrementChildrenNumberForPrevAsyncContext()
	if err != nil {
		return err
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
func (serializableContext *serializableAsyncContext) PostprocessCrossShardCallback(host arwen.VMHost, callID []byte, asyncCallIdentifier *arwen.AsyncCallIdentifier) error {
	prevAsyncContext, err := serializableContext.decrementChildrenNumberForPrevAsyncContext(host.Storage())
	if err != nil {
		return err
	}

	context := fromSerializable(serializableContext)
	context.host = host

	if prevAsyncContext.areAllChildrenComplete() {
		prevAsyncContext.callCallbackOfOriginalCaller(host, context)
	}

	currentGroupID := asyncCallIdentifier.GroupIdentifier
	asyncCallIndex := asyncCallIdentifier.IndexInGroup

	defer func() {
		context.Save()
	}()

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

	return nil
}

func (context *serializableAsyncContext) callCallbackOfOriginalCaller(host arwen.VMHost, originalContext *asyncContext) error {
	sender := context.Address
	destination := context.CallerAddr

	data := GetEncodedDataForAsyncCallbackTransfer(host,
		context.CallerCallID,
		context.CallerCallAsyncIdentifierAsBytes)

	sameShard := host.AreInSameShard(sender, destination)
	if !sameShard {
		return sendCrossShardCallback(host, sender, destination, data)
	}

	callerAsyncContext, err := NewSerializedAsyncContextFromStore(host.Storage(),
		context.CallerAddr,
		context.CallerCallID)

	asyncCallIdentifier, err := arwen.ReadAsyncCallIdentifierFromBytes(context.CallerCallAsyncIdentifierAsBytes)
	if err != nil {
		return err
	}

	asyncCall, err := callerAsyncContext.getCallByAsyncIdentifier(asyncCallIdentifier)
	if err != nil {
		return err
	}
	fmt.Println(asyncCall)

	// TODO matei-p err should be the err returned by function()
	originalContext.executeSyncCallbackAndAccumulateGas(
		asyncCall,
		originalContext.rootAsyncInfo.ChildResults[0],
		err)

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
	areAllChildrenComplete, err := context.areAllChildrenComplete()
	if err != nil {
		return err
	}
	if areAllChildrenComplete {
		return nil
	}
	return context.toSerializable().writeToStorage(context.host.Storage(), context.callID)
}

func (context *serializableAsyncContext) writeToStorage(storage arwen.StorageContext, callID []byte) error {
	return context.writeToStorageWithAddress(storage, storage.GetAddress(), callID)
}

func (context *serializableAsyncContext) writeToStorageWithAddress(storage arwen.StorageContext, address []byte, callID []byte) error {
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, callID)
	data, err := json.Marshal(context)
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

// Load restores the internal state of the AsyncContext from the storage of the contract.
func (context *asyncContext) Load(callID []byte) error {
	storage := context.host.Storage()
	address := context.host.Runtime().GetSCAddress()
	deserializedContext, err := NewSerializedAsyncContextFromStore(storage, address, callID)
	if err != nil {
		return err
	}
	loadedContext := fromSerializable(deserializedContext)

	context.callerAddr = loadedContext.callerAddr
	context.returnData = loadedContext.returnData
	context.asyncCallGroups = loadedContext.asyncCallGroups
	context.rootAsyncInfo = loadedContext.rootAsyncInfo

	return nil
}

// NewSerializedAsyncContextFromStore read a new serializableAsyncContext from storage of the address account
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

// // NewRootAsyncInfoFromStore read a new rootAsyncContextInfo from storage of the address account
// func NewRootAsyncInfoFromStore(storage arwen.StorageContext, address []byte, callID []byte) (*rootAsyncContextInfo, error) {
// 	storageKey := arwen.CustomStorageKey(arwen.AsyncRootDataPrefix, callID)
// 	data := storage.GetStorageFromAddressNoChecks(address, storageKey)
// 	if len(data) == 0 {
// 		return nil, arwen.ErrNoStoredAsyncContextFound
// 	}

// 	deserializedRootAsyncInfo, err := deserializeRootAsyncInfo(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return deserializedRootAsyncInfo, nil
// }

// func (rootAsyncInfo *rootAsyncContextInfo) ToBytes() []byte {
// 	data, err := json.Marshal(rootAsyncInfo)
// 	if err != nil {
// 		return nil
// 	}
// 	return data
// }

// func (rootAsyncInfo *rootAsyncContextInfo) setCallBackInfoFromAsyncCall(destination []byte, asyncCall *arwen.AsyncCall) {
// 	rootAsyncInfo.CallBackInfo.SuccessCallback = asyncCall.SuccessCallback
// 	rootAsyncInfo.CallBackInfo.ErrorCallback = asyncCall.ErrorCallback
// 	rootAsyncInfo.CallBackInfo.GasLocked = asyncCall.GasLocked
// }

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
	callData.Bytes(runtime.GenerateNewCallID())
	// next two are for the callback to identify this async call in callers async context
	callData.Bytes(runtime.GetCallID())
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
	callbackVMOutput, callBackErr := context.host.ExecuteOnDestContext(callbackCallInput)
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

func GetEncodedDataForAsyncCallbackTransfer(host arwen.VMHost, callerCallID []byte, callerCallAsyncIdentifierAsBytes []byte) []byte {
	output := host.Output()

	retCode := output.ReturnCode()

	transferData := txDataBuilder.NewBuilder()
	transferData.Func("<calback>")
	transferData.Bytes(host.Runtime().GenerateNewCallID())
	transferData.Bytes(callerCallID)
	transferData.Bytes(callerCallAsyncIdentifierAsBytes)

	// TODO matei-p should include return code / data from all children
	transferData.Int64(int64(retCode))
	if retCode == vmcommon.Ok {
		for _, data := range output.ReturnData() {
			transferData.Bytes(data)
		}
	} else {
		transferData.Str(output.ReturnMessage())
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
	stateStack := context.host.Runtime().GetFirstAsyncOrCallbackOnStack()

	return &serializableAsyncContext{
		Address:                          context.address,
		CallID:                           context.callID,
		CallerAddr:                       context.callerAddr,
		CallerCallID:                     context.callerCallID,
		CallerCallAsyncIdentifierAsBytes: context.callAsyncIdentifierAsBytes,
		Callback:                         context.callback,
		CallbackData:                     context.callbackData,
		GasPrice:                         context.gasPrice,
		GasAccumulated:                   context.gasAccumulated,
		ReturnData:                       context.returnData,
		AsyncCallGroups:                  context.asyncCallGroups,
		RootAsyncInfo:                    context.rootAsyncInfo,
		AddressStack:                     stateStack,
	}
}

func fromSerializable(serializedContext *serializableAsyncContext) *asyncContext {
	return &asyncContext{
		host:                       nil,
		stateStack:                 nil,
		address:                    serializedContext.Address,
		callID:                     serializedContext.CallID,
		callerAddr:                 serializedContext.CallerAddr,
		callerCallID:               serializedContext.CallerCallID,
		callAsyncIdentifierAsBytes: serializedContext.CallerCallAsyncIdentifierAsBytes,
		callback:                   serializedContext.Callback,
		callbackData:               serializedContext.CallbackData,
		gasPrice:                   serializedContext.GasPrice,
		gasAccumulated:             serializedContext.GasAccumulated,
		returnData:                 serializedContext.ReturnData,
		asyncCallGroups:            serializedContext.AsyncCallGroups,
		rootAsyncInfo:              serializedContext.RootAsyncInfo,
	}
}

func (rootAsyncInfo *rootAsyncContextInfo) Serialize() ([]byte, error) {
	return json.Marshal(rootAsyncInfo)
}

func deserializeRootAsyncInfo(data []byte) (*rootAsyncContextInfo, error) {
	deserializedRootAsyncInfo := &rootAsyncContextInfo{}
	err := json.Unmarshal(data, deserializedRootAsyncInfo)
	if err != nil {
		return nil, err
	}
	return deserializedRootAsyncInfo, nil
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
