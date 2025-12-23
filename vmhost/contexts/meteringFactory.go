package contexts

import (
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type meteringContextFactory struct{}

// NewMeteringContextFactory creates a new metering context factory
func NewMeteringContextFactory() *meteringContextFactory {
	return &meteringContextFactory{}
}

// CreateMeteringContext create the metering context
func (mcf *meteringContextFactory) CreateMeteringContext(host vmhost.VMHost, gasMap config.GasScheduleMap, blockGasLimit uint64, gasScheduleFactory config.GasScheduleFactory) (vmhost.MeteringContext, error) {
	return NewMeteringContext(host, gasMap, blockGasLimit, gasScheduleFactory)
}

// IsInterfaceNil returns true if underlying object is nil
func (mcf *meteringContextFactory) IsInterfaceNil() bool {
	return mcf == nil
}
