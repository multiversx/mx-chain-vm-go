package subcontexts

import (
	"unsafe"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Runtime struct {
	blockChainHook  vmcommon.BlockchainHook
	instance        *wasmer.Instance
	instanceContext *wasmer.InstanceContext
	vmInput         *vmcommon.VMInput
	scAddress       []byte
	callFunction    string
	vmType          []byte
}

func NewRuntimeSubcontext(blockChainHook vmcommon.BlockchainHook) (*Runtime, error) {
	runtime := &Runtime{
		blockChainHook:  blockChainHook,
		instance:        nil,
		instanceContext: nil,
		vmInput:         nil,
	}
	return runtime, nil
}

func (runtime *Runtime) InitializeFromInput(input *vmcommon.ContractCallInput) error {
	runtime.vmInput = &input.VMInput
	runtime.scAddress = input.RecipientAddr
	runtime.callFunction = input.Function

	contract, err := runtime.blockChainHook.GetCode(runtime.scAddress)
	gasProvided := runtime.vmInput.GasProvided
	instance, err := wasmer.NewMeteredInstance(contract, gasProvided)
	if err != nil {
		return err
	}
	runtime.instance = instance
	return nil
}

func (runtime *Runtime) GetVMInput() *vmcommon.VMInput {
	return runtime.vmInput
}

func (runtime *Runtime) GetSCAddress() []byte {
	panic("not implemented")
}

func (runtime *Runtime) Function() string {
	panic("not implemented")
}

func (runtime *Runtime) Arguments() [][]byte {
	panic("not implemented")
}

func (runtime *Runtime) SignalUserError() {
	panic("not implemented")
}

func (runtime *Runtime) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	panic("not implemented")
}

func (runtime *Runtime) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	panic("not implemented")
}

func (runtime *Runtime) CallData() []byte {
	panic("not implemented")
}

func (runtime *Runtime) SetReadOnly(readOnly bool) {
	panic("not implemented")
}

func (runtime *Runtime) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	panic("not implemented")
}

func (runtime *Runtime) ExecuteOnDestContext(input *vmcommon.ContractCallInput) error {
	panic("not implemented")
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

func (runtime *Runtime) MemStore(offset int32, data []byte) error {
	memory := runtime.instanceContext.Memory()
	return arwen.StoreBytes(memory, offset, data)
}

func (runtime *Runtime) MemLoad(offset int32, length int32) ([]byte, error) {
	memory := runtime.instanceContext.Memory()
	return arwen.LoadBytes(memory, offset, length)
}

func (runtime *Runtime) Clean() {
	runtime.instance.Clean()
}
