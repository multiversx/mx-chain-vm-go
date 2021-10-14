package elrondgo_exporter

import (
	"math/big"
	"testing"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
	"github.com/stretchr/testify/require"
)

var adderCode = []byte{0, 97, 115, 109, 1, 0, 0, 0, 1, 39, 8, 96, 0, 1, 127, 96, 2, 127, 127, 0, 96, 2, 127, 127, 1, 127, 96, 1, 126, 1, 127, 96, 0, 0, 96, 1, 127, 0, 96, 3, 127, 127, 127, 0, 96, 1, 127, 1, 127, 2, 173, 2, 13, 3, 101, 110, 118, 15, 103, 101, 116, 78, 117, 109, 65, 114, 103, 117, 109, 101, 110, 116, 115, 0, 0, 3, 101, 110, 118, 11, 115, 105, 103, 110, 97, 108, 69, 114, 114, 111, 114, 0, 1, 3, 101, 110, 118, 10, 109, 66, 117, 102, 102, 101, 114, 78, 101, 119, 0, 0, 3, 101, 110, 118, 18, 109, 66, 117, 102, 102, 101, 114, 83, 116, 111, 114, 97, 103, 101, 76, 111, 97, 100, 0, 2, 3, 101, 110, 118, 9, 98, 105, 103, 73, 110, 116, 78, 101, 119, 0, 3, 3, 101, 110, 118, 21, 109, 66, 117, 102, 102, 101, 114, 84, 111, 66, 105, 103, 73, 110, 116, 83, 105, 103, 110, 101, 100, 0, 2, 3, 101, 110, 118, 23, 109, 66, 117, 102, 102, 101, 114, 70, 114, 111, 109, 66, 105, 103, 73, 110, 116, 83, 105, 103, 110, 101, 100, 0, 2, 3, 101, 110, 118, 19, 109, 66, 117, 102, 102, 101, 114, 83, 116, 111, 114, 97, 103, 101, 83, 116, 111, 114, 101, 0, 2, 3, 101, 110, 118, 19, 109, 66, 117, 102, 102, 101, 114, 78, 101, 119, 70, 114, 111, 109, 66, 121, 116, 101, 115, 0, 2, 3, 101, 110, 118, 14, 99, 104, 101, 99, 107, 78, 111, 80, 97, 121, 109, 101, 110, 116, 0, 4, 3, 101, 110, 118, 18, 98, 105, 103, 73, 110, 116, 70, 105, 110, 105, 115, 104, 83, 105, 103, 110, 101, 100, 0, 5, 3, 101, 110, 118, 23, 98, 105, 103, 73, 110, 116, 71, 101, 116, 83, 105, 103, 110, 101, 100, 65, 114, 103, 117, 109, 101, 110, 116, 0, 1, 3, 101, 110, 118, 9, 98, 105, 103, 73, 110, 116, 65, 100, 100, 0, 6, 3, 9, 8, 5, 7, 1, 0, 4, 4, 4, 4, 5, 3, 1, 0, 17, 6, 25, 3, 127, 1, 65, 128, 128, 192, 0, 11, 127, 0, 65, 156, 128, 192, 0, 11, 127, 0, 65, 160, 128, 192, 0, 11, 7, 70, 7, 6, 109, 101, 109, 111, 114, 121, 2, 0, 6, 103, 101, 116, 83, 117, 109, 0, 17, 4, 105, 110, 105, 116, 0, 18, 3, 97, 100, 100, 0, 19, 8, 99, 97, 108, 108, 66, 97, 99, 107, 0, 20, 10, 95, 95, 100, 97, 116, 97, 95, 101, 110, 100, 3, 1, 11, 95, 95, 104, 101, 97, 112, 95, 98, 97, 115, 101, 3, 2, 10, 156, 3, 8, 32, 0, 2, 64, 16, 128, 128, 128, 128, 0, 32, 0, 71, 13, 0, 15, 11, 65, 128, 128, 192, 128, 0, 65, 25, 16, 129, 128, 128, 128, 0, 0, 11, 45, 1, 1, 127, 32, 0, 40, 2, 0, 16, 130, 128, 128, 128, 0, 34, 0, 16, 131, 128, 128, 128, 0, 26, 32, 0, 66, 0, 16, 132, 128, 128, 128, 0, 34, 1, 16, 133, 128, 128, 128, 0, 26, 32, 1, 11, 35, 1, 1, 127, 16, 130, 128, 128, 128, 0, 34, 2, 32, 1, 16, 134, 128, 128, 128, 0, 26, 32, 0, 40, 2, 0, 32, 2, 16, 135, 128, 128, 128, 0, 26, 11, 16, 0, 65, 153, 128, 192, 128, 0, 65, 3, 16, 136, 128, 128, 128, 0, 11, 74, 1, 1, 127, 35, 128, 128, 128, 128, 0, 65, 16, 107, 34, 0, 36, 128, 128, 128, 128, 0, 16, 137, 128, 128, 128, 0, 65, 0, 16, 141, 128, 128, 128, 0, 32, 0, 16, 144, 128, 128, 128, 0, 54, 2, 12, 32, 0, 65, 12, 106, 16, 142, 128, 128, 128, 0, 16, 138, 128, 128, 128, 0, 32, 0, 65, 16, 106, 36, 128, 128, 128, 128, 0, 11, 88, 1, 2, 127, 35, 128, 128, 128, 128, 0, 65, 16, 107, 34, 0, 36, 128, 128, 128, 128, 0, 16, 137, 128, 128, 128, 0, 65, 1, 16, 141, 128, 128, 128, 0, 65, 0, 66, 0, 16, 132, 128, 128, 128, 0, 34, 1, 16, 139, 128, 128, 128, 0, 32, 0, 16, 144, 128, 128, 128, 0, 54, 2, 12, 32, 0, 65, 12, 106, 32, 1, 16, 143, 128, 128, 128, 0, 32, 0, 65, 16, 106, 36, 128, 128, 128, 128, 0, 11, 111, 1, 3, 127, 35, 128, 128, 128, 128, 0, 65, 16, 107, 34, 0, 36, 128, 128, 128, 128, 0, 16, 137, 128, 128, 128, 0, 65, 1, 16, 141, 128, 128, 128, 0, 65, 0, 66, 0, 16, 132, 128, 128, 128, 0, 34, 1, 16, 139, 128, 128, 128, 0, 32, 0, 16, 144, 128, 128, 128, 0, 54, 2, 12, 32, 0, 65, 12, 106, 16, 142, 128, 128, 128, 0, 34, 2, 32, 2, 32, 1, 16, 140, 128, 128, 128, 0, 32, 0, 65, 12, 106, 32, 2, 16, 143, 128, 128, 128, 0, 32, 0, 65, 16, 106, 36, 128, 128, 128, 128, 0, 11, 2, 0, 11, 11, 37, 1, 0, 65, 128, 128, 192, 0, 11, 28, 119, 114, 111, 110, 103, 32, 110, 117, 109, 98, 101, 114, 32, 111, 102, 32, 97, 114, 103, 117, 109, 101, 110, 116, 115, 115, 117, 109}

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
	scAccount := SetNewAccount(0, addressAdder, big.NewInt(0), make(map[string][]byte), adderCode, addressOwner)
	expectedAccs = append(expectedAccs, ownerAccount, scAccount)
	require.Equal(t, 2, len(expectedAccs))

	transaction := CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*mj.ESDTTxData, 0), accounts[0].address, accounts[1].address, 5000000, 0)
	expectedTxs = append(expectedTxs, transaction)

	require.Nil(t, err)
	require.Equal(t, expectedAccs, accounts)
	require.Equal(t, expectedTxs, transactions)
}