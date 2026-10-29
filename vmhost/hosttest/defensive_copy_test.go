package hostCoretest

import (
	"encoding/binary"
	"math/big"
	"testing"

	vm "github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/require"
)

func TestDefensiveCopy_StorageUpdateKeyMutation(t *testing.T) {
	parentBalance := int64(1000000)
	gasProvided := uint64(1000000)

	originalKey := []byte("original_storage_key")
	mutatedKey := []byte("mutated_storage_keys")
	originalVal := []byte("original_storage_val")

	var capturedOutput *vmcommon.VMOutput

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testStorage", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						keyBuf := managedTypes.NewManagedBufferFromBytes(originalKey)
						valBuf := managedTypes.NewManagedBufferFromBytes(originalVal)

						res, _ := handleHookErr(host)(hooks.MBufferStorageStore(keyBuf, valBuf))
						require.Equal(t, int32(0), res)

						// Mutate keyBuf in place
						handleHookErr(host)(vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, keyBuf, 0, int32(len(mutatedKey)), mutatedKey))

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("testStorage").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			capturedOutput = verify.VmOutput
		})

	require.Nil(t, err)
	require.NotNil(t, capturedOutput)

	parentAcc, hasAcc := capturedOutput.OutputAccounts[string(test.ParentAddress)]
	require.True(t, hasAcc)

	update, hasUpdate := parentAcc.StorageUpdates[string(originalKey)]
	require.True(t, hasUpdate, "Original storage key must be present in StorageUpdates")
	require.Equal(t, originalKey, update.Offset, "StorageUpdate.Offset must NOT be mutated by subsequent key buffer changes")
	require.NotEqual(t, mutatedKey, update.Offset)
}

func TestDefensiveCopy_WriteLogTopicAndDataMutation(t *testing.T) {
	parentBalance := int64(1000000)
	gasProvided := uint64(1000000)

	originalTopic := []byte("topic_original_001")
	mutatedTopic := []byte("topic_mutated__002")
	originalData := []byte("data_original_value")
	mutatedData := []byte("data_mutated__value")

	var capturedOutput *vmcommon.VMOutput

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testLog", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						topicBuf := managedTypes.NewManagedBufferFromBytes(originalTopic)
						dataBuf := managedTypes.NewManagedBufferFromBytes(originalData)

						// Create managed vector with topicBuf
						vecBuf := managedTypes.NewManagedBuffer()
						handleBytes := make([]byte, 4)
						binary.BigEndian.PutUint32(handleBytes, uint32(topicBuf))
						managedTypes.SetBytes(vecBuf, handleBytes)

						handleHookError(host, hooks.ManagedWriteLog(vecBuf, dataBuf))

						// Mutate topicBuf and dataBuf in place
						handleHookErr(host)(vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, topicBuf, 0, int32(len(mutatedTopic)), mutatedTopic))
						handleHookErr(host)(vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, dataBuf, 0, int32(len(mutatedData)), mutatedData))

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("testLog").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			capturedOutput = verify.VmOutput
		})

	require.Nil(t, err)
	require.NotNil(t, capturedOutput)
	require.NotEmpty(t, capturedOutput.Logs)

	lastLog := capturedOutput.Logs[len(capturedOutput.Logs)-1]
	require.NotEmpty(t, lastLog.Topics)
	require.Equal(t, originalTopic, lastLog.Topics[0], "Log topic must retain original value and not be mutated")
	require.NotEqual(t, mutatedTopic, lastLog.Topics[0])

	require.NotEmpty(t, lastLog.Data)
	require.Equal(t, originalData, lastLog.Data[0], "Log data must retain original value and not be mutated")
	require.NotEqual(t, mutatedData, lastLog.Data[0])
}

