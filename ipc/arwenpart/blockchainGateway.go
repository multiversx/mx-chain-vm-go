package arwenpart

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookGateway)(nil)

// BlockchainHookGateway forwards requests to the actual hook
type BlockchainHookGateway struct {
	messenger *ChildMessenger
}

// NewBlockchainHookGateway creates a new gateway
func NewBlockchainHookGateway(messenger *ChildMessenger) *BlockchainHookGateway {
	return &BlockchainHookGateway{messenger: messenger}
}

// AccountExists forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) AccountExists(address []byte) (bool, error) {
	request := common.NewMessageBlockchainAccountExistsRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return false, err
	}

	response := rawResponse.(*common.MessageBlockchainAccountExistsResponse)
	return response.Result, response.GetError()
}

// NewAddress forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	request := common.NewMessageBlockchainNewAddressRequest(creatorAddress, creatorNonce, vmType)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainNewAddressResponse)
	return response.Result, response.GetError()
}

// GetBalance forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetBalance(address []byte) (*big.Int, error) {
	request := common.NewMessageBlockchainGetBalanceRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainGetBalanceResponse)
	return response.Balance, response.GetError()
}

// GetNonce forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetNonce(address []byte) (uint64, error) {
	request := common.NewMessageBlockchainGetNonceRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0, err
	}

	response := rawResponse.(*common.MessageBlockchainGetNonceResponse)
	return response.Nonce, response.GetError()
}

// GetStorageData forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetStorageData(address []byte, index []byte) ([]byte, error) {
	request := common.NewMessageBlockchainGetStorageDataRequest(address, index)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainGetStorageDataResponse)
	return response.Data, response.GetError()
}

// IsCodeEmpty forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) IsCodeEmpty(address []byte) (bool, error) {
	request := common.NewMessageBlockchainIsCodeEmptyRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return false, err
	}

	response := rawResponse.(*common.MessageBlockchainIsCodeEmptyResponse)
	return response.Result, response.GetError()
}

// GetCode forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetCode(address []byte) ([]byte, error) {
	request := common.NewMessageBlockchainGetCodeRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainGetCodeResponse)
	return response.Code, response.GetError()
}

// GetBlockhash forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetBlockhash(nonce uint64) ([]byte, error) {
	request := common.NewMessageBlockchainGetBlockhashRequest(nonce)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*common.MessageBlockchainGetBlockhashResponse)
	return response.Result, response.GetError()
}

// LastNonce forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastNonce() uint64 {
	request := common.NewMessageBlockchainLastNonceRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastNonceResponse)
	return response.Result
}

// LastRound forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastRound() uint64 {
	request := common.NewMessageBlockchainLastRoundRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastRoundResponse)
	return response.Result
}

// LastTimeStamp forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastTimeStamp() uint64 {
	request := common.NewMessageBlockchainLastTimeStampRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastTimeStampResponse)
	return response.Result
}

// LastRandomSeed forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastRandomSeed() []byte {
	request := common.NewMessageBlockchainLastRandomSeedRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil
	}

	response := rawResponse.(*common.MessageBlockchainLastRandomSeedResponse)
	return response.Result
}

// LastEpoch forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) LastEpoch() uint32 {
	request := common.NewMessageBlockchainLastEpochRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainLastEpochResponse)
	return response.Result
}

// GetStateRootHash forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetStateRootHash() []byte {
	request := common.NewMessageBlockchainGetStateRootHashRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil
	}

	response := rawResponse.(*common.MessageBlockchainGetStateRootHashResponse)
	return response.Result
}

// CurrentNonce forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentNonce() uint64 {
	request := common.NewMessageBlockchainCurrentNonceRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentNonceResponse)
	return response.Result
}

// CurrentRound forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentRound() uint64 {
	request := common.NewMessageBlockchainCurrentRoundRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentRoundResponse)
	return response.Result
}

// CurrentTimeStamp forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentTimeStamp() uint64 {
	request := common.NewMessageBlockchainCurrentTimeStampRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentTimeStampResponse)
	return response.Result
}

// CurrentRandomSeed forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentRandomSeed() []byte {
	request := common.NewMessageBlockchainCurrentRandomSeedRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil
	}

	response := rawResponse.(*common.MessageBlockchainCurrentRandomSeedResponse)
	return response.Result
}

// CurrentEpoch forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) CurrentEpoch() uint32 {
	request := common.NewMessageBlockchainCurrentEpochRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentEpochResponse)
	return response.Result
}
