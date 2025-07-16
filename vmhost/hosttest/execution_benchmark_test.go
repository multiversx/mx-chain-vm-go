package hostCoretest

import (
	"math"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	gasSchedules "github.com/multiversx/mx-chain-vm-go/scenario/gasSchedules"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Address is a type alias for []byte
type Address = []byte

var owner = Address("owner")
var receiver = Address("receiver")
var scAddress = Address("erc20")
var gasProvided = uint64(5_000_000_000)

var logBenchmark = logger.GetOrCreate("vm/benchmark")

func Test_RunERC20Benchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}
	runERC20Benchmark(t, 100, 100, false)
}

func Test_RunERC20BenchmarkFail(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runERC20Benchmark(t, 10, 100, true)
}

func Test_WarmInstancesMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runMemoryUsageBenchmark(t, 100, 200)
}

func Test_WarmInstancesFuzzyMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runMemoryUsageFuzzyBenchmark(t, 100, 100)
}

func runERC20Benchmark(tb testing.TB, nTransfers int, nRuns int, failTransaction bool) {
	totalTokenSupply := big.NewInt(int64(nTransfers * nRuns))
	mockWorld, ownerAccount, host, err := prepare(tb, owner)
	require.Nil(tb, err)

	code := testcommon.GetTestSCCode("erc20", "../../")
	deploy(tb, host, mockWorld, ownerAccount, totalTokenSupply, code)

	// Prepare ERC20 transfer call input
	transferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: owner,
			Arguments: [][]byte{
				receiver,
				big.NewInt(1).Bytes(),
			},
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasPrice:    100000000000000,
			GasProvided: gasProvided,
		},
		RecipientAddr: scAddress,
		Function:      "transferToken",
	}
	wrongArguments := [][]byte{receiver, big.NewInt(1).Bytes(), []byte("fail")}
	goodArguments := [][]byte{receiver, big.NewInt(1).Bytes()}

	// Perform ERC20 transfers
	for r := 0; r < nRuns; r++ {
		start := time.Now()
		if failTransaction {
			if r%2 == 0 {
				transferInput.Arguments = wrongArguments
			} else {
				transferInput.Arguments = goodArguments
			}
		}

		for i := 0; i < nTransfers; i++ {
			transferInput.GasProvided = gasProvided
			vmOutput, err := host.RunSmartContractCall(transferInput)
			require.Nil(tb, err)
			require.NotNil(tb, vmOutput)
			if !(failTransaction && r%2 == 0) {
				require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)
				require.Equal(tb, "", vmOutput.ReturnMessage)

				_ = mockWorld.UpdateAccounts(vmOutput.OutputAccounts, nil)
			} else {
				isProblem := checkLogsHaveDefinedString(vmOutput.Logs, "unknown")
				if isProblem {
					_ = logger.SetLogLevel("*:ERROR")
				}
				assert.False(tb, isProblem)
			}
		}
		elapsedTime := time.Since(start)
		logBenchmark.Trace("Executing ERC20 transfers", "batch", r, "transfers", nTransfers, "time", elapsedTime.String())

	}

	if !failTransaction {
		verifyTransfers(tb, mockWorld, totalTokenSupply, scAddress)
	}

	defer func() {
		err := host.Runtime().ValidateInstances()
		require.Nil(tb, err)
		host.Reset()
	}()
}

