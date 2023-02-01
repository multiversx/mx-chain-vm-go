package wasmer

import (
	"fmt"
	"unsafe"
)

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

func retrieveExportedMemory(wasmExports *cWasmerExportsT) (Memory, bool, error) {
	var numberOfExports = int(cWasmerExportsLen(wasmExports))

	var memory Memory
	var hasMemory = false

	for nth := 0; nth < numberOfExports; nth++ {
		var wasmExport = cWasmerExportsGet(wasmExports, cInt(nth))
		var wasmExportKind = cWasmerExportKind(wasmExport)

		if wasmExportKind == cWasmMemory {
			var wasmMemory *cWasmerMemoryT

			if cWasmerExportToMemory(wasmExport, &wasmMemory) != cWasmerOk {
				var emptyMemory Memory
				return emptyMemory, false, NewInstanceError("Failed to extract the exported memory.")
			}

			memory = newMemory(wasmMemory)
			hasMemory = true
		}
	}

	return memory, hasMemory, nil
}

func retrieveExportedFunctions(
	cInstance *cWasmerInstanceT,
	wasmExports *cWasmerExportsT,
) (ExportsMap, ExportSignaturesMap, error) {
	var exports = make(ExportsMap)
	var signatures = make(ExportSignaturesMap)

	var numberOfExports = int(cWasmerExportsLen(wasmExports))

	for nth := 0; nth < numberOfExports; nth++ {
		var wasmExport = cWasmerExportsGet(wasmExports, cInt(nth))
		var wasmExportKind = cWasmerExportKind(wasmExport)

		if wasmExportKind != cWasmFunction {
			continue
		}

		var wasmExportName = cWasmerExportName(wasmExport)
		var wasmFunction = cWasmerExportToFunc(wasmExport)
		var exportedFunctionName = cGoStringN((*cChar)(unsafe.Pointer(wasmExportName.bytes)), (cInt)(wasmExportName.bytes_len))

		callInfo, err := createExportedFunctionCallInfo(wasmFunction, exportedFunctionName)
		if err != nil {
			return nil, nil, err
		}

		signature := &ExportedFunctionSignature{
			InputArity:  int(callInfo.InputArity),
			OutputArity: int(callInfo.OutputArity),
		}

		exports[exportedFunctionName] = callInfo
		signatures[exportedFunctionName] = signature
	}

	return exports, signatures, nil
}

func createExportedFunctionCallInfo(
	wasmFunction *cWasmerExportFuncT,
	exportedFunctionName string,
) (*ExportedFunctionCallInfo, error) {

	wasmFunctionInputSignatures, wasmFunctionInputsArity, err := getExportedFunctionSignature(wasmFunction, exportedFunctionName)
	if err != nil {
		return nil, err
	}

	wasmFunctionOutputsArity, err := getExportedFunctionOutputArity(wasmFunction, exportedFunctionName)
	if err != nil {
		return nil, err
	}

	callInfo := &ExportedFunctionCallInfo{
		FuncName:       exportedFunctionName,
		InputArity:     wasmFunctionInputsArity,
		InputSignature: wasmFunctionInputSignatures,
		OutputArity:    wasmFunctionOutputsArity,
	}

	return callInfo, nil
}

func callExportedFunction(
	cInstance *cWasmerInstanceT,
	callInfo *ExportedFunctionCallInfo,
	arguments ...interface{},
) (Value, error) {
	err := validateGivenArguments(callInfo.FuncName, arguments, callInfo.InputArity)
	if err != nil {
		return Void(), err
	}

	var wasmInputs []cWasmerValueT
	wasmInputs, err = createWasmInputsFromArguments(
		arguments,
		callInfo.InputArity,
		callInfo.InputSignature,
		callInfo.FuncName)
	if err != nil {
		return Void(), err
	}

	wasmOutputs, callResult := callWasmFunction(
		cInstance,
		callInfo.FuncName,
		callInfo.InputArity,
		callInfo.OutputArity,
		wasmInputs,
	)

	if callResult != cWasmerOk {
		err = fmt.Errorf("failed to call the `%s` exported function", callInfo.FuncName)
		return Void(), newWrappedError(err)
	}

	var value Value
	value, err = convertWasmOutputToValue(callInfo.OutputArity, wasmOutputs, callInfo.FuncName)
	return value, err
}
