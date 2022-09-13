package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t		v1_5_bigIntNew(void* context, long long smallValue);
//
// extern int32_t		v1_5_bigIntUnsignedByteLength(void* context, int32_t reference);
// extern int32_t		v1_5_bigIntSignedByteLength(void* context, int32_t reference);
// extern int32_t		v1_5_bigIntGetUnsignedBytes(void* context, int32_t reference, int32_t byteOffset);
// extern int32_t		v1_5_bigIntGetSignedBytes(void* context, int32_t reference, int32_t byteOffset);
// extern void			v1_5_bigIntSetUnsignedBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
// extern void			v1_5_bigIntSetSignedBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
//
// extern int32_t		v1_5_bigIntIsInt64(void* context, int32_t reference);
// extern long long		v1_5_bigIntGetInt64(void* context, int32_t reference);
// extern void			v1_5_bigIntSetInt64(void* context, int32_t destination, long long value);
//
// extern void			v1_5_bigIntAdd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntSub(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntMul(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntTDiv(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntTMod(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntEDiv(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntEMod(void* context, int32_t destination, int32_t op1, int32_t op2);
//
// extern void			v1_5_bigIntPow(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern int32_t		v1_5_bigIntLog2(void* context, int32_t op);
// extern void			v1_5_bigIntSqrt(void* context, int32_t destination, int32_t op);
//
// extern void			v1_5_bigIntAbs(void* context, int32_t destination, int32_t op);
// extern void			v1_5_bigIntNeg(void* context, int32_t destination, int32_t op);
// extern int32_t		v1_5_bigIntSign(void* context, int32_t op);
// extern int32_t		v1_5_bigIntCmp(void* context, int32_t op1, int32_t op2);
//
// extern void			v1_5_bigIntNot(void* context, int32_t destination, int32_t op);
// extern void			v1_5_bigIntAnd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntOr(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntXor(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void			v1_5_bigIntShr(void* context, int32_t destination, int32_t op, int32_t bits);
// extern void			v1_5_bigIntShl(void* context, int32_t destination, int32_t op, int32_t bits);
//
// extern void			v1_5_bigIntFinishUnsigned(void* context, int32_t reference);
// extern void			v1_5_bigIntFinishSigned(void* context, int32_t reference);
// extern int32_t		v1_5_bigIntStorageStoreUnsigned(void *context, int32_t keyOffset, int32_t keyLength, int32_t source);
// extern int32_t		v1_5_bigIntStorageLoadUnsigned(void *context, int32_t keyOffset, int32_t keyLength, int32_t destination);
// extern void			v1_5_bigIntGetUnsignedArgument(void *context, int32_t id, int32_t destination);
// extern void			v1_5_bigIntGetSignedArgument(void *context, int32_t id, int32_t destination);
// extern void			v1_5_bigIntGetCallValue(void *context, int32_t destination);
// extern void			v1_5_bigIntGetESDTCallValue(void *context, int32_t destination);
// extern void			v1_5_bigIntGetESDTCallValueByIndex(void *context, int32_t destination, int32_t index);
// extern void			v1_5_bigIntGetESDTExternalBalance(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t result);
// extern void			v1_5_bigIntGetExternalBalance(void *context, int32_t addressOffset, int32_t result);
// extern void			v1_5_bigIntToString(void *context, int32_t bigIntHandle, int32_t destinaitonHandle);
import "C"

