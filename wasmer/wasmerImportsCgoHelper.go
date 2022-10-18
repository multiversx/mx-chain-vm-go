package wasmer

import (
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

func getVMHooksFromContextRawPtr(contextPtr unsafe.Pointer) executor.VMHooks {
	instCtx := IntoInstanceContext(contextPtr)
	//ctx := (*cWasmerInstanceContextT)(contextPtr)
	//data := cWasmerInstanceContextDataGet(ctx)
	data := *(*uintptr)(instCtx.Data())
	//fmt.Printf("instance context data <------- %x\n", instCtx.Data())
	//fmt.Printf("resolved <-------- %x\n", unsafe.Pointer(data))
	return *(*executor.VMHooks)(unsafe.Pointer(data))
}

func injectCgoFunctionPointers() (vmcommon.FunctionNames, error) {
	imports := newWasmerImports()
	populateWasmerImports(imports)
	wasmImportsCPointer, numberOfImports := generateWasmerImports(imports)

	var result = cWasmerCacheImportObjectFromImports(
		wasmImportsCPointer,
		cInt(numberOfImports),
	)

	if result != cWasmerOk {
		return nil, newWrappedError(ErrFailedCacheImports)
	}

	return extractImportNames(imports), nil
}

func extractImportNames(imports *wasmerImports) vmcommon.FunctionNames {
	names := make(vmcommon.FunctionNames)
	var empty struct{}
	for _, env := range imports.imports {
		for name := range env {
			names[name] = empty
		}
	}
	return names
}

func generateWasmerImports(imports *wasmerImports) (*cWasmerImportT, int) {
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
