package wasmer

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const OPCODE_COUNT = 448

var logWasmer = logger.GetOrCreate("arwen/wasmer")

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

// SetRkyvSerializationEnabled enables or disables RKYV serialization of
// instances in Wasmer
func SetRkyvSerializationEnabled(enabled bool) {
	if enabled {
		cWasmerInstanceEnableRkyv()
	} else {
		cWasmerInstanceDisableRkyv()
	}
}

// SetSIGSEGVPassthrough instructs Wasmer to never register a handler for
// SIGSEGV. Only has effect if called before creating the first Wasmer instance
// since the process started. Calling this function after the first Wasmer
// instance will not unregister the signal handler set by Wasmer.
func SetSIGSEGVPassthrough() {
	cWasmerSetSIGSEGVPassthrough()
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
	Memory MemoryHandler

	Data        *uintptr
	DataPointer unsafe.Pointer

	InstanceCtx InstanceContext

	AlreadyClean bool
}

type CompilationOptions struct {
	GasLimit           uint64
	UnmeteredLocals    uint64
	MaxMemoryGrow      uint64
	MaxMemoryGrowDelta uint64
	OpcodeTrace        bool
	Metering           bool
	RuntimeBreakpoints bool
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

func SetImports(imports *Imports) error {
	wasmImportsCPointer, numberOfImports := generateWasmerImports(imports)

	var result = cWasmerCacheImportObjectFromImports(
		wasmImportsCPointer,
		cInt(numberOfImports),
	)

	if result != cWasmerOk {
		return newWrappedError(ErrFailedCacheImports)
	}
	return nil
}

func SetOpcodeCosts(opcodeCosts *[OPCODE_COUNT]uint32) {
	cWasmerSetOpcodeCosts(opcodeCosts)
}

func NewInstanceWithOptions(
	bytes []byte,
	options CompilationOptions,
) (*Instance, error) {
	var cInstance *cWasmerInstanceT

	if len(bytes) == 0 {
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
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
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	instance, err := newInstance(cInstance)
	if instance != nil && instance.Memory != nil {
		cInstanceContext := cWasmerInstanceContextGet(cInstance)
		instance.InstanceCtx = IntoInstanceContextDirect(cInstanceContext)
	}
	return instance, err
}

func newInstance(cInstance *cWasmerInstanceT) (*Instance, error) {
	var emptyInstance = &Instance{instance: nil, Exports: nil, Signatures: nil, Memory: nil}

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
		return &Instance{instance: cInstance, Exports: exports, Signatures: signatures, Memory: nil}, nil
	}

	return &Instance{instance: cInstance, Exports: exports, Signatures: signatures, Memory: &memory}, nil
}

// HasMemory checks whether the instance has at least one exported memory.
func (instance *Instance) HasMemory() bool {
	return nil != instance.Memory
}

func NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options CompilationOptions,
) (*Instance, error) {
	var cInstance *cWasmerInstanceT

	if len(compiledCode) == 0 {
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
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
		var emptyInstance = &Instance{instance: nil, Exports: nil, Memory: nil}
		return emptyInstance, newWrappedError(ErrFailedInstantiation)
	}

	instance, err := newInstance(cInstance)
	if instance != nil && instance.Memory != nil {
		cInstanceContext := cWasmerInstanceContextGet(cInstance)
		instance.InstanceCtx = IntoInstanceContextDirect(cInstanceContext)
	}

	return instance, err
}

// SetContextData assigns a data that can be used by all imported
// functions. Indeed, each imported function receives as its first
// argument an instance context (see `InstanceContext`). An instance
// context can hold a pointer to any kind of data. It is important to
// understand that this data is shared by all imported function, it's
// global to the instance.
func (instance *Instance) SetContextData(data uintptr) {
	instance.Data = &data
	instance.DataPointer = unsafe.Pointer(instance.Data)
	cWasmerInstanceContextDataSet(instance.instance, instance.DataPointer)
}

// Clean cleans instance
func (instance *Instance) Clean() bool {
	logWasmer.Trace("cleaning instance", "id", instance.Id())
	if instance.AlreadyClean {
		logWasmer.Trace("clean: already cleaned instance", "id", instance.Id())
		return false
	}

	if instance.instance != nil {
		cWasmerInstanceDestroy(instance.instance)

		if instance.Memory != nil {
			instance.Memory.Destroy()
		}

		instance.AlreadyClean = true
		logWasmer.Trace("cleaned instance", "id", instance.Id())

		return true
	}

	return false
}

// AlreadyCleaned returns the internal field AlreadyClean
func (instance *Instance) AlreadyCleaned() bool {
	return instance.AlreadyClean
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

// SetBreakpointValue sets the breakpoint value for the instance
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

// GetExports returns the exports map for the current instance
func (instance *Instance) GetExports() ExportsMap {
	return instance.Exports
}

// GetSignature returns the signature for the given functionName
func (instance *Instance) GetSignature(functionName string) (*ExportedFunctionSignature, bool) {
	signature, ok := instance.Signatures[functionName]
	return signature, ok
}

// GetData returns a pointer for the current instance's data
func (instance *Instance) GetData() uintptr {
	return *instance.Data
}

// GetInstanceCtxMemory returns the memory for the instance context
func (instance *Instance) GetInstanceCtxMemory() MemoryHandler {
	return instance.InstanceCtx.Memory()
}

// Id returns an identifier for the instance, unique at runtime
func (instance *Instance) Id() string {
	return fmt.Sprintf("%p", instance.instance)
}

// GetMemory returns the memory for the instance
func (instance *Instance) GetMemory() MemoryHandler {
	return instance.Memory
}

// Reset resets the instance memories and globals
func (instance *Instance) Reset() bool {
	if instance.AlreadyClean {
		logWasmer.Trace("reset: already cleaned instance", "id", instance.Id())
		return false
	}

	result := cWasmerInstanceReset(instance.instance)
	logWasmer.Trace("reset: warm instance", "id", instance.Id())
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
