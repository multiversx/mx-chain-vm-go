package hosttest

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/wasm-vm/arwen/host"
	"github.com/ElrondNetwork/wasm-vm/arwen/mock"
	gasSchedules "github.com/ElrondNetwork/wasm-vm/arwenmandos/gasSchedules"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	"github.com/ElrondNetwork/wasm-vm/testcommon"
	"github.com/stretchr/testify/require"
)

type Address = []byte

var owner = Address("owner")
var receiver = Address("receiver")
var scAddress = Address("erc20")
var gasProvided = uint64(5000000000)

func Test_RunERC20Benchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runERC20Benchmark(t, 1000, 4)
}

func Test_WarmInstancesMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runMemoryUsageBenchmark(t, 20, 1_000_000)
}

func runMemoryUsageBenchmark(tb testing.TB, nContracts int, nTransfers int) {
	totalTokenSupply := big.NewInt(int64(nTransfers))
	mockWorld, ownerAccount, host, err := prepare(tb)
	require.Nil(tb, err)

	defer func() {
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
			require.Equal(tb, "", vmOutput.ReturnMessage)

			_ = mockWorld.UpdateAccounts(vmOutput.OutputAccounts, nil)
		}
		fmt.Printf("Executing %d ERC20 transfers for contract %d\n", nTransfers, i)
	}
	for j := 0; j < nContracts; j++ {
		verifyTransfers(tb, mockWorld, totalTokenSupply, createAddress(j))
	}
}

func createAddress(i int) Address {
	address := make(Address, 0)
	address = append(address, scAddress...)
	address = append(address, '0'+byte(i))
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
			GasPrice:    100000000000000,
			GasProvided: gasProvided,
		},
		RecipientAddr: createAddress(i),
		Function:      "transferToken",
	}
	return transferInput
}

func deployNContracts(tb testing.TB, nContracts int, mockWorld *worldmock.MockWorld, ownerAccount *worldmock.Account, host arwen.VMHost, totalTokenSupply *big.Int) {
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

func runERC20Benchmark(tb testing.TB, nTransfers int, nRuns int) {
	totalTokenSupply := big.NewInt(int64(nTransfers * nRuns))
	mockWorld, ownerAccount, host, err := prepare(tb)
	require.Nil(tb, err)

	code := testcommon.GetTestSCCode("erc20", "../../")
	deploy(tb, host, mockWorld, ownerAccount, totalTokenSupply, code)

	gasProvided := uint64(5000000000)

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
			GasPrice:    100000000000000,
			GasProvided: gasProvided,
		},
		RecipientAddr: scAddress,
		Function:      "transferToken",
	}

	// Perform ERC20 transfers
	for r := 0; r < nRuns; r++ {
		start := time.Now()
		for i := 0; i < nTransfers; i++ {
			transferInput.GasProvided = gasProvided
			vmOutput, err := host.RunSmartContractCall(transferInput)
			require.Nil(tb, err)
			require.NotNil(tb, vmOutput)
			require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)
			require.Equal(tb, "", vmOutput.ReturnMessage)

			_ = mockWorld.UpdateAccounts(vmOutput.OutputAccounts, nil)
		}
		elapsedTime := time.Since(start)
		fmt.Printf("Executing %d ERC20 transfers: %s\n", nTransfers, elapsedTime.String())
	}

	verifyTransfers(tb, mockWorld, totalTokenSupply, scAddress)
	defer func() {
		host.Reset()
	}()
}

func deploy(
	tb testing.TB,
	host arwen.VMHost,
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
			GasProvided: 0xFFFFFFFFFFFFFFFF,
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

func prepare(tb testing.TB) (*worldmock.MockWorld, *worldmock.Account, arwen.VMHost, error) {
	// Prepare the host
	mockWorld := worldmock.NewMockWorld()
	ownerAccount := &worldmock.Account{
		Address: owner,
		Nonce:   1024,
		Balance: big.NewInt(0),
	}
	mockWorld.AcctMap.PutAccount(ownerAccount)

	gasMap, err := gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV3())
	require.Nil(tb, err)

	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	host, err := arwenHost.NewArwenVM(mockWorld, &arwen.VMHostParameters{
		VMType:                   testcommon.DefaultVMType,
		BlockGasLimit:            uint64(1000),
		GasSchedule:              gasMap,
		BuiltInFuncContainer:     builtInFunctions.NewBuiltInFunctionContainer(),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
		ESDTTransferParser:       esdtTransferParser,
		EpochNotifier:            &mock.EpochNotifierStub{},
		EnableEpochsHandler:      &mock.EnableEpochsHandlerStub{},
		WasmerSIGSEGVPassthrough: false,
	})
	require.Nil(tb, err)
	return mockWorld, ownerAccount, host, err
}

func verifyTransfers(tb testing.TB, mockWorld *worldmock.MockWorld, totalTokenSupply *big.Int, address Address) {
	ownerKey := createERC20Key("owner")
	receiverKey := createERC20Key("receiver")

	scStorage := mockWorld.AcctMap.GetAccount(address).Storage
	ownerTokens := big.NewInt(0).SetBytes(scStorage[ownerKey])
	receiverTokens := big.NewInt(0).SetBytes(scStorage[receiverKey])
	require.Equal(tb, arwen.Zero, ownerTokens)
	require.Equal(tb, totalTokenSupply, receiverTokens)
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
