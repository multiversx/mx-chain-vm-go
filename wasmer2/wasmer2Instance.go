package wasmer2

// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

var _ = (executor.Instance)((*Wasmer2Instance)(nil))

// Wasmer2Instance represents a WebAssembly instance.
type Wasmer2Instance struct {
	// The underlying WebAssembly instance.
	cgoInstance *cWasmerInstanceT

	// // The exported memory of a WebAssembly instance.
	// Memory MemoryHandler

	callbacks       executor.VMHooks
	callbacksPtr    uintptr
	callbacksPtrPtr unsafe.Pointer

	// InstanceCtx InstanceContext
}

func emptyInstance() *Wasmer2Instance {
	// return &Wasmer2Instance{instance: nil, Exports: nil, Signatures: nil, Memory: nil}
	return &Wasmer2Instance{cgoInstance: nil}
}

func newInstance(c_instance *cWasmerInstanceT) (*Wasmer2Instance, error) {
	// var hasMemory bool

	// memory, hasMemory, err := retrieveExportedMemory(wasmExports)
	// if err != nil {
	// 	return emptyInstance, err
	// }

	// if !hasMemory {
	// 	return emptyInstance(), nil
	// }

	return &Wasmer2Instance{cgoInstance: c_instance}, nil
}

// // HasMemory checks whether the instance has at least one exported memory.
func (instance *Wasmer2Instance) HasMemory() bool {
	// return nil != instance.Memory
	return true
}

// func NewInstanceFromCompiledCodeWithOptions(
// 	compiledCode []byte,
// 	options CompilationOptions,
// ) (*Wasmer2Instance, error) {
// 	var c_instance *cWasmerInstanceT

// 	if len(compiledCode) == 0 {
// 		var emptyInstance = &Wasmer2Instance{instance: nil, Exports: nil, Memory: nil}
// 		return emptyInstance, newWrappedError(ErrInvalidBytecode)
// 	}

// 	cOptions := unsafe.Pointer(&options)
// 	var instantiateResult = cWasmerInstanceFromCache(
// 		&c_instance,
// 		(*cUchar)(unsafe.Pointer(&compiledCode[0])),
// 		cUint32T(len(compiledCode)),
// 		(*cWasmerCompilationOptions)(cOptions),
// 	)

// 	if instantiateResult != cWasmerOk {
// 		var emptyInstance = &Wasmer2Instance{instance: nil, Exports: nil, Memory: nil}
// 		return emptyInstance, newWrappedError(ErrFailedInstantiation)
// 	}

// 	instance, err := newInstance(c_instance)
// 	if instance != nil && instance.Memory != nil {
// 		c_instance_context := cWasmerInstanceContextGet(c_instance)
// 		instance.InstanceCtx = IntoInstanceContextDirect(c_instance_context)
// 	}

// 	return instance, err
// }

// // SetVMHooks assigns a data that can be used by all imported
// // functions. Indeed, each imported function receives as its first
// // argument an instance context (see `InstanceContext`). An instance
// // context can hold a pointer to any kind of data. It is important to
// // understand that this data is shared by all imported function, it's
// // global to the instance.
// func (instance *Wasmer2Instance) SetVMHooks(callbacks executor.VMHooks) {
// 	instance.callbacks = callbacks
// 	// This has to be a local variable, to fool Go into thinking this has nothing to do with the other structures.
// 	localPtr := uintptr(unsafe.Pointer(&instance.callbacks))
// 	instance.callbacksPtr = localPtr
// 	instance.callbacksPtrPtr = unsafe.Pointer(&localPtr)
// 	cWasmerInstanceContextDataSet(instance.cgoInstance, instance.callbacksPtrPtr)
// }

// GetVMHooks returns a pointer for the current instance's data
func (instance *Wasmer2Instance) GetVMHooks() executor.VMHooks {
	return instance.callbacks
}

func (instance *Wasmer2Instance) Clean() {
	if instance.cgoInstance != nil {
		cWasmerInstanceDestroy(instance.cgoInstance)

		// if instance.Memory != nil {
		// 	instance.Memory.Destroy()
		// }
	}
	// panic("Wasmer2Instance Clean")
}

func (instance *Wasmer2Instance) GetPointsUsed() uint64 {
	// return cWasmerInstanceGetPointsUsed(instance.instance)
	return 0
}

func (instance *Wasmer2Instance) SetPointsUsed(points uint64) {
	// cWasmerInstanceSetPointsUsed(instance.instance, points)
}

func (instance *Wasmer2Instance) SetGasLimit(gasLimit uint64) {
	// cWasmerInstanceSetGasLimit(instance.instance, gasLimit)
}

