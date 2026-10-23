package vmhookstest

import (
	"io"
	"math"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmMath "github.com/multiversx/mx-chain-vm-go/math"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/assert"

	twoscomplement "github.com/multiversx/mx-components-big-int/twos-complement"
)

var mBufferKey = []byte("mBuffer")
var managedBuffer = []byte{0xff, 0x2a, 0x26, 0x5f, 0x8b, 0xcb, 0xdc, 0xaf,
	0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24,
	0xcb, 0x40, 0xd6, 0x4a, 0x50, 0x1f, 0xba, 0x9c,
	0x11, 0x84, 0x7b, 0x28, 0x96, 0x5b, 0xc7, 0x37,
	0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24,
	0xcb, 0x40, 0xd6, 0x4a, 0x50, 0x1f, 0xba, 0x9c,
	0x11, 0x84, 0x7b, 0x28, 0x96, 0x5b, 0xc7, 0x37,
	0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24}
var numberOfReps = 100
var lengthOfBuffer = 64

func buildRandomizer(host vmhost.VMHost) io.Reader {
	// building the randomizer
	blockchainContext := host.Blockchain()
	previousRandomSeed := blockchainContext.LastRandomSeed()
	currentRandomSeed := blockchainContext.CurrentRandomSeed()
	txHash := host.Runtime().GetCurrentTxHash()

	blocksRandomSeed := append(previousRandomSeed, currentRandomSeed...)
	randomSeed := append(blocksRandomSeed, txHash...)
	randReader := vmMath.NewSeedRandReader(randomSeed)
	return randReader
}

func TestManBuffers_MixedFunctions(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferMethod").
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			expectedStorageEntry := test.CreateStoreEntry(test.ParentAddress).WithKey(mBufferKey).WithValue(managedBuffer)
			verify.Ok().
				ReturnData(managedBuffer, []byte("succ")).
				Storage(expectedStorageEntry)
		})
}

func TestManBuffers_New(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferNewTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData([]byte{byte(numberOfReps - 1)})
		})
}

func TestManBuffers_NewFromBytes(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferNewFromBytesTest").
			WithArguments([]byte{byte(numberOfReps)}, []byte{byte(lengthOfBuffer)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData(managedBuffer)
		})
}

func TestManBuffers_SetRandom(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferSetRandomTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			randomBuffer := make([]byte, numberOfReps)
			for i := 0; i < numberOfReps; i++ {
				_, _ = randReader.Read(randomBuffer)
			}
			verify.Ok().
				ReturnData(randomBuffer)
		})
}

func TestManBuffers_GetLength(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferGetLengthTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData([]byte{byte(numberOfReps)})
		})
}

func TestManBuffers_GetBytes(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferGetBytesTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			randomBuffer := make([]byte, numberOfReps)
			for i := 0; i < numberOfReps; i++ {
				_, _ = randReader.Read(randomBuffer)
			}
			verify.Ok().
				ReturnData(randomBuffer, randomBuffer[:numberOfReps])
		})
}

func TestManBuffers_AppendBytes(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferAppendTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			finalBuffer := make([]byte, 0)
			randomBuffer := make([]byte, numberOfReps)
			for i := 0; i < numberOfReps; i++ {
				_, _ = randReader.Read(randomBuffer)
				finalBuffer = append(finalBuffer, randomBuffer...)
			}
			verify.Ok().
				ReturnData(finalBuffer)
		})
}

func TestManBuffers_mBufferToBigIntUnsigned(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferToBigIntUnsignedTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			randomBuffer := make([]byte, numberOfReps)
			for i := 0; i < numberOfReps; i++ {
				_, _ = randReader.Read(randomBuffer)
			}
			verify.Ok().
				ReturnData(randomBuffer, randomBuffer)
		})
}

func TestManBuffers_mBufferToBigIntSigned(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferToBigIntSignedTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			randomBuffer := make([]byte, numberOfReps)
			for i := 0; i < numberOfReps; i++ {
				_, _ = randReader.Read(randomBuffer)
			}
			expectedBuffer := twoscomplement.ToBytes(big.NewInt(0).SetBytes(randomBuffer))[1:]
			verify.Ok().
				ReturnData(randomBuffer, expectedBuffer)
		})
}

func TestManBuffers_mBufferFromBigIntUnsigned(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(1000000).
			WithFunction("mBufferFromBigIntUnsignedTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			randomBuffer := make([]byte, numberOfReps)
			for i := 0; i < numberOfReps; i++ {
				_, _ = randReader.Read(randomBuffer)
			}
			verify.Ok().
				ReturnData(randomBuffer, randomBuffer)
		})
}

