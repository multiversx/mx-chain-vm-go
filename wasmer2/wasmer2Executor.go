package wasmer2

import (
	"errors"
	"fmt"
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

var _ = (executor.Executor)((*Wasmer2Executor)(nil))

// WasmerExecutor oversees the creation of Wasmer instances and execution.
type Wasmer2Executor struct {
	cgoExecutor *cWasmerExecutorT

	vmHookPointers *cWasmerVmHookPointers
	vmHooks        executor.VMHooks
	vmHooksPtr     uintptr
	vmHooksPtrPtr  unsafe.Pointer
}

// NewExecutor creates a new wasmer executor.
func NewExecutor() (*Wasmer2Executor, error) {
	vmHookPointers := populateCgoFunctionPointers()
	fmt.Printf("\nget_gas_left_func_ptr %x\n", uintptr(unsafe.Pointer(vmHookPointers.get_gas_left_func_ptr)))
	fmt.Printf("get_sc_address_func_ptr %x\n", uintptr(unsafe.Pointer(vmHookPointers.get_sc_address_func_ptr)))
	fmt.Printf("check_no_payment_func_ptr %x\n", uintptr(unsafe.Pointer(vmHookPointers.check_no_payment_func_ptr)))
	localPtr := uintptr(unsafe.Pointer(vmHookPointers))
	fmt.Printf("localPtr %x\n", localPtr)
	localPtrPtr := unsafe.Pointer(&localPtr)

	// wasmerExecutor.vmHooks = vmHooks
	// wasmerExecutor.vmHooksPtr = localPtr
	// wasmerExecutor.vmHooksPtrPtr = unsafe.Pointer(&localPtr)
	// cWasmerExecutorContextDataSet(wasmerExecutor.cgoExecutor, localPtrPtr)

	var c_executor *cWasmerExecutorT

	var result = cWasmerNewExecutor(
		&c_executor,
		localPtrPtr,
	)

	if result != cWasmerOk {
		var emptyInstance = &Wasmer2Executor{cgoExecutor: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	executor := &Wasmer2Executor{
		cgoExecutor:    c_executor,
		vmHookPointers: vmHookPointers,
	}

	return executor, nil
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (wasmerExecutor *Wasmer2Executor) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	// SetOpcodeCosts(opcodeCosts)
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (wasmerExecutor *Wasmer2Executor) SetRkyvSerializationEnabled(enabled bool) {
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
	var c_instance *cWasmerInstanceT

	if len(contractCode) == 0 {
		var emptyInstance = &Wasmer2Instance{cgoInstance: nil}
		return emptyInstance, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var compileResult = cWasmerInstantiateWithOptions(
		wasmerExecutor.cgoExecutor,
		&c_instance,
		(*cUchar)(unsafe.Pointer(&contractCode[0])),
		cUint(len(contractCode)),
		(*cWasmerCompilationOptions)(cOptions),
	)

	if compileResult != cWasmerOk {
		var emptyInstance = &Wasmer2Instance{cgoInstance: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	instance, err := newInstance(c_instance)
	// if instance != nil && instance.Memory != nil {
	// 	c_instance_context := cWasmerInstanceContextGet(c_instance)
	// 	instance.InstanceCtx = IntoInstanceContextDirect(c_instance_context)
	// }

	executionInfo, _ := GetExecutionInfo()
	fmt.Println(executionInfo)

	return instance, err
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (wasmerExecutor *Wasmer2Executor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	// return NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	return nil, errors.New("NewInstanceFromCompiledCodeWithOptions not implemented")
}

// InitVMHooks inits the VM hooks
func (wasmerExecutor *Wasmer2Executor) InitVMHooks(vmHooks executor.VMHooks) {
	wasmerExecutor.vmHooks = vmHooks
	localPtr := uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
	wasmerExecutor.vmHooksPtr = localPtr
	wasmerExecutor.vmHooksPtrPtr = unsafe.Pointer(&localPtr)
	fmt.Printf("InitVMHooks vmHooksPtr %x\n", wasmerExecutor.vmHooksPtr)
	fmt.Printf("InitVMHooks vmHooksPtrPtr %x\n", wasmerExecutor.vmHooksPtrPtr)
	cWasmerExecutorContextDataSet(wasmerExecutor.cgoExecutor, wasmerExecutor.vmHooksPtrPtr)
}

// GetVMHooks returns the VM hooks
func (wasmerExecutor *Wasmer2Executor) GetVMHooks() executor.VMHooks {
	return wasmerExecutor.vmHooks
}