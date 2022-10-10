package elrondapi

import (
	"math/big"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/executor"
	"github.com/ElrondNetwork/wasm-vm/math"
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

// BigIntImports populates imports with the BigInt API methods.
func BigIntImports(imports executor.ImportFunctionReceiver) error {
	imports.Namespace("env")

	return nil
}

// BigIntGetUnsignedArgument EIInterface implementation
func (context *EICallbacks) BigIntGetUnsignedArgument(id int32, destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedArgument
	metering.UseGasAndAddTracedGas(bigIntGetUnsignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	value.SetBytes(args[id])
}

// BigIntGetSignedArgument EIInterface implementation
func (context *EICallbacks) BigIntGetSignedArgument(id int32, destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedArgument
	metering.UseGasAndAddTracedGas(bigIntGetSignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	twos.SetBytes(value, args[id])
}

// BigIntStorageStoreUnsigned EIInterface implementation
func (context *EICallbacks) BigIntStorageStoreUnsigned(keyOffset int32, keyLength int32, sourceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageStoreUnsigned
	metering.UseGasAndAddTracedGas(bigIntStorageStoreUnsignedName, gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
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

// BigIntStorageLoadUnsigned EIInterface implementation
func (context *EICallbacks) BigIntStorageLoadUnsigned(keyOffset int32, keyLength int32, destinationHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes, usedCache := storage.GetStorage(key)
	storage.UseGasForStorageLoad(bigIntStorageLoadUnsignedName, metering.GasSchedule().BigIntAPICost.BigIntStorageLoadUnsigned, usedCache)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)

	return int32(len(bytes))
}

// BigIntGetCallValue EIInterface implementation
func (context *EICallbacks) BigIntGetCallValue(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	metering.UseGasAndAddTracedGas(bigIntGetCallValueName, gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.Set(runtime.GetVMInput().CallValue)
}

// BigIntGetESDTCallValue EIInterface implementation
func (context *EICallbacks) BigIntGetESDTCallValue(destination int32) {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return
	}
	context.BigIntGetESDTCallValueByIndex(destination, 0)
}

// BigIntGetESDTCallValueByIndex EIInterface implementation
func (context *EICallbacks) BigIntGetESDTCallValueByIndex(destinationHandle int32, index int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	metering.UseGasAndAddTracedGas(bigIntGetESDTCallValueByIndexName, gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if esdtTransfer != nil {
		value.Set(esdtTransfer.ESDTValue)
	} else {
		value.Set(big.NewInt(0))
	}
}

// BigIntGetExternalBalance EIInterface implementation
func (context *EICallbacks) BigIntGetExternalBalance(addressOffset int32, result int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	metering.UseGasAndAddTracedGas(bigIntGetExternalBalanceName, gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)
	value := managedType.GetBigIntOrCreate(result)

	value.SetBytes(balance)
}

// BigIntGetESDTExternalBalance EIInterface implementation
func (context *EICallbacks) BigIntGetESDTExternalBalance(addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64, resultHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetESDTExternalBalanceName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	metering.UseAndTraceGas(gasToUse)

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

// BigIntNew EIInterface implementation
func (context *EICallbacks) BigIntNew(smallValue int64) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNew
	metering.UseGasAndAddTracedGas(bigIntNewName, gasToUse)

	return managedType.NewBigIntFromInt64(smallValue)
}

// BigIntUnsignedByteLength EIInterface implementation
func (context *EICallbacks) BigIntUnsignedByteLength(referenceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntUnsignedByteLength
	metering.UseGasAndAddTracedGas(bigIntUnsignedByteLengthName, gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes := value.Bytes()
	return int32(len(bytes))
}

// BigIntSignedByteLength EIInterface implementation
func (context *EICallbacks) BigIntSignedByteLength(referenceHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSignedByteLength
	metering.UseGasAndAddTracedGas(bigIntSignedByteLengthName, gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes := twos.ToBytes(value)
	return int32(len(bytes))
}

// BigIntGetUnsignedBytes EIInterface implementation
func (context *EICallbacks) BigIntGetUnsignedBytes(referenceHandle int32, byteOffset int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedBytes
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	bytes := value.Bytes()

	err = runtime.MemStore(byteOffset, bytes)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseAndTraceGas(gasToUse)

	return int32(len(bytes))
}

// BigIntGetSignedBytes EIInterface implementation
func (context *EICallbacks) BigIntGetSignedBytes(referenceHandle int32, byteOffset int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedBytes
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	bytes := twos.ToBytes(value)

	err = runtime.MemStore(byteOffset, bytes)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseAndTraceGas(gasToUse)

	return int32(len(bytes))
}

// BigIntSetUnsignedBytes EIInterface implementation
func (context *EICallbacks) BigIntSetUnsignedBytes(destinationHandle int32, byteOffset int32, byteLength int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetUnsignedBytes
	metering.UseAndTraceGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseGas(gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)
}

// BigIntSetSignedBytes EIInterface implementation
func (context *EICallbacks) BigIntSetSignedBytes(destinationHandle int32, byteOffset int32, byteLength int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetSignedBytes
	metering.UseAndTraceGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseGas(gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	twos.SetBytes(value, bytes)
}

// BigIntIsInt64 EIInterface implementation
func (context *EICallbacks) BigIntIsInt64(destinationHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntIsInt64
	metering.UseGasAndAddTracedGas(bigIntIsInt64Name, gasToUse)

	value, err := managedType.GetBigInt(destinationHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	if value.IsInt64() {
		return 1
	}
	return 0
}

// BigIntGetInt64 EIInterface implementation
func (context *EICallbacks) BigIntGetInt64(destinationHandle int32) int64 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGasAndAddTracedGas(bigIntGetInt64Name, gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	return value.Int64()
}

// BigIntSetInt64 EIInterface implementation
func (context *EICallbacks) BigIntSetInt64(destinationHandle int32, value int64) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetInt64
	metering.UseGasAndAddTracedGas(bigIntSetInt64Name, gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	dest.SetInt64(value)
}

// BigIntAdd EIInterface implementation
func (context *EICallbacks) BigIntAdd(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntAddName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAdd
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	dest.Add(a, b)
}

// BigIntSub EIInterface implementation
func (context *EICallbacks) BigIntSub(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntSubName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	dest.Sub(a, b)
}

// BigIntMul EIInterface implementation
func (context *EICallbacks) BigIntMul(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntMulName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntMul
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)

	dest.Mul(a, b)
}

// BigIntTDiv EIInterface implementation
func (context *EICallbacks) BigIntTDiv(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntTDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTDiv
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = context.WithFault(arwen.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Quo(a, b) // Quo implements truncated division (like Go)
}

// BigIntTMod EIInterface implementation
func (context *EICallbacks) BigIntTMod(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntTModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTMod
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = context.WithFault(arwen.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rem(a, b) // Rem implements truncated modulus (like Go)
}

// BigIntEDiv EIInterface implementation
func (context *EICallbacks) BigIntEDiv(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntEDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEDiv
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = context.WithFault(arwen.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Div(a, b) // Div implements Euclidean division (unlike Go)
}

// BigIntEMod EIInterface implementation
func (context *EICallbacks) BigIntEMod(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntEModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEMod
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = context.WithFault(arwen.ErrDivZero, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Mod(a, b) // Mod implements Euclidean division (unlike Go)
}

// BigIntSqrt EIInterface implementation
func (context *EICallbacks) BigIntSqrt(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntSqrtName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSqrt
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	if a.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBadLowerBounds, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Sqrt(a)
}

// BigIntPow EIInterface implementation
func (context *EICallbacks) BigIntPow(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntPowName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntPow
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	//this calculates the length of the result in bytes
	lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(b, big.NewInt(int64(a.BitLen()))), big.NewInt(8))

	managedType.ConsumeGasForThisBigIntNumberOfBytes(lengthOfResult)
	managedType.ConsumeGasForBigIntCopy(a, b)

	if b.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBadLowerBounds, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}

	dest.Exp(a, b, nil)
}

// BigIntLog2 EIInterface implementation
func (context *EICallbacks) BigIntLog2(op1Handle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntLog2Name)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntLog
	metering.UseAndTraceGas(gasToUse)

	a, err := managedType.GetBigInt(op1Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	managedType.ConsumeGasForBigIntCopy(a)
	if a.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBadLowerBounds, runtime.BigIntAPIErrorShouldFailExecution())
		return -1
	}

	return int32(a.BitLen() - 1)
}

// BigIntAbs EIInterface implementation
func (context *EICallbacks) BigIntAbs(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntAbsName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAbs
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	dest.Abs(a)
}

// BigIntNeg EIInterface implementation
func (context *EICallbacks) BigIntNeg(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntNegName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNeg
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	dest.Neg(a)
}

// BigIntSign EIInterface implementation
func (context *EICallbacks) BigIntSign(opHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntSignName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSign
	metering.UseAndTraceGas(gasToUse)

	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}
	managedType.ConsumeGasForBigIntCopy(a)
	return int32(a.Sign())
}

// BigIntCmp EIInterface implementation
func (context *EICallbacks) BigIntCmp(op1Handle, op2Handle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntCmpName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntCmp
	metering.UseAndTraceGas(gasToUse)

	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	return int32(a.Cmp(b))
}

// BigIntNot EIInterface implementation
func (context *EICallbacks) BigIntNot(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntNotName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNot
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	if a.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Not(a)
}

// BigIntAnd EIInterface implementation
func (context *EICallbacks) BigIntAnd(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntAndName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAnd
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	if a.Sign() < 0 || b.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.And(a, b)
}

// BigIntOr EIInterface implementation
func (context *EICallbacks) BigIntOr(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntOrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntOr
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	if a.Sign() < 0 || b.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Or(a, b)
}

// BigIntXor EIInterface implementation
func (context *EICallbacks) BigIntXor(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntXorName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntXor
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	if a.Sign() < 0 || b.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBitwiseNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Xor(a, b)
}

// BigIntShr EIInterface implementation
func (context *EICallbacks) BigIntShr(destinationHandle, opHandle, bits int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntShrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShr
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a)
	if a.Sign() < 0 || bits < 0 {
		_ = context.WithFault(arwen.ErrShiftNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rsh(a, uint(bits))
	managedType.ConsumeGasForBigIntCopy(dest)
}

// BigIntShl EIInterface implementation
func (context *EICallbacks) BigIntShl(destinationHandle, opHandle, bits int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntShlName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShl
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a)
	if a.Sign() < 0 || bits < 0 {
		_ = context.WithFault(arwen.ErrShiftNegative, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Lsh(a, uint(bits))
	managedType.ConsumeGasForBigIntCopy(dest)

}

// BigIntFinishUnsigned EIInterface implementation
func (context *EICallbacks) BigIntFinishUnsigned(referenceHandle int32) {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntFinishUnsignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishUnsigned
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	bigIntBytes := value.Bytes()
	output.Finish(bigIntBytes)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(value.Bytes())))
	metering.UseAndTraceGas(gasToUse)
}

// BigIntFinishSigned EIInterface implementation
func (context *EICallbacks) BigIntFinishSigned(referenceHandle int32) {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigIntFinishSignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	bigInt2cBytes := twos.ToBytes(value)
	output.Finish(bigInt2cBytes)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(bigInt2cBytes)))
	metering.UseAndTraceGas(gasToUse)
}

// BigIntToString EIInterface implementation
func (context *EICallbacks) BigIntToString(bigIntHandle int32, destinationHandle int32) {
	host := context.GetVMHost()
	BigIntToStringWithHost(host, bigIntHandle, destinationHandle)
}

func BigIntToStringWithHost(host arwen.VMHost, bigIntHandle int32, destinationHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	metering.UseGasAndAddTracedGas(bigIntToStringName, gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if WithFaultAndHost(host, err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	resultStr := value.String()
	managedType.SetBytes(destinationHandle, []byte(resultStr))
	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(resultStr)))
	metering.UseAndTraceGas(gasToUse)
}
