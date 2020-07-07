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

func (actx *AsyncContext) HasWaitingCallGroups() bool {
	return len(actx.AsyncCallGroups) > 0
}

func (actx *AsyncContext) FindAsyncCallByDestination(destination []byte) (string, int, error) {
	var groupID string
	var asyncCallPosition int
	for id, asyncCallGroup := range actx.AsyncCallGroups {
		for position, asyncCall := range asyncCallGroup.AsyncCalls {
			if bytes.Equal(destination, asyncCall.Destination) {
				asyncCallPosition = position
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

	return groupID, asyncCallPosition, nil
}
