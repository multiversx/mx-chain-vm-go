package contexts

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type runtimeContext struct {
	host            arwen.VMHost
	instance        *wasmer.Instance
	instanceContext *wasmer.InstanceContext
	vmInput         *vmcommon.VMInput
	scAddress       []byte
	callFunction    string
	vmType          []byte
	readOnly        bool

	stateStack    []*runtimeContext
	instanceStack []*wasmer.Instance
}

func NewRuntimeContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
	vmType []byte,
) (*runtimeContext, error) {
	context := &runtimeContext{
		host:          host,
		vmType:        vmType,
		stateStack:    make([]*runtimeContext, 0),
		instanceStack: make([]*wasmer.Instance, 0),
	}

	context.InitState()

	return context, nil
}

func (context *runtimeContext) InitState() {
	context.vmInput = &vmcommon.VMInput{}
	context.scAddress = make([]byte, 0)
	context.callFunction = ""
	context.readOnly = false
}

func (context *runtimeContext) CreateWasmerInstance(contract []byte, gasLimit uint64) error {
	var err error
	context.instance, err = wasmer.NewMeteredInstance(contract, gasLimit)
	if err != nil {
		return err
	}
	context.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	return nil
}

func (context *runtimeContext) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	context.vmInput = &input.VMInput
	context.scAddress = input.RecipientAddr
	context.callFunction = input.Function
}

func (context *runtimeContext) PushState() {
	newState := &runtimeContext{
		vmInput:      context.vmInput,
		scAddress:    context.scAddress,
		callFunction: context.callFunction,
		readOnly:     context.readOnly,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *runtimeContext) PushInstance() {
	context.instanceStack = append(context.instanceStack, context.instance)
}

func (context *runtimeContext) PopInstance() error {
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

func (context *runtimeContext) PopState() error {
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

func (context *runtimeContext) GetVMType() []byte {
	return context.vmType
}

func (context *runtimeContext) GetVMInput() *vmcommon.VMInput {
	return context.vmInput
}

func (context *runtimeContext) SetVMInput(vmInput *vmcommon.VMInput) {
	context.vmInput = vmInput
}

func (context *runtimeContext) GetSCAddress() []byte {
	return context.scAddress
}

func (context *runtimeContext) SetSCAddress(scAddress []byte) {
	context.scAddress = scAddress
}

func (context *runtimeContext) Function() string {
	return context.callFunction
}

func (context *runtimeContext) Arguments() [][]byte {
	return context.vmInput.Arguments
}

func (context *runtimeContext) SignalExit(exitCode int) {
	context.host.Output().SetReturnCode(vmcommon.Ok)
	message := strconv.Itoa(exitCode)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointSignalExit)
}

func (context *runtimeContext) SignalUserError(message string) {
	// SignalUserError() remains in Runtime, and won't be moved into Output,
	// because there will be extra handling added here later, which requires
	// information from Runtime (e.g. runtime breakpoints)
	context.host.Output().SetReturnCode(vmcommon.UserError)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointSignalError)
}

func (context *runtimeContext) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	context.instance.SetBreakpointValue(uint64(value))
}

func (context *runtimeContext) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(context.instance.GetBreakpointValue())
}

func (context *runtimeContext) GetPointsUsed() uint64 {
	return context.instance.GetPointsUsed()
}

func (context *runtimeContext) SetPointsUsed(gasPoints uint64) {
	context.instance.SetPointsUsed(gasPoints)
}

func (context *runtimeContext) ReadOnly() bool {
	return context.readOnly
}

func (context *runtimeContext) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

func (context *runtimeContext) SetInstanceContextId(id int) {
	context.instance.SetContextData(unsafe.Pointer(&id))
}

func (context *runtimeContext) SetInstanceContext(instCtx *wasmer.InstanceContext) {
	context.instanceContext = instCtx
}

func (context *runtimeContext) GetInstanceContext() *wasmer.InstanceContext {
	return context.instanceContext
}

func (context *runtimeContext) GetInstanceExports() wasmer.ExportsMap {
	return context.instance.Exports
}

func (context *runtimeContext) CleanInstance() {
	context.instance.Clean()
	context.instance = nil
}

func (context *runtimeContext) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
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

func (context *runtimeContext) GetInitFunction() wasmer.ExportedFunctionCallback {
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

func (context *runtimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
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

func (context *runtimeContext) MemStore(offset int32, data []byte) error {
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
