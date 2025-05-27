package hostCoretest

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/mock/contracts"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupReadOnlyTest creates a common setup for testing ManagedExecuteReadOnly.
// ContractA calls a function in ContractB via ManagedExecuteReadOnly.
// ContractB's function then attempts a restricted operation.
func setupReadOnlyTest(
	t *testing.T,
	targetContractAddress []byte,
	targetFunctionName string,
	targetSetup func(host vmhost.VMHost, instance *mock.InstanceMock), // To set up ContractB's method
	initialGas uint64,
) (vmhost.VMHost, *vmcommon.ContractCallInput) {
	host := testcommon.NewTestHostBuilder(t).Build()

	// ContractA (Caller)
	callerContract := testcommon.CreateMockContract(testcommon.CallerAddress)
	callerContract.SetMethods(func(instance *mock.InstanceMock, cfg interface{}) {
		instance.AddMockMethod("callTargetReadOnly", func() *mock.InstanceMock {
			h := instance.Host
			mt := h.ManagedTypes()

			targetAddrHandle := mt.NewManagedBufferFromBytes(targetContractAddress)
			targetFuncHandle := mt.NewManagedBufferFromBytes([]byte(targetFunctionName))
			argsHandle := mt.NewManagedBuffer() // No args for simplicity

			// Handles for results are not strictly needed if we expect failure
			resultBufHandle := mt.NewManagedBuffer()

			vmh := vmhooks.NewVMHooksImpl(h)
			ret := vmh.ManagedExecuteReadOnly(
				initialGas/2,      // Gas for the read-only call
				targetAddrHandle,  // Target contract address
				targetFuncHandle,  // Target function name
				argsHandle,        // Arguments
				resultBufHandle,   // Result handle
			)

			// If ManagedExecuteReadOnly itself returns an error (e.g. -1),
			// it means the call failed as expected at that level.
			// If it returns 0 (success), the error should be in Runtime errors.
			if ret != 0 && h.GetRuntimeErrors() == nil {
				h.Runtime().SignalUserError("ManagedExecuteReadOnly failed unexpectedly at EEI level but no runtime error")
			}

			// Propagate any error message from the read-only context if it was set
			if h.GetRuntimeErrors() != nil {
				errBytes, _ := mt.GetBytes(resultBufHandle)
				if len(errBytes) > 0 {
					// This might be tricky if the error isn't directly put into resultHandle
					// For now, we rely on host.GetRuntimeErrors()
				}
			}
			return instance
		})
	})

	// ContractB (Target)
	targetContract := testcommon.CreateMockContract(targetContractAddress)
	targetContract.SetMethods(func(instance *mock.InstanceMock, cfg interface{}) {
		targetSetup(instance.Host, instance)
	})

	host.ContextsHolder().AddContract(callerContract)
	host.ContextsHolder().AddContract(targetContract)

	callInput := testcommon.CreateTestContractCallInputBuilder().
		WithCallerAddr(testcommon.UserAccountAddress). // EOA calling ContractA
		WithRecipientAddr(testcommon.CallerAddress).   // ContractA
		WithGasProvided(initialGas).
		WithFunction("callTargetReadOnly").
		Build()

	return host, callInput
}

func TestManagedExecuteReadOnly_AttemptMBufferSetBytes_ShouldFail(t *testing.T) {
	targetAddr := testcommon.MakeSimpleAddress(2)
	targetFunc := "attemptMBufferSetBytes"
	initialGas := uint64(5_000_000)

	host, callInput := setupReadOnlyTest(t, targetAddr, targetFunc,
		func(h vmhost.VMHost, instance *mock.InstanceMock) {
			instance.AddMockMethod(targetFunc, func() *mock.InstanceMock {
				mt := h.ManagedTypes()
				// Assume bufferHandle 0 is passed or created if not passed by caller.
				// For this test, we'll create one inside, assuming it's a buffer owned by this context.
				// The principle is that even if it's a new buffer, modification should fail.
				bufferHandle := mt.NewManagedBufferFromBytes([]byte("original"))
				dataToSetHandle := mt.NewManagedBufferFromBytes([]byte("modified"))

				vmh := vmhooks.NewVMHooksImpl(h)
				// This is the call that should fail
				vmh.MBufferSetBytes(bufferHandle, 0, int32(len("modified"))) // Params for MBufferSetBytes are (handle, offset, length)
				// Actually, MBufferSetBytes takes (handle, dataOffset, dataLength) from WASM memory.
				// Let's simplify for the EEI call directly.
				// The direct EEI call is MBufferSetBytes(mBufferHandle int32, dataOffset executor.MemPtr, dataLength executor.MemLength)
				// We'll simulate the effect by calling the WithHost variant or directly the context method if possible.
				// For this test, we can directly use the context method if it's easier, or mock the MemLoad.

				// Simpler: Call the EEI that takes a handle to data
				// Let's assume we have an EEI that sets from another managed buffer for simplicity here,
				// or directly try to modify via context if that's what MBufferSetBytes does.
				// The actual MBufferSetBytes EEI reads from WASM memory.
				// We can directly call the ManagedTypesContext method for testing the ReadOnly check.
				// However, EEIs are the public API to SCs.

				// Let's use the EEI: MBufferSetBytes(mBufferHandle int32, dataOffset executor.MemPtr, dataLength executor.MemLength)
				// We need to put "modified" into WASM memory first.
				wasmMemDataOffset := int32(0)
				err := h.Runtime().GetInstance().MemStore(0, []byte("modified"))
				require.Nil(t, err)

				vmh.MBufferSetBytes(bufferHandle, 0, int32(len("modified"))) // Attempt to set bytes

				return instance
			})
		}, initialGas)

	defer host.Reset()

	vmOutput, err := host.RunSmartContractCall(callInput)
	require.Nil(t, err)
	assert.NotEqual(t, vmcommon.Ok, vmOutput.ReturnCode, "Execution should have failed")
	assert.Contains(t, vmOutput.ReturnMessage, vmhost.ErrInvalidCallOnReadOnlyMode.Error())
}

