package subcontexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

type Metering struct {
	gasCostConfig *config.GasCost
}

func (metering *Metering) GasSchedule() *config.GasCost {
	return metering.gasCostConfig
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