import (
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/math"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
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

// BigIntImports creates a new wasmer.Imports populated with the BigInt API methods
func BigIntImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("bigIntNew", v1_5_bigIntNew, C.v1_5_bigIntNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntUnsignedByteLength", v1_5_bigIntUnsignedByteLength, C.v1_5_bigIntUnsignedByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSignedByteLength", v1_5_bigIntSignedByteLength, C.v1_5_bigIntSignedByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetUnsignedBytes", v1_5_bigIntGetUnsignedBytes, C.v1_5_bigIntGetUnsignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetSignedBytes", v1_5_bigIntGetSignedBytes, C.v1_5_bigIntGetSignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetUnsignedBytes", v1_5_bigIntSetUnsignedBytes, C.v1_5_bigIntSetUnsignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetSignedBytes", v1_5_bigIntSetSignedBytes, C.v1_5_bigIntSetSignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntIsInt64", v1_5_bigIntIsInt64, C.v1_5_bigIntIsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetInt64", v1_5_bigIntGetInt64, C.v1_5_bigIntGetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetInt64", v1_5_bigIntSetInt64, C.v1_5_bigIntSetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAdd", v1_5_bigIntAdd, C.v1_5_bigIntAdd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSub", v1_5_bigIntSub, C.v1_5_bigIntSub)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntMul", v1_5_bigIntMul, C.v1_5_bigIntMul)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntTDiv", v1_5_bigIntTDiv, C.v1_5_bigIntTDiv)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntTMod", v1_5_bigIntTMod, C.v1_5_bigIntTMod)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntEDiv", v1_5_bigIntEDiv, C.v1_5_bigIntEDiv)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntEMod", v1_5_bigIntEMod, C.v1_5_bigIntEMod)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSqrt", v1_5_bigIntSqrt, C.v1_5_bigIntSqrt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntPow", v1_5_bigIntPow, C.v1_5_bigIntPow)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntLog2", v1_5_bigIntLog2, C.v1_5_bigIntLog2)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAbs", v1_5_bigIntAbs, C.v1_5_bigIntAbs)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntNeg", v1_5_bigIntNeg, C.v1_5_bigIntNeg)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSign", v1_5_bigIntSign, C.v1_5_bigIntSign)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntCmp", v1_5_bigIntCmp, C.v1_5_bigIntCmp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntNot", v1_5_bigIntNot, C.v1_5_bigIntNot)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAnd", v1_5_bigIntAnd, C.v1_5_bigIntAnd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntOr", v1_5_bigIntOr, C.v1_5_bigIntOr)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntXor", v1_5_bigIntXor, C.v1_5_bigIntXor)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntShr", v1_5_bigIntShr, C.v1_5_bigIntShr)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntShl", v1_5_bigIntShl, C.v1_5_bigIntShl)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntFinishUnsigned", v1_5_bigIntFinishUnsigned, C.v1_5_bigIntFinishUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntFinishSigned", v1_5_bigIntFinishSigned, C.v1_5_bigIntFinishSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntStorageStoreUnsigned", v1_5_bigIntStorageStoreUnsigned, C.v1_5_bigIntStorageStoreUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntStorageLoadUnsigned", v1_5_bigIntStorageLoadUnsigned, C.v1_5_bigIntStorageLoadUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetUnsignedArgument", v1_5_bigIntGetUnsignedArgument, C.v1_5_bigIntGetUnsignedArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetSignedArgument", v1_5_bigIntGetSignedArgument, C.v1_5_bigIntGetSignedArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetCallValue", v1_5_bigIntGetCallValue, C.v1_5_bigIntGetCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetESDTCallValue", v1_5_bigIntGetESDTCallValue, C.v1_5_bigIntGetESDTCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetESDTExternalBalance", v1_5_bigIntGetESDTExternalBalance, C.v1_5_bigIntGetESDTExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetExternalBalance", v1_5_bigIntGetExternalBalance, C.v1_5_bigIntGetExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetESDTCallValueByIndex", v1_5_bigIntGetESDTCallValueByIndex, C.v1_5_bigIntGetESDTCallValueByIndex)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntToString", v1_5_bigIntToString, C.v1_5_bigIntToString)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_5_bigIntGetUnsignedArgument
func v1_5_bigIntGetUnsignedArgument(context unsafe.Pointer, id int32, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedArgument
	metering.UseGasAndAddTracedGas(bigIntGetUnsignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	value.SetBytes(args[id])
}

//export v1_5_bigIntGetSignedArgument
func v1_5_bigIntGetSignedArgument(context unsafe.Pointer, id int32, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedArgument
	metering.UseGasAndAddTracedGas(bigIntGetSignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id || id < 0 {
		return
	}

	value := managedType.GetBigIntOrCreate(destinationHandle)

	twos.SetBytes(value, args[id])
}

//export v1_5_bigIntStorageStoreUnsigned
func v1_5_bigIntStorageStoreUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, sourceHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageStoreUnsigned
	metering.UseGasAndAddTracedGas(bigIntStorageStoreUnsignedName, gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	value := managedType.GetBigIntOrCreate(sourceHandle)
	bytes := value.Bytes()

	storageStatus, err := storage.SetStorage(key, bytes)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export v1_5_bigIntStorageLoadUnsigned
func v1_5_bigIntStorageLoadUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, destinationHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes, usedCache := storage.GetStorage(key)
	storage.UseGasForStorageLoad(bigIntStorageLoadUnsignedName, metering.GasSchedule().BigIntAPICost.BigIntStorageLoadUnsigned, usedCache)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)

	return int32(len(bytes))
}

//export v1_5_bigIntGetCallValue
func v1_5_bigIntGetCallValue(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	metering.UseGasAndAddTracedGas(bigIntGetCallValueName, gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.Set(runtime.GetVMInput().CallValue)
}

//export v1_5_bigIntGetESDTCallValue
func v1_5_bigIntGetESDTCallValue(context unsafe.Pointer, destination int32) {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return
	}
	v1_5_bigIntGetESDTCallValueByIndex(context, destination, 0)
}

//export v1_5_bigIntGetESDTCallValueByIndex
func v1_5_bigIntGetESDTCallValueByIndex(context unsafe.Pointer, destinationHandle int32, index int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	metering.UseGasAndAddTracedGas(bigIntGetESDTCallValueByIndexName, gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(arwen.GetVMHost(context), index)
	if esdtTransfer != nil {
		value.Set(esdtTransfer.ESDTValue)
	} else {
		value.Set(big.NewInt(0))
	}
}

//export v1_5_bigIntGetExternalBalance
func v1_5_bigIntGetExternalBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	metering.UseGasAndAddTracedGas(bigIntGetExternalBalanceName, gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)
	value := managedType.GetBigIntOrCreate(result)

	value.SetBytes(balance)
}

//export v1_5_bigIntGetESDTExternalBalance
func v1_5_bigIntGetESDTExternalBalance(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64, resultHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(bigIntGetESDTExternalBalanceName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	metering.UseAndTraceGas(gasToUse)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	if esdtData == nil {
		return
	}

	value := managedType.GetBigIntOrCreate(resultHandle)
	value.Set(esdtData.Value)
}

//export v1_5_bigIntNew
func v1_5_bigIntNew(context unsafe.Pointer, smallValue int64) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNew
	metering.UseGasAndAddTracedGas(bigIntNewName, gasToUse)

	return managedType.NewBigIntFromInt64(smallValue)
}

//export v1_5_bigIntUnsignedByteLength
func v1_5_bigIntUnsignedByteLength(context unsafe.Pointer, referenceHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntUnsignedByteLength
	metering.UseGasAndAddTracedGas(bigIntUnsignedByteLengthName, gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes := value.Bytes()
	return int32(len(bytes))
}

//export v1_5_bigIntSignedByteLength
func v1_5_bigIntSignedByteLength(context unsafe.Pointer, referenceHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSignedByteLength
	metering.UseGasAndAddTracedGas(bigIntSignedByteLengthName, gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	bytes := twos.ToBytes(value)
	return int32(len(bytes))
}

//export v1_5_bigIntGetUnsignedBytes
func v1_5_bigIntGetUnsignedBytes(context unsafe.Pointer, referenceHandle int32, byteOffset int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(bigIntGetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedBytes
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	bytes := value.Bytes()

	err = runtime.MemStore(byteOffset, bytes)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseAndTraceGas(gasToUse)

	return int32(len(bytes))
}

//export v1_5_bigIntGetSignedBytes
func v1_5_bigIntGetSignedBytes(context unsafe.Pointer, referenceHandle int32, byteOffset int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(bigIntGetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedBytes
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	bytes := twos.ToBytes(value)

	err = runtime.MemStore(byteOffset, bytes)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseAndTraceGas(gasToUse)

	return int32(len(bytes))
}

//export v1_5_bigIntSetUnsignedBytes
func v1_5_bigIntSetUnsignedBytes(context unsafe.Pointer, destinationHandle int32, byteOffset int32, byteLength int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(bigIntSetUnsignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetUnsignedBytes
	metering.UseAndTraceGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseGas(gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	value.SetBytes(bytes)
}

//export v1_5_bigIntSetSignedBytes
func v1_5_bigIntSetSignedBytes(context unsafe.Pointer, destinationHandle int32, byteOffset int32, byteLength int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(bigIntSetSignedBytesName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetSignedBytes
	metering.UseAndTraceGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(bytes)))
	metering.UseGas(gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	twos.SetBytes(value, bytes)
}

//export v1_5_bigIntIsInt64
func v1_5_bigIntIsInt64(context unsafe.Pointer, destinationHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntIsInt64
	metering.UseGasAndAddTracedGas(bigIntIsInt64Name, gasToUse)

	value, err := managedType.GetBigInt(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	if value.IsInt64() {
		return 1
	}
	return 0
}

//export v1_5_bigIntGetInt64
func v1_5_bigIntGetInt64(context unsafe.Pointer, destinationHandle int32) int64 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGasAndAddTracedGas(bigIntGetInt64Name, gasToUse)

	value := managedType.GetBigIntOrCreate(destinationHandle)
	return value.Int64()
}

//export v1_5_bigIntSetInt64
func v1_5_bigIntSetInt64(context unsafe.Pointer, destinationHandle int32, value int64) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetInt64
	metering.UseGasAndAddTracedGas(bigIntSetInt64Name, gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	dest.SetInt64(value)
}

//export v1_5_bigIntAdd
func v1_5_bigIntAdd(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntAddName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAdd
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	dest.Add(a, b)
}

//export v1_5_bigIntSub
func v1_5_bigIntSub(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntSubName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	dest.Sub(a, b)
}

//export v1_5_bigIntMul
func v1_5_bigIntMul(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntMulName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntMul
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)

	dest.Mul(a, b)
}

//export v1_5_bigIntTDiv
func v1_5_bigIntTDiv(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntTDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTDiv
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Quo(a, b) // Quo implements truncated division (like Go)
}

//export v1_5_bigIntTMod
func v1_5_bigIntTMod(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntTModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTMod
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rem(a, b) // Rem implements truncated modulus (like Go)
}

//export v1_5_bigIntEDiv
func v1_5_bigIntEDiv(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntEDivName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEDiv
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Div(a, b) // Div implements Euclidean division (unlike Go)
}

//export v1_5_bigIntEMod
func v1_5_bigIntEMod(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntEModName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntEMod
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a, b)
	if b.Sign() == 0 {
		_ = arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Mod(a, b) // Mod implements Euclidean division (unlike Go)
}

//export v1_5_bigIntSqrt
func v1_5_bigIntSqrt(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntSqrtName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSqrt
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	if a.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Sqrt(a)
}

//export v1_5_bigIntPow
func v1_5_bigIntPow(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntPowName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntPow
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	//this calculates the length of the result in bytes
	lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(b, big.NewInt(int64(a.BitLen()))), big.NewInt(8))

	managedType.ConsumeGasForThisBigIntNumberOfBytes(lengthOfResult)
	managedType.ConsumeGasForBigIntCopy(a, b)

	if b.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}

	dest.Exp(a, b, nil)
}

//export v1_5_bigIntLog2
func v1_5_bigIntLog2(context unsafe.Pointer, op1Handle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntLog2Name)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntLog
	metering.UseAndTraceGas(gasToUse)

	a, err := managedType.GetBigInt(op1Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}
	managedType.ConsumeGasForBigIntCopy(a)
	if a.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigIntAPIErrorShouldFailExecution())
		return -1
	}

	return int32(a.BitLen() - 1)
}

//export v1_5_bigIntAbs
func v1_5_bigIntAbs(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntAbsName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAbs
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	dest.Abs(a)
}

//export v1_5_bigIntNeg
func v1_5_bigIntNeg(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntNegName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNeg
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	dest.Neg(a)
}

//export v1_5_bigIntSign
func v1_5_bigIntSign(context unsafe.Pointer, opHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntSignName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSign
	metering.UseAndTraceGas(gasToUse)

	a, err := managedType.GetBigInt(opHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}
	managedType.ConsumeGasForBigIntCopy(a)
	return int32(a.Sign())
}

//export v1_5_bigIntCmp
func v1_5_bigIntCmp(context unsafe.Pointer, op1Handle, op2Handle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntCmpName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntCmp
	metering.UseAndTraceGas(gasToUse)

	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -2
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	return int32(a.Cmp(b))
}

//export v1_5_bigIntNot
func v1_5_bigIntNot(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntNotName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNot
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(dest, a)
	if a.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Not(a)
}

//export v1_5_bigIntAnd
func v1_5_bigIntAnd(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntAndName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAnd
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	if a.Sign() < 0 || b.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.And(a, b)
}

//export v1_5_bigIntOr
func v1_5_bigIntOr(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntOrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntOr
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	if a.Sign() < 0 || b.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Or(a, b)
}

//export v1_5_bigIntXor
func v1_5_bigIntXor(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntXorName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntXor
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, b, err := managedType.GetTwoBigInt(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a, b)
	if a.Sign() < 0 || b.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Xor(a, b)
}

//export v1_5_bigIntShr
func v1_5_bigIntShr(context unsafe.Pointer, destinationHandle, opHandle, bits int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntShrName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShr
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a)
	if a.Sign() < 0 || bits < 0 {
		_ = arwen.WithFault(arwen.ErrShiftNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rsh(a, uint(bits))
	managedType.ConsumeGasForBigIntCopy(dest)
}

//export v1_5_bigIntShl
func v1_5_bigIntShl(context unsafe.Pointer, destinationHandle, opHandle, bits int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntShlName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntShl
	metering.UseAndTraceGas(gasToUse)

	dest := managedType.GetBigIntOrCreate(destinationHandle)
	a, err := managedType.GetBigInt(opHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBigIntCopy(a)
	if a.Sign() < 0 || bits < 0 {
		_ = arwen.WithFault(arwen.ErrShiftNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Lsh(a, uint(bits))
	managedType.ConsumeGasForBigIntCopy(dest)

}

//export v1_5_bigIntFinishUnsigned
func v1_5_bigIntFinishUnsigned(context unsafe.Pointer, referenceHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntFinishUnsignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishUnsigned
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	bigIntBytes := value.Bytes()
	output.Finish(bigIntBytes)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(value.Bytes())))
	metering.UseAndTraceGas(gasToUse)
}

//export v1_5_bigIntFinishSigned
func v1_5_bigIntFinishSigned(context unsafe.Pointer, referenceHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigIntFinishSignedName)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	metering.UseAndTraceGas(gasToUse)

	value, err := managedType.GetBigInt(referenceHandle)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	bigInt2cBytes := twos.ToBytes(value)
	output.Finish(bigInt2cBytes)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(len(bigInt2cBytes)))
	metering.UseAndTraceGas(gasToUse)
}

//export v1_5_bigIntToString
func v1_5_bigIntToString(context unsafe.Pointer, bigIntHandle int32, destinationHandle int32) {
	host := arwen.GetVMHost(context)
	BigIntToStringWithHost(host, bigIntHandle, destinationHandle)
}

func BigIntToStringWithHost(host arwen.VMHost, bigIntHandle int32, destinationHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	metering.UseGasAndAddTracedGas(bigIntToStringName, gasToUse)

	value, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFaultAndHost(host, err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	resultStr := value.String()
	managedType.SetBytes(destinationHandle, []byte(resultStr))
	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(resultStr)))
	metering.UseAndTraceGas(gasToUse)
}
