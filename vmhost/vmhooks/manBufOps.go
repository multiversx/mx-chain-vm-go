package vmhooks

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_4_mBufferNew(void* context);
// extern int32_t 	v1_4_mBufferNewFromBytes(void* context, int32_t dataOffset, int32_t dataLength);
//
// extern int32_t 	v1_4_mBufferGetLength(void* context, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferGetBytes(void* context, int32_t mBufferHandle, int32_t resultOffset);
// extern int32_t	v1_4_mBufferGetByteSlice(void* context, int32_t sourceHandle, int32_t startingPosition, int32_t sliceLength, int32_t resultOffset);
// extern int32_t	v1_4_mBufferCopyByteSlice(void* context, int32_t sourceHandle, int32_t startingPosition, int32_t sliceLength, int32_t destinationHandle);
// extern int32_t	v1_4_mBufferEq(void* context, int32_t mBufferHandle1, int32_t mBufferHandle2);
//
// extern int32_t	v1_4_mBufferSetBytes(void* context, int32_t mBufferHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t v1_4_mBufferSetByteSlice(void* context, int32_t mBufferHandle, int32_t startingPosition, int32_t dataLength, int32_t dataOffset);
// extern int32_t	v1_4_mBufferAppend(void* context, int32_t accumulatorHandle, int32_t dataHandle);
// extern int32_t	v1_4_mBufferAppendBytes(void* context, int32_t accumulatorHandle, int32_t dataOffset, int32_t dataLength);
//
// extern int32_t	v1_4_mBufferToBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t 	v1_4_mBufferToBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferFromBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferFromBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t	v1_4_mBufferToBigFloat(void* context, int32_t mBufferHandle, int32_t bigFloatHandle);
// extern int32_t	v1_4_mBufferFromBigFloat(void* context, int32_t mBufferHandle, int32_t bigFloatHandle);
//
// extern int32_t	v1_4_mBufferStorageStore(void* context, int32_t keyHandle ,int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferStorageLoad(void* context, int32_t keyHandle, int32_t mBufferHandle);
// extern void  	v1_4_mBufferStorageLoadFromAddress(void* context, int32_t addressHandle, int32_t keyHandle, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferGetArgument(void* context, int32_t id, int32_t mBufferHandle);
// extern int32_t	v1_4_mBufferFinish(void* context, int32_t mBufferHandle);
//
// extern int32_t	v1_4_mBufferSetRandom(void* context, int32_t destinationHandle, int32_t length);
import "C"
import (
	"bytes"
	"math/big"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooksmeta"
	"github.com/multiversx/mx-chain-vm-v1_4-go/math"
	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

const (
	mBufferNewName                = "mBufferNew"
	mBufferNewFromBytesName       = "mBufferNewFromBytes"
	mBufferGetLengthName          = "mBufferGetLength"
	mBufferGetBytesName           = "mBufferGetBytes"
	mBufferGetByteSliceName       = "mBufferGetByteSlice"
	mBufferCopyByteSliceName      = "mBufferCopyByteSlice"
	mBufferEqName                 = "mBufferEq"
	mBufferSetBytesName           = "mBufferSetBytes"
	mBufferAppendName             = "mBufferAppend"
	mBufferAppendBytesName        = "mBufferAppendBytes"
	mBufferToBigIntUnsignedName   = "mBufferToBigIntUnsigned"
	mBufferToBigIntSignedName     = "mBufferToBigIntSigned"
	mBufferFromBigIntUnsignedName = "mBufferFromBigIntUnsigned"
	mBufferFromBigIntSignedName   = "mBufferFromBigIntSigned"
	mBufferStorageStoreName       = "mBufferStorageStore"
	mBufferStorageLoadName        = "mBufferStorageLoad"
	mBufferGetArgumentName        = "mBufferGetArgument"
	mBufferFinishName             = "mBufferFinish"
	mBufferSetRandomName          = "mBufferSetRandom"
	mBufferToBigFloatName         = "mBufferToBigFloat"
	mBufferFromBigFloatName       = "mBufferFromBigFloat"
)

// ManagedBufferImports creates a new wasmer.Imports populated with the ManagedBuffer API methods
func ManagedBufferImports(imports vmhooksmeta.EIFunctionReceiver) error {
	imports.Namespace("env")

	err := imports.Append("mBufferNew", v1_4_mBufferNew, C.v1_4_mBufferNew)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferNewFromBytes", v1_4_mBufferNewFromBytes, C.v1_4_mBufferNewFromBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetLength", v1_4_mBufferGetLength, C.v1_4_mBufferGetLength)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferSetBytes", v1_4_mBufferSetBytes, C.v1_4_mBufferSetBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferSetByteSlice", v1_4_mBufferSetByteSlice, C.v1_4_mBufferSetByteSlice)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetBytes", v1_4_mBufferGetBytes, C.v1_4_mBufferGetBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetByteSlice", v1_4_mBufferGetByteSlice, C.v1_4_mBufferGetByteSlice)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferCopyByteSlice", v1_4_mBufferCopyByteSlice, C.v1_4_mBufferCopyByteSlice)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferEq", v1_4_mBufferEq, C.v1_4_mBufferEq)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferAppend", v1_4_mBufferAppend, C.v1_4_mBufferAppend)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferAppendBytes", v1_4_mBufferAppendBytes, C.v1_4_mBufferAppendBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferToBigIntUnsigned", v1_4_mBufferToBigIntUnsigned, C.v1_4_mBufferToBigIntUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferToBigIntSigned", v1_4_mBufferToBigIntSigned, C.v1_4_mBufferToBigIntSigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFromBigIntUnsigned", v1_4_mBufferFromBigIntUnsigned, C.v1_4_mBufferFromBigIntUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFromBigIntSigned", v1_4_mBufferFromBigIntSigned, C.v1_4_mBufferFromBigIntSigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferToBigFloat", v1_4_mBufferToBigFloat, C.v1_4_mBufferToBigFloat)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFromBigFloat", v1_4_mBufferFromBigFloat, C.v1_4_mBufferFromBigFloat)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferStorageStore", v1_4_mBufferStorageStore, C.v1_4_mBufferStorageStore)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferStorageLoad", v1_4_mBufferStorageLoad, C.v1_4_mBufferStorageLoad)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferStorageLoadFromAddress", v1_4_mBufferStorageLoadFromAddress, C.v1_4_mBufferStorageLoadFromAddress)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetArgument", v1_4_mBufferGetArgument, C.v1_4_mBufferGetArgument)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFinish", v1_4_mBufferFinish, C.v1_4_mBufferFinish)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferSetRandom", v1_4_mBufferSetRandom, C.v1_4_mBufferSetRandom)
	if err != nil {
		return err
	}

	return nil
}

