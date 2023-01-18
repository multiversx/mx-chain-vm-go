package vmhooks

import (
	"math"
	"math/big"

	"github.com/multiversx/wasm-vm/arwen"
	arwenMath "github.com/multiversx/wasm-vm/math"
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
		_ = WithFaultAndHost(host, arwen.ErrInfinityFloatOperation, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}

	exponent := result.MantExp(nil)
	if managedType.BigFloatExpIsNotValid(exponent) {
		_ = WithFaultAndHost(host, arwen.ErrExponentTooBigOrTooSmall, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if WithFaultAndHost(host, err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	dest.Set(result)
}

// BigFloatNewFromParts VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatNewFromParts(integralPart, fractionalPart, exponent int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGasAndAddTracedGas(bigFloatNewFromPartsName, gasToUse)

	if exponent > 0 {
		_ = context.WithFault(arwen.ErrPositiveExponent, runtime.BigFloatAPIErrorShouldFailExecution())
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
		if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
	}

	var value *big.Float
	if integralPart >= 0 {
		value, err = arwenMath.AddBigFloat(big.NewFloat(float64(integralPart)), bigFractional)
		if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
	} else {
		value, err = arwenMath.SubBigFloat(big.NewFloat(float64(integralPart)), bigFractional)
		if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
	}
	handle, err := managedType.PutBigFloat(value)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	return handle
}

// BigFloatNewFromFrac VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatNewFromFrac(numerator, denominator int64) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGasAndAddTracedGas(bigFloatNewFromFracName, gasToUse)

	if denominator == 0 {
		_ = context.WithFault(arwen.ErrDivZero, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}

	bigNumerator := big.NewFloat(float64(numerator))
	bigDenominator := big.NewFloat(float64(denominator))
	value, err := arwenMath.QuoBigFloat(bigNumerator, bigDenominator)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	handle, err := managedType.PutBigFloat(value)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	return handle
}

// BigFloatNewFromSci VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatNewFromSci(significand, exponent int64) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	metering.UseGasAndAddTracedGas(bigFloatNewFromSciName, gasToUse)

	if exponent > 0 {
		_ = context.WithFault(arwen.ErrPositiveExponent, runtime.BigFloatAPIErrorShouldFailExecution())
		return -1
	}
	if exponent < -322 {
		handle, err := managedType.PutBigFloat(big.NewFloat(0))
		if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
			return -1
		}
		return handle
	}

	bigSignificand := big.NewFloat(float64(significand))
	bigExponentMultiplier := big.NewFloat(math.Pow10(int(exponent)))
	value, err := arwenMath.MulBigFloat(bigSignificand, bigExponentMultiplier)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	handle, err := managedType.PutBigFloat(value)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	return handle
}

// BigFloatAdd VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatAdd(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatAddName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAdd
	metering.UseGasAndAddTracedGas(bigFloatAddName, gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultAdd, err := arwenMath.AddBigFloat(op1, op2)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	setResultIfNotInfinity(context.GetVMHost(), resultAdd, destinationHandle)
}

// BigFloatSub VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatSub(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatSubName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSub
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultSub, err := arwenMath.SubBigFloat(op1, op2)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(context.GetVMHost(), resultSub, destinationHandle)
}

// BigFloatMul VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatMul(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatMulName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatMul
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}

	resultMul, err := arwenMath.MulBigFloat(op1, op2)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(context.GetVMHost(), resultMul, destinationHandle)
}

// BigFloatDiv VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatDiv(destinationHandle, op1Handle, op2Handle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatDivName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatDiv
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if areAllZero(op1, op2) {
		_ = context.WithFault(arwen.ErrAllOperandsAreEqualToZero, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}

	resultDiv, err := arwenMath.QuoBigFloat(op1, op2)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(context.GetVMHost(), resultDiv, destinationHandle)
}

// BigFloatNeg VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatNeg(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatNegName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNeg
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Neg(op)
}

// BigFloatClone VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatClone(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatCloneName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatClone
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Copy(op)
}

// BigFloatCmp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatCmp(op1Handle, op2Handle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatCmpName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCmp
	metering.UseAndTraceGas(gasToUse)

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -2
	}
	return int32(op1.Cmp(op2))
}

// BigFloatAbs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatAbs(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatAbsName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Abs(op)
}

// BigFloatSign VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatSign(opHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	metering.UseGasAndAddTracedGas(bigFloatSignName, gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -2
	}
	return int32(op.Sign())
}

// BigFloatSqrt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatSqrt(destinationHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatSqrtName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSqrt
	metering.UseAndTraceGas(gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	op, err := managedType.GetBigFloat(opHandle)

	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	if op.Sign() < 0 {
		_ = context.WithFault(arwen.ErrBadLowerBounds, runtime.BigFloatAPIErrorShouldFailExecution())
		return
	}
	resultSqrt, err := arwenMath.SqrtBigFloat(op)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.Set(resultSqrt)
}

// BigFloatPow VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatPow(destinationHandle, opHandle, exponent int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatPowName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatPow
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
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

	powResult, err := context.pow(op, exponent)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	setResultIfNotInfinity(context.GetVMHost(), powResult, destinationHandle)
}

func (context *ElrondApi) pow(base *big.Float, exp int32) (*big.Float, error) {
	result := big.NewFloat(1)
	result.SetPrec(base.Prec())
	managedType := context.GetManagedTypesContext()

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

// BigFloatFloor VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatFloor(destBigIntHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatFloorName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatFloor
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
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

// BigFloatCeil VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatCeil(destBigIntHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatCeilName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCeil
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
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

// BigFloatTruncate VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatTruncate(destBigIntHandle, opHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatTruncateName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatTruncate
	metering.UseAndTraceGas(gasToUse)

	op, err := managedType.GetBigFloat(opHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	bigIntValue := managedType.GetBigIntOrCreate(destBigIntHandle)

	op.Int(bigIntValue)
	managedType.ConsumeGasForBigIntCopy(bigIntValue)
}

// BigFloatSetInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatSetInt64(destinationHandle int32, value int64) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetInt64
	metering.UseGasAndAddTracedGas(bigFloatSetInt64Name, gasToUse)

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	dest.SetInt64(value)
}

// BigFloatIsInt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatIsInt(opHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatIsIntName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatIsInt
	metering.UseAndTraceGas(gasToUse)
	op, err := managedType.GetBigFloat(opHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return -1
	}
	if op.IsInt() {
		return 1
	}
	return 0
}

// BigFloatSetBigInt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatSetBigInt(destinationHandle, bigIntHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(bigFloatSetBigIntName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetBigInt
	metering.UseAndTraceGas(gasToUse)

	bigIntValue, err := managedType.GetBigInt(bigIntHandle)
	managedType.ConsumeGasForBigIntCopy(bigIntValue)
	if context.WithFault(err, runtime.BigIntAPIErrorShouldFailExecution()) {
		return
	}
	resultSetInt := big.NewFloat(0).SetInt(bigIntValue)
	setResultIfNotInfinity(context.GetVMHost(), resultSetInt, destinationHandle)
}

// BigFloatGetConstPi VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatGetConstPi(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatGetConst
	metering.UseGasAndAddTracedGas(bigFloatGetConstPiName, gasToUse)

	pi, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	pi.SetFloat64(math.Pi)
}

// BigFloatGetConstE VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) BigFloatGetConstE(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatGetConst
	metering.UseGasAndAddTracedGas(bigFloatGetConstEName, gasToUse)

	e, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if context.WithFault(err, runtime.BigFloatAPIErrorShouldFailExecution()) {
		return
	}
	e.SetFloat64(math.E)
}
