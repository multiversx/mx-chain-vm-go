package nodepart

import "github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"

func (part *NodePart) replyToBlockchainAccountExists(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainAccountExistsRequest)
	result, err := part.blockchain.AccountExists(typedRequest.Address)
	response := common.NewMessageBlockchainAccountExistsResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainNewAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainNewAddressRequest)
	result, err := part.blockchain.NewAddress(typedRequest.CreatorAddress, typedRequest.CreatorNonce, typedRequest.VmType)
	response := common.NewMessageBlockchainNewAddressResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainGetBalance(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetBalanceRequest)
	balance, err := part.blockchain.GetBalance(typedRequest.Address)
	response := common.NewMessageBlockchainGetBalanceResponse(balance, err)
	return response
}

func (part *NodePart) replyToBlockchainGetNonce(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetNonceRequest)
	nonce, err := part.blockchain.GetNonce(typedRequest.Address)
	response := common.NewMessageBlockchainGetNonceResponse(nonce, err)
	return response
}

func (part *NodePart) replyToBlockchainGetStorageData(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetStorageDataRequest)
	data, err := part.blockchain.GetStorageData(typedRequest.Address, typedRequest.Index)
	response := common.NewMessageBlockchainGetStorageDataResponse(data, err)
	return response
}

func (part *NodePart) replyToBlockchainIsCodeEmpty(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainIsCodeEmptyRequest)
	result, err := part.blockchain.IsCodeEmpty(typedRequest.Address)
	response := common.NewMessageBlockchainIsCodeEmptyResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainGetCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetCodeRequest)
	code, err := part.blockchain.GetCode(typedRequest.Address)
	response := common.NewMessageBlockchainGetCodeResponse(code, err)
	return response
}

func (part *NodePart) replyToBlockchainGetBlockhash(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetBlockhashRequest)
	result, err := part.blockchain.GetBlockhash(typedRequest.Nonce)
	response := common.NewMessageBlockchainGetBlockhashResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainLastNonce(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastNonce()
	response := common.NewMessageBlockchainLastNonceResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastRound(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastRound()
	response := common.NewMessageBlockchainLastRoundResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastTimeStamp(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastTimeStamp()
	response := common.NewMessageBlockchainLastTimeStampResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastRandomSeed(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastRandomSeed()
	response := common.NewMessageBlockchainLastRandomSeedResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastEpoch(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastEpoch()
	response := common.NewMessageBlockchainLastEpochResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainGetStateRootHash(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.GetStateRootHash()
	response := common.NewMessageBlockchainGetStateRootHashResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentNonce(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentNonce()
	response := common.NewMessageBlockchainCurrentNonceResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentRound(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentRound()
	response := common.NewMessageBlockchainCurrentRoundResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentTimeStamp(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentTimeStamp()
	response := common.NewMessageBlockchainCurrentTimeStampResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentRandomSeed(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentRandomSeed()
	response := common.NewMessageBlockchainCurrentRandomSeedResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentEpoch(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentEpoch()
	response := common.NewMessageBlockchainCurrentEpochResponse(result)
	return response
}
