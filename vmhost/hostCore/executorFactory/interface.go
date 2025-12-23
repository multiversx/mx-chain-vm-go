package executorFactory

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// ExecutorCreator defines the executor factory creator
type ExecutorCreator interface {
	CreateExecutor(vmHost vmhost.VMHost, hostParameters *vmhost.VMHostParameters) (executor.Executor, error)
	IsInterfaceNil() bool
}
