package arwen

// AsyncCallGroup is a structure containing a group of async calls and a callback
// that should be called when all these async calls are resolved
type AsyncCallGroup struct {
	// TODO re-enable AsyncCallGroup.Callback after the rest of the functionality works
	// Callback string
	AsyncCalls []*AsyncCall
}

func (acg *AsyncCallGroup) HasPendingCalls() bool {
	return len(acg.AsyncCalls) > 0
}

func (acg *AsyncCallGroup) IsCompleted() bool {
	return len(acg.AsyncCalls) == 0
}

func (acg *AsyncCallGroup) DeleteAsyncCall(index int) {
	asyncCalls := acg.AsyncCalls
	asyncCalls[index] = asyncCalls[len(asyncCalls)-1]
	asyncCalls[len(asyncCalls)-1] = nil
	asyncCalls = asyncCalls[:len(asyncCalls)-1]
	acg.AsyncCalls = asyncCalls
}

func (acg *AsyncCallGroup) DeleteCompletedAsyncCalls() {
	remainingAsyncCalls := make([]*AsyncCall, 0)
	for _, asyncCall := range acg.AsyncCalls {
		if asyncCall.Status == AsyncCallPending {
			remainingAsyncCalls = append(remainingAsyncCalls, asyncCall)
		}
	}

	acg.AsyncCalls = remainingAsyncCalls
}
