package mandosTests

import (
	"math/big"
	"testing"

	mge "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/elrondgo-exporter"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
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
	accounts, deployedAccounts, transactions, deployTxs, benchmarkTxPos, err := mge.GetAccountsAndTransactionsFromMandos("adder.scen.json")
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

	transaction := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[0].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	expectedTxs = append(expectedTxs, transaction, transaction)

	require.Nil(t, err)
	require.Equal(t, expectedBenchmarkTxPos, benchmarkTxPos)
	require.Equal(t, expectedAccs, accounts)
	require.Equal(t, expectedDeployedAccs, deployedAccounts)
	require.Equal(t, expectedDeployTxs, deployTxs)
	require.Equal(t, expectedTxs, transactions)
}

func TestGetAccountsAndTransactionsFrom_AdderWithExternalSteps(t *testing.T) {
	accounts, _, transactions, deployTxs, benchmarkTxPos, err := mge.GetAccountsAndTransactionsFromMandos("adder_with_external_steps.scen.json")
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
	require.Equal(t, expectedAccs, accounts)

	transactionAlice := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[0].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	transactionBob := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[2].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	transactionOwner := mge.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[3].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	expectedTxs = append(expectedTxs, transactionBob, transactionAlice, transactionOwner)
	require.Equal(t, expectedBenchmarkTxPos, benchmarkTxPos)
	require.Equal(t, expectedTxs, transactions)
	require.Equal(t, expectedDeployTxs, deployTxs)
}
