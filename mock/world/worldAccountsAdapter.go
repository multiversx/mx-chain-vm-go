package worldmock

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen"
)

// ErrTrieHandlingNotImplemented indicates that no trie-related operations are
// currently implemented.
var ErrTrieHandlingNotImplemented = errors.New("trie handling not implemented")

// MockAccountsAdapter is an implementation of AccountsAdapter based on
// MockWorld and the accounts within it.
type MockAccountsAdapter struct {
	World     *MockWorld
	Snapshots []AccountMap
}

// NewMockAccountsAdapter instantiates a new MockAccountsAdapter.
func NewMockAccountsAdapter(world *MockWorld) *MockAccountsAdapter {
	return &MockAccountsAdapter{
		World:     world,
		Snapshots: make([]AccountMap, 0),
	}
}

// GetExistingAccount -
func (m *MockAccountsAdapter) GetExistingAccount(address []byte) (vmcommon.AccountHandler, error) {
	account, exists := m.World.AcctMap[string(address)]
	if !exists {
		return nil, arwen.ErrInvalidAccount
	}

	return account, nil
}

// LoadAccount -
func (m *MockAccountsAdapter) LoadAccount(address []byte) (vmcommon.AccountHandler, error) {
	account, exists := m.World.AcctMap[string(address)]
	if !exists {
		account = m.World.AcctMap.CreateAccount(address, m.World)
	}

	return account, nil
}

// SaveAccount -
func (m *MockAccountsAdapter) SaveAccount(account vmcommon.AccountHandler) error {
	mockAccount, ok := account.(*Account)
	if !ok {
		return errors.New("invalid account to save")
	}

	m.World.AcctMap.PutAccount(mockAccount)
	return nil
}

// RemoveAccount -
func (m *MockAccountsAdapter) RemoveAccount(address []byte) error {
	_, exists := m.World.AcctMap[string(address)]
	if !exists {
		return arwen.ErrInvalidAccount
	}

	m.World.AcctMap.DeleteAccount(address)
	return nil
}

// Commit -
func (m *MockAccountsAdapter) Commit() ([]byte, error) {
	m.Snapshots = make([]AccountMap, 0)
	return nil, nil
}

// JournalLen -
func (m *MockAccountsAdapter) JournalLen() int {
	return len(m.Snapshots) - 1
}

// RevertToSnapshot -
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

	// TODO should probably set BalanceDelta of all accounts to 0 as well?
	return m.World.AcctMap.LoadAccountStorageFrom(snapshot)
}

// GetNumCheckpoints -
func (m *MockAccountsAdapter) GetNumCheckpoints() uint32 {
	return uint32(len(m.Snapshots))
}

// GetCode -
func (m *MockAccountsAdapter) GetCode(codeHash []byte) []byte {
	for _, account := range m.World.AcctMap {
		if bytes.Equal(account.GetCodeHash(), codeHash) {
			return account.GetCode()
		}
	}

	return nil
}

// RootHash -
func (m *MockAccountsAdapter) RootHash() ([]byte, error) {
	return nil, ErrTrieHandlingNotImplemented
}

// RecreateTrie -
func (m *MockAccountsAdapter) RecreateTrie(_ []byte) error {
	return ErrTrieHandlingNotImplemented
}

// SnapshotState -
func (m *MockAccountsAdapter) SnapshotState(_ []byte, _ context.Context) {
	snapshot := m.World.AcctMap.Clone()
	m.Snapshots = append(m.Snapshots, snapshot)
}

// SetStateCheckpoint -
func (m *MockAccountsAdapter) SetStateCheckpoint(_ []byte, _ context.Context) {
}

// IsPruningEnabled -
func (m *MockAccountsAdapter) IsPruningEnabled() bool {
	return false
}

// IsInterfaceNil -
func (m *MockAccountsAdapter) IsInterfaceNil() bool {
	return m == nil
}
