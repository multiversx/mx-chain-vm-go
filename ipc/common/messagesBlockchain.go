package common

// MessageBlockchainAccountExistsRequest represents a message
type MessageBlockchainAccountExistsRequest struct {
	Message
}

// MessageBlockchainAccountExistsResponse represents a message
type MessageBlockchainAccountExistsResponse struct {
	Message
}

// MessageBlockchainNewAddressRequest represents a message
type MessageBlockchainNewAddressRequest struct {
	Message
	CreatorAddress []byte
	CreatorNonce   uint64
	VMType         []byte
}

// NewMessageBlockchainNewAddressRequest creates a message
func NewMessageBlockchainNewAddressRequest() *MessageBlockchainNewAddressRequest {
	message := &MessageBlockchainNewAddressRequest{}
	message.Kind = BlockchainNewAddressRequest
	return message
}

// MessageBlockchainNewAddressResponse represents a message
type MessageBlockchainNewAddressResponse struct {
	Message
	Address []byte
}

// NewMessageBlockchainNewAddressResponse creates a message
func NewMessageBlockchainNewAddressResponse(err error) *MessageBlockchainNewAddressResponse {
	message := &MessageBlockchainNewAddressResponse{}
	message.Kind = BlockchainNewAddressResponse
	message.SetError(err)
	return message
}

// MessageBlockchainGetBalanceRequest represents a message
type MessageBlockchainGetBalanceRequest struct {
	Message
}

// MessageBlockchainGetBalanceResponse represents a message
type MessageBlockchainGetBalanceResponse struct {
	Message
}

// MessageBlockchainGetNonceRequest represents a message
type MessageBlockchainGetNonceRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetNonceRequest creates a message
func NewMessageBlockchainGetNonceRequest() *MessageBlockchainGetNonceRequest {
	message := &MessageBlockchainGetNonceRequest{}
	message.Kind = BlockchainGetNonceRequest
	return message
}

// MessageBlockchainGetNonceResponse represents a message
type MessageBlockchainGetNonceResponse struct {
	Message
	Nonce uint64
}

// NewMessageBlockchainGetNonceResponse creates a message
func NewMessageBlockchainGetNonceResponse(err error) *MessageBlockchainGetNonceResponse {
	message := &MessageBlockchainGetNonceResponse{}
	message.Kind = BlockchainGetNonceResponse
	message.SetError(err)
	return message
}

// MessageBlockchainGetStorageDataRequest represents a message
type MessageBlockchainGetStorageDataRequest struct {
	Message
	Address []byte
	Index   []byte
}

// NewMessageBlockchainGetStorageDataRequest creates a message
func NewMessageBlockchainGetStorageDataRequest() *MessageBlockchainGetStorageDataRequest {
	message := &MessageBlockchainGetStorageDataRequest{}
	message.Kind = BlockchainGetStorageDataRequest
	return message
}

// MessageBlockchainGetStorageDataResponse represents a message
type MessageBlockchainGetStorageDataResponse struct {
	Message
	Data []byte
}

// NewMessageBlockchainGetStorageDataResponse creates a message
func NewMessageBlockchainGetStorageDataResponse(err error) *MessageBlockchainGetStorageDataResponse {
	message := &MessageBlockchainGetStorageDataResponse{}
	message.Kind = BlockchainGetStorageDataResponse
	message.SetError(err)
	return message
}

// MessageBlockchainIsCodeEmptyRequest represents a message
type MessageBlockchainIsCodeEmptyRequest struct {
	Message
}

// MessageBlockchainIsCodeEmptyResponse represents a message
type MessageBlockchainIsCodeEmptyResponse struct {
	Message
}

// MessageBlockchainGetCodeRequest represents a message
type MessageBlockchainGetCodeRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetCodeRequest creates a message
func NewMessageBlockchainGetCodeRequest() *MessageBlockchainGetCodeRequest {
	message := &MessageBlockchainGetCodeRequest{}
	message.Kind = BlockchainGetCodeRequest
	return message
}

// MessageBlockchainGetCodeResponse represents a message
type MessageBlockchainGetCodeResponse struct {
	Message
	Code []byte
}

// NewMessageBlockchainGetCodeResponse creates a message
func NewMessageBlockchainGetCodeResponse(err error) *MessageBlockchainGetCodeResponse {
	message := &MessageBlockchainGetCodeResponse{}
	message.Kind = BlockchainGetCodeResponse
	message.SetError(err)
	return message
}

// MessageBlockchainGetBlockhashRequest represents a message
type MessageBlockchainGetBlockhashRequest struct {
	Message
}

// MessageBlockchainGetBlockhashResponse represents a message
type MessageBlockchainGetBlockhashResponse struct {
	Message
}

// MessageBlockchainLastNonceRequest represents a message
type MessageBlockchainLastNonceRequest struct {
	Message
}

// MessageBlockchainLastNonceResponse represents a message
type MessageBlockchainLastNonceResponse struct {
	Message
}

// MessageBlockchainLastRoundRequest represents a message
type MessageBlockchainLastRoundRequest struct {
	Message
}

// MessageBlockchainLastRoundResponse represents a message
type MessageBlockchainLastRoundResponse struct {
	Message
}

// MessageBlockchainLastTimeStampRequest represents a message
type MessageBlockchainLastTimeStampRequest struct {
	Message
}

// MessageBlockchainLastTimeStampResponse represents a message
type MessageBlockchainLastTimeStampResponse struct {
	Message
}

// MessageBlockchainLastRandomSeedRequest represents a message
type MessageBlockchainLastRandomSeedRequest struct {
	Message
}

// MessageBlockchainLastRandomSeedResponse represents a message
type MessageBlockchainLastRandomSeedResponse struct {
	Message
}

// MessageBlockchainLastEpochRequest represents a message
type MessageBlockchainLastEpochRequest struct {
	Message
}

// MessageBlockchainLastEpochResponse represents a message
type MessageBlockchainLastEpochResponse struct {
	Message
}

// MessageBlockchainGetStateRootHashRequest represents a message
type MessageBlockchainGetStateRootHashRequest struct {
	Message
}

// MessageBlockchainGetStateRootHashResponse represents a message
type MessageBlockchainGetStateRootHashResponse struct {
	Message
}

// MessageBlockchainCurrentNonceRequest represents a message
type MessageBlockchainCurrentNonceRequest struct {
	Message
}

// MessageBlockchainCurrentNonceResponse represents a message
type MessageBlockchainCurrentNonceResponse struct {
	Message
}

// MessageBlockchainCurrentRoundRequest represents a message
type MessageBlockchainCurrentRoundRequest struct {
	Message
}

// MessageBlockchainCurrentRoundResponse represents a message
type MessageBlockchainCurrentRoundResponse struct {
	Message
}

// MessageBlockchainCurrentTimeStampRequest represents a message
type MessageBlockchainCurrentTimeStampRequest struct {
	Message
}

// MessageBlockchainCurrentTimeStampResponse represents a message
type MessageBlockchainCurrentTimeStampResponse struct {
	Message
}

// MessageBlockchainCurrentRandomSeedRequest represents a message
type MessageBlockchainCurrentRandomSeedRequest struct {
	Message
}

// MessageBlockchainCurrentRandomSeedResponse represents a message
type MessageBlockchainCurrentRandomSeedResponse struct {
	Message
}

// MessageBlockchainCurrentEpochRequest represents a message
type MessageBlockchainCurrentEpochRequest struct {
	Message
}

// MessageBlockchainCurrentEpochResponse represents a message
type MessageBlockchainCurrentEpochResponse struct {
	Message
}
