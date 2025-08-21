package vmhooks

import (
	"github.com/multiversx/mx-chain-crypto-go/zk/groth16"
	"github.com/multiversx/mx-chain-crypto-go/zk/lowLevelFeatures"
	"github.com/multiversx/mx-chain-crypto-go/zk/plonk"
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

func managedECOperationWithHost(
	host vmhost.VMHost,
	operationName string,
	gasCost uint64,
	failureError error,
	curveID int32,
	groupID int32,
	inputHandles []int32,
	resultHandle int32,
	execute func(definedEC lowLevelFeatures.ECGroup, inputs [][]byte) ([]byte, error),
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

	// TODO: use more gas depending on scalar and curve type. This would require changes to the gas schedule and the VM's core logic.
	// The gas cost should be dependent on the curve type and other parameters. To implement this, we would need to:
	// 1. Define new gas cost parameters in the `CryptoAPICost` struct for each curve type and operation.
	// 2. Update the `FillGasMapCryptoAPICosts` function to initialize these new parameters.
	// 3. Update this function to use the new gas cost parameters based on the curve ID.

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

func addEC(definedEC lowLevelFeatures.ECGroup, inputs [][]byte) ([]byte, error) {
	if len(inputs) != 2 {
		return nil, vmhost.ErrArgIndexOutOfRange
	}
	return definedEC.Add(inputs[0], inputs[1])
}

func mulEC(definedEC lowLevelFeatures.ECGroup, inputs [][]byte) ([]byte, error) {
	if len(inputs) != 2 {
		return nil, vmhost.ErrArgIndexOutOfRange
	}
	return definedEC.Mul(inputs[0], inputs[1])
}

func mapToCurveEC(definedEC lowLevelFeatures.ECGroup, inputs [][]byte) ([]byte, error) {
	if len(inputs) != 1 {
		return nil, vmhost.ErrArgIndexOutOfRange
	}
	return definedEC.MapToCurve(inputs[0])
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
	return managedECOperationWithHost(
		host,
		managedAddEC,
		host.Metering().GasSchedule().CryptoAPICost.AddECC,
		vmhost.ErrEllipticCurveAddFailed,
		curveID,
		groupID,
		[]int32{point1Handle, point2Handle},
		resultHandle,
		addEC,
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
	return managedECOperationWithHost(
		host,
		managedMulEC,
		host.Metering().GasSchedule().CryptoAPICost.ScalarMultECC,
		vmhost.ErrEllipticCurveMulFailed,
		curveID,
		groupID,
		[]int32{pointHandle, scalarHandle},
		resultHandle,
		mulEC,
	)
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
	return managedECOperationWithHost(
		host,
		managedMapToCurveEC,
		host.Metering().GasSchedule().CryptoAPICost.AddECC,
		vmhost.ErrEllipticCurveMapToCurveFailed,
		curveID,
		groupID,
		[]int32{elementHandle},
		resultHandle,
		mapToCurveEC,
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

	err := metering.UseGasBoundedAndAddTracedGas(managedMultiExpEC, metering.GasSchedule().CryptoAPICost.AddECC)
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

// ManagedPairingChecksEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedPairingChecksEC(
	curveID int32,
	pointsG1Handle, pointsG2Handle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedPairingChecksECWithHost(host, curveID, pointsG1Handle, pointsG2Handle)
}

// ManagedPairingChecksECWithHost implements the pairing checks elliptic curves operation on the set of defined curves
func ManagedPairingChecksECWithHost(
	host vmhost.VMHost,
	curveID int32,
	pointsG1Handle, pointsG2Handle int32,
) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	err := metering.UseGasBoundedAndAddTracedGas(managedPairingCheckEC, metering.GasSchedule().CryptoAPICost.AddECC)
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
