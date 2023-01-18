package wasmer

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/multiversx/wasm-vm/executor"
)

// ImportedFunctionError represents any kind of errors related to a
// WebAssembly imported function. It is returned by `Import` or `Imports`
// functions only.
type ImportedFunctionError struct {
	functionName string
	message      string
}

// NewImportedFunctionError constructs a new `ImportedFunctionError`,
// where `functionName` is the name of the imported function, and
// `message` is the error message. If the error message contains `%s`,
// then this parameter will be replaced by `functionName`.
func NewImportedFunctionError(functionName string, message string) *ImportedFunctionError {
	return &ImportedFunctionError{functionName, message}
}

// ImportedFunctionError is an actual error. The `Error` function
// returns the error message.
func (error *ImportedFunctionError) Error() string {
	return fmt.Sprintf(error.message, error.functionName)
}

// wasmerImport represents an WebAssembly instance imported function.
type wasmerImport struct {
	// An implementation must be of type:
	// `func(context unsafe.Pointer, arguments ...interface{}) interface{}`.
	// It represents the real function implementation written in Go.
	implementation interface{}

	// The pointer to the cgo function implementation, something
	// like `C.foo`.
	cgoPointer unsafe.Pointer

	// The pointer to the Wasmer imported function.
	importedFunctionPointer *cWasmerImportFuncT

	// The function implementation signature as a WebAssembly signature.
	wasmInputs []cWasmerValueTag

	// The function implementation signature as a WebAssembly signature.
	wasmOutputs []cWasmerValueTag

	// The namespace of the imported function.
	namespace string
}

// wasmerImports represents a set of imported functions for a WebAssembly instance.
type wasmerImports struct {
	// All imports.
	imports map[string]map[string]wasmerImport

	// Current namespace where to register the import.
	currentNamespace string
}

// newWasmerImports constructs a new empty `wasmerImports`.
func newWasmerImports() *wasmerImports {
	var imports = make(map[string]map[string]wasmerImport)
	var currentNamespace = "env"

	return &wasmerImports{imports, currentNamespace}
}

// Count returns the number of imports
func (imports *wasmerImports) Count() int {
	count := 0
	for _, namespacedImports := range imports.imports {
		count += len(namespacedImports)
	}
	return count
}

// Append adds a new imported function to the current set.
func (imports *wasmerImports) append(importName string, implementation interface{}, cgoPointer unsafe.Pointer) error {
	var importType = reflect.TypeOf(implementation)

	if importType.Kind() != reflect.Func {
		return NewImportedFunctionError(importName, fmt.Sprintf("Imported function `%%s` must be a function; given `%s`.", importType.Kind()))
	}

	var importInputsArity = importType.NumIn()

	if importInputsArity < 1 {
		return NewImportedFunctionError(importName, "Imported function `%s` must at least have one argument for the instance context.")
	}

	if importType.In(0).Kind() != reflect.UnsafePointer {
		return NewImportedFunctionError(importName, fmt.Sprintf("The instance context of the `%%s` imported function must be of kind `unsafe.Pointer`; given `%s`; is it missing?", importType.In(0).Kind()))
	}

	importInputsArity--
	var importOutputsArity = importType.NumOut()
	var wasmInputs = make([]cWasmerValueTag, importInputsArity)
	var wasmOutputs = make([]cWasmerValueTag, importOutputsArity)

	for nth := 0; nth < importInputsArity; nth++ {
		var importInput = importType.In(nth + 1)

		switch importInput.Kind() {
		case reflect.Int32:
			wasmInputs[nth] = cWasmI32
		case reflect.Int64:
			wasmInputs[nth] = cWasmI64
		default:
			return NewImportedFunctionError(importName, fmt.Sprintf("Invalid input type for the `%%s` imported function; given `%s`; only accept `int32`, `int64`, `float32`, and `float64`.", importInput.Kind()))
		}
	}

	if importOutputsArity > 1 {
		return NewImportedFunctionError(importName, "The `%s` imported function must have at most one output value.")
	} else if importOutputsArity == 1 {
		switch importType.Out(0).Kind() {
		case reflect.Int32:
			wasmOutputs[0] = cWasmI32
		case reflect.Int64:
			wasmOutputs[0] = cWasmI64
		default:
			return NewImportedFunctionError(importName, fmt.Sprintf("Invalid output type for the `%%s` imported function; given `%s`; only accept `int32`, `int64`, `float32`, and `float64`.", importType.Out(0).Kind()))
		}
	}

	var importedFunctionPointer *cWasmerImportFuncT
	var namespace = imports.currentNamespace

	if imports.imports[namespace] == nil {
		imports.imports[namespace] = make(map[string]wasmerImport)
	}

	imports.imports[namespace][importName] = wasmerImport{
		implementation,
		cgoPointer,
		importedFunctionPointer,
		wasmInputs,
		wasmOutputs,
		namespace,
	}

	return nil
}

// Close closes/frees all imported functions that have been registered by Wasmer.
func (imports *wasmerImports) Close() {
	for _, namespacedImports := range imports.imports {
		for _, importFunction := range namespacedImports {
			if nil != importFunction.importedFunctionPointer {
				cWasmerImportFuncDestroy(importFunction.importedFunctionPointer)
			}
		}
	}
}

// InstanceContext represents a way to access instance API from within
// an imported context.
type InstanceContext struct {
	context *cWasmerInstanceContextT
	memory  executor.Memory
}

// IntoInstanceContext casts the first `context unsafe.Pointer`
// argument of an imported function into an `InstanceContext`.
func IntoInstanceContext(instanceContext unsafe.Pointer) InstanceContext {
	context := (*cWasmerInstanceContextT)(instanceContext)
	memory := newMemory(cWasmerInstanceContextMemory(context))

	return InstanceContext{context, &memory}
}

// IntoInstanceContextDirect retrieves the Wasmer instance context directly
// from the Wasmer instance. This context can be stored as long as the instance itself.
func IntoInstanceContextDirect(instanceContext *cWasmerInstanceContextT) InstanceContext {
	memory := newMemory(cWasmerInstanceContextMemory(instanceContext))
	return InstanceContext{instanceContext, &memory}
}

// Memory returns the current instance memory.
func (instanceContext *InstanceContext) Memory() executor.Memory {
	return instanceContext.memory
}

// Data returns the instance context data as an `unsafe.Pointer`. It's
// up to the user to cast it appropriately as a pointer to a data.
func (instanceContext *InstanceContext) Data() unsafe.Pointer {
	return cWasmerInstanceContextDataGet(instanceContext.context)
}
