package mandosTests

import (
	"math/big"
	"testing"

	mge "github.com/multiversx/mx-chain-vm-v1_4-go/mandos-go/elrondgo-exporter"

	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	mj "github.com/multiversx/mx-chain-vm-v1_4-go/mandos-go/model"
	"github.com/stretchr/testify/require"
)

// address:owner
var addressOwner = []byte("owner___________________________")

// address:adder
var addressAdder = []byte("adder___________________________")

// address:alice
var addressAlice = []byte("alice___________________________")

// address:bob
var addressBob = []byte("bob_____________________________")

// sc:deployedAdder
var addressDeployedAdder = []byte("deployedAdder___________________")

func TestGetAccountsAndTransactionsFrom_Adder(t *testing.T) {
	sbi, err := mge.GetAccountsAndTransactionsFromMandos("adder.scen.json")
	require.Nil(t, err)
	expectedAccs := make([]*mge.TestAccount, 0)
	expectedDeployedAccs := make([]*mge.TestAccount, 0)
	expectedTxs := make([]*mge.Transaction, 0)
	expectedDeployTxs := make([]*mge.Transaction, 0)
	expectedBenchmarkTxPos := 1

	ownerAccount := mge.SetNewAccount(1, addressOwner, big.NewInt(48), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	scAccount := mge.SetNewAccount(0, append(mge.ScAddressPrefix, addressAdder[mge.ScAddressPrefixLength:]...), big.NewInt(0), make(map[string][]byte), arwen.GetSCCode("../../../test/adder/output/adder.wasm"), addressOwner)
	deployedScAccount := mge.SetNewAccount(0, append(mge.ScAddressPrefix, addressDeployedAdder[mge.ScAddressPrefixLength:]...), big.NewInt(0), make(map[string][]byte), make([]byte, 0), addressOwner)
	expectedAccs = append(expectedAccs, ownerAccount, scAccount)
	expectedDeployedAccs = append(expectedDeployedAccs, deployedScAccount)

	transaction := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), sbi.Accs[0].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	expectedTxs = append(expectedTxs, transaction, transaction)

	require.Nil(t, err)
	require.Equal(t, expectedBenchmarkTxPos, sbi.BenchmarkTxPos)
	require.Equal(t, expectedAccs, sbi.Accs)
	require.Equal(t, expectedDeployedAccs, sbi.DeployedAccs)
	require.Equal(t, expectedDeployTxs, sbi.DeployTxs)
	require.Equal(t, expectedTxs, sbi.Txs)
}

func TestGetAccountsAndTransactionsFrom_AdderWithExternalSteps(t *testing.T) {
	sbi, err := mge.GetAccountsAndTransactionsFromMandos("adder_with_external_steps.scen.json")
	require.Nil(t, err)
	expectedAccs := make([]*mge.TestAccount, 0)
	expectedTxs := make([]*mge.Transaction, 0)
	expectedDeployTxs := make([]*mge.Transaction, 0)
	expectedBenchmarkTxPos := 1

	ownerAccount := mge.SetNewAccount(1, addressOwner, big.NewInt(48), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	scAccount := mge.SetNewAccount(0, append(mge.ScAddressPrefix, addressAdder[mge.ScAddressPrefixLength:]...), big.NewInt(0), make(map[string][]byte), arwen.GetSCCode("../../../test/adder/output/adder.wasm"), addressOwner)
	aliceAccount := mge.SetNewAccount(5, addressAlice, big.NewInt(284), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	bobAccount := mge.SetNewAccount(3, addressBob, big.NewInt(11), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	expectedAccs = append(expectedAccs, aliceAccount, scAccount, bobAccount, ownerAccount)
	require.Equal(t, expectedAccs, sbi.Accs)

	transactionAlice := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), sbi.Accs[0].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	transactionBob := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), sbi.Accs[2].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	transactionOwner := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), sbi.Accs[3].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	expectedTxs = append(expectedTxs, transactionBob, transactionAlice, transactionOwner)
	require.Equal(t, expectedBenchmarkTxPos, sbi.BenchmarkTxPos)
	require.Equal(t, expectedTxs, sbi.Txs)
	require.Equal(t, expectedDeployTxs, sbi.DeployTxs)
}
