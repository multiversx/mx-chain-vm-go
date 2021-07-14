package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_3_manBufNew(void* context);
// extern int32_t 	v1_3_manBufNewFromBytes(void* context, int32_t dataOffset, int32_t dataLength);
// extern int32_t	v1_3_manBufSetBytes(void* context, int32_t manBufHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t 	v1_3_manBufGetLength(void* context, int32_t manBufHandle);
// extern int32_t	v1_3_manBufGetBytes(void* context, int32_t manBufHandle, int32_t resultOffset);
// extern int32_t	v1_3_manBufExtendFromSlice(void* context, int32_t manBufHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t	v1_3_manBufToBigIntUnsigned(void* context, int32_t manBufHandle, int32_t bigIntHandle);
// extern int32_t 	v1_3_manBufToBigIntSigned(void* context, int32_t manBufHandle, int32_t bigIntHandle);
// extern int32_t	v1_3_manBufFromBigIntUnsigned(void* context, int32_t manBufHandle, int32_t bigIntHandle);
// extern int32_t	v1_3_manBufFromBigIntSigned(void* context, int32_t manBufHandle, int32_t bigIntHandle);
// extern int32_t	v1_3_manBufStorageStore(void* context, int32_t keyOffset, int32_t keyLength,int32_t manBufHandle);
// extern int32_t	v1_3_manBufStorageLoad(void* context, int32_t keyOffset, int32_t keyLength, int32_t manBufHandle);
// extern int32_t	v1_3_manBufGetArgument(void* context, int32_t id, int32_t manBufHandle);
// extern int32_t	v1_3_manBufFinish(void* context, int32_t manBufHandle);
import "C"
import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

// ManagedBufferImports creates a new wasmer.Imports populated with the ManagedBuffer API methods
func ManagedBufferImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("manBufNew", v1_3_manBufNew, C.v1_3_manBufNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufNewFromBytes", v1_3_manBufNewFromBytes, C.v1_3_manBufNewFromBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufSetBytes", v1_3_manBufSetBytes, C.v1_3_manBufSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufGetLength", v1_3_manBufGetLength, C.v1_3_manBufGetLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufGetBytes", v1_3_manBufGetBytes, C.v1_3_manBufGetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufExtendFromSlice", v1_3_manBufExtendFromSlice, C.v1_3_manBufExtendFromSlice)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufToBigIntUnsigned", v1_3_manBufToBigIntUnsigned, C.v1_3_manBufToBigIntUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufToBigIntSigned", v1_3_manBufToBigIntSigned, C.v1_3_manBufToBigIntSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufFromBigIntUnsigned", v1_3_manBufFromBigIntUnsigned, C.v1_3_manBufFromBigIntUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufFromBigIntSigned", v1_3_manBufFromBigIntSigned, C.v1_3_manBufFromBigIntSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufStorageStore", v1_3_manBufStorageStore, C.v1_3_manBufStorageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufStorageLoad", v1_3_manBufStorageLoad, C.v1_3_manBufStorageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufGetArgument", v1_3_manBufGetArgument, C.v1_3_manBufGetArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("manBufFinish", v1_3_manBufFinish, C.v1_3_manBufFinish)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_3_manBufNew
func v1_3_manBufNew(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufNew
	metering.UseGas(gasToUse)

	return managedType.NewManagedBuffer()
}

//export v1_3_manBufNewFromBytes
func v1_3_manBufNewFromBytes(context unsafe.Pointer, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufNewFromBytes
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	return managedType.NewManagedBufferFromBytes(data)
}

//export v1_3_manBufSetBytes
func v1_3_manBufSetBytes(context unsafe.Pointer, manBufHandle int32, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufSetBytes
	metering.UseGas(gasToUse)
	managedType.ConsumeGasForThisIntNumberOfBytes(int(dataLength))

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	isSuccess := managedType.SetBytesForThisManagedBuffer(manBufHandle, data)
	if !isSuccess {
		arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_3_manBufGetLength
func v1_3_manBufGetLength(context unsafe.Pointer, manBufHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufGetLength
	metering.UseGas(gasToUse)

	length := managedType.GetLengthForThisManagedBuffer(manBufHandle)
	if length == -1 {
		arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
	}

	return length
}

//export v1_3_manBufGetBytes
func v1_3_manBufGetBytes(context unsafe.Pointer, manBufHandle int32, resultOffset int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufGetBytes
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytesForThisManagedBuffer(manBufHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(managedBuffer))

	err = runtime.MemStore(resultOffset, managedBuffer)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(managedBuffer)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_3_manBufExtendFromSlice
func v1_3_manBufExtendFromSlice(context unsafe.Pointer, manBufHandle int32, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufSetBytes
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	isSuccess := managedType.AppendBytesToThisManagedBuffer(manBufHandle, data)
	if !isSuccess {
		arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_3_manBufToBigIntUnsigned
func v1_3_manBufToBigIntUnsigned(context unsafe.Pointer, manBufHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufToBigIntUnsigned
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytesForThisManagedBuffer(manBufHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	bigInt.SetBytes(managedBuffer)

	return 0
}

//export v1_3_manBufToBigIntSigned
func v1_3_manBufToBigIntSigned(context unsafe.Pointer, manBufHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufToBigIntSigned
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytesForThisManagedBuffer(manBufHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	twos.SetBytes(bigInt, managedBuffer)

	return 0
}

//export v1_3_manBufFromBigIntUnsigned
func v1_3_manBufFromBigIntUnsigned(context unsafe.Pointer, manBufHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufFromBigIntUnsigned
	metering.UseGas(gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}
	bytes := value.Bytes()

	isSuccess := managedType.SetBytesForThisManagedBuffer(manBufHandle, bytes)
	if !isSuccess {
		arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
	}
	return 0
}

//export v1_3_manBufFromBigIntSigned
func v1_3_manBufFromBigIntSigned(context unsafe.Pointer, manBufHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufFromBigIntSigned
	metering.UseGas(gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}
	bytes := twos.ToBytes(value)

	isSuccess := managedType.SetBytesForThisManagedBuffer(manBufHandle, bytes)
	if !isSuccess {
		arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
	}
	return 0
}

//export v1_3_manBufStorageStore
func v1_3_manBufStorageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, manBufHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufStorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	managedBuffer, err := managedType.GetBytesForThisManagedBuffer(manBufHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	_, err = storage.SetStorage(key, managedBuffer)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_3_manBufStorageLoad
func v1_3_manBufStorageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32, manBufHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufStorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	bytes := storage.GetStorage(key)

	isSuccess := managedType.SetBytesForThisManagedBuffer(manBufHandle, bytes)
	if !isSuccess {
		arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	return 0
}

//export v1_3_manBufGetArgument
func v1_3_manBufGetArgument(context unsafe.Pointer, id int32, manBufHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufGetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id {
		return 1
	}

	isSuccess := managedType.SetBytesForThisManagedBuffer(manBufHandle, args[id])
	if !isSuccess {
		arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
	}
	return 0
}

//export v1_3_manBufFinish
func v1_3_manBufFinish(context unsafe.Pointer, manBufHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.ManBufFinish
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytesForThisManagedBuffer(manBufHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	output.Finish(managedBuffer)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(managedBuffer)))
	metering.UseGas(gasToUse)
	return 0
}
