package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestElrondEI_CallValue(t *testing.T) {
	code := GetTestSCCode("elrondei", "../../")

	// 1-byte call value
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_1byte"
	input.CallValue = big.NewInt(64)

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 3)
	data := vmOutput.ReturnData
	require.Equal(t, []byte("ok"), data[0])
	require.Equal(t, []byte{1, 0, 0, 0}, data[1])
	require.Equal(t, []byte{64}, data[2])

	// 4-byte call value
	host, _ = DefaultTestArwenForCall(t, code)
	input = DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_4bytes"
	input.CallValue = big.NewInt(0).SetBytes([]byte{64, 12, 16, 99})

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 3)
	data = vmOutput.ReturnData
	require.Equal(t, []byte("ok"), data[0])
	require.Equal(t, []byte{4, 0, 0, 0}, data[1])
	require.Equal(t, []byte{64, 12, 16, 99}, data[2])
}
