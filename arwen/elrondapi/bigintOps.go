package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t bigIntNew(void* context, long long smallValue);
//
// extern int32_t bigIntUnsignedByteLength(void* context, int32_t reference);
// extern int32_t bigIntSignedByteLength(void* context, int32_t reference);
// extern int32_t bigIntGetUnsignedBytes(void* context, int32_t reference, int32_t byteOffset);
// extern int32_t bigIntGetSignedBytes(void* context, int32_t reference, int32_t byteOffset);
// extern void bigIntSetUnsignedBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
// extern void bigIntSetSignedBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
//
// extern int32_t bigIntIsInt64(void* context, int32_t reference);
// extern long long bigIntGetInt64(void* context, int32_t reference);
// extern void bigIntSetInt64(void* context, int32_t destination, long long value);
//
// extern void bigIntAdd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntSub(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntMul(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntTDiv(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntTMod(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntEDiv(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntEMod(void* context, int32_t destination, int32_t op1, int32_t op2);
//
// extern void bigIntAbs(void* context, int32_t destination, int32_t op);
// extern void bigIntNeg(void* context, int32_t destination, int32_t op);
// extern int32_t bigIntSign(void* context, int32_t op);
// extern int32_t bigIntCmp(void* context, int32_t op1, int32_t op2);
//
// extern void bigIntNot(void* context, int32_t destination, int32_t op);
// extern void bigIntAnd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntOr(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntXor(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntShr(void* context, int32_t destination, int32_t op, int32_t bits);
// extern void bigIntShl(void* context, int32_t destination, int32_t op, int32_t bits);
//
// extern void bigIntFinishUnsigned(void* context, int32_t reference);
// extern void bigIntFinishSigned(void* context, int32_t reference);
// extern int32_t bigIntStorageStoreUnsigned(void *context, int32_t keyOffset, int32_t keyLength, int32_t source);
// extern int32_t bigIntStorageLoadUnsigned(void *context, int32_t keyOffset, int32_t keyLength, int32_t destination);
// extern void bigIntGetUnsignedArgument(void *context, int32_t id, int32_t destination);
// extern void bigIntGetSignedArgument(void *context, int32_t id, int32_t destination);
// extern void bigIntGetCallValue(void *context, int32_t destination);
// extern void bigIntGetExternalBalance(void *context, int32_t addressOffset, int32_t result);
import "C"

import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

// BigIntImports creates a new wasmer.Imports populated with the BigInt API methods
func BigIntImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("bigIntNew", bigIntNew, C.bigIntNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntUnsignedByteLength", bigIntUnsignedByteLength, C.bigIntUnsignedByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSignedByteLength", bigIntSignedByteLength, C.bigIntSignedByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetUnsignedBytes", bigIntGetUnsignedBytes, C.bigIntGetUnsignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetSignedBytes", bigIntGetSignedBytes, C.bigIntGetSignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetUnsignedBytes", bigIntSetUnsignedBytes, C.bigIntSetUnsignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetSignedBytes", bigIntSetSignedBytes, C.bigIntSetSignedBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntIsInt64", bigIntIsInt64, C.bigIntIsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetInt64", bigIntGetInt64, C.bigIntGetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetInt64", bigIntSetInt64, C.bigIntSetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAdd", bigIntAdd, C.bigIntAdd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSub", bigIntSub, C.bigIntSub)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntMul", bigIntMul, C.bigIntMul)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntTDiv", bigIntTDiv, C.bigIntTDiv)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntTMod", bigIntTMod, C.bigIntTMod)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntEDiv", bigIntEDiv, C.bigIntEDiv)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntEMod", bigIntEMod, C.bigIntEMod)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAbs", bigIntAbs, C.bigIntAbs)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntNeg", bigIntNeg, C.bigIntNeg)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSign", bigIntSign, C.bigIntSign)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntCmp", bigIntCmp, C.bigIntCmp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntNot", bigIntNot, C.bigIntNot)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAnd", bigIntAnd, C.bigIntAnd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntOr", bigIntOr, C.bigIntOr)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntXor", bigIntXor, C.bigIntXor)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntShr", bigIntShr, C.bigIntShr)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntShl", bigIntShl, C.bigIntShl)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntFinishUnsigned", bigIntFinishUnsigned, C.bigIntFinishUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntFinishSigned", bigIntFinishSigned, C.bigIntFinishSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntStorageStoreUnsigned", bigIntStorageStoreUnsigned, C.bigIntStorageStoreUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntStorageLoadUnsigned", bigIntStorageLoadUnsigned, C.bigIntStorageLoadUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetUnsignedArgument", bigIntGetUnsignedArgument, C.bigIntGetUnsignedArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetSignedArgument", bigIntGetSignedArgument, C.bigIntGetSignedArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetCallValue", bigIntGetCallValue, C.bigIntGetCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetExternalBalance", bigIntGetExternalBalance, C.bigIntGetExternalBalance)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export bigIntGetUnsignedArgument
func bigIntGetUnsignedArgument(context unsafe.Pointer, id int32, destination int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id {
		return
	}

	value := bigInt.GetOne(destination)

	value.SetBytes(args[id])
}

