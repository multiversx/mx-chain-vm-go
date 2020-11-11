package host

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

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
 * which means that the corresponding AsyncCall must be deleted from the
 * current AsyncContext.

 * Moreover, because individual AsyncCalls are contained by AsyncCallGroups, we
 * must verify whether the containing AsyncCallGroup has any remaining calls
 * pending. If not, the final callback of the containing AsyncCallGroup must be
 * executed as well (note, though, that callbacks of AsyncCallGroups are
 * currently disabled).
 */
func (host *vmHost) postprocessCrossShardCallback() error {
	runtime := host.Runtime()

	asyncContext, err := host.getCurrentAsyncContext()
	if err != nil {
		return err
	}

	// TODO FindAsyncCallByDestination() only returns the first matched AsyncCall
	// by destination, but there could be multiple matches in an AsyncContext.
	vmInput := runtime.GetVMInput()
	currentGroupID, asyncCallIndex, err := asyncContext.FindAsyncCallByDestination(vmInput.CallerAddr)
	if err != nil {
		return arwen.ErrCallBackFuncNotExpected
	}

	currentCallGroup := asyncContext.AsyncCallGroups[currentGroupID]
	currentCallGroup.DeleteAsyncCall(asyncCallIndex)

	if currentCallGroup.HasPendingCalls() {
		return nil
	}

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
	callbackCallInput, err := host.createSyncCallbackInput(
		host.Output().GetVMOutput(),
		asyncContext.CallerAddr,
		arwen.CallbackDefault,
		nil,
	)
	if err != nil {
		return err
	}

	callbackVMOutput, callBackErr := host.ExecuteOnDestContext(callbackCallInput)
	host.finishSyncExecution(callbackVMOutput, callBackErr)

	return nil
}

// TODO move into RuntimeContext
func (host *vmHost) getFunctionByCallType(callType vmcommon.CallType) (wasmer.ExportedFunctionCallback, error) {
	runtime := host.Runtime()

	if callType != vmcommon.AsynchronousCallBack {
		return runtime.GetFunctionToCall()
	}

	asyncContext, err := host.getCurrentAsyncContext()
	if err != nil {
		return nil, err
	}

	vmInput := runtime.GetVMInput()

	customCallback := false
	for _, asyncCallGroup := range asyncContext.AsyncCallGroups {
		for _, asyncCall := range asyncCallGroup.AsyncCalls {
			if bytes.Equal(vmInput.CallerAddr, asyncCall.Destination) {
				customCallback = true
				// TODO why asyncCall.SuccessCallback? why not check for error as well,
				// and set asyncCall.ErrorCallback?
				runtime.SetCustomCallFunction(asyncCall.SuccessCallback)
				break
			}
		}

		if customCallback {
			break
		}
	}

	return runtime.GetFunctionToCall()
}

func (host *vmHost) getCurrentAsyncContext() (*arwen.AsyncContext, error) {
	runtime := host.Runtime()
	storage := host.Storage()

	asyncContext := &arwen.AsyncContext{}
	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
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

	storageKey := arwen.CustomStorageKey(arwen.AsyncDataPrefix, runtime.GetOriginalTxHash())
	_, err := storage.SetStorage(storageKey, nil)
	return err
}
