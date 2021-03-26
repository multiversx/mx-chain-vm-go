package worldmock

import (
	"bytes"
	"context"
	"errors"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/state"
)

var ErrSnapshotsNotImplemented = errors.New("snapshots not implemented")
var ErrTrieHandlingNotImplemented = errors.New("trie handling not implemented")

type MockAccountsAdapter struct {
	CurrentSnapshot AccountMap
	Snapshots       []AccountMap
}

func NewMockAccountsAdapter(accounts AccountMap) *MockAccountsAdapter {
	return &MockAccountsAdapter{
		CurrentSnapshot: accounts,
		Snapshots:       []AccountMap{accounts},
	}
}

func (a *MockAccountsAdapter) GetExistingAccount(address []byte) (state.AccountHandler, error) {
	account, exists := a.CurrentSnapshot[string(address)]
	if !exists {
		return nil, arwen.ErrInvalidAccount
	}

	return account, nil
}

func (a *MockAccountsAdapter) LoadAccount(address []byte) (state.AccountHandler, error) {
	return a.GetExistingAccount(address)
}

func (a *MockAccountsAdapter) SaveAccount(account state.AccountHandler) error {
	mockAccount, ok := account.(*Account)
	if !ok {
		return errors.New("invalid account to save")
	}

	a.CurrentSnapshot.PutAccount(mockAccount)
	return nil
}

func (a *MockAccountsAdapter) RemoveAccount(address []byte) error {
	_, exists := a.CurrentSnapshot[string(address)]
	if !exists {
		return arwen.ErrInvalidAccount
	}

	a.CurrentSnapshot.DeleteAccount(address)
	return nil
}

func (a *MockAccountsAdapter) Commit() ([]byte, error) {
	return nil, ErrSnapshotsNotImplemented
}

func (a *MockAccountsAdapter) JournalLen() int {
	return 1
}

func (a *MockAccountsAdapter) RevertToSnapshot(snapshot int) error {
	return ErrSnapshotsNotImplemented
}

func (a *MockAccountsAdapter) GetNumCheckpoints() uint32 {
	return 1
}

func (m *MockAccountsAdapter) GetCode(codeHash []byte) []byte {
	for _, account := range m.CurrentSnapshot {
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
