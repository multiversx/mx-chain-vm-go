package contexts

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const callbackNamePlaceholder = "<callback>"

// SendCrossShardCallback creates a transfer for a cross shard callback
func (context *asyncContext) SendCrossShardCallback() error {
	output := context.host.Output()
	_, lastTransfers := context.extractLastTransferToCaller(context.callerAddr, output.GetOutputAccounts())
	sender := context.address
	destination := context.callerAddr
	asyncData, data := context.createDataForCrossShardCallback(lastTransfers, output.ReturnCode(), output.ReturnData(), output.ReturnMessage())
	return sendCrossShardCallback(context.host, sender, destination, asyncData, data)
}

func (context *asyncContext) sendAsyncCallCrossShard(asyncCall *vmhost.AsyncCall) error {
	host := context.host
	runtime := host.Runtime()
	output := host.Output()

	function, arguments, err := context.callArgsParser.ParseData(string(asyncCall.GetData()))
	if err != nil {
		return err
	}

	context.incrementCallsCounter()

	newCallID := context.generateNewCallID()
	asyncCall.CallID = newCallID

	asyncData := createAsyncDataForAsyncCall(newCallID, context.GetCallID())

	callData := txDataBuilder.NewBuilder()
	callData.Func(function)
	for _, argument := range arguments {
		callData.Bytes(argument)
	}

	return output.Transfer(
		asyncCall.GetDestination(),
		runtime.GetContextAddress(),
		asyncCall.GetGasLimit(),
		asyncCall.GetGasLocked(),
		big.NewInt(0).SetBytes(asyncCall.GetValue()),
		asyncData,
		callData.ToBytes(),
		vm.AsynchronousCall,
	)
}

func createAsyncDataForAsyncCall(newCallID []byte, currentCallID []byte) []byte {
	asyncData := txDataBuilder.NewBuilder()
	asyncData.Bytes(newCallID)
	asyncData.Bytes(currentCallID)
	return asyncData.ToBytes()
}

func sendCrossShardCallback(host vmhost.VMHost, sender []byte, destination []byte, asyncData []byte, data []byte) error {
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
		asyncData,
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

func (context *asyncContext) createDataForCrossShardCallback(
	lastTransfers []byte,
	returnCode vmcommon.ReturnCode,
	returnData [][]byte,
	returnMessage string,
) ([]byte, []byte) {
	asyncData := txDataBuilder.NewBuilder()
	asyncData.Bytes(context.generateNewCallID())
	asyncData.Bytes(context.callID)
	asyncData.Bytes(context.callerCallID)
	asyncData.Bytes(big.NewInt(int64(context.gasAccumulated)).Bytes())

	transferData := txDataBuilder.NewBuilder()
	// This is just a placeholder, necessary not to break decoding, it's not used anywhere.
	transferData.Func(callbackNamePlaceholder)
	if lastTransfers != nil {
		transferData.Bytes(lastTransfers)
	}
	transferData.Bytes(ReturnCodeToBytes(returnCode))
	if returnCode == vmcommon.Ok {
		for _, data := range returnData {
			transferData.Bytes(data)
		}
	} else {
		transferData.Str(returnMessage)
	}
	return asyncData.ToBytes(), transferData.ToBytes()
}