//export bigIntGetSignedArgument
func bigIntGetSignedArgument(context unsafe.Pointer, id int32, destination int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id {
		return
	}

	value := bigInt.GetOne(destination)

	twos.SetBytes(value, args[id])
}

//export bigIntStorageStoreUnsigned
func bigIntStorageStoreUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, source int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageStoreUnsigned
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 0
	}

	value := bigInt.GetOne(source)
	bytes := value.Bytes()

	storageStatus, err := storage.SetStorage(key, bytes)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export bigIntStorageLoadUnsigned
func bigIntStorageLoadUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, destination int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageLoadUnsigned
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 0
	}

	bytes := storage.GetStorage(key)

	value := bigInt.GetOne(destination)
	value.SetBytes(bytes)

	return int32(len(bytes))
}

//export bigIntGetCallValue
func bigIntGetCallValue(context unsafe.Pointer, destination int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(destination)
	value.Set(runtime.GetVMInput().CallValue)
}

//export bigIntGetExternalBalance
func bigIntGetExternalBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)
	value := bigInt.GetOne(result)

	value.SetBytes(balance)
}

//export bigIntNew
func bigIntNew(context unsafe.Pointer, smallValue int64) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNew
	metering.UseGas(gasToUse)

	return bigInt.Put(smallValue)
}

//export bigIntUnsignedByteLength
func bigIntUnsignedByteLength(context unsafe.Pointer, reference int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntUnsignedByteLength
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(reference)

	bytes := value.Bytes()
	return int32(len(bytes))
}

//export bigIntSignedByteLength
func bigIntSignedByteLength(context unsafe.Pointer, reference int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSignedByteLength
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(reference)

	bytes := twos.ToBytes(value)
	return int32(len(bytes))
}

//export bigIntGetUnsignedBytes
func bigIntGetUnsignedBytes(context unsafe.Pointer, reference int32, byteOffset int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetUnsignedBytes
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(reference)
	bytes := value.Bytes()

	err := runtime.MemStore(byteOffset, bytes)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 0
	}

	gasToUse = metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	metering.UseGas(gasToUse)

	return int32(len(bytes))
}

//export bigIntGetSignedBytes
func bigIntGetSignedBytes(context unsafe.Pointer, reference int32, byteOffset int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetSignedBytes
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(reference)
	bytes := twos.ToBytes(value)

	err := runtime.MemStore(byteOffset, bytes)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return 0
	}

	gasToUse = metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	metering.UseGas(gasToUse)

	return int32(len(bytes))
}

//export bigIntSetUnsignedBytes
func bigIntSetUnsignedBytes(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetUnsignedBytes
	metering.UseGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value := bigInt.GetOne(destination)
	value.SetBytes(bytes)

	gasToUse = metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	metering.UseGas(gasToUse)
}

//export bigIntSetSignedBytes
func bigIntSetSignedBytes(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetSignedBytes
	metering.UseGas(gasToUse)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}

	value := bigInt.GetOne(destination)
	twos.SetBytes(value, bytes)

	gasToUse = metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	metering.UseGas(gasToUse)
}

