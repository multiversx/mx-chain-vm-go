package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var marshalizer = &marshal.GogoProtoMarshalizer{}

func (context *asyncContext) Serialize() ([]byte, error) {
	serializableContext := context.toSerializable()
	return marshalizer.Marshal(serializableContext)
}

func deserializeAsyncContext(data []byte) (*SerializableAsyncContext, error) {
	deserializedContext := &SerializableAsyncContext{}
	err := marshalizer.Unmarshal(deserializedContext, data)
	if err != nil {
		return nil, err
	}
	return deserializedContext, nil
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
