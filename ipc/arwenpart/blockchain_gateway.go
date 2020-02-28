package arwenpart

import (
	"log"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
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
	log.Fatal("not implemented")
	return false, nil
}

// NewAddress forwards
func (blockchain *BlockchainHookGateway) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	request := common.NewHookCallRequest("blockchain", "NewAddress")
	request.Bytes1 = creatorAddress
	request.Uint64_1 = creatorNonce
	request.Bytes2 = vmType
	results, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	return results.Bytes1, nil
}

// GetBalance forwards
func (blockchain *BlockchainHookGateway) GetBalance(address []byte) (*big.Int, error) {
	log.Fatal("not implemented")
	return nil, nil
}

// GetNonce forwards
func (blockchain *BlockchainHookGateway) GetNonce(address []byte) (uint64, error) {
	log.Fatal("not implemented")
	return 0, nil
}

// GetStorageData forwards
func (blockchain *BlockchainHookGateway) GetStorageData(accountAddress []byte, index []byte) ([]byte, error) {
	log.Fatal("not implemented")
	return nil, nil
}

// IsCodeEmpty forwards
func (blockchain *BlockchainHookGateway) IsCodeEmpty(address []byte) (bool, error) {
	log.Fatal("not implemented")
	return false, nil
}

// GetCode forwards
func (blockchain *BlockchainHookGateway) GetCode(address []byte) ([]byte, error) {
	request := common.NewHookCallRequest("blockchain", "GetCode")
	request.Bytes1 = address
	results, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	return results.Bytes1, nil
}

// GetBlockhash forwards
func (blockchain *BlockchainHookGateway) GetBlockhash(nonce uint64) ([]byte, error) {
	log.Fatal("not implemented")
	return nil, nil
}

// LastNonce forwards
func (blockchain *BlockchainHookGateway) LastNonce() uint64 {
	log.Fatal("not implemented")
	return 0
}

// LastRound forwards
func (blockchain *BlockchainHookGateway) LastRound() uint64 {
	log.Fatal("not implemented")
	return 0
}

// LastTimeStamp forwards
func (blockchain *BlockchainHookGateway) LastTimeStamp() uint64 {
	log.Fatal("not implemented")
	return 0
}

// LastRandomSeed forwards
func (blockchain *BlockchainHookGateway) LastRandomSeed() []byte {
	log.Fatal("not implemented")
	return nil
}

// LastEpoch forwards
func (blockchain *BlockchainHookGateway) LastEpoch() uint32 {
	log.Fatal("not implemented")
	return 0
}

// GetStateRootHash forwards
func (blockchain *BlockchainHookGateway) GetStateRootHash() []byte {
	log.Fatal("not implemented")
	return nil
}

// CurrentNonce forwards
func (blockchain *BlockchainHookGateway) CurrentNonce() uint64 {
	log.Fatal("not implemented")
	return 0
}

// CurrentRound forwards
func (blockchain *BlockchainHookGateway) CurrentRound() uint64 {
	log.Fatal("not implemented")
	return 0
}

// CurrentTimeStamp forwards
func (blockchain *BlockchainHookGateway) CurrentTimeStamp() uint64 {
	log.Fatal("not implemented")
	return 0
}

// CurrentRandomSeed forwards
func (blockchain *BlockchainHookGateway) CurrentRandomSeed() []byte {
	log.Fatal("not implemented")
	return nil
}

// CurrentEpoch forwards
func (blockchain *BlockchainHookGateway) CurrentEpoch() uint32 {
	log.Fatal("not implemented")
	return 0
}
