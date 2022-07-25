package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_5/crypto"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
)

func RemoveAsyncContextArguments(input *vmcommon.VMInput) ([][]byte, error) {
	var err error
	if IsCallAsync(input.CallType) {
		var callID, callerCallID, callbackAsyncInitiatorCallID, gasAccumulated []byte
		callID, err = PopFirstArgumentFromVMInput(input)
		if err != nil {
			return nil, err
		}

		callerCallID, err = PopFirstArgumentFromVMInput(input)
		if err != nil {
			return nil, err
		}

		if IsCallback(input.CallType) {
			callbackAsyncInitiatorCallID, err = PopFirstArgumentFromVMInput(input)
			if err != nil {
				return nil, err
			}
			gasAccumulated, err = PopFirstArgumentFromVMInput(input)
			if err != nil {
				return nil, err
			}
			return [][]byte{callID, callerCallID, callbackAsyncInitiatorCallID, gasAccumulated}, nil
		} else {
			return [][]byte{callID, callerCallID}, nil
		}
	}

	return nil, nil
}

func AddAsyncParamsToVmOutput(
	address []byte,
	asyncParams [][]byte,
	callType vm.CallType,
	parseDataFunc func(data string) (string, [][]byte, error),
	vmOutput *vmcommon.VMOutput) error {
	if asyncParams == nil {
		return nil
	}
	for _, outAcc := range vmOutput.OutputAccounts {
		if !bytes.Equal(address, outAcc.Address) {
			continue
		}

		for t, outTransfer := range outAcc.OutputTransfers {
			if outTransfer.CallType != callType {
				continue
			}

			newData, err := AppendAsyncParamsToCallData(
				asyncParams,
				outTransfer.Data,
				parseDataFunc)

			if err != nil {
				return err
			}

			outAcc.OutputTransfers[t] = vmcommon.OutputTransfer{
				Value:         outTransfer.Value,
				GasLimit:      outTransfer.GasLimit,
				GasLocked:     outTransfer.GasLocked,
				Data:          newData,
				CallType:      outTransfer.CallType,
				SenderAddress: outTransfer.SenderAddress,
			}
		}
	}

	return nil
}

func AppendAsyncParamsToCallData(
	asyncParams [][]byte,
	data []byte,
	parseDataFunc func(data string) (string, [][]byte, error)) ([]byte, error) {
	function, args, err := parseDataFunc(string(data))
	if err != nil {
		return nil, err
	}

	callData := txDataBuilder.NewBuilder()
	callData.Func(function)
	for _, asyncParam := range asyncParams {
		callData.Bytes(asyncParam)
	}

	for _, arg := range args {
		callData.Bytes(arg)
	}

	return callData.ToBytes(), nil
}

func AppendAsyncParamsToArguments(
	asyncParams [][]byte,
	data []byte,
	parseArgumentsFunc func(data string) ([][]byte, error)) ([]byte, error) {
	args, err := parseArgumentsFunc(string(data))
	if err != nil {
		return nil, err
	}

	callData := txDataBuilder.NewBuilder()
	for _, asyncParam := range asyncParams {
		callData.Bytes(asyncParam)
	}

	// args string start with a @ so first parsed argument will be empty always
	if len(args) != 0 {
		for _, arg := range args[1:] {
			callData.Bytes(arg)
		}
	} else {
		callData.Bytes([]byte{})
	}

	return callData.ToBytes(), nil
}

func GenerateNewCallID(hasher crypto.Hasher, parentCallID []byte, suffix []byte) []byte {
	newCallID := append(parentCallID, suffix...)
	newCallID, err := hasher.Sha256(newCallID)
	if err != nil {
		return []byte{}
	}
	return newCallID
}

func CreateCallbackAsyncParams(hasher crypto.Hasher, asyncParams [][]byte) [][]byte {
	if asyncParams == nil {
		return nil
	}
	newAsyncParams := make([][]byte, 4)
	newAsyncParams[0] = GenerateNewCallID(hasher, asyncParams[0], []byte{0})
	newAsyncParams[1] = asyncParams[0]
	newAsyncParams[2] = asyncParams[1]
	newAsyncParams[3] = []byte{0}
	return newAsyncParams
}
