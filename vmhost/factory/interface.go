package factory

import (
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost/contexts"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/execution"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/executorFactory"
)

type runTypeComponentsCreator interface {
	Create() *RunTypeComponents
	IsInterfaceNil() bool
}

// ComponentHandler defines the actions common to all component handlers
type ComponentHandler interface {
	Create() error
	Close() error
	CheckSubcomponents() error
	String() string
}

// RunTypeComponentsHandler defines the run type components handler actions
type RunTypeComponentsHandler interface {
	ComponentHandler
	RunTypeComponentsHolder
}

// RunTypeComponentsHolder holds the run type components
type RunTypeComponentsHolder interface {
	BlockchainContextCreator() contexts.BlockchainContextCreator
	ExecutorCreator() executorFactory.ExecutorCreator
	RuntimeContextCreator() contexts.RuntimeContextCreator
	MeteringContextCreator() contexts.MeteringContextCreator
	OutputContextCreator() contexts.OutputContextCreator
	ExecuteOnSameContextHandler() execution.ExecuteOnSameContextHandler
	CreateNewContractHandler() execution.CreateNewContractHandler
	GasScheduleFactory() config.GasScheduleFactory
	Create() error
	Close() error
	CheckSubcomponents() error
	String() string
	IsInterfaceNil() bool
}
