package contexts

import (
	"math/big"
)

type bigIntMap map[int32]*big.Int

type bigIntContext struct {
	values     bigIntMap
	stateStack []bigIntMap
}

// NewBigIntContext creates a new bigIntContext
func NewBigIntContext() (*bigIntContext, error) {
	context := &bigIntContext{
		values:     make(bigIntMap),
		stateStack: make([]bigIntMap, 0),
	}

	return context, nil
}

// InitState initializes the underlying values map
func (context *bigIntContext) InitState() {
	context.values = make(bigIntMap)
}

// PushState appends the values map to the state stack
func (context *bigIntContext) PushState() {
	newState := context.clone()
	context.stateStack = append(context.stateStack, newState)
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current values map
func (context *bigIntContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevValues := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.values = prevValues
}

// PopDiscard removes the latest entry from the state stack
func (context *bigIntContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack initializes the state stack
func (context *bigIntContext) ClearStateStack() {
	context.stateStack = make([]bigIntMap, 0)
}

func (context *bigIntContext) clone() bigIntMap {
	newState := make(bigIntMap, len(context.values))
	for handle, bigInt := range context.values {
		newState[handle] = big.NewInt(0).Set(bigInt)
	}
	return newState
}

// Put adds the given value to the current values map and returns the handle
func (context *bigIntContext) Put(value int64) int32 {
	newHandle := int32(len(context.values))
	for {
		if _, ok := context.values[newHandle]; !ok {
			break
		}
		newHandle++
	}

	context.values[newHandle] = big.NewInt(value)

	return newHandle
}

// GetOne returns the value at the given handle. If there is no value under that handle, it will return 0
func (context *bigIntContext) GetOne(handle int32) *big.Int {
	if _, ok := context.values[handle]; !ok {
		context.values[handle] = big.NewInt(0)
	}

	return context.values[handle]
}

// GetTwo returns the values at the given handles.
func (context *bigIntContext) GetTwo(handle1 int32, handle2 int32) (*big.Int, *big.Int) {
	return context.GetOne(handle1), context.GetOne(handle2)
}

// GetThree returns the values at the given handles.
func (context *bigIntContext) GetThree(handle1 int32, handle2 int32, handle3 int32) (*big.Int, *big.Int, *big.Int) {
	return context.GetOne(handle1), context.GetOne(handle2), context.GetOne(handle3)
}

// IsInterfaceNil returns true if there is no value under the interface
func (context *bigIntContext) IsInterfaceNil() bool {
	return context == nil
}
