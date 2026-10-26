package vmhooks

import (
	"bytes"
	"math/big"
	"strings"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

const (
	mBufferNewName                  = "mBufferNew"
	mBufferNewFromBytesName         = "mBufferNewFromBytes"
	mBufferGetLengthName            = "mBufferGetLength"
	mBufferGetBytesName             = "mBufferGetBytes"
	mBufferGetByteSliceName         = "mBufferGetByteSlice"
	mBufferCopyByteSliceName        = "mBufferCopyByteSlice"
	mBufferEqName                   = "mBufferEq"
	mBufferSetBytesName             = "mBufferSetBytes"
	mBufferAppendName               = "mBufferAppend"
	mBufferAppendBytesName          = "mBufferAppendBytes"
	mBufferToBigIntUnsignedName     = "mBufferToBigIntUnsigned"
	mBufferToBigIntSignedName       = "mBufferToBigIntSigned"
	mBufferFromBigIntUnsignedName   = "mBufferFromBigIntUnsigned"
	mBufferFromBigIntSignedName     = "mBufferFromBigIntSigned"
	mBufferToSmallIntUnsignedName   = "mBufferToSmallIntUnsigned"
	mBufferToSmallIntSignedName     = "mBufferToSmallIntSigned"
	mBufferFromSmallIntUnsignedName = "mBufferFromSmallIntUnsigned"
	mBufferFromSmallIntSignedName   = "mBufferFromSmallIntSigned"
	mBufferStorageStoreName         = "mBufferStorageStore"
	mBufferStorageLoadName          = "mBufferStorageLoad"
	mBufferGetArgumentName          = "mBufferGetArgument"
	mBufferFinishName               = "mBufferFinish"
	mBufferSetRandomName            = "mBufferSetRandom"
	mBufferToBigFloatName           = "mBufferToBigFloat"
	mBufferFromBigFloatName         = "mBufferFromBigFloat"
)

const maxBufferCopyBytes = 10 * 1024 * 1024

// MBufferNew VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferNew() (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNew
	err := metering.UseGasBoundedAndAddTracedGas(mBufferNewName, gasToUse)
	if err != nil {
		return -1, err
	}

	return managedType.NewManagedBuffer(), nil
}

// MBufferNewFromBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferNewFromBytes(dataOffset executor.MemPtr, dataLength executor.MemLength) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferNewFromBytes
	err := metering.UseGasBoundedAndAddTracedGas(mBufferNewFromBytesName, gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := context.loadBoundedBuffer(dataOffset, dataLength)
	if err != nil {
		return -1, err
	}

	if context.host.Runtime().FixCorrectSmallIntAndGasLimits() {
		err = managedType.ConsumeGasForBytes(data)
		if err != nil {
			return -1, err
		}
	}

	return managedType.NewManagedBufferFromBytes(data), nil
}

// MBufferGetLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetLength(mBufferHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetLength
	err := metering.UseGasBoundedAndAddTracedGas(mBufferGetLengthName, gasToUse)
	if err != nil {
		return -1, err
	}

	length := managedType.GetLength(mBufferHandle)
	if length == -1 {
		return -1, vmhost.ErrNoManagedBufferUnderThisHandle
	}

	return length, nil
}

// MBufferGetBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetBytes(mBufferHandle int32, resultOffset executor.MemPtr) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferGetBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	mBufferBytes, err := managedType.GetBytes(mBufferHandle)
	if err != nil {
		return 1, err
	}
	err = managedType.ConsumeGasForBytes(mBufferBytes)
	if err != nil {
		return 1, err
	}

	err = context.MemStore(resultOffset, mBufferBytes)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// MBufferGetByteSlice VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetByteSlice(
	sourceHandle int32,
	startingPosition int32,
	sliceLength int32,
	resultOffset executor.MemPtr) (int32, error) {

	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferGetByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetByteSlice
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		return 1, err
	}
	err = managedType.ConsumeGasForBytes(sourceBytes)
	if err != nil {
		return 1, err
	}

	err = math.CheckBoundsInt32(startingPosition, sliceLength, len(sourceBytes))
	if err != nil {
		// does not fail execution if slice exceeds bounds
		return 1, nil
	}

	slice := sourceBytes[startingPosition : startingPosition+sliceLength]
	err = context.MemStore(resultOffset, slice)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// MBufferCopyByteSlice VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferCopyByteSlice(sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) (int32, error) {
	host := context.GetVMHost()
	return ManagedBufferCopyByteSliceWithHost(host, sourceHandle, startingPosition, sliceLength, destinationHandle)
}

