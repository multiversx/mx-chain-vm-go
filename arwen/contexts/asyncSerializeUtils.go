package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var marshalizer = &marshal.GogoProtoMarshalizer{}

func deserializeAsyncContext(data []byte) (*asyncContext, error) {
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
		CallerCallID:                 context.callerCallID,
		CallbackAsyncInitiatorCallID: context.callbackAsyncInitiatorCallID,
		Callback:                     context.callback,
		CallbackData:                 context.callbackData,
		GasPrice:                     context.gasPrice,
		GasAccumulated:               context.gasAccumulated,
		ReturnData:                   context.returnData,
		AsyncCallGroups:              arwen.ToSerializableAsyncCallGroups(context.asyncCallGroups),
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
		callerCallID:                 serializedContext.CallerCallID,
		callType:                     vm.CallType(serializedContext.CallType),
		callbackAsyncInitiatorCallID: serializedContext.CallbackAsyncInitiatorCallID,
		callback:                     serializedContext.Callback,
		callbackData:                 serializedContext.CallbackData,
		gasPrice:                     serializedContext.GasPrice,
		gasAccumulated:               serializedContext.GasAccumulated,
		returnData:                   serializedContext.ReturnData,
		asyncCallGroups:              arwen.FromSerializableAsyncCallGroups(serializedContext.AsyncCallGroups),
		childResults:                 fromSerializableVMOutput(serializedContext.ChildResults),
	}
}

// GetCallByAsyncIdentifier -
func (context *SerializableAsyncContext) GetCallByAsyncIdentifier(asyncCallIdentifier []byte) (*arwen.AsyncCall, int, int, error) {
	for groupIndex, group := range context.AsyncCallGroups {
		for callIndex, callInGroup := range group.AsyncCalls {
			if bytes.Equal(callInGroup.CallID, asyncCallIdentifier) {
				return callInGroup.FromSerializable(), groupIndex, callIndex, nil
			}
		}
	}

	return nil, -1, -1, arwen.ErrAsyncCallNotFound
}

// IsComplete -
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
