package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type meteringContext struct {
	gasSchedule           *config.GasCost
	blockGasLimit         uint64
	gasLockedForAsyncStep bool
	host                  arwen.VMHost
}

func NewMeteringContext(
	host arwen.VMHost,
	gasSchedule map[string]map[string]uint64,
	blockGasLimit uint64,
) (*meteringContext, error) {

	gasCostConfig, err := config.CreateGasConfig(gasSchedule)
	if err != nil {
		return nil, err
	}

	context := &meteringContext{
		gasSchedule:           gasCostConfig,
		blockGasLimit:         blockGasLimit,
		gasLockedForAsyncStep: false,
		host:                  host,
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

func (context *meteringContext) FreeGas(gas uint64) {
	refund := context.host.Output().GetRefund() + gas
	context.host.Output().SetRefund(refund)
}

func (context *meteringContext) GasLeft() uint64 {
	gasProvided := context.host.Runtime().GetVMInput().GasProvided
	gasUsed := context.host.Runtime().GetPointsUsed()
	return gasProvided - gasUsed
}

func (context *meteringContext) BoundGasLimit(value int64) uint64 {
	gasLeft := context.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	} else {
		return limit
	}
}

func (context *meteringContext) DeductAndLockGasIfAsyncStep() error {
	input := context.host.Runtime().GetVMInput()

	if input.CallType == vmcommon.AsynchronousCall {
		gasSchedule := context.GasSchedule().ElrondAPICost

		input.GasProvided -= gasSchedule.AsyncCallStep
		if input.GasProvided <= 0 {
			return arwen.ErrNotEnoughGas
		}

		input.GasProvided -= gasSchedule.AsyncCallStep + gasSchedule.AsyncCallbackGasLock
		if input.GasProvided <= 0 {
			return arwen.ErrNotEnoughGas
		}

		context.gasLockedForAsyncStep = true
	} else {
		context.gasLockedForAsyncStep = false
	}

	return nil
}

func (context *meteringContext) UnlockGasIfAsyncStep() {
	if context.gasLockedForAsyncStep {
		gasSchedule := context.GasSchedule().ElrondAPICost
		input := context.host.Runtime().GetVMInput()
		input.GasProvided += gasSchedule.AsyncCallStep + gasSchedule.AsyncCallbackGasLock
		context.gasLockedForAsyncStep = false
	}
}

func (context *meteringContext) BlockGasLimit() uint64 {
	return context.blockGasLimit
}

func (context *meteringContext) DeductInitialGasForExecution(contract []byte) error {
	err := context.deductInitialGas(
		contract,
		0,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
	return err
}

func (context *meteringContext) DeductInitialGasForDirectDeployment(input *vmcommon.ContractCreateInput) error {
	err := context.deductInitialGas(
		input.ContractCode,
		context.gasSchedule.ElrondAPICost.CreateContract,
		context.gasSchedule.BaseOperationCost.StorePerByte,
	)
	return err
}

func (context *meteringContext) DeductInitialGasForIndirectDeployment(input *vmcommon.ContractCreateInput) error {
	err := context.deductInitialGas(
		input.ContractCode,
		0,
		context.gasSchedule.BaseOperationCost.StorePerByte,
	)
	return err
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
