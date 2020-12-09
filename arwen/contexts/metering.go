package contexts

import (
	builtinMath "math"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type meteringContext struct {
	gasSchedule   *config.GasCost
	blockGasLimit uint64
	host          arwen.VMHost
}

// NewMeteringContext creates a new meteringContext
func NewMeteringContext(
	host arwen.VMHost,
	gasMap config.GasScheduleMap,
	blockGasLimit uint64,
) (*meteringContext, error) {

	gasSchedule, err := config.CreateGasConfig(gasMap)
	if err != nil {
		return nil, err
	}

	context := &meteringContext{
		gasSchedule:   gasSchedule,
		blockGasLimit: blockGasLimit,
		host:          host,
	}

	return context, nil
}

// GasSchedule returns the current gas schedule
func (context *meteringContext) GasSchedule() *config.GasCost {
	return context.gasSchedule
}

// SetGasSchedule sets the gas schedule to the given gas map
func (context *meteringContext) SetGasSchedule(gasMap config.GasScheduleMap) {
	gasSchedule, err := config.CreateGasConfig(gasMap)
	if err != nil {
		log.Error("SetGasSchedule createGasConfig", "error", err)
		return
	}
	context.gasSchedule = gasSchedule
}

// UseGas sets in the runtime context the given gas as gas used
func (context *meteringContext) UseGas(gas uint64) {
	gasUsed, err := math.AddUint64(context.host.Runtime().GetPointsUsed(), gas)
	if err != nil || gasUsed > builtinMath.MaxUint64 {
		log.Error("UseGas overflow",
			"gasUsed = ", context.host.Runtime().GetPointsUsed(),
			"gasToUse = ", gas,
		)
		context.host.Runtime().SetPointsUsed(builtinMath.MaxUint64)
		return
	}

	context.host.Runtime().SetPointsUsed(gasUsed)
}

// RestoreGas subtracts the given gas from the gas used that is set in the runtime context.
func (context *meteringContext) RestoreGas(gas uint64) {
	gasUsed := context.host.Runtime().GetPointsUsed()
	if gas <= gasUsed {
		gasUsed -= gas
		context.host.Runtime().SetPointsUsed(gasUsed)
	}
}

// FreeGas adds the given gas to the refunded gas.
func (context *meteringContext) FreeGas(gas uint64) {
	refund, err := math.AddUint64(context.host.Output().GetRefund(), gas)
	if err != nil || refund > builtinMath.MaxUint64 {
		log.Error("FreeGas overflow",
			"gasUsed = ", context.host.Runtime().GetPointsUsed(),
			"gasToUse = ", gas,
		)
		context.host.Runtime().SetPointsUsed(builtinMath.MaxUint64)
		return
	}

	context.host.Output().SetRefund(refund)
}

// GasLeft returns how much gas is left.
func (context *meteringContext) GasLeft() uint64 {
	gasProvided := context.host.Runtime().GetVMInput().GasProvided
	gasUsed := context.host.Runtime().GetPointsUsed()

	if gasProvided < gasUsed {
		return 0
	}

	return gasProvided - gasUsed
}

// BoundGasLimit returns the gas left if it is less than the given limit, or the given value otherwise.
func (context *meteringContext) BoundGasLimit(value int64) uint64 {
	gasLeft := context.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

// UseGasForAsyncStep consumes the AsyncCallStep gas cost on the currently
// running Wasmer instance
func (context *meteringContext) UseGasForAsyncStep() error {
	gasSchedule := context.GasSchedule().ElrondAPICost
	gasToDeduct := gasSchedule.AsyncCallStep
	return context.UseGasBounded(gasToDeduct)
}

// UseGasBounded returns an error if the given gasToUse is less than the available gas,
// otherwise it uses the given gas
func (context *meteringContext) UseGasBounded(gasToUse uint64) error {
	if context.GasLeft() <= gasToUse {
		return arwen.ErrNotEnoughGas
	}
	context.UseGas(gasToUse)
	return nil
}

// ComputeGasLockedForAsync calculates the minimum amount of gas to lock for async callbacks
func (context *meteringContext) ComputeGasLockedForAsync() uint64 {
	baseGasSchedule := context.GasSchedule().BaseOperationCost
	apiGasSchedule := context.GasSchedule().ElrondAPICost
	codeSize := context.host.Runtime().GetSCCodeSize()

	costPerByte := baseGasSchedule.CompilePerByte
	if context.host.IsAheadOfTimeCompileEnabled() {
		costPerByte = baseGasSchedule.AoTPreparePerByte
	}

	// Exact amount of gas required to compile this SC again, to execute the callback
	compilationGasLock := uint64(0)
	if context.host.IsDynamicGasLockingEnabled() {
		var err error
		compilationGasLock, err = math.MulUint64(codeSize, costPerByte)
		if err != nil {
			log.Error("ComputeGasLockedForAsync overflow",
				"codeSize = ", codeSize,
				"costPerByte = ", costPerByte,
			)

			return builtinMath.MaxUint64
		}
	}

	// Minimum amount required to execute the callback
	executionGasLock, err := math.AddUint64(apiGasSchedule.AsyncCallStep, apiGasSchedule.AsyncCallbackGasLock)
	if err != nil {
		log.Error("ComputeGasLockedForAsync overflow",
			"AsyncCallStep = ", apiGasSchedule.AsyncCallStep,
			"AsyncCallbackGasLock = ", apiGasSchedule.AsyncCallbackGasLock,
		)

		return builtinMath.MaxUint64
	}

	gasLockedForAsync, err := math.AddUint64(compilationGasLock, executionGasLock)
	if err != nil {
		log.Error("ComputeGasLockedForAsync overflow",
			"compilationGasLock = ", compilationGasLock,
			"executionGasLock = ", executionGasLock,
		)

		return builtinMath.MaxUint64
	}

	return gasLockedForAsync
}

// UnlockGasIfAsyncCallback unlocks the locked gas if the call type is async callback
func (context *meteringContext) UnlockGasIfAsyncCallback() {
	input := context.host.Runtime().GetVMInput()
	if input.CallType != vmcommon.AsynchronousCallBack {
		return
	}

	gasProvided, err := math.AddUint64(input.GasProvided, input.GasLocked)
	if err != nil {
		log.Error("UnlockGasIfAsyncCallback overflow",
			"GasProvided = ", input.GasProvided,
			"GasLocked = ", input.GasLocked,
		)

		gasProvided = builtinMath.MaxUint64
	}

	input.GasProvided = gasProvided
	input.GasLocked = 0
}

// GetGasLocked returns the locked gas
func (context *meteringContext) GetGasLocked() uint64 {
	input := context.host.Runtime().GetVMInput()
	return input.GasLocked
}

// BlockGasLimit returns the gas limit for the current block
func (context *meteringContext) BlockGasLimit() uint64 {
	return context.blockGasLimit
}

// DeductInitialGasForExecution deducts gas for compilation and locks gas if the execution is an asynchronous call
func (context *meteringContext) DeductInitialGasForExecution(contract []byte) error {
	costPerByte := context.gasSchedule.BaseOperationCost.CompilePerByte
	if context.host.IsAheadOfTimeCompileEnabled() {
		costPerByte = context.gasSchedule.BaseOperationCost.AoTPreparePerByte
	}
	err := context.deductInitialGas(contract, 0, costPerByte)
	if err != nil {
		return err
	}

	return nil
}

// DeductInitialGasForDirectDeployment deducts gas for the deployment of a contract initiated by a Transaction
func (context *meteringContext) DeductInitialGasForDirectDeployment(input arwen.CodeDeployInput) error {
	return context.deductInitialGas(
		input.ContractCode,
		context.gasSchedule.ElrondAPICost.CreateContract,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
}

// DeductInitialGasForIndirectDeployment deducts gas for the deployment of a contract initiated by another SmartContract
func (context *meteringContext) DeductInitialGasForIndirectDeployment(input arwen.CodeDeployInput) error {
	return context.deductInitialGas(
		input.ContractCode,
		0,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
}

func (context *meteringContext) deductInitialGas(
	code []byte,
	baseCost uint64,
	costPerByte uint64,
) error {
	input := context.host.Runtime().GetVMInput()
	codeLength := uint64(len(code))
	codeCost, err := math.MulUint64(codeLength, costPerByte)
	if err != nil {
		log.Error("deductInitialGas overflow",
			"codeLength = ", codeLength,
			"costPerByte = ", costPerByte,
		)

		return arwen.ErrNotEnoughGas
	}

	initialCost, err := math.AddUint64(baseCost, codeCost)
	if err != nil {
		log.Error("deductInitialGas overflow",
			"baseCost = ", baseCost,
			"codeCost = ", codeCost,
		)

		return arwen.ErrNotEnoughGas
	}

	if initialCost > input.GasProvided {
		return arwen.ErrNotEnoughGas
	}

	input.GasProvided -= initialCost
	return nil
}
