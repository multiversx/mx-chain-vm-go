package subcontexts

import (
	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

type Metering struct {
	gasSchedule       *config.GasCost
	blockGasLimit     uint64
	host arwen.VMContext
}

func NewMeteringSubcontext(
	gasSchedule map[string]map[string]uint64,
	blockGasLimit uint64,
	host arwen.VMContext,
) (*Metering, error) {

	gasCostConfig, err := config.CreateGasConfig(gasSchedule)
	if err != nil {
		return nil, err
	}

	metering := &Metering{
		gasSchedule:       gasCostConfig,
		blockGasLimit:     blockGasLimit,
		host: host,
	}

	return metering, nil
}

func (metering *Metering) GasSchedule() *config.GasCost {
	return metering.gasSchedule
}

func (metering *Metering) UseGas(gas uint64) {
	gasUsed := metering.host.Runtime().GetPointsUsed() + gas
	metering.host.Runtime().SetPointsUsed(gasUsed)
}

func (metering *Metering) FreeGas(gas uint64) {
	refund := metering.host.Output().GetRefund() + gas
	metering.host.Output().SetRefund(refund)
}

func (metering *Metering) GasLeft() uint64 {
	gasProvided := metering.host.Runtime().GetVMInput().GasProvided
	gasUsed := metering.host.Runtime().GetPointsUsed()
	return gasProvided - gasUsed
}

func (metering *Metering) BoundGasLimit(value int64) uint64 {
	gasLeft := metering.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	} else {
		return limit
	}
}

func (metering *Metering) BlockGasLimit() uint64 {
	return metering.blockGasLimit
}
