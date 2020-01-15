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
		gasSchedule:   gasCostConfig,
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

func (context *meteringContext) BlockGasLimit() uint64 {
	return context.blockGasLimit
}

func (context *meteringContext) DeductInitialGasForExecution(input *vmcommon.ContractCallInput, contract []byte) (uint64, error) {
	remainingGas, err := context.deductInitialGas(
		input.GasProvided,
		contract,
		0,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
	return remainingGas, err
}

func (context *meteringContext) DeductInitialGasForDirectDeployment(input *vmcommon.ContractCreateInput) (uint64, error) {
	remainingGas, err := context.deductInitialGas(
		input.GasProvided,
		input.ContractCode,
		context.gasSchedule.ElrondAPICost.CreateContract,
		context.gasSchedule.BaseOperationCost.StorePerByte,
	)
	return remainingGas, err
}

func (context *meteringContext) DeductInitialGasForIndirectDeployment(input *vmcommon.ContractCreateInput) (uint64, error) {
	remainingGas, err := context.deductInitialGas(
		input.GasProvided,
		input.ContractCode,
		0,
		context.gasSchedule.BaseOperationCost.StorePerByte,
	)
	return remainingGas, err
}

func (context *meteringContext) deductInitialGas(
	gasProvided uint64,
	code []byte,
	baseCost uint64,
	costPerByte uint64,
) (uint64, error) {
	codeLength := uint64(len(code))
	codeCost := codeLength * costPerByte
	initialCost := baseCost + codeCost

	if initialCost > gasProvided {
		return 0, arwen.ErrNotEnoughGas
	}

	return gasProvided - initialCost, nil
}