func runMemoryUsageFuzzyBenchmark(tb testing.TB, nContracts int, nTransfers int) {
	totalTokenSupply := big.NewInt(int64(nTransfers))
	mockWorld, ownerAccount, host, err := prepare(tb, owner)
	require.Nil(tb, err)

	defer func() {
		err := host.Runtime().ValidateInstances()
		require.Nil(tb, err)
		host.Reset()
	}()

	deployNContracts(tb, nContracts, mockWorld, ownerAccount, host, totalTokenSupply)

	availableContracts := make([]int, nContracts)
	remainingTransfers := make(map[int]int, nContracts)
	for i := 0; i < nContracts; i++ {
		availableContracts[i] = i
		remainingTransfers[i] = nTransfers
	}

	seed := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(seed)

	for len(availableContracts) != 0 {
		contract := availableContracts[randomizer.Intn(len(availableContracts))]
		transfers := randomizer.Intn(remainingTransfers[contract]) + 1

		for i := 0; i < transfers; i++ {
			transferInput := createTransferInput(contract)

			vmOutput, errRun := host.RunSmartContractCall(transferInput)
			require.Nil(tb, errRun)
			require.NotNil(tb, vmOutput)
			require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)

			_ = mockWorld.UpdateAccounts(vmOutput.OutputAccounts, nil)
		}

		remainingTransfers[contract] -= transfers
		if remainingTransfers[contract] == 0 {
			remainingContracts := make([]int, len(availableContracts)-1)
			j := 0
			for _, v := range availableContracts {
				if v == contract {
					continue
				}
				remainingContracts[j] = v
				j += 1
			}
			availableContracts = remainingContracts
			delete(remainingTransfers, contract)
		}
	}
	for j := 0; j < nContracts; j++ {
		verifyTransfers(tb, mockWorld, totalTokenSupply, createAddress(j))
	}
}

func runMemoryUsageBenchmark(tb testing.TB, nContracts int, nTransfers int) {
	totalTokenSupply := big.NewInt(int64(nTransfers))
	mockWorld, ownerAccount, host, err := prepare(tb, owner)
	require.Nil(tb, err)

	defer func() {
		err := host.Runtime().ValidateInstances()
		require.Nil(tb, err)
		host.Reset()
	}()

	deployNContracts(tb, nContracts, mockWorld, ownerAccount, host, totalTokenSupply)

	for i := 0; i < nContracts; i++ {
		for j := 0; j < nTransfers; j++ {
			transferInput := createTransferInput(i)

			vmOutput, err := host.RunSmartContractCall(transferInput)
			require.Nil(tb, err)
			require.NotNil(tb, vmOutput)
			require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)

			_ = mockWorld.UpdateAccounts(vmOutput.OutputAccounts, nil)
		}
	}
	for j := 0; j < nContracts; j++ {
		verifyTransfers(tb, mockWorld, totalTokenSupply, createAddress(j))
	}
}

func prepare(tb testing.TB, ownerAddress []byte) (*worldmock.MockWorld, *worldmock.Account, vmhost.VMHost, error) {
	gasMap, err := gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV3())
	require.Nil(tb, err)

	mockWorld := worldmock.NewMockWorld()
	err = mockWorld.InitBuiltinFunctions(gasMap)
	require.Nil(tb, err)

	ownerAccount := &worldmock.Account{
		Address: ownerAddress,
		Nonce:   1024,
		Balance: big.NewInt(0),
	}
	mockWorld.AcctMap.PutAccount(ownerAccount)

	require.Nil(tb, err)

	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	host, err := hostCore.NewVMHost(
		mockWorld,
		&vmhost.VMHostParameters{
			VMType:                    testcommon.DefaultVMType,
			BlockGasLimit:             uint64(1000),
			GasSchedule:               gasMap,
			BuiltInFuncContainer:      mockWorld.BuiltinFuncs.Container,
			ProtectedKeyPrefix:        []byte("E" + "L" + "R" + "O" + "N" + "D"),
			ESDTTransferParser:        esdtTransferParser,
			EpochNotifier:             &mock.EpochNotifierStub{},
			EnableEpochsHandler:       worldmock.EnableEpochsHandlerStubNoFlags(),
			WasmerSIGSEGVPassthrough:  false,
			Hasher:                    worldmock.DefaultHasher,
			MapOpcodeAddressIsAllowed: map[string]map[string]struct{}{},
		})
	require.Nil(tb, err)
	return mockWorld, ownerAccount, host, err
}

