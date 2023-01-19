package hostCoretest

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseOpsAPI_CallValue(t *testing.T) {
	code := testcommon.GetTestSCCode("baseOps", "../../")

	// 1-byte call value
	host, _ := testcommon.DefaultTestVMForCall(t, code, nil)
	defer func() {
		host.Reset()
	}()

	input := testcommon.DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_1byte"
	input.CallValue = big.NewInt(64)

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	assert.Equal(t, "", vmOutput.ReturnMessage)
	data := vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{32, 0, 0, 0}, data[1])
	assert.Equal(t,
		[]byte{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 64,
		},
		data[2])

	// 4-byte call value
	host, _ = testcommon.DefaultTestVMForCall(t, code, nil)
	input = testcommon.DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_4bytes"
	input.CallValue = big.NewInt(0).SetBytes([]byte{64, 12, 16, 99})

	vmOutput, err = host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	data = vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{32, 0, 0, 0}, data[1])
	assert.Equal(t,
		[]byte{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 64, 12, 16, 99,
		},
		data[2])

	// BigInt call value
	host, _ = testcommon.DefaultTestVMForCall(t, code, nil)
	input = testcommon.DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_bigInt_to_Bytes"
	input.CallValue = big.NewInt(19*256 + 233)

	vmOutput, err = host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 4)
	data = vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{32, 0, 0, 0}, data[1])
	assert.Equal(t,
		[]byte{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 19, 233,
		},
		data[2])

	val12345 := big.NewInt(0).SetBytes(data[3])
	assert.Equal(t, big.NewInt(12345), val12345)
}

func TestBaseOpsAPI_int64getArgument(t *testing.T) {
	code := testcommon.GetTestSCCode("baseOps", "../../")
	host, _ := testcommon.DefaultTestVMForCall(t, code, nil)
	defer func() {
		host.Reset()
	}()

	input := testcommon.DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_int64getArgument"
	input.Arguments = [][]byte{big.NewInt(12345).Bytes()}

	vmOutput, err := host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	data := vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{57, 48, 0, 0}, data[1])

	invBytes := vmhost.InverseBytes(data[1])
	val12345 := big.NewInt(0).SetBytes(invBytes)
	assert.Equal(t, big.NewInt(12345), val12345)

	i64val12345 := big.NewInt(0).SetBytes(data[2])
	assert.Equal(t, big.NewInt(12345), i64val12345)

	// Take the result of the SC method (the number 12345 as bytes, received from
	// the SC in data[2]) and feed it back into the SC method.
	input.Arguments = [][]byte{data[2]}

	vmOutput, err = host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	data = vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{57, 48, 0, 0}, data[1])

	invBytes = vmhost.InverseBytes(data[1])
	val12345 = big.NewInt(0).SetBytes(invBytes)
	assert.Equal(t, big.NewInt(12345), val12345)

	i64val12345 = big.NewInt(0).SetBytes(data[2])
	assert.Equal(t, big.NewInt(12345), i64val12345)
}
