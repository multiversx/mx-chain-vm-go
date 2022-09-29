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

// // InstanceError represents any kind of errors related to a WebAssembly instance. It
// // is returned by `Instance` functions only.
// type InstanceError struct {
// 	// Error message.
// 	message string
// }

// // NewInstanceError constructs a new `InstanceError`.
// func NewInstanceError(message string) *InstanceError {
// 	return &InstanceError{message}
// }

// // `InstanceError` is an actual error. The `Error` function returns
// // the error message.
// func (error *InstanceError) Error() string {
// 	return error.message
// }

// // NewExportedFunctionError constructs a new `ExportedFunctionError`,
// // where `functionName` is the name of the exported function, and
// // `message` is the error message. If the error message contains `%s`,
// // then this parameter will be replaced by `functionName`.
// func NewExportedFunctionError(functionName string, message string) *ExportedFunctionError {
// 	return &ExportedFunctionError{functionName, message}
// }

// // ExportedFunctionError is an actual error. The `Error` function
// // returns the error message.
// func (error *ExportedFunctionError) Error() string {
// 	return error.message
// }

// type ExportedFunctionCallback func(...interface{}) (Value, error)
// type ExportsMap map[string]ExportedFunctionCallback
// type ExportSignaturesMap map[string]*ExportedFunctionSignature

// Instance represents a WebAssembly instance.
type Instance struct {
	// The underlying WebAssembly instance.
	instance *cWasmerInstanceT

	// // The exported memory of a WebAssembly instance.
	// Memory MemoryHandler

	// Data        *uintptr
	// DataPointer unsafe.Pointer

	// InstanceCtx InstanceContext
}

func emptyInstance() *Instance {
	// return &Instance{instance: nil, Exports: nil, Signatures: nil, Memory: nil}
	return &Instance{instance: nil}
}

