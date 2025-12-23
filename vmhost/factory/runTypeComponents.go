package factory

import (
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost/contexts"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/execution"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/executorFactory"
)

// RunTypeComponents struct holds the components needed for run type
type RunTypeComponents struct {
	BlockchainContextFactoryCreator contexts.BlockchainContextCreator
	ExecutorFactoryCreator          executorFactory.ExecutorCreator
	RuntimeContextFactoryCreator    contexts.RuntimeContextCreator
	MeteringContextFactoryCreator   contexts.MeteringContextCreator
	OutputContextFactoryCreator     contexts.OutputContextCreator
	ExecuteOnSameContextHandler     execution.ExecuteOnSameContextHandler
	CreateNewContractHandler        execution.CreateNewContractHandler
	GasScheduleFactoryCreator       config.GasScheduleFactory
}

// Close does nothing
func (rtc *RunTypeComponents) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rtc *RunTypeComponents) IsInterfaceNil() bool {
	return rtc == nil
}
