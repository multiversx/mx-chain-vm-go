package execution

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type executeOnSameContextHandler struct{}

// NewExecuteOnSameContextHandler creates new execute on same context handler
func NewExecuteOnSameContextHandler() *executeOnSameContextHandler {
	return &executeOnSameContextHandler{}
}

// PrepareSameContext updates recipient address and returns the code address
func (h *executeOnSameContextHandler) PrepareSameContext(
	input vmcommon.ContractCallInputHandler,
) ([]byte, error) {
	librarySCAddress := make([]byte, len(input.GetRecipientAddr()))
	copy(librarySCAddress, input.GetRecipientAddr())

	input.SetRecipientAddr(input.GetVMInput().CallerAddr)

	return librarySCAddress, nil
}

// ExecuteOnSameContextTransferValue will execute on same context transfer value
func (h *executeOnSameContextHandler) ExecuteOnSameContextTransferValue(
	input vmcommon.ContractCallInputHandler,
	output vmhost.OutputContext,
	runtime vmhost.RuntimeContext,
) error {
	return ExecuteSameContextTransferValue(input, output, runtime)
}