func (instance *Wasmer2Instance) SetBreakpointValue(value uint64) {
	// cWasmerInstanceSetBreakpointValue(instance.instance, value)
}

func (instance *Wasmer2Instance) GetBreakpointValue() uint64 {
	return 0
	// return cWasmerInstanceGetBreakpointValue(instance.instance)
}

func (instance *Wasmer2Instance) Cache() ([]byte, error) {
	// return []byte{}, errors.New("instance Cache not implemented")
	// var cacheBytes *cUchar
	// var cacheLen cUint32T

	// var cacheResult = cWasmerInstanceCache(
	// 	instance.instance,
	// 	&cacheBytes,
	// 	&cacheLen,
	// )

	// if cacheResult != cWasmerOk {
	// 	return nil, ErrCachingFailed
	// }

	// goBytes := C.GoBytes(unsafe.Pointer(cacheBytes), C.int(cacheLen))

	// C.free(unsafe.Pointer(cacheBytes))
	// cacheBytes = nil
	// return goBytes, nil

	return nil, ErrCachingFailed
}

// IsFunctionImported returns true if the instance imports the specified function
func (instance *Wasmer2Instance) IsFunctionImported(name string) bool {
	// return cWasmerInstanceIsFunctionImported(instance.instance, name)
	return false
}

func (instance *Wasmer2Instance) CallFunction(functionName string) error {
	var wasmFunctionName = cCString(functionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	var callResult = cWasmerInstanceCall(
		instance.cgoInstance,
		wasmFunctionName,
	)

	if callResult != cWasmerOk {
		err := fmt.Errorf("failed to call the `%s` exported function", functionName)
		return newWrappedError(err)
	}

	return nil
}

func (instance *Wasmer2Instance) HasFunction(functionName string) bool {
	var wasmFunctionName = cCString(functionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	result := cWasmerInstanceHasFunction(
		instance.cgoInstance,
		wasmFunctionName,
	)

	return result == 1
}

// GetLastError returns the last error message if any, otherwise returns an error.
func (instance *Wasmer2Instance) getFunctionNamesConcat() (string, error) {
	var bufferLength = cWasmerInstanceExportedFunctionNamesLength(instance.cgoInstance)

	if bufferLength == 0 {
		return "", nil
	}

	var buffer = make([]cChar, bufferLength)
	var bufferPointer = (*cChar)(unsafe.Pointer(&buffer[0]))

	var result = cWasmerInstanceExportedFunctionNames(instance.cgoInstance, bufferPointer, bufferLength)

	if result == -1 {
		return "", errors.New("cannot read function names")
	}

	return cGoString(bufferPointer), nil
}

func (instance *Wasmer2Instance) GetFunctionNames() []string {
	buffer, err := instance.getFunctionNamesConcat()
	if err != nil {
		return nil
	}
	return strings.Split(buffer, "|")
}

func (instance *Wasmer2Instance) ValidateVoidFunction(functionName string) error {
	// return instance.verifyVoidFunction(functionName)
	return nil
}

// GetData returns a pointer for the current instance's data
func (instance *Wasmer2Instance) GetData() uintptr {
	panic("instance GetData")
}

// GetInstanceCtxMemory returns the memory for the instance context
func (instance *Wasmer2Instance) GetInstanceCtxMemory() executor.Memory {
	// return instance.InstanceCtx.Memory()
	panic("instance GetInstanceCtxMemory")
}

// GetMemory returns the memory for the instance
func (instance *Wasmer2Instance) GetMemory() executor.Memory {
	// return instance.Memory
	panic("instance GetMemory")
}

// SetMemory sets the memory for the instance returns true if success
func (instance *Wasmer2Instance) SetMemory(cleanMemory []byte) bool {
	// instanceMemory := instance.GetMemory().Data()
	// if len(instanceMemory) != len(cleanMemory) {
	// 	// TODO shrink the instance memory instead and return true
	// 	return false
	// }

	// copy(instanceMemory, cleanMemory)
	// return true
	panic("instance SetMemory")
}

// Reset resets the instance memories and globals
func (instance *Wasmer2Instance) Reset() bool {
	// result := cWasmerInstanceReset(instance.instance)
	// return result == cWasmerOk
	return true
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *Wasmer2Instance) IsInterfaceNil() bool {
	return instance == nil
}

func (instance *Wasmer2Instance) SetVMHooksPtr(vmHooksPtr uintptr) {
}

func (instance *Wasmer2Instance) GetVMHooksPtr() uintptr {
	return uintptr(0)
}
