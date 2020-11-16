package arwen

import "bytes"

// AsyncCallGroup is a structure containing a group of async calls and a callback
// that should be called when all these async calls are resolved
type AsyncCallGroup struct {
	// TODO re-enable AsyncCallGroup.Callback after the rest of the functionality works
	// Callback string
	Identifier string
	AsyncCalls []*AsyncCall
}

// NewAsyncCallGroup creates a new instance of AsyncCallGroup
func NewAsyncCallGroup(identifier string) *AsyncCallGroup {
	return &AsyncCallGroup{
		Identifier: identifier,
		AsyncCalls: make([]*AsyncCall, 0),
	}
}

// AddAsyncCall adds a given AsyncCall to the AsyncCallGroup
func (acg *AsyncCallGroup) AddAsyncCall(call *AsyncCall) {
	acg.AsyncCalls = append(acg.AsyncCalls, call)
}

// HasPendingCalls verifies whether the AsyncCallGroup has any AsyncCalls left
// to return from the destination call
func (acg *AsyncCallGroup) HasPendingCalls() bool {
	return len(acg.AsyncCalls) > 0
}

// IsCompleted verifies whether all AsyncCalls have been completed
func (acg *AsyncCallGroup) IsCompleted() bool {
	return len(acg.AsyncCalls) == 0
}

// FindByDestination returns the index of an AsyncCall in this AsyncCallGroup
// that matches the provided destination
func (acg *AsyncCallGroup) FindByDestination(destination []byte) (int, bool) {
	for index, call := range acg.AsyncCalls {
		if bytes.Equal(destination, call.Destination) {
			return index, true
		}
	}
	return -1, false
}

// DeleteAsyncCall removes an AsyncCall from this AsyncCallGroup, given its index
func (acg *AsyncCallGroup) DeleteAsyncCall(index int) {
	asyncCalls := acg.AsyncCalls
	if len(asyncCalls) == 0 {
		return
	}

	last := len(asyncCalls) - 1
	if index < 0 || index > last {
		return
	}

	asyncCalls[index] = asyncCalls[last]
	asyncCalls = asyncCalls[:last]
	acg.AsyncCalls = asyncCalls
}

// DeleteCompletedAsyncCalls removes all completed AsyncCalls, keeping only
// those with status AsyncCallPending
func (acg *AsyncCallGroup) DeleteCompletedAsyncCalls() {
	remainingAsyncCalls := make([]*AsyncCall, 0)
	for _, asyncCall := range acg.AsyncCalls {
		if asyncCall.Status == AsyncCallPending {
			remainingAsyncCalls = append(remainingAsyncCalls, asyncCall)
		}
	}

	acg.AsyncCalls = remainingAsyncCalls
}
