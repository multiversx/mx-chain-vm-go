package wasmer

import (
	"fmt"
	"unsafe"
)

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
) error {
	err := validateGivenArguments(callInfo.FuncName, arguments, callInfo.InputArity)
	if err != nil {
		return err
	}

	var wasmInputs []cWasmerValueT
	wasmInputs, err = createWasmInputsFromArguments(
		arguments,
		callInfo.InputArity,
		callInfo.InputSignature,
		callInfo.FuncName)
	if err != nil {
		return err
	}

	_, callResult := callWasmFunction(
		cInstance,
		callInfo.FuncName,
		callInfo.InputArity,
		callInfo.OutputArity,
		wasmInputs,
	)

	if callResult != cWasmerOk {
		err = fmt.Errorf("failed to call the `%s` exported function", callInfo.FuncName)
		return newWrappedError(err)
	}

	return nil
}
