package host

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

func TestAsync_NoAsyncCalls(t *testing.T) {
	code := arwen.GetTestSCCode("async-alice", "../../")
	host, _ := DefaultTestArwenForCall(t, code, nil)

	input := DefaultTestContractCallInput()
	input.GasProvided = 1000
	input.Function = "no_async"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{{42}}, vmOutput.ReturnData)

	async := host.Async()
	require.Equal(t, input.CallerAddr, async.GetCallerAddress())
	require.Equal(t, input.GasPrice, async.GetGasPrice())
	require.Empty(t, async.GetReturnData())
	require.True(t, async.IsComplete())
}

func TestAsync_OneAsyncCall(t *testing.T) {
	code := arwen.GetTestSCCode("async-alice", "../../")
	host, _ := DefaultTestArwenForCall(t, code, nil)

	input := DefaultTestContractCallInput()
	input.GasProvided = 1000
	input.Function = "one_async_call_no_cb"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	async := host.Async()
	require.Equal(t, input.CallerAddr, async.GetCallerAddress())
	require.Equal(t, input.GasPrice, async.GetGasPrice())
	require.Empty(t, async.GetReturnData())
	require.False(t, async.IsComplete())
}
