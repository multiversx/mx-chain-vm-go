package arwen

// AsyncCallGroup is a structure containing a group of async calls and a callback
// that should be called when all these async calls are resolved
type AsyncCallGroup struct {
	Callback   string
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
