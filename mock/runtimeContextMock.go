package mock

import (
	"fmt"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type RuntimeContextMock struct {
	Err error
}

func (r *RuntimeContextMock) InitState() {
}

func (r *RuntimeContextMock) CreateWasmerInstance(contract []byte, gasLimit uint64) error {
	if Err != nil {
		return 
	return nil
}

func (r *RuntimeContextMock) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
}

func (r *RuntimeContextMock) PushState() {
}

func (r *RuntimeContextMock) PushInstance() {
}

func (r *RuntimeContextMock) PopInstance() error {
	return nil
}

func (r *RuntimeContextMock) PopState() error {
	stateStackLen := len(r.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

	prevState := r.stateStack[stateStackLen-1]
	r.stateStack = r.stateStack[:stateStackLen-1]

	r.vmInput = prevState.vmInput
	r.scAddress = prevState.scAddress
	r.callFunction = prevState.callFunction
	r.readOnly = prevState.readOnly

	return nil
}

func (r *RuntimeContextMock) GetVMType() []byte {
	return []byte("type")
}

func (r *RuntimeContextMock) GetVMInput() *vmcommon.VMInput {
	return r.vmInput
}

func (r *RuntimeContextMock) SetVMInput(vmInput *vmcommon.VMInput) {
	r.vmInput = vmInput
}

func (r *RuntimeContextMock) GetSCAddress() []byte {
	return r.scAddress
}

func (r *RuntimeContextMock) SetSCAddress(scAddress []byte) {
	r.scAddress = scAddress
}

func (r *RuntimeContextMock) Function() string {
	return r.callFunction
}

func (r *RuntimeContextMock) Arguments() [][]byte {
	return r.vmInput.Arguments
}

func (r *RuntimeContextMock) SignalExit(_ int) {
	r.SetRuntimeBreakpointValue(arwen.BreakpointSignalExit)
}

func (r *RuntimeContextMock) SignalUserError(_ string) {
}

func (r *RuntimeContextMock) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
}

func (r *RuntimeContextMock) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(r.instance.GetBreakpointValue())
}

func (r *RuntimeContextMock) GetPointsUsed() uint64 {
	return r.pointsUsed
}

func (r *RuntimeContextMock) SetPointsUsed(gasPoints uint64) {
	r.pointsUsed = gasPoints
}

func (r *RuntimeContextMock) ReadOnly() bool {
	return r.readOnly
}

func (r *RuntimeContextMock) SetReadOnly(readOnly bool) {
	r.readOnly = readOnly
}

func (r *RuntimeContextMock) SetInstanceContextId(id int) {
	r.instance.SetContextData(unsafe.Pointer(&id))
}

func (r *RuntimeContextMock) SetInstanceContext(instCtx *wasmer.InstanceContext) {
	r.instanceContext = instCtx
}

func (r *RuntimeContextMock) GetInstanceContext() *wasmer.InstanceContext {
	return r.instanceContext
}

func (r *RuntimeContextMock) GetInstanceExports() wasmer.ExportsMap {
	return r.instance.Exports
}

func (r *RuntimeContextMock) CleanInstance() {
	r.instance.Clean()
	r.instance = nil
}

func (r *RuntimeContextMock) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	exports := r.instance.Exports
	function, ok := exports[r.callFunction]

	if !ok {
		function, ok = exports["main"]
	}

	if !ok {
		return nil, arwen.ErrFuncNotFound
	}

	return function, nil
}

func (r *RuntimeContextMock) GetInitFunction() wasmer.ExportedFunctionCallback {
	exports := r.instance.Exports
	init, ok := exports[arwen.InitFunctionName]

	if !ok {
		init, ok = exports[arwen.InitFunctionNameEth]
	}

	if !ok {
		init = nil
	}

	return init
}

func (r *RuntimeContextMock) MemLoad(offset int32, length int32) ([]byte, error) {
	memory := r.instanceContext.Memory()
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

func (r *RuntimeContextMock) MemStore(offset int32, data []byte) error {
	memory := r.instanceContext.Memory()
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
