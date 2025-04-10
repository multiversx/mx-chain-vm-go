package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ = (StateDB)((*EVMStateDB)(nil))

type EVMStateDB struct {
	executor.EVMHooks
	transientStorage transientStorage
}

func CreateEVMStateDB(evmHooks executor.EVMHooks) *EVMStateDB {
	return &EVMStateDB{
		EVMHooks:         evmHooks,
		transientStorage: newTransientStorage(),
	}
}

func (evmState *EVMStateDB) GetTransientState(address common.Address, key common.Hash) common.Hash {
	return evmState.transientStorage.Get(address, key)
}

func (evmState *EVMStateDB) SetTransientState(address common.Address, key common.Hash, value common.Hash) {
	prev := evmState.GetTransientState(address, key)
	if prev == value {
		return
	}
	evmState.transientStorage.Set(address, key, value)
}
