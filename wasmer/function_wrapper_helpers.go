package wasmer

import (
	"fmt"
	"unsafe"
)

func getExportedFunctionSignature(
	wasmFunction *cWasmerExportFuncT,
	exportedFunctionName string,
) ([]cWasmerValueTag, cUint32T, error) {
	var wasmFunctionInputsArity cUint32T
	if cWasmerExportFuncParamsArity(wasmFunction, &wasmFunctionInputsArity) != cWasmerOk {
		return nil, 0, NewExportedFunctionError(exportedFunctionName, "Failed to read the input arity of the `%s` exported function.")
	}

	var wasmFunctionInputSignatures = make([]cWasmerValueTag, int(wasmFunctionInputsArity))

	if wasmFunctionInputsArity > 0 {
		var wasmFunctionInputSignaturesCPointer = (*cWasmerValueTag)(unsafe.Pointer(&wasmFunctionInputSignatures[0]))

		if cWasmerExportFuncParams(wasmFunction, wasmFunctionInputSignaturesCPointer, wasmFunctionInputsArity) != cWasmerOk {
			return nil, 0, NewExportedFunctionError(exportedFunctionName, "Failed to read the signature of the `%s` exported function.")
		}
	}

	return wasmFunctionInputSignatures, wasmFunctionInputsArity, nil
}

func getExportedFunctionOutputArity(
	wasmFunction *cWasmerExportFuncT,
	exportedFunctionName string,
) (cUint32T, error) {
	var wasmFunctionOutputsArity cUint32T
	if cWasmerExportFuncResultsArity(wasmFunction, &wasmFunctionOutputsArity) != cWasmerOk {
		return cUint32T(0), NewExportedFunctionError(exportedFunctionName, "Failed to read the output arity of the `%s` exported function.")
	}
	return wasmFunctionOutputsArity, nil
}

func validateGivenArguments(
	exportedFunctionName string,
	arguments []interface{},
	wasmFunctionInputsArity cUint32T,
) error {
	var numberOfGivenArguments = len(arguments)
	var diff = int(wasmFunctionInputsArity) - numberOfGivenArguments

	if diff > 0 {
		return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Missing %d argument(s) when calling the `%%s` exported function; Expect %d argument(s), given %d.", diff, int(wasmFunctionInputsArity), numberOfGivenArguments))
	} else if diff < 0 {
		return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Given %d extra argument(s) when calling the `%%s` exported function; Expect %d argument(s), given %d.", -diff, int(wasmFunctionInputsArity), numberOfGivenArguments))
	}
	return nil
}

