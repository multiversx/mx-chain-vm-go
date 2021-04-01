package worldmock

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/state"
)

var ErrSnapshotsNotImplemented = errors.New("snapshots not implemented")
var ErrTrieHandlingNotImplemented = errors.New("trie handling not implemented")

type MockAccountsAdapter struct {
	World     *MockWorld
	Snapshots []AccountMap
}

func NewMockAccountsAdapter(world *MockWorld) *MockAccountsAdapter {
	return &MockAccountsAdapter{
		World:     world,
		Snapshots: make([]AccountMap, 0),
	}
}

func (m *MockAccountsAdapter) GetExistingAccount(address []byte) (state.AccountHandler, error) {
	account, exists := m.World.AcctMap[string(address)]
	if !exists {
		return nil, arwen.ErrInvalidAccount
	}

	return account, nil
}

func (m *MockAccountsAdapter) LoadAccount(address []byte) (state.AccountHandler, error) {
	return m.GetExistingAccount(address)
}

func (m *MockAccountsAdapter) SaveAccount(account state.AccountHandler) error {
	mockAccount, ok := account.(*Account)
	if !ok {
		return errors.New("invalid account to save")
	}

	m.World.AcctMap.PutAccount(mockAccount)
	return nil
}

func (m *MockAccountsAdapter) RemoveAccount(address []byte) error {
	_, exists := m.World.AcctMap[string(address)]
	if !exists {
		return arwen.ErrInvalidAccount
	}

	m.World.AcctMap.DeleteAccount(address)
	return nil
}

func (m *MockAccountsAdapter) Commit() ([]byte, error) {
	return nil, nil
}

func (m *MockAccountsAdapter) JournalLen() int {
	return 1
}

func (m *MockAccountsAdapter) RevertToSnapshot(snapshotIndex int) error {
	if len(m.Snapshots) == 0 {
		return errors.New("no snapshots")
	}

	if snapshotIndex >= len(m.Snapshots) || snapshotIndex < 0 {
		return fmt.Errorf(
			"snapshot %d out of bounds (min 0, max %d)",
			snapshotIndex,
			len(m.Snapshots)-1)
	}

	snapshot := m.Snapshots[snapshotIndex]
	m.Snapshots = m.Snapshots[:snapshotIndex]

	m.World.AcctMap.LoadAccountStorageFrom(snapshot)

	// TODO should probably set BalanceDelta of all accounts to 0 as well
	return nil
}

func (m *MockAccountsAdapter) GetNumCheckpoints() uint32 {
	return uint32(len(m.Snapshots))
}

func (m *MockAccountsAdapter) GetCode(codeHash []byte) []byte {
	for _, account := range m.World.AcctMap {
		if bytes.Equal(account.GetCodeHash(), codeHash) {
			return account.GetCode()
		}
	}

	return nil
}

func (m *MockAccountsAdapter) RootHash() ([]byte, error) {
	return nil, ErrTrieHandlingNotImplemented
}

func (m *MockAccountsAdapter) RecreateTrie(rootHash []byte) error {
	return ErrTrieHandlingNotImplemented
}

func (m *MockAccountsAdapter) PruneTrie(rootHash []byte, identifier data.TriePruningIdentifier) {
}

func (m *MockAccountsAdapter) CancelPrune(rootHash []byte, identifier data.TriePruningIdentifier) {
}

func (m *MockAccountsAdapter) SnapshotState(rootHash []byte, ctx context.Context) {
	snapshot := m.World.AcctMap.Clone()
	m.Snapshots = append(m.Snapshots, snapshot)
}

func (m *MockAccountsAdapter) SetStateCheckpoint(rootHash []byte, ctx context.Context) {
}

func (m *MockAccountsAdapter) IsPruningEnabled() bool {
	return false
}

func (m *MockAccountsAdapter) GetAllLeaves(rootHash []byte, ctx context.Context) (chan core.KeyValueHolder, error) {
	return nil, ErrTrieHandlingNotImplemented
}

func (m *MockAccountsAdapter) RecreateAllTries(rootHash []byte, ctx context.Context) (map[string]data.Trie, error) {
	return nil, ErrTrieHandlingNotImplemented
}

func (m *MockAccountsAdapter) IsInterfaceNil() bool {
	return m == nil
}