func TestDefensiveCopy_FinishReturnDataMutation(t *testing.T) {
	parentBalance := int64(1000000)
	gasProvided := uint64(1000000)

	originalRet := []byte("original_return_data_1234")
	mutatedRet := []byte("mutated_return__data_5678")

	var capturedOutput *vmcommon.VMOutput

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testFinish", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						retBuf := managedTypes.NewManagedBufferFromBytes(originalRet)
						res, _ := handleHookErr(host)(hooks.MBufferFinish(retBuf))
						require.Equal(t, int32(0), res)

						// Mutate retBuf in place
						handleHookErr(host)(vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, retBuf, 0, int32(len(mutatedRet)), mutatedRet))

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("testFinish").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			capturedOutput = verify.VmOutput
		})

	require.Nil(t, err)
	require.NotNil(t, capturedOutput)
	require.NotEmpty(t, capturedOutput.ReturnData)

	require.Equal(t, originalRet, capturedOutput.ReturnData[0], "ReturnData must retain original value and not be mutated")
	require.NotEqual(t, mutatedRet, capturedOutput.ReturnData[0])
}

func TestDefensiveCopy_TransferMutation(t *testing.T) {
	parentBalance := int64(1000000)
	gasProvided := uint64(5000000)

	recipient := []byte("recipient_address_12345678901234")
	sender := make([]byte, len(test.ParentAddress))
	copy(sender, test.ParentAddress)

	originalInput := []byte("transfer_call_data_input")
	originalAsync := []byte("async_data_payload_info")

	mutableInput := make([]byte, len(originalInput))
	copy(mutableInput, originalInput)

	mutableAsync := make([]byte, len(originalAsync))
	copy(mutableAsync, originalAsync)

	mutableSender := make([]byte, len(sender))
	copy(mutableSender, sender)

	var capturedOutput *vmcommon.VMOutput

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testTransfer", func() *contextmock.InstanceMock {
						host := parentInstance.Host

						err := host.Output().Transfer(
							recipient,
							mutableSender,
							0,
							0,
							big.NewInt(100),
							mutableAsync,
							mutableInput,
							vm.DirectCall,
						)
						require.Nil(t, err)

						// Now mutate the slices in-place
						mutableInput[0] = 'X'
						mutableAsync[0] = 'Y'
						mutableSender[0] = 'Z'

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("testTransfer").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			capturedOutput = verify.VmOutput
		})

	require.Nil(t, err)
	require.NotNil(t, capturedOutput)

	recAcc, hasRec := capturedOutput.OutputAccounts[string(recipient)]
	require.True(t, hasRec)
	require.NotEmpty(t, recAcc.OutputTransfers)

	outTransfer := recAcc.OutputTransfers[0]
	require.Equal(t, originalInput, outTransfer.Data, "OutputTransfer.Data must not be mutated")
	require.Equal(t, originalAsync, outTransfer.AsyncData, "OutputTransfer.AsyncData must not be mutated")
	require.Equal(t, sender, outTransfer.SenderAddress, "OutputTransfer.SenderAddress must not be mutated")
}

func TestDefensiveCopy_ReadManagedVec(t *testing.T) {
	parentBalance := int64(1000000)
	gasProvided := uint64(1000000)

	originalData1 := []byte("buffer_data_1")
	originalData2 := []byte("buffer_data_2")
	mutatedData := []byte("mutated_buf__")

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testReadVec", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						managedTypes := host.ManagedTypes()

						buf1 := managedTypes.NewManagedBufferFromBytes(originalData1)
						buf2 := managedTypes.NewManagedBufferFromBytes(originalData2)

						vecBuf := managedTypes.NewManagedBuffer()
						handleBytes := make([]byte, 8)
						binary.BigEndian.PutUint32(handleBytes[0:4], uint32(buf1))
						binary.BigEndian.PutUint32(handleBytes[4:8], uint32(buf2))
						managedTypes.SetBytes(vecBuf, handleBytes)

						vecSlices, _, err := managedTypes.ReadManagedVecOfManagedBuffers(vecBuf)
						require.Nil(t, err)
						require.Len(t, vecSlices, 2)
						require.Equal(t, originalData1, vecSlices[0])
						require.Equal(t, originalData2, vecSlices[1])

						// Mutate buf1 in place
						handleHookErr(host)(vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, buf1, 0, int32(len(mutatedData)), mutatedData))

						// Verify vecSlices[0] was NOT affected
						require.Equal(t, originalData1, vecSlices[0], "ReadManagedVecOfManagedBuffers returned slice must not be mutated")
						require.NotEqual(t, mutatedData, vecSlices[0])

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("testReadVec").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

	require.Nil(t, err)
}
