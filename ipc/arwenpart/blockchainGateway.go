package arwenpart

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookGateway)(nil)

// BlockchainHookGateway forwards requests to the actual hook
type BlockchainHookGateway struct {
	messenger *ArwenMessenger
}

// NewBlockchainHookGateway creates a new gateway
func NewBlockchainHookGateway(messenger *ArwenMessenger) *BlockchainHookGateway {
	return &BlockchainHookGateway{messenger: messenger}
}

// AccountExists forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) AccountExists(address []byte) (bool, error) {
	request := common.NewMessageBlockchainAccountExistsRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return false, err
	}

	if rawResponse.GetKind() != common.BlockchainAccountExistsResponse {
		return false, common.ErrBadHookResponseFromNode
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

	if rawResponse.GetKind() != common.BlockchainNewAddressResponse {
		return nil, common.ErrBadHookResponseFromNode
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

	if rawResponse.GetKind() != common.BlockchainGetBalanceResponse {
		return nil, common.ErrBadHookResponseFromNode
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

	if rawResponse.GetKind() != common.BlockchainGetNonceResponse {
		return 0, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainGetNonceResponse)
	return response.Nonce, response.GetError()
}

// GetStorageFullData forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetStorageData(address []byte, index []byte) ([]byte, error) {
	request := common.NewMessageBlockchainGetStorageDataRequest(address, index)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	storageData := make([]byte, 0)
	if err != nil {
		return storageData, err
	}

	if rawResponse.GetKind() != common.BlockchainGetStorageDataResponse {
		return storageData, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainGetStorageDataResponse)
	storageData = response.Data
	return storageData, response.GetError()
}

// IsCodeEmpty forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) IsCodeEmpty(address []byte) (bool, error) {
	request := common.NewMessageBlockchainIsCodeEmptyRequest(address)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return false, err
	}

	if rawResponse.GetKind() != common.BlockchainIsCodeEmptyResponse {
		return false, common.ErrBadHookResponseFromNode
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

	if rawResponse.GetKind() != common.BlockchainGetCodeResponse {
		return nil, common.ErrBadHookResponseFromNode
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

	if rawResponse.GetKind() != common.BlockchainGetBlockhashResponse {
		return nil, common.ErrBadHookResponseFromNode
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

	if rawResponse.GetKind() != common.BlockchainLastNonceResponse {
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

	if rawResponse.GetKind() != common.BlockchainLastRoundResponse {
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

	if rawResponse.GetKind() != common.BlockchainLastTimeStampResponse {
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

	if rawResponse.GetKind() != common.BlockchainLastRandomSeedResponse {
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

	if rawResponse.GetKind() != common.BlockchainLastEpochResponse {
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

	if rawResponse.GetKind() != common.BlockchainGetStateRootHashResponse {
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

	if rawResponse.GetKind() != common.BlockchainCurrentNonceResponse {
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

	if rawResponse.GetKind() != common.BlockchainCurrentRoundResponse {
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

	if rawResponse.GetKind() != common.BlockchainCurrentTimeStampResponse {
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

	if rawResponse.GetKind() != common.BlockchainCurrentRandomSeedResponse {
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

	if rawResponse.GetKind() != common.BlockchainCurrentEpochResponse {
		return 0
	}

	response := rawResponse.(*common.MessageBlockchainCurrentEpochResponse)
	return response.Result
}

// ProcessBuiltInFunction forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	request := common.NewMessageBlockchainProcessBuiltinFunctionRequest(*input)
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return nil, err
	}

	if rawResponse.GetKind() != common.BlockchainProcessBuiltinFunctionResponse {
		return nil, common.ErrBadHookResponseFromNode
	}

	response := rawResponse.(*common.MessageBlockchainProcessBuiltinFunctionResponse)
	return response.VMOutput, response.GetError()
}

// GetBuiltinFunctionNames forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	request := common.NewMessageBlockchainGetBuiltinFunctionNamesRequest()
	rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
	if err != nil {
		return make(vmcommon.FunctionNames)
	}

	if rawResponse.GetKind() != common.BlockchainGetBuiltinFunctionNamesResponse {
		return make(vmcommon.FunctionNames)
	}

	response := rawResponse.(*common.MessageBlockchainGetBuiltinFunctionNamesResponse)
	return response.FunctionNames
}

// GetAllState forwards a message to the actual hook
func (blockchain *BlockchainHookGateway) GetAllState(address []byte) (map[string][]byte, error) {
	//TODO implement this
	return nil, nil
}
