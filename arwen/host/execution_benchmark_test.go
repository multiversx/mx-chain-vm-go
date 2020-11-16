package host

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

var owner = []byte("owner")
var receiver = []byte("receiver")
var scAddress = []byte("erc20")

func Test_RunERC20Benchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runERC20Benchmark(t, 1000, 4)
}

func runERC20Benchmark(tb testing.TB, nTransfers int, nRuns int) {
	totalTokenSupply := big.NewInt(int64(nTransfers * nRuns))
	host, mockBlockchainHook := deploy(tb, totalTokenSupply)

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
			CallType:    vmcommon.DirectCall,
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

			mockBlockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
		}
		elapsedTime := time.Since(start)
		fmt.Printf("Executing %d ERC20 transfers: %s\n", nTransfers, elapsedTime.String())
	}

	verifyTransfers(tb, mockBlockchainHook, totalTokenSupply)
}

func deploy(tb testing.TB, totalTokenSupply *big.Int) (*vmHost, *mock.BlockchainHookMock) {
	// Prepare the host
	mockBlockchainHook := mock.NewBlockchainHookMock()
	mockBlockchainHook.AddAccount(&mock.AccountMock{
		Address: owner,
		Nonce:   1024,
		Balance: big.NewInt(0),
	})

	gasMap, err := LoadGasScheduleConfig("../../test/gasSchedule.toml")
	require.Nil(tb, err)

	host, err := NewArwenVM(mockBlockchainHook, &arwen.VMHostParameters{
		VMType:                   defaultVMType,
		BlockGasLimit:            uint64(1000),
		GasSchedule:              gasMap,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
	})
	require.Nil(tb, err)

	// Deploy ERC20
	deployInput := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: owner,
			Arguments: [][]byte{
				totalTokenSupply.Bytes(),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0xFFFFFFFFFFFFFFFF,
		},
		ContractCode: GetTestSCCode("erc20", "../../"),
	}

	mockBlockchainHook.NewAddr = scAddress
	vmOutput, err := host.RunSmartContractCreate(deployInput)
	require.Nil(tb, err)
	require.NotNil(tb, vmOutput)
	require.Equal(tb, "", vmOutput.ReturnMessage)
	require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)

	// Ensure the deployment persists in the mock BlockchainHook
	mockBlockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	return host, mockBlockchainHook
}

func verifyTransfers(tb testing.TB, mockBlockchainHook *mock.BlockchainHookMock, totalTokenSupply *big.Int) {
	ownerKey := createERC20Key("owner")
	receiverKey := createERC20Key("receiver")

	scStorage := mockBlockchainHook.Accounts[string(scAddress)].Storage
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