func TestManBuffers_mBufferFromBigIntSigned(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferFromBigIntSignedTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			randomBuffer := make([]byte, numberOfReps)
			for i := 0; i < numberOfReps; i++ {
				_, _ = randReader.Read(randomBuffer)
			}
			expectedBuffer := twoscomplement.ToBytes(big.NewInt(0).SetBytes(randomBuffer))[1:]
			verify.Ok().
				ReturnData(randomBuffer, expectedBuffer)
		})
}

func TestManBuffers_StorageStore(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferStorageStoreTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			lastRandomBuffer := make([]byte, numberOfReps)
			lastRandomKey := make([]byte, 5)
			storage := make([]test.StoreEntry, 0)
			for i := 0; i < numberOfReps; i++ {
				keyBuffer := make([]byte, 5)
				randomBuffer := make([]byte, numberOfReps)
				_, _ = randReader.Read(keyBuffer)
				_, _ = randReader.Read(randomBuffer)
				entry := test.CreateStoreEntry(test.ParentAddress).WithKey(keyBuffer).WithValue(randomBuffer)
				storage = append(storage, entry)
				if i == numberOfReps-1 {
					lastRandomBuffer = randomBuffer
					lastRandomKey = keyBuffer
				}
			}
			verify.Ok().
				ReturnData(lastRandomBuffer, lastRandomKey).
				Storage(storage...)
		})
}

func TestManBuffers_StorageLoad(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("managed-buffers", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("mBufferStorageLoadTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			randReader := buildRandomizer(host)

			lastRandomBuffer := make([]byte, numberOfReps)
			lastRandomKey := make([]byte, 5)
			storage := make([]test.StoreEntry, 0)
			for i := 0; i < numberOfReps; i++ {
				keyBuffer := make([]byte, 5)
				randomBuffer := make([]byte, numberOfReps)
				_, _ = randReader.Read(keyBuffer)
				_, _ = randReader.Read(randomBuffer)
				entry := test.CreateStoreEntry(test.ParentAddress).WithKey(keyBuffer).WithValue(randomBuffer)
				storage = append(storage, entry)
				if i == numberOfReps-1 {
					lastRandomBuffer = randomBuffer
					lastRandomKey = keyBuffer
				}
			}
			verify.Ok().
				ReturnData(lastRandomBuffer, lastRandomKey).
				Storage(storage...)
		})
}

func TestManBuffers_Ints_WithNoBarnardOpcodesFlag(t *testing.T) {
	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						instance := contextmock.GetMockInstance(host)
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						sourceBuffer := managedTypes.NewManagedBufferFromBytes([]byte{0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff})

						bigIntDest := managedTypes.NewBigIntFromInt64(0)
						result := hooks.MBufferToBigIntUnsigned(sourceBuffer, bigIntDest)
						assert.True(t, result == 0)

						bigIntResult, err := managedTypes.GetBigInt(bigIntDest)
						assert.Nil(t, err)
						assert.True(t, bigIntResult.Cmp(big.NewInt(0x00_ff_00_ff_00_ff_00_ff)) == 0)

						result = hooks.MBufferToBigIntSigned(sourceBuffer, bigIntDest)
						assert.True(t, result == 0)

						bigIntResult, err = managedTypes.GetBigInt(bigIntDest)
						assert.Nil(t, err)
						assert.True(t, bigIntResult.Cmp(big.NewInt(0x00_ff_00_ff_00_ff_00_ff)) == 0)

						return instance
					})
				})).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("test").
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
			enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
				return flag != vmhost.BarnardOpcodesFlag
			}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(gasProvided-35).
				GasUsed(test.ParentAddress, 35)
		})
	assert.Nil(t, err)
}

func TestManBuffers_Ints_WithBarnardOpcodesFlag(t *testing.T) {
	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						instance := contextmock.GetMockInstance(host)
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						sourceBuffer := managedTypes.NewManagedBufferFromBytes([]byte{0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff})

						bigIntDest := managedTypes.NewBigIntFromInt64(0)
						result := hooks.MBufferToBigIntUnsigned(sourceBuffer, bigIntDest)
						assert.True(t, result == 0)

						bigIntResult, err := managedTypes.GetBigInt(bigIntDest)
						assert.Nil(t, err)
						assert.True(t, bigIntResult.Cmp(big.NewInt(0x00_ff_00_ff_00_ff_00_ff)) == 0)

						result = hooks.MBufferToBigIntSigned(sourceBuffer, bigIntDest)
						assert.True(t, result == 0)

						bigIntResult, err = managedTypes.GetBigInt(bigIntDest)
						assert.Nil(t, err)
						assert.True(t, bigIntResult.Cmp(big.NewInt(0x00_ff_00_ff_00_ff_00_ff)) == 0)

						smallInt := hooks.MBufferToSmallIntSigned(sourceBuffer)
						assert.True(t, smallInt == 0x00_ff_00_ff_00_ff_00_ff)

						smallInt = hooks.MBufferToSmallIntUnsigned(sourceBuffer)
						assert.True(t, smallInt == 0x00_ff_00_ff_00_ff_00_ff)

						return instance
					})
				})).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("test").
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
			enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
				return true
			}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(gasProvided-69).
				GasUsed(test.ParentAddress, 69)
		})
	assert.Nil(t, err)
}

