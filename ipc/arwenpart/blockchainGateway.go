package arwenpart

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookGateway)(nil)

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
	common.LogError("not implemented: AccountExists")
	return false, nil
}

// NewAddress forwards
func (blockchain *BlockchainHookGateway) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	request := common.NewMessageBlockchainNewAddressRequest(creatorAddress, creatorNonce, vmType)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainNewAddressResponse)
	return response.Result, response.GetError()
}

// GetBalance forwards
func (blockchain *BlockchainHookGateway) GetBalance(address []byte) (*big.Int, error) {
	common.LogError("not implemented: GetBalance")
	return nil, nil
}

// GetNonce forwards
func (blockchain *BlockchainHookGateway) GetNonce(address []byte) (uint64, error) {
	request := common.NewMessageBlockchainGetNonceRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0, err
	}

	response := rawResponse.(*common.MessageBlockchainGetNonceResponse)
	return response.Nonce, response.GetError()
}

// GetStorageData forwards
func (blockchain *BlockchainHookGateway) GetStorageData(accountAddress []byte, index []byte) ([]byte, error) {
	request := common.NewMessageBlockchainGetStorageDataRequest(accountAddress, index)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainGetStorageDataResponse)
	return response.Data, response.GetError()
}

// IsCodeEmpty forwards
func (blockchain *BlockchainHookGateway) IsCodeEmpty(address []byte) (bool, error) {
	common.LogError("not implemented: IsCodeEmpty")
	return false, nil
}

// GetCode forwards
func (blockchain *BlockchainHookGateway) GetCode(address []byte) ([]byte, error) {
	request := common.NewMessageBlockchainGetCodeRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainGetCodeResponse)
	return response.Code, response.GetError()
}

// GetBlockhash forwards
func (blockchain *BlockchainHookGateway) GetBlockhash(nonce uint64) ([]byte, error) {
	common.LogError("not implemented: GetBlockhash")
	return nil, nil
}

// LastNonce forwards
func (blockchain *BlockchainHookGateway) LastNonce() uint64 {
	common.LogError("not implemented: LastNonce")
	return 0
}

// LastRound forwards
func (blockchain *BlockchainHookGateway) LastRound() uint64 {
	common.LogError("not implemented: LastRound")
	return 0
}

// LastTimeStamp forwards
func (blockchain *BlockchainHookGateway) LastTimeStamp() uint64 {
	common.LogError("not implemented: LastTimeStamp")
	return 0
}

// LastRandomSeed forwards
func (blockchain *BlockchainHookGateway) LastRandomSeed() []byte {
	common.LogError("not implemented: LastRandomSeed")
	return nil
}

// LastEpoch forwards
func (blockchain *BlockchainHookGateway) LastEpoch() uint32 {
	common.LogError("not implemented: LastEpoch")
	return 0
}

// GetStateRootHash forwards
func (blockchain *BlockchainHookGateway) GetStateRootHash() []byte {
	common.LogError("not implemented: GetStateRootHash")
	return nil
}

// CurrentNonce forwards
func (blockchain *BlockchainHookGateway) CurrentNonce() uint64 {
	common.LogError("not implemented: CurrentNonce")
	return 0
}

// CurrentRound forwards
func (blockchain *BlockchainHookGateway) CurrentRound() uint64 {
	common.LogError("not implemented: CurrentRound")
	return 0
}

// CurrentTimeStamp forwards
func (blockchain *BlockchainHookGateway) CurrentTimeStamp() uint64 {
	common.LogError("not implemented: CurrentTimeStamp")
	return 0
}

// CurrentRandomSeed forwards
func (blockchain *BlockchainHookGateway) CurrentRandomSeed() []byte {
	common.LogError("not implemented: CurrentRandomSeed")
	return nil
}

// CurrentEpoch forwards
func (blockchain *BlockchainHookGateway) CurrentEpoch() uint32 {
	common.LogError("not implemented: CurrentEpoch")
	return 0
}
