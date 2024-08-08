package vmhooks

import (
	"bytes"
	"math/big"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
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

// MBufferNew VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferNew() int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNew
	err := metering.UseGasBoundedAndAddTracedGas(mBufferNewName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	return managedType.NewManagedBuffer()
}

// MBufferNewFromBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferNewFromBytes(dataOffset executor.MemPtr, dataLength executor.MemLength) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNewFromBytes
	err := metering.UseGasBoundedAndAddTracedGas(mBufferNewFromBytesName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	return managedType.NewManagedBufferFromBytes(data)
}

// MBufferGetLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetLength(mBufferHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetLength
	err := metering.UseGasBoundedAndAddTracedGas(mBufferGetLengthName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	length := managedType.GetLength(mBufferHandle)
	if length == -1 {
		_ = context.WithFault(vmhost.ErrNoManagedBufferUnderThisHandle, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return -1
	}

	return length
}

// MBufferGetBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetBytes(mBufferHandle int32, resultOffset executor.MemPtr) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferGetBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	mBufferBytes, err := managedType.GetBytes(mBufferHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(mBufferBytes)

	err = context.MemStore(resultOffset, mBufferBytes)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// MBufferGetByteSlice VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetByteSlice(
	sourceHandle int32,
	startingPosition int32,
	sliceLength int32,
	resultOffset executor.MemPtr) int32 {

	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferGetByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetByteSlice
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sourceBytes)

	if startingPosition < 0 || sliceLength < 0 || int(startingPosition+sliceLength) > len(sourceBytes) {
		// does not fail execution if slice exceeds bounds
		return 1
	}

	slice := sourceBytes[startingPosition : startingPosition+sliceLength]
	err = context.MemStore(resultOffset, slice)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// MBufferCopyByteSlice VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferCopyByteSlice(sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedBufferCopyByteSliceWithHost(host, sourceHandle, startingPosition, sliceLength, destinationHandle)
}

// ManagedBufferCopyByteSliceWithHost VMHooks implementation.
func ManagedBufferCopyByteSliceWithHost(host vmhost.VMHost, sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) int32 {
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(mBufferCopyByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferCopyByteSlice
	err := metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
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
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// MBufferEq VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferEq(mBufferHandle1 int32, mBufferHandle2 int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferEqName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferCopyByteSlice
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bytes1, err := managedType.GetBytes(mBufferHandle1)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}
	managedType.ConsumeGasForBytes(bytes1)

	bytes2, err := managedType.GetBytes(mBufferHandle2)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}
	managedType.ConsumeGasForBytes(bytes2)

	if bytes.Equal(bytes1, bytes2) {
		return 1
	}

	return 0
}

// MBufferSetBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferSetBytes(mBufferHandle int32, dataOffset executor.MemPtr, dataLength executor.MemLength) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferSetBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(data)
	managedType.SetBytes(mBufferHandle, data)

	return 0
}

// MBufferSetByteSlice VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferSetByteSlice(
	mBufferHandle int32,
	startingPosition int32,
	dataLength executor.MemLength,
	dataOffset executor.MemPtr) int32 {

	host := context.GetVMHost()
	return context.ManagedBufferSetByteSliceWithHost(host, mBufferHandle, startingPosition, dataLength, dataOffset)
}

