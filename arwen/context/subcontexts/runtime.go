package subcontexts

import (
	"strconv"
	"unsafe"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Runtime struct {
	host            arwen.VMContext
	instance        *wasmer.Instance
	instanceContext *wasmer.InstanceContext
	vmInput         *vmcommon.VMInput
	scAddress       []byte
	callFunction    string
	vmType          []byte
	readOnly        bool

	stateStack    []*Runtime
	instanceStack []*wasmer.Instance
}

func NewRuntimeSubcontext(
	host arwen.VMContext,
	blockChainHook vmcommon.BlockchainHook,
	vmType []byte,
) (*Runtime, error) {
	runtime := &Runtime{
		host:          host,
		vmType:        vmType,
		stateStack:    make([]*Runtime, 0),
		instanceStack: make([]*wasmer.Instance, 0),
	}

	runtime.InitState()

	return runtime, nil
}

func (runtime *Runtime) InitState() {
	runtime.vmInput = &vmcommon.VMInput{}
	runtime.scAddress = make([]byte, 0)
	runtime.callFunction = ""
	runtime.readOnly = false
}

func (runtime *Runtime) CreateWasmerInstance(contract []byte) error {
	if runtime.instance != nil {
		runtime.CleanInstance()
	}

	var err error
	runtime.instance, err = wasmer.NewMeteredInstance(contract, runtime.vmInput.GasProvided)
	if err != nil {
		return err
	}
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	return nil
}

func (runtime *Runtime) CreateWasmerInstanceWithGasLimit(contract []byte, gasLimit uint64) error {
	var err error
	runtime.instance, err = wasmer.NewMeteredInstance(contract, gasLimit)
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	return err
}

func (runtime *Runtime) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	runtime.vmInput = &input.VMInput
	runtime.scAddress = input.RecipientAddr
	runtime.callFunction = input.Function
}

func (runtime *Runtime) PushState() {
	newState := &Runtime{
		vmInput:      runtime.vmInput,
		scAddress:    runtime.scAddress,
		callFunction: runtime.callFunction,
		readOnly:     runtime.readOnly,
	}

	runtime.stateStack = append(runtime.stateStack, newState)
}

func (runtime *Runtime) PushInstance() {
	runtime.instanceStack = append(runtime.instanceStack, runtime.instance)
}

func (runtime *Runtime) PopInstance() error {
	instanceStackLen := len(runtime.instanceStack)
	if instanceStackLen < 1 {
		return InstanceStackUnderflow
	}

	prevInstance := runtime.instanceStack[instanceStackLen-1]
	runtime.instanceStack = runtime.instanceStack[:instanceStackLen-1]

	runtime.instance.Clean()
	runtime.instance = prevInstance

	return nil
}

func (runtime *Runtime) PopState() error {
	stateStackLen := len(runtime.stateStack)
	if stateStackLen < 1 {
		return StateStackUnderflow
	}

	prevState := runtime.stateStack[stateStackLen-1]
	runtime.stateStack = runtime.stateStack[:stateStackLen-1]

	runtime.vmInput = prevState.vmInput
	runtime.scAddress = prevState.scAddress
	runtime.callFunction = prevState.callFunction
	runtime.readOnly = prevState.readOnly

	return nil
}

func (runtime *Runtime) GetVMType() []byte {
	return runtime.vmType
}

func (runtime *Runtime) GetVMInput() *vmcommon.VMInput {
	return runtime.vmInput
}

func (runtime *Runtime) SetVMInput(vmInput *vmcommon.VMInput) {
	runtime.vmInput = vmInput
}

func (runtime *Runtime) GetSCAddress() []byte {
	return runtime.scAddress
}

func (runtime *Runtime) SetSCAddress(scAddress []byte) {
	runtime.scAddress = scAddress
}

func (runtime *Runtime) Function() string {
	return runtime.callFunction
}

func (runtime *Runtime) Arguments() [][]byte {
	return runtime.vmInput.Arguments
}

func (runtime *Runtime) SignalExit(exitCode int) {
	runtime.host.Output().SetReturnCode(vmcommon.Ok)
	message := strconv.Itoa(exitCode)
	runtime.host.Output().SetReturnMessage(message)
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointSignalExit)
}

func (runtime *Runtime) SignalUserError(message string) {
	// SignalUserError() remains in Runtime, and won't be moved into Output,
	// because there will be extra handling added here later, which requires
	// information from Runtime (e.g. runtime breakpoints)
	runtime.host.Output().SetReturnCode(vmcommon.UserError)
	runtime.host.Output().SetReturnMessage(message)
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointSignalError)
}

func (runtime *Runtime) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	runtime.instance.SetBreakpointValue(uint64(value))
}

func (runtime *Runtime) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(runtime.instance.GetBreakpointValue())
}

func (runtime *Runtime) GetPointsUsed() uint64 {
	return runtime.instance.GetPointsUsed()
}

func (runtime *Runtime) SetPointsUsed(gasPoints uint64) {
	runtime.instance.SetPointsUsed(gasPoints)
}

func (runtime *Runtime) CallData() []byte {
	panic("not implemented")
}

func (runtime *Runtime) ReadOnly() bool {
	return runtime.readOnly
}

func (runtime *Runtime) SetReadOnly(readOnly bool) {
	runtime.readOnly = readOnly
}

func (runtime *Runtime) SetInstanceContextId(id int) {
	runtime.instance.SetContextData(unsafe.Pointer(&id))
}

func (runtime *Runtime) SetInstanceContext(instCtx *wasmer.InstanceContext) {
	runtime.instanceContext = instCtx
}

func (runtime *Runtime) GetInstanceContext() *wasmer.InstanceContext {
	return runtime.instanceContext
}

func (runtime *Runtime) GetInstanceExports() wasmer.ExportsMap {
	return runtime.instance.Exports
}

func (runtime *Runtime) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	exports := runtime.instance.Exports
	function, ok := exports[runtime.callFunction]

	if !ok {
		function, ok = exports["main"]
	}

	if !ok {
		return nil, ErrFuncNotFound
	}

	return function, nil
}

func (runtime *Runtime) GetInitFunction() wasmer.ExportedFunctionCallback {
	exports := runtime.instance.Exports
	init, ok := exports[arwen.InitFunctionName]

	if !ok {
		init, ok = exports[arwen.InitFunctionNameEth]
	}

	if !ok {
		init = nil
	}

	return init
}

func (runtime *Runtime) MemStore(offset int32, data []byte) error {
	memory := runtime.instanceContext.Memory()
	return arwen.StoreBytes(memory, offset, data)
}

func (runtime *Runtime) MemLoad(offset int32, length int32) ([]byte, error) {
	memory := runtime.instanceContext.Memory()
	return arwen.LoadBytes(memory, offset, length)
}

func (runtime *Runtime) CleanInstance() {
	runtime.instance.Clean()
	runtime.instance = nil
}
