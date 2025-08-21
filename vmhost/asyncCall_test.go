package vmhost

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/stretchr/testify/require"
)

func TestAsyncCall_Clone(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{
		CallID:          []byte("call_id"),
		Status:          AsyncCallPending,
		ExecutionMode:   SyncExecution,
		Destination:     []byte("destination"),
		Data:            []byte("data"),
		GasLimit:        1000,
		GasLocked:       2000,
		ValueBytes:      []byte("value"),
		SuccessCallback: "success",
		ErrorCallback:   "error",
		CallbackClosure: []byte("closure"),
	}

	clone := ac.Clone()
	require.NotSame(t, ac, clone)
	require.Equal(t, ac, clone)

	// Check that changing the clone does not affect the original
	clone.GasLimit = 500
	require.NotEqual(t, ac.GasLimit, clone.GasLimit)
}

func TestAsyncCall_Getters(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{
		CallID:          []byte("call_id"),
		Destination:     []byte("destination"),
		Data:            []byte("data"),
		GasLimit:        1000,
		GasLocked:       2000,
		ValueBytes:      []byte("value"),
		SuccessCallback: "success",
		ErrorCallback:   "error",
	}

	require.Equal(t, []byte("call_id"), ac.GetIdentifier())
	require.Equal(t, []byte("destination"), ac.GetDestination())
	require.Equal(t, []byte("data"), ac.GetData())
	require.Equal(t, uint64(1000), ac.GetGasLimit())
	require.Equal(t, uint64(2000), ac.GetGasLocked())
	require.Equal(t, uint64(3000), ac.GetTotalGas())
	require.Equal(t, []byte("value"), ac.GetValue())
}

func TestAsyncCall_IsLocalRemote(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{ExecutionMode: SyncExecution}
	require.True(t, ac.IsLocal())
	require.False(t, ac.IsRemote())

	ac.ExecutionMode = AsyncBuiltinFuncCrossShard
	require.False(t, ac.IsLocal())
	require.True(t, ac.IsRemote())

	ac.ExecutionMode = AsyncUnknown
	require.False(t, ac.IsLocal())
	require.True(t, ac.IsRemote())
}

func TestAsyncCall_HasCallback(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{Status: AsyncCallResolved, SuccessCallback: "success"}
	require.True(t, ac.HasCallback())

	ac.Status = AsyncCallRejected
	ac.ErrorCallback = "error"
	require.True(t, ac.HasCallback())

	ac.SuccessCallback = ""
	ac.ErrorCallback = ""
	require.False(t, ac.HasCallback())
}

func TestAsyncCall_HasDefinedAnyCallback(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{}
	require.False(t, ac.HasDefinedAnyCallback())

	ac.SuccessCallback = "success"
	require.True(t, ac.HasDefinedAnyCallback())

	ac.SuccessCallback = ""
	ac.ErrorCallback = "error"
	require.True(t, ac.HasDefinedAnyCallback())

	ac.SuccessCallback = "success"
	ac.ErrorCallback = "error"
	require.True(t, ac.HasDefinedAnyCallback())
}

func TestAsyncCall_UpdateStatus(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{}
	ac.UpdateStatus(vmcommon.Ok)
	require.Equal(t, AsyncCallResolved, ac.Status)

	ac.UpdateStatus(1) // Using a non-OK status code
	require.Equal(t, AsyncCallRejected, ac.Status)
}

func TestAsyncCall_Reject(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{}
	ac.Reject()
	require.Equal(t, AsyncCallRejected, ac.Status)
}

func TestAsyncCall_GetCallbackName(t *testing.T) {
	t.Parallel()

	ac := &AsyncCall{SuccessCallback: "success", ErrorCallback: "error"}

	ac.Status = AsyncCallResolved
	require.Equal(t, "success", ac.GetCallbackName())

	ac.Status = AsyncCallRejected
	require.Equal(t, "error", ac.GetCallbackName())

	ac.Status = AsyncCallPending
	require.Equal(t, "error", ac.GetCallbackName())
}

func TestAsyncCall_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var ac *AsyncCall
	require.True(t, ac.IsInterfaceNil())

	ac = &AsyncCall{}
	require.False(t, ac.IsInterfaceNil())
}
