package arwenpart

import (
	"math/big"
)

// BlockchainHookGateway is
type BlockchainHookGateway struct {
	messenger *ChildMessenger
}

// NewBlockchainHookGateway creates
func NewBlockchainHookGateway(messenger *ChildMessenger) *BlockchainHookGateway {
	return &BlockchainHookGateway{messenger: messenger}
}

// AccountExists forwards
func (blockchain *BlockchainHookGateway) AccountExists(address []byte) (bool, error) {
	return false, nil
}

// NewAddress forwards
func (blockchain *BlockchainHookGateway) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	return nil, nil
}

// GetBalance forwards
func (blockchain *BlockchainHookGateway) GetBalance(address []byte) (*big.Int, error) {
	return nil, nil
}

// GetNonce forwards
func (blockchain *BlockchainHookGateway) GetNonce(address []byte) (uint64, error) {
	return 0, nil
}

// GetStorageData forwards
func (blockchain *BlockchainHookGateway) GetStorageData(accountAddress []byte, index []byte) ([]byte, error) {
	return nil, nil
}

// IsCodeEmpty forwards
func (blockchain *BlockchainHookGateway) IsCodeEmpty(address []byte) (bool, error) {
	return false, nil
}

// GetCode forwards
func (blockchain *BlockchainHookGateway) GetCode(address []byte) ([]byte, error) {
	return nil, nil
}

// GetBlockhash forwards
func (blockchain *BlockchainHookGateway) GetBlockhash(nonce uint64) ([]byte, error) {
	return nil, nil
}

// LastNonce forwards
func (blockchain *BlockchainHookGateway) LastNonce() uint64 {
	return 0
}

// LastRound forwards
func (blockchain *BlockchainHookGateway) LastRound() uint64 {
	return 0
}

// LastTimeStamp forwards
func (blockchain *BlockchainHookGateway) LastTimeStamp() uint64 {
	return 0
}

// LastRandomSeed forwards
func (blockchain *BlockchainHookGateway) LastRandomSeed() []byte { return nil }

// LastEpoch forwards
func (blockchain *BlockchainHookGateway) LastEpoch() uint32 { return 0 }

// GetStateRootHash forwards
func (blockchain *BlockchainHookGateway) GetStateRootHash() []byte { return nil }

// CurrentNonce forwards
func (blockchain *BlockchainHookGateway) CurrentNonce() uint64 { return 0 }

// CurrentRound forwards
func (blockchain *BlockchainHookGateway) CurrentRound() uint64 { return 0 }

// CurrentTimeStamp forwards
func (blockchain *BlockchainHookGateway) CurrentTimeStamp() uint64 { return 0 }

// CurrentRandomSeed forwards
func (blockchain *BlockchainHookGateway) CurrentRandomSeed() []byte { return nil }

// CurrentEpoch forwards
func (blockchain *BlockchainHookGateway) CurrentEpoch() uint32 { return 0 }
