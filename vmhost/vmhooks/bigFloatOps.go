package vmhooks

import (
	"math"
	"math/big"

	vmMath "github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
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

func setResultIfNotInfinity(host vmhost.VMHost, result *big.Float, destinationHandle int32) error {
	managedType := host.ManagedTypes()
	if result.IsInf() {
		return vmhost.ErrInfinityFloatOperation
	}

	exponent := result.MantExp(nil)
	if managedType.BigFloatExpIsNotValid(exponent) {
		return vmhost.ErrExponentTooBigOrTooSmall
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}

	dest.Set(result)
	return nil
}

// BigFloatNewFromParts VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatNewFromParts(integralPart, fractionalPart, exponent int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatNewFromPartsName, gasToUse)
	if err != nil {
		return -1, err
	}

	if exponent > 0 {
		return -1, vmhost.ErrPositiveExponent
	}

	var bigFractional *big.Float
	if exponent < -322 {
		bigFractional = big.NewFloat(0)
	} else {
		bigFractionalPart := big.NewFloat(float64(fractionalPart))
		bigExponentMultiplier := big.NewFloat(math.Pow10(int(exponent)))
		bigFractional, err = vmMath.MulBigFloat(bigFractionalPart, bigExponentMultiplier)
		if err != nil {
			return -1, err
		}
	}

	var value *big.Float
	if integralPart >= 0 {
		value, err = vmMath.AddBigFloat(big.NewFloat(float64(integralPart)), bigFractional)
		if err != nil {
			return -1, err
		}
	} else {
		value, err = vmMath.SubBigFloat(big.NewFloat(float64(integralPart)), bigFractional)
		if err != nil {
			return -1, err
		}
	}
	handle, err := managedType.PutBigFloat(value)
	if err != nil {
		return -1, err
	}
	return handle, nil
}

// BigFloatNewFromFrac VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatNewFromFrac(numerator, denominator int64) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatNewFromFracName, gasToUse)
	if err != nil {
		return -1, err
	}

	if denominator == 0 {
		return -1, vmhost.ErrDivZero
	}

	bigNumerator := big.NewFloat(float64(numerator))
	bigDenominator := big.NewFloat(float64(denominator))
	value, err := vmMath.QuoBigFloat(bigNumerator, bigDenominator)
	if err != nil {
		return -1, err
	}
	handle, err := managedType.PutBigFloat(value)
	if err != nil {
		return -1, err
	}
	return handle, nil
}

// BigFloatNewFromSci VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatNewFromSci(significand, exponent int64) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNewFromParts
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatNewFromSciName, gasToUse)
	if err != nil {
		return -1, err
	}

	if exponent > 0 {
		return -1, vmhost.ErrPositiveExponent
	}
	if exponent < -322 {
		handle, err := managedType.PutBigFloat(big.NewFloat(0))
		if err != nil {
			return -1, err
		}
		return handle, nil
	}

	bigSignificand := big.NewFloat(float64(significand))
	bigExponentMultiplier := big.NewFloat(math.Pow10(int(exponent)))
	value, err := vmMath.MulBigFloat(bigSignificand, bigExponentMultiplier)
	if err != nil {
		return -1, err
	}
	handle, err := managedType.PutBigFloat(value)
	if err != nil {
		return -1, err
	}
	return handle, nil
}

// BigFloatAdd VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatAdd(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatAddName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAdd
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatAddName, gasToUse)
	if err != nil {
		return err
	}

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	resultAdd, err := vmMath.AddBigFloat(op1, op2)
	if err != nil {
		return err
	}

	return setResultIfNotInfinity(context.GetVMHost(), resultAdd, destinationHandle)
}

// BigFloatSub VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatSub(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatSubName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSub
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)
	if err != nil {
		return err
	}

	resultSub, err := vmMath.SubBigFloat(op1, op2)
	if err != nil {
		return err
	}
	return setResultIfNotInfinity(context.GetVMHost(), resultSub, destinationHandle)
}