func callWasmFunction(
	c_instance *cWasmerInstanceT,
	exportedFunctionName string,
	wasmFunctionInputsArity cUint32T,
	wasmFunctionOutputsArity cUint32T,
	wasmInputs []cWasmerValueT,
) ([]cWasmerValueT, cWasmerResultT) {
	var wasmFunctionName = cCString(exportedFunctionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	var wasmInputsCPointer *cWasmerValueT
	if wasmFunctionInputsArity > 0 {
		wasmInputsCPointer = (*cWasmerValueT)(unsafe.Pointer(&wasmInputs[0]))
	} else {
		wasmInputsCPointer = (*cWasmerValueT)(unsafe.Pointer(&wasmInputs))
	}

	var wasmOutputs = make([]cWasmerValueT, wasmFunctionOutputsArity)
	var wasmOutputsCPointer *cWasmerValueT
	if wasmFunctionOutputsArity > 0 {
		wasmOutputsCPointer = (*cWasmerValueT)(unsafe.Pointer(&wasmOutputs[0]))
	} else {
		wasmOutputsCPointer = (*cWasmerValueT)(unsafe.Pointer(&wasmOutputs))
	}

	var callResult = cWasmerInstanceCall(
		c_instance,
		wasmFunctionName,
		wasmInputsCPointer,
		wasmFunctionInputsArity,
		wasmOutputsCPointer,
		wasmFunctionOutputsArity,
	)

	return wasmOutputs, callResult
}

func createWasmInputsFromArguments(
	arguments []interface{},
	wasmFunctionInputsArity cUint32T,
	wasmFunctionInputSignatures []cWasmerValueTag,
	exportedFunctionName string,
) ([]cWasmerValueT, error) {
	var err error
	var wasmInputs = make([]cWasmerValueT, wasmFunctionInputsArity)
	for index, value := range arguments {
		var wasmInputType = wasmFunctionInputSignatures[index]

		switch wasmInputType {
		case cWasmI32:
			err = writeInt32ToWasmInputs(wasmInputs, index, value, exportedFunctionName)
			if err != nil {
				return nil, err
			}
		case cWasmI64:
			err = writeInt64ToWasmInputs(wasmInputs, index, value, exportedFunctionName)
			if err != nil {
				return nil, err
			}
		case cWasmF32:
			err = writeFloat32ToWasmInputs(wasmInputs, index, value, exportedFunctionName)
			if err != nil {
				return nil, err
			}
		case cWasmF64:
			err = writeFloat64ToWasmInputs(wasmInputs, index, value, exportedFunctionName)
			if err != nil {
				return nil, err
			}
		default:
			return nil, NewExportedFunctionError(exportedFunctionName, "Invalid arguments type when calling the `%s` exported function.")
		}
	}

	return wasmInputs, nil
}

func convertWasmOutputToValue(
	wasmFunctionOutputsArity cUint32T,
	wasmOutputs []cWasmerValueT,
	exportedFunctionName string,
) (Value, error) {
	if wasmFunctionOutputsArity > 0 {
		var result = wasmOutputs[0]

		switch result.tag {
		case cWasmI32:
			pointer := (*int32)(unsafe.Pointer(&result.value))

			return I32(*pointer), nil
		case cWasmI64:
			pointer := (*int64)(unsafe.Pointer(&result.value))

			return I64(*pointer), nil
		case cWasmF32:
			pointer := (*float32)(unsafe.Pointer(&result.value))

			return F32(*pointer), nil
		case cWasmF64:
			pointer := (*float64)(unsafe.Pointer(&result.value))

			return F64(*pointer), nil
		default:
			return Void(), NewExportedFunctionError(exportedFunctionName, "Invalid output type retrieved from function `%s`.")
		}
	}

	return Void(), nil
}

func writeInt32ToWasmInputs(wasmInputs []cWasmerValueT, index int, value interface{}, exportedFunctionName string) error {
	wasmInputs[index].tag = cWasmI32
	var pointer = (*int32)(unsafe.Pointer(&wasmInputs[index].value))

	switch value.(type) {
	case int8:
		*pointer = int32(value.(int8))
	case uint8:
		*pointer = int32(value.(uint8))
	case int16:
		*pointer = int32(value.(int16))
	case uint16:
		*pointer = int32(value.(uint16))
	case int32:
		*pointer = int32(value.(int32))
	case int:
		*pointer = int32(value.(int))
	case uint:
		*pointer = int32(value.(uint))
	case Value:
		var value = value.(Value)

		if value.GetType() != TypeI32 {
			return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `i32`, cannot cast given value to this type.", index+1))
		}

		*pointer = value.ToI32()
	default:
		return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `i32`, cannot cast given value to this type.", index+1))
	}

	return nil
}

func writeInt64ToWasmInputs(wasmInputs []cWasmerValueT, index int, value interface{}, exportedFunctionName string) error {
	wasmInputs[index].tag = cWasmI64
	var pointer = (*int64)(unsafe.Pointer(&wasmInputs[index].value))

	switch value.(type) {
	case int8:
		*pointer = int64(value.(int8))
	case uint8:
		*pointer = int64(value.(uint8))
	case int16:
		*pointer = int64(value.(int16))
	case uint16:
		*pointer = int64(value.(uint16))
	case int32:
		*pointer = int64(value.(int32))
	case uint32:
		*pointer = int64(value.(uint32))
	case int64:
		*pointer = int64(value.(int64))
	case int:
		*pointer = int64(value.(int))
	case uint:
		*pointer = int64(value.(uint))
	case Value:
		var value = value.(Value)

		if value.GetType() != TypeI64 {
			return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `i64`, cannot cast given value to this type.", index+1))
		}

		*pointer = value.ToI64()
	default:
		return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `i64`, cannot cast given value to this type.", index+1))
	}

	return nil
}

func writeFloat32ToWasmInputs(wasmInputs []cWasmerValueT, index int, value interface{}, exportedFunctionName string) error {
	wasmInputs[index].tag = cWasmF32
	var pointer = (*float32)(unsafe.Pointer(&wasmInputs[index].value))

	switch value.(type) {
	case float32:
		*pointer = value.(float32)
	case Value:
		var value = value.(Value)

		if value.GetType() != TypeF32 {
			return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `f32`, cannot cast given value to this type.", index+1))
		}

		*pointer = value.ToF32()
	default:
		return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `f32`, cannot cast given value to this type.", index+1))
	}

	return nil
}

func writeFloat64ToWasmInputs(wasmInputs []cWasmerValueT, index int, value interface{}, exportedFunctionName string) error {
	wasmInputs[index].tag = cWasmF64
	var pointer = (*float64)(unsafe.Pointer(&wasmInputs[index].value))

	switch value.(type) {
	case float32:
		*pointer = float64(value.(float32))
	case float64:
		*pointer = value.(float64)
	case Value:
		var value = value.(Value)

		if value.GetType() != TypeF64 {
			return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `f64`, cannot cast given value to this type.", index+1))
		}

		*pointer = value.ToF64()
	default:
		return NewExportedFunctionError(exportedFunctionName, fmt.Sprintf("Argument #%d of the `%%s` exported function must be of type `f64`, cannot cast given value to this type.", index+1))
	}

	return nil
}