func TestManagedExecuteReadOnly_AttemptManagedMapPut_ShouldFail(t *testing.T) {
	targetAddr := testcommon.MakeSimpleAddress(2)
	targetFunc := "attemptManagedMapPut"
	initialGas := uint64(5_000_000)

	host, callInput := setupReadOnlyTest(t, targetAddr, targetFunc,
		func(h vmhost.VMHost, instance *mock.InstanceMock) {
			instance.AddMockMethod(targetFunc, func() *mock.InstanceMock {
				mt := h.ManagedTypes()
				mapHandle := mt.NewManagedMap()
				keyHandle := mt.NewManagedBufferFromBytes([]byte("key"))
				valueHandle := mt.NewManagedBufferFromBytes([]byte("value"))

				vmh := vmhooks.NewVMHooksImpl(h)
				vmh.ManagedMapPut(mapHandle, keyHandle, valueHandle) // Attempt to put into map
				return instance
			})
		}, initialGas)

	defer host.Reset()

	vmOutput, err := host.RunSmartContractCall(callInput)
	require.Nil(t, err)
	assert.NotEqual(t, vmcommon.Ok, vmOutput.ReturnCode, "Execution should have failed")
	assert.Contains(t, vmOutput.ReturnMessage, vmhost.ErrInvalidCallOnReadOnlyMode.Error())
}

func TestManagedExecuteReadOnly_AttemptAsyncCall_ShouldFail(t *testing.T) {
	targetAddr := testcommon.MakeSimpleAddress(2)
	targetFunc := "attemptAsyncCall"
	initialGas := uint64(5_000_000)

	host, callInput := setupReadOnlyTest(t, targetAddr, targetFunc,
		func(h vmhost.VMHost, instance *mock.InstanceMock) {
			instance.AddMockMethod(targetFunc, func() *mock.InstanceMock {
				mt := h.ManagedTypes()
				destHandle := mt.NewManagedBufferFromBytes(testcommon.MakeSimpleAddress(3))
				valueHandle := mt.NewBigInt(big.NewInt(0))
				funcHandle := mt.NewManagedBufferFromBytes([]byte("someFunc"))
				argsHandle := mt.NewManagedBuffer()

				vmh := vmhooks.NewVMHooksImpl(h)
				// ManagedAsyncCall(destHandle int32, valueHandle int32, functionHandle int32, argumentsHandle int32)
				vmh.ManagedAsyncCall(destHandle, valueHandle, funcHandle, argsHandle)
				return instance
			})
		}, initialGas)

	defer host.Reset()

	vmOutput, err := host.RunSmartContractCall(callInput)
	require.Nil(t, err)
	assert.NotEqual(t, vmcommon.Ok, vmOutput.ReturnCode, "Execution should have failed")
	assert.Contains(t, vmOutput.ReturnMessage, vmhost.ErrInvalidCallOnReadOnlyMode.Error())
}

func TestManagedExecuteReadOnly_AttemptCreateAsyncCall_ShouldFail(t *testing.T) {
	targetAddr := testcommon.MakeSimpleAddress(2)
	targetFunc := "attemptCreateAsyncCall"
	initialGas := uint64(5_000_000)

	host, callInput := setupReadOnlyTest(t, targetAddr, targetFunc,
		func(h vmhost.VMHost, instance *mock.InstanceMock) {
			instance.AddMockMethod(targetFunc, func() *mock.InstanceMock {
				mt := h.ManagedTypes()
				destHandle := mt.NewManagedBufferFromBytes(testcommon.MakeSimpleAddress(3))
				valueHandle := mt.NewBigInt(big.NewInt(0))
				funcHandle := mt.NewManagedBufferFromBytes([]byte("someFunc"))
				argsHandle := mt.NewManagedBuffer()
				closureHandle := mt.NewManagedBuffer() // Empty closure

				vmh := vmhooks.NewVMHooksImpl(h)
				// ManagedCreateAsyncCall(destHandle, valueHandle, funcHandle, argumentsHandle,
				// successOffset, successLength, errorOffset, errorLength, gas, extraGasForCallback, callbackClosureHandle)
				// For simplicity, using 0 for offsets/lengths for callbacks as it should fail before that.
				vmh.ManagedCreateAsyncCall(destHandle, valueHandle, funcHandle, argsHandle,
					0, 0, 0, 0, 100000, 0, closureHandle)
				return instance
			})
		}, initialGas)

	defer host.Reset()

	vmOutput, err := host.RunSmartContractCall(callInput)
	require.Nil(t, err)
	assert.NotEqual(t, vmcommon.Ok, vmOutput.ReturnCode, "Execution should have failed")
	assert.Contains(t, vmOutput.ReturnMessage, vmhost.ErrInvalidCallOnReadOnlyMode.Error())
}

// TODO: Add similar tests for baseOps.AsyncCall and baseOps.CreateAsyncCall if they can be invoked
// from a ManagedExecuteReadOnly context through a target contract that uses these base EEIs.
// This would require the target contract to be written in a language that compiles to WASM and uses these.
// For now, testing the Managed* versions covers the direct EEI access from Rust/C++.