func deploy(
	tb testing.TB,
	host vmhost.VMHost,
	mockWorld *worldmock.MockWorld,
	ownerAccount *worldmock.Account,
	totalTokenSupply *big.Int,
	code []byte) {

	// Deploy ERC20
	deployInput := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: owner,
			Arguments: [][]byte{
				totalTokenSupply.Bytes(),
			},
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasPrice:    0,
			GasProvided: math.MaxInt64,
		},
		ContractCode: code,
	}

	mockWorld.NewAddressMocks = append(mockWorld.NewAddressMocks, &worldmock.NewAddressMock{
		CreatorAddress: owner,
		CreatorNonce:   ownerAccount.Nonce,
		NewAddress:     scAddress,
	})
	ownerAccount.Nonce++ // nonce increases before deploy
	vmOutput, err := host.RunSmartContractCreate(deployInput)
	require.Nil(tb, err)
	require.NotNil(tb, vmOutput)
	require.Equal(tb, "", vmOutput.ReturnMessage)
	require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)

	// Ensure the deployment persists in the mock BlockchainHook
	_ = mockWorld.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func deployNContracts(tb testing.TB, nContracts int, mockWorld *worldmock.MockWorld, ownerAccount *worldmock.Account, host vmhost.VMHost, totalTokenSupply *big.Int) {
	code := testcommon.GetTestSCCode("erc20", "../../")
	for i := 0; i < nContracts; i++ {
		modifyERC20BytecodeWithCustomTransferEvent(code, []byte{byte(i)})
		mockWorld.NewAddressMocks = append(mockWorld.NewAddressMocks, &worldmock.NewAddressMock{
			CreatorAddress: owner,
			CreatorNonce:   ownerAccount.Nonce,
			NewAddress:     createAddress(i),
		})
		deploy(tb, host, mockWorld, ownerAccount, totalTokenSupply, code)
	}
}

func createERC20Key(accountName string) string {
	keyLength := 32
	key := make([]byte, keyLength)
	key[0] = 1
	key[1] = 0
	i := 2
	for _, c := range accountName {
		key[i] = byte(c)
		i++
		if i == keyLength {
			break
		}
	}
	for q := i; q < keyLength; q++ {
		key[q] = 0
	}

	return string(key)
}

func createAddress(i int) Address {
	address := make(Address, 0)
	address = append(address, scAddress...)
	bytes := big.NewInt(int64(i)).Bytes()
	address = append(address, bytes...)
	return address
}

func createTransferInput(i int) *vmcommon.ContractCallInput {
	// Prepare ERC20 transfer call input
	transferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: owner,
			Arguments: [][]byte{
				receiver,
				big.NewInt(1).Bytes(),
			},
			CallValue:   big.NewInt(10),
			CallType:    vm.DirectCall,
			GasPrice:    10_000_000_000_000,
			GasProvided: gasProvided,
		},
		RecipientAddr: createAddress(i),
		Function:      "transferToken",
	}
	return transferInput
}

func verifyTransfers(tb testing.TB, mockWorld *worldmock.MockWorld, totalTokenSupply *big.Int, address Address) {
	ownerKey := createERC20Key("owner")
	receiverKey := createERC20Key("receiver")

	scStorage := mockWorld.AcctMap.GetAccount(address).Storage
	ownerTokens := big.NewInt(0).SetBytes(scStorage[ownerKey])
	receiverTokens := big.NewInt(0).SetBytes(scStorage[receiverKey])
	require.Equal(tb, vmhost.Zero, ownerTokens)
	require.Equal(tb, totalTokenSupply, receiverTokens)
}

func checkLogsHaveDefinedString(logs []*vmcommon.LogEntry, str string) bool {
	for _, log := range logs {
		if strings.Contains(string(log.Data[0]), str) {
			return true
		}
	}
	return false
}
