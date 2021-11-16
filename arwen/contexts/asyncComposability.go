package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (context *asyncContext) NotifyChildIsComplete(callID []byte, gasToAccumulate uint64) (arwen.AsyncContext, error) {
	if logAsync.GetLevel() == logger.LogTrace {
		logAsync.Trace("NofityChildIsComplete")
		logAsync.Trace("", "address", string(context.address))
		logAsync.Trace("", "callID", context.callID) // DebugCallIDAsString
		logAsync.Trace("", "callerAddr", string(context.callerAddr))
		logAsync.Trace("", "callerCallID", context.callerCallID)
		logAsync.Trace("", "callID", callID)
		logAsync.Trace("", "gasToAccumulate", gasToAccumulate)
	}

	context.CompleteChild(callID, gasToAccumulate)

	if !context.IsComplete() {
		// store changes in context made by CompleteChild()
		err := context.Save()
		if err != nil {
			return nil, err
		}
	} else {
		// There are no more callbacks to return from other shards. The context can
		// be deleted from storage.
		err := context.DeleteFromAddress(context.address)
		// err := context.Delete()
		if err != nil {
			return nil, err
		}

		// if we reached first call, stop notification chain
		if context.IsFirstCall() {
			return context, nil
		}

		currentCallID := context.GetCallID()
		gasAccumulatedInNotifingContext := context.gasAccumulated
		if context.callType == vm.AsynchronousCall {
			vmOutput := context.childResults
			isComplete, _, err := context.callCallback(currentCallID, vmOutput, nil)
			if err != nil {
				return nil, err
			}
			if isComplete {
				return context.NotifyChildIsComplete(currentCallID, 0)
			}
		} else if context.callType == vm.AsynchronousCallBack {
			currentCallID := context.GetCallerCallID()
			context.LoadParentContext()
			return context.NotifyChildIsComplete(currentCallID, gasAccumulatedInNotifingContext)
		} else if context.callType == vm.DirectCall {
			context.LoadParentContext()
			return context.NotifyChildIsComplete(nil, gasAccumulatedInNotifingContext)
		}
	}

	return context, nil
}

func (context *asyncContext) CompleteChild(callID []byte, gasToAccumulate uint64) error {
	return context.CompleteChildConditional(true, callID, gasToAccumulate)
}

func (context *asyncContext) CompleteChildConditional(isComplete bool, callID []byte, gasToAccumulate uint64) error {
	if !isComplete {
		return nil
	}
	context.DecrementCallsCounter()
	context.accumulateGas(gasToAccumulate)
	if callID != nil {
		err := context.DeleteAsyncCallAndCleanGroup(callID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (context *asyncContext) removeAsyncCallIfCompleted(
	callID []byte,
	returnCode vmcommon.ReturnCode,
) error {
	asyncCall, _, _, err := context.GetAsyncCallByCallID(callID)
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
