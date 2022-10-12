package wasmer

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/executor"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
)

// InstanceError represents any kind of errors related to a WebAssembly instance. It
// is returned by `Instance` functions only.
type InstanceError struct {
	// Error message.
	message string
}

// NewInstanceError constructs a new `InstanceError`.
func NewInstanceError(message string) *InstanceError {
	return &InstanceError{message}
}

// `InstanceError` is an actual error. The `Error` function returns
// the error message.
func (error *InstanceError) Error() string {
	return error.message
}

// ExportedFunctionError represents any kind of errors related to a
// WebAssembly exported function. It is returned by `Instance`
// functions only.
type ExportedFunctionError struct {
	functionName string
	message      string
}

// ExportedFunctionSignature holds information about the input/output arities
// of an exported function
type ExportedFunctionSignature struct {
	InputArity  int
	OutputArity int
}

// NewExportedFunctionError constructs a new `ExportedFunctionError`,
// where `functionName` is the name of the exported function, and
// `message` is the error message. If the error message contains `%s`,
// then this parameter will be replaced by `functionName`.
func NewExportedFunctionError(functionName string, message string) *ExportedFunctionError {
	return &ExportedFunctionError{functionName, message}
}

// ExportedFunctionError is an actual error. The `Error` function
// returns the error message.
func (error *ExportedFunctionError) Error() string {
	return error.message
}

type ExportedFunctionCallback func(...interface{}) (Value, error)
type ExportsMap map[string]ExportedFunctionCallback
type ExportSignaturesMap map[string]*ExportedFunctionSignature

// Instance represents a WebAssembly instance.
type Instance struct {
	// The underlying WebAssembly instance.
	instance *cWasmerInstanceT

	// All functions exported by the WebAssembly instance, indexed
	// by their name as a string. An exported function is a
	// regular variadic Go closure. Arguments are untyped. Since
	// WebAssembly only supports: `i32`, `i64`, `f32` and `f64`
	// types, the accepted Go types are: `int8`, `uint8`, `int16`,
	// `uint16`, `int32`, `uint32`, `int64`, `int`, `uint`, `float32`
	// and `float64`. In addition to those types, the `Value` type
	// (from this project) is accepted. The conversion from a Go
	// value to a WebAssembly value is done automatically except for
	// the `Value` type (where type is coerced, that's the intent
	// here). The WebAssembly type is automatically inferred. Note
	// that the returned value is of kind `Value`, and not a
	// standard Go type.
	Exports ExportsMap

	Signatures ExportSignaturesMap

	// The exported memory of a WebAssembly instance.
	Memory executor.MemoryHandler

	callbacks       executor.ImportsInterface
	callbacksPtr    uintptr
	callbacksPtrPtr unsafe.Pointer

	InstanceCtx InstanceContext
}

func newWrappedError(target error) error {
	var lastError string
	var err error
	lastError, err = GetLastError()

	if err != nil {
		lastError = "unknown details"
	}

	return fmt.Errorf("%w: %s", target, lastError)
}

