package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

func TestAsync_NoAsyncCalls(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	hostGasTester := mock.NewHostGasTester(t, host, 32)

	parentSC := ibm.CreateAndStoreInstanceMock(parentAddress, 0)
	parentSC.AddMockMethod("no_async", func() {
		hostGasTester.UseGas(100)

		finish := []byte("forty two")
		host.Output().Finish(finish)
		hostGasTester.UseGasForLastFinish()
	})
	hostGasTester.UseGasForContractCode()

	vmInput := DefaultTestContractCallInput()
	vmInput.GasProvided = 1000
	vmInput.Function = "no_async"

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{[]byte("forty two")}, vmOutput.ReturnData)

	hostGasTester.Validate(vmInput, vmOutput)

	async := host.Async()
	require.Equal(t, vmInput.CallerAddr, async.GetCallerAddress())
	require.Equal(t, vmInput.GasPrice, async.GetGasPrice())
	require.Empty(t, async.GetReturnData())
	require.True(t, async.IsComplete())
}

func TestAsync_SimpleAsyncCall_NoCallbacks(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	hostGasTester := mock.NewHostGasTester(t, host, 32)

	parentSC := ibm.CreateAndStoreInstanceMock(parentAddress, 40)
	parentSC.AddMockMethod("register_async_call_no_callbacks", func() {
		hostGasTester.UseGas(300)

		host.Async().RegisterAsyncCall("testgroup", &arwen.AsyncCall{
			Destination:     childAddress,
			Data:            []byte("childmethod"),
			GasLimit:        100,
			ValueBytes:      big.NewInt(10).Bytes(),
			SuccessCallback: "",
			ErrorCallback:   "",
		})
		hostGasTester.UseGasForAPI()
	})
	hostGasTester.UseGasForContractCode()

	childSC := ibm.CreateAndStoreInstanceMock(childAddress, 0)
	childSC.AddMockMethod("childmethod", func() {
		hostGasTester.UseGasForAPI() // AsyncCallStep

		require.Nil(t, host.Runtime().Arguments())
		require.Equal(t, big.NewInt(10), host.Runtime().GetVMInput().CallValue)

		hostGasTester.UseGas(100)
	})
	hostGasTester.UseGasForContractCode()

	vmInput := DefaultTestContractCallInput()
	vmInput.GasProvided = 1000
	vmInput.Function = "register_async_call_no_callbacks"

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.True(t, host.Async().IsComplete())

	hostGasTester.Validate(vmInput, vmOutput)
}

func TestAsync_AsyncCall_NoCallbacks(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	hostGasTester := mock.NewHostGasTester(t, host, 32)

	parentSC := ibm.CreateAndStoreInstanceMock(parentAddress, 40)
	parentSC.AddMockMethod("register_async_call_no_callbacks", func() {

		host.Async().RegisterAsyncCall("testgroup", &arwen.AsyncCall{
			Destination:     childAddress,
			Data:            []byte("childmethod@02"),
			GasLimit:        100,
			ValueBytes:      big.NewInt(10).Bytes(),
			SuccessCallback: "",
			ErrorCallback:   "",
		})

		hostGasTester.UseGas(300)
	})
	hostGasTester.UseGasForContractCode()

	// This method must exist, but must not be called.
	parentSC.AddMockMethod("callBack", func() {
		host.Output().Finish([]byte("callback called but shouldn't have been"))
		hostGasTester.UseGasForLastFinish()
		hostGasTester.UseGas(29)
	})

	childSC := ibm.CreateAndStoreInstanceMock(childAddress, 0)
	childSC.AddMockMethod("childmethod", func() {
		require.Equal(t, [][]byte{{2}}, host.Runtime().Arguments())
		require.Equal(t, big.NewInt(10), host.Runtime().GetVMInput().CallValue)

		host.Output().Finish([]byte("childmethod called"))
		hostGasTester.UseGasForLastFinish()
		hostGasTester.UseGas(20)
	})
	hostGasTester.UseGasForContractCode()

	vmInput := DefaultTestContractCallInput()
	vmInput.GasProvided = 1000
	vmInput.Function = "register_async_call_no_callbacks"

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.True(t, host.Async().IsComplete())

	hostGasTester.Validate(vmInput, vmOutput)
}
