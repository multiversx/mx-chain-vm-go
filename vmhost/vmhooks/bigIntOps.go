package vmhooks

import (
	"math/big"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

const (
	bigIntNewName                     = "bigIntNew"
	bigIntUnsignedByteLengthName      = "bigIntUnsignedByteLength"
	bigIntSignedByteLengthName        = "bigIntSignedByteLength"
	bigIntGetUnsignedBytesName        = "bigIntGetUnsignedBytes"
	bigIntGetSignedBytesName          = "bigIntGetSignedBytes"
	bigIntSetUnsignedBytesName        = "bigIntSetUnsignedBytes"
	bigIntSetSignedBytesName          = "bigIntSetSignedBytes"
	bigIntIsInt64Name                 = "bigIntIsInt64"
	bigIntGetInt64Name                = "bigIntGetInt64"
	bigIntSetInt64Name                = "bigIntSetInt64"
	bigIntAddName                     = "bigIntAdd"
	bigIntSubName                     = "bigIntSub"
	bigIntMulName                     = "bigIntMul"
	bigIntTDivName                    = "bigIntTDiv"
	bigIntTModName                    = "bigIntTMod"
	bigIntEDivName                    = "bigIntEDiv"
	bigIntEModName                    = "bigIntEMod"
	bigIntPowName                     = "bigIntPow"
	bigIntLog2Name                    = "bigIntLog2"
	bigIntSqrtName                    = "bigIntSqrt"
	bigIntAbsName                     = "bigIntAbs"
	bigIntNegName                     = "bigIntNeg"
	bigIntSignName                    = "bigIntSign"
	bigIntCmpName                     = "bigIntCmp"
	bigIntNotName                     = "bigIntNot"
	bigIntAndName                     = "bigIntAnd"
	bigIntOrName                      = "bigIntOr"
	bigIntXorName                     = "bigIntXor"
	bigIntShrName                     = "bigIntShr"
	bigIntShlName                     = "bigIntShl"
	bigIntFinishUnsignedName          = "bigIntFinishUnsigned"
	bigIntFinishSignedName            = "bigIntFinishSigned"
	bigIntStorageStoreUnsignedName    = "bigIntStorageStoreUnsigned"
	bigIntStorageLoadUnsignedName     = "bigIntStorageLoadUnsigned"
	bigIntGetUnsignedArgumentName     = "bigIntGetUnsignedArgument"
	bigIntGetSignedArgumentName       = "bigIntGetSignedArgument"
	bigIntGetCallValueName            = "bigIntGetCallValue"
	bigIntGetESDTCallValueName        = "bigIntGetESDTCallValue"
	bigIntGetESDTCallValueByIndexName = "bigIntGetESDTCallValueByIndex"
	bigIntGetESDTExternalBalanceName  = "bigIntGetESDTExternalBalance"
	bigIntGetExternalBalanceName      = "bigIntGetExternalBalance"
	bigIntToStringName                = "bigIntToString"
)

// BigIntGetUnsignedArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetUnsignedArgument(id int32, destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedArgument
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetUnsignedArgumentName, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	value.SetBytes(args[id])
}

// BigIntGetSignedArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetSignedArgument(id int32, destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedArgument
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetSignedArgumentName, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	twos.SetBytes(value, args[id])
}

// BigIntStorageStoreUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntStorageStoreUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength, sourceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageStoreUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(bigIntStorageStoreUnsignedName, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value := managedType.GetBigIntOrCreate(sourceHandle)
	bytes := value.Bytes()

	storageStatus, err := storage.SetStorage(key, bytes)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

// BigIntStorageLoadUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntStorageLoadUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength, destinationHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes, trieDepth, usedCache, err := storage.GetStorage(key)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	err = storage.UseGasForStorageLoad(bigIntStorageLoadUnsignedName,
		int64(trieDepth),
		metering.GasSchedule().BigIntAPICost.BigIntStorageLoadUnsigned,
		usedCache)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)

	return int32(len(bytes))
}

// BigIntGetCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetCallValue(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetCallValueName, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.Set(runtime.GetVMInput().CallValue)
}

// BigIntGetESDTCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetESDTCallValue(destination int32) {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return
	}
	context.BigIntGetESDTCallValueByIndex(destination, 0)
}

// BigIntGetESDTCallValueByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetESDTCallValueByIndex(destinationHandle int32, index int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetESDTCallValueByIndexName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if esdtTransfer != nil {
		value.Set(esdtTransfer.ESDTValue)
	} else {
		value.Set(big.NewInt(0))
	}
}

// BigIntGetExternalBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetExternalBalance(addressOffset executor.MemPtr, result int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetExternalBalanceName, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)
	value := managedType.GetBigIntOrCreate(result)

	value.SetBytes(balance)
}

// BigIntGetESDTExternalBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetESDTExternalBalance(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
	resultHandle int32) {

	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetESDTExternalBalanceName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	if esdtData == nil {
		return
	}

	value := managedType.GetBigIntOrCreate(resultHandle)
	value.Set(esdtData.Value)
}

// BigIntNew VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntNew(smallValue int64) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNew
	err := metering.UseGasBoundedAndAddTracedGas(bigIntNewName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	return managedType.NewBigIntFromInt64(smallValue)
}

// BigIntUnsignedByteLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntUnsignedByteLength(referenceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntUnsignedByteLength
	err := metering.UseGasBoundedAndAddTracedGas(bigIntUnsignedByteLengthName, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes := value.Bytes()
	return int32(len(bytes))
}

// BigIntSignedByteLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSignedByteLength(referenceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSignedByteLength
	err := metering.UseGasBoundedAndAddTracedGas(bigIntSignedByteLengthName, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes := twos.ToBytes(value)
	return int32(len(bytes))
}

// BigIntGetUnsignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetUnsignedBytes(referenceHandle int32, byteOffset executor.MemPtr) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	bytes := value.Bytes()

	err = context.MemStore(byteOffset, bytes)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(bytes))
}

// BigIntGetSignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetSignedBytes(referenceHandle int32, byteOffset executor.MemPtr) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	bytes := twos.ToBytes(value)

	err = context.MemStore(byteOffset, bytes)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(bytes))
}

// BigIntSetUnsignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSetUnsignedBytes(destinationHandle int32, byteOffset executor.MemPtr, byteLength executor.MemLength) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetUnsignedBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	bytes, err := context.MemLoad(byteOffset, byteLength)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)
}

// BigIntSetSignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSetSignedBytes(destinationHandle int32, byteOffset executor.MemPtr, byteLength executor.MemLength) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetSignedBytes
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	bytes, err := context.MemLoad(byteOffset, byteLength)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	twos.SetBytes(value, bytes)
}

// BigIntIsInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntIsInt64(destinationHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntIsInt64
	err := metering.UseGasBoundedAndAddTracedGas(bigIntIsInt64Name, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value, err := managedType.GetBigInt(destinationHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	if value.IsInt64() {
		return 1
	}
	return 0
}

// BigIntGetInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetInt64(destinationHandle int32) int64 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetInt64Name, gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	if !value.IsInt64() {
		if context.WithFault(vmhost.ErrBigIntCannotBeRepresentedAsInt64, runtime.BigIntAPIErrorShouldFailExecution()) {
			return -1
		}
	}
	return value.Int64()
}

// BigIntSetInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSetInt64(destinationHandle int32, value int64) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetInt64
	err := metering.UseGasBoundedAndAddTracedGas(bigIntSetInt64Name, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	dest.SetInt64(value)
}

// BigIntAdd VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntAdd(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntAddName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAdd
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest.Add(a, b)
}

// BigIntSub VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSub(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntSubName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest.Sub(a, b)
}

// BigIntMul VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntMul(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntMulName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntMul
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest.Mul(a, b)
}

