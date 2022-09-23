package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t	v1_5_bigFloatNewFromParts(void* context, int32_t integralPart, int32_t fractionalPart, int32_t exponent);
// extern int32_t	v1_5_bigFloatNewFromFrac(void* context, long long numerator, long long denominator);
// extern int32_t	v1_5_bigFloatNewFromSci(void* context, long long significand, long long exponent);
//
// extern void		v1_5_bigFloatAdd(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_5_bigFloatSub(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_5_bigFloatMul(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void		v1_5_bigFloatDiv(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
//
// extern void		v1_5_bigFloatAbs(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_5_bigFloatNeg(void* context, int32_t destinationHandle, int32_t opHandle);
// extern int32_t	v1_5_bigFloatCmp(void* context, int32_t op1Handle, int32_t op2Handle);
// extern int32_t	v1_5_bigFloatSign(void* context, int32_t opHandle);
// extern void 		v1_5_bigFloatClone(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_5_bigFloatSqrt(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void		v1_5_bigFloatPow(void* context, int32_t destinationHandle, int32_t opHandle, int32_t exponent);
//
// extern void		v1_5_bigFloatFloor(void* context, int32_t destBigIntHandle, int32_t opHandle);
// extern void		v1_5_bigFloatCeil(void* context, int32_t destBigIntHandle, int32_t opHandle);
// extern void		v1_5_bigFloatTruncate(void* context, int32_t destBigIntHandle, int32_t opHandle);
//
// extern int32_t	v1_5_bigFloatIsInt(void* context, int32_t opHandle);
// extern void		v1_5_bigFloatSetInt64(void* context, int32_t destinationHandle, long long value);
// extern void		v1_5_bigFloatSetBigInt(void* context, int32_t destinationHandle, int32_t bigIntHandle);
//
// extern void		v1_5_bigFloatGetConstPi(void* context, int32_t destinationHandle);
// extern void		v1_5_bigFloatGetConstE(void* context, int32_t destinationHandle);
import "C"
import (
	"math"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/executor"
	arwenMath "github.com/ElrondNetwork/wasm-vm/math"
)

const (
	bigFloatNewFromPartsName = "bigFloatNewFromParts"
	bigFloatNewFromFracName  = "bigFloatNewFromFrac"
	bigFloatNewFromSciName   = "bigFloatNewFromSci"
	bigFloatAddName          = "bigFloatAdd"
	bigFloatSubName          = "bigFloatSub"
	bigFloatMulName          = "bigFloatMul"
	bigFloatDivName          = "bigFloatDiv"
	bigFloatAbsName          = "bigFloatAbs"
	bigFloatNegName          = "bigFloatNeg"
	bigFloatCmpName          = "bigFloatCmp"
	bigFloatSignName         = "bigFloatSign"
	bigFloatCloneName        = "bigFloatClone"
	bigFloatSqrtName         = "bigFloatSqrt"
	bigFloatPowName          = "bigFloatPow"
	bigFloatFloorName        = "bigFloatFloor"
	bigFloatCeilName         = "bigFloatCeil"
	bigFloatTruncateName     = "bigFloatTruncate"
	bigFloatIsIntName        = "bigFloatIsInt"
	bigFloatSetInt64Name     = "bigFloatSetInt64"
	bigFloatSetBigIntName    = "bigFloatSetBigInt"
	bigFloatGetConstPiName   = "bigFloatGetConstPi"
	bigFloatGetConstEName    = "bigFloatGetConstE"
)

