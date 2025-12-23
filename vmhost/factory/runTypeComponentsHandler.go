package factory

import (
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"

	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost/contexts"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/execution"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore/executorFactory"
)

const runTypeComponentsName = "managedRunTypeComponents"

var _ ComponentHandler = (*managedRunTypeComponents)(nil)
var _ RunTypeComponentsHandler = (*managedRunTypeComponents)(nil)
var _ RunTypeComponentsHolder = (*managedRunTypeComponents)(nil)

type managedRunTypeComponents struct {
	*RunTypeComponents
	factory              runTypeComponentsCreator
	mutRunTypeComponents sync.Mutex
}

// NewManagedRunTypeComponents returns a news instance of managed runType components
func NewManagedRunTypeComponents(rtc runTypeComponentsCreator) (*managedRunTypeComponents, error) {
	if rtc == nil {
		return nil, errNilRunTypeComponentsFactory
	}

	return &managedRunTypeComponents{
		RunTypeComponents: nil,
		factory:           rtc,
	}, nil
}

// Create will create the managed components
func (mrtc *managedRunTypeComponents) Create() error {
	vtc := mrtc.factory.Create()

	mrtc.mutRunTypeComponents.Lock()
	mrtc.RunTypeComponents = vtc
	mrtc.mutRunTypeComponents.Unlock()

	return nil
}

// Close will close all underlying subcomponents
func (mrtc *managedRunTypeComponents) Close() error {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}

	err := mrtc.RunTypeComponents.Close()
	if err != nil {
		return err
	}
	mrtc.RunTypeComponents = nil

	return nil
}

// CheckSubcomponents verifies all subcomponents
func (mrtc *managedRunTypeComponents) CheckSubcomponents() error {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return ErrNilRunTypeComponents
	}
	if check.IfNil(mrtc.BlockchainContextFactoryCreator) {
		return contexts.ErrNilBlockchainContextFactory
	}
	if check.IfNil(mrtc.ExecutorFactoryCreator) {
		return executorFactory.ErrNilCreateExecutorFactory
	}
	if check.IfNil(mrtc.RuntimeContextFactoryCreator) {
		return contexts.ErrNilRuntimeContextFactory
	}
	if check.IfNil(mrtc.MeteringContextFactoryCreator) {
		return contexts.ErrNilMeteringContextFactory
	}
	if check.IfNil(mrtc.OutputContextFactoryCreator) {
		return contexts.ErrNilOutputContextFactory
	}

	return nil
}

// BlockchainContextCreator returns the blockchain context factory creator
func (mrtc *managedRunTypeComponents) BlockchainContextCreator() contexts.BlockchainContextCreator {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.BlockchainContextFactoryCreator
}

// ExecutorCreator returns the executor factory creator
func (mrtc *managedRunTypeComponents) ExecutorCreator() executorFactory.ExecutorCreator {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.ExecutorFactoryCreator
}

// RuntimeContextCreator returns the runtime context factory creator
func (mrtc *managedRunTypeComponents) RuntimeContextCreator() contexts.RuntimeContextCreator {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.RuntimeContextFactoryCreator
}

// MeteringContextCreator returns the metering context factory creator
func (mrtc *managedRunTypeComponents) MeteringContextCreator() contexts.MeteringContextCreator {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.MeteringContextFactoryCreator
}

// OutputContextCreator returns the output context factory creator
func (mrtc *managedRunTypeComponents) OutputContextCreator() contexts.OutputContextCreator {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.OutputContextFactoryCreator
}

// ExecuteOnSameContextHandler returns vm execution same context handler
func (mrtc *managedRunTypeComponents) ExecuteOnSameContextHandler() execution.ExecuteOnSameContextHandler {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.ExecuteOnSameContextHandler
}

func (mrtc *managedRunTypeComponents) CreateNewContractHandler() execution.CreateNewContractHandler {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.CreateNewContractHandler
}

func (mrtc *managedRunTypeComponents) GasScheduleFactory() config.GasScheduleFactory {
	mrtc.mutRunTypeComponents.Lock()
	defer mrtc.mutRunTypeComponents.Unlock()

	if check.IfNil(mrtc.RunTypeComponents) {
		return nil
	}
	return mrtc.RunTypeComponents.GasScheduleFactoryCreator
}

// IsInterfaceNil returns true if underlying object is nil
func (mrtc *managedRunTypeComponents) IsInterfaceNil() bool {
	return mrtc == nil
}

// String returns the name of the component
func (mrtc *managedRunTypeComponents) String() string {
	return runTypeComponentsName
}
