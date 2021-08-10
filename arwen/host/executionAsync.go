package host

import (
	"encoding/hex"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	async := host.Async()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	legacyGroupID := arwen.LegacyAsyncCallGroupID
	legacyGroup, exists := async.GetCallGroup(legacyGroupID)
	if !exists {
		return arwen.ErrLegacyAsyncCallNotFound

	}

	if legacyGroup.IsComplete() {
		return arwen.ErrLegacyAsyncCallInvalid
	}

	return nil
}

// TODO compare with asyncContext.sendContextCallbackToOriginalCaller()
func (host *vmHost) sendAsyncCallbackToCaller() error {
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	currentCall := runtime.GetVMInput()

	retCode := output.ReturnCode()
	retCodeBytes := big.NewInt(int64(retCode)).Bytes()
	retData := []byte("@" + hex.EncodeToString(runtime.GetPrevTxHash()))
	retData = append(retData, []byte("@"+hex.EncodeToString(retCodeBytes))...)
	if retCode == vmcommon.Ok {
		for _, data := range output.ReturnData() {
			retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
		}
	} else {
		retMessage := []byte(output.ReturnMessage())
		retData = append(retData, []byte("@"+hex.EncodeToString(retMessage))...)
	}

	gasLeft := metering.GasLeft()

	err := output.Transfer(
		currentCall.CallerAddr,
		runtime.GetSCAddress(),
		gasLeft,
		0,
		currentCall.CallValue,
		retData,
		vm.AsynchronousCallBack,
	)
	metering.UseGas(gasLeft)
	if err != nil {
		runtime.FailExecution(err)
		return err
	}

	log.Trace(
		"sendAsyncCallbackToCaller",
		"caller", currentCall.CallerAddr,
		"data", retData,
		"gas", gasLeft)

	return nil
}
