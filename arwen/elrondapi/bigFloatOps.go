package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_4_bigFloatNew(void* context, int32_t intBase, int32_t subIntBase, int32_t exponent);
//
// extern void		v1_4_bigFloatAdd(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatSub(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatMul(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatDiv(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatRoundDiv(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatMod(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
//
// extern void		v1_4_bigFloatAbs(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_4_bigFloatNeg(void* context, int32_t destinationHandle, int32_t opHandle);
// extern int32_t	v1_4_bigFloatCmp(void* context, int32_t op1Handle, int32_t op2Handle);
// extern int32_t	v1_4_bigFloatSign(void* context, int32_t opHandle);
// extern void 		v1_4_bigFloatCopy(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_4_bigFloatSqrt(void* context, int32_t destinationHandle, int32_t opHandle);
//
// extern int32_t	v1_4_bigFloatIsInt(void* context, int32_t opHandle);
// extern void		v1_4_bigFloatSetInt64(void* context, int32_t destinationHandle, long long value);
// extern void		v1_4_bigFloatSetBigInt(void* context, int32_t destinationHandle, int32_t bigIntHandle);
import "C"
import (
	"math"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
)

func BigFloatImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("bigFloatNew", v1_4_bigFloatNew, C.v1_4_bigFloatNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatAdd", v1_4_bigFloatAdd, C.v1_4_bigFloatAdd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatSub", v1_4_bigFloatSub, C.v1_4_bigFloatSub)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatMul", v1_4_bigFloatMul, C.v1_4_bigFloatMul)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatDiv", v1_4_bigFloatDiv, C.v1_4_bigFloatDiv)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatRoundDiv", v1_4_bigFloatRoundDiv, C.v1_4_bigFloatRoundDiv)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatMod", v1_4_bigFloatMod, C.v1_4_bigFloatMod)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatNeg", v1_4_bigFloatNeg, C.v1_4_bigFloatNeg)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatCopy", v1_4_bigFloatCopy, C.v1_4_bigFloatCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatCmp", v1_4_bigFloatCmp, C.v1_4_bigFloatCmp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatAbs", v1_4_bigFloatAbs, C.v1_4_bigFloatAbs)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatSign", v1_4_bigFloatSign, C.v1_4_bigFloatSign)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatSqrt", v1_4_bigFloatSqrt, C.v1_4_bigFloatSqrt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatSetInt64", v1_4_bigFloatSetInt64, C.v1_4_bigFloatSetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatIsInt", v1_4_bigFloatIsInt, C.v1_4_bigFloatIsInt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatSetBigInt", v1_4_bigFloatSetBigInt, C.v1_4_bigFloatSetBigInt)
	if err != nil {
		return nil, err
	}

	return imports, err
}

//export v1_4_bigFloatNew
func v1_4_bigFloatNew(context unsafe.Pointer, intBase, subIntBase, exponent int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNew
	metering.UseGas(gasToUse)

	if exponent >= 0 {
		_ = arwen.WithFault(arwen.ErrPositiveExponent, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}

	value := float64(intBase) + float64(subIntBase)*math.Pow10(int(exponent))
	return managedType.PutBigFloat(value)
}

//export v1_4_bigFloatAdd
func v1_4_bigFloatAdd(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext((context))

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAdd
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Add(op1, op2)
}

//export v1_4_bigFloatSub
func v1_4_bigFloatSub(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSub
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Sub(op1, op2)
}

//export v1_4_bigFloatMul
func v1_4_bigFloatMul(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatMul
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Mul(op1, op2)
}

//export v1_4_bigFloatDiv
func v1_4_bigFloatDiv(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatDiv
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Quo(op1, op2)
}

//export v1_4_bigFloatNeg
func v1_4_bigFloatNeg(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNeg
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Neg(op)
}

//export v1_4_bigFloatCopy
func v1_4_bigFloatCopy(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCopy
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Copy(op)
}

//export v1_4_bigFloatCmp
func v1_4_bigFloatCmp(context unsafe.Pointer, op1Handle, op2Handle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCmp
	metering.UseGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -2
	}
	return int32(op1.Cmp(op2))
}

//export v1_4_bigFloatAbs
func v1_4_bigFloatAbs(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Abs(op)
}

//export v1_4_bigFloatSign
func v1_4_bigFloatSign(context unsafe.Pointer, opHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -2
	}
	return int32(op.Sign())
}

//export v1_4_bigFloatSqrt
func v1_4_bigFloatSqrt(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSqrt
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if op.Sign() < 0 {
		arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigFloatAPIErrorShouldFailExecution())
	}
	dest.Sqrt(op)
}

//export v1_4_bigFloatSetInt64
func v1_4_bigFloatSetInt64(context unsafe.Pointer, destinationHandle int32, value int64) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSub
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	dest.SetInt64(value)
}

//export v1_4_bigFloatIsInt
func v1_4_bigFloatIsInt(context unsafe.Pointer, opHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatIsInt
	metering.UseGas(gasToUse)
	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	if op.IsInt() {
		return 1
	}
	return 0
}

//export v1_4_bigFloatSetBigInt
func v1_4_bigFloatSetBigInt(context unsafe.Pointer, destinationHandle, bigIntHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetBigInt
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	bigIntValue, err := managedType.GetBigInt(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.SetInt(bigIntValue)
}

//export v1_4_bigFloatRoundDiv
func v1_4_bigFloatRoundDiv(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatRoundDiv
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Quo(op1, op2)
	rDiv := big.NewInt(0)
	dest.Int(rDiv)
	dest.SetInt(rDiv)
}

//export v1_4_bigFloatMod
func v1_4_bigFloatMod(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatMod
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Quo(op1, op2)
	rDiv := big.NewInt(0)
	dest.Int(rDiv)
	dest.Sub(dest, new(big.Float).SetInt(rDiv))
}