// ManagedBufferCopyByteSliceWithHost VMHooks implementation.
func ManagedBufferCopyByteSliceWithHost(host vmhost.VMHost, sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) (int32, error) {
	managedType := host.ManagedTypes()
	metering := host.Metering()
	metering.StartGasTracing(mBufferCopyByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferCopyByteSlice
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		return 1, err
	}
	err = managedType.ConsumeGasForBytes(sourceBytes)
	if err != nil {
		return 1, err
	}

	err = math.CheckBoundsInt32(startingPosition, sliceLength, len(sourceBytes))
	if err != nil {
		// does not fail execution if slice exceeds bounds
		return 1, nil
	}

	slice := sourceBytes[startingPosition : startingPosition+sliceLength]
	managedType.SetBytes(destinationHandle, slice)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(slice)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// MBufferEq VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferEq(mBufferHandle1 int32, mBufferHandle2 int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferEqName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferCopyByteSlice
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	bytes1, err := managedType.GetBytes(mBufferHandle1)
	if err != nil {
		return -1, err
	}
	err = managedType.ConsumeGasForBytes(bytes1)
	if err != nil {
		return 1, err
	}

	bytes2, err := managedType.GetBytes(mBufferHandle2)
	if err != nil {
		return -1, err
	}
	err = managedType.ConsumeGasForBytes(bytes2)
	if err != nil {
		return 1, err
	}

	if bytes.Equal(bytes1, bytes2) {
		return 1, nil
	}

	return 0, nil
}

// MBufferSetBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferSetBytes(mBufferHandle int32, dataOffset executor.MemPtr, dataLength executor.MemLength) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferSetBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	data, err := context.loadBoundedBuffer(dataOffset, dataLength)
	if err != nil {
		return 1, err
	}

	err = managedType.ConsumeGasForBytes(data)
	if err != nil {
		return 1, err
	}

	managedType.SetBytes(mBufferHandle, data)

	return 0, nil
}

// MBufferSetByteSlice VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferSetByteSlice(
	mBufferHandle int32,
	startingPosition int32,
	dataLength executor.MemLength,
	dataOffset executor.MemPtr) (int32, error) {

	host := context.GetVMHost()
	return context.ManagedBufferSetByteSliceWithHost(host, mBufferHandle, startingPosition, dataLength, dataOffset)
}

// ManagedBufferSetByteSliceWithHost VMHooks implementation.
func (context *VMHooksImpl) ManagedBufferSetByteSliceWithHost(
	host vmhost.VMHost,
	mBufferHandle int32,
	startingPosition int32,
	dataLength executor.MemLength,
	dataOffset executor.MemPtr) (int32, error) {

	metering := host.Metering()
	metering.StartGasTracing(mBufferGetByteSliceName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if err != nil {
		return 1, err
	}

	return ManagedBufferSetByteSliceWithTypedArgs(host, mBufferHandle, startingPosition, dataLength, data)
}

// ManagedBufferSetByteSliceWithTypedArgs VMHooks implementation.
func ManagedBufferSetByteSliceWithTypedArgs(host vmhost.VMHost, mBufferHandle int32, startingPosition int32, dataLength int32, data []byte) (int32, error) {
	managedType := host.ManagedTypes()
	metering := host.Metering()
	metering.StartGasTracing(mBufferGetByteSliceName)

	err := managedType.ConsumeGasForBytes(data)
	if err != nil {
		return 1, err
	}

	bufferBytes, err := managedType.GetBytes(mBufferHandle)
	if err != nil {
		return 1, err
	}

	err = math.CheckBoundsInt32(startingPosition, dataLength, len(bufferBytes))
	if err != nil {
		// does not fail execution if slice exceeds bounds
		return 1, nil
	}

	start := int(startingPosition)
	length := int(dataLength)

	clonedBytes := bytes.Clone(bufferBytes)
	destination := clonedBytes[start : start+length]

	copy(destination, data)

	managedType.SetBytes(mBufferHandle, clonedBytes)

	return 0, nil
}

// MBufferAppend VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferAppend(accumulatorHandle int32, dataHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferAppendName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppend
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	dataBufferBytes, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return 1, err
	}

	err = managedType.ConsumeGasForBytes(dataBufferBytes)
	if err != nil {
		return 1, err
	}

	isSuccess := managedType.AppendBytes(accumulatorHandle, dataBufferBytes)
	if !isSuccess {
		return 1, vmhost.ErrNoManagedBufferUnderThisHandle
	}

	return 0, nil
}

// MBufferAppendBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferAppendBytes(accumulatorHandle int32, dataOffset executor.MemPtr, dataLength executor.MemLength) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferAppendBytesName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferAppendBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	data, err := context.loadBoundedBuffer(dataOffset, dataLength)
	if err != nil {
		return 1, err
	}

	isSuccess := managedType.AppendBytes(accumulatorHandle, data)
	if !isSuccess {
		return 1, vmhost.ErrNoManagedBufferUnderThisHandle
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// MBufferToBigIntUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToBigIntUnsigned(mBufferHandle int32, bigIntHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	enableEpochsHandler := context.GetEnableEpochsHandler()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigIntUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferToBigIntUnsignedName, gasToUse)
	if err != nil {
		return -1, err
	}

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if err != nil {
		return 1, err
	}

	if enableEpochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag) {
		gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(managedBuffer)))
		err = metering.UseGasBounded(gasToUse)
		if err != nil {
			return 1, err
		}
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	bigInt.SetBytes(managedBuffer)

	return 0, nil
}

// MBufferToBigIntSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToBigIntSigned(mBufferHandle int32, bigIntHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	enableEpochsHandler := context.GetEnableEpochsHandler()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigIntSigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferToBigIntSignedName, gasToUse)
	if err != nil {
		return -1, err
	}

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if err != nil {
		return 1, err
	}

	if enableEpochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag) {
		gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(managedBuffer)))
		err = metering.UseGasBounded(gasToUse)
		if err != nil {
			return 1, err
		}
	}

	bigInt := managedType.GetBigIntOrCreate(bigIntHandle)
	twos.SetBytes(bigInt, managedBuffer)

	return 0, nil
}

// MBufferFromBigIntUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromBigIntUnsigned(mBufferHandle int32, bigIntHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigIntUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferFromBigIntUnsignedName, gasToUse)
	if err != nil {
		return -1, err
	}

	value, err := managedType.GetBigInt(bigIntHandle)
	if err != nil {
		return 1, err
	}

	byteLen := (value.BitLen() + 7) / 8
	err = chargeGasForExtraLength(byteLen, context.GetVMHost())
	if err != nil {
		return -1, err
	}

	managedType.SetBytes(mBufferHandle, value.Bytes())

	return 0, nil
}

// MBufferFromBigIntSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromBigIntSigned(mBufferHandle int32, bigIntHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigIntSigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferFromBigIntSignedName, gasToUse)
	if err != nil {
		return -1, err
	}

	value, err := managedType.GetBigInt(bigIntHandle)
	if err != nil {
		return 1, err
	}

	byteLen := (value.BitLen() + 7) / 8
	err = chargeGasForExtraLength(byteLen, context.GetVMHost())
	if err != nil {
		return -1, err
	}

	managedType.SetBytes(mBufferHandle, twos.ToBytes(value))
	return 0, nil
}

// MBufferToSmallIntUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToSmallIntUnsigned(mBufferHandle int32) (int64, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToSmallIntUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferToSmallIntUnsignedName, gasToUse)
	if err != nil {
		return 0, err
	}

	data, err := managedType.GetBytes(mBufferHandle)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	bigInt := big.NewInt(0).SetBytes(data)
	if !bigInt.IsUint64() {
		return 0, vmhost.ErrBytesExceedUint64
	}
	return int64(bigInt.Uint64()), nil
}

// MBufferToSmallIntSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToSmallIntSigned(mBufferHandle int32) (int64, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToSmallIntSigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferToSmallIntSignedName, gasToUse)
	if err != nil {
		return 0, err
	}

	data, err := managedType.GetBytes(mBufferHandle)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	bigInt := twos.SetBytes(big.NewInt(0), data)
	if !bigInt.IsInt64() {
		return 0, vmhost.ErrBytesExceedInt64
	}
	return bigInt.Int64(), nil
}

// MBufferFromSmallIntUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromSmallIntUnsigned(mBufferHandle int32, value int64) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromSmallIntUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferFromSmallIntUnsignedName, gasToUse)
	if err != nil {
		return err
	}

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	managedType.SetBytes(mBufferHandle, valueBytes)
	return nil
}

// MBufferFromSmallIntSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromSmallIntSigned(mBufferHandle int32, value int64) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromSmallIntSigned
	err := metering.UseGasBoundedAndAddTracedGas(mBufferFromSmallIntSignedName, gasToUse)
	if err != nil {
		return err
	}

	var valueBytes []byte
	if context.GetRuntimeContext().FixCorrectSmallIntAndGasLimits() {
		valueBytes = twos.ToBytes(big.NewInt(value))
	} else {
		valueBytes = big.NewInt(0).SetInt64(value).Bytes()
	}
	managedType.SetBytes(mBufferHandle, valueBytes)
	return nil
}

// MBufferToBigFloat VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferToBigFloat(mBufferHandle, bigFloatHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferToBigFloatName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferToBigFloat
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	managedBuffer, err := managedType.GetBytes(mBufferHandle)
	if err != nil {
		return 1, err
	}

	err = managedType.ConsumeGasForBytes(managedBuffer)
	if err != nil {
		return 1, err
	}

	if managedType.EncodedBigFloatIsNotValid(managedBuffer) {
		return 1, vmhost.ErrBigFloatWrongPrecision
	}

	value, err := managedType.GetBigFloatOrCreate(bigFloatHandle)
	if err != nil {
		return 1, err
	}

	bigFloat := new(big.Float)
	err = bigFloat.GobDecode(managedBuffer)
	if err != nil {
		if !enableEpochsHandler.IsFlagEnabled(vmhost.ValidationOnGobDecodeFlag) &&
			isGobDecodeValidationError(err) {

		} else {
			if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
				err = vmhost.ErrBigFloatDecode
			}

			return 1, err
		}
	}

	if bigFloat.IsInf() {
		return 1, vmhost.ErrInfinityFloatOperation
	}

	value.Set(bigFloat)
	return 0, nil
}

func isGobDecodeValidationError(err error) bool {
	if err == nil {
		return false
	}

	validationErrors := []string{
		"nonzero finite number with empty mantissa",
		"msb not set in last word",
		"zero precision finite number",
	}

	for _, validationError := range validationErrors {
		if strings.Contains(err.Error(), validationError) {
			return true
		}
	}

	return false
}

// MBufferFromBigFloat VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFromBigFloat(mBufferHandle, bigFloatHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferFromBigFloatName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFromBigFloat
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	value, err := managedType.GetBigFloat(bigFloatHandle)
	if err != nil {
		return 1, err
	}

	encodedFloat, err := value.GobEncode()
	if err != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = vmhost.ErrBigFloatEncode
		}
		return 1, err
	}

	err = managedType.ConsumeGasForBytes(encodedFloat)
	if err != nil {
		return 1, err
	}

	managedType.SetBytes(mBufferHandle, encodedFloat)

	return 0, nil
}

