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

func retrieveExportedFunctions(c_instance *cWasmerInstanceT, wasmExports *cWasmerExportsT) (ExportsMap, ExportSignaturesMap, error) {
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

		wrappedWasmFunction, signature, err := createExportedFunctionWrapper(c_instance, wasmFunction, exportedFunctionName)
		if err != nil {
			return nil, nil, err
		}

		exports[exportedFunctionName] = wrappedWasmFunction
		signatures[exportedFunctionName] = signature
	}

	return exports, signatures, nil
}

func createExportedFunctionWrapper(
	c_instance *cWasmerInstanceT,
	wasmFunction *cWasmerExportFuncT,
	exportedFunctionName string,
) (func(...interface{}) (Value, error), *ExportedFunctionSignature, error) {
	wasmFunctionInputSignatures, wasmFunctionInputsArity, err := getExportedFunctionSignature(wasmFunction, exportedFunctionName)
	if err != nil {
		return nil, nil, err
	}

	wasmFunctionOutputsArity, err := getExportedFunctionOutputArity(wasmFunction, exportedFunctionName)
	if err != nil {
		return nil, nil, err
	}

	signature := &ExportedFunctionSignature{
		InputArity:  int(wasmFunctionInputsArity),
		OutputArity: int(wasmFunctionOutputsArity),
	}

	wrapper := func(arguments ...interface{}) (Value, error) {
		err := validateGivenArguments(exportedFunctionName, arguments, wasmFunctionInputsArity)
		if err != nil {
			return Void(), err
		}

		wasmInputs, err := createWasmInputsFromArguments(arguments, wasmFunctionInputsArity, wasmFunctionInputSignatures, exportedFunctionName)
		if err != nil {
			return Void(), err
		}

		wasmOutputs, callResult := callWasmFunction(
			c_instance,
			exportedFunctionName,
			wasmFunctionInputsArity,
			wasmFunctionOutputsArity,
			wasmInputs,
		)

		if callResult != cWasmerOk {
			err := fmt.Errorf("Failed to call the `%s` exported function.", exportedFunctionName)
			return Void(), newWrappedError(err)
		}

		value, err := convertWasmOutputToValue(wasmFunctionOutputsArity, wasmOutputs, exportedFunctionName)
		return value, err
	}
	return wrapper, signature, nil
}
