package hostCoretest

import (
	"math/big"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseOpsAPI_CallValue(t *testing.T) {
	code := testcommon.GetTestSCCode("baseOps", "../../")

	// 1-byte call value
	host := testcommon.NewTestHostBuilder(t).
		WithBlockchainHook(testcommon.BlockchainHookStubForCall(code, nil)).
		Build()
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
	host = testcommon.NewTestHostBuilder(t).
		WithBlockchainHook(testcommon.BlockchainHookStubForCall(code, nil)).
		Build()
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
	host = testcommon.NewTestHostBuilder(t).
		WithBlockchainHook(testcommon.BlockchainHookStubForCall(code, nil)).
		Build()
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
	host := testcommon.NewTestHostBuilder(t).
		WithBlockchainHook(testcommon.BlockchainHookStubForCall(code, nil)).
		Build()
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

// setupReadOnlyTestForBaseOps creates a common setup for testing base EEIs from a read-only context.
// ContractA (caller, using managed EEIs) calls ManagedExecuteReadOnly to execute a function in ContractB.
// ContractB's function (target, using base EEIs) then attempts a restricted operation.
func setupReadOnlyTestForBaseOps(
	t *testing.T,
	targetContractAddress []byte,
	targetFunctionName string,
	targetSetup func(host vmhost.VMHost, instance *testcommon.InstanceMock), // To set up ContractB's method using base EEIs
	initialGas uint64,
) (vmhost.VMHost, *vmcommon.ContractCallInput) {
	host := testcommon.NewTestHostBuilder(t).Build()

	// ContractA (Caller) - This will use ManagedExecuteReadOnly
	callerContractVM := testcommon.CreateMockContract(testcommon.CallerAddress)
	callerContractVM.SetMethods(func(instance *testcommon.InstanceMock, cfg interface{}) {
		instance.AddMockMethod("callTargetReadOnly_base", func() *testcommon.InstanceMock {
			h := instance.Host
			mt := h.ManagedTypes()

			targetAddrHandle := mt.NewManagedBufferFromBytes(targetContractAddress)
			targetFuncHandle := mt.NewManagedBufferFromBytes([]byte(targetFunctionName))
			argsHandle := mt.NewManagedBuffer()
			resultBufHandle := mt.NewManagedBuffer()

			vmh_managed := vmhooks.NewVMHooksImpl(h) // Managed EEI for ReadOnly call
			ret := vmh_managed.ManagedExecuteReadOnly(
				initialGas/2,
				targetAddrHandle,
				targetFuncHandle,
				argsHandle,
				resultBufHandle,
			)
			if ret != 0 && h.GetRuntimeErrors() == nil {
				h.Runtime().SignalUserError("ManagedExecuteReadOnly failed unexpectedly at EEI level but no runtime error")
			}
			return instance
		})
	})

	// ContractB (Target) - This will attempt base EEI calls
	targetContractVM := testcommon.CreateMockContract(targetContractAddress)
	targetContractVM.SetMethods(func(instance *testcommon.InstanceMock, cfg interface{}) {
		// The targetSetup will add mock methods to this instance that use base EEIs
		targetSetup(instance.Host, instance)
	})

	host.ContextsHolder().AddContract(callerContractVM)
	host.ContextsHolder().AddContract(targetContractVM)

	callInput := testcommon.CreateTestContractCallInputBuilder().
		WithCallerAddr(testcommon.UserAccountAddress).       // EOA calling ContractA
		WithRecipientAddr(testcommon.CallerAddress).       // ContractA
		WithGasProvided(initialGas).
		WithFunction("callTargetReadOnly_base").
		Build()

	return host, callInput
}

func TestBaseOps_AsyncCall_FromReadOnlyContext_ShouldFail(t *testing.T) {
	targetAddr := testcommon.MakeSimpleAddress(2)
	targetFunc := "attemptBaseAsyncCall"
	initialGas := uint64(5_000_000)

	host, callInput := setupReadOnlyTestForBaseOps(t, targetAddr, targetFunc,
		func(h vmhost.VMHost, instance *testcommon.InstanceMock) {
			instance.AddMockMethod(targetFunc, func() *testcommon.InstanceMock {
				// Simulate calling the base AsyncCall EEI
				// For base EEIs, arguments are direct memory pointers
				// We need to put data into WASM memory for dest, value, data
				destOffset := int32(0)
				valueOffset := int32(32) // Assuming address is 32 bytes
				dataOffset := int32(64)  // Assuming value is 32 bytes

				err := instance.MemStore(uint32(destOffset), testcommon.MakeSimpleAddress(3))
				require.Nil(t, err)
				err = instance.MemStore(uint32(valueOffset), big.NewInt(0).Bytes()) // Zero value
				require.Nil(t, err)
				err = instance.MemStore(uint32(dataOffset), []byte("test"))
				require.Nil(t, err)

				vmh_base := vmhooks.NewVMHooksImpl(h)
				vmh_base.AsyncCall(
					uint32(destOffset),
					uint32(valueOffset),
					uint32(dataOffset),
					int32(len("test")),
				)
				return instance
			})
		}, initialGas)

	defer host.Reset()

	vmOutput, err := host.RunSmartContractCall(callInput)
	require.Nil(t, err)
	assert.NotEqual(t, vmcommon.Ok, vmOutput.ReturnCode, "Execution should have failed")
	assert.Contains(t, vmOutput.ReturnMessage, vmhost.ErrInvalidCallOnReadOnlyMode.Error())
}

func TestBaseOps_CreateAsyncCall_FromReadOnlyContext_ShouldFail(t *testing.T) {
	targetAddr := testcommon.MakeSimpleAddress(2)
	targetFunc := "attemptBaseCreateAsyncCall"
	initialGas := uint64(5_000_000)

	host, callInput := setupReadOnlyTestForBaseOps(t, targetAddr, targetFunc,
		func(h vmhost.VMHost, instance *testcommon.InstanceMock) {
			instance.AddMockMethod(targetFunc, func() *testcommon.InstanceMock {
				// Simulate calling the base CreateAsyncCall EEI
				destOffset := int32(0)
				valueOffset := int32(32)
				dataOffset := int32(64)
				// For callbacks, assuming empty strings, so 0 length and offset can be same as data for simplicity
				// as it should fail before trying to resolve these.
				cbSuccessOffset := dataOffset + int32(len("test_data"))
				cbErrorOffset := cbSuccessOffset

				err := instance.MemStore(uint32(destOffset), testcommon.MakeSimpleAddress(3))
				require.Nil(t, err)
				err = instance.MemStore(uint32(valueOffset), big.NewInt(0).Bytes())
				require.Nil(t, err)
				err = instance.MemStore(uint32(dataOffset), []byte("test_data"))
				require.Nil(t, err)
				// Store empty strings for callback names
				err = instance.MemStore(uint32(cbSuccessOffset), []byte(""))
				require.Nil(t, err)


				vmh_base := vmhooks.NewVMHooksImpl(h)
				ret := vmh_base.CreateAsyncCall(
					uint32(destOffset),
					uint32(valueOffset),
					uint32(dataOffset), int32(len("test_data")),
					uint32(cbSuccessOffset), 0, // success callback (empty)
					uint32(cbErrorOffset), 0,   // error callback (empty)
					100000, 0,                  // gas, extraGas
				)
				// ret value for CreateAsyncCall is 0 on success, 1 on error in this context
				// but the actual error is set on the host by FailExecution
				return instance
			})
		}, initialGas)

	defer host.Reset()

	vmOutput, err := host.RunSmartContractCall(callInput)
	require.Nil(t, err)
	assert.NotEqual(t, vmcommon.Ok, vmOutput.ReturnCode, "Execution should have failed")
	assert.Contains(t, vmOutput.ReturnMessage, vmhost.ErrInvalidCallOnReadOnlyMode.Error())
}