// MBufferStorageStore VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferStorageStore(keyHandle int32, sourceHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferStorageStore
	err := metering.UseGasBoundedAndAddTracedGas(mBufferStorageStoreName, gasToUse)
	if err != nil {
		return -1, err
	}

	key, err := managedType.GetBytes(keyHandle)
	if err != nil {
		return 1, err
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		return 1, err
	}

	_, err = storage.SetStorage(key, sourceBytes)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// MBufferStorageLoad VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferStorageLoad(keyHandle int32, destinationHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := managedType.GetBytes(keyHandle)
	if err != nil {
		return 1, err
	}

	storageBytes, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return 0, err
	}

	err = storage.UseGasForStorageLoad(
		mBufferStorageLoadName,
		int64(trieDepth),
		metering.GasSchedule().ManagedBufferAPICost.MBufferStorageLoad,
		usedCache)
	if err != nil {
		return -1, err
	}

	if context.host.Runtime().FixCorrectSmallIntAndGasLimits() {
		err = managedType.ConsumeGasForBytes(storageBytes)
		if err != nil {
			return -1, err
		}
	}

	managedType.SetBytes(destinationHandle, storageBytes)

	return 0, nil
}

// MBufferStorageLoadFromAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferStorageLoadFromAddress(addressHandle, keyHandle, destinationHandle int32) error {
	host := context.GetVMHost()
	managedType := context.GetManagedTypesContext()

	key, err := managedType.GetBytes(keyHandle)
	if err != nil {
		return err
	}

	address, err := loadAddressFromHandle(context.host, managedType, addressHandle)
	if err != nil {
		return err
	}

	storageBytes, err := StorageLoadFromAddressWithTypedArgs(host, address, key)
	if err != nil {
		return err
	}

	if context.GetRuntimeContext().FixCorrectSmallIntAndGasLimits() {
		err = managedType.ConsumeGasForBytes(storageBytes)
		if err != nil {
			return err
		}
	}

	managedType.SetBytes(destinationHandle, storageBytes)
	return nil
}

// MBufferGetArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferGetArgument(id int32, destinationHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetArgument
	err := metering.UseGasBoundedAndAddTracedGas(mBufferGetArgumentName, gasToUse)
	if err != nil {
		return -1, err
	}

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return 1, vmhost.ErrArgOutOfRange
	}

	if context.host.Runtime().FixCorrectSmallIntAndGasLimits() {
		err = managedType.ConsumeGasForBytes(args[id])
		if err != nil {
			return -1, err
		}
	}

	managedType.SetBytes(destinationHandle, args[id])
	return 0, nil
}

// MBufferFinish VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferFinish(sourceHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(mBufferFinishName)

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferFinish
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	sourceBytes, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(sourceBytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	output.Finish(sourceBytes)
	return 0, nil
}

// MBufferSetRandom VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MBufferSetRandom(destinationHandle int32, length int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	if length < 1 {
		return -1, vmhost.ErrLengthOfBufferNotCorrect
	}

	baseGasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetRandom
	lengthDependentGasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(baseGasToUse, lengthDependentGasToUse)
	err := metering.UseGasBoundedAndAddTracedGas(mBufferSetRandomName, gasToUse)
	if err != nil {
		return -1, err
	}

	randomizer := managedType.GetRandReader()
	buffer := make([]byte, length)
	_, err = randomizer.Read(buffer)
	if err != nil {
		return -1, err
	}

	managedType.SetBytes(destinationHandle, buffer)
	return 0, nil
}

func (context *VMHooksImpl) loadBoundedBuffer(dataOffset executor.MemPtr, dataLength executor.MemLength) ([]byte, error) {
	if context.host.Runtime().FixCorrectSmallIntAndGasLimits() {
		if int64(dataLength) > maxBufferCopyBytes || dataLength < 0 {
			return nil, vmhost.ErrArgOutOfRange
		}
	}

	return context.MemLoad(dataOffset, dataLength)
}