//export v1_4_mBufferNew
func v1_4_mBufferNew(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNew
	metering.UseGasAndAddTracedGas(mBufferNewName, gasToUse)

	return managedType.NewManagedBuffer()
}

//export v1_4_mBufferNewFromBytes
func v1_4_mBufferNewFromBytes(context unsafe.Pointer, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNewFromBytes
	metering.UseGasAndAddTracedGas(mBufferNewFromBytesName, gasToUse)

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
	metering.UseGasAndAddTracedGas(mBufferGetLengthName, gasToUse)

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
	metering.StartGasTracing(mBufferGetBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetBytes
	metering.UseAndTraceGas(gasToUse)

	mBufferBytes, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(mBufferBytes)

	err = runtime.MemStore(resultOffset, mBufferBytes)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	storage := arwen.GetStorageContext(context)
	if !storage.IsUseDifferentGasCostFlagSet() {
		gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(mBufferBytes)))
		metering.UseAndTraceGas(gasToUse)
	}

	return 0
}

//export v1_4_mBufferGetByteSlice
func v1_4_mBufferGetByteSlice(context unsafe.Pointer, sourceHandle int32, startingPosition int32, sliceLength int32, resultOffset int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(mBufferGetByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetByteSlice
	metering.UseAndTraceGas(gasToUse)

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sourceBytes)

	if startingPosition < 0 || sliceLength < 0 || int(startingPosition+sliceLength) > len(sourceBytes) {
		// does not fail execution if slice exceeds bounds
		return 1
	}

	slice := sourceBytes[startingPosition : startingPosition+sliceLength]
	err = runtime.MemStore(resultOffset, slice)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	storage := arwen.GetStorageContext(context)
	if !storage.IsUseDifferentGasCostFlagSet() {
		gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(sourceBytes)))
		metering.UseAndTraceGas(gasToUse)
	}

	return 0
}

