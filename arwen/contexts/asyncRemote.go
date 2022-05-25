package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_5/arwen"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
)

const callbackNamePlaceholder = "<callback>"

// SendCrossShardCallback creates a transfer for a cross shard callback
func (context *asyncContext) SendCrossShardCallback(
	returnCode vmcommon.ReturnCode,
	returnData [][]byte,
	returnMessage string,
) error {
	sender := context.address
	destination := context.callerAddr
	data := context.createCallbackArgumentsForCrossShardCallback(returnCode, returnData, returnMessage)
	return sendCrossShardCallback(context.host, sender, destination, data)
}

func (context *asyncContext) sendAsyncCallCrossShard(asyncCall *arwen.AsyncCall) error {
	host := context.host
	runtime := host.Runtime()
	output := host.Output()

	function, arguments, err := context.callArgsParser.ParseData(string(asyncCall.GetData()))
	if err != nil {
		return err
	}

	context.incrementCallsCounter()
	newCallID := context.generateNewCallID()
	callData := txDataBuilder.NewBuilder()
	callData.Func(function)
	callData.Bytes(newCallID)
	callData.Bytes(context.GetCallID())

	asyncCall.CallID = newCallID

	for _, argument := range arguments {
		callData.Bytes(argument)
	}

	return output.Transfer(
		asyncCall.GetDestination(),
		runtime.GetSCAddress(),
		asyncCall.GetGasLimit(),
		asyncCall.GetGasLocked(),
		big.NewInt(0).SetBytes(asyncCall.GetValue()),
		callData.ToBytes(),
		vm.AsynchronousCall,
	)
}

func sendCrossShardCallback(host arwen.VMHost, sender []byte, destination []byte, data []byte) error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	gasLeft := metering.GasLeft()
	metering.UseGas(gasLeft)
	err := output.Transfer(
		destination,
		sender,
		gasLeft,
		0,
		big.NewInt(0),
		data,
		vm.AsynchronousCallBack,
	)
	if err != nil {
		runtime.FailExecution(err)
		return err
	}

	logAsync.Trace(
		"sendCrossShardCallback",
		"caller", currentCall.CallerAddr,
		"data", data,
		"gas", gasLeft)

	return nil
}

func (context *asyncContext) createCallbackArgumentsForCrossShardCallback(
	returnCode vmcommon.ReturnCode,
	returnData [][]byte,
	returnMessage string,
) []byte {
	transferData := txDataBuilder.NewBuilder()

	// This is just a placeholder, necessary not to break decoding, it's not used anywhere.
	transferData.Func(callbackNamePlaceholder)

	transferData.Bytes(context.generateNewCallID())
	transferData.Bytes(context.callID)
	transferData.Bytes(context.callerCallID)
	transferData.Bytes(big.NewInt(int64(context.gasAccumulated)).Bytes())

	transferData.Int64(int64(returnCode))
	if returnCode == vmcommon.Ok {
		for _, data := range returnData {
			transferData.Bytes(data)
		}
	} else {
		transferData.Str(returnMessage)
	}
	return transferData.ToBytes()
}
