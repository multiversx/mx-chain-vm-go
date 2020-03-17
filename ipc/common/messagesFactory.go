package common

// CreateMessage creates a message given its kind
func CreateMessage(kind MessageKind) MessageHandler {
	kindIndex := uint32(kind)
	length := uint32(len(messageCreators))
	if kindIndex < length {
		message := messageCreators[kindIndex]()
		message.SetKind(kind)
		return message
	}

	return createUndefinedMessage()
}

type messageCreator func() MessageHandler

var messageCreators = make([]messageCreator, LastKind)

func init() {
	for i := 0; i < len(messageCreators); i++ {
		messageCreators[i] = createUndefinedMessage
	}

	messageCreators[Initialize] = createMessageInitialize
	messageCreators[Stop] = createMessageStop
	messageCreators[ContractDeployRequest] = createMessageContractDeployRequest
	messageCreators[ContractCallRequest] = createMessageContractCallRequest
	messageCreators[ContractResponse] = createMessageContractResponse
	messageCreators[DiagnoseWaitRequest] = createMessageDiagnoseWaitRequest
	messageCreators[DiagnoseWaitResponse] = createMessageDiagnoseWaitResponse

	messageCreators[BlockchainAccountExistsRequest] = createMessageBlockchainAccountExistsRequest
	messageCreators[BlockchainAccountExistsResponse] = createMessageBlockchainAccountExistsResponse
	messageCreators[BlockchainNewAddressRequest] = createMessageBlockchainNewAddressRequest
	messageCreators[BlockchainNewAddressResponse] = createMessageBlockchainNewAddressResponse
	messageCreators[BlockchainGetBalanceRequest] = createMessageBlockchainGetBalanceRequest
	messageCreators[BlockchainGetBalanceResponse] = createMessageBlockchainGetBalanceResponse
	messageCreators[BlockchainGetNonceRequest] = createMessageBlockchainGetNonceRequest
	messageCreators[BlockchainGetNonceResponse] = createMessageBlockchainGetNonceResponse
	messageCreators[BlockchainGetStorageDataRequest] = createMessageBlockchainGetStorageDataRequest
	messageCreators[BlockchainGetStorageDataResponse] = createMessageBlockchainGetStorageDataResponse
	messageCreators[BlockchainIsCodeEmptyRequest] = createMessageBlockchainIsCodeEmptyRequest
	messageCreators[BlockchainIsCodeEmptyResponse] = createMessageBlockchainIsCodeEmptyResponse
	messageCreators[BlockchainGetCodeRequest] = createMessageBlockchainGetCodeRequest
	messageCreators[BlockchainGetCodeResponse] = createMessageBlockchainGetCodeResponse
	messageCreators[BlockchainGetBlockhashRequest] = createMessageBlockchainGetBlockhashRequest
	messageCreators[BlockchainGetBlockhashResponse] = createMessageBlockchainGetBlockhashResponse
	messageCreators[BlockchainLastNonceRequest] = createMessageBlockchainLastNonceRequest
	messageCreators[BlockchainLastNonceResponse] = createMessageBlockchainLastNonceResponse
	messageCreators[BlockchainLastRoundRequest] = createMessageBlockchainLastRoundRequest
	messageCreators[BlockchainLastRoundResponse] = createMessageBlockchainLastRoundResponse
	messageCreators[BlockchainLastTimeStampRequest] = createMessageBlockchainLastTimeStampRequest
	messageCreators[BlockchainLastTimeStampResponse] = createMessageBlockchainLastTimeStampResponse
	messageCreators[BlockchainLastRandomSeedRequest] = createMessageBlockchainLastRandomSeedRequest
	messageCreators[BlockchainLastRandomSeedResponse] = createMessageBlockchainLastRandomSeedResponse
	messageCreators[BlockchainLastEpochRequest] = createMessageBlockchainLastEpochRequest
	messageCreators[BlockchainLastEpochResponse] = createMessageBlockchainLastEpochResponse
	messageCreators[BlockchainGetStateRootHashRequest] = createMessageBlockchainGetStateRootHashRequest
	messageCreators[BlockchainGetStateRootHashResponse] = createMessageBlockchainGetStateRootHashResponse
	messageCreators[BlockchainCurrentNonceRequest] = createMessageBlockchainCurrentNonceRequest
	messageCreators[BlockchainCurrentNonceResponse] = createMessageBlockchainCurrentNonceResponse
	messageCreators[BlockchainCurrentRoundRequest] = createMessageBlockchainCurrentRoundRequest
	messageCreators[BlockchainCurrentRoundResponse] = createMessageBlockchainCurrentRoundResponse
	messageCreators[BlockchainCurrentTimeStampRequest] = createMessageBlockchainCurrentTimeStampRequest
	messageCreators[BlockchainCurrentTimeStampResponse] = createMessageBlockchainCurrentTimeStampResponse
	messageCreators[BlockchainCurrentRandomSeedRequest] = createMessageBlockchainCurrentRandomSeedRequest
	messageCreators[BlockchainCurrentRandomSeedResponse] = createMessageBlockchainCurrentRandomSeedResponse
	messageCreators[BlockchainCurrentEpochRequest] = createMessageBlockchainCurrentEpochRequest
	messageCreators[BlockchainCurrentEpochResponse] = createMessageBlockchainCurrentEpochResponse

}

func createMessageInitialize() MessageHandler {
	return &MessageInitialize{}
}

func createMessageStop() MessageHandler {
	return &MessageStop{}
}

func createMessageContractDeployRequest() MessageHandler {
	return &MessageContractDeployRequest{}
}

func createMessageContractCallRequest() MessageHandler {
	return &MessageContractCallRequest{}
}

