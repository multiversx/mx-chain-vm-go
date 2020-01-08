package subcontexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

type Metering struct {
	gasSchedule *config.GasCost
	blockGasLimit uint64
}

func NewMeteringSubcontext(
	gasSchedule map[string]map[string]uint64,
	blockGasLimit uint64,
) (*Metering, error) {

	gasCostConfig, err := config.CreateGasConfig(gasSchedule)
	if err != nil {
		return nil, err
	}

	metering := &Metering{
		gasSchedule: gasCostConfig,
		blockGasLimit: blockGasLimit,
	}

	return metering, nil
}

func (metering *Metering) GasSchedule() *config.GasCost {
	return metering.gasSchedule
}

func (metering *Metering) UseGas(gas uint64) {
	panic("not implemented")
}

func (metering *Metering) FreeGas(gas uint64) {
	panic("not implemented")
}

func (metering *Metering) GasLeft() uint64 {
	panic("not implemented")
}

func (metering *Metering) BoundGasLimit(value int64) uint64 {
	panic("not implemented")
}

func (metering *Metering) BlockGasLimit() uint64 {
	panic("not implemented")
}
