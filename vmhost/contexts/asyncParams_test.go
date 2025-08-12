package contexts

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/require"
)

func TestAddAsyncArgumentsToOutputTransfers(t *testing.T) {
	t.Parallel()

	t.Run("nil async params", func(t *testing.T) {
		t.Parallel()
		err := AddAsyncArgumentsToOutputTransfers(nil, nil, nil, 0, nil)
		require.Nil(t, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()
		output := &mockery.MockOutputContext{}
		vmOutput := &vmcommon.VMOutput{
			OutputAccounts: map[string]*vmcommon.OutputAccount{
				"addr1": {
					OutputTransfers: []vmcommon.OutputTransfer{
						{
							CallType: vm.AsynchronousCall,
						},
					},
				},
			},
		}
		asyncParams := &vmcommon.AsyncArguments{
			CallID: []byte("callID"),
		}

		err := AddAsyncArgumentsToOutputTransfers(output, []byte("addr1"), asyncParams, vm.AsynchronousCall, vmOutput)
		require.Nil(t, err)
		require.NotNil(t, vmOutput.OutputAccounts["addr1"].OutputTransfers[0].AsyncData)
	})
}

func TestCreateDataFromAsyncParams(t *testing.T) {
	t.Parallel()

	t.Run("nil async params", func(t *testing.T) {
		t.Parallel()
		data, err := createDataFromAsyncParams(nil, vm.DirectCall)
		require.Nil(t, err)
		require.Nil(t, data)
	})

	t.Run("async call", func(t *testing.T) {
		t.Parallel()
		asyncParams := &vmcommon.AsyncArguments{
			CallID:       []byte("callID"),
			CallerCallID: []byte("callerCallID"),
		}
		data, err := createDataFromAsyncParams(asyncParams, vm.AsynchronousCall)
		require.Nil(t, err)
		require.NotNil(t, data)
	})

	t.Run("async callback", func(t *testing.T) {
		t.Parallel()
		asyncParams := &vmcommon.AsyncArguments{
			CallID:                       []byte("callID"),
			CallerCallID:                 []byte("callerCallID"),
			CallbackAsyncInitiatorCallID: []byte("initiator"),
			GasAccumulated:               100,
		}
		data, err := createDataFromAsyncParams(asyncParams, vm.AsynchronousCallBack)
		require.Nil(t, err)
		require.NotNil(t, data)
		// a bit of a hack to check if the gas was encoded
		require.Contains(t, string(data), hex.EncodeToString(big.NewInt(100).Bytes()))
	})
}
