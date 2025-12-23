package factory

import (
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost/contexts"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/execution"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/executorFactory"
)

type runTypeComponentsFactory struct{}

// NewRunTypeComponentsFactory will return a new instance of runType components factory
func NewRunTypeComponentsFactory() *runTypeComponentsFactory {
	return &runTypeComponentsFactory{}
}

// Create will create the runType components
func (rtcf *runTypeComponentsFactory) Create() *RunTypeComponents {
	return &RunTypeComponents{
		BlockchainContextFactoryCreator: contexts.NewBlockchainContextFactory(),
		ExecutorFactoryCreator:          executorFactory.NewExecutorFactory(),
		RuntimeContextFactoryCreator:    contexts.NewRuntimeContextFactory(),
		MeteringContextFactoryCreator:   contexts.NewMeteringContextFactory(),
		OutputContextFactoryCreator:     contexts.NewOutputContextFactory(),
		ExecuteOnSameContextHandler:     execution.NewExecuteOnSameContextHandler(),
		CreateNewContractHandler:        execution.NewCreateNewContractHandler(),
		GasScheduleFactoryCreator:       config.NewGasScheduleFactory(),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (rtcf *runTypeComponentsFactory) IsInterfaceNil() bool {
	return rtcf == nil
}
