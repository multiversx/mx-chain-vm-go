package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (context *asyncContext) executeSynchronousCalls() error {
	for groupIndex, group := range context.AsyncCallGroups {
		for _, call := range group.AsyncCalls {
			if call.ExecutionMode != arwen.SyncExecution {
				continue
			}

			err := context.executeSyncCall(call)
			if err != nil {
				return err
			}
		}

		group.DeleteCompletedAsyncCalls()

		// If all the AsyncCalls in the AsyncCallGroup were executed synchronously,
		// then the AsyncCallGroup can have its callback executed.
		if group.IsCompleted() {
			// TODO reenable this, after allowing a gas limit for it and deciding what
			// arguments it receives (this method is currently a NOP and returns nil)
			err := context.executeAsyncCallGroupCallback(group)
			if err != nil {
				return err
			}

			context.DeleteCallGroup(groupIndex)
		}
	}
	return nil
}

func (context *asyncContext) executeSyncCall(asyncCall *arwen.AsyncCall) error {
	destinationCallInput, err := context.createSyncCallInput(asyncCall)
	if err != nil {
		return err
	}

	vmOutput, err := context.host.ExecuteOnDestContext(destinationCallInput)

	// The vmOutput instance returned by host.executeSyncCall() is never nil,
	// by design. Using it without checking for err is safe here.
	asyncCall.UpdateStatus(vmOutput.ReturnCode)

	// TODO host.executeSyncCallback() returns a vmOutput produced by executing
	// the callback. Information from this vmOutput should be preserved in the
	// pending AsyncCallGroup, and made available to the callback of the
	// AsyncCallGroup (currently not implemented).
	callbackVMOutput, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, err)
	context.finishSyncExecution(callbackVMOutput, callbackErr)

	return nil
}

func (context *asyncContext) executeSyncCallback(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	err error,
) (*vmcommon.VMOutput, error) {

	callbackInput, err := context.createSyncCallbackInput(asyncCall, vmOutput, err)
	if err != nil {
		return nil, err
	}

	return context.host.ExecuteOnDestContext(callbackInput)
}

func (context *asyncContext) finishSyncExecution(vmOutput *vmcommon.VMOutput, err error) {
	if err == nil {
		return
	}

	runtime := context.host.Runtime()
	output := context.host.Output()

	runtime.GetVMInput().GasProvided = 0

	if vmOutput == nil {
		vmOutput = output.CreateVMOutputInCaseOfError(err)
	}

	output.SetReturnMessage(vmOutput.ReturnMessage)

	output.Finish([]byte(vmOutput.ReturnCode.String()))
	output.Finish(runtime.GetCurrentTxHash())
}

func (context *asyncContext) createSyncCallInput(asyncCall arwen.AsyncCallHandler) (*vmcommon.ContractCallInput, error) {
	host := context.host
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()

	function, arguments, err := host.CallArgsParser().ParseData(string(asyncCall.GetData()))
	if err != nil {
		return nil, err
	}

	gasLimit := asyncCall.GetGasLimit()
	gasToUse := host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0).SetBytes(asyncCall.GetValueBytes()),
			CallType:       vmcommon.AsynchronousCall,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: asyncCall.GetDestination(),
		Function:      function,
	}

	return contractCallInput, nil
}

func (context *asyncContext) createSyncCallbackInput(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.ContractCallInput, error) {
	metering := context.host.Metering()
	runtime := context.host.Runtime()

	// always provide return code as the first argument to callback function
	arguments := [][]byte{
		big.NewInt(int64(vmOutput.ReturnCode)).Bytes(),
	}
	if destinationErr == nil {
		// when execution went Ok, callBack arguments are:
		// [0, result1, result2, ....]
		arguments = append(arguments, vmOutput.ReturnData...)
	} else {
		// when execution returned error, callBack arguments are:
		// [error code, error message]
		arguments = append(arguments, []byte(vmOutput.ReturnMessage))
	}

	callbackFunction := asyncCall.GetCallbackName()

	gasLimit := vmOutput.GasRemaining + asyncCall.GetGasLocked()
	dataLength := computeDataLengthFromArguments(callbackFunction, arguments)

	gasToUse := metering.GasSchedule().ElrondAPICost.AsyncCallStep
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(dataLength)
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	// Return to the sender SC, calling its specified callback method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     asyncCall.Destination,
			Arguments:      arguments,
			CallValue:      big.NewInt(0),
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      callbackFunction,
	}

	return contractCallInput, nil
}

func (context *asyncContext) createSyncContextCallbackInput() *vmcommon.ContractCallInput {
	host := context.host
	runtime := host.Runtime()
	metering := host.Metering()

	_, arguments, err := host.CallArgsParser().ParseData(string(context.ReturnData))
	if err != nil {
		arguments = [][]byte{context.ReturnData}
	}

	// TODO ensure a new value for VMInput.CurrentTxHash
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     runtime.GetSCAddress(),
			Arguments:      arguments,
			CallValue:      runtime.GetVMInput().CallValue,
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    metering.GasLeft(),
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: context.CallerAddr,
		Function:      arwen.CallbackFunctionName, // TODO currently default; will customize in AsynContext
	}
	return input
}