// BigFloatMul VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatMul(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatMulName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatMul
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if err != nil {
		return err
	}

	resultMul, err := vmMath.MulBigFloat(op1, op2)
	if err != nil {
		return err
	}
	return setResultIfNotInfinity(context.GetVMHost(), resultMul, destinationHandle)
}

// BigFloatDiv VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatDiv(destinationHandle, op1Handle, op2Handle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatDivName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatDiv
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if err != nil {
		return err
	}
	if areAllZero(op1, op2) {
		return vmhost.ErrAllOperandsAreEqualToZero
	}

	resultDiv, err := vmMath.QuoBigFloat(op1, op2)
	if err != nil {
		return err
	}
	return setResultIfNotInfinity(context.GetVMHost(), resultDiv, destinationHandle)
}

// BigFloatNeg VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatNeg(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatNegName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatNeg
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}
	op, err := managedType.GetBigFloat(opHandle)

	if err != nil {
		return err
	}
	dest.Neg(op)
	return nil
}

// BigFloatClone VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatClone(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatCloneName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatClone
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}
	op, err := managedType.GetBigFloat(opHandle)

	if err != nil {
		return err
	}
	dest.Copy(op)
	return nil
}

// BigFloatCmp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatCmp(op1Handle, op2Handle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatCmpName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCmp
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -2, err
	}

	op1, op2, err := managedType.GetTwoBigFloats(op1Handle, op2Handle)

	if err != nil {
		return -2, err
	}
	return int32(op1.Cmp(op2)), nil
}

// BigFloatAbs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatAbs(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatAbsName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}
	op, err := managedType.GetBigFloat(opHandle)

	if err != nil {
		return err
	}
	dest.Abs(op)
	return nil
}

// BigFloatSign VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatSign(opHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatAbs
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatSignName, gasToUse)
	if err != nil {
		return -2, err
	}

	op, err := managedType.GetBigFloat(opHandle)
	if err != nil {
		return -2, err
	}
	return int32(op.Sign()), nil
}

// BigFloatSqrt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatSqrt(destinationHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatSqrtName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSqrt
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}
	op, err := managedType.GetBigFloat(opHandle)

	if err != nil {
		return err
	}
	if op.Sign() < 0 {
		return vmhost.ErrBadLowerBounds
	}
	resultSqrt, err := vmMath.SqrtBigFloat(op)
	if err != nil {
		return err
	}
	dest.Set(resultSqrt)
	return nil
}

// BigFloatPow VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatPow(destinationHandle, opHandle, exponent int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatPowName)
	enableEpochsHandler := context.GetEnableEpochsHandler()

	if context.GetRuntimeContext().FixCorrectSmallIntAndGasLimits() && exponent < 0 {
		return vmhost.ErrExponentTooBigOrTooSmall
	}

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatPow
	if enableEpochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag) {
		gasToUse += vmMath.MulUint64(metering.GasSchedule().BigFloatAPICost.BigFloatPowPerIteration, uint64(exponent))
	}

	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	op, err := managedType.GetBigFloat(opHandle)
	if err != nil {
		return err
	}

	opBigInt := big.NewInt(0)
	op.Int(opBigInt)
	op2BigInt := big.NewInt(int64(exponent))
	if opBigInt.Sign() > 0 {
		opBigInt.Add(opBigInt, big.NewInt(1))
	}

	//this calculates the length of the result in bytes
	lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(op2BigInt, big.NewInt(int64(opBigInt.BitLen()))), big.NewInt(8))
	err = managedType.ConsumeGasForThisBigIntNumberOfBytes(lengthOfResult)
	if err != nil {
		return err
	}

	powResult, err := context.pow(op, exponent)
	if err != nil {
		return err
	}
	return setResultIfNotInfinity(context.GetVMHost(), powResult, destinationHandle)
}

func (context *VMHooksImpl) pow(base *big.Float, exp int32) (*big.Float, error) {
	result := big.NewFloat(1)
	result.SetPrec(base.Prec())
	managedType := context.GetManagedTypesContext()

	for i := 0; i < int(exp); i++ {
		resultMul, err := vmMath.MulBigFloat(result, base)
		if err != nil {
			return nil, err
		}
		exponent := resultMul.MantExp(nil)
		if managedType.BigFloatExpIsNotValid(exponent) {
			return nil, vmhost.ErrExponentTooBigOrTooSmall
		}
		result.Set(resultMul)
	}
	return result, nil
}

