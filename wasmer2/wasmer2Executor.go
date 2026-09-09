package wasmer2

import (
	"unsafe"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.Executor = (*Wasmer2Executor)(nil)

// Wasmer2Executor oversees the creation of Wasmer instances and execution.
type Wasmer2Executor struct {
	cgoExecutor *cWasmerExecutorT

	vmHookPointers *cWasmerVmHookPointers
	vmHooks        executor.VMHooks
	vmHooksPtr     uintptr
	vmHooksPtrPtr  unsafe.Pointer

	opcodeCost *OpcodeCost
}

// CreateExecutor creates a new wasmer executor.
func CreateExecutor() (*Wasmer2Executor, error) {
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
func (wasmerExecutor *Wasmer2Executor) SetOpcodeConfig(opcodeVersion executor.OpcodeVersion, wasmOps *executor.WASMOpcodeCost) {
	// extract only wasmer2 opcodes
	wasmerExecutor.opcodeCost = wasmerExecutor.extractOpcodeCost(wasmOps)
	cWasmerExecutorSetOpcodeConfig(
		wasmerExecutor.cgoExecutor,
		(int32)(opcodeVersion),
		(*cWasmerOpcodeCostT)(unsafe.Pointer(wasmerExecutor.opcodeCost)),
	)
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
	var cInstance *cWasmerInstanceT

	if len(contractCode) == 0 {
		return nil, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var compileResult = cWasmerInstantiateWithOptions(
		wasmerExecutor.cgoExecutor,
		&cInstance,
		(*cUchar)(unsafe.Pointer(&contractCode[0])),
		cUint(len(contractCode)),
		(*cWasmerCompilationOptions)(cOptions),
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
	var cInstance *cWasmerInstanceT

	if len(compiledCode) == 0 {
		return nil, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var compileResult = cWasmerInstanceFromCache(
		wasmerExecutor.cgoExecutor,
		&cInstance,
		(*cUchar)(unsafe.Pointer(&compiledCode[0])),
		cUint32T(len(compiledCode)),
		(*cWasmerCompilationOptions)(cOptions),
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

// InitVMHooks inits the VM hooks
func (wasmerExecutor *Wasmer2Executor) initVMHooks(vmHooks executor.VMHooks) {
	wasmerExecutor.vmHooks = vmHooks
	localPtr := uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
	wasmerExecutor.vmHooksPtr = localPtr
	wasmerExecutor.vmHooksPtrPtr = unsafe.Pointer(&localPtr)
	cWasmerExecutorContextDataSet(wasmerExecutor.cgoExecutor, wasmerExecutor.vmHooksPtrPtr)
}
