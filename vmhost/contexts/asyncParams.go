//nolint:all
package contexts

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	"github.com/multiversx/wasm-vm/crypto"
)

/*
	Called to process OutputTransfers created by a
	direct call (on dest) builtin function call by the VM
    TODO(fix) this function
*/
func AddAsyncArgumentsToOutputTransfers(
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
	Called when a SCR for a callback is created outside the VM
	(by createAsyncCallBackSCRFromVMOutput())
	This is the case
	A)	after an async call executed following a builtin function call,
	B)	other cases where processing the output trasnfers of a VMOutput did
		not produce a SCR of type AsynchronousCallBack
    TODO(check): function not used?
*/
func AppendAsyncArgumentsToCallbackCallData(
	hasher crypto.Hasher,
	data []byte,
	asyncArguments *vmcommon.AsyncArguments,
	parseArgumentsFunc func(data string) ([][]byte, error)) ([]byte, error) {

	return appendAsyncParamsToCallData(
		CreateCallbackAsyncParams(hasher, asyncArguments),
		data,
		false,
		parseArgumentsFunc)
}

/*
	Called when a SCR is created from VMOutput in order to recompose
	async data and call data into a transfer data ready for the SCR
	(by preprocessOutTransferToSCR())
    TODO(check): function not used?
*/
func AppendTransferAsyncDataToCallData(
	callData []byte,
	asyncData []byte,
	parseArgumentsFunc func(data string) ([][]byte, error)) ([]byte, error) {

	var asyncParams [][]byte
	if asyncData != nil {
		asyncParams, _ = parseArgumentsFunc(string(asyncData))
		// string start with a @ so first parsed argument will be empty always
		asyncParams = asyncParams[1:]
	} else {
		return callData, nil
	}

	return appendAsyncParamsToCallData(
		asyncParams,
		callData,
		true,
		parseArgumentsFunc)
}

func appendAsyncParamsToCallData(
	asyncParams [][]byte,
	data []byte,
	hasFunction bool,
	parseArgumentsFunc func(data string) ([][]byte, error)) ([]byte, error) {

	if data == nil {
		return nil, nil
	}

	args, err := parseArgumentsFunc(string(data))
	if err != nil {
		return nil, err
	}

	var functionName string
	if hasFunction {
		functionName = string(args[0])
	}

	// check if there is only one argument and that is 0
	if len(args) != 0 {
		args = args[1:]
	}

	callData := txDataBuilder.NewBuilder()

	if functionName != "" {
		callData.Func(functionName)
	}

	if len(args) != 0 {
		for _, arg := range args {
			callData.Bytes(arg)
		}
	} else {
		if !hasFunction {
			callData.Bytes([]byte{})
		}
	}

	for _, asyncParam := range asyncParams {
		callData.Bytes(asyncParam)
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