func NewInstanceWithOptions(
	bytes []byte,
	options executor.CompilationOptions,
) (*Instance, error) {
	var c_instance *cWasmerInstanceT

	if len(bytes) == 0 {
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
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
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	instance, err := newInstance(c_instance)
	if instance != nil && instance.Memory != nil {
		c_instance_context := cWasmerInstanceContextGet(c_instance)
		instance.InstanceCtx = IntoInstanceContextDirect(c_instance_context)
	}
	return instance, err
}

func newInstance(c_instance *cWasmerInstanceT) (*Instance, error) {
	var emptyInstance = &Instance{instance: nil, Exports: nil, Signatures: nil, Memory: nil}

	var wasmExports *cWasmerExportsT
	var hasMemory bool

	cWasmerInstanceExports(c_instance, &wasmExports)
	defer cWasmerExportsDestroy(wasmExports)

	exports, signatures, err := retrieveExportedFunctions(c_instance, wasmExports)
	if err != nil {
		return emptyInstance, err
	}

	memory, hasMemory, err := retrieveExportedMemory(wasmExports)
	if err != nil {
		return emptyInstance, err
	}

	if !hasMemory {
		return &Instance{instance: c_instance, Exports: exports, Signatures: signatures, Memory: nil}, nil
	}

	return &Instance{instance: c_instance, Exports: exports, Signatures: signatures, Memory: &memory}, nil
}

// HasMemory checks whether the instance has at least one exported memory.
func (instance *Instance) HasMemory() bool {
	return nil != instance.Memory
}

func NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (*Instance, error) {
	var c_instance *cWasmerInstanceT

	if len(compiledCode) == 0 {
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var instantiateResult = cWasmerInstanceFromCache(
		&c_instance,
		(*cUchar)(unsafe.Pointer(&compiledCode[0])),
		cUint32T(len(compiledCode)),
		(*cWasmerCompilationOptions)(cOptions),
	)

	if instantiateResult != cWasmerOk {
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	instance, err := newInstance(c_instance)
	if instance != nil && instance.Memory != nil {
		c_instance_context := cWasmerInstanceContextGet(c_instance)
		instance.InstanceCtx = IntoInstanceContextDirect(c_instance_context)
	}

	return instance, err
}

// SetCallbacks assigns a data that can be used by all imported
// functions. Indeed, each imported function receives as its first
// argument an instance context (see `InstanceContext`). An instance
// context can hold a pointer to any kind of data. It is important to
// understand that this data is shared by all imported function, it's
// global to the instance.
func (instance *Instance) SetCallbacks(callbacks executor.ImportsInterface) {
	instance.callbacks = callbacks
	// This has to be a local variable, to fool Go into thinking this has nothing to do with the other structures.
	localPtr := uintptr(unsafe.Pointer(&instance.callbacks))
	instance.callbacksPtr = localPtr
	instance.callbacksPtrPtr = unsafe.Pointer(&localPtr)
	cWasmerInstanceContextDataSet(instance.instance, instance.callbacksPtrPtr)
}

// GetCallbacks returns a pointer for the current instance's data
func (instance *Instance) GetCallbacks() executor.ImportsInterface {
	return instance.callbacks
}

// Clean cleans instance
func (instance *Instance) Clean() {
	if instance.instance != nil {
		cWasmerInstanceDestroy(instance.instance)

		if instance.Memory != nil {
			instance.Memory.Destroy()
		}
	}
}

// GetPointsUsed returns the internal instance gas counter
func (instance *Instance) GetPointsUsed() uint64 {
	return cWasmerInstanceGetPointsUsed(instance.instance)
}

// SetPointsUsed sets the internal instance gas counter
func (instance *Instance) SetPointsUsed(points uint64) {
	cWasmerInstanceSetPointsUsed(instance.instance, points)
}

// SetGasLimit sets the gas limit for the instance
func (instance *Instance) SetGasLimit(gasLimit uint64) {
	cWasmerInstanceSetGasLimit(instance.instance, gasLimit)
}

// SetBreakpoints sets the breakpoint value for the instance
func (instance *Instance) SetBreakpointValue(value uint64) {
	cWasmerInstanceSetBreakpointValue(instance.instance, value)
}

// GetBreakpointValue returns the breakpoint value
func (instance *Instance) GetBreakpointValue() uint64 {
	return cWasmerInstanceGetBreakpointValue(instance.instance)
}

// Cache caches the instance
func (instance *Instance) Cache() ([]byte, error) {
	var cacheBytes *cUchar
	var cacheLen cUint32T

	var cacheResult = cWasmerInstanceCache(
		instance.instance,
		&cacheBytes,
		&cacheLen,
	)

	if cacheResult != cWasmerOk {
		return nil, ErrCachingFailed
	}

	goBytes := C.GoBytes(unsafe.Pointer(cacheBytes), C.int(cacheLen))

	C.free(unsafe.Pointer(cacheBytes))
	cacheBytes = nil
	return goBytes, nil
}

// IsFunctionImported returns true if the instance imports the specified function
func (instance *Instance) IsFunctionImported(name string) bool {
	return cWasmerInstanceIsFunctionImported(instance.instance, name)
}

// CallFunction executes given function from loaded contract.
func (instance *Instance) CallFunction(functionName string) error {
	if function, ok := instance.Exports[functionName]; ok {
		_, err := function()
		return err
	}

	return executor.ErrFuncNotFound
}

// HasFunction checks if loaded contract has a function (endpoint) with given name.
func (instance *Instance) HasFunction(functionName string) bool {
	_, ok := instance.Exports[functionName]
	return ok
}

// GetFunctionNames loads a list of contract function (endpoint) names. Required for validating reserved names.
func (instance *Instance) GetFunctionNames() []string {
	var functionNames []string
	for functionName := range instance.Exports {
		functionNames = append(functionNames, functionName)
	}
	return functionNames
}

// ValidateVoidFunction checks that no function (endpoint) of the given contract has any parameters or returns any result.
// All arguments and results should be transferred via the import functions.
func (instance *Instance) ValidateVoidFunction(functionName string) error {
	return instance.verifyVoidFunction(functionName)
}

// GetInstanceCtxMemory returns the memory for the instance context
func (instance *Instance) GetInstanceCtxMemory() executor.MemoryHandler {
	return instance.InstanceCtx.Memory()
}

// GetMemory returns the memory for the instance
func (instance *Instance) GetMemory() executor.MemoryHandler {
	return instance.Memory
}

// Reset resets the instance memories and globals
func (instance *Instance) Reset() bool {
	result := cWasmerInstanceReset(instance.instance)
	return result == cWasmerOk
}

// SetMemory sets the memory for the instance returns true if success
func (instance *Instance) SetMemory(data []byte) bool {
	if instance.instance == nil {
		return false
	}

	if check.IfNil(instance.GetMemory()) {
		return false
	}

	memory := instance.GetMemory().Data()
	if len(memory) != len(data) {
		return false
	}

	copy(memory, data)
	return true
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *Instance) IsInterfaceNil() bool {
	return instance == nil
}
