package common

import (
	"fmt"
)

// CreateMessage creates a message given its kind
func CreateMessage(kind MessageKind) MessageHandler {
	var message MessageHandler

	switch kind {
	case Stop:
		message = &MessageStop{}
	case ContractDeployRequest:
		message = &MessageContractDeployRequest{}
	case ContractCallRequest:
		message = &MessageContractCallRequest{}
	case ContractResponse:
		message = &MessageContractResponse{}
	case BlockchainAccountExistsRequest:
		message = &MessageBlockchainAccountExistsRequest{}
	case BlockchainAccountExistsResponse:
		message = &MessageBlockchainAccountExistsResponse{}
	case BlockchainNewAddressRequest:
		message = &MessageBlockchainNewAddressRequest{}
	case BlockchainNewAddressResponse:
		message = &MessageBlockchainNewAddressResponse{}
	case BlockchainGetBalanceRequest:
		message = &MessageBlockchainGetBalanceRequest{}
	case BlockchainGetBalanceResponse:
		message = &MessageBlockchainGetBalanceResponse{}
	case BlockchainGetNonceRequest:
		message = &MessageBlockchainGetNonceRequest{}
	case BlockchainGetNonceResponse:
		message = &MessageBlockchainGetNonceResponse{}
	case BlockchainGetStorageDataRequest:
		message = &MessageBlockchainGetStorageDataRequest{}
	case BlockchainGetStorageDataResponse:
		message = &MessageBlockchainGetStorageDataResponse{}
	case BlockchainIsCodeEmptyRequest:
		message = &MessageBlockchainIsCodeEmptyRequest{}
	case BlockchainIsCodeEmptyResponse:
		message = &MessageBlockchainIsCodeEmptyResponse{}
	case BlockchainGetCodeRequest:
		message = &MessageBlockchainGetCodeRequest{}
	case BlockchainGetCodeResponse:
		message = &MessageBlockchainGetCodeResponse{}
	case BlockchainGetBlockhashRequest:
		message = &MessageBlockchainGetBlockhashRequest{}
	case BlockchainGetBlockhashResponse:
		message = &MessageBlockchainGetBlockhashResponse{}
	case BlockchainLastNonceRequest:
		message = &MessageBlockchainLastNonceRequest{}
	case BlockchainLastNonceResponse:
		message = &MessageBlockchainLastNonceResponse{}
	case BlockchainLastRoundRequest:
		message = &MessageBlockchainLastRoundRequest{}
	case BlockchainLastRoundResponse:
		message = &MessageBlockchainLastRoundResponse{}
	case BlockchainLastTimeStampRequest:
		message = &MessageBlockchainLastTimeStampRequest{}
	case BlockchainLastTimeStampResponse:
		message = &MessageBlockchainLastTimeStampResponse{}
	case BlockchainLastRandomSeedRequest:
		message = &MessageBlockchainLastRandomSeedRequest{}
	case BlockchainLastRandomSeedResponse:
		message = &MessageBlockchainLastRandomSeedResponse{}
	case BlockchainLastEpochRequest:
		message = &MessageBlockchainLastEpochRequest{}
	case BlockchainLastEpochResponse:
		message = &MessageBlockchainLastEpochResponse{}
	case BlockchainGetStateRootHashRequest:
		message = &MessageBlockchainGetStateRootHashRequest{}
	case BlockchainGetStateRootHashResponse:
		message = &MessageBlockchainGetStateRootHashResponse{}
	case BlockchainCurrentNonceRequest:
		message = &MessageBlockchainCurrentNonceRequest{}
	case BlockchainCurrentNonceResponse:
		message = &MessageBlockchainCurrentNonceResponse{}
	case BlockchainCurrentRoundRequest:
		message = &MessageBlockchainCurrentRoundRequest{}
	case BlockchainCurrentRoundResponse:
		message = &MessageBlockchainCurrentRoundResponse{}
	case BlockchainCurrentTimeStampRequest:
		message = &MessageBlockchainCurrentTimeStampRequest{}
	case BlockchainCurrentTimeStampResponse:
		message = &MessageBlockchainCurrentTimeStampResponse{}
	case BlockchainCurrentRandomSeedRequest:
		message = &MessageBlockchainCurrentRandomSeedRequest{}
	case BlockchainCurrentRandomSeedResponse:
		message = &MessageBlockchainCurrentRandomSeedResponse{}
	case BlockchainCurrentEpochRequest:
		message = &MessageBlockchainCurrentEpochRequest{}
	case BlockchainCurrentEpochResponse:
		message = &MessageBlockchainCurrentEpochResponse{}
	case DiagnoseWaitRequest:
		message = &MessageDiagnoseWaitRequest{}
	case DiagnoseWaitResponse:
		message = &MessageDiagnoseWaitResponse{}
	default:
		panic(fmt.Sprintf("Unknown message kind [%d]", kind))
	}

	message.SetKind(kind)
	return message
}
