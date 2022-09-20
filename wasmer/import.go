package wasmer

import (
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executorinterface"
)

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

// Namespace changes the current namespace of the next imported functions.
func (imports *Imports) Namespace(namespace string) *Imports {
	imports.currentNamespace = namespace

	return imports
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

func convertArgType(argType executorinterface.EIFunctionValue) cWasmerValueTag {
	switch argType {
	case executorinterface.EIFunctionValueInt32:
		return cWasmI32
	case executorinterface.EIFunctionValueInt64:
		return cWasmI64
	}
	return cWasmI32 // unreachable, but might consider adding an error
}

func ConvertImports(eiFunctions *executorinterface.EIFunctions) *Imports {
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

// InstanceContext represents a way to access instance API from within
// an imported context.
type InstanceContext struct {
	context *cWasmerInstanceContextT
	memory  MemoryHandler
}

// NewInstanceContext creates a new wasmer context given a cWasmerInstance and a memory
func NewInstanceContext(ctx *cWasmerInstanceContextT, mem Memory) *InstanceContext {
	return &InstanceContext{
		context: ctx,
		memory:  &mem,
	}
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
func (instanceContext *InstanceContext) Memory() MemoryHandler {
	return instanceContext.memory
}

// Data returns the instance context data as an `unsafe.Pointer`. It's
// up to the user to cast it appropriately as a pointer to a data.
func (instanceContext *InstanceContext) Data() unsafe.Pointer {
	return cWasmerInstanceContextDataGet(instanceContext.context)
}
