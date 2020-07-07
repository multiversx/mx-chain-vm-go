package async

// AsyncCallGroup is a structure containing a group of async calls and a callback
// that should be called when all these async calls are resolved
type AsyncCallGroup struct {
	Callback   string
	AsyncCalls []*AsyncCall
}

func (acg *AsyncCallGroup) HasWaitingCalls() bool {
	return len(acg.AsyncCalls) > 0
}
