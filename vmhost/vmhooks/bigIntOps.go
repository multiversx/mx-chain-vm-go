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
func (context *VMHooksImpl) BigIntGetUnsignedArgument(id int32, destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedArgument
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetUnsignedArgumentName, gasToUse)
	if err != nil {
		return err
	}

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return nil
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	value.SetBytes(args[id])
	return nil
}

// BigIntGetSignedArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetSignedArgument(id int32, destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	enableEpochsHandler := context.GetEnableEpochsHandler()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedArgument
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetSignedArgumentName, gasToUse)
	if err != nil {
		return err
	}

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return nil
	}

	if enableEpochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag) {
		gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(args[id])))
		err = metering.UseGasBounded(gasToUse)
		if err != nil {
			return err
		}
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	twos.SetBytes(value, args[id])
	return nil
}

// BigIntStorageStoreUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntStorageStoreUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength, sourceHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageStoreUnsigned
	err := metering.UseGasBoundedAndAddTracedGas(bigIntStorageStoreUnsignedName, gasToUse)
	if err != nil {
		return -1, err
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	value := managedType.GetBigIntOrCreate(sourceHandle)
	bytes := value.Bytes()

	storageStatus, err := storage.SetStorage(key, bytes)
	if err != nil {
		return -1, err
	}

	return int32(storageStatus), nil
}

// BigIntStorageLoadUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntStorageLoadUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength, destinationHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	bytes, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return -1, err
	}

	err = storage.UseGasForStorageLoad(bigIntStorageLoadUnsignedName,
		int64(trieDepth),
		metering.GasSchedule().BigIntAPICost.BigIntStorageLoadUnsigned,
		usedCache)
	if err != nil {
		return -1, err
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)

	return int32(len(bytes)), nil
}

// BigIntGetCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetCallValue(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetCallValueName, gasToUse)
	if err != nil {
		return err
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.Set(runtime.GetVMInput().CallValue)
	return nil
}

// BigIntGetESDTCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetESDTCallValue(destination int32) error {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return nil
	}
	return context.BigIntGetESDTCallValueByIndex(destination, 0)
}

// BigIntGetESDTCallValueByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetESDTCallValueByIndex(destinationHandle int32, index int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetESDTCallValueByIndexName, gasToUse)
	if err != nil {
		return err
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if esdtTransfer != nil {
		value.Set(esdtTransfer.ESDTValue)
	} else {
		value.Set(big.NewInt(0))
	}
	return nil
}

// BigIntGetExternalBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetExternalBalance(addressOffset executor.MemPtr, result int32) error {
	managedType := context.GetManagedTypesContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetExternalBalanceName, gasToUse)
	if err != nil {
		return err
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return err
	}

	balance := blockchain.GetBalance(address)
	value := managedType.GetBigIntOrCreate(result)

	value.SetBytes(balance)
	return nil
}

// BigIntGetESDTExternalBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetESDTExternalBalance(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
	resultHandle int32) error {

	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetESDTExternalBalanceName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if err != nil {
		return err
	}
	if esdtData == nil {
		return nil
	}

	value := managedType.GetBigIntOrCreate(resultHandle)
	value.Set(esdtData.Value)
	return nil
}

// BigIntNew VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntNew(smallValue int64) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNew
	err := metering.UseGasBoundedAndAddTracedGas(bigIntNewName, gasToUse)
	if err != nil {
		return -1, err
	}

	return managedType.NewBigIntFromInt64(smallValue), nil
}

// BigIntUnsignedByteLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntUnsignedByteLength(referenceHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntUnsignedByteLength
	err := metering.UseGasBoundedAndAddTracedGas(bigIntUnsignedByteLengthName, gasToUse)
	if err != nil {
		return -1, err
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if err != nil {
		return -1, err
	}

	bytes := value.Bytes()
	return int32(len(bytes)), nil
}

// BigIntSignedByteLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSignedByteLength(referenceHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSignedByteLength
	err := metering.UseGasBoundedAndAddTracedGas(bigIntSignedByteLengthName, gasToUse)
	if err != nil {
		return -1, err
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if err != nil {
		return -1, err
	}

	bytes := twos.ToBytes(value)
	return int32(len(bytes)), nil
}

// BigIntGetUnsignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetUnsignedBytes(referenceHandle int32, byteOffset executor.MemPtr) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if err != nil {
		return -1, err
	}
	bytes := value.Bytes()

	err = context.MemStore(byteOffset, bytes)
	if err != nil {
		return -1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	return int32(len(bytes)), nil
}

// BigIntGetSignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetSignedBytes(referenceHandle int32, byteOffset executor.MemPtr) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntGetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if err != nil {
		return -1, err
	}
	bytes := twos.ToBytes(value)

	err = context.MemStore(byteOffset, bytes)
	if err != nil {
		return -1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	return int32(len(bytes)), nil
}

// BigIntSetUnsignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSetUnsignedBytes(destinationHandle int32, byteOffset executor.MemPtr, byteLength executor.MemLength) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetUnsignedBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	bytes, err := context.MemLoad(byteOffset, byteLength)
	if err != nil {
		return err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)
	return nil
}

// BigIntSetSignedBytes VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSetSignedBytes(destinationHandle int32, byteOffset executor.MemPtr, byteLength executor.MemLength) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetSignedBytes
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	bytes, err := context.MemLoad(byteOffset, byteLength)
	if err != nil {
		return err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	twos.SetBytes(value, bytes)
	return nil
}

// BigIntIsInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntIsInt64(destinationHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntIsInt64
	err := metering.UseGasBoundedAndAddTracedGas(bigIntIsInt64Name, gasToUse)
	if err != nil {
		return -1, err
	}

	value, err := managedType.GetBigInt(destinationHandle)
	if err != nil {
		return -1, err
	}
	if value.IsInt64() {
		return 1, nil
	}
	return 0, nil
}

// BigIntGetInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntGetInt64(destinationHandle int32) (int64, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	err := metering.UseGasBoundedAndAddTracedGas(bigIntGetInt64Name, gasToUse)
	if err != nil {
		return -1, err
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)
	if !value.IsInt64() {
		return -1, vmhost.ErrBigIntCannotBeRepresentedAsInt64
	}
	return value.Int64(), nil
}

// BigIntSetInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSetInt64(destinationHandle int32, value int64) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetInt64
	err := metering.UseGasBoundedAndAddTracedGas(bigIntSetInt64Name, gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	dest.SetInt64(value)
	return nil
}

// BigIntAdd VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntAdd(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntAddName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAdd
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if err != nil {
		return err
	}

	dest.Add(a, b)
	return nil
}

// BigIntSub VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSub(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSubName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if err != nil {
		return err
	}

	dest.Sub(a, b)
	return nil
}

// BigIntMul VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntMul(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntMulName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntMul
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if err != nil {
		return err
	}

	dest.Mul(a, b)
	return nil
}

// BigIntTDiv VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntTDiv(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntTDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTDiv
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if err != nil {
		return err
	}

	if b.Sign() == 0 {
		return vmhost.ErrDivZero
	}
	dest.Quo(a, b) // Quo implements truncated division (like Go)
	return nil
}

// BigIntTMod VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntTMod(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntTModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTMod
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if err != nil {
		return err
	}

	if b.Sign() == 0 {
		return vmhost.ErrDivZero
	}
	dest.Rem(a, b) // Rem implements truncated modulus (like Go)
	return nil
}

// BigIntEDiv VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntEDiv(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntEDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEDiv
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if err != nil {
		return err
	}

	if b.Sign() == 0 {
		return vmhost.ErrDivZero
	}
	dest.Div(a, b) // Div implements Euclidean division (unlike Go)
	return nil
}

// BigIntEMod VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntEMod(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntEModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEMod
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if err != nil {
		return err
	}

	if b.Sign() == 0 {
		return vmhost.ErrDivZero
	}
	dest.Mod(a, b) // Mod implements Euclidean division (unlike Go)
	return nil
}

const maxSqrtLen = 8000

// BigIntSqrt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSqrt(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSqrtName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSqrt
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if err != nil {
		return err
	}

	if a.Sign() < 0 {
		return vmhost.ErrBadLowerBounds
	}

	if context.GetRuntimeContext().AttributeExtraGasUsage() && a.BitLen() > maxSqrtLen {
		return vmhost.ErrBadUpperBounds
	}

	dest.Sqrt(a)
	return nil
}

// BigIntPow VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntPow(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntPowName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntPow
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	//this calculates the length of the result in bytes
	lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(b, big.NewInt(int64(a.BitLen()))), big.NewInt(8))

	err = managedType.ConsumeGasForThisBigIntNumberOfBytes(lengthOfResult)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if err != nil {
		return err
	}

	if b.Sign() < 0 {
		return vmhost.ErrBadLowerBounds
	}

	dest.Exp(a, b, nil)
	return nil
}

// BigIntLog2 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntLog2(op1Handle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntLog2Name)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntLog
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	a, err := managedType.GetBigInt(op1Handle)
	if err != nil {
		return -1, err
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if err != nil {
		return -1, err
	}

	if a.Sign() < 0 {
		return -1, vmhost.ErrBadLowerBounds
	}

	return int32(a.BitLen() - 1), nil
}