//export v1_4_mBufferCopyByteSlice
func v1_4_mBufferCopyByteSlice(context unsafe.Pointer, sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) int32 {
	host := arwen.GetVMHost(context)
	return ManagedBufferCopyByteSliceWithHost(host, sourceHandle, startingPosition, sliceLength, destinationHandle)
}

func ManagedBufferCopyByteSliceWithHost(host arwen.VMHost, sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) int32 {
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(mBufferCopyByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferCopyByteSlice
	metering.UseAndTraceGas(gasToUse)

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sourceBytes)

	if startingPosition < 0 || sliceLength < 0 || int(startingPosition+sliceLength) > len(sourceBytes) {
		// does not fail execution if slice exceeds bounds
		return 1
	}

	slice := sourceBytes[startingPosition : startingPosition+sliceLength]
	managedType.SetBytes(destinationHandle, slice)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(slice)))
	metering.UseAndTraceGas(gasToUse)

	return 0
}

//export v1_4_mBufferEq
func v1_4_mBufferEq(context unsafe.Pointer, mBufferHandle1 int32, mBufferHandle2 int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(mBufferEqName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferCopyByteSlice
	metering.UseAndTraceGas(gasToUse)

	bytes1, err := managedType.GetBytes(mBufferHandle1)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}
	managedType.ConsumeGasForBytes(bytes1)

	bytes2, err := managedType.GetBytes(mBufferHandle2)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}
	managedType.ConsumeGasForBytes(bytes2)

	if bytes.Equal(bytes1, bytes2) {
		return 1
	}

	return 0
}

