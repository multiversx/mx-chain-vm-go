package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_4_bigFloatNew(void* context, int32_t intBase, int32_t subIntBase, int32_t exponent);
// extern int32_t	v1_4_bigFloatNewFromFrac(void* context, long long numerator, long long denominator);
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
// extern int32_t	v1_4_bigFloatLog2(void* context, int32_t opHandle);
// extern void		v1_4_bigFloatPow(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
//
// extern void		v1_4_bigFloatFloor(void* context, int32_t opHandle, int32_t bigIntHandle);
// extern void		v1_4_bigFloatCeil(void* context, int32_t opHandle, int32_t bigIntHandle);
//
// extern int32_t	v1_4_bigFloatIsInt(void* context, int32_t opHandle);
// extern void		v1_4_bigFloatSetInt64(void* context, int32_t destinationHandle, long long value);
// extern void		v1_4_bigFloatSetBigInt(void* context, int32_t destinationHandle, int32_t bigIntHandle);
//
// extern void		v1_4_bigFloatGetConstPi(void* context, int32_t destinationHandle);
// extern void		v1_4_bigFloatGetConstE(void* context, int32_t destinationHandle);
//
// extern void		v1_4_bigFloatSetBytes(void* context, int32_t destinationHandle, int32_t dataOffset, int32_t dataLength);
// extern void		v1_4_bigFloatGetBytes(void* context, int32_t destinationHandle, int32_t dataOffset);
import "C"
import (
	"math"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	arwenMath "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
)

func BigFloatImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("bigFloatNew", v1_4_bigFloatNew, C.v1_4_bigFloatNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatNewFromFrac", v1_4_bigFloatNewFromFrac, C.v1_4_bigFloatNewFromFrac)
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

	imports, err = imports.Append("bigFloatPow", v1_4_bigFloatPow, C.v1_4_bigFloatPow)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatLog2", v1_4_bigFloatLog2, C.v1_4_bigFloatLog2)
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

	imports, err = imports.Append("bigFloatSetBytes", v1_4_bigFloatSetBytes, C.v1_4_bigFloatSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigFloatGetBytes", v1_4_bigFloatGetBytes, C.v1_4_bigFloatGetBytes)
	if err != nil {
		return nil, err
	}

	return imports, err
}

func oneIsInfinity(values ...*big.Float) bool {
	for _, val := range values {
		if val.IsInf() {
			return true
		}
	}
	return false
}

func allAreEqualToZero(values ...*big.Float) bool {
	for _, val := range values {
		if val.Sign() != 0 {
			return false
		}
	}
	return true
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

//export v1_4_bigFloatNewFromFrac
func v1_4_bigFloatNewFromFrac(context unsafe.Pointer, numerator, denominator int64) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNew
	metering.UseGas(gasToUse)

	if denominator == 0 {
		_ = arwen.WithFault(arwen.ErrDivZero, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}

	value := float64(numerator) / float64(denominator)
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
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(dest, op1, op2)
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
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(dest, op1, op2)
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
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(dest, op1, op2)
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
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	if allAreEqualToZero(op1, op2) {
		_ = arwen.WithFault(arwen.ErrAllOperandsAreEqualToZero, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(dest, op1, op2)
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
	managedType.ConsumeGasForBigFloatCopy(dest, op)
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
	managedType.ConsumeGasForBigFloatCopy(dest, op)
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
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -2
	}
	managedType.ConsumeGasForBigFloatCopy(op1, op2)
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
	managedType.ConsumeGasForBigFloatCopy(dest, op)
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
	managedType.ConsumeGasForBigFloatCopy(op)
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
		_ = arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	if oneIsInfinity(op) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(dest, op)
	dest.Sqrt(op)
}

//export v1_4_bigFloatLog2
func v1_4_bigFloatLog2(context unsafe.Pointer, opHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatLog2
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	if op.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}
	if oneIsInfinity(op) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}
	bigIntOp := new(big.Int)
	op.Int(bigIntOp)
	managedType.ConsumeGasForBigFloatCopy(op)
	managedType.ConsumeGasForBigIntCopy(bigIntOp)
	return int32(bigIntOp.BitLen() - 1)
}

//export v1_4_bigFloatPow
func v1_4_bigFloatPow(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntPow
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	op1, op2, err := managedType.GetTwoBigFloat(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(op1, op2, dest)

	op1BigInt := big.NewInt(0)
	op1.Int(op1BigInt)
	op2BigInt := big.NewInt(0)
	op1.Int(op2BigInt)
	if op1BigInt.Sign() > 0 {
		op1BigInt.Add(op1BigInt, big.NewInt(1))
	}
	if op2BigInt.Sign() > 0 {
		op2BigInt.Add(op2BigInt, big.NewInt(1))
	}

	//this calculates the length of the result in bytes
	lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(op2BigInt, big.NewInt(int64(op1BigInt.BitLen()))), big.NewInt(8))
	managedType.ConsumeGasForThisBigIntNumberOfBytes(lengthOfResult)

	dest.Set(pow(op1, op2))
}

func pow(base *big.Float, exp *big.Float) *big.Float {
	result := new(big.Float).Copy(base)
	counter := big.NewFloat(1)
	for counter.Cmp(exp) < 0 {
		result.Mul(result, base)
		counter.Add(counter, big.NewFloat(1))
	}
	return result
}

//export v1_4_bigFloatFloor
func v1_4_bigFloatFloor(context unsafe.Pointer, opHandle, bigIntHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatFloor
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	bigIntOp := managedType.GetBigIntOrCreate(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if oneIsInfinity(op) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(op)
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
func v1_4_bigFloatCeil(context unsafe.Pointer, opHandle, bigIntHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCeil
	metering.UseGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	bigIntOp := managedType.GetBigIntOrCreate(bigIntHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if oneIsInfinity(op) {
		_ = arwen.WithFault(arwen.ErrInfinityFloatOperation, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(op)
	op.Int(bigIntOp)
	managedType.ConsumeGasForBigIntCopy(bigIntOp)
	if op.IsInt() {
		return
	}
	if bigIntOp.Sign() > 0 {
		bigIntOp.Add(bigIntOp, big.NewInt(1))
	}
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
	managedType.ConsumeGasForBigFloatCopy(op)
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
	managedType.ConsumeGasForBigFloatCopy(dest)
	managedType.ConsumeGasForBigIntCopy(bigIntValue)
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
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(dest, op1, op2)
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
	if oneIsInfinity(op1, op2) {
		_ = arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	managedType.ConsumeGasForBigFloatCopy(dest, op1, op2)
	dest.Quo(op1, op2)
	rDiv := big.NewInt(0)
	dest.Int(rDiv)
	dest.Sub(dest, new(big.Float).SetInt(rDiv))
}

//export v1_4_bigFloatGetConstPi
func v1_4_bigFloatGetConstPi(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNew
	metering.UseGas(gasToUse)

	pi := managedType.GetBigFloatOrCreate(destinationHandle)
	pi.SetFloat64(math.Pi)
}

//export v1_4_bigFloatGetConstE
func v1_4_bigFloatGetConstE(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNew
	metering.UseGas(gasToUse)

	pi := managedType.GetBigFloatOrCreate(destinationHandle)
	pi.SetFloat64(math.E)
}

//export v1_4_bigFloatSetBytes
func v1_4_bigFloatSetBytes(context unsafe.Pointer, destinationHandle, dataOffset, dataLength int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetBytes
	metering.UseGas(gasToUse)

	dest := managedType.GetBigFloatOrCreate(destinationHandle)
	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	gasToUse = arwenMath.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseGas(gasToUse)

	floatToBeSet := new(big.Float)
	err = floatToBeSet.GobDecode(data)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	dest.Set(floatToBeSet)
}

//export v1_4_bigFloatGetBytes
func v1_4_bigFloatGetBytes(context unsafe.Pointer, destinationHandle, dataOffset int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatGetBytes
	metering.UseGas(gasToUse)

	dest, err := managedType.GetBigFloat(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	data, err := dest.GobEncode()
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	err = runtime.MemStore(dataOffset, data)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	gasToUse = arwenMath.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseGas(gasToUse)
}
