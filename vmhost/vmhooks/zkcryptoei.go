package vmhooks

import (
	"github.com/multiversx/mx-chain-crypto-go/zk/groth16"
	"github.com/multiversx/mx-chain-crypto-go/zk/lowLevelFeatures"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const (
	ManagedVerifyGroth16 = "ManagedVerifyGroth16"
	ManagedVerifyPlonk   = "ManagedVerifyPlonk"
)

/*
	BN254
	BLS12_377
	BLS12_381
	BLS24_315
	BLS24_317
	BW6_761
	BW6_633
	STARK_CURVE
	SECP256K1
	GRUMPKIN
*/

// ManagedVerifyGroth16 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyGroth16(
	curveID uint16, proofHandle, vkHandle, pubWitnessHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyZKFunctionWithHost(
		host, ManagedVerifyGroth16, curveID, proofHandle, vkHandle, pubWitnessHandle)
}

// ManagedVerifyPlonk VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyPlonk(
	curveID uint16, proofHandle, vkHandle, pubWitnessHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyZKFunctionWithHost(
		host, ManagedVerifyPlonk, curveID, proofHandle, vkHandle, pubWitnessHandle)
}

func getBytesAndConsumeGas(managedType vmhost.ManagedTypesContext, handle int32) ([]byte, error) {
	bytesVec, err := managedType.GetBytes(handle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBytes(bytesVec)
	if err != nil {
		return nil, err
	}

	return bytesVec, nil
}

// ManagedVerifyZKFunctionWithHost VMHooks implementation with host
func ManagedVerifyZKFunctionWithHost(
	host vmhost.VMHost,
	zkFunc string,
	curveID uint16,
	proofHandle, vkHandle, pubWitnessHandle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas(zkFunc, metering.GasSchedule().CryptoAPICost.VerifyBLSMultiSig)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	proofBytes, err := getBytesAndConsumeGas(managedType, proofHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	vkBytes, err := getBytesAndConsumeGas(managedType, vkHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	pubWitnessBytes, err := getBytesAndConsumeGas(managedType, pubWitnessHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	verified := false
	invalidSigErr := vmhost.ErrInvalidArgument
	switch zkFunc {
	case ManagedVerifyGroth16:
		verified, invalidSigErr = groth16.VerifyGroth16(curveID, proofBytes, vkBytes, pubWitnessBytes)
	case ManagedVerifyPlonk:
		verified, invalidSigErr = groth16.VerifyGroth16(curveID, proofBytes, vkBytes, pubWitnessBytes)
	}

	if invalidSigErr != nil || !verified {
		FailExecutionConditionally(host, vmhost.ErrZKVerify)
		return -1
	}

	return 0
}

// ManagedAddEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedAddEC(
	curveID uint16,
	groupID uint16,
	point1Handle, point2Handle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedAddECWithHost(host, curveID, groupID, point1Handle, point2Handle, resultHandle)
}

// ManagedAddECWithHost implements the Add elliptic curves operation on the set of defined curves and group
func ManagedAddECWithHost(
	host vmhost.VMHost,
	curveID uint16,
	groupID uint16,
	point1Handle, point2Handle int32,
	resultHandle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas(addECName, metering.GasSchedule().CryptoAPICost.AddECC)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	point1Bytes, err := getBytesAndConsumeGas(managedType, point1Handle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	point2Bytes, err := getBytesAndConsumeGas(managedType, point2Handle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	definedECParam := lowLevelFeatures.ECParams{Curve: lowLevelFeatures.ID(curveID), Group: lowLevelFeatures.GroupID(groupID)}
	definedEC, ok := lowLevelFeatures.EcRegistry[definedECParam]
	if !ok {
		FailExecution(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// TODO: use more gas depending on type

	result, err := definedEC.Add(point1Bytes, point2Bytes)
	if err != nil {
		FailExecutionConditionally(host, vmhost.ErrEllipticCurveAddFailed)
		return -1
	}

	err = managedType.ConsumeGasForBytes(result)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	managedType.SetBytes(resultHandle, result)

	return 0
}

// ManagedMulEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMulEC(
	curveID uint16,
	groupID uint16,
	pointHandle, scalarHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMulECWithHost(host, curveID, groupID, pointHandle, scalarHandle, resultHandle)
}

// ManagedMulECWithHost implements the Multiply elliptic curves operation on the set of defined curves and group
func ManagedMulECWithHost(
	host vmhost.VMHost,
	curveID uint16,
	groupID uint16,
	pointHandle, scalarHandle int32,
	resultHandle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas("multiplyEC", metering.GasSchedule().CryptoAPICost.AddECC)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	pointBytes, err := getBytesAndConsumeGas(managedType, pointHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	scalarBytes, err := getBytesAndConsumeGas(managedType, scalarHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	definedECParam := lowLevelFeatures.ECParams{Curve: lowLevelFeatures.ID(curveID), Group: lowLevelFeatures.GroupID(groupID)}
	definedEC, ok := lowLevelFeatures.EcRegistry[definedECParam]
	if !ok {
		FailExecution(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// TODO: use more gas depending on scalar and curve type

	result, err := definedEC.Mul(pointBytes, scalarBytes)
	if err != nil {
		FailExecutionConditionally(host, vmhost.ErrEllipticCurveMulFailed)
		return -1
	}

	err = managedType.ConsumeGasForBytes(result)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	managedType.SetBytes(resultHandle, result)

	return 0
}

// ManagedMultiExpEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiExpEC(
	curveID uint16,
	groupID uint16,
	pointsHandle, scalarsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMultiExpECWithHost(host, curveID, groupID, pointsHandle, scalarsHandle, resultHandle)
}

// ManagedMultiExpECWithHost implements the MultiExp elliptic curves operation on the set of defined curves and group
func ManagedMultiExpECWithHost(
	host vmhost.VMHost,
	curveID uint16,
	groupID uint16,
	pointsHandle, scalarsHandle int32,
	resultHandle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas("multiplyEC", metering.GasSchedule().CryptoAPICost.AddECC)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	pointsVec, actualLenPoints, err := managedType.ReadManagedVecOfManagedBuffers(pointsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	scalarsVec, actualLenScalars, err := managedType.ReadManagedVecOfManagedBuffers(scalarsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLenPoints+actualLenScalars)
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	definedECParam := lowLevelFeatures.ECParams{Curve: lowLevelFeatures.ID(curveID), Group: lowLevelFeatures.GroupID(groupID)}
	definedEC, ok := lowLevelFeatures.EcRegistry[definedECParam]
	if !ok {
		FailExecution(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// TODO: use more gas depending on scalar and curve type

	result, err := definedEC.MultiExp(pointsVec, scalarsVec)
	if err != nil {
		FailExecutionConditionally(host, vmhost.ErrEllipticCurveMultiExpFailed)
		return -1
	}

	err = managedType.ConsumeGasForBytes(result)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	managedType.SetBytes(resultHandle, result)

	return 0
}

// ManagedMultiExpEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiExpEC(
	curveID uint16,
	groupID uint16,
	elementHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMapToCurveWithHost(host, curveID, groupID, elementHandle, resultHandle)
}

// ManagedMapToCurveWithHost implements the map to curve elliptic curves operation on the set of defined curves and group
func ManagedMapToCurveWithHost(
	host vmhost.VMHost,
	curveID uint16,
	groupID uint16,
	elementHandle int32,
	resultHandle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas("multiplyEC", metering.GasSchedule().CryptoAPICost.AddECC)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	element, err := getBytesAndConsumeGas(managedType, elementHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	definedECParam := lowLevelFeatures.ECParams{Curve: lowLevelFeatures.ID(curveID), Group: lowLevelFeatures.GroupID(groupID)}
	definedEC, ok := lowLevelFeatures.EcRegistry[definedECParam]
	if !ok {
		FailExecution(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// TODO: use more gas depending on scalar and curve type

	result, err := definedEC.MapToCurve(element)
	if err != nil {
		FailExecutionConditionally(host, vmhost.ErrEllipticCurveMapToCurveFailed)
		return -1
	}

	err = managedType.ConsumeGasForBytes(result)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	managedType.SetBytes(resultHandle, result)

	return 0
}

// ManagedPairingChecksEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedPairingChecksEC(
	curveID uint16,
	pointsG1Handle, pointsG2Handle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedPairingChecksECWithHost(host, curveID, pointsG1Handle, pointsG2Handle)
}

// ManagedPairingChecksECWithHost implements the pairing checks elliptic curves operation on the set of defined curves
func ManagedPairingChecksECWithHost(
	host vmhost.VMHost,
	curveID uint16,
	pointsG1Handle, pointsG2Handle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas("multiplyEC", metering.GasSchedule().CryptoAPICost.AddECC)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	pointsG1Vec, actualLenPoints, err := managedType.ReadManagedVecOfManagedBuffers(pointsG1Handle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	pointsG2Vec, actualLenScalars, err := managedType.ReadManagedVecOfManagedBuffers(pointsG2Handle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLenPoints+actualLenScalars)
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	definedPairingRegistry, ok := lowLevelFeatures.PairingRegistry[lowLevelFeatures.ID(curveID)]
	if !ok {
		FailExecution(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// TODO: use more gas depending on scalar and curve type

	verified, err := definedPairingRegistry.PairingCheck(pointsG1Vec, pointsG2Vec)
	if err != nil || !verified {
		FailExecutionConditionally(host, vmhost.ErrEllipticCurvePairingCheckFailed)
		return -1
	}

	return 0
}
