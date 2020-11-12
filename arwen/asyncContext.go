package arwen

// AsyncContext is the structure resulting after a smart contract call that has initiated
// one or more async calls. It contains all of the async calls produced by the
// smart contract method.
type AsyncContext struct {
	CallerAddr      []byte
	ReturnData      []byte
	AsyncCallGroups []*AsyncCallGroup
}

func NewAsyncContext() *AsyncContext {
	return &AsyncContext{
		CallerAddr:      nil,
		ReturnData:      nil,
		AsyncCallGroups: make([]*AsyncCallGroup, 0),
	}
}

func (actx *AsyncContext) AddAsyncGroup(group *AsyncCallGroup) {
	actx.AsyncCallGroups = append(actx.AsyncCallGroups, group)
}

func (actx *AsyncContext) HasPendingCallGroups() bool {
	return len(actx.AsyncCallGroups) > 0
}

func (actx *AsyncContext) IsCompleted() bool {
	return len(actx.AsyncCallGroups) == 0
}

func (actx *AsyncContext) MakeAsyncContextWithPendingCalls() *AsyncContext {
	pendingGroups := make([]*AsyncCallGroup, 0)
	var pendingGroup *AsyncCallGroup
	for _, group := range actx.AsyncCallGroups {
		pendingGroup = nil
		for _, asyncCall := range group.AsyncCalls {
			if asyncCall.Status != AsyncCallPending {
				continue
			}

			if pendingGroup == nil {
				pendingGroup = NewAsyncCallGroup(group.Identifier)
				pendingGroups = append(pendingGroups, pendingGroup)
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
	for _, group := range actx.AsyncCallGroups {
		callIndex, ok := group.FindByDestination(destination)
		if ok {
			return group.Identifier, callIndex, nil
		}
	}

	return "", -1, ErrAsyncCallNotFound
}

func (actx *AsyncContext) GetAsyncCallGroup(groupID string) (*AsyncCallGroup, bool) {
	index, ok := actx.FindAsyncCallGroup(groupID)
	if ok {
		return actx.AsyncCallGroups[index], true
	}

	return nil, false
}

func (actx *AsyncContext) FindAsyncCallGroup(groupID string) (int, bool) {
	return findGroupByID(actx.AsyncCallGroups, groupID)
}

func (actx *AsyncContext) DeleteAsyncCallGroupByID(groupID string) {
	index, ok := actx.FindAsyncCallGroup(groupID)
	if !ok {
		return
	}
	actx.DeleteAsyncCallGroup(index)
}

func (actx *AsyncContext) DeleteAsyncCallGroup(index int) {
	groups := actx.AsyncCallGroups
	last := len(groups) - 1
	groups[index] = groups[last]
	groups[last] = nil
	groups = groups[:last]

	actx.AsyncCallGroups = groups
}

func findGroupByID(groups []*AsyncCallGroup, groupID string) (int, bool) {
	for index, group := range groups {
		if group.Identifier == groupID {
			return index, true
		}
	}
	return -1, false
}
