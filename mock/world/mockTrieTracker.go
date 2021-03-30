package worldmock

import (
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/state"
)

var _ state.DataTrieTracker = (*MockTrieTracker)(nil)

type MockTrieTracker struct {
	TrieData MockTrieData
}

func NewMockTrieTracker(data MockTrieData) *MockTrieTracker {
	return &MockTrieTracker{
		TrieData: data,
	}
}

func (tt *MockTrieTracker) RetrieveValue(key []byte) ([]byte, error) {
	return tt.TrieData[string(key)], nil
}

func (tt *MockTrieTracker) SaveKeyValue(key []byte, value []byte) error {
	tt.TrieData[string(key)] = value
	return nil
}

func (tt *MockTrieTracker) SetDataTrie(tr data.Trie) {
}

func (tt *MockTrieTracker) DataTrie() data.Trie {
	return nil
}

func (tt *MockTrieTracker) ClearDataCaches() {
}

func (tt *MockTrieTracker) DirtyData() MockTrieData {
	return tt.TrieData
}

func (tt *MockTrieTracker) IsInterfaceNil() bool {
	return tt == nil
}
