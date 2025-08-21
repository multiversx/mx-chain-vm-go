package vmhooks

import (
	"github.com/multiversx/mx-chain-crypto-go/zk/groth16"
	"github.com/multiversx/mx-chain-crypto-go/zk/plonk"
	"github.com/multiversx/mx-chain-crypto-go/zk/lowLevelFeatures"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const (
	managedVerifyGroth16  = "ManagedVerifyGroth16"
	managedVerifyPlonk    = "ManagedVerifyPlonk"
	managedAddEC          = "ManagedAddEC"
	managedMulEC          = "ManagedMulEC"
	managedMultiExpEC     = "ManagedMultiExpEC"
	managedMapToCurveEC   = "ManagedMapToCurveEC"
	managedPairingCheckEC = "ManagedPairingCheckEC"
)


// ManagedVerifyGroth16 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyGroth16(
	curveID int32, proofHandle, vkHandle, pubWitnessHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyZKFunctionWithHost(
		host, managedVerifyGroth16, curveID, proofHandle, vkHandle, pubWitnessHandle)
}

// ManagedVerifyPlonk VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyPlonk(
	curveID int32, proofHandle, vkHandle, pubWitnessHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyZKFunctionWithHost(
		host, managedVerifyPlonk, curveID, proofHandle, vkHandle, pubWitnessHandle)
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
	curveID int32,
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
	case managedVerifyGroth16:
		verified, invalidSigErr = groth16.VerifyGroth16(uint16(curveID), proofBytes, vkBytes, pubWitnessBytes)
	case managedVerifyPlonk:
		verified, invalidSigErr = plonk.VerifyPlonk(uint16(curveID), proofBytes, vkBytes, pubWitnessBytes)
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
	curveID int32,
	groupID int32,
	point1Handle, point2Handle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedAddECWithHost(host, curveID, groupID, point1Handle, point2Handle, resultHandle)
}

