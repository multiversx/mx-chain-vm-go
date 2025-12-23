package contexts

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type runtimeContextFactory struct{}

// NewRuntimeContextFactory create a new runtime context factory
func NewRuntimeContextFactory() *runtimeContextFactory {
	return &runtimeContextFactory{}
}

// CreateRuntimeContext create the runtime context
func (rcf *runtimeContextFactory) CreateRuntimeContext(
	host vmhost.VMHost,
	vmType []byte,
	builtInFuncContainer vmcommon.BuiltInFunctionContainer,
	vmExecutor executor.Executor,
	hasher vmhost.HashComputer,
) (vmhost.RuntimeContext, error) {
	scAPINames := vmExecutor.FunctionNames()
	validator := NewWASMValidator(scAPINames, builtInFuncContainer, host.EnableEpochsHandler())
	inputFactory := NewVMInputFactory()

	return NewRuntimeContext(host, vmType, vmExecutor, hasher, validator, inputFactory)
}

// IsInterfaceNil returns true if underlying object is nil
func (rcf *runtimeContextFactory) IsInterfaceNil() bool {
	return rcf == nil
}