func TestManBuffers_Int32OverflowBoundsCheck(t *testing.T) {
	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						sourceBuffer := managedTypes.NewManagedBufferFromBytes([]byte{0x01, 0x02, 0x03})
						destBuffer := managedTypes.NewManagedBufferFromBytes([]byte{0x00, 0x00, 0x00})

						// startingPosition + sliceLength = MaxInt32 + 1 overflows int32 to negative
						res := hooks.MBufferGetByteSlice(sourceBuffer, math.MaxInt32, 1, 0)
						assert.Equal(t, int32(1), res)

						res = hooks.MBufferCopyByteSlice(sourceBuffer, math.MaxInt32, 1, destBuffer)
						assert.Equal(t, int32(1), res)

						res = vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, destBuffer, math.MaxInt32, 1, []byte{0x99})
						assert.Equal(t, int32(1), res)

						// negative starting position
						res = hooks.MBufferGetByteSlice(sourceBuffer, -1, 1, 0)
						assert.Equal(t, int32(1), res)

						res = hooks.MBufferCopyByteSlice(sourceBuffer, -1, 1, destBuffer)
						assert.Equal(t, int32(1), res)

						res = vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, destBuffer, -1, 1, []byte{0x99})
						assert.Equal(t, int32(1), res)

						// negative slice length
						res = hooks.MBufferGetByteSlice(sourceBuffer, 0, -1, 0)
						assert.Equal(t, int32(1), res)

						res = hooks.MBufferCopyByteSlice(sourceBuffer, 0, -1, destBuffer)
						assert.Equal(t, int32(1), res)

						res = vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, destBuffer, 0, -1, []byte{0x99})
						assert.Equal(t, int32(1), res)

						return contextmock.GetMockInstance(host)
					})
				})).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("test").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
		})
	assert.Nil(t, err)
}

func TestManBuffers_ManagedBufferSetByteSlice_CopyOnWrite(t *testing.T) {
	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()

						initialData := []byte{0xAA, 0xBB, 0xCC, 0xDD}
						bufHandle := managedTypes.NewManagedBufferFromBytes(initialData)

						// Get reference slice
						refSlice, err := managedTypes.GetBytes(bufHandle)
						assert.Nil(t, err)
						assert.Equal(t, []byte{0xAA, 0xBB, 0xCC, 0xDD}, refSlice)

						// Overwrite part of the buffer via SetByteSlice
						res := vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, bufHandle, 1, 2, []byte{0x11, 0x22})
						assert.Equal(t, int32(0), res)

						// Live buffer should have new content
						newSlice, err := managedTypes.GetBytes(bufHandle)
						assert.Nil(t, err)
						assert.Equal(t, []byte{0xAA, 0x11, 0x22, 0xDD}, newSlice)

						// Prior reference MUST NOT be mutated (Copy-on-Write isolation)
						assert.Equal(t, []byte{0xAA, 0xBB, 0xCC, 0xDD}, refSlice)

						return contextmock.GetMockInstance(host)
					})
				})).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("test").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
		})
	assert.Nil(t, err)
}

func TestManBuffers_MBufferStorageLoad_InvalidHandleReturnsOne(t *testing.T) {
	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						destHandle := managedTypes.NewManagedBuffer()
						res := hooks.MBufferStorageLoad(-999, destHandle)
						assert.Equal(t, int32(1), res)

						return contextmock.GetMockInstance(host)
					})
				})).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("test").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().ReturnMessage(vmhost.ErrNoManagedBufferUnderThisHandle.Error())
		})
	assert.Nil(t, err)
}

func TestManBuffers_MBufferStorageLoad_ValidKey(t *testing.T) {
	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						destHandle := managedTypes.NewManagedBuffer()
						keyHandle := managedTypes.NewManagedBufferFromBytes([]byte("test_key"))
						res := hooks.MBufferStorageLoad(keyHandle, destHandle)
						assert.Equal(t, int32(0), res)

						destBytes, err := managedTypes.GetBytes(destHandle)
						assert.Nil(t, err)
						assert.Equal(t, []byte{}, destBytes)

						valHandle := managedTypes.NewManagedBufferFromBytes([]byte("test_value"))
						resStore := hooks.MBufferStorageStore(keyHandle, valHandle)
						assert.Equal(t, int32(0), resStore)

						res = hooks.MBufferStorageLoad(keyHandle, destHandle)
						assert.Equal(t, int32(0), res)

						destBytes, err = managedTypes.GetBytes(destHandle)
						assert.Nil(t, err)
						assert.Equal(t, []byte("test_value"), destBytes)

						return contextmock.GetMockInstance(host)
					})
				})).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("test").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
		})
	assert.Nil(t, err)
}

