package host

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	// TODO ensure the default async call is deleted either by
	// executeSyncCallback() or by postprocessCrossShardCallback()
	asyncCall := runtime.GetDefaultAsyncCall()
	err := host.executeAsyncCall(asyncCall, false)
	if err != nil {
		return err
	}

	return nil
}

func (host *vmHost) sendAsyncCallbackToCaller() error {
	runtime := host.Runtime()
	if !runtime.GetAsyncContext().IsCompleted() {
		return nil
	}

	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	retData := []byte("@" + hex.EncodeToString([]byte(output.ReturnCode().String())))
	for _, data := range output.ReturnData() {
		retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
	}

	err := output.Transfer(
		currentCall.CallerAddr,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		0,
		currentCall.CallValue,
		retData,
		vmcommon.AsynchronousCallBack,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

func (host *vmHost) sendContextCallbackToOriginalCaller(asyncContext *arwen.AsyncContext) error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	err := output.Transfer(
		asyncContext.CallerAddr,
		runtime.GetSCAddress(),
		metering.GasLeft(),
		0,
		currentCall.CallValue,
		asyncContext.ReturnData,
		vmcommon.AsynchronousCallBack,
	)
	if err != nil {
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
		return err
	}

	return nil
}

/**
 * postprocessCrossShardCallback() is called by host.callSCMethod() after it
 * has locally executed the callback of a returning cross-shard AsyncCall,
 * which means that the AsyncContext corresponding to the original transaction
 * must be loaded from storage, and then the corresponding AsyncCall must be
 * deleted from the current AsyncContext.

 * TODO because individual AsyncCalls are contained by AsyncCallGroups, we
 * must verify whether the containing AsyncCallGroup has any remaining calls
 * pending. If not, the final callback of the containing AsyncCallGroup must be
 * executed as well.
 */
func (host *vmHost) postprocessCrossShardCallback() error {
	asyncContext, err := host.loadCurrentAsyncContext()
	if err != nil {
		return err
	}

	// TODO FindAsyncCallByDestination() only returns the first matched AsyncCall
	// by destination, but there could be multiple matches in an AsyncContext.
	vmInput := host.Runtime().GetVMInput()
	currentGroupID, asyncCallIndex, err := asyncContext.FindAsyncCallByDestination(vmInput.CallerAddr)
	if err != nil {
		return arwen.ErrCallBackFuncNotExpected
	}

	currentCallGroup, ok := asyncContext.GetAsyncCallGroup(currentGroupID)
	if !ok {
		return arwen.ErrCallBackFuncNotExpected
	}

	currentCallGroup.DeleteAsyncCall(asyncCallIndex)
	if currentCallGroup.HasPendingCalls() {
		return nil
	}

	asyncContext.DeleteAsyncCallGroupByID(currentGroupID)
	// Are we still waiting for callbacks to return?
	if asyncContext.HasPendingCallGroups() {
		return nil
	}

	err = host.deleteCurrentAsyncContext()
	if err != nil {
		return err
	}

	return host.executeAsyncContextCallback(asyncContext)
}

// executeAsyncContextCallback will either execute a sync call (in-shard) to
// the original caller by invoking its callback directly, or will dispatch a
// cross-shard callback to it.
func (host *vmHost) executeAsyncContextCallback(asyncContext *arwen.AsyncContext) error {
	execMode, err := host.determineExecutionMode(asyncContext.CallerAddr, asyncContext.ReturnData)
	if err != nil {
		return err
	}

	if execMode != arwen.SyncExecution {
		return host.sendContextCallbackToOriginalCaller(asyncContext)
	}

	// The caller is in the same shard, execute its callback
	callbackCallInput := host.createSyncContextCallbackInput(asyncContext)

	callbackVMOutput, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	host.finishSyncExecution(callbackVMOutput, callBackErr)

	return nil
}

func (host *vmHost) createSyncContextCallbackInput(asyncContext *arwen.AsyncContext) *vmcommon.ContractCallInput {
	runtime := host.Runtime()
	metering := host.Metering()

	_, arguments, err := host.CallArgsParser().ParseData(string(asyncContext.ReturnData))
	if err != nil {
		arguments = [][]byte{asyncContext.ReturnData}
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
		RecipientAddr: asyncContext.CallerAddr,
		Function:      arwen.CallbackFunctionName, // TODO currently default; will customize in AsynContext
	}
	return input
}

func (host *vmHost) loadCurrentAsyncContext() (*arwen.AsyncContext, error) {
	runtime := host.Runtime()
	storage := host.Storage()

	asyncContext := &arwen.AsyncContext{}
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetPrevTxHash())
	buff := storage.GetStorage(storageKey)
	if len(buff) == 0 {
		return asyncContext, nil
	}

	err := json.Unmarshal(buff, &asyncContext)
	if err != nil {
		return nil, err
	}

	return asyncContext, nil
}

func (host *vmHost) deleteCurrentAsyncContext() error {
	runtime := host.Runtime()
	storage := host.Storage()

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetPrevTxHash())
	_, err := storage.SetStorage(storageKey, nil)
	return err
}
