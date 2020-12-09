package worldmock

// NewAddressMock allows tests to specify what new addresses to generate
type NewAddressMock struct {
	CreatorAddress []byte
	CreatorNonce   uint64
	NewAddress     []byte
}

// BlockInfo contains metadata about a mocked block
type BlockInfo struct {
	BlockTimestamp uint64
	BlockNonce     uint64
	BlockRound     uint64
	BlockEpoch     uint32
	RandomSeed     []byte
}

// MockWorld provides a mock representation of the blockchain to be used in VM tests.
type MockWorld struct {
	AcctMap                    AccountMap
	PreviousBlockInfo          *BlockInfo
	CurrentBlockInfo           *BlockInfo
	Blockhashes                [][]byte
	NewAddressMocks            []*NewAddressMock
	StateRootHash              []byte
	Err                        error
	LastCreatedContractAddress []byte
	CompiledCode               map[string][]byte
}

// NewMockWorld creates a new mock instance
func NewMockWorld() *MockWorld {
	return &MockWorld{
		AcctMap:           NewAccountMap(),
		PreviousBlockInfo: nil,
		CurrentBlockInfo:  nil,
		Blockhashes:       nil,
		NewAddressMocks:   nil,
		CompiledCode:      make(map[string][]byte),
	}
}

// Clear resets all mock data between tests.
func (b *MockWorld) Clear() {
	b.AcctMap = NewAccountMap()
	b.PreviousBlockInfo = nil
	b.CurrentBlockInfo = nil
	b.Blockhashes = nil
	b.NewAddressMocks = nil
	b.CompiledCode = make(map[string][]byte)
}

func (b *MockWorld) SetCurrentBlockHash(blockHash []byte) {
	if b.CurrentBlockInfo == nil {
		b.CurrentBlockInfo = &BlockInfo{}
	}
	b.Blockhashes = [][]byte{blockHash}
}