//export v1_4_mBufferSetBytes
func v1_4_mBufferSetBytes(context unsafe.Pointer, mBufferHandle int32, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(mBufferSetBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(data)
	managedType.SetBytes(mBufferHandle, data)

	storage := arwen.GetStorageContext(context)
	if !storage.IsUseDifferentGasCostFlagSet() {
		gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
		metering.UseAndTraceGas(gasToUse)
	}

	return 0
}

//export v1_4_mBufferSetByteSlice
func v1_4_mBufferSetByteSlice(context unsafe.Pointer, mBufferHandle int32, startingPosition int32, dataLength int32, dataOffset int32) int32 {
	host := arwen.GetVMHost(context)
	return ManagedBufferSetByteSliceWithHost(host, mBufferHandle, startingPosition, dataLength, dataOffset)
}

func ManagedBufferSetByteSliceWithHost(host arwen.VMHost, mBufferHandle int32, startingPosition int32, dataLength int32, dataOffset int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(mBufferGetByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return ManagedBufferSetByteSliceWithTypedArgs(host, mBufferHandle, startingPosition, dataLength, data)
}

func ManagedBufferSetByteSliceWithTypedArgs(host arwen.VMHost, mBufferHandle int32, startingPosition int32, dataLength int32, data []byte) int32 {
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(mBufferGetByteSliceName)

	managedType.ConsumeGasForBytes(data)

	bufferBytes, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	if startingPosition < 0 || dataLength < 0 || int(startingPosition+dataLength) > len(bufferBytes) {
		// does not fail execution if slice exceeds bounds
		return 1
	}

	start := int(startingPosition)
	length := int(dataLength)
	destination := bufferBytes[start : start+length]

	copy(destination, data)

	managedType.SetBytes(mBufferHandle, bufferBytes)

	return 0
}

//export v1_4_mBufferAppend
func v1_4_mBufferAppend(context unsafe.Pointer, accumulatorHandle int32, dataHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(mBufferAppendName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppend
	metering.UseAndTraceGas(gasToUse)

	dataBufferBytes, err := managedType.GetBytes(dataHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(dataBufferBytes)

	isSuccess := managedType.AppendBytes(accumulatorHandle, dataBufferBytes)
	if !isSuccess {
		_ = arwen.WithFault(arwen.ErrNoManagedBufferUnderThisHandle, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	storage := arwen.GetStorageContext(context)
	if !storage.IsUseDifferentGasCostFlagSet() {
		gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(dataBufferBytes)))
		metering.UseAndTraceGas(gasToUse)
	}

	return 0
}

//export v1_4_mBufferAppendBytes
func v1_4_mBufferAppendBytes(context unsafe.Pointer, accumulatorHandle int32, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(mBufferAppendBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppendBytes
	metering.UseAndTraceGas(gasToUse)

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
	metering.UseAndTraceGas(gasToUse)

	return 0
}

//export v1_4_mBufferToBigIntUnsigned
func v1_4_mBufferToBigIntUnsigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigIntUnsigned
	metering.UseGasAndAddTracedGas(mBufferToBigIntUnsignedName, gasToUse)

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
	metering.UseGasAndAddTracedGas(mBufferToBigIntSignedName, gasToUse)

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
	metering.UseGasAndAddTracedGas(mBufferFromBigIntUnsignedName, gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(mBufferHandle, value.Bytes())

	return 0
}

//export v1_4_mBufferFromBigIntSigned
func v1_4_mBufferFromBigIntSigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigIntSigned
	metering.UseGasAndAddTracedGas(mBufferFromBigIntSignedName, gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(mBufferHandle, twos.ToBytes(value))
	return 0
}

//export v1_4_mBufferToBigFloat
func v1_4_mBufferToBigFloat(context unsafe.Pointer, mBufferHandle, bigFloatHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(mBufferToBigFloatName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigFloat
	metering.UseAndTraceGas(gasToUse)

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBytes(managedBuffer)
	if managedType.EncodedBigFloatIsNotValid(managedBuffer) {
		_ = arwen.WithFault(arwen.ErrBigFloatWrongPrecision, context, runtime.BigFloatAPIErrorShouldFailExecution())
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
	metering.StartGasTracing(mBufferFromBigFloatName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigFloat
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigFloat(bigFloatHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return 1
	}

	encodedFloat, err := value.GobEncode()
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(encodedFloat)

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
	metering.UseGasAndAddTracedGas(mBufferStorageStoreName, gasToUse)

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

	key, err := managedType.GetBytes(keyHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	storageBytes, usedCache, err := storage.GetStorage(key)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 0
	}
	storage.UseGasForStorageLoad(mBufferStorageLoadName, metering.GasSchedule().ManagedBufferAPICost.MBufferStorageLoad, usedCache)

	managedType.SetBytes(destinationHandle, storageBytes)

	return 0
}

//export v1_4_mBufferStorageLoadFromAddress
func v1_4_mBufferStorageLoadFromAddress(context unsafe.Pointer, addressHandle, keyHandle, destinationHandle int32) {
	host := arwen.GetVMHost(context)
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)

	key, err := managedType.GetBytes(keyHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	storageBytes, err := StorageLoadFromAddressWithTypedArgs(host, address, key)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, storageBytes)
}

//export v1_4_mBufferGetArgument
func v1_4_mBufferGetArgument(context unsafe.Pointer, id int32, destinationHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetArgument
	metering.UseGasAndAddTracedGas(mBufferGetArgumentName, gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		arwen.WithFaultAndHostIfFailAlwaysActive(arwen.ErrArgOutOfRange, arwen.GetVMHost(context), runtime.ElrondAPIErrorShouldFailExecution())
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
	metering.StartGasTracing(mBufferFinishName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFinish
	metering.UseAndTraceGas(gasToUse)

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(sourceBytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		_ = arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	output.Finish(sourceBytes)
	return 0
}

//export v1_4_mBufferSetRandom
func v1_4_mBufferSetRandom(context unsafe.Pointer, destinationHandle int32, length int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	if length < 1 {
		_ = arwen.WithFault(arwen.ErrLengthOfBufferNotCorrect, context, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return -1
	}

	baseGasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetRandom
	lengthDependentGasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(baseGasToUse, lengthDependentGasToUse)
	metering.UseGasAndAddTracedGas(mBufferSetRandomName, gasToUse)

	randomizer := managedType.GetRandReader()
	buffer := make([]byte, length)
	_, err := randomizer.Read(buffer)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	managedType.SetBytes(destinationHandle, buffer)
	return 0
}
