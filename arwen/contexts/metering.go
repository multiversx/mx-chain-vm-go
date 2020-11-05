package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
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

func (context *meteringContext) GasSchedule() *config.GasCost {
	return context.gasSchedule
}

func (context *meteringContext) UseGas(gas uint64) {
	gasUsed := context.host.Runtime().GetPointsUsed() + gas
	context.host.Runtime().SetPointsUsed(gasUsed)
}

func (context *meteringContext) RestoreGas(gas uint64) {
	gasUsed := context.host.Runtime().GetPointsUsed()
	if gas <= gasUsed {
		gasUsed -= gas
		context.host.Runtime().SetPointsUsed(gasUsed)
	}
}

func (context *meteringContext) FreeGas(gas uint64) {
	refund := context.host.Output().GetRefund() + gas
	context.host.Output().SetRefund(refund)
}

func (context *meteringContext) GasLeft() uint64 {
	gasProvided := context.host.Runtime().GetVMInput().GasProvided
	gasUsed := context.host.Runtime().GetPointsUsed()

	if gasProvided < gasUsed {
		return 0
	}

	return gasProvided - gasUsed
}

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
// TODO the warmInstance will be affected if this function is misused
func (context *meteringContext) UseGasForAsyncStep() error {
	gasSchedule := context.GasSchedule().ElrondAPICost
	gasToDeduct := gasSchedule.AsyncCallStep
	return context.UseGasBounded(gasToDeduct)
}

func (context *meteringContext) UseGasBounded(gasToUse uint64) error {
	if context.GasLeft() <= gasToUse {
		return arwen.ErrNotEnoughGas
	}
	context.UseGas(gasToUse)
	return nil
}

// ComputeGasToLockForAsync calculates the minimum amount of gas to lock for async callbacks
func (context *meteringContext) ComputeGasLockedForAsync() uint64 {
	baseGasSchedule := context.GasSchedule().BaseOperationCost
	apiGasSchedule := context.GasSchedule().ElrondAPICost
	codeSize := context.host.Runtime().GetSCCodeSize()

	// Exact amount of gas required to compile this SC again, to execute the
	// callback
	compilationGasLock := codeSize * baseGasSchedule.CompilePerByte

	// Minimum amount required to execute the callback
	executionGasLock := apiGasSchedule.AsyncCallStep + apiGasSchedule.AsyncCallbackGasLock

	return compilationGasLock + executionGasLock
}

func (context *meteringContext) UnlockGasIfAsyncCallback() {
	input := context.host.Runtime().GetVMInput()
	if input.CallType != vmcommon.AsynchronousCallBack {
		return
	}

	input.GasProvided += input.GasLocked
	input.GasLocked = 0
}

func (context *meteringContext) GetGasLocked() uint64 {
	input := context.host.Runtime().GetVMInput()
	return input.GasLocked
}

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
	codeCost := codeLength * costPerByte
	initialCost := baseCost + codeCost

	if initialCost > input.GasProvided {
		return arwen.ErrNotEnoughGas
	}

	input.GasProvided -= initialCost
	return nil
}
