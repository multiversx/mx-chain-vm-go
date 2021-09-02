package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_4_bigFloatNewFromParts(void* context, int32_t integralPart, int32_t fractionalPart, int32_t exponent);
// extern int32_t	v1_4_bigFloatNewFromFrac(void* context, long long numerator, long long denominator);
// extern int32_t	v1_4_bigFloatNewFromSci(void* context, long long significand, long long exponent);
//
// extern void		v1_4_bigFloatAdd(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatSub(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatMul(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_4_bigFloatDiv(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
//
// extern void		v1_4_bigFloatAbs(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_4_bigFloatNeg(void* context, int32_t destinationHandle, int32_t opHandle);
// extern int32_t	v1_4_bigFloatCmp(void* context, int32_t op1Handle, int32_t op2Handle);
// extern int32_t	v1_4_bigFloatSign(void* context, int32_t opHandle);
// extern void 		v1_4_bigFloatClone(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_4_bigFloatSqrt(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_4_bigFloatPow(void* context, int32_t destinationHandle, int32_t opHandle, int32_t exponent);
//
// extern void		v1_4_bigFloatFloor(void* context, int32_t opHandle, int32_t destBigIntHandle);
// extern void		v1_4_bigFloatCeil(void* context, int32_t opHandle, int32_t destBigIntHandle);
// extern void		v1_4_bigFloatTruncate(void* context, int32_t opHandle, int32_t destBigIntHandle);
//
// extern int32_t	v1_4_bigFloatIsInt(void* context, int32_t opHandle);
// extern void		v1_4_bigFloatSetInt64(void* context, int32_t destinationHandle, long long value);
// extern void		v1_4_bigFloatSetBigInt(void* context, int32_t destinationHandle, int32_t bigIntHandle);
//
// extern void		v1_4_bigFloatGetConstPi(void* context, int32_t destinationHandle);
// extern void		v1_4_bigFloatGetConstE(void* context, int32_t destinationHandle);
import "C"
import (
	"math"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
)

// BigFloatImports creates a new wasmer.Imports populated with the BigFloat API methods
func BigFloatImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("bigFloatNewFromParts", v1_4_bigFloatNewFromParts, C.v1_4_bigFloatNewFromParts)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatNewFromFrac", v1_4_bigFloatNewFromFrac, C.v1_4_bigFloatNewFromFrac)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatNewFromSci", v1_4_bigFloatNewFromSci, C.v1_4_bigFloatNewFromSci)
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

	imports, err = imports.Append("bigFloatTruncate", v1_4_bigFloatTruncate, C.v1_4_bigFloatTruncate)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatNeg", v1_4_bigFloatNeg, C.v1_4_bigFloatNeg)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatClone", v1_4_bigFloatClone, C.v1_4_bigFloatClone)
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

	imports, err = imports.Append("bigFloatPow", v1_4_bigFloatPow, C.v1_4_bigFloatPow)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatFloor", v1_4_bigFloatFloor, C.v1_4_bigFloatFloor)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatCeil", v1_4_bigFloatCeil, C.v1_4_bigFloatCeil)
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

	imports, err = imports.Append("bigFloatGetConstPi", v1_4_bigFloatGetConstPi, C.v1_4_bigFloatGetConstPi)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatGetConstE", v1_4_bigFloatGetConstE, C.v1_4_bigFloatGetConstE)
	if err != nil {
		return nil, err
	}

	return imports, err
}

func areAllZero(values ...*big.Float) bool {
	for _, val := range values {
		if val.Sign() != 0 {
			return false
		}
	}
	return true
}

