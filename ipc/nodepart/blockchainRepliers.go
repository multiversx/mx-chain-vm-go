package nodepart

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/ipc/common"
)

func (part *NodePart) replyToBlockchainNewAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainNewAddressRequest)
	result, err := part.blockchain.NewAddress(typedRequest.CreatorAddress, typedRequest.CreatorNonce, typedRequest.VmType)
	response := common.NewMessageBlockchainNewAddressResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainGetStorageData(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetStorageDataRequest)
	data, err := part.blockchain.GetStorageData(typedRequest.AccountAddress, typedRequest.Index)
	response := common.NewMessageBlockchainGetStorageDataResponse(data, err)
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

func (part *NodePart) replyToBlockchainProcessBuiltInFunction(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainProcessBuiltInFunctionRequest)
	vmOutput, err := part.blockchain.ProcessBuiltInFunction(typedRequest.Input)
	response := common.NewMessageBlockchainProcessBuiltInFunctionResponse(vmOutput, err)
	return response
}

func (part *NodePart) replyToBlockchainGetBuiltinFunctionNames(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.GetBuiltinFunctionNames()
	response := common.NewMessageBlockchainGetBuiltinFunctionNamesResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainGetAllState(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetAllStateRequest)
	result, err := part.blockchain.GetAllState(typedRequest.Address)
	response := common.NewMessageBlockchainGetAllStateResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainGetUserAccount(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetUserAccountRequest)
	result, err := part.blockchain.GetUserAccount(typedRequest.Address)
	response := common.NewMessageBlockchainGetUserAccountResponse(&common.Account{
		Nonce:           result.GetNonce(),
		Balance:         result.GetBalance(),
		CodeHash:        result.GetCodeHash(),
		RootHash:        result.GetRootHash(),
		Address:         result.AddressBytes(),
		DeveloperReward: result.GetDeveloperReward(),
		OwnerAddress:    result.GetOwnerAddress(),
		UserName:        result.GetUserName(),
		CodeMetadata:    result.GetCodeMetadata(),
	}, err)
	return response
}

func (part *NodePart) replyToBlockchainGetCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetCodeRequest)
	code := part.blockchain.GetCode(typedRequest.Account)
	response := common.NewMessageBlockchainGetCodeResponse(code)
	return response
}

func (part *NodePart) replyToBlockchainGetShardOfAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetShardOfAddressRequest)
	result := part.blockchain.GetShardOfAddress(typedRequest.Address)
	response := common.NewMessageBlockchainGetShardOfAddressResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainIsSmartContract(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainIsSmartContractRequest)
	result := part.blockchain.IsSmartContract(typedRequest.Address)
	response := common.NewMessageBlockchainIsSmartContractResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainIsPayable(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainIsPayableRequest)
	result, err := part.blockchain.IsPayable(typedRequest.Address)
	response := common.NewMessageBlockchainIsPayableResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainSaveCompiledCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainSaveCompiledCodeRequest)
	part.blockchain.SaveCompiledCode(typedRequest.CodeHash, typedRequest.Code)
	response := common.NewMessageBlockchainSaveCompiledCodeResponse()
	return response
}

func (part *NodePart) replyToBlockchainGetCompiledCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetCompiledCodeRequest)
	found, code := part.blockchain.GetCompiledCode(typedRequest.CodeHash)
	response := common.NewMessageBlockchainGetCompiledCodeResponse(found, code)
	return response
}

func (part *NodePart) replyToBlockchainClearCompiledCodes(request common.MessageHandler) common.MessageHandler {
	part.blockchain.ClearCompiledCodes()
	response := common.NewMessageBlockchainClearCompiledCodesResponse()
	return response
}

func (part *NodePart) replyToBlockchainGetESDTToken(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetESDTTokenRequest)
	result, err := part.blockchain.GetESDTToken(typedRequest.Address, typedRequest.TokenID, typedRequest.Nonce)
	response := common.NewMessageBlockchainGetESDTTokenResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainIsInterfaceNil(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.IsInterfaceNil()
	response := common.NewMessageBlockchainIsInterfaceNilResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainGetSnapshot(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.GetSnapshot()
	response := common.NewMessageBlockchainGetSnapshotResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainRevertToSnapshot(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainRevertToSnapshotRequest)
	err := part.blockchain.RevertToSnapshot(typedRequest.Snapshot)
	response := common.NewMessageBlockchainRevertToSnapshotResponse(err)
	return response
}
