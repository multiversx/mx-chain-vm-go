package context

import (
	"math/big"
	"testing"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/assert"
)

func TestMergingTwoOutputAccounts(t *testing.T) {
	left := &vmcommon.OutputAccount{
		Address:        []byte("account"),
		Nonce:          0,
		Balance:        nil,
		BalanceDelta:   big.NewInt(2),
		StorageUpdates: []*vmcommon.StorageUpdate {
      &vmcommon.StorageUpdate {
        Offset: []byte("mykey-1"),
        Data: []byte("mydata-1"),
      },
    },
		Code:           []byte{42, 99, 101},
		Data:           make([]byte, 0),
		GasLimit:       16,
	}

	right := &vmcommon.OutputAccount{
		Address:        []byte("account"),
		Nonce:          1,
		Balance:        nil,
		BalanceDelta:   big.NewInt(2),
		StorageUpdates: make([]*vmcommon.StorageUpdate, 0),
		Code:           make([]byte, 0),
		Data:           []byte("some data"),
		GasLimit:       17,
	}

	expected := &vmcommon.OutputAccount{
		Address:        []byte("account"),
		Nonce:          1,
		Balance:        nil,
		BalanceDelta:   big.NewInt(4),
		StorageUpdates: []*vmcommon.StorageUpdate {
      &vmcommon.StorageUpdate {
        Offset: []byte("mykey-1"),
        Data: []byte("mydata-1"),
      },
    },
		Code:           make([]byte, 0),
		Data:           []byte("some data"),
		GasLimit:       17,
	}

	mergedAccount := mergeOutputAccounts(left, right)

	assert.NotNil(t, mergedAccount)
	assert.Equal(t, expected, mergedAccount)
}
