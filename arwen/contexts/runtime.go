package contexts

import (
	"fmt"
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

	asyncCallInfo *arwen.AsyncCallInfo

	argParser arwen.ArgumentsParser
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

func (runtime *runtimeContext) InitState() {
	runtime.vmInput = &vmcommon.VMInput{}
	runtime.scAddress = make([]byte, 0)
	runtime.callFunction = ""
	runtime.readOnly = false
	runtime.argParser = vmcommon.NewAtArgumentParser()
	runtime.asyncCallInfo = nil
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
		vmInput:       context.vmInput,
		scAddress:     context.scAddress,
		callFunction:  context.callFunction,
		readOnly:      context.readOnly,
		asyncCallInfo: context.asyncCallInfo,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (runtime *runtimeContext) PopState() {
	stateStackLen := len(runtime.stateStack)

	prevState := runtime.stateStack[stateStackLen-1]
	runtime.stateStack = runtime.stateStack[:stateStackLen-1]

	runtime.vmInput = prevState.vmInput
	runtime.scAddress = prevState.scAddress
	runtime.callFunction = prevState.callFunction
	runtime.readOnly = prevState.readOnly
	runtime.asyncCallInfo = prevState.asyncCallInfo
}

func (runtime *runtimeContext) PushInstance() {
	runtime.instanceStack = append(runtime.instanceStack, runtime.instance)
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

func (runtime *runtimeContext) ArgParser() arwen.ArgumentsParser {
	return runtime.argParser
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

func (context *runtimeContext) FailExecution(err error) {
	context.host.Output().SetReturnCode(vmcommon.ExecutionFailed)
	if err != nil {
		context.host.Output().SetReturnMessage(err.Error())
	}
	context.SetRuntimeBreakpointValue(arwen.BreakpointExecutionFailed)
}

func (context *runtimeContext) SignalUserError(message string) {
	// SignalUserError() remains in runtimeContext, and won't be moved into Output,
	// because there will be extra handling added here later, which requires
	// information from runtimeContext (e.g. runtime breakpoints)
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

func (context *runtimeContext) ElrondAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) BigIntAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) CryptoAPIErrorShouldFailExecution() bool {
	return true
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

func (runtime *runtimeContext) SetAsyncCallInfo(asyncCallInfo *arwen.AsyncCallInfo) {
	runtime.asyncCallInfo = asyncCallInfo
}

func (runtime *runtimeContext) GetAsyncCallInfo() *arwen.AsyncCallInfo {
	return runtime.asyncCallInfo
}

func (runtime *runtimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
	memory := runtime.instanceContext.Memory()
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
