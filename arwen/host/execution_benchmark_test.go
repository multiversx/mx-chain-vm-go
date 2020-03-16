package host

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var owner = []byte("owner")
var receiver = []byte("receiver")
var scAddress = []byte("erc20")
var nTransfers = 100
var nRuns = 1
var totalTokenSupply = big.NewInt(int64(nTransfers * nRuns))

func Test_RunERC20Benchmark(t *testing.T) {
	runERC20Benchmark(t)
}

func runERC20Benchmark(tb testing.TB) {
	host, mockBlockchainHook := deploy(tb)

	// Prepare ERC20 transfer call input
	transferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: owner,
			Arguments: [][]byte{
				receiver,
				big.NewInt(1).Bytes(),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 100000,
		},
		RecipientAddr: scAddress,
		Function:      "transferToken",
	}

	// Perform ERC20 transfers
	for r := 0; r < nRuns; r++ {
		TotalWasmerExecution = 0
		start := time.Now()
		for i := 0; i < nTransfers; i++ {
			transferInput.GasProvided = 100000
			vmOutput, err := host.RunSmartContractCall(transferInput)
			require.Nil(tb, err)
			require.NotNil(tb, vmOutput)
			require.Equal(tb, "", vmOutput.ReturnMessage)
			require.Equal(tb, vmcommon.Ok, vmOutput.ReturnCode)

			mockBlockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
		}
		elapsedTime := time.Since(start)
		fmt.Printf("Executing %d ERC20 transfers: %s\n", nTransfers, elapsedTime.String())
	}

	verifyTransfers(tb, mockBlockchainHook)
}

func deploy(tb testing.TB) (*vmHost, *mock.BlockchainHookMock) {
	// Prepare the host
	mockBlockchainHook := mock.NewBlockchainHookMock()
	mockBlockchainHook.AddAccount(&mock.Account{
		Address: owner,
		Nonce:   1024,
		Balance: big.NewInt(88000),
	})

	host, err := NewArwenVM(mockBlockchainHook, &mock.CryptoHookMock{}, defaultVmType, uint64(1000), config.MakeGasMap(1))
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
			GasProvided: 100000,
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

func verifyTransfers(tb testing.TB, mockBlockchainHook *mock.BlockchainHookMock) {
	ownerKey := string([]byte{
		1, 0, 'o', 'w', 'n', 'e', 'r', 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})
	receiverKey := string([]byte{
		1, 0, 'r', 'e', 'c', 'e', 'i', 'v',
		'e', 'r', 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})

	scStorage := mockBlockchainHook.Accounts[string(scAddress)].Storage
	ownerTokens := big.NewInt(0).SetBytes(scStorage[ownerKey])
	receiverTokens := big.NewInt(0).SetBytes(scStorage[receiverKey])
	require.Equal(tb, arwen.Zero, ownerTokens)
	require.Equal(tb, totalTokenSupply, receiverTokens)
}
