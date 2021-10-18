package elrondgo_exporter

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
	"github.com/stretchr/testify/require"
)

// address:owner
var addressOwner = []byte{111, 119, 110, 101, 114, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95}

// address:adder
var addressAdder = []byte{97, 100, 100, 101, 114, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95, 95}

func TestGetAccountsAndTransactionsFromAdder(t *testing.T) {
	accounts, transactions, err := GetAccountsAndTransactionsFromMandos("./mandosTests/adder.scen.json")
	require.Nil(t, err)
	expectedAccs := make([]*TestAccount, 0)
	expectedTxs := make([]*Transaction, 0)

	ownerAccount := SetNewAccount(1, addressOwner, big.NewInt(48), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	scAccount := SetNewAccount(0, addressAdder, big.NewInt(0), make(map[string][]byte), arwen.GetSCCode("../../test/adder/output/adder.wasm"), addressOwner)
	expectedAccs = append(expectedAccs, ownerAccount, scAccount)
	require.Equal(t, 2, len(expectedAccs))

	transaction := CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[0].address, accounts[1].address, 5000000, 0)
	expectedTxs = append(expectedTxs, transaction)

	require.Nil(t, err)
	require.Equal(t, expectedAccs, accounts)
	require.Equal(t, expectedTxs, transactions)
}
