package wasmer

import (
	"unsafe"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.Executor = (*WasmerExecutor)(nil)

// WasmerExecutor oversees the creation of Wasmer instances and execution.
type WasmerExecutor struct {
	eiFunctionNames vmcommon.FunctionNames
	vmHooks         executor.VMHooks
	vmHooksPtr      uintptr
}

// CreateExecutor creates a new wasmer executor.
func CreateExecutor() (*WasmerExecutor, error) {
	functionNames, err := injectCgoFunctionPointers()
	if err != nil {
		return nil, err
	}

	ForceInstallSighandlers()

	return &WasmerExecutor{
		eiFunctionNames: functionNames,
	}, nil
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (wasmerExecutor *WasmerExecutor) SetOpcodeCosts(opcodeCosts *executor.WASMOpcodeCost) {
	SetOpcodeCosts(opcodeCosts)
}

// HasFunctionNameChecks returns true if the instance requires function name checks.
func (wasmerExecutor *WasmerExecutor) HasFunctionNameChecks() bool {
	return true
}

// FunctionNames returns the function names
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

// initVMHooks inits the VM hooks
func (wasmerExecutor *WasmerExecutor) initVMHooks(vmHooks executor.VMHooks) {
	wasmerExecutor.vmHooks = vmHooks
	wasmerExecutor.vmHooksPtr = uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
}

// IsInterfaceNil returns true if there is no value under the interface
func (wasmerExecutor *WasmerExecutor) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}
