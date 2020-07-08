package arwen

import "bytes"

// AsyncContext is the structure resulting after a smart contract call that has initiated
// one or more async calls. It contains all of the async calls produced by the
// smart contract method.
type AsyncContext struct {
	CallerAddr      []byte
	ReturnData      []byte
	AsyncCallGroups map[string]*AsyncCallGroup
}

func (actx *AsyncContext) HasPendingCallGroups() bool {
	return len(actx.AsyncCallGroups) > 0
}

func (actx *AsyncContext) IsCompleted() bool {
	return len(actx.AsyncCallGroups) == 0
}

func (actx *AsyncContext) MakeAsyncContextWithPendingCalls() *AsyncContext {
	pendingGroups := make(map[string]*AsyncCallGroup)
	for groupID, asyncCallGroup := range actx.AsyncCallGroups {
		for _, asyncCall := range asyncCallGroup.AsyncCalls {
			if asyncCall.Status != AsyncCallPending {
				continue
			}

			pendingGroup, ok := pendingGroups[groupID]
			if !ok {
				pendingGroup = &AsyncCallGroup{
					Callback:   asyncCallGroup.Callback,
					AsyncCalls: make([]*AsyncCall, 0),
				}
				pendingGroups[groupID] = pendingGroup
			}
			pendingGroup.AsyncCalls = append(pendingGroup.AsyncCalls, asyncCall)
		}
	}

	return &AsyncContext{
		CallerAddr:      actx.CallerAddr,
		ReturnData:      actx.ReturnData,
		AsyncCallGroups: pendingGroups,
	}
}

func (actx *AsyncContext) FindAsyncCallByDestination(destination []byte) (string, int, error) {
	var groupID string
	var asyncCallIndex int
	// TODO ranging over a map has unpredictable order, which can be exploited as
	// a randomness source by malicious smart contracts
	for id, asyncCallGroup := range actx.AsyncCallGroups {
		for position, asyncCall := range asyncCallGroup.AsyncCalls {
			if bytes.Equal(destination, asyncCall.Destination) {
				asyncCallIndex = position
				groupID = id
				break
			}
		}

		if len(groupID) > 0 {
			break
		}
	}

	if len(groupID) == 0 {
		return "", -1, ErrAsyncCallNotFound
	}

	return groupID, asyncCallIndex, nil
}

func (actx *AsyncContext) DeleteAsyncCallGroup(groupID string) {
	delete(actx.AsyncCallGroups, groupID)
}