// BigFloatFloor VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatFloor(destBigIntHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatFloorName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatFloor
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	op, err := managedType.GetBigFloat(opHandle)
	if err != nil {
		return err
	}
	bigIntOp := managedType.GetBigIntOrCreate(destBigIntHandle)

	err = managedType.ConsumeGasForBigIntCopy(bigIntOp)
	if err != nil {
		return err
	}

	op.Int(bigIntOp)
	if op.IsInt() {
		return nil
	}
	if context.GetRuntimeContext().FixCorrectSmallIntAndGasLimits() {
		if op.Sign() < 0 {
			bigIntOp.Sub(bigIntOp, big.NewInt(1))
		}
	} else {
		if bigIntOp.Sign() < 0 {
			bigIntOp.Sub(bigIntOp, big.NewInt(1))
		}
	}
	return nil
}

// BigFloatCeil VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatCeil(destBigIntHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatCeilName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatCeil
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	op, err := managedType.GetBigFloat(opHandle)
	if err != nil {
		return err
	}
	bigIntOp := managedType.GetBigIntOrCreate(destBigIntHandle)

	err = managedType.ConsumeGasForBigIntCopy(bigIntOp)
	if err != nil {
		return err
	}

	op.Int(bigIntOp)
	if op.IsInt() {
		return nil
	}
	if context.GetRuntimeContext().FixCorrectSmallIntAndGasLimits() {
		if op.Sign() > 0 {
			bigIntOp.Add(bigIntOp, big.NewInt(1))
		}
	} else {
		if bigIntOp.Sign() > 0 {
			bigIntOp.Add(bigIntOp, big.NewInt(1))
		}
	}
	return nil
}

// BigFloatTruncate VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatTruncate(destBigIntHandle, opHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatTruncateName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatTruncate
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	op, err := managedType.GetBigFloat(opHandle)
	if err != nil {
		return err
	}
	bigIntValue := managedType.GetBigIntOrCreate(destBigIntHandle)

	err = managedType.ConsumeGasForBigIntCopy(bigIntValue)
	if err != nil {
		return err
	}

	op.Int(bigIntValue)
	return nil
}

// BigFloatSetInt64 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatSetInt64(destinationHandle int32, value int64) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetInt64
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatSetInt64Name, gasToUse)
	if err != nil {
		return err
	}

	dest, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}
	dest.SetInt64(value)
	return nil
}

// BigFloatIsInt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatIsInt(opHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatIsIntName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatIsInt
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	op, err := managedType.GetBigFloat(opHandle)
	if err != nil {
		return -1, err
	}
	if op.IsInt() {
		return 1, nil
	}
	return 0, nil
}

// BigFloatSetBigInt VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatSetBigInt(destinationHandle, bigIntHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(bigFloatSetBigIntName)

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatSetBigInt
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	bigIntValue, err := managedType.GetBigInt(bigIntHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBigIntCopy(bigIntValue)
	if err != nil {
		return err
	}

	resultSetInt := big.NewFloat(0).SetInt(bigIntValue)
	return setResultIfNotInfinity(context.GetVMHost(), resultSetInt, destinationHandle)
}

// BigFloatGetConstPi VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatGetConstPi(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatGetConst
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatGetConstPiName, gasToUse)
	if err != nil {
		return err
	}

	pi, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}
	pi.SetFloat64(math.Pi)
	return nil
}

// BigFloatGetConstE VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) BigFloatGetConstE(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigFloatAPICost.BigFloatGetConst
	err := metering.UseGasBoundedAndAddTracedGas(bigFloatGetConstEName, gasToUse)
	if err != nil {
		return err
	}

	e, err := managedType.GetBigFloatOrCreate(destinationHandle)
	if err != nil {
		return err
	}
	e.SetFloat64(math.E)
	return nil
}