func NewInstanceWithOptions(
	bytes []byte,
	options executor.CompilationOptions,
) (*Instance, error) {
	var c_instance *cWasmerInstanceT

	if len(bytes) == 0 {
		var emptyInstance = &Instance{instance: nil}
		return emptyInstance, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var compileResult = cWasmerInstantiateWithOptions(
		&c_instance,
		(*cUchar)(unsafe.Pointer(&bytes[0])),
		cUint(len(bytes)),
		(*cWasmerCompilationOptions)(cOptions),
	)

	if compileResult != cWasmerOk {
		var emptyInstance = &Instance{instance: nil}
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

func newInstance(c_instance *cWasmerInstanceT) (*Instance, error) {
	// var hasMemory bool

	// memory, hasMemory, err := retrieveExportedMemory(wasmExports)
	// if err != nil {
	// 	return emptyInstance, err
	// }

	// if !hasMemory {
	// 	return emptyInstance(), nil
	// }

	return &Instance{instance: c_instance}, nil
}

// // HasMemory checks whether the instance has at least one exported memory.
func (instance *Instance) HasMemory() bool {
	// return nil != instance.Memory
	return true
}

// func NewInstanceFromCompiledCodeWithOptions(
// 	compiledCode []byte,
// 	options CompilationOptions,
// ) (*Instance, error) {
// 	var c_instance *cWasmerInstanceT

// 	if len(compiledCode) == 0 {
// 		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
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
// 		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
// 		return emptyInstance, newWrappedError(ErrFailedInstantiation)
// 	}

// 	instance, err := newInstance(c_instance)
// 	if instance != nil && instance.Memory != nil {
// 		c_instance_context := cWasmerInstanceContextGet(c_instance)
// 		instance.InstanceCtx = IntoInstanceContextDirect(c_instance_context)
// 	}

// 	return instance, err
// }

// SetContextData assigns a data that can be used by all imported
// functions. Indeed, each imported function receives as its first
// argument an instance context (see `InstanceContext`). An instance
// context can hold a pointer to any kind of data. It is important to
// understand that this data is shared by all imported function, it's
// global to the instance.
func (instance *Instance) SetContextData(data uintptr) {
	// instance.Data = &data
	// instance.DataPointer = unsafe.Pointer(instance.Data)
	// cWasmerInstanceContextDataSet(instance.instance, instance.DataPointer)
	// panic("Instance SetContextData")
}

func (instance *Instance) Clean() {
	if instance.instance != nil {
		cWasmerInstanceDestroy(instance.instance)

		// if instance.Memory != nil {
		// 	instance.Memory.Destroy()
		// }
	}
	// panic("Instance Clean")
}

func (instance *Instance) GetPointsUsed() uint64 {
	// return cWasmerInstanceGetPointsUsed(instance.instance)
	return 0
}

func (instance *Instance) SetPointsUsed(points uint64) {
	// cWasmerInstanceSetPointsUsed(instance.instance, points)
}

func (instance *Instance) SetGasLimit(gasLimit uint64) {
	// cWasmerInstanceSetGasLimit(instance.instance, gasLimit)
}

func (instance *Instance) SetBreakpointValue(value uint64) {
	// cWasmerInstanceSetBreakpointValue(instance.instance, value)
}

func (instance *Instance) GetBreakpointValue() uint64 {
	return 0
	// return cWasmerInstanceGetBreakpointValue(instance.instance)
}

func (instance *Instance) Cache() ([]byte, error) {
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
func (instance *Instance) IsFunctionImported(name string) bool {
	// return cWasmerInstanceIsFunctionImported(instance.instance, name)
	return false
}

func (instance *Instance) CallFunction(functionName string) error {
	var wasmFunctionName = cCString(functionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	var callResult = cWasmerInstanceCall(
		instance.instance,
		wasmFunctionName,
	)

	if callResult != cWasmerOk {
		err := fmt.Errorf("failed to call the `%s` exported function", functionName)
		return newWrappedError(err)
	}

	return nil
}

func (instance *Instance) HasFunction(functionName string) bool {
	var wasmFunctionName = cCString(functionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	result := cWasmerInstanceHasFunction(
		instance.instance,
		wasmFunctionName,
	)

	return result == 1
}

// GetLastError returns the last error message if any, otherwise returns an error.
func (instance *Instance) getFunctionNamesConcat() (string, error) {
	var bufferLength = cWasmerInstanceExportedFunctionNamesLength(instance.instance)

	if bufferLength == 0 {
		return "", nil
	}

	var buffer = make([]cChar, bufferLength)
	var bufferPointer = (*cChar)(unsafe.Pointer(&buffer[0]))

	var result = cWasmerInstanceExportedFunctionNames(instance.instance, bufferPointer, bufferLength)

	if result == -1 {
		return "", errors.New("cannot read function names")
	}

	return cGoString(bufferPointer), nil
}

func (instance *Instance) GetFunctionNames() []string {
	buffer, err := instance.getFunctionNamesConcat()
	if err != nil {
		return nil
	}
	return strings.Split(buffer, "|")
}

func (instance *Instance) ValidateVoidFunction(functionName string) error {
	// return instance.verifyVoidFunction(functionName)
	return nil
}

// GetData returns a pointer for the current instance's data
func (instance *Instance) GetData() uintptr {
	panic("instance GetData")
}

// GetInstanceCtxMemory returns the memory for the instance context
func (instance *Instance) GetInstanceCtxMemory() executor.MemoryHandler {
	// return instance.InstanceCtx.Memory()
	panic("instance GetInstanceCtxMemory")
}

// GetMemory returns the memory for the instance
func (instance *Instance) GetMemory() executor.MemoryHandler {
	// return instance.Memory
	panic("instance GetMemory")
}

// SetMemory sets the memory for the instance returns true if success
func (instance *Instance) SetMemory(cleanMemory []byte) bool {
	// instanceMemory := instance.GetMemory().Data()
	// if len(instanceMemory) != len(cleanMemory) {
	// 	// TODO shrink the instance memory instead and return true
	// 	return false
	// }

	// copy(instanceMemory, cleanMemory)
	// return true
	panic("instance SetMemory")
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *Instance) IsInterfaceNil() bool {
	return instance == nil
}
