package execution

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type ExecuteOnSameContextHandler interface {
	PrepareSameContext(
		input vmcommon.ContractCallInputHandler,
	) ([]byte, error)
	ExecuteOnSameContextTransferValue(
		input vmcommon.ContractCallInputHandler,
		output vmhost.OutputContext,
		runtime vmhost.RuntimeContext,
	) error
}

type CreateNewContractHandler interface {
	DeductGasForContractDeployment(
		input vmcommon.ContractCreateInputHandler,
		meteringContext vmhost.MeteringContext,
		codeDeployInput vmhost.CodeDeployInput,
	) error
}