func TestManBuffers_MBufferFromSmallIntSigned_Negatives(t *testing.T) {
	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	testCases := []struct {
		name        string
		flagEnabled bool
		input       int64
		expectedOut int64
	}{
		{
			name:        "flag enabled: negative 5 preserves sign",
			flagEnabled: true,
			input:       -5,
			expectedOut: -5,
		},
		{
			name:        "flag enabled: negative 1 preserves sign",
			flagEnabled: true,
			input:       -1,
			expectedOut: -1,
		},
		{
			name:        "flag enabled: negative 123456 preserves sign",
			flagEnabled: true,
			input:       -123456,
			expectedOut: -123456,
		},
		{
			name:        "flag disabled: negative 5 returns positive magnitude 5 (old behavior)",
			flagEnabled: false,
			input:       -5,
			expectedOut: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(parentBalance).
						WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
							instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
								host := instanceMock.Host
								managedTypes := host.ManagedTypes()
								hooks := vmhooks.NewVMHooksImpl(host)

								bufHandle := managedTypes.NewManagedBufferFromBytes([]byte{})
								hooks.MBufferFromSmallIntSigned(bufHandle, tc.input)

								readBack := hooks.MBufferToSmallIntSigned(bufHandle)
								assert.Equal(t, tc.expectedOut, readBack)

								return contextmock.GetMockInstance(host)
							})
						})).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(gasProvided).
					WithFunction("test").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return tc.flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
				})
			assert.Nil(t, err)
		})
	}
}

func TestManBuffers_MBufferNewFromBytes_GasUsage(t *testing.T) {
	t.Parallel()

	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	testCases := []struct {
		name        string
		flagEnabled bool
	}{
		{
			name:        "flag enabled: charges base cost + per-byte cost",
			flagEnabled: true,
		},
		{
			name:        "flag disabled: charges only base cost",
			flagEnabled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(parentBalance).
						WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
							instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
								host := instanceMock.Host
								hooks := vmhooks.NewVMHooksImpl(host)

								testBytes := []byte("01234567890123456789") // 20 bytes
								err := instanceMock.MemStore(0, testBytes)
								assert.Nil(t, err)

								baseCost := host.Metering().GasSchedule().ManagedBufferAPICost.MBufferNewFromBytes
								copyPerByte := host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte

								gasBefore := host.Metering().GasLeft()
								handle := hooks.MBufferNewFromBytes(0, 20)
								assert.True(t, handle >= 0)
								gasUsed := gasBefore - host.Metering().GasLeft()

								if tc.flagEnabled {
									assert.Equal(t, baseCost+20*copyPerByte, gasUsed)
								} else {
									assert.Equal(t, baseCost, gasUsed)
								}

								return contextmock.GetMockInstance(host)
							})
						})).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(gasProvided).
					WithFunction("test").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return tc.flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
				})
			assert.Nil(t, err)
		})
	}
}

func TestManBuffers_ManagedBufferToHex_GasUsage(t *testing.T) {
	t.Parallel()

	parentBalance := int64(100000)
	gasProvided := uint64(100000)

	testCases := []struct {
		name        string
		flagEnabled bool
	}{
		{
			name:        "flag enabled: charges base cost + per-byte cost for hex encoded bytes",
			flagEnabled: true,
		},
		{
			name:        "flag disabled: charges only base cost",
			flagEnabled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(parentBalance).
						WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
							instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
								host := instanceMock.Host
								managedTypes := host.ManagedTypes()
								hooks := vmhooks.NewVMHooksImpl(host)

								testBytes := []byte("0123456789") // 10 bytes -> hex is 20 bytes
								srcHandle := managedTypes.NewManagedBufferFromBytes(testBytes)
								destHandle := managedTypes.NewManagedBuffer()

								baseCost := host.Metering().GasSchedule().ManagedBufferAPICost.MBufferSetBytes
								copyPerByte := host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte

								gasBefore := host.Metering().GasLeft()
								hooks.ManagedBufferToHex(srcHandle, destHandle)
								gasUsed := gasBefore - host.Metering().GasLeft()

								if tc.flagEnabled {
									assert.Equal(t, baseCost+20*copyPerByte, gasUsed)
								} else {
									assert.Equal(t, baseCost, gasUsed)
								}

								destBytes, err := managedTypes.GetBytes(destHandle)
								assert.Nil(t, err)
								assert.Equal(t, "30313233343536373839", string(destBytes))

								return contextmock.GetMockInstance(host)
							})
						})).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(gasProvided).
					WithFunction("test").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return tc.flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
				})
			assert.Nil(t, err)
		})
	}
}
