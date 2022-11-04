package wasmer

import (
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/config"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

var _ = (executor.Executor)((*WasmerExecutor)(nil))

// WasmerExecutor oversees the creation of Wasmer instances and execution.
type WasmerExecutor struct {
	eiFunctionNames vmcommon.FunctionNames
	vmHooks         executor.VMHooks
	vmHooksPtr      uintptr
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
func (wasmerExecutor *WasmerExecutor) SetOpcodeCosts(opcodeCosts *config.WASMOpcodeCost) {
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
		instance.SetVMHooksPtr(wasmerExecutor.vmHooksPtr)
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
		instance.SetVMHooksPtr(wasmerExecutor.vmHooksPtr)
	}
	return instance, err
}

// InitVMHooks inits the VM hooks
func (wasmerExecutor *WasmerExecutor) InitVMHooks(vmHooks executor.VMHooks) {
	wasmerExecutor.vmHooks = vmHooks
	wasmerExecutor.vmHooksPtr = uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
}

// GetVMHooks returns the VM hooks
func (wasmerExecutor *WasmerExecutor) GetVMHooks() executor.VMHooks {
	return wasmerExecutor.vmHooks
}