// BigIntTDiv VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntTDiv(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntTDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTDiv
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if b.Sign() == 0 {
		_ = context.WithFault(vmhost.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Quo(a, b) // Quo implements truncated division (like Go)
}

// BigIntTMod VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntTMod(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntTModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTMod
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if b.Sign() == 0 {
		_ = context.WithFault(vmhost.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rem(a, b) // Rem implements truncated modulus (like Go)
}

// BigIntEDiv VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntEDiv(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntEDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEDiv
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if b.Sign() == 0 {
		_ = context.WithFault(vmhost.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Div(a, b) // Div implements Euclidean division (unlike Go)
}

// BigIntEMod VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntEMod(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntEModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEMod
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if b.Sign() == 0 {
		_ = context.WithFault(vmhost.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Mod(a, b) // Mod implements Euclidean division (unlike Go)
}

// BigIntSqrt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSqrt(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntSqrtName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSqrt
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if a.Sign() < 0 {
		_ = context.WithFault(vmhost.ErrBadLowerBounds, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Sqrt(a)
}

// BigIntPow VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntPow(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntPowName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntPow
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	//this calculates the length of the result in bytes
	lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(b, big.NewInt(int64(a.BitLen()))), big.NewInt(8))

	err = managedType.ConsumeGasForThisBigIntNumberOfBytes(lengthOfResult)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if b.Sign() < 0 {
		_ = context.WithFault(vmhost.ErrBadLowerBounds, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}

	dest.Exp(a, b, nil)
}

// BigIntLog2 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntLog2(op1Handle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntLog2Name)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntLog
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	a, err := managedType.GetBigInt(op1Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	if a.Sign() < 0 {
		_ = context.WithFault(vmhost.ErrBadLowerBounds, runtime.BigIntAPIErrorShouldFailExecution())
		return -1
	}

	return int32(a.BitLen() - 1)
}

// BigIntAbs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntAbs(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntAbsName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAbs
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest.Abs(a)
}

// BigIntNeg VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntNeg(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntNegName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNeg
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest.Neg(a)
}

// BigIntSign VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSign(opHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntSignName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSign
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}

	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}

	return int32(a.Sign())
}

// BigIntCmp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntCmp(op1Handle, op2Handle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntCmpName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntCmp
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}

	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}

	return int32(a.Cmp(b))
}

// BigIntNot VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntNot(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntNotName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNot
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if a.Sign() < 0 {
		_ = context.WithFault(vmhost.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Not(a)
}

// BigIntAnd VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntAnd(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntAndName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAnd
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if a.Sign() < 0 || b.Sign() < 0 {
		_ = context.WithFault(vmhost.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.And(a, b)
}

// BigIntOr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntOr(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntOrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntOr
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if a.Sign() < 0 || b.Sign() < 0 {
		_ = context.WithFault(vmhost.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Or(a, b)
}

// BigIntXor VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntXor(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntXorName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntXor
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if a.Sign() < 0 || b.Sign() < 0 {
		_ = context.WithFault(vmhost.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Xor(a, b)
}

// BigIntShr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntShr(destinationHandle, opHandle, bits int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntShrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShr
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if a.Sign() < 0 || bits < 0 {
		_ = context.WithFault(vmhost.ErrShiftNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rsh(a, uint(bits))

	err = managedType.ConsumeGasForBigIntCopy(dest)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
}

// BigIntShl VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntShl(destinationHandle, opHandle, bits int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntShlName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShl
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	if a.Sign() < 0 || bits < 0 {
		_ = context.WithFault(vmhost.ErrShiftNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Lsh(a, uint(bits))

	err = managedType.ConsumeGasForBigIntCopy(dest)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
}

// BigIntFinishUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntFinishUnsigned(referenceHandle int32) {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntFinishUnsignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishUnsigned
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	bigIntBytes := value.Bytes()

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(value.Bytes())))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	output.Finish(bigIntBytes)
}

// BigIntFinishSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntFinishSigned(referenceHandle int32) {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntFinishSignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	err := metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	bigInt2cBytes := twos.ToBytes(value)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(bigInt2cBytes)))
	err = metering.UseGasBounded(gasToUse)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	output.Finish(bigInt2cBytes)
}

// BigIntToString VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntToString(bigIntHandle int32, destinationHandle int32) {
	host := context.GetVMHost()
	BigIntToStringWithHost(host, bigIntHandle, destinationHandle)
}

func BigIntToStringWithHost(host vmhost.VMHost, bigIntHandle int32, destinationHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	err := metering.UseGasBoundedAndAddTracedGas(bigIntToStringName, gasToUse)
	if WithFaultAndHost(host, err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value, err := managedType.GetBigInt(bigIntHandle)
	if WithFaultAndHost(host, err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	resultStr := value.String()

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(resultStr)))
	err = metering.UseGasBounded(gasToUse)
	if WithFaultAndHost(host, err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, []byte(resultStr))
}