// BigIntAbs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntAbs(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntAbsName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAbs
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if err != nil {
		return err
	}

	dest.Abs(a)
	return nil
}

// BigIntNeg VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntNeg(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntNegName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNeg
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if err != nil {
		return err
	}

	dest.Neg(a)
	return nil
}

// BigIntSign VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntSign(opHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntSignName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSign
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -2, err
	}

	a, err := managedType.GetBigInt(opHandle)
	if err != nil {
		return -2, err
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if err != nil {
		return -2, err
	}

	return int32(a.Sign()), nil
}

// BigIntCmp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntCmp(op1Handle, op2Handle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntCmpName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntCmp
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -2, err
	}

	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return -2, err
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if err != nil {
		return -2, err
	}

	return int32(a.Cmp(b)), nil
}

// BigIntNot VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntNot(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntNotName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNot
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(dest, a)
	if err != nil {
		return err
	}

	if a.Sign() < 0 {
		return vmhost.ErrBitwiseNegative
	}
	dest.Not(a)
	return nil
}

// BigIntAnd VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntAnd(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntAndName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAnd
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if err != nil {
		return err
	}

	if a.Sign() < 0 || b.Sign() < 0 {
		return vmhost.ErrBitwiseNegative
	}
	dest.And(a, b)
	return nil
}

// BigIntOr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntOr(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntOrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntOr
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if err != nil {
		return err
	}

	if a.Sign() < 0 || b.Sign() < 0 {
		return vmhost.ErrBitwiseNegative
	}
	dest.Or(a, b)
	return nil
}

// BigIntXor VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntXor(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntXorName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntXor
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(a, b)
	if err != nil {
		return err
	}

	if a.Sign() < 0 || b.Sign() < 0 {
		return vmhost.ErrBitwiseNegative
	}
	dest.Xor(a, b)
	return nil
}

// BigIntShr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntShr(destinationHandle, opHandle, bits int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntShrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShr
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if err != nil {
		return err
	}

	if a.Sign() < 0 || bits < 0 {
		return vmhost.ErrShiftNegative
	}
	dest.Rsh(a, uint(bits))

	err = managedType.ConsumeGasForBigIntCopy(dest)
	if err != nil {
		return err
	}
	return nil
}

// BigIntShl VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntShl(destinationHandle, opHandle, bits int32) error {
	managedType := context.GetManagedTypesContext()
	enableEpochsHandler := context.GetEnableEpochsHandler()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntShlName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShl
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(a)
	if err != nil {
		return err
	}

	if a.Sign() < 0 || bits < 0 {
		return vmhost.ErrShiftNegative
	}

	if !enableEpochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag) {
		dest.Lsh(a, uint(bits))

		err = managedType.ConsumeGasForBigIntCopy(dest)
		if err != nil {
			return err
		}

		return nil
	}

	//this calculates the length of the result in bytes
	resultBits := a.BitLen() + int(bits)
	resultBytes := big.NewInt(0).Div(big.NewInt(0).Add(big.NewInt(int64(resultBits)), big.NewInt(int64(a.BitLen()))), big.NewInt(8))

	err = managedType.ConsumeGasForThisBigIntNumberOfBytes(resultBytes)
	if err != nil {
		return err
	}

	dest.Lsh(a, uint(bits))
	return nil
}

// BigIntFinishUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntFinishUnsigned(referenceHandle int32) error {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntFinishUnsignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishUnsigned
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if err != nil {
		return err
	}
	bigIntBytes := value.Bytes()

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(value.Bytes())))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	output.Finish(bigIntBytes)
	return nil
}

// BigIntFinishSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntFinishSigned(referenceHandle int32) error {
	managedType := context.GetManagedTypesContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigIntFinishSignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	value, err := managedType.GetBigInt(referenceHandle)
	if err != nil {
		return err
	}
	bigInt2cBytes := twos.ToBytes(value)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(bigInt2cBytes)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	output.Finish(bigInt2cBytes)
	return nil
}

// BigIntToString VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigIntToString(bigIntHandle int32, destinationHandle int32) error {
	host := context.GetVMHost()
	return BigIntToStringWithHost(host, bigIntHandle, destinationHandle)
}

func BigIntToStringWithHost(host vmhost.VMHost, bigIntHandle int32, destinationHandle int32) error {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	err := metering.UseGasBoundedAndAddTracedGas(bigIntToStringName, gasToUse)
	if err != nil {
		return err
	}

	value, err := managedType.GetBigInt(bigIntHandle)
	if err != nil {
		return err
	}

	resultStr := value.String()

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(resultStr)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	managedType.SetBytes(destinationHandle, []byte(resultStr))
	return nil
}
