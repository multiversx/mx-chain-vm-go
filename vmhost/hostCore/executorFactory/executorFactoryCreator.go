package executorFactory

import (
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

type executorFactory struct{}

// NewExecutorFactory create a new executor factory
func NewExecutorFactory() *executorFactory {
	return &executorFactory{}
}

// CreateExecutor creates a new executor
func (ef *executorFactory) CreateExecutor(host vmhost.VMHost, hostParameters *vmhost.VMHostParameters) (executor.Executor, error) {
	vmHooks := vmhooks.NewVMHooksImpl(host)
	gasCostConfig, err := config.CreateGasConfig(host.GetGasScheduleMap())
	if err != nil {
		return nil, err
	}

	var vmExecutorFactory executor.ExecutorAbstractFactory

	if hostParameters.OverrideVMExecutor != nil {
		vmExecutorFactory = hostParameters.OverrideVMExecutor
	} else {
		vmExecutorFactory = wasmer2.ExecutorFactory()
	}
	vmExecutorFactoryArgs := executor.ExecutorFactoryArgs{
		VMHooks:                  vmHooks,
		OpcodeCosts:              gasCostConfig.WASMOpcodeCost,
		RkyvSerializationEnabled: true,
		WasmerSIGSEGVPassthrough: hostParameters.WasmerSIGSEGVPassthrough,
	}
	return vmExecutorFactory.CreateExecutor(vmExecutorFactoryArgs)
}

// IsInterfaceNil returns true if underlying object is nil
func (ef *executorFactory) IsInterfaceNil() bool {
	return ef == nil
}