// BigFloatImports populates imports with the BigFloat API methods
func BigFloatImports(imports executor.ImportFunctionReceiver) error {
	imports.Namespace("env")

	err := imports.Append("bigFloatNewFromParts", v1_5_bigFloatNewFromParts, C.v1_5_bigFloatNewFromParts)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatNewFromFrac", v1_5_bigFloatNewFromFrac, C.v1_5_bigFloatNewFromFrac)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatNewFromSci", v1_5_bigFloatNewFromSci, C.v1_5_bigFloatNewFromSci)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatAdd", v1_5_bigFloatAdd, C.v1_5_bigFloatAdd)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSub", v1_5_bigFloatSub, C.v1_5_bigFloatSub)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatMul", v1_5_bigFloatMul, C.v1_5_bigFloatMul)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatDiv", v1_5_bigFloatDiv, C.v1_5_bigFloatDiv)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatTruncate", v1_5_bigFloatTruncate, C.v1_5_bigFloatTruncate)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatNeg", v1_5_bigFloatNeg, C.v1_5_bigFloatNeg)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatClone", v1_5_bigFloatClone, C.v1_5_bigFloatClone)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatCmp", v1_5_bigFloatCmp, C.v1_5_bigFloatCmp)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatAbs", v1_5_bigFloatAbs, C.v1_5_bigFloatAbs)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSign", v1_5_bigFloatSign, C.v1_5_bigFloatSign)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSqrt", v1_5_bigFloatSqrt, C.v1_5_bigFloatSqrt)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatPow", v1_5_bigFloatPow, C.v1_5_bigFloatPow)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatFloor", v1_5_bigFloatFloor, C.v1_5_bigFloatFloor)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatCeil", v1_5_bigFloatCeil, C.v1_5_bigFloatCeil)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSetInt64", v1_5_bigFloatSetInt64, C.v1_5_bigFloatSetInt64)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatIsInt", v1_5_bigFloatIsInt, C.v1_5_bigFloatIsInt)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSetBigInt", v1_5_bigFloatSetBigInt, C.v1_5_bigFloatSetBigInt)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatGetConstPi", v1_5_bigFloatGetConstPi, C.v1_5_bigFloatGetConstPi)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatGetConstE", v1_5_bigFloatGetConstE, C.v1_5_bigFloatGetConstE)
	if err != nil {
		return err
	}

	return err
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

	exponent := result.MantExp(nil)
	if managedType.BigFloatExpIsNotValid(exponent) {
		_ = arwen.WithFaultAndHost(host, arwen.ErrExponentTooBigOrTooSmall, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFaultAndHost(host, err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	dest.Set(result)
}

//export v1_5_bigFloatNewFromParts
func v1_5_bigFloatNewFromParts(context unsafe.Pointer, integralPart, fractionalPart, exponent int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGasAndAddTracedGas(bigFloatNewFromPartsName, gasToUse)

	if exponent > 0 {
		_ = arwen.WithFault(arwen.ErrPositiveExponent, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}
	var err error
	var bigFractional *big.Float
	if exponent < -322 {
		bigFractional = big.NewFloat(0)
	} else {
		bigFractionalPart := big.NewFloat(float64(fractionalPart))
		bigExponentMultiplier := big.NewFloat(math.Pow10(int(exponent)))
		bigFractional, err = arwenMath.MulBigFloat(bigFractionalPart, bigExponentMultiplier)
		if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
	}

	var value *big.Float
	if integralPart >= 0 {
		value, err = arwenMath.AddBigFloat(big.NewFloat(float64(integralPart)), bigFractional)
		if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
	} else {
		value, err = arwenMath.SubBigFloat(big.NewFloat(float64(integralPart)), bigFractional)
		if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
	}
	handle, err := managedType.PutBigFloat(value)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	return handle
}

//export v1_5_bigFloatNewFromFrac
func v1_5_bigFloatNewFromFrac(context unsafe.Pointer, numerator, denominator int64) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGasAndAddTracedGas(bigFloatNewFromFracName, gasToUse)

	if denominator == 0 {
		_ = arwen.WithFault(arwen.ErrDivZero, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}

	bigNumerator := big.NewFloat(float64(numerator))
	bigDenominator := big.NewFloat(float64(denominator))
	value, err := arwenMath.QuoBigFloat(bigNumerator, bigDenominator)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	handle, err := managedType.PutBigFloat(value)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	return handle
}

//export v1_5_bigFloatNewFromSci
func v1_5_bigFloatNewFromSci(context unsafe.Pointer, significand, exponent int64) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGasAndAddTracedGas(bigFloatNewFromSciName, gasToUse)

	if exponent > 0 {
		_ = arwen.WithFault(arwen.ErrPositiveExponent, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}
	if exponent < -322 {
		handle, err := managedType.PutBigFloat(big.NewFloat(0))
		if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
		return handle
	}

	bigSignificand := big.NewFloat(float64(significand))
	bigExponentMultiplier := big.NewFloat(math.Pow10(int(exponent)))
	value, err := arwenMath.MulBigFloat(bigSignificand, bigExponentMultiplier)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	handle, err := managedType.PutBigFloat(value)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	return handle
}

//export v1_5_bigFloatAdd
func v1_5_bigFloatAdd(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext((context))
	metering.StartGasTracing(bigFloatAddName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAdd
	metering.UseGasAndAddTracedGas(bigFloatAddName, gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultAdd, err := arwenMath.AddBigFloat(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	setResultIfNotInfinity(arwen.GetVMHost(context), resultAdd, destinationHandle)
}

//export v1_5_bigFloatSub
func v1_5_bigFloatSub(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatSubName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSub
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultSub, err := arwenMath.SubBigFloat(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(arwen.GetVMHost(context), resultSub, destinationHandle)
}

//export v1_5_bigFloatMul
func v1_5_bigFloatMul(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatMulName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatMul
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultMul, err := arwenMath.MulBigFloat(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(arwen.GetVMHost(context), resultMul, destinationHandle)
}

//export v1_5_bigFloatDiv
func v1_5_bigFloatDiv(context unsafe.Pointer, destinationHandle, op1Handle, op2Handle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatDivName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatDiv
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if areAllZero(op1, op2) {
		_ = arwen.WithFault(arwen.ErrAllOperandsAreEqualToZero, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}

	resultDiv, err := arwenMath.QuoBigFloat(op1, op2)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(arwen.GetVMHost(context), resultDiv, destinationHandle)
}

//export v1_5_bigFloatNeg
func v1_5_bigFloatNeg(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatNegName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNeg
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Neg(op)
}

//export v1_5_bigFloatClone
func v1_5_bigFloatClone(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatCloneName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatClone
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Copy(op)
}

//export v1_5_bigFloatCmp
func v1_5_bigFloatCmp(context unsafe.Pointer, op1Handle, op2Handle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatCmpName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCmp
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -2
	}
	return int32(op1.Cmp(op2))
}

//export v1_5_bigFloatAbs
func v1_5_bigFloatAbs(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatAbsName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Abs(op)
}

//export v1_5_bigFloatSign
func v1_5_bigFloatSign(context unsafe.Pointer, opHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	metering.UseGasAndAddTracedGas(bigFloatSignName, gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -2
	}
	return int32(op.Sign())
}

//export v1_5_bigFloatSqrt
func v1_5_bigFloatSqrt(context unsafe.Pointer, destinationHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatSqrtName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSqrt
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if op.Sign() < 0 {
		_ = arwen.WithFault(arwen.ErrBadLowerBounds, context, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	resultSqrt, err := arwenMath.SqrtBigFloat(op)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Set(resultSqrt)
}

//export v1_5_bigFloatPow
func v1_5_bigFloatPow(context unsafe.Pointer, destinationHandle, opHandle, exponent int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatPowName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatPow
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
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

	powResult, err := pow(context, op, exponent)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(arwen.GetVMHost(context), powResult, destinationHandle)
}

func pow(context unsafe.Pointer, base *big.Float, exp int32) (*big.Float, error) {
	result := big.NewFloat(1)
	result.SetPrec(base.Prec())
	managedType := arwen.GetManagedTypesContext(context)

	for i := 0; i < int(exp); i++ {
		resultMul, err := arwenMath.MulBigFloat(result, base)
		if err != nil {
			return nil, err
		}
		exponent := resultMul.MantExp(nil)
		if managedType.BigFloatExpIsNotValid(exponent) {
			return nil, arwen.ErrExponentTooBigOrTooSmall
		}
		result.Set(resultMul)
	}
	return result, nil
}

//export v1_5_bigFloatFloor
func v1_5_bigFloatFloor(context unsafe.Pointer, destBigIntHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatFloorName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatFloor
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
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

//export v1_5_bigFloatCeil
func v1_5_bigFloatCeil(context unsafe.Pointer, destBigIntHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatCeilName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCeil
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
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

//export v1_5_bigFloatTruncate
func v1_5_bigFloatTruncate(context unsafe.Pointer, destBigIntHandle, opHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatTruncateName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatTruncate
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	bigIntValue := managedType.GetBigIntOrCreate(destBigIntHandle)

	op.Int(bigIntValue)
	managedType.ConsumeGasForBigIntCopy(bigIntValue)
}

//export v1_5_bigFloatSetInt64
func v1_5_bigFloatSetInt64(context unsafe.Pointer, destinationHandle int32, value int64) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetInt64
	metering.UseGasAndAddTracedGas(bigFloatSetInt64Name, gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.SetInt64(value)
}

//export v1_5_bigFloatIsInt
func v1_5_bigFloatIsInt(context unsafe.Pointer, opHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatIsIntName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatIsInt
	metering.UseAndTraceGas(gasToUse)
	op, err := managedType.GetBigFloat(opHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	if op.IsInt() {
		return 1
	}
	return 0
}

//export v1_5_bigFloatSetBigInt
func v1_5_bigFloatSetBigInt(context unsafe.Pointer, destinationHandle, bigIntHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(bigFloatSetBigIntName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetBigInt
	metering.UseAndTraceGas(gasToUse)

	bigIntValue, err := managedType.GetBigInt(bigIntHandle)
	managedType.ConsumeGasForBigIntCopy(bigIntValue)
	if arwen.WithFault(err, context, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	resultSetInt := big.NewFloat(0).SetInt(bigIntValue)
	setResultIfNotInfinity(arwen.GetVMHost(context), resultSetInt, destinationHandle)
}

//export v1_5_bigFloatGetConstPi
func v1_5_bigFloatGetConstPi(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatGetConst
	metering.UseGasAndAddTracedGas(bigFloatGetConstPiName, gasToUse)

	pi, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	pi.SetFloat64(math.Pi)
}

//export v1_5_bigFloatGetConstE
func v1_5_bigFloatGetConstE(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatGetConst
	metering.UseGasAndAddTracedGas(bigFloatGetConstEName, gasToUse)

	e, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if arwen.WithFault(err, context, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	e.SetFloat64(math.E)
}
