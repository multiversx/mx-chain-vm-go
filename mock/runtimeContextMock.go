package mock

import (
	"fmt"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type RuntimeContextMock struct {
	instance        *wasmer.Instance
	instanceContext *wasmer.InstanceContext
	vmInput         *vmcommon.VMInput
	scAddress       []byte
	callFunction    string
	readOnly        bool

	stateStack    []*RuntimeContextMock
	instanceStack []*wasmer.Instance
	pointsUsed    uint64
}

func NewRuntimeContextMock() *RuntimeContextMock {
	context := &RuntimeContextMock{
		stateStack:    make([]*RuntimeContextMock, 0),
		instanceStack: make([]*wasmer.Instance, 0),
		pointsUsed:    0,
	}

	context.InitState()

	return context
}

func (context *RuntimeContextMock) InitState() {
	context.vmInput = &vmcommon.VMInput{}
	context.scAddress = make([]byte, 0)
	context.callFunction = ""
	context.readOnly = false
}

func (context *RuntimeContextMock) CreateWasmerInstance(contract []byte, gasLimit uint64) error {
	var err error
	context.instance, err = wasmer.NewMeteredInstance(contract, gasLimit)
	if err != nil {
		return err
	}
	context.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	return nil
}

func (context *RuntimeContextMock) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	context.vmInput = &input.VMInput
	context.scAddress = input.RecipientAddr
	context.callFunction = input.Function
}

func (context *RuntimeContextMock) PushState() {
	newState := &RuntimeContextMock{
		vmInput:      context.vmInput,
		scAddress:    context.scAddress,
		callFunction: context.callFunction,
		readOnly:     context.readOnly,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *RuntimeContextMock) PushInstance() {
	context.instanceStack = append(context.instanceStack, context.instance)
}

func (context *RuntimeContextMock) PopInstance() error {
	instanceStackLen := len(context.instanceStack)
	if instanceStackLen < 1 {
		return arwen.InstanceStackUnderflow
	}

	prevInstance := context.instanceStack[instanceStackLen-1]
	context.instanceStack = context.instanceStack[:instanceStackLen-1]

	context.instance.Clean()
	context.instance = prevInstance

	return nil
}

func (context *RuntimeContextMock) PopState() error {
	stateStackLen := len(context.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.vmInput = prevState.vmInput
	context.scAddress = prevState.scAddress
	context.callFunction = prevState.callFunction
	context.readOnly = prevState.readOnly

	return nil
}

func (context *RuntimeContextMock) GetVMType() []byte {
	return []byte("type")
}

func (context *RuntimeContextMock) GetVMInput() *vmcommon.VMInput {
	return context.vmInput
}

func (context *RuntimeContextMock) SetVMInput(vmInput *vmcommon.VMInput) {
	context.vmInput = vmInput
}

func (context *RuntimeContextMock) GetSCAddress() []byte {
	return context.scAddress
}

func (context *RuntimeContextMock) SetSCAddress(scAddress []byte) {
	context.scAddress = scAddress
}

func (context *RuntimeContextMock) Function() string {
	return context.callFunction
}

func (context *RuntimeContextMock) Arguments() [][]byte {
	return context.vmInput.Arguments
}

func (context *RuntimeContextMock) SignalExit(_ int) {
	context.SetRuntimeBreakpointValue(arwen.BreakpointSignalExit)
}

func (context *RuntimeContextMock) SignalUserError(_ string) {
	// SignalUserError() remains in Runtime, and won't be moved into Output,
	// because there will be extra handling added here later, which requires
	// information from Runtime (e.g. runtime breakpoints)
	context.SetRuntimeBreakpointValue(arwen.BreakpointSignalError)
}

func (context *RuntimeContextMock) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	context.instance.SetBreakpointValue(uint64(value))
}

func (context *RuntimeContextMock) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(context.instance.GetBreakpointValue())
}

func (context *RuntimeContextMock) GetPointsUsed() uint64 {
	return context.pointsUsed
}

func (context *RuntimeContextMock) SetPointsUsed(gasPoints uint64) {
	context.pointsUsed = gasPoints
}

func (context *RuntimeContextMock) ReadOnly() bool {
	return context.readOnly
}

func (context *RuntimeContextMock) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

func (context *RuntimeContextMock) SetInstanceContextId(id int) {
	context.instance.SetContextData(unsafe.Pointer(&id))
}

func (context *RuntimeContextMock) SetInstanceContext(instCtx *wasmer.InstanceContext) {
	context.instanceContext = instCtx
}

func (context *RuntimeContextMock) GetInstanceContext() *wasmer.InstanceContext {
	return context.instanceContext
}

func (context *RuntimeContextMock) GetInstanceExports() wasmer.ExportsMap {
	return context.instance.Exports
}

func (context *RuntimeContextMock) CleanInstance() {
	context.instance.Clean()
	context.instance = nil
}

func (context *RuntimeContextMock) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	exports := context.instance.Exports
	function, ok := exports[context.callFunction]

	if !ok {
		function, ok = exports["main"]
	}

	if !ok {
		return nil, arwen.ErrFuncNotFound
	}

	return function, nil
}

func (context *RuntimeContextMock) GetInitFunction() wasmer.ExportedFunctionCallback {
	exports := context.instance.Exports
	init, ok := exports[arwen.InitFunctionName]

	if !ok {
		init, ok = exports[arwen.InitFunctionNameEth]
	}

	if !ok {
		init = nil
	}

	return init
}

func (context *RuntimeContextMock) MemLoad(offset int32, length int32) ([]byte, error) {
	memory := context.instanceContext.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + length)
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > memoryLength
	isRequestedEndTooLarge := requestedEnd > memoryLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("LoadBytes: bad bounds")
	}
	if isLengthNegative {
		return nil, fmt.Errorf("LoadBytes: negative length")
	}

	result := make([]byte, length)
	if isRequestedEndTooLarge {
		copy(result, memoryView[offset:])
	} else {
		copy(result, memoryView[offset:requestedEnd])
	}

	return result, nil
}

func (context *RuntimeContextMock) MemStore(offset int32, data []byte) error {
	memory := context.instanceContext.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	dataLength := int32(len(data))
	requestedEnd := uint32(offset + dataLength)
	isOffsetTooSmall := offset < 0
	isNewPageNecessary := requestedEnd > memoryLength

	if isOffsetTooSmall {
		return fmt.Errorf("StoreBytes: bad lower bounds")
	}
	if isNewPageNecessary {
		err := memory.Grow(1)
		if err != nil {
			return err
		}

		memoryView = memory.Data()
		memoryLength = memory.Length()
	}

	isRequestedEndTooLarge := requestedEnd > memoryLength
	if isRequestedEndTooLarge {
		return fmt.Errorf("StoreBytes: bad upper bounds")
	}

	copy(memoryView[offset:requestedEnd], data)
	return nil
}
