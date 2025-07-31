//nolint:all
package contexts

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

/*
	Called to process OutputTransfers created by a
	direct call (on dest) builtin function call by the VM
*/
func AddAsyncArgumentsToOutputTransfers(
	output vmhost.OutputContext,
	address []byte,
	asyncParams *vmcommon.AsyncArguments,
	callType vm.CallType,
	vmOutput *vmcommon.VMOutput) error {
	if asyncParams == nil {
		return nil
	}
	for _, outAcc := range vmOutput.OutputAccounts {
		// if !bytes.Equal(address, outAcc.Address) {
		// 	continue
		// }

		for t, outTransfer := range outAcc.OutputTransfers {
			// if !bytes.Equal(address, outTransfer.SenderAddress) {
			// 	continue
			// }
			if outTransfer.CallType != callType {
				continue
			}

			asyncData, err := createDataFromAsyncParams(
				asyncParams,
				callType)

			if err != nil {
				return err
			}

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
		}
	}

	return nil
}

func createDataFromAsyncParams(
	asyncParams *vmcommon.AsyncArguments,
	callType vm.CallType,
) ([]byte, error) {
	if asyncParams == nil {
		if callType == vm.AsynchronousCall || callType == vm.AsynchronousCallBack {
			return nil, vmcommon.ErrAsyncParams
		} else {
			return nil, nil
		}
	}

	callData := txDataBuilder.NewBuilder()
	callData.Bytes(asyncParams.CallID)
	callData.Bytes(asyncParams.CallerCallID)
	if callType == vm.AsynchronousCallBack {
		callData.Bytes(asyncParams.CallbackAsyncInitiatorCallID)
		callData.Bytes(big.NewInt(int64(asyncParams.GasAccumulated)).Bytes())
	}

	return callData.ToBytes(), nil
}


/*
	Used by when a callback SCR is created
	1)	after a failure of an async call
		Async data is extracted (by extractAsyncCallParamsFromTxData()) and then
		reappended to the new SCR's callback data (by reapendAsyncParamsToTxData())
	2)	from the last transfer (see useLastTransferAsAsyncCallBackWhenNeeded())
*/
func CreateCallbackAsyncParams(hasher crypto.Hasher, asyncParams *vmcommon.AsyncArguments) [][]byte {
	if asyncParams == nil {
		return nil
	}
	newAsyncParams := make([][]byte, 4)
	newAsyncParams[0] = GenerateNewCallID(hasher, asyncParams.CallID, []byte{0})
	newAsyncParams[1] = asyncParams.CallID
	newAsyncParams[2] = asyncParams.CallerCallID
	newAsyncParams[3] = []byte{0}
	return newAsyncParams
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