func createMessageContractResponse() MessageHandler {
	return &MessageContractResponse{}
}

func createMessageDiagnoseWaitRequest() MessageHandler {
	return &MessageDiagnoseWaitRequest{}
}

func createMessageDiagnoseWaitResponse() MessageHandler {
	return &MessageDiagnoseWaitResponse{}
}

func createUndefinedMessage() MessageHandler {
	return NewUndefinedMessage()
}

func createMessageBlockchainAccountExistsRequest() MessageHandler {
	return &MessageBlockchainAccountExistsRequest{}
}

func createMessageBlockchainAccountExistsResponse() MessageHandler {
	return &MessageBlockchainAccountExistsResponse{}
}

func createMessageBlockchainNewAddressRequest() MessageHandler {
	return &MessageBlockchainNewAddressRequest{}
}

func createMessageBlockchainNewAddressResponse() MessageHandler {
	return &MessageBlockchainNewAddressResponse{}
}

func createMessageBlockchainGetBalanceRequest() MessageHandler {
	return &MessageBlockchainGetBalanceRequest{}
}

func createMessageBlockchainGetBalanceResponse() MessageHandler {
	return &MessageBlockchainGetBalanceResponse{}
}

func createMessageBlockchainGetNonceRequest() MessageHandler {
	return &MessageBlockchainGetNonceRequest{}
}

func createMessageBlockchainGetNonceResponse() MessageHandler {
	return &MessageBlockchainGetNonceResponse{}
}

func createMessageBlockchainGetStorageDataRequest() MessageHandler {
	return &MessageBlockchainGetStorageDataRequest{}
}

func createMessageBlockchainGetStorageDataResponse() MessageHandler {
	return &MessageBlockchainGetStorageDataResponse{}
}

func createMessageBlockchainIsCodeEmptyRequest() MessageHandler {
	return &MessageBlockchainIsCodeEmptyRequest{}
}

func createMessageBlockchainIsCodeEmptyResponse() MessageHandler {
	return &MessageBlockchainIsCodeEmptyResponse{}
}

func createMessageBlockchainGetCodeRequest() MessageHandler {
	return &MessageBlockchainGetCodeRequest{}
}

func createMessageBlockchainGetCodeResponse() MessageHandler {
	return &MessageBlockchainGetCodeResponse{}
}

func createMessageBlockchainGetBlockhashRequest() MessageHandler {
	return &MessageBlockchainGetBlockhashRequest{}
}

func createMessageBlockchainGetBlockhashResponse() MessageHandler {
	return &MessageBlockchainGetBlockhashResponse{}
}

func createMessageBlockchainLastNonceRequest() MessageHandler {
	return &MessageBlockchainLastNonceRequest{}
}

func createMessageBlockchainLastNonceResponse() MessageHandler {
	return &MessageBlockchainLastNonceResponse{}
}

func createMessageBlockchainLastRoundRequest() MessageHandler {
	return &MessageBlockchainLastRoundRequest{}
}

func createMessageBlockchainLastRoundResponse() MessageHandler {
	return &MessageBlockchainLastRoundResponse{}
}

func createMessageBlockchainLastTimeStampRequest() MessageHandler {
	return &MessageBlockchainLastTimeStampRequest{}
}

func createMessageBlockchainLastTimeStampResponse() MessageHandler {
	return &MessageBlockchainLastTimeStampResponse{}
}

func createMessageBlockchainLastRandomSeedRequest() MessageHandler {
	return &MessageBlockchainLastRandomSeedRequest{}
}

func createMessageBlockchainLastRandomSeedResponse() MessageHandler {
	return &MessageBlockchainLastRandomSeedResponse{}
}

func createMessageBlockchainLastEpochRequest() MessageHandler {
	return &MessageBlockchainLastEpochRequest{}
}

func createMessageBlockchainLastEpochResponse() MessageHandler {
	return &MessageBlockchainLastEpochResponse{}
}

func createMessageBlockchainGetStateRootHashRequest() MessageHandler {
	return &MessageBlockchainGetStateRootHashRequest{}
}

func createMessageBlockchainGetStateRootHashResponse() MessageHandler {
	return &MessageBlockchainGetStateRootHashResponse{}
}

func createMessageBlockchainCurrentNonceRequest() MessageHandler {
	return &MessageBlockchainCurrentNonceRequest{}
}

func createMessageBlockchainCurrentNonceResponse() MessageHandler {
	return &MessageBlockchainCurrentNonceResponse{}
}

func createMessageBlockchainCurrentRoundRequest() MessageHandler {
	return &MessageBlockchainCurrentRoundRequest{}
}

func createMessageBlockchainCurrentRoundResponse() MessageHandler {
	return &MessageBlockchainCurrentRoundResponse{}
}

func createMessageBlockchainCurrentTimeStampRequest() MessageHandler {
	return &MessageBlockchainCurrentTimeStampRequest{}
}

func createMessageBlockchainCurrentTimeStampResponse() MessageHandler {
	return &MessageBlockchainCurrentTimeStampResponse{}
}

func createMessageBlockchainCurrentRandomSeedRequest() MessageHandler {
	return &MessageBlockchainCurrentRandomSeedRequest{}
}

func createMessageBlockchainCurrentRandomSeedResponse() MessageHandler {
	return &MessageBlockchainCurrentRandomSeedResponse{}
}

func createMessageBlockchainCurrentEpochRequest() MessageHandler {
	return &MessageBlockchainCurrentEpochRequest{}
}

func createMessageBlockchainCurrentEpochResponse() MessageHandler {
	return &MessageBlockchainCurrentEpochResponse{}
}