// ManagedBufferSetByteSliceWithHost VMHooks implementation.
func (context *VMHooksImpl) ManagedBufferSetByteSliceWithHost(
	host vmhost.VMHost,
	mBufferHandle int32,
	startingPosition int32,
	dataLength executor.MemLength,
	dataOffset executor.MemPtr) int32 {

	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(mBufferGetByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return ManagedBufferSetByteSliceWithTypedArgs(host, mBufferHandle, startingPosition, dataLength, data)
}

// ManagedBufferSetByteSliceWithTypedArgs VMHooks implementation.
func ManagedBufferSetByteSliceWithTypedArgs(host vmhost.VMHost, mBufferHandle int32, startingPosition int32, dataLength int32, data []byte) int32 {
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(mBufferGetByteSliceName)

	managedType.ConsumeGasForBytes(data)

	bufferBytes, err := managedType.GetBytes(mBufferHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
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

// MBufferAppend VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferAppend(accumulatorHandle int32, dataHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferAppendName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppend
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	dataBufferBytes, err := managedType.GetBytes(dataHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(dataBufferBytes)

	isSuccess := managedType.AppendBytes(accumulatorHandle, dataBufferBytes)
	if !isSuccess {
		_ = context.WithFault(vmhost.ErrNoManagedBufferUnderThisHandle, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	return 0
}

// MBufferAppendBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferAppendBytes(accumulatorHandle int32, dataOffset executor.MemPtr, dataLength executor.MemLength) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferAppendBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppendBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	isSuccess := managedType.AppendBytes(accumulatorHandle, data)
	if !isSuccess {
		_ = context.WithFault(vmhost.ErrNoManagedBufferUnderThisHandle, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// MBufferToBigIntUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToBigIntUnsigned(mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigIntUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferToBigIntUnsignedName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	bigInt.SetBytes(managedBuffer)

	return 0
}

// MBufferToBigIntSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToBigIntSigned(mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigIntSigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferToBigIntSignedName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	twos.SetBytes(bigInt, managedBuffer)

	return 0
}

// MBufferFromBigIntUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromBigIntUnsigned(mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigIntUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferFromBigIntUnsignedName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	value, err := managedType.GetBigInt(bigIntHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(mBufferHandle, value.Bytes())

	return 0
}

// MBufferFromBigIntSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromBigIntSigned(mBufferHandle int32, bigIntHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigIntSigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferFromBigIntSignedName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	value, err := managedType.GetBigInt(bigIntHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(mBufferHandle, twos.ToBytes(value))
	return 0
}

// MBufferToBigFloat VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToBigFloat(mBufferHandle, bigFloatHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferToBigFloatName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigFloat
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBytes(managedBuffer)
	if managedType.EncodedBigFloatIsNotValid(managedBuffer) {
		_ = context.WithFault(vmhost.ErrBigFloatWrongPrecision, runtime.BigFloatAPIErrorShouldFailExecution())
		return 1
	}

	value, err := managedType.GetBigFloatOrCreate(bigFloatHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	bigFloat := new(big.Float)
	err = bigFloat.GobDecode(managedBuffer)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	if bigFloat.IsInf() {
		_ = context.WithFault(vmhost.ErrInfinityFloatOperation, runtime.BigFloatAPIErrorShouldFailExecution())
		return 1
	}

	value.Set(bigFloat)
	return 0
}

// MBufferFromBigFloat VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromBigFloat(mBufferHandle, bigFloatHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferFromBigFloatName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigFloat
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	value, err := managedType.GetBigFloat(bigFloatHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return 1
	}

	encodedFloat, err := value.GobEncode()
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(encodedFloat)

	managedType.SetBytes(mBufferHandle, encodedFloat)

	return 0
}

// MBufferStorageStore VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferStorageStore(keyHandle int32, sourceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferStorageStore
	err := metering.UseGasBoundedAndAddTracedGas(mBufferStorageStoreName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	key, err := managedType.GetBytes(keyHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	_, err = storage.SetStorage(key, sourceBytes)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// MBufferStorageLoad VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferStorageLoad(keyHandle int32, destinationHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := managedType.GetBytes(keyHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	storageBytes, trieDepth, usedCache, err := storage.GetStorage(key)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 0
	}

	err = storage.UseGasForStorageLoad(
		mBufferStorageLoadName,
		int64(trieDepth),
		metering.GasSchedule().ManagedBufferAPICost.MBufferStorageLoad,
		usedCache)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	managedType.SetBytes(destinationHandle, storageBytes)

	return 0
}

// MBufferStorageLoadFromAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferStorageLoadFromAddress(addressHandle, keyHandle, destinationHandle int32) {
	host := context.GetVMHost()
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()

	key, err := managedType.GetBytes(keyHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	storageBytes, err := StorageLoadFromAddressWithTypedArgs(host, address, key)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, storageBytes)
}

// MBufferGetArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetArgument(id int32, destinationHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetArgument
	err := metering.UseGasBoundedAndAddTracedGas(mBufferGetArgumentName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}
	managedType.SetBytes(destinationHandle, args[id])
	return 0
}

// MBufferFinish VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFinish(sourceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(mBufferFinishName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFinish
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(sourceBytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		_ = context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return 1
	}

	output.Finish(sourceBytes)
	return 0
}

// MBufferSetRandom VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferSetRandom(destinationHandle int32, length int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	if length < 1 {
		_ = context.WithFault(vmhost.ErrLengthOfBufferNotCorrect, runtime.ManagedBufferAPIErrorShouldFailExecution())
		return -1
	}

	baseGasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetRandom
	lengthDependentGasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(baseGasToUse, lengthDependentGasToUse)
	err := metering.UseGasBoundedAndAddTracedGas(mBufferSetRandomName, gasToUse)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	randomizer := managedType.GetRandReader()
	buffer := make([]byte, length)
	_, err = randomizer.Read(buffer)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return -1
	}

	managedType.SetBytes(destinationHandle, buffer)
	return 0
}
