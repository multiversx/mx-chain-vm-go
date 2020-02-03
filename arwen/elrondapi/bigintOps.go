package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t bigIntNew(void* context, long long smallValue);
// extern int32_t bigIntByteLength(void* context, int32_t reference);
// extern int32_t bigIntGetBytes(void* context, int32_t reference, int32_t byteOffset);
// extern void bigIntSetBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
// extern int32_t bigIntIsInt64(void* context, int32_t reference);
// extern long long bigIntGetInt64(void* context, int32_t reference);
// extern void bigIntSetInt64(void* context, int32_t destination, long long value);
// extern void bigIntAdd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntSub(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntMul(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern int32_t bigIntCmp(void* context, int32_t op1, int32_t op2);
// extern void bigIntFinish(void* context, int32_t reference);
// extern int32_t bigIntStorageStore(void *context, int32_t keyOffset, int32_t source);
// extern int32_t bigIntStorageLoad(void *context, int32_t keyOffset, int32_t destination);
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

func BigIntImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("bigIntNew", bigIntNew, C.bigIntNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntByteLength", bigIntByteLength, C.bigIntByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetBytes", bigIntGetBytes, C.bigIntGetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetBytes", bigIntSetBytes, C.bigIntSetBytes)
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

	imports, err = imports.Append("bigIntCmp", bigIntCmp, C.bigIntCmp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntFinish", bigIntFinish, C.bigIntFinish)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntStorageStore", bigIntStorageStore, C.bigIntStorageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntStorageLoad", bigIntStorageLoad, C.bigIntStorageLoad)
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

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetArgument
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

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if int32(len(args)) <= id {
		return
	}

	value := bigInt.GetOne(destination)

	twos.SetBytes(value, args[id])
}

//export bigIntStorageStore
func bigIntStorageStore(context unsafe.Pointer, keyOffset int32, source int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, arwen.HashLen)
	if withFault(err, context) {
		return 0
	}

	value := bigInt.GetOne(source)
	bytes := value.Bytes()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageStore
	metering.UseGas(gasToUse)

	return storage.SetStorage(runtime.GetSCAddress(), key, bytes)
}

//export bigIntStorageLoad
func bigIntStorageLoad(context unsafe.Pointer, keyOffset int32, destination int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, arwen.HashLen)
	if withFault(err, context) {
		return 0
	}

	bytes := storage.GetStorage(runtime.GetSCAddress(), key)

	value := bigInt.GetOne(destination)
	value.SetBytes(bytes)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntStorageLoad
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	metering.UseGas(gasToUse)

	return int32(len(bytes))
}

//export bigIntGetCallValue
func bigIntGetCallValue(context unsafe.Pointer, destination int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	value := bigInt.GetOne(destination)
	value.Set(runtime.GetVMInput().CallValue)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetCallValue
	metering.UseGas(gasToUse)
}

//export bigIntGetExternalBalance
func bigIntGetExternalBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return
	}

	balance := blockchain.GetBalance(address)
	value := bigInt.GetOne(result)

	value.SetBytes(balance)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	metering.UseGas(gasToUse)
}

//export bigIntNew
func bigIntNew(context unsafe.Pointer, smallValue int64) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntNew
	metering.UseGas(gasToUse)

	return bigInt.Put(smallValue)
}

//export bigIntByteLength
func bigIntByteLength(context unsafe.Pointer, reference int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	value := bigInt.GetOne(reference)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntByteLength
	metering.UseGas(gasToUse)

	return int32(len(value.Bytes()))
}

//export bigIntGetBytes
func bigIntGetBytes(context unsafe.Pointer, reference int32, byteOffset int32) int32 {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	bytes := bigInt.GetOne(reference).Bytes()

	err := runtime.MemStore(byteOffset, bytes)
	if withFault(err, context) {
		return 0
	}

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetBytes
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	metering.UseGas(gasToUse)

	return int32(len(bytes))
}

//export bigIntSetBytes
func bigIntSetBytes(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	bigInt := arwen.GetBigIntContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	bytes, err := runtime.MemLoad(byteOffset, byteLength)
	if withFault(err, context) {
		return
	}

	value := bigInt.GetOne(destination)
	value.SetBytes(bytes)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetBytes
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
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

	dest := bigInt.GetOne(destination)
	dest.SetInt64(value)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSetInt64
	metering.UseGas(gasToUse)
}

//export bigIntAdd
func bigIntAdd(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	dest.Add(a, b)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntAdd
	metering.UseGas(gasToUse)
}

//export bigIntSub
func bigIntSub(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	dest.Sub(a, b)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntSub
	metering.UseGas(gasToUse)
}

//export bigIntMul
func bigIntMul(context unsafe.Pointer, destination, op1, op2 int32) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	dest, a, b := bigInt.GetThree(destination, op1, op2)
	dest.Mul(a, b)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntMul
	metering.UseGas(gasToUse)
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

//export bigIntFinish
func bigIntFinish(context unsafe.Pointer, reference int32) {
	bigInt := arwen.GetBigIntContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	value := bigInt.GetOne(reference)
	bigIntBytes := value.Bytes()
	if len(bigIntBytes) == 0 {
		// send one byte of "0", otherwise nothing gets saved when we "return 0"
		bigIntBytes = []byte{0}
	}
	output.Finish(bigIntBytes)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntFinish
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(value.Bytes()))
	metering.UseGas(gasToUse)
}
