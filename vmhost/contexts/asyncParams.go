//nolint:all
package contexts

import (
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// AddAsyncArgumentsToOutputTransfers called to process OutputTransfers created by a direct call (on dest) builtin function call by the VM
// it will add the asyncContext to the one output transfer it finds. Only one must exist.
func AddAsyncArgumentsToOutputTransfers(
	asyncParams *vmcommon.AsyncArguments,
	vmOutput *vmcommon.VMOutput,
) error {
	if asyncParams == nil {
		return nil
	}

	foundTransfer := false

	for _, outAcc := range vmOutput.OutputAccounts {
		for t, outTransfer := range outAcc.OutputTransfers {
			if outTransfer.CallType != vm.AsynchronousCall {
				continue
			}

			if foundTransfer {
				return vmhost.ErrTooManyTransfersFromBuiltInFunction
			}

			asyncData := createDataFromAsyncParams(asyncParams)
			outAcc.OutputTransfers[t] = vmcommon.OutputTransfer{
				Index:         outTransfer.Index,
				Value:         outTransfer.Value,
				GasLimit:      outTransfer.GasLimit,
				GasLocked:     outTransfer.GasLocked,
				AsyncData:     asyncData,
				Data:          outTransfer.Data,
				CallType:      outTransfer.CallType,
				SenderAddress: outTransfer.SenderAddress,
			}

			foundTransfer = true
		}
	}

	return nil
}

func createDataFromAsyncParams(asyncParams *vmcommon.AsyncArguments) []byte {
	callData := txDataBuilder.NewBuilder()
	callData.Bytes(asyncParams.CallID)
	callData.Bytes(asyncParams.CallerCallID)
	return callData.ToBytes()
}

// GenerateNewCallID will generate a new call ID as byte slice
func GenerateNewCallID(hasher crypto.Hasher, parentCallID []byte, suffix []byte) []byte {
	newCallID := append(parentCallID, suffix...)
	newCallID, err := hasher.Sha256(newCallID)
	if err != nil {
		return []byte{}
	}
	return newCallID
}
