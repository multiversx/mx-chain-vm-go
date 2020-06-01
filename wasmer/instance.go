package wasmer

import (
	"fmt"
	"unsafe"
)

const OPCODE_COUNT = 447

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
	Memory *Memory

	Data unsafe.Pointer
}

type CompilationOptions struct {
	GasLimit           uint64
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

func SetOpcodeCosts(opcode_costs *[OPCODE_COUNT]uint32) {
	cWasmerSetOpcodeCosts(opcode_costs)
}

func NewInstanceWithOptions(
	bytes []byte,
	options CompilationOptions,
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

// SetContextData assigns a data that can be used by all imported
// functions. Indeed, each imported function receives as its first
// argument an instance context (see `InstanceContext`). An instance
// context can hold a pointer to any kind of data. It is important to
// understand that this data is shared by all imported function, it's
// global to the instance.
func (instance *Instance) SetContextData(data unsafe.Pointer) {
	cWasmerInstanceContextDataSet(instance.instance, data)
	instance.Data = data
}

func (instance *Instance) Clean() {
	if instance.instance != nil {
		cWasmerInstanceDestroy(instance.instance)

		if instance.Memory != nil {
			instance.Memory.Destroy()
		}
	}
}

func (instance *Instance) GetPointsUsed() uint64 {
	return cWasmerInstanceGetPointsUsed(instance.instance)
}

func (instance *Instance) SetPointsUsed(points uint64) {
	cWasmerInstanceSetPointsUsed(instance.instance, points)
}

func (instance *Instance) SetBreakpointValue(value uint64) {
	cWasmerInstanceSetBreakpointValue(instance.instance, value)
}

func (instance *Instance) GetBreakpointValue() uint64 {
	return cWasmerInstanceGetBreakpointValue(instance.instance)
}
