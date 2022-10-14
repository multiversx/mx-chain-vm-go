package wasmer

import (
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WasmerExecutor oversees the creation of Wasmer instances and execution.
type WasmerExecutor struct {
	eiFunctionNames    vmcommon.FunctionNames
	vmHooks            executor.VMHooks
	vmHooksData        uintptr
	vmHooksDataPointer unsafe.Pointer
}

// NewExecutor creates a new wasmer executor.
func NewExecutor() (*WasmerExecutor, error) {
	functionNames, err := injectCgoFunctionPointers()
	if err != nil {
		return nil, err
	}
	return &WasmerExecutor{
		eiFunctionNames: functionNames,
	}, nil
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (wasmerExecutor *WasmerExecutor) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	SetOpcodeCosts(opcodeCosts)
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (wasmerExecutor *WasmerExecutor) SetRkyvSerializationEnabled(enabled bool) {
	SetRkyvSerializationEnabled(enabled)
}

// SetSIGSEGVPassthrough controls a Wasmer flag.
func (wasmerExecutor *WasmerExecutor) SetSIGSEGVPassthrough() {
	SetSIGSEGVPassthrough()
}

func (wasmerExecutor *WasmerExecutor) FunctionNames() vmcommon.FunctionNames {
	return wasmerExecutor.eiFunctionNames
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (wasmerExecutor *WasmerExecutor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, err := NewInstanceWithOptions(contractCode, options)
	if err == nil {
		wasmerExecutor.setVMHooksPtrs(instance)
	}

	return instance, err
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (wasmerExecutor *WasmerExecutor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, err := NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err == nil {
		wasmerExecutor.setVMHooksPtrs(instance)
	}
	return instance, err
}

// SetVMHooks sets the VM hooks that will be used by the executor for current instance.
func (wasmerExecutor *WasmerExecutor) SetVMHooks(vmHooks executor.VMHooks) {
	wasmerExecutor.vmHooks = vmHooks
}

// SetVMHooksForInstance replaces the VM hooks for chosen instance with new ones.
func (wasmerExecutor *WasmerExecutor) SetVMHooksForInstance(instance executor.Instance, vmHooks executor.VMHooks) {
	wasmerExecutor.SetVMHooks(vmHooks)
	wasmerExecutor.setVMHooksPtrs(instance)
}

// GetVMHooks returns the VM hooks.
func (wasmerExecutor *WasmerExecutor) GetVMHooks() executor.VMHooks {
	return wasmerExecutor.vmHooks
}

func (wasmerExecutor *WasmerExecutor) setVMHooksPtrs(instance executor.Instance) {
	data := uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
	wasmerExecutor.vmHooksData = data
	wasmerExecutor.vmHooksDataPointer = unsafe.Pointer(&data)
	instance.SetContextData(unsafe.Pointer(&data))
}
