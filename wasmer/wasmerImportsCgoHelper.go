package wasmer

import (
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

func importsInterfaceFromRaw(contextPtr unsafe.Pointer) executor.ImportsInterface {
	// instCtx := wasmer.IntoInstanceContext(vmHostPtr)
	// var ptr = *(*uintptr)(context)
	return *(*executor.ImportsInterface)(contextPtr)
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
