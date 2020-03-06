package nodepart

import "github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"

func (part *NodePart) replyToBlockchainNewAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainNewAddressRequest)
	address, err := part.blockchain.NewAddress(typedRequest.CreatorAddress, typedRequest.CreatorNonce, typedRequest.VMType)
	response := common.NewMessageBlockchainNewAddressResponse(err)
	response.Address = address
	return response
}

func (part *NodePart) replyToBlockchainGetNonce(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetNonceRequest)
	nonce, err := part.blockchain.GetNonce(typedRequest.Address)
	response := common.NewMessageBlockchainGetNonceResponse(err)
	response.Nonce = nonce
	return response
}

func (part *NodePart) replyToBlockchainGetStorageData(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetStorageDataRequest)
	data, err := part.blockchain.GetStorageData(typedRequest.Address, typedRequest.Index)
	response := common.NewMessageBlockchainGetStorageDataResponse(err)
	response.Data = data
	return response
}

func (part *NodePart) replyToBlockchainGetCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetCodeRequest)
	code, err := part.blockchain.GetCode(typedRequest.Address)
	response := common.NewMessageBlockchainGetCodeResponse(err)
	response.Code = code
	return response
}
