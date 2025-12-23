package execution

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type createNewContractHandler struct{}

// NewCreateNewContractHandler creates new contract handler
func NewCreateNewContractHandler() *createNewContractHandler {
	return &createNewContractHandler{}
}

// DeductGasForContractDeployment will deduct gas for contract deployment
func (f *createNewContractHandler) DeductGasForContractDeployment(
	_ vmcommon.ContractCreateInputHandler,
	meteringContext vmhost.MeteringContext,
	codeDeployInput vmhost.CodeDeployInput,
) error {
	return meteringContext.DeductInitialGasForIndirectDeployment(codeDeployInput)
}
