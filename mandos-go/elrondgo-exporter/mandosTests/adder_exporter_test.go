package mandosTests

import (
	"math/big"
	"testing"

	elrondgo_exporter "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwenmandos/elrondgo-exporter"

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

func TestGetAccountsAndTransactionsFrom_Adder(t *testing.T) {
	accounts, transactions, err := elrondgo_exporter.GetAccountsAndTransactionsFromMandos("adder.scen.json")
	require.Nil(t, err)
	expectedAccs := make([]*elrondgo_exporter.TestAccount, 0)
	expectedTxs := make([]*elrondgo_exporter.Transaction, 0)

	ownerAccount := elrondgo_exporter.SetNewAccount(1, addressOwner, big.NewInt(48), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	scAccount := elrondgo_exporter.SetNewAccount(0, addressAdder, big.NewInt(0), make(map[string][]byte), arwen.GetSCCode("../../../test/adder/output/adder.wasm"), addressOwner)
	expectedAccs = append(expectedAccs, ownerAccount, scAccount)

	transaction := elrondgo_exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[0].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	expectedTxs = append(expectedTxs, transaction)

	require.Nil(t, err)
	require.Equal(t, expectedAccs, accounts)
	require.Equal(t, expectedTxs, transactions)
}

func TestGetAccountsAndTransactionsFrom_AdderWithExternalSteps(t *testing.T) {
	accounts, transactions, err := elrondgo_exporter.GetAccountsAndTransactionsFromMandos("adder_with_external_steps.scen.json")
	require.Nil(t, err)
	expectedAccs := make([]*elrondgo_exporter.TestAccount, 0)
	expectedTxs := make([]*elrondgo_exporter.Transaction, 0)

	ownerAccount := elrondgo_exporter.SetNewAccount(1, addressOwner, big.NewInt(48), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	scAccount := elrondgo_exporter.SetNewAccount(0, addressAdder, big.NewInt(0), make(map[string][]byte), arwen.GetSCCode("../../../test/adder/output/adder.wasm"), addressOwner)
	aliceAccount := elrondgo_exporter.SetNewAccount(5, addressAlice, big.NewInt(284), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	bobAccount := elrondgo_exporter.SetNewAccount(3, addressBob, big.NewInt(11), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	expectedAccs = append(expectedAccs, aliceAccount, scAccount, bobAccount, ownerAccount)
	require.Equal(t, expectedAccs, accounts)

	transactionAlice := elrondgo_exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[0].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	transactionBob := elrondgo_exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[2].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	transactionOwner := elrondgo_exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[3].GetAddress(), accounts[1].GetAddress(), 5000000, 0)
	expectedTxs = append(expectedTxs, transactionBob, transactionAlice, transactionOwner)
	require.Equal(t, expectedTxs, transactions)
}