//export bigIntIsInt64
func bigIntIsInt64(context unsafe.Pointer, destination int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntIsInt64
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(destination)
	if value.IsInt64() {
		return 1
	}
	return 0
}

//export bigIntGetInt64
func bigIntGetInt64(context unsafe.Pointer, destination int32) int64 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(destination)
	return value.Int64()
}

//export bigIntSetInt64
func bigIntSetInt64(context unsafe.Pointer, destination int32, value int64) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest := bigInt.GetOne(destination)
	dest.SetInt64(value)
}

//export bigIntAdd
func bigIntAdd(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	dest.Add(a, b)
}

//export bigIntSub
func bigIntSub(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	dest.Sub(a, b)
}

//export bigIntMul
func bigIntMul(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntMul
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	dest.Mul(a, b)
}

//export bigIntTDiv
func bigIntTDiv(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntTDiv
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	if b.Sign() == 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Quo(a, b) // Quo implements truncated division (like Go)
}

//export bigIntTMod
func bigIntTMod(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	if b.Sign() == 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rem(a, b) // Rem implements truncated modulus (like Go)
}

//export bigIntEDiv
func bigIntEDiv(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	if b.Sign() == 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Div(a, b) // Div implements Euclidean division (unlike Go)
}

//export bigIntEMod
func bigIntEMod(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	if b.Sign() == 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrDivZero, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Mod(a, b) // Mod implements Euclidean division (unlike Go)
}

//export bigIntAbs
func bigIntAbs(context unsafe.Pointer, destination, op int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a := bigInt.GetTwo(destination, op)
	dest.Abs(a)
}

//export bigIntNeg
func bigIntNeg(context unsafe.Pointer, destination, op int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a := bigInt.GetTwo(destination, op)
	dest.Neg(a)
}

//export bigIntSign
func bigIntSign(context unsafe.Pointer, op int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSign
	metering.UseGas(gasToUse)

	a := bigInt.GetOne(op)
	return int32(a.Sign())
}

//export bigIntCmp
func bigIntCmp(context unsafe.Pointer, op1, op2 int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntCmp
	metering.UseGas(gasToUse)

	a, b := bigInt.GetTwo(op1, op2)
	return int32(a.Cmp(b))
}

//export bigIntNot
func bigIntNot(context unsafe.Pointer, destination, op int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a := bigInt.GetTwo(destination, op)
	if a.Sign() < 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Not(a)
}

//export bigIntAnd
func bigIntAnd(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	if a.Sign() < 0 || b.Sign() < 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.And(a, b)
}

//export bigIntOr
func bigIntOr(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	if a.Sign() < 0 || b.Sign() < 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Or(a, b)
}

//export bigIntXor
func bigIntXor(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	if a.Sign() < 0 || b.Sign() < 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrBitwiseNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Xor(a, b)
}

//export bigIntShr
func bigIntShr(context unsafe.Pointer, destination, op, bits int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a := bigInt.GetTwo(destination, op)
	if a.Sign() < 0 || bits < 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrShiftNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Rsh(a, uint(bits))
}

//export bigIntShl
func bigIntShl(context unsafe.Pointer, destination, op, bits int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)

	dest, a := bigInt.GetTwo(destination, op)
	if a.Sign() < 0 || bits < 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrShiftNegative, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	dest.Lsh(a, uint(bits))
}

//export bigIntFinishUnsigned
func bigIntFinishUnsigned(context unsafe.Pointer, reference int32) {
	bigInt := arwen.GetBigIntContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishUnsigned
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(reference)
	bigIntBytes := value.Bytes()
	output.Finish(bigIntBytes)

	gasToUse = metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(value.Bytes()))
	metering.UseGas(gasToUse)
}

//export bigIntFinishSigned
func bigIntFinishSigned(context unsafe.Pointer, reference int32) {
	bigInt := arwen.GetBigIntContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinishSigned
	metering.UseGas(gasToUse)

	value := bigInt.GetOne(reference)
	bigInt2cBytes := twos.ToBytes(value)
	output.Finish(bigInt2cBytes)

	gasToUse = metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(bigInt2cBytes))
	metering.UseGas(gasToUse)
}
