package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_4_mBufferNew(void* context);
// extern int32_t 	v1_4_mBufferNewFromBytes(void* context, int32_t dataOffset, int32_t dataLength);
// extern int32_t	v1_4_mBufferSetRandom(void* context, int32_t mBufferHandle, int32_t length);
// extern int32_t	v1_4_mBufferSetBytes(void* context, int32_t mBufferHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t 	v1_4_mBufferGetLength(void* context, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferGetBytes(void* context, int32_t mBufferHandle, int32_t resultOffset);
// extern int32_t	v1_4_mBufferAppend(void* context, int32_t mBufferHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t	v1_4_mBufferToBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t 	v1_4_mBufferToBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferFromBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferFromBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferStorageStore(void* context, int32_t keyHandle ,int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferStorageLoad(void* context, int32_t keyHandle, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferGetArgument(void* context, int32_t id, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferFinish(void* context, int32_t mBufferHandle);
import "C"
import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

// ManagedBufferImports creates a new wasmer.Imports populated with the ManagedBuffer API methods
func ManagedBufferImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("mBufferNew", v1_4_mBufferNew, C.v1_4_mBufferNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferNewFromBytes", v1_4_mBufferNewFromBytes, C.v1_4_mBufferNewFromBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferSetRandom", v1_4_mBufferSetRandom, C.v1_4_mBufferSetRandom)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferSetBytes", v1_4_mBufferSetBytes, C.v1_4_mBufferSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferGetLength", v1_4_mBufferGetLength, C.v1_4_mBufferGetLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferGetBytes", v1_4_mBufferGetBytes, C.v1_4_mBufferGetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferAppend", v1_4_mBufferAppend, C.v1_4_mBufferAppend)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferToBigIntUnsigned", v1_4_mBufferToBigIntUnsigned, C.v1_4_mBufferToBigIntUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferToBigIntSigned", v1_4_mBufferToBigIntSigned, C.v1_4_mBufferToBigIntSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferFromBigIntUnsigned", v1_4_mBufferFromBigIntUnsigned, C.v1_4_mBufferFromBigIntUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferFromBigIntSigned", v1_4_mBufferFromBigIntSigned, C.v1_4_mBufferFromBigIntSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferStorageStore", v1_4_mBufferStorageStore, C.v1_4_mBufferStorageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferStorageLoad", v1_4_mBufferStorageLoad, C.v1_4_mBufferStorageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferGetArgument", v1_4_mBufferGetArgument, C.v1_4_mBufferGetArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferFinish", v1_4_mBufferFinish, C.v1_4_mBufferFinish)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_4_mBufferNew
func v1_4_mBufferNew(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNew
	metering.UseGas(gasToUse)

	return managedType.NewManagedBuffer()
}

//export v1_4_mBufferNewFromBytes
func v1_4_mBufferNewFromBytes(context unsafe.Pointer, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNewFromBytes
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	return managedType.NewManagedBufferFromBytes(data)
}

//export v1_4_mBufferSetBytes
func v1_4_mBufferSetBytes(context unsafe.Pointer, mBufferHandle int32, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	metering.UseGas(gasToUse)
	managedType.ConsumeGasForThisIntNumberOfBytes(int(dataLength))

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(mBufferHandle, data)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_4_mBufferGetLength
func v1_4_mBufferGetLength(context unsafe.Pointer, mBufferHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetLength
	metering.UseGas(gasToUse)

	length := managedType.GetLength(mBufferHandle)
	if length == -1 {
		_ = arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return -1
	}

	return length
}

//export v1_4_mBufferGetBytes
func v1_4_mBufferGetBytes(context unsafe.Pointer, mBufferHandle int32, resultOffset int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetBytes
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
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

//export v1_4_mBufferAppend
func v1_4_mBufferAppend(context unsafe.Pointer, mBufferHandle int32, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	isSuccess := managedType.AppendBytes(mBufferHandle, data)
	if !isSuccess {
		_ = arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_4_mBufferToBigIntUnsigned
func v1_4_mBufferToBigIntUnsigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigIntUnsigned
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	bigInt.SetBytes(managedBuffer)

	return 0
}

//export v1_4_mBufferToBigIntSigned
func v1_4_mBufferToBigIntSigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigIntSigned
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	twos.SetBytes(bigInt, managedBuffer)

	return 0
}

//export v1_4_mBufferFromBigIntUnsigned
func v1_4_mBufferFromBigIntUnsigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigIntUnsigned
	metering.UseGas(gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}
	bytes := value.Bytes()

	managedType.SetBytes(mBufferHandle, bytes)

	return 0
}

//export v1_4_mBufferFromBigIntSigned
func v1_4_mBufferFromBigIntSigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigIntSigned
	metering.UseGas(gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}
	bytes := twos.ToBytes(value)

	managedType.SetBytes(mBufferHandle, bytes)
	return 0
}

//export v1_4_mBufferStorageStore
func v1_4_mBufferStorageStore(context unsafe.Pointer, keyHandle int32, mBufferHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferStorageStore
	metering.UseGas(gasToUse)

	key, err := managedType.GetBytes(keyHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	_, err = storage.SetStorage(key, managedBuffer)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_mBufferStorageLoad
func v1_4_mBufferStorageLoad(context unsafe.Pointer, keyHandle int32, mBufferHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferStorageLoad
	metering.UseGas(gasToUse)

	key, err := managedType.GetBytes(keyHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	bytes := storage.GetStorage(key)

	managedType.SetBytes(mBufferHandle, bytes)

	return 0
}

//export v1_4_mBufferGetArgument
func v1_4_mBufferGetArgument(context unsafe.Pointer, id int32, mBufferHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id {
		return 1
	}
	managedType.SetBytes(mBufferHandle, args[id])
	return 0
}

//export v1_4_mBufferFinish
func v1_4_mBufferFinish(context unsafe.Pointer, mBufferHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFinish
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	output.Finish(managedBuffer)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(managedBuffer)))
	metering.UseGas(gasToUse)
	return 0
}

//export v1_4_mBufferSetRandom
func v1_4_mBufferSetRandom(context unsafe.Pointer, mBufferHandle int32, length int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	if length < 1 {
		_ = arwen.WithFault(arwen.ErrLengthOfBufferNotCorrect, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	baseGasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetRandom
	lengthDependentGasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	metering.UseGas(math.AddUint64(baseGasToUse, lengthDependentGasToUse))

	randomizer := managedType.GetRandReader()
	buffer := make([]byte, length)
	_, err := randomizer.Read(buffer)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	managedType.SetBytes(mBufferHandle, buffer)
	return 0
}
