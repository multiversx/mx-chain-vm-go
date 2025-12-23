package execution

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

func ExecuteSameContextTransferValue(
	input vmcommon.ContractCallInputHandler,
	output vmhost.OutputContext,
	runtime vmhost.RuntimeContext,
) error {
	err := output.TransferValueOnly(input.GetRecipientAddr(), input.GetVMInput().CallerAddr, input.GetVMInput().CallValue, false)
	if err != nil {
		runtime.AddError(err, input.GetFunction())
		return err
	}

	return nil
}
