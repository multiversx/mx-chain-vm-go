package contexts

import (
	"errors"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-core-go/marshal"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// Save serializes and saves the AsyncContext to the storage of the contract, under a protected key.
func (context *asyncContext) Save() error {
	address := context.address
	callID := context.callID
	storage := context.host.Storage()

	if len(callID) > vmhost.AddressLen {
		return errors.New("callID must be 32 bytes")
	}

	storageKey := getAsyncContextStorageKey(context.asyncStorageDataPrefix, callID)
	data, err := context.marshalizer.Marshal(context.toSerializable())
	if err != nil {
		return err
	}

	_, err = storage.SetProtectedStorageToAddressUnmetered(address, storageKey, data)
	if err != nil {
		return err
	}

	return nil
}

// LoadParentContext loads AsyncContext from the storage of the contract using the caller id
func (context *asyncContext) LoadParentContext() error {
	switch context.callType {
	case vm.DirectCall:
		return context.loadSpecificContext(context.callerAddr, context.callerCallID)
	case vm.AsynchronousCallBack:
		// parent is the same as the callback, and the id is callbackAsyncInitiatorCallID
		return context.loadSpecificContext(context.address, context.callbackAsyncInitiatorCallID)
	default:
		return vmhost.ErrNoAsyncParentContext
	}
}

// DeleteFromCallID deletes the persisted state of the AsyncContext from the contract storage.
func (context *asyncContext) DeleteFromCallID(callID []byte) error {
	storage := context.host.Storage()

	// Delete AsyncContext
	storageKey := getAsyncContextStorageKey(context.asyncStorageDataPrefix, callID)
	_, err := storage.SetProtectedStorageToAddressUnmetered(context.address, storageKey, nil)
	if err != nil {
		return err
	}

	// Delete AsyncResults
	resultsKey := getAsyncContextStorageKey(storage.GetVmProtectedPrefix(vmhost.AsyncResultsPrefix), callID)
	_, err = storage.SetProtectedStorageToAddressUnmetered(context.address, resultsKey, nil)
	return err
}

// LoadParentContextFromStackOrStorage loads the parent async comtext
func (context *asyncContext) LoadParentContextFromStackOrStorage() (vmhost.AsyncContext, error) {
	if context.callType != vm.AsynchronousCallBack {
		return context.loadFromStackOrStorage(context.callerAddr, context.callerCallID)
	}
	return context.loadFromStackOrStorage(context.address, context.callbackAsyncInitiatorCallID)
}

func (context *asyncContext) loadFromStackOrStorage(address []byte, callID []byte) (*asyncContext, error) {
	stackContext := context.getContextFromStack(address, callID)
	if stackContext != nil {
		return stackContext, nil
	}
	err := context.loadSpecificContext(address, callID)
	return context, err
}

// Load restores the internal state of the AsyncContext from the storage of the contract.
func (context *asyncContext) loadSpecificContext(address []byte, callID []byte) error {
	loadedContext, err := readAsyncContextFromStorage(
		context.host.Storage(),
		address,
		callID,
		context.marshalizer)
	if err != nil {
		return err
	}

	context.address = loadedContext.address
	context.callID = loadedContext.callID
	context.callerAddr = loadedContext.callerAddr
	context.parentAddr = loadedContext.parentAddr
	context.callerCallID = loadedContext.callerCallID
	context.callbackAsyncInitiatorCallID = loadedContext.callbackAsyncInitiatorCallID
	context.callType = loadedContext.callType
	context.returnData = loadedContext.returnData
	context.asyncCallGroups = loadedContext.asyncCallGroups
	context.callsCounter = loadedContext.callsCounter
	context.totalCallsCounter = loadedContext.totalCallsCounter
	context.childResults = loadedContext.childResults
	context.gasAccumulated = loadedContext.gasAccumulated

	return nil
}

func readAsyncContextFromStorage(
	storage vmhost.StorageContext,
	address []byte,
	callID []byte,
	marshalizer *marshal.GogoProtoMarshalizer,
) (*asyncContext, error) {
	storageKey := getAsyncContextStorageKey(storage.GetVmProtectedPrefix(vmhost.AsyncDataPrefix), callID)
	data, _, _, _ := storage.GetStorageFromAddressNoChecks(address, storageKey)
	if len(data) == 0 {
		return nil, vmhost.ErrNoStoredAsyncContextFound
	}

	return deserializeAsyncContext(data, marshalizer)
}

func deserializeAsyncContext(data []byte, marshalizer *marshal.GogoProtoMarshalizer) (*asyncContext, error) {
	deserializedAsyncContext := &SerializableAsyncContext{}
	err := marshalizer.Unmarshal(deserializedAsyncContext, data)
	if err != nil {
		return nil, err
	}

	return fromSerializable(deserializedAsyncContext), nil
}

func (context *asyncContext) toSerializable() *SerializableAsyncContext {
	return &SerializableAsyncContext{
		Address:                      context.address,
		CallID:                       context.callID,
		CallType:                     SerializableCallType(context.callType),
		CallerAddr:                   context.callerAddr,
		ParentAddr:                   context.parentAddr,
		CallerCallID:                 context.callerCallID,
		CallbackAsyncInitiatorCallID: context.callbackAsyncInitiatorCallID,
		Callback:                     context.callback,
		CallbackData:                 context.callbackData,
		GasAccumulated:               context.gasAccumulated,
		ReturnData:                   context.returnData,
		AsyncCallGroups:              vmhost.ToSerializableAsyncCallGroups(context.asyncCallGroups),
		CallsCounter:                 context.callsCounter,
		TotalCallsCounter:            context.totalCallsCounter,
		ChildResults:                 toSerializableVMOutput(context.childResults),
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
		parentAddr:                   serializedContext.ParentAddr,
		callerCallID:                 serializedContext.CallerCallID,
		callType:                     vm.CallType(serializedContext.CallType),
		callbackAsyncInitiatorCallID: serializedContext.CallbackAsyncInitiatorCallID,
		callback:                     serializedContext.Callback,
		callbackData:                 serializedContext.CallbackData,
		gasAccumulated:               serializedContext.GasAccumulated,
		returnData:                   serializedContext.ReturnData,
		asyncCallGroups:              vmhost.FromSerializableAsyncCallGroups(serializedContext.AsyncCallGroups),
		childResults:                 fromSerializableVMOutput(serializedContext.ChildResults),
	}
}

// IsComplete returns true if no more async calls are pending
func (context *SerializableAsyncContext) IsComplete() bool {
	return context.CallsCounter == 0 && len(context.AsyncCallGroups) == 0
}

func toSerializableVMOutput(vmOutput *vmcommon.VMOutput) *SerializableVMOutput {
	if vmOutput == nil {
		return nil
	}

	return &SerializableVMOutput{
		ReturnData:    vmOutput.ReturnData,
		ReturnCode:    uint64(vmOutput.ReturnCode),
		ReturnMessage: vmOutput.ReturnMessage,
		GasRemaining:  vmOutput.GasRemaining,
	}
}

func fromSerializableVMOutput(serializedVMOutput *SerializableVMOutput) *vmcommon.VMOutput {
	if serializedVMOutput == nil {
		return nil
	}
	return &vmcommon.VMOutput{
		ReturnData:    serializedVMOutput.ReturnData,
		ReturnCode:    vmcommon.ReturnCode(serializedVMOutput.ReturnCode),
		ReturnMessage: serializedVMOutput.ReturnMessage,
		GasRemaining:  serializedVMOutput.GasRemaining,
		GasRefund:     serializedVMOutput.GasRefund,
	}
}

func getAsyncContextStorageKey(prefix []byte, callID []byte) []byte {
	return vmhost.CustomStorageKey(string(prefix), callID)
}
