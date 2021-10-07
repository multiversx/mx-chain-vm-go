package elrondgo_exporter

import (
	"math/big"
	"testing"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
	"github.com/stretchr/testify/require"
)

// address:owner
var addressOwner = []byte{111, 119, 110, 101, 114, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95}

// address:adder
var addressAdder = []byte{0, 0, 0, 0, 0, 0, 0, 0, 97, 100, 100, 101, 114, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95}

func TestGetAccountsAndTransactionsFromAdder(t *testing.T) {
	accounts, scAccounts, transactions, err := GetAccountsAndTransactionsFromMandos("./mandos/adder.scen.json")
	expectedAccs := make([]*TestAccount, 0)
	expectedScAccs := make([]*TestAccount, 0)
	expectedTxs := make([]*Transaction, 0)
	adderSCAcc := NewTestAccount().WithAddress(addressAdder)

	expectedAccs = append(expectedAccs, NewTestAccount().WithAddress(addressOwner).WithNonce(1))
	expectedScAccs = append(expectedScAccs, adderSCAcc)

	transaction1 := CreateDeployTransaction("file:../../../test/adder/output/adder.wasm", [][]byte{{5}}, 0, big.NewInt(0), addressOwner, 5000000, 0)
	transaction2 := CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[0].address, scAccounts[0].address, 5000000, 0)
	expectedTxs = append(expectedTxs, transaction1, transaction2)

	require.Nil(t, err)
	require.Equal(t, expectedAccs, accounts)
	require.Equal(t, scAccounts, expectedScAccs)
	require.Equal(t, expectedTxs, transactions)
}
