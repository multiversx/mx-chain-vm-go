package callbackblockchain

// NewAddressMock allows tests to specify what new addresses to generate
type NewAddressMock struct {
	CreatorAddress []byte
	CreatorNonce   uint64
	NewAddress     []byte
}

// BlockInfo contains mock data about the corent block
type BlockInfo struct {
	BlockTimestamp uint64
	BlockNonce     uint64
	BlockRound     uint64
	BlockEpoch     uint32
}

// BlockchainHookMock provides a mock representation of the blockchain to be used in VM tests.
type BlockchainHookMock struct {
	AcctMap                      AccountMap
	PreviousBlockInfo            *BlockInfo
	CurrentBlockInfo             *BlockInfo
	Blockhashes                  [][]byte
	mockAddressGenerationEnabled bool
	NewAddressMocks              []*NewAddressMock
}

// NewMock creates a new mock instance
func NewMock() *BlockchainHookMock {
	return &BlockchainHookMock{
		AcctMap:                      NewAccountMap(),
		PreviousBlockInfo:            nil,
		CurrentBlockInfo:             nil,
		Blockhashes:                  nil,
		mockAddressGenerationEnabled: false,
	}
}

// Clear resets all mock data between tests.
func (b *BlockchainHookMock) Clear() {
	b.AcctMap = NewAccountMap()
	b.Blockhashes = nil
}

// EnableMockAddressGeneration causes the mock to generate its own new addresses.
func (b *BlockchainHookMock) EnableMockAddressGeneration() {
	b.mockAddressGenerationEnabled = true
}
