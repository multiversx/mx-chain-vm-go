package context

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"math/big"
)

type BigIntContainer interface {
	Clean()
	Put(value int64) int32
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
	IsInterfaceNil() bool
}

// ArgumentsParser defines the functionality to parse transaction data into arguments and code for smart contracts
type ArgumentsParser interface {
	GetArguments() ([][]byte, error)
	GetCode() ([]byte, error)
	GetFunction() (string, error)
	ParseData(data string) error

	CreateDataFromStorageUpdate(storageUpdates []*vmcommon.StorageUpdate) string
	GetStorageUpdates(data string) ([]*vmcommon.StorageUpdate, error)
	IsInterfaceNil() bool
}