func setResultIfNotInfinity(host arwen.VMHost, result *big.Float, destinationHandle int32) {
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	if result.IsInf() {
		_ = arwen.WithFaultAndHost(host, arwen.ErrInfinityFloatOperation, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFaultAndHost(host, err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	managedType.ConsumeGasForBigFloatCopy(result)
	dest.Set(result)
}

//export v1_4_bigFloatNewFromParts
func v1_4_bigFloatNewFromParts(context unsafe.Pointer, integralPart, fractionalPart, exponent int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGas(gasToUse)

	if exponent > 0 {
		_ = arwen.WithFault(arwen.ErrPositiveExponent, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}
	var fractional float64
	if exponent < -322 {
		fractional = float64(0)
	} else {
		fractional = float64(fractionalPart) * math.Pow10(int(exponent))
	}

	value := float64(integralPart) + fractional
	return managedType.PutBigFloat(value)
}

//export v1_4_bigFloatNewFromFrac
func v1_4_bigFloatNewFromFrac(context unsafe.Pointer, numerator, denominator int64) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGas(gasToUse)

	if denominator == 0 {
		_ = arwen.WithFault(arwen.ErrDivZero, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}

	value := float64(numerator) / float64(denominator)
	return managedType.PutBigFloat(value)
}

//export v1_4_bigFloatNewFromSci
func v1_4_bigFloatNewFromSci(context unsafe.Pointer, significand, exponent int64) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGas(gasToUse)

	if exponent > 0 {
		_ = arwen.WithFault(arwen.ErrPositiveExponent, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}
	if exponent < -322 {
		return managedType.PutBigFloat(float64(0))
	}

	value := float64(significand) * math.Pow10(int(exponent))
	return managedType.PutBigFloat(value)
}

//export v1_4_bigFloatAdd
func v1_4_bigFloatAdd(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext((context))

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAdd
	metering.UseGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	managedType.ConsumeGasForBigFloatCopy(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultAdd := new(big.Float).Add(op1, op2)
	setResultIfNotInfinity(arwen.GetVMHost(context), resultAdd, destinationHandle)
}

//export v1_4_bigFloatSub
func v1_4_bigFloatSub(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSub
	metering.UseGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	managedType.ConsumeGasForBigFloatCopy(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultSub := new(big.Float).Sub(op1, op2)
	setResultIfNotInfinity(arwen.GetVMHost(context), resultSub, destinationHandle)
}

//export v1_4_bigFloatMul
func v1_4_bigFloatMul(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatMul
	metering.UseGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	managedType.ConsumeGasForBigFloatCopy(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultMul := new(big.Float).Mul(op1, op2)
	setResultIfNotInfinity(arwen.GetVMHost(context), resultMul, destinationHandle)
}

//export v1_4_bigFloatDiv
func v1_4_bigFloatDiv(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatDiv
	metering.UseGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	managedType.ConsumeGasForBigFloatCopy(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if areAllZero(op1, op2) {
		_ = arwen.WithFault(arwen.ErrAllOperandsAreEqualToZero, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}

	resultDiv := new(big.Float).Quo(op1, op2)
	setResultIfNotInfinity(arwen.GetVMHost(context), resultDiv, destinationHandle)
}

//export v1_4_bigFloatNeg
func v1_4_bigFloatNeg(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNeg
	metering.UseGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(dest, op)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Neg(op)
}

//export v1_4_bigFloatClone
func v1_4_bigFloatClone(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatClone
	metering.UseGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(dest, op)
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

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	managedType.ConsumeGasForBigFloatCopy(op1, op2)
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

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(dest, op)
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

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(dest, op)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if op.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	dest.Sqrt(op)
}

//export v1_4_bigFloatPow
func v1_4_bigFloatPow(context unsafe.Pointer, destinationHandle, opHandle, exponent int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntPow
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(op)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	opBigInt := big.NewInt(0)
	op.Int(opBigInt)
	op2BigInt := big.NewInt(int64(exponent))
	if opBigInt.Sign() > 0 {
		opBigInt.Add(opBigInt, big.NewInt(1))
	}

	//this calculates the length of the result in bytes
	lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(op2BigInt, big.NewInt(int64(opBigInt.BitLen()))), big.NewInt(8))
	managedType.ConsumeGasForThisBigIntNumberOfBytes(lengthOfResult)

	powResult := new(big.Float).Set(pow(op, exponent))
	setResultIfNotInfinity(arwen.GetVMHost(context), powResult, destinationHandle)
}

func pow(base *big.Float, exp int32) *big.Float {
	result := new(big.Float).Copy(base)
	for i := 0; i < int(exp); i++ {
		result.Mul(result, base)
	}
	return result
}

//export v1_4_bigFloatFloor
func v1_4_bigFloatFloor(context unsafe.Pointer, opHandle, destBigIntHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatFloor
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(op)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	bigIntOp := managedType.GetBigIntOrCreate(destBigIntHandle)

	op.Int(bigIntOp)
	managedType.ConsumeGasForBigIntCopy(bigIntOp)
	if op.IsInt() {
		return
	}
	if bigIntOp.Sign() < 0 {
		bigIntOp.Sub(bigIntOp, big.NewInt(1))
	}
}

//export v1_4_bigFloatCeil
func v1_4_bigFloatCeil(context unsafe.Pointer, opHandle, destBigIntHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCeil
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(op)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	bigIntOp := managedType.GetBigIntOrCreate(destBigIntHandle)

	op.Int(bigIntOp)
	managedType.ConsumeGasForBigIntCopy(bigIntOp)
	if op.IsInt() {
		return
	}
	if bigIntOp.Sign() > 0 {
		bigIntOp.Add(bigIntOp, big.NewInt(1))
	}
}

//export v1_4_bigFloatTruncate
func v1_4_bigFloatTruncate(context unsafe.Pointer, opHandle, destBigIntHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatTruncate
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	managedType.ConsumeGasForBigFloatCopy(op)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	bigIntValue := managedType.GetBigIntOrCreate(destBigIntHandle)

	op.Int(bigIntValue)
	managedType.ConsumeGasForBigIntCopy(bigIntValue)
}

//export v1_4_bigFloatSetInt64
func v1_4_bigFloatSetInt64(context unsafe.Pointer, destinationHandle int32, value int64) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetInt64
	metering.UseGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
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
	managedType.ConsumeGasForBigFloatCopy(op)
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

	bigIntValue, err := managedType.GetBigInt(bigIntHandle)
	managedType.ConsumeGasForBigIntCopy(bigIntValue)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	resultSetInt := big.NewFloat(0).SetInt(bigIntValue)
	setResultIfNotInfinity(arwen.GetVMHost(context), resultSetInt, destinationHandle)
}

//export v1_4_bigFloatGetConstPi
func v1_4_bigFloatGetConstPi(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGas(gasToUse)

	pi, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	pi.SetFloat64(math.Pi)
}

//export v1_4_bigFloatGetConstE
func v1_4_bigFloatGetConstE(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGas(gasToUse)

	e, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	e.SetFloat64(math.E)
}
