package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var logAsync = logger.GetOrCreate("arwen/async")

var marshalizer = &marshal.GogoProtoMarshalizer{}

func (context *asyncContext) Serialize() ([]byte, error) {
	serializableContext := context.toSerializable()
	return marshalizer.Marshal(serializableContext)
}

func deserializeAsyncContext(data []byte) (*SerializableAsyncContextProto, error) {
	deserializedContext := &SerializableAsyncContextProto{}
	err := marshalizer.Unmarshal(deserializedContext, data)
	if err != nil {
		return nil, err
	}
	return deserializedContext, nil
}

func (context *asyncContext) toSerializable() *SerializableAsyncContextProto {
	return &SerializableAsyncContextProto{
		Address:                      context.address,
		CallID:                       context.callID,
		CallerAddr:                   context.callerAddr,
		CallerCallID:                 context.callerCallID,
		CallType:                     SerializableCallType(context.callType),
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

func fromSerializable(serializedContext *SerializableAsyncContextProto) *asyncContext {
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
	return &vmcommon.VMOutput{
		ReturnData:    serializedVMOutput.ReturnData,
		ReturnCode:    vmcommon.ReturnCode(serializedVMOutput.ReturnCode),
		ReturnMessage: serializedVMOutput.ReturnMessage,
		GasRemaining:  serializedVMOutput.GasRemaining,
		// TODO matei-p update async.proto for all big.Int fields
		// GasRefund:     serVMOutput.GasRefund,
	}
}
