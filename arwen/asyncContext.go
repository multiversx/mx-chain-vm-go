package arwen

// AsyncContext is the structure resulting after a smart contract call that has initiated
// one or more async calls. It contains all of the async calls produced by the
// smart contract method.
type AsyncContext struct {
	CallerAddr      []byte
	ReturnData      []byte
	AsyncCallGroups []*AsyncCallGroup
}

// NewAsyncContext creates a new instance of AsyncContext
func NewAsyncContext() *AsyncContext {
	return &AsyncContext{
		CallerAddr:      nil,
		ReturnData:      nil,
		AsyncCallGroups: make([]*AsyncCallGroup, 0),
	}
}

// AddAsyncGroup adds the provided AsyncCallGroup to the AsyncContext
func (actx *AsyncContext) AddAsyncGroup(group *AsyncCallGroup) {
	actx.AsyncCallGroups = append(actx.AsyncCallGroups, group)
}

// HasPendingCallGroups verifies whether the AsyncContext has any
// AsyncCallGroups yet to complete
func (actx *AsyncContext) HasPendingCallGroups() bool {
	return len(actx.AsyncCallGroups) > 0
}

// IsCompleted verifies whether all the AsyncCallGroups in the AsyncContext
// have been completed
func (actx *AsyncContext) IsCompleted() bool {
	return len(actx.AsyncCallGroups) == 0
}

// MakeAsyncContextWithPendingCalls creates a new AsyncContext containing only
// the pending AsyncCallGroups, without deleting anything from the initial AsyncContext
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

// FindAsyncCallByDestination retrieves the AsyncCall which matches the given
// destination, from within the AsyncCallGroups
func (actx *AsyncContext) FindAsyncCallByDestination(destination []byte) (string, int, error) {
	for _, group := range actx.AsyncCallGroups {
		callIndex, ok := group.FindByDestination(destination)
		if ok {
			return group.Identifier, callIndex, nil
		}
	}

	return "", -1, ErrAsyncCallNotFound
}

// GetAsyncCallGroup retrieves an AsyncCallGroup by its Identifier
func (actx *AsyncContext) GetAsyncCallGroup(groupID string) (*AsyncCallGroup, bool) {
	index, ok := findGroupByID(actx.AsyncCallGroups, groupID)
	if ok {
		return actx.AsyncCallGroups[index], true
	}

	return nil, false
}

// DeleteAsyncCallGroupByID deletes an AsyncCallGroup by its Identifier
func (actx *AsyncContext) DeleteAsyncCallGroupByID(groupID string) {
	index, ok := findGroupByID(actx.AsyncCallGroups, groupID)
	if !ok {
		return
	}
	actx.DeleteAsyncCallGroup(index)
}

// DeleteAsyncCallGroup deletes an AsyncCallGroup by its index
func (actx *AsyncContext) DeleteAsyncCallGroup(index int) {
	groups := actx.AsyncCallGroups
	if len(groups) == 0 {
		return
	}

	last := len(groups) - 1
	if index < 0 || index > last {
		return
	}

	groups[index] = groups[last]
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
