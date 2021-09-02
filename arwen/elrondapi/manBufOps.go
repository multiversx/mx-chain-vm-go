package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_4_mBufferNew(void* context);
// extern int32_t 	v1_4_mBufferNewFromBytes(void* context, int32_t dataOffset, int32_t dataLength);
// extern int32_t 	v1_4_mBufferGetLength(void* context, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferGetBytes(void* context, int32_t mBufferHandle, int32_t resultOffset);
// extern int32_t	v1_4_mBufferGetByteSlice(void* context, int32_t sourceHandle, int32_t startingPosition, int32_t sliceLength, int32_t resultOffset);
// extern int32_t	v1_4_mBufferCopyByteSlice(void* context, int32_t sourceHandle, int32_t startingPosition, int32_t sliceLength, int32_t destinationHandle);
// extern int32_t	v1_4_mBufferSetBytes(void* context, int32_t mBufferHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t	v1_4_mBufferAppend(void* context, int32_t accumulatorHandle, int32_t dataHandle);
// extern int32_t	v1_4_mBufferAppendBytes(void* context, int32_t accumulatorHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t	v1_4_mBufferToBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t 	v1_4_mBufferToBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferFromBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferFromBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferToBigFloat(void* context, int32_t mBufferHandle, int32_t bigFloatHandle);
// extern int32_t	v1_4_mBufferFromBigFloat(void* context, int32_t mBufferHandle, int32_t bigFloatHandle);
// extern int32_t	v1_4_mBufferStorageStore(void* context, int32_t keyHandle ,int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferStorageLoad(void* context, int32_t keyHandle, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferGetArgument(void* context, int32_t id, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferFinish(void* context, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferSetRandom(void* context, int32_t destinationHandle, int32_t length);
import "C"
import (
	"math/big"
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

	imports, err = imports.Append("mBufferGetLength", v1_4_mBufferGetLength, C.v1_4_mBufferGetLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferSetBytes", v1_4_mBufferSetBytes, C.v1_4_mBufferSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferGetBytes", v1_4_mBufferGetBytes, C.v1_4_mBufferGetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferGetByteSlice", v1_4_mBufferGetByteSlice, C.v1_4_mBufferGetByteSlice)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferCopyByteSlice", v1_4_mBufferCopyByteSlice, C.v1_4_mBufferCopyByteSlice)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferSetBytes", v1_4_mBufferSetBytes, C.v1_4_mBufferSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferAppend", v1_4_mBufferAppend, C.v1_4_mBufferAppend)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("mBufferAppendBytes", v1_4_mBufferAppendBytes, C.v1_4_mBufferAppendBytes)
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

	imports, err = imports.Append("mBufferToBigFloat", v1_4_mBufferToBigFloat, C.v1_4_mBufferToBigFloat)
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

	imports, err = imports.Append("mBufferSetRandom", v1_4_mBufferSetRandom, C.v1_4_mBufferSetRandom)
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

	mBufferBytes, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(mBufferBytes))

	err = runtime.MemStore(resultOffset, mBufferBytes)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(mBufferBytes)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_4_mBufferGetByteSlice
func v1_4_mBufferGetByteSlice(context unsafe.Pointer, sourceHandle int32, startingPosition int32, sliceLength int32, resultOffset int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetByteSlice
	metering.UseGas(gasToUse)

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(sourceBytes))

	if startingPosition < 0 || sliceLength < 0 || int(startingPosition+sliceLength) > len(sourceBytes) {
		// does not fail execution if slice exceeds bounds
		return 1
	}

	slice := sourceBytes[startingPosition : startingPosition+sliceLength]
	err = runtime.MemStore(resultOffset, slice)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(sourceBytes)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_4_mBufferCopyByteSlice
func v1_4_mBufferCopyByteSlice(context unsafe.Pointer, sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferCopyByteSlice
	metering.UseGas(gasToUse)

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(sourceBytes))

	if startingPosition < 0 || sliceLength < 0 || int(startingPosition+sliceLength) > len(sourceBytes) {
		// does not fail execution if slice exceeds bounds
		return 1
	}

	slice := sourceBytes[startingPosition : startingPosition+sliceLength]
	managedType.SetBytes(destinationHandle, slice)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(slice)))
	metering.UseGas(gasToUse)

	return 0
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

//export v1_4_mBufferAppend
func v1_4_mBufferAppend(context unsafe.Pointer, accumulatorHandle int32, dataHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppend
	metering.UseGas(gasToUse)

	dataBufferBytes, err := managedType.GetBytes(dataHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(dataBufferBytes))

	isSuccess := managedType.AppendBytes(accumulatorHandle, dataBufferBytes)
	if !isSuccess {
		_ = arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(dataBufferBytes)))
	metering.UseGas(gasToUse)

	return 0
}

//export v1_4_mBufferAppendBytes
func v1_4_mBufferAppendBytes(context unsafe.Pointer, accumulatorHandle int32, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppendBytes
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	isSuccess := managedType.AppendBytes(accumulatorHandle, data)
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

//export v1_4_mBufferToBigFloat
func v1_4_mBufferToBigFloat(context unsafe.Pointer, mBufferHandle, bigFloatHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigFloat
	metering.UseGas(gasToUse)

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	value, err := managedType.GetBigFloatOrCreate(bigFloatHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigFloat := new(big.Float)
	err = bigFloat.GobDecode(managedBuffer)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	if bigFloat.Prec() != managedType.GetBigFloatPrecision() {
		_ = arwen.WithFault(arwen.ErrBigFloatWrongPrecision, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return 1
	}
	if bigFloat.IsInf() {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return 1
	}

	value.Set(bigFloat)
	return 0
}

//export v1_4_mBufferFromBigFloat
func v1_4_mBufferFromBigFloat(context unsafe.Pointer, mBufferHandle, bigFloatHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigFloat
	metering.UseGas(gasToUse)

	value, err := managedType.GetBigFloat(bigFloatHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return 1
	}

	encodedFloat, err := value.GobEncode()
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(mBufferHandle, encodedFloat)

	return 0
}

//export v1_4_mBufferStorageStore
func v1_4_mBufferStorageStore(context unsafe.Pointer, keyHandle int32, sourceHandle int32) int32 {
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

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	_, err = storage.SetStorage(key, sourceBytes)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_mBufferStorageLoad
func v1_4_mBufferStorageLoad(context unsafe.Pointer, keyHandle int32, destinationHandle int32) int32 {
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

	managedType.SetBytes(destinationHandle, bytes)

	return 0
}

//export v1_4_mBufferGetArgument
func v1_4_mBufferGetArgument(context unsafe.Pointer, id int32, destinationHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id {
		return 1
	}
	managedType.SetBytes(destinationHandle, args[id])
	return 0
}

//export v1_4_mBufferFinish
func v1_4_mBufferFinish(context unsafe.Pointer, sourceHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFinish
	metering.UseGas(gasToUse)

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	output.Finish(sourceBytes)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(sourceBytes)))
	metering.UseGas(gasToUse)
	return 0
}

//export v1_4_mBufferSetRandom
func v1_4_mBufferSetRandom(context unsafe.Pointer, destinationHandle int32, length int32) int32 {
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

	managedType.SetBytes(destinationHandle, buffer)
	return 0
}
