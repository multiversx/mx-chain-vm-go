package wasmer2

import (
	"fmt"
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executor"
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

// Import represents an WebAssembly instance imported function.
type Import struct {
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

// Imports represents a set of imported functions for a WebAssembly instance.
type Imports struct {
	// All imports.
	imports map[string]map[string]Import

	// Current namespace where to register the import.
	currentNamespace string
}

// NewImports constructs a new empty `Imports`.
func NewImports() *Imports {
	var imports = make(map[string]map[string]Import)
	var currentNamespace = "env"

	return &Imports{imports, currentNamespace}
}

func (imports *Imports) Count() int {
	count := 0
	for _, namespacedImports := range imports.imports {
		count += len(namespacedImports)
	}
	return count
}

func (imports *Imports) Names() vmcommon.FunctionNames {
	names := make(vmcommon.FunctionNames)
	var empty struct{}
	for _, env := range imports.imports {
		for name := range env {
			names[name] = empty
		}
	}
	return names
}

func convertArgType(argType executor.ImportFunctionValue) cWasmerValueTag {
	switch argType {
	case executor.ImportFunctionValueInt32:
		return cWasmI32
	case executor.ImportFunctionValueInt64:
		return cWasmI64
	}
	return cWasmI32 // unreachable, but might consider adding an error
}

func ConvertImports(eiFunctions *executor.ImportFunctions) *Imports {
	imports := NewImports()

	for funcName, funcData := range eiFunctions.FunctionMap {
		implementation := funcData.Implementation
		cgoPointer := funcData.CgoPointer
		var importedFunctionPointer *cWasmerImportFuncT
		var namespace = funcData.Namespace

		var wasmInputs = make([]cWasmerValueTag, len(funcData.FunctionInputs))
		for i, input := range funcData.FunctionInputs {
			wasmInputs[i] = convertArgType(input)
		}
		var wasmOutputs = make([]cWasmerValueTag, len(funcData.FunctionOutputs))
		for i, output := range funcData.FunctionOutputs {
			wasmOutputs[i] = convertArgType(output)
		}

		if imports.imports[namespace] == nil {
			imports.imports[namespace] = make(map[string]Import)
		}

		imports.imports[namespace][funcName] = Import{
			implementation,
			cgoPointer,
			importedFunctionPointer,
			wasmInputs,
			wasmOutputs,
			namespace,
		}
	}

	return imports
}

// Close closes/frees all imported functions that have been registered by Wasmer.
func (imports *Imports) Close() {
	for _, namespacedImports := range imports.imports {
		for _, importFunction := range namespacedImports {
			if nil != importFunction.importedFunctionPointer {
				cWasmerImportFuncDestroy(importFunction.importedFunctionPointer)
			}
		}
	}
}

func generateWasmerImports(imports *Imports) (*cWasmerImportT, int) {
	var numberOfImports = imports.Count()
	var wasmImports = make([]cWasmerImportT, numberOfImports)
	var importFunctionNth = 0

	for _, namespacedImports := range imports.imports {
		for importName, importFunction := range namespacedImports {
			var wasmInputsArity = len(importFunction.wasmInputs)
			var wasmOutputsArity = len(importFunction.wasmOutputs)

			var importFunctionInputsCPointer *cWasmerValueTag
			var importFunctionOutputsCPointer *cWasmerValueTag

			if wasmInputsArity > 0 {
				importFunctionInputsCPointer = (*cWasmerValueTag)(unsafe.Pointer(&importFunction.wasmInputs[0]))
			}

			if wasmOutputsArity > 0 {
				importFunctionOutputsCPointer = (*cWasmerValueTag)(unsafe.Pointer(&importFunction.wasmOutputs[0]))
			}

			importFunction.importedFunctionPointer = cWasmerImportFuncNew(
				importFunction.cgoPointer,
				importFunctionInputsCPointer,
				cUint(wasmInputsArity),
				importFunctionOutputsCPointer,
				cUint(wasmOutputsArity),
			)

			var importedFunction = cNewWasmerImportT(
				importFunction.namespace,
				importName,
				importFunction.importedFunctionPointer,
			)

			wasmImports[importFunctionNth] = importedFunction
			importFunctionNth++
		}
	}

	var wasmImportsCPointer *cWasmerImportT

	if numberOfImports > 0 {
		wasmImportsCPointer = (*cWasmerImportT)(unsafe.Pointer(&wasmImports[0]))
	}

	return wasmImportsCPointer, numberOfImports
}

func SetImports(imports *Imports) error {
	wasmImportsCPointer, numberOfImports := generateWasmerImports(imports)

	var result = cWasmerSetImports(
		wasmImportsCPointer,
		cInt(numberOfImports),
	)

	if result != cWasmerOk {
		return newWrappedError(ErrFailedCacheImports)
	}
	return nil
}