func managedECOperationWithHost(
	host vmhost.VMHost,
	operationName string,
	gasCost uint64,
	failureError error,
	curveID int32,
	groupID int32,
	inputHandles []int32,
	resultHandle int32,
	execute func(definedEC lowLevelFeatures.EllipticCurve, inputs [][]byte) ([]byte, error),
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas(operationName, gasCost)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	var inputsBytes [][]byte
	for _, handle := range inputHandles {
		bytes, err := getBytesAndConsumeGas(managedType, handle)
		if err != nil {
			FailExecution(host, err)
			return -1
		}
		inputsBytes = append(inputsBytes, bytes)
	}

	definedECParam := lowLevelFeatures.ECParams{Curve: lowLevelFeatures.ID(curveID), Group: lowLevelFeatures.GroupID(groupID)}
	definedEC, ok := lowLevelFeatures.EcRegistry[definedECParam]
	if !ok {
		FailExecutionConditionally(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// Gas cost is a placeholder. A more accurate gas cost would require changes to the gas schedule.
	result, err := execute(definedEC, inputsBytes)
	if err != nil {
		FailExecutionConditionally(host, failureError)
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

// ManagedAddECWithHost implements the Add elliptic curves operation on the set of defined curves and group
func ManagedAddECWithHost(
	host vmhost.VMHost,
	curveID int32,
	groupID int32,
	point1Handle, point2Handle int32,
	resultHandle int32,
) int32 {
	return managedECOperationWithHost(
		host,
		managedAddEC,
		host.Metering().GasSchedule().CryptoAPICost.AddECC,
		vmhost.ErrEllipticCurveAddFailed,
		curveID,
		groupID,
		[]int32{point1Handle, point2Handle},
		resultHandle,
		func(definedEC lowLevelFeatures.EllipticCurve, inputs [][]byte) ([]byte, error) {
			return definedEC.Add(inputs[0], inputs[1])
		},
	)
}

// ManagedMulEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMulEC(
	curveID int32,
	groupID int32,
	pointHandle, scalarHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMulECWithHost(host, curveID, groupID, pointHandle, scalarHandle, resultHandle)
}

// ManagedMulECWithHost implements the Multiply elliptic curves operation on the set of defined curves and group
func ManagedMulECWithHost(
	host vmhost.VMHost,
	curveID int32,
	groupID int32,
	pointHandle, scalarHandle int32,
	resultHandle int32,
) int32 {
	return managedECOperationWithHost(
		host,
		managedMulEC,
		host.Metering().GasSchedule().CryptoAPICost.ScalarMultECC,
		vmhost.ErrEllipticCurveMulFailed,
		curveID,
		groupID,
		[]int32{pointHandle, scalarHandle},
		resultHandle,
		func(definedEC lowLevelFeatures.EllipticCurve, inputs [][]byte) ([]byte, error) {
			return definedEC.Mul(inputs[0], inputs[1])
		},
	)
}

// ManagedMultiExpEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiExpEC(
	curveID int32,
	groupID int32,
	pointsHandle, scalarsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMultiExpECWithHost(host, curveID, groupID, pointsHandle, scalarsHandle, resultHandle)
}

func readManagedVectorsAndConsumeGas(
	host vmhost.VMHost,
	handle1, handle2 int32,
) ([][]byte, [][]byte, error) {
	managedType := host.ManagedTypes()
	metering := host.Metering()

	vec1, len1, err := managedType.ReadManagedVecOfManagedBuffers(handle1)
	if err != nil {
		return nil, nil, err
	}

	vec2, len2, err := managedType.ReadManagedVecOfManagedBuffers(handle2)
	if err != nil {
		return nil, nil, err
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, len1+len2)
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return nil, nil, err
	}

	return vec1, vec2, nil
}

// ManagedMultiExpECWithHost implements the MultiExp elliptic curves operation on the set of defined curves and group
func ManagedMultiExpECWithHost(
	host vmhost.VMHost,
	curveID int32,
	groupID int32,
	pointsHandle, scalarsHandle int32,
	resultHandle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas(managedMultiExpEC, metering.GasSchedule().CryptoAPICost.VerifyBLSMultiSig)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	pointsVec, scalarsVec, err := readManagedVectorsAndConsumeGas(host, pointsHandle, scalarsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	definedECParam := lowLevelFeatures.ECParams{Curve: lowLevelFeatures.ID(curveID), Group: lowLevelFeatures.GroupID(groupID)}
	definedEC, ok := lowLevelFeatures.EcRegistry[definedECParam]
	if !ok {
		FailExecutionConditionally(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// Gas cost is a placeholder. A more accurate gas cost would require changes to the gas schedule.
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

// ManagedMapToCurveEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMapToCurveEC(
	curveID int32,
	groupID int32,
	elementHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMapToCurveECWithHost(host, curveID, groupID, elementHandle, resultHandle)
}

// ManagedMapToCurveECWithHost implements the map to curve elliptic curves operation on the set of defined curves and group
func ManagedMapToCurveECWithHost(
	host vmhost.VMHost,
	curveID int32,
	groupID int32,
	elementHandle int32,
	resultHandle int32,
) int32 {
	return managedECOperationWithHost(
		host,
		managedMapToCurveEC,
		host.Metering().GasSchedule().CryptoAPICost.AddECC,
		vmhost.ErrEllipticCurveMapToCurveFailed,
		curveID,
		groupID,
		[]int32{elementHandle},
		resultHandle,
		func(definedEC lowLevelFeatures.EllipticCurve, inputs [][]byte) ([]byte, error) {
			return definedEC.MapToCurve(inputs[0])
		},
	)
}

// ManagedPairingCheckEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedPairingCheckEC(
	curveID int32,
	pointsG1Handle, pointsG2Handle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedPairingCheckECWithHost(host, curveID, pointsG1Handle, pointsG2Handle)
}

// ManagedPairingCheckECWithHost implements the pairing checks elliptic curves operation on the set of defined curves
func ManagedPairingCheckECWithHost(
	host vmhost.VMHost,
	curveID int32,
	pointsG1Handle, pointsG2Handle int32,
) int32 {
	metering := host.Metering()

	err := metering.UseGasBoundedAndAddTracedGas(managedPairingCheckEC, metering.GasSchedule().CryptoAPICost.VerifyBLSMultiSig)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	pointsG1Vec, pointsG2Vec, err := readManagedVectorsAndConsumeGas(host, pointsG1Handle, pointsG2Handle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	definedPairingRegistry, ok := lowLevelFeatures.PairingRegistry[lowLevelFeatures.ID(curveID)]
	if !ok {
		FailExecutionConditionally(host, vmhost.ErrNoEllipticCurveUnderThisHandle)
		return -1
	}

	// Gas cost is a placeholder. A more accurate gas cost would require changes to the gas schedule.
	verified, err := definedPairingRegistry.PairingCheck(pointsG1Vec, pointsG2Vec)
	if err != nil || !verified {
		FailExecutionConditionally(host, vmhost.ErrEllipticCurvePairingCheckFailed)
		return -1
	}

	return 0
}
