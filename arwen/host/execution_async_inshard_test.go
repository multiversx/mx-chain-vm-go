package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

func TestAsync_NoAsyncCalls(t *testing.T) {
	code := arwen.GetTestSCCodeModule("promises/parent-simple", "parent-simple", "../../")
	host, _ := defaultTestArwenForCall(t, code, nil)

	input := DefaultTestContractCallInput()
	input.GasProvided = 10000000
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
	parentCode := arwen.GetTestSCCodeModule("promises/parent-simple", "parent-simple", "../../")
	childCode := arwen.GetTestSCCodeModule("promises/child-simple", "child-simple", "../../")
	balance := big.NewInt(100)
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, balance, balance)

	input := DefaultTestContractCallInput()
	input.GasProvided = 10000000
	input.Function = "one_async_call_no_cb_with_call_value"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	async := host.Async()
	require.Equal(t, input.CallerAddr, async.GetCallerAddress())
	require.Equal(t, input.GasPrice, async.GetGasPrice())
	require.Empty(t, async.GetReturnData())
	require.True(t, async.IsComplete())
}
