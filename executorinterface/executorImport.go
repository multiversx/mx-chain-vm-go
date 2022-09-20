package executorinterface

import (
	"fmt"
	"reflect"
	"unsafe"
)

// EIFunctionValue represents EI argument types
type EIFunctionValue int

const (
	EIFunctionValueInt32 EIFunctionValue = iota
	EIFunctionValueInt64
)

// EIFunctionReceiver abstracts an EI imports container, where EI functions are registered.
type EIFunctionReceiver interface {
	Namespace(namespace string)
	Append(importName string, implementation interface{}, cgoPointer unsafe.Pointer) error
}

// EIFunction represents a EI function that gets imported in a constract WASM module.
type EIFunction struct {
	// An implementation must be of type:
	// `func(context unsafe.Pointer, arguments ...interface{}) interface{}`.
	// It represents the real function implementation written in Go.
	Implementation interface{}

	// The pointer to the cgo function implementation, something
	// like `C.foo`.
	CgoPointer unsafe.Pointer

	// The function implementation signature.
	FunctionInputs []EIFunctionValue

	// The function implementation signature.
	FunctionOutputs []EIFunctionValue

	// The namespace of the imported function.
	Namespace string
}

// EIFunctions holds a collection of EI functions.
type EIFunctions struct {
	FunctionMap map[string]EIFunction

	// Current namespace where to register the import.
	CurrentNamespace string
}

// NewImportFunctions constructs a new empty `EIFunctions`.
func NewImportFunctions() *EIFunctions {
	return &EIFunctions{
		FunctionMap:      make(map[string]EIFunction),
		CurrentNamespace: "env",
	}
}

// Namespace changes the current namespace of the next imported functions.
func (imports *EIFunctions) Namespace(namespace string) {
	imports.CurrentNamespace = namespace
}

// Append validates and adds a new imported function to the current structure.
func (imports *EIFunctions) Append(importName string, implementation interface{}, cgoPointer unsafe.Pointer) error {
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
	var wasmInputs = make([]EIFunctionValue, importInputsArity)
	var wasmOutputs = make([]EIFunctionValue, importOutputsArity)

	for nth := 0; nth < importInputsArity; nth++ {
		var importInput = importType.In(nth + 1)

		switch importInput.Kind() {
		case reflect.Int32:
			wasmInputs[nth] = EIFunctionValueInt32
		case reflect.Int64:
			wasmInputs[nth] = EIFunctionValueInt64
		default:
			return NewImportFunctionError(importName, fmt.Sprintf("Invalid input type for the `%%s` imported function; given `%s`; only accept `int32`, `int64`, `float32`, and `float64`.", importInput.Kind()))
		}
	}

	if importOutputsArity > 1 {
		return NewImportFunctionError(importName, "The `%s` imported function must have at most one output value.")
	} else if importOutputsArity == 1 {
		switch importType.Out(0).Kind() {
		case reflect.Int32:
			wasmOutputs[0] = EIFunctionValueInt32
		case reflect.Int64:
			wasmOutputs[0] = EIFunctionValueInt64
		default:
			return NewImportFunctionError(importName, fmt.Sprintf("Invalid output type for the `%%s` imported function; given `%s`; only accept `int32`, `int64`, `float32`, and `float64`.", importType.Out(0).Kind()))
		}
	}

	var namespace = imports.CurrentNamespace

	if _, duplicate := imports.FunctionMap[importName]; duplicate {
		return NewImportFunctionError(importName, "Duplicate imported function `%s`.")
	}

	imports.FunctionMap[importName] = EIFunction{
		implementation,
		cgoPointer,
		wasmInputs,
		wasmOutputs,
		namespace,
	}

	return nil
}
