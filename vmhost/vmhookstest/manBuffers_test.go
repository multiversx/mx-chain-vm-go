package vmhookstest

import (
	"io"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
	arwenMath "github.com/multiversx/mx-chain-vm-go/math"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"

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

func buildRandomizer(host arwen.VMHost) io.Reader {
	// building the randomizer
	blockchainContext := host.Blockchain()
	previousRandomSeed := blockchainContext.LastRandomSeed()
	currentRandomSeed := blockchainContext.CurrentRandomSeed()
	txHash := host.Runtime().GetCurrentTxHash()

	blocksRandomSeed := append(previousRandomSeed, currentRandomSeed...)
	randomSeed := append(blocksRandomSeed, txHash...)
	randReader := arwenMath.NewSeedRandReader(randomSeed)
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
			WithGasProvided(100000).
			WithFunction("mBufferFromBigIntUnsignedTest").
			WithArguments([]byte{byte(numberOfReps)}).
			Build()).
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
