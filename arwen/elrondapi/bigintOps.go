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
// extern void bigIntGetArgument(void *context, int32_t id, int32_t destination);
// extern void bigIntGetCallValue(void *context, int32_t destination);
// extern void bigIntGetExternalBalance(void *context, int32_t addressOffset, int32_t result);
import "C"

import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
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

	imports, err = imports.Append("bigIntGetArgument", bigIntGetArgument, C.bigIntGetArgument)
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

//export bigIntGetArgument
func bigIntGetArgument(context unsafe.Pointer, id int32, destination int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntGetArgument
	hostContext.UseGas(gasToUse)

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return
	}

	value := hostContext.GetOne(destination)
	value.Set(args[id])
}

//export bigIntStorageStore
func bigIntStorageStore(context unsafe.Pointer, keyOffset int32, source int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	value := hostContext.GetOne(source)
	bytes := value.Bytes()

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntStorageStore
	gasToUse += hostContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(bytes))
	hostContext.UseGas(gasToUse)

	return hostContext.SetStorage(hostContext.GetSCAddress(), key, bytes)
}

//export bigIntStorageLoad
func bigIntStorageLoad(context unsafe.Pointer, keyOffset int32, destination int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	bytes := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	value := hostContext.GetOne(destination)
	value.SetBytes(bytes)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntStorageLoad
	gasToUse += hostContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	hostContext.UseGas(gasToUse)

	return int32(len(bytes))
}

//export bigIntGetCallValue
func bigIntGetCallValue(context unsafe.Pointer, destination int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	value := hostContext.GetOne(destination)
	value.Set(hostContext.GetVMInput().CallValue)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntGetCallValue
	hostContext.UseGas(gasToUse)
}

//export bigIntGetExternalBalance
func bigIntGetExternalBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	balance := hostContext.GetBalance(address)
	value := hostContext.GetOne(result)

	value.SetBytes(balance)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntGetExternalBalance
	hostContext.UseGas(gasToUse)
}

//export bigIntNew
func bigIntNew(context unsafe.Pointer, smallValue int64) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntNew
	hostContext.UseGas(gasToUse)

	return hostContext.Put(smallValue)
}

//export bigIntByteLength
func bigIntByteLength(context unsafe.Pointer, reference int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	value := hostContext.GetOne(reference)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntByteLength
	hostContext.UseGas(gasToUse)

	return int32(len(value.Bytes()))
}

//export bigIntGetBytes
func bigIntGetBytes(context unsafe.Pointer, reference int32, byteOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	bytes := hostContext.GetOne(reference).Bytes()

	_ = arwen.StoreBytes(instCtx.Memory(), byteOffset, bytes)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntGetBytes
	gasToUse += hostContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	hostContext.UseGas(gasToUse)

	return int32(len(bytes))
}

//export bigIntSetBytes
func bigIntSetBytes(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	bytes := arwen.LoadBytes(instCtx.Memory(), byteOffset, byteLength)

	value := hostContext.GetOne(destination)
	value.SetBytes(bytes)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntSetBytes
	gasToUse += hostContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(bytes))
	hostContext.UseGas(gasToUse)
}

//export bigIntIsInt64
func bigIntIsInt64(context unsafe.Pointer, destination int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntIsInt64
	hostContext.UseGas(gasToUse)

	value := hostContext.GetOne(destination)
	if value.IsInt64() {
		return 1
	}
	return 0
}

//export bigIntGetInt64
func bigIntGetInt64(context unsafe.Pointer, destination int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntGetInt64
	hostContext.UseGas(gasToUse)

	value := hostContext.GetOne(destination)
	return value.Int64()
}

//export bigIntSetInt64
func bigIntSetInt64(context unsafe.Pointer, destination int32, value int64) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	dest := hostContext.GetOne(destination)
	dest.SetInt64(value)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntSetInt64
	hostContext.UseGas(gasToUse)
}

//export bigIntAdd
func bigIntAdd(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	dest, a, b := hostContext.GetThree(destination, op1, op2)
	dest.Add(a, b)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntAdd
	hostContext.UseGas(gasToUse)
}

//export bigIntSub
func bigIntSub(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	dest, a, b := hostContext.GetThree(destination, op1, op2)
	dest.Sub(a, b)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntSub
	hostContext.UseGas(gasToUse)
}

//export bigIntMul
func bigIntMul(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	dest, a, b := hostContext.GetThree(destination, op1, op2)
	dest.Mul(a, b)

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntMul
	hostContext.UseGas(gasToUse)
}

//export bigIntCmp
func bigIntCmp(context unsafe.Pointer, op1, op2 int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntCmp
	hostContext.UseGas(gasToUse)

	a, b := hostContext.GetTwo(op1, op2)
	return int32(a.Cmp(b))
}

//export bigIntFinish
func bigIntFinish(context unsafe.Pointer, reference int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetBigIntContext(instCtx.Data())

	value := hostContext.GetOne(reference)
	hostContext.Finish(value.Bytes())

	gasToUse := hostContext.GasSchedule().BigIntAPICost.BigIntFinish
	gasToUse += hostContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(value.Bytes()))
	hostContext.UseGas(gasToUse)
}
