package wasmer

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
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

// ExportsMap is a map of names to ExportedFunctionCallInfo values
type ExportsMap map[string]*ExportedFunctionCallInfo

// ExportSignaturesMap is a map of names to ExportedFunctionSignatures
type ExportSignaturesMap map[string]*ExportedFunctionSignature

// ExportedFunctionCallInfo contains information required to call an exported WASM function
type ExportedFunctionCallInfo struct {
	FuncName       string
	InputArity     cUint32T
	InputSignature []cWasmerValueTag
	OutputArity    cUint32T
}

// WasmerInstance represents a WebAssembly instance.
type WasmerInstance struct {
	// The underlying WebAssembly instance.
	instance *cWasmerInstanceT

	AlreadyClean bool

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
	exports ExportsMap

	signatures ExportSignaturesMap

	// The exported memory of a WebAssembly instance.
	Memory executor.Memory

	// The instance context.
	InstanceCtx InstanceContext

	vmHooksPtr unsafe.Pointer
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

// NewInstanceWithOptions creates a new instance from provided bytes & options
func NewInstanceWithOptions(
	bytes []byte,
	options executor.CompilationOptions,
) (*WasmerInstance, error) {
	var cInstance *cWasmerInstanceT

	if len(bytes) == 0 {
		var emptyInstance = &WasmerInstance{instance: nil, exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var compileResult = cWasmerInstantiateWithOptions(
		&cInstance,
		(*cUchar)(unsafe.Pointer(&bytes[0])),
		cUint(len(bytes)),
		(*cWasmerCompilationOptions)(cOptions),
	)

	if compileResult != cWasmerOk {
		var emptyInstance = &WasmerInstance{instance: nil, exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	instance, err := newInstance(cInstance)
	if instance != nil && instance.Memory != nil {
		cInstanceContext := cWasmerInstanceContextGet(cInstance)
		instance.InstanceCtx = IntoInstanceContextDirect(cInstanceContext)
	}

	logWasmer.Trace("new instance created", "id", instance.ID())
	return instance, err
}

func newInstance(cInstance *cWasmerInstanceT) (*WasmerInstance, error) {
	var emptyInstance = &WasmerInstance{instance: nil, exports: nil, signatures: nil, Memory: nil}

	var wasmExports *cWasmerExportsT
	var hasMemory bool

	cWasmerInstanceExports(cInstance, &wasmExports)
	defer cWasmerExportsDestroy(wasmExports)

	exports, signatures, err := retrieveExportedFunctions(cInstance, wasmExports)
	if err != nil {
		return emptyInstance, err
	}

	memory, hasMemory, err := retrieveExportedMemory(wasmExports)
	if err != nil {
		return emptyInstance, err
	}

	if !hasMemory {
		return &WasmerInstance{instance: cInstance, exports: exports, signatures: signatures, Memory: nil}, nil
	}

	return &WasmerInstance{instance: cInstance, exports: exports, signatures: signatures, Memory: &memory}, nil
}

// NewInstanceFromCompiledCodeWithOptions creates a new instance from compiled code
func NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (*WasmerInstance, error) {
	var cInstance *cWasmerInstanceT

	if len(compiledCode) == 0 {
		var emptyInstance = &WasmerInstance{instance: nil, exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var instantiateResult = cWasmerInstanceFromCache(
		&cInstance,
		(*cUchar)(unsafe.Pointer(&compiledCode[0])),
		cUint32T(len(compiledCode)),
		(*cWasmerCompilationOptions)(cOptions),
	)

	if instantiateResult != cWasmerOk {
		var emptyInstance = &WasmerInstance{instance: nil, exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	instance, err := newInstance(cInstance)
	if instance != nil && instance.Memory != nil {
		cInstanceContext := cWasmerInstanceContextGet(cInstance)
		instance.InstanceCtx = IntoInstanceContextDirect(cInstanceContext)
	}

	return instance, err
}

// Clean cleans instance
func (instance *WasmerInstance) Clean() bool {
	logWasmer.Trace("cleaning instance", "id", instance.ID())
	if instance.AlreadyClean {
		logWasmer.Trace("clean: already cleaned instance", "id", instance.ID())
		return false
	}

	if instance.instance != nil {
		cWasmerInstanceDestroy(instance.instance)

		if instance.Memory != nil {
			instance.Memory.Destroy()
		}

		instance.AlreadyClean = true
		logWasmer.Trace("cleaned instance", "id", instance.ID())

		return true
	}

	return false
}

// IsAlreadyCleaned returns the internal field AlreadyClean
func (instance *WasmerInstance) IsAlreadyCleaned() bool {
	return instance.AlreadyClean
}

// GetPointsUsed returns the internal instance gas counter
func (instance *WasmerInstance) GetPointsUsed() uint64 {
	return cWasmerInstanceGetPointsUsed(instance.instance)
}

// SetPointsUsed sets the internal instance gas counter
func (instance *WasmerInstance) SetPointsUsed(points uint64) {
	cWasmerInstanceSetPointsUsed(instance.instance, points)
}

// SetGasLimit sets the gas limit for the instance
func (instance *WasmerInstance) SetGasLimit(gasLimit uint64) {
	cWasmerInstanceSetGasLimit(instance.instance, gasLimit)
}

// SetBreakpointValue sets the breakpoint value for the instance
func (instance *WasmerInstance) SetBreakpointValue(value uint64) {
	cWasmerInstanceSetBreakpointValue(instance.instance, value)
}

// GetBreakpointValue returns the breakpoint value
func (instance *WasmerInstance) GetBreakpointValue() uint64 {
	return cWasmerInstanceGetBreakpointValue(instance.instance)
}

// HasCompiledCode specifies if the code is compiled
func (instance *WasmerInstance) HasCompiledCode() bool {
	return true
}

// Cache caches the instance
func (instance *WasmerInstance) Cache() ([]byte, error) {
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
func (instance *WasmerInstance) IsFunctionImported(name string) bool {
	return cWasmerInstanceIsFunctionImported(instance.instance, name)
}

// CallFunction executes given function from loaded contract.
func (instance *WasmerInstance) CallFunction(functionName string) error {
	callInfo, found := instance.exports[functionName]
	if !found {
		return executor.ErrFuncNotFound
	}

	return callExportedFunction(instance.instance, callInfo)
}

// HasFunction checks if loaded contract has a function (endpoint) with given name.
func (instance *WasmerInstance) HasFunction(functionName string) bool {
	_, ok := instance.exports[functionName]
	return ok
}

// GetFunctionNames loads a list of contract function (endpoint) names. Required for validating reserved names.
func (instance *WasmerInstance) GetFunctionNames() []string {
	var functionNames []string
	for functionName := range instance.exports {
		functionNames = append(functionNames, functionName)
	}
	return functionNames
}

// ValidateFunctionArities checks that no function (endpoint) of the given contract has any parameters or returns any result.
// All arguments and results should be transferred via the import functions.
func (instance *WasmerInstance) ValidateFunctionArities() error {
	for functionName := range instance.exports {
		err := instance.verifyVoidFunction(functionName)
		if err != nil {
			return err
		}
	}
	return nil
}

// HasMemory checks whether the instance has at least one exported memory.
func (instance *WasmerInstance) HasMemory() bool {
	return nil != instance.Memory
}

// MemLoad returns the contents from the given offset of the WASM memory.
func (instance *WasmerInstance) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return executor.MemLoadFromMemory(instance.Memory, memPtr, length)
}

// MemStore stores the given data in the WASM memory at the given offset.
func (instance *WasmerInstance) MemStore(memPtr executor.MemPtr, data []byte) error {
	return executor.MemStoreToMemory(instance.Memory, memPtr, data)
}

// MemLength returns the length of the allocated memory. Only called directly in tests.
func (instance *WasmerInstance) MemLength() uint32 {
	return instance.Memory.Length()
}

// MemGrow allocates more pages to the current memory. Only called directly in tests.
func (instance *WasmerInstance) MemGrow(pages uint32) error {
	return instance.Memory.Grow(pages)
}

// MemDump yields the entire contents of the memory. Only used in tests.
func (instance *WasmerInstance) MemDump() []byte {
	return instance.Memory.Data()
}

// Reset resets the instance memories and globals
func (instance *WasmerInstance) Reset() bool {
	if instance.AlreadyClean {
		logWasmer.Trace("reset: already cleaned instance", "id", instance.ID())
		return false
	}

	result := cWasmerInstanceReset(instance.instance)
	ok := result == cWasmerOk

	logWasmer.Trace("reset: warm instance", "id", instance.ID(), "ok", ok)
	return ok
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *WasmerInstance) IsInterfaceNil() bool {
	return instance == nil
}

// SetVMHooksPtr sets the VM hooks pointer
func (instance *WasmerInstance) SetVMHooksPtr(vmHooksPtr uintptr) {
	localVMHooksPointer := unsafe.Pointer(&vmHooksPtr)
	instance.vmHooksPtr = localVMHooksPointer
	cWasmerInstanceContextDataSet(instance.instance, localVMHooksPointer)
}

// GetVMHooksPtr returns the VM hooks pointer
func (instance *WasmerInstance) GetVMHooksPtr() uintptr {
	return *(*uintptr)(instance.vmHooksPtr)
}

// ID returns an identifier for the instance, unique at runtime
func (instance *WasmerInstance) ID() string {
	return fmt.Sprintf("%p", instance.instance)
}
