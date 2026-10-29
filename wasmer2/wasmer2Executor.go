package wasmer2

import (
	"runtime"
	"unsafe"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.Executor = (*Wasmer2Executor)(nil)

// Wasmer2Executor oversees the creation of Wasmer instances and execution.
type Wasmer2Executor struct {
	cgoExecutor *cWasmerExecutorT

	vmHookPointers *cWasmerVmHookPointers
	vmHooks        executor.VMHooksLegacy
	vmHooksPtr     uintptr
	vmHooksPtrPtr  unsafe.Pointer

	opcodeCost *OpcodeCost
}

// CreateExecutor creates a new wasmer executor.
func CreateExecutor() (*Wasmer2Executor, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	vmHookPointers := populateCgoFunctionPointers()
	localPtr := uintptr(unsafe.Pointer(vmHookPointers))
	localPtrPtr := unsafe.Pointer(&localPtr)

	var cExecutor *cWasmerExecutorT

	var result = cWasmerNewExecutor(
		&cExecutor,
		localPtrPtr,
	)

	if result != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	cWasmerForceInstallSighandlers()

	wasmerExecutor := &Wasmer2Executor{
		cgoExecutor:    cExecutor,
		vmHookPointers: vmHookPointers,
	}

	return wasmerExecutor, nil
}

// SetOpcodeConfig sets the opcode version and opcode costs, based on the gas schedule.
func (wasmerExecutor *Wasmer2Executor) SetOpcodeConfig(opcodeVersion executor.OpcodeVersion, wasmOps *executor.WASMOpcodeCost) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// extract only wasmer2 opcodes
	wasmerExecutor.opcodeCost = wasmerExecutor.extractOpcodeCost(wasmOps)
	result := cWasmerExecutorSetOpcodeConfig(
		wasmerExecutor.cgoExecutor,
		(int32)(opcodeVersion),
		(*cWasmerOpcodeCostT)(unsafe.Pointer(wasmerExecutor.opcodeCost)),
	)
	if result != cWasmerOk {
		return newWrappedError(ErrFailedInstantiation)
	}

	return nil
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (wasmerExecutor *Wasmer2Executor) SetRkyvSerializationEnabled(_ bool) {
}

// SetSIGSEGVPassthrough controls a Wasmer flag.
func (wasmerExecutor *Wasmer2Executor) SetSIGSEGVPassthrough() {
}

func (wasmerExecutor *Wasmer2Executor) FunctionNames() vmcommon.FunctionNames {
	return functionNames
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (wasmerExecutor *Wasmer2Executor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var cInstance *cWasmerInstanceT

	if len(contractCode) == 0 {
		return nil, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := toCWasmerCompilationOptions(options)
	var compileResult = cWasmerInstantiateWithOptions(
		wasmerExecutor.cgoExecutor,
		&cInstance,
		(*cUchar)(unsafe.Pointer(&contractCode[0])),
		cUint(len(contractCode)),
		&cOptions,
	)

	if compileResult != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	return newInstance(cInstance)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (wasmerExecutor *Wasmer2Executor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var cInstance *cWasmerInstanceT

	if len(compiledCode) == 0 {
		return nil, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := toCWasmerCompilationOptions(options)
	var compileResult = cWasmerInstanceFromCache(
		wasmerExecutor.cgoExecutor,
		&cInstance,
		(*cUchar)(unsafe.Pointer(&compiledCode[0])),
		cUint32T(len(compiledCode)),
		&cOptions,
	)

	if compileResult != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	return newInstance(cInstance)
}

// IsInterfaceNil returns true if underlying object is nil
func (wasmerExecutor *Wasmer2Executor) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}

// initVMHooks inits the VM hooks.
//
// The hooks are legacy-shaped because the functions this backend imports are plain C
// functions: they return a value and have nowhere to report that the hook wants execution to
// stop. The adapter the factory built applies that as a side effect instead.
//
// The context data is set once, for the executor, so every instance it makes shares these
// hooks; a hook resolves the execution it serves through the host, not through the caller.
func (wasmerExecutor *Wasmer2Executor) initVMHooks(vmHooks executor.VMHooksLegacy) {
	wasmerExecutor.vmHooks = vmHooks
	localPtr := uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
	wasmerExecutor.vmHooksPtr = localPtr
	wasmerExecutor.vmHooksPtrPtr = unsafe.Pointer(&localPtr)
	cWasmerExecutorContextDataSet(wasmerExecutor.cgoExecutor, wasmerExecutor.vmHooksPtrPtr)
}
