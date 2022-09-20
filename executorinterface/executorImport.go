package executorinterface

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ImportFunctionValue represents EI argument types
type ImportFunctionValue int

const (
	ImportFunctionValueInt32 ImportFunctionValue = iota
	ImportFunctionValueInt64
)

// ImportFunctionReceiver abstracts an EI imports container, where EI functions are registered.
type ImportFunctionReceiver interface {
	Namespace(namespace string)
	Append(importName string, implementation interface{}, cgoPointer unsafe.Pointer) error
}

// ImportFunction represents a EI function that gets imported in a constract WASM module.
type ImportFunction struct {
	// An implementation must be of type:
	// `func(context unsafe.Pointer, arguments ...interface{}) interface{}`.
	// It represents the real function implementation written in Go.
	Implementation interface{}

	// The pointer to the cgo function implementation, something
	// like `C.foo`.
	CgoPointer unsafe.Pointer

	// The function implementation signature.
	FunctionInputs []ImportFunctionValue

	// The function implementation signature.
	FunctionOutputs []ImportFunctionValue

	// The namespace of the imported function.
	Namespace string
}

// ImportFunctions holds a collection of EI functions.
type ImportFunctions struct {
	FunctionMap map[string]ImportFunction

	// Current namespace where to register the import.
	CurrentNamespace string
}

// NewImportFunctions constructs a new empty `EIFunctions`.
func NewImportFunctions() *ImportFunctions {
	return &ImportFunctions{
		FunctionMap:      make(map[string]ImportFunction),
		CurrentNamespace: "env",
	}
}

// Namespace changes the current namespace of the next imported functions.
func (imports *ImportFunctions) Namespace(namespace string) {
	imports.CurrentNamespace = namespace
}

// Append validates and adds a new imported function to the current structure.
func (imports *ImportFunctions) Append(importName string, implementation interface{}, cgoPointer unsafe.Pointer) error {
	var importType = reflect.TypeOf(implementation)

	if importType.Kind() != reflect.Func {
		return NewImportFunctionError(importName, fmt.Sprintf("Imported function `%%s` must be a function; given `%s`.", importType.Kind()))
	}

	var importInputsArity = importType.NumIn()

	if importInputsArity < 1 {
		return NewImportFunctionError(importName, "Imported function `%s` must at least have one argument for the instance context.")
	}

	if importType.In(0).Kind() != reflect.UnsafePointer {
		return NewImportFunctionError(importName, fmt.Sprintf("The instance context of the `%%s` imported function must be of kind `unsafe.Pointer`; given `%s`; is it missing?", importType.In(0).Kind()))
	}

	importInputsArity--
	var importOutputsArity = importType.NumOut()
	var wasmInputs = make([]ImportFunctionValue, importInputsArity)
	var wasmOutputs = make([]ImportFunctionValue, importOutputsArity)

	for nth := 0; nth < importInputsArity; nth++ {
		var importInput = importType.In(nth + 1)

		switch importInput.Kind() {
		case reflect.Int32:
			wasmInputs[nth] = ImportFunctionValueInt32
		case reflect.Int64:
			wasmInputs[nth] = ImportFunctionValueInt64
		default:
			return NewImportFunctionError(importName, fmt.Sprintf("Invalid input type for the `%%s` imported function; given `%s`; only accept `int32`, `int64`, `float32`, and `float64`.", importInput.Kind()))
		}
	}

	if importOutputsArity > 1 {
		return NewImportFunctionError(importName, "The `%s` imported function must have at most one output value.")
	} else if importOutputsArity == 1 {
		switch importType.Out(0).Kind() {
		case reflect.Int32:
			wasmOutputs[0] = ImportFunctionValueInt32
		case reflect.Int64:
			wasmOutputs[0] = ImportFunctionValueInt64
		default:
			return NewImportFunctionError(importName, fmt.Sprintf("Invalid output type for the `%%s` imported function; given `%s`; only accept `int32`, `int64`, `float32`, and `float64`.", importType.Out(0).Kind()))
		}
	}

	var namespace = imports.CurrentNamespace

	if _, duplicate := imports.FunctionMap[importName]; duplicate {
		return NewImportFunctionError(importName, "Duplicate imported function `%s`.")
	}

	imports.FunctionMap[importName] = ImportFunction{
		implementation,
		cgoPointer,
		wasmInputs,
		wasmOutputs,
		namespace,
	}

	return nil
}
