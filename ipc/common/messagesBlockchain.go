package common

import (
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// MessageBlockchainAccountExistsRequest represents a request message
type MessageBlockchainAccountExistsRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainAccountExistsRequest creates a request message
func NewMessageBlockchainAccountExistsRequest(address []byte) *MessageBlockchainAccountExistsRequest {
	message := &MessageBlockchainAccountExistsRequest{}
	message.Kind = BlockchainAccountExistsRequest
	message.Address = address
	return message
}

// MessageBlockchainAccountExistsResponse represents a response message
type MessageBlockchainAccountExistsResponse struct {
	Message
	Result bool
}

// NewMessageBlockchainAccountExistsResponse creates a response message
func NewMessageBlockchainAccountExistsResponse(result bool, err error) *MessageBlockchainAccountExistsResponse {
	message := &MessageBlockchainAccountExistsResponse{}
	message.Kind = BlockchainAccountExistsResponse
	message.Result = result
	message.SetError(err)
	return message
}

// MessageBlockchainNewAddressRequest represents a request message
type MessageBlockchainNewAddressRequest struct {
	Message
	CreatorAddress []byte
	CreatorNonce   uint64
	VmType         []byte
}

// NewMessageBlockchainNewAddressRequest creates a request message
func NewMessageBlockchainNewAddressRequest(creatorAddress []byte, creatorNonce uint64, vmType []byte) *MessageBlockchainNewAddressRequest {
	message := &MessageBlockchainNewAddressRequest{}
	message.Kind = BlockchainNewAddressRequest
	message.CreatorAddress = creatorAddress
	message.CreatorNonce = creatorNonce
	message.VmType = vmType
	return message
}

// MessageBlockchainNewAddressResponse represents a response message
type MessageBlockchainNewAddressResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainNewAddressResponse creates a response message
func NewMessageBlockchainNewAddressResponse(result []byte, err error) *MessageBlockchainNewAddressResponse {
	message := &MessageBlockchainNewAddressResponse{}
	message.Kind = BlockchainNewAddressResponse
	message.Result = result
	message.SetError(err)
	return message
}

// MessageBlockchainGetBalanceRequest represents a request message
type MessageBlockchainGetBalanceRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetBalanceRequest creates a request message
func NewMessageBlockchainGetBalanceRequest(address []byte) *MessageBlockchainGetBalanceRequest {
	message := &MessageBlockchainGetBalanceRequest{}
	message.Kind = BlockchainGetBalanceRequest
	message.Address = address
	return message
}

// MessageBlockchainGetBalanceResponse represents a response message
type MessageBlockchainGetBalanceResponse struct {
	Message
	Balance *big.Int
}

// NewMessageBlockchainGetBalanceResponse creates a response message
func NewMessageBlockchainGetBalanceResponse(balance *big.Int, err error) *MessageBlockchainGetBalanceResponse {
	message := &MessageBlockchainGetBalanceResponse{}
	message.Kind = BlockchainGetBalanceResponse
	message.Balance = balance
	message.SetError(err)
	return message
}

// MessageBlockchainGetNonceRequest represents a request message
type MessageBlockchainGetNonceRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetNonceRequest creates a request message
func NewMessageBlockchainGetNonceRequest(address []byte) *MessageBlockchainGetNonceRequest {
	message := &MessageBlockchainGetNonceRequest{}
	message.Kind = BlockchainGetNonceRequest
	message.Address = address
	return message
}

// MessageBlockchainGetNonceResponse represents a response message
type MessageBlockchainGetNonceResponse struct {
	Message
	Nonce uint64
}

// NewMessageBlockchainGetNonceResponse creates a response message
func NewMessageBlockchainGetNonceResponse(nonce uint64, err error) *MessageBlockchainGetNonceResponse {
	message := &MessageBlockchainGetNonceResponse{}
	message.Kind = BlockchainGetNonceResponse
	message.Nonce = nonce
	message.SetError(err)
	return message
}

// MessageBlockchainGetStorageDataRequest represents a request message
type MessageBlockchainGetStorageDataRequest struct {
	Message
	Address []byte
	Index   []byte
}

// NewMessageBlockchainGetStorageDataRequest creates a request message
func NewMessageBlockchainGetStorageDataRequest(address []byte, index []byte) *MessageBlockchainGetStorageDataRequest {
	message := &MessageBlockchainGetStorageDataRequest{}
	message.Kind = BlockchainGetStorageDataRequest
	message.Address = address
	message.Index = index
	return message
}

// MessageBlockchainGetStorageDataResponse represents a response message
type MessageBlockchainGetStorageDataResponse struct {
	Message
	Data []byte
}

// NewMessageBlockchainGetStorageDataResponse creates a response message
func NewMessageBlockchainGetStorageDataResponse(data []byte, err error) *MessageBlockchainGetStorageDataResponse {
	message := &MessageBlockchainGetStorageDataResponse{}
	message.Kind = BlockchainGetStorageDataResponse
	message.Data = data
	message.SetError(err)
	return message
}

// MessageBlockchainIsCodeEmptyRequest represents a request message
type MessageBlockchainIsCodeEmptyRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainIsCodeEmptyRequest creates a request message
func NewMessageBlockchainIsCodeEmptyRequest(address []byte) *MessageBlockchainIsCodeEmptyRequest {
	message := &MessageBlockchainIsCodeEmptyRequest{}
	message.Kind = BlockchainIsCodeEmptyRequest
	message.Address = address
	return message
}

// MessageBlockchainIsCodeEmptyResponse represents a response message
type MessageBlockchainIsCodeEmptyResponse struct {
	Message
	Result bool
}

// NewMessageBlockchainIsCodeEmptyResponse creates a response message
func NewMessageBlockchainIsCodeEmptyResponse(result bool, err error) *MessageBlockchainIsCodeEmptyResponse {
	message := &MessageBlockchainIsCodeEmptyResponse{}
	message.Kind = BlockchainIsCodeEmptyResponse
	message.Result = result
	message.SetError(err)
	return message
}

// MessageBlockchainGetCodeRequest represents a request message
type MessageBlockchainGetCodeRequest struct {
	Message
	Address []byte
}

// NewMessageBlockchainGetCodeRequest creates a request message
func NewMessageBlockchainGetCodeRequest(address []byte) *MessageBlockchainGetCodeRequest {
	message := &MessageBlockchainGetCodeRequest{}
	message.Kind = BlockchainGetCodeRequest
	message.Address = address
	return message
}

// MessageBlockchainGetCodeResponse represents a response message
type MessageBlockchainGetCodeResponse struct {
	Message
	Code []byte
}

// NewMessageBlockchainGetCodeResponse creates a response message
func NewMessageBlockchainGetCodeResponse(code []byte, err error) *MessageBlockchainGetCodeResponse {
	message := &MessageBlockchainGetCodeResponse{}
	message.Kind = BlockchainGetCodeResponse
	message.Code = code
	message.SetError(err)
	return message
}

// MessageBlockchainGetBlockhashRequest represents a request message
type MessageBlockchainGetBlockhashRequest struct {
	Message
	Nonce uint64
}

// NewMessageBlockchainGetBlockhashRequest creates a request message
func NewMessageBlockchainGetBlockhashRequest(nonce uint64) *MessageBlockchainGetBlockhashRequest {
	message := &MessageBlockchainGetBlockhashRequest{}
	message.Kind = BlockchainGetBlockhashRequest
	message.Nonce = nonce
	return message
}

// MessageBlockchainGetBlockhashResponse represents a response message
type MessageBlockchainGetBlockhashResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainGetBlockhashResponse creates a response message
func NewMessageBlockchainGetBlockhashResponse(result []byte, err error) *MessageBlockchainGetBlockhashResponse {
	message := &MessageBlockchainGetBlockhashResponse{}
	message.Kind = BlockchainGetBlockhashResponse
	message.Result = result
	message.SetError(err)
	return message
}

// MessageBlockchainLastNonceRequest represents a request message
type MessageBlockchainLastNonceRequest struct {
	Message
}

// NewMessageBlockchainLastNonceRequest creates a request message
func NewMessageBlockchainLastNonceRequest() *MessageBlockchainLastNonceRequest {
	message := &MessageBlockchainLastNonceRequest{}
	message.Kind = BlockchainLastNonceRequest

	return message
}

// MessageBlockchainLastNonceResponse represents a response message
type MessageBlockchainLastNonceResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainLastNonceResponse creates a response message
func NewMessageBlockchainLastNonceResponse(result uint64) *MessageBlockchainLastNonceResponse {
	message := &MessageBlockchainLastNonceResponse{}
	message.Kind = BlockchainLastNonceResponse
	message.Result = result
	return message
}

// MessageBlockchainLastRoundRequest represents a request message
type MessageBlockchainLastRoundRequest struct {
	Message
}

// NewMessageBlockchainLastRoundRequest creates a request message
func NewMessageBlockchainLastRoundRequest() *MessageBlockchainLastRoundRequest {
	message := &MessageBlockchainLastRoundRequest{}
	message.Kind = BlockchainLastRoundRequest

	return message
}

// MessageBlockchainLastRoundResponse represents a response message
type MessageBlockchainLastRoundResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainLastRoundResponse creates a response message
func NewMessageBlockchainLastRoundResponse(result uint64) *MessageBlockchainLastRoundResponse {
	message := &MessageBlockchainLastRoundResponse{}
	message.Kind = BlockchainLastRoundResponse
	message.Result = result
	return message
}

// MessageBlockchainLastTimeStampRequest represents a request message
type MessageBlockchainLastTimeStampRequest struct {
	Message
}

// NewMessageBlockchainLastTimeStampRequest creates a request message
func NewMessageBlockchainLastTimeStampRequest() *MessageBlockchainLastTimeStampRequest {
	message := &MessageBlockchainLastTimeStampRequest{}
	message.Kind = BlockchainLastTimeStampRequest

	return message
}

// MessageBlockchainLastTimeStampResponse represents a response message
type MessageBlockchainLastTimeStampResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainLastTimeStampResponse creates a response message
func NewMessageBlockchainLastTimeStampResponse(result uint64) *MessageBlockchainLastTimeStampResponse {
	message := &MessageBlockchainLastTimeStampResponse{}
	message.Kind = BlockchainLastTimeStampResponse
	message.Result = result
	return message
}

// MessageBlockchainLastRandomSeedRequest represents a request message
type MessageBlockchainLastRandomSeedRequest struct {
	Message
}

// NewMessageBlockchainLastRandomSeedRequest creates a request message
func NewMessageBlockchainLastRandomSeedRequest() *MessageBlockchainLastRandomSeedRequest {
	message := &MessageBlockchainLastRandomSeedRequest{}
	message.Kind = BlockchainLastRandomSeedRequest

	return message
}

// MessageBlockchainLastRandomSeedResponse represents a response message
type MessageBlockchainLastRandomSeedResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainLastRandomSeedResponse creates a response message
func NewMessageBlockchainLastRandomSeedResponse(result []byte) *MessageBlockchainLastRandomSeedResponse {
	message := &MessageBlockchainLastRandomSeedResponse{}
	message.Kind = BlockchainLastRandomSeedResponse
	message.Result = result
	return message
}

// MessageBlockchainLastEpochRequest represents a request message
type MessageBlockchainLastEpochRequest struct {
	Message
}

// NewMessageBlockchainLastEpochRequest creates a request message
func NewMessageBlockchainLastEpochRequest() *MessageBlockchainLastEpochRequest {
	message := &MessageBlockchainLastEpochRequest{}
	message.Kind = BlockchainLastEpochRequest

	return message
}

// MessageBlockchainLastEpochResponse represents a response message
type MessageBlockchainLastEpochResponse struct {
	Message
	Result uint32
}

// NewMessageBlockchainLastEpochResponse creates a response message
func NewMessageBlockchainLastEpochResponse(result uint32) *MessageBlockchainLastEpochResponse {
	message := &MessageBlockchainLastEpochResponse{}
	message.Kind = BlockchainLastEpochResponse
	message.Result = result
	return message
}

// MessageBlockchainGetStateRootHashRequest represents a request message
type MessageBlockchainGetStateRootHashRequest struct {
	Message
}

// NewMessageBlockchainGetStateRootHashRequest creates a request message
func NewMessageBlockchainGetStateRootHashRequest() *MessageBlockchainGetStateRootHashRequest {
	message := &MessageBlockchainGetStateRootHashRequest{}
	message.Kind = BlockchainGetStateRootHashRequest

	return message
}

// MessageBlockchainGetStateRootHashResponse represents a response message
type MessageBlockchainGetStateRootHashResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainGetStateRootHashResponse creates a response message
func NewMessageBlockchainGetStateRootHashResponse(result []byte) *MessageBlockchainGetStateRootHashResponse {
	message := &MessageBlockchainGetStateRootHashResponse{}
	message.Kind = BlockchainGetStateRootHashResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentNonceRequest represents a request message
type MessageBlockchainCurrentNonceRequest struct {
	Message
}

// NewMessageBlockchainCurrentNonceRequest creates a request message
func NewMessageBlockchainCurrentNonceRequest() *MessageBlockchainCurrentNonceRequest {
	message := &MessageBlockchainCurrentNonceRequest{}
	message.Kind = BlockchainCurrentNonceRequest

	return message
}

// MessageBlockchainCurrentNonceResponse represents a response message
type MessageBlockchainCurrentNonceResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainCurrentNonceResponse creates a response message
func NewMessageBlockchainCurrentNonceResponse(result uint64) *MessageBlockchainCurrentNonceResponse {
	message := &MessageBlockchainCurrentNonceResponse{}
	message.Kind = BlockchainCurrentNonceResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentRoundRequest represents a request message
type MessageBlockchainCurrentRoundRequest struct {
	Message
}

// NewMessageBlockchainCurrentRoundRequest creates a request message
func NewMessageBlockchainCurrentRoundRequest() *MessageBlockchainCurrentRoundRequest {
	message := &MessageBlockchainCurrentRoundRequest{}
	message.Kind = BlockchainCurrentRoundRequest

	return message
}

// MessageBlockchainCurrentRoundResponse represents a response message
type MessageBlockchainCurrentRoundResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainCurrentRoundResponse creates a response message
func NewMessageBlockchainCurrentRoundResponse(result uint64) *MessageBlockchainCurrentRoundResponse {
	message := &MessageBlockchainCurrentRoundResponse{}
	message.Kind = BlockchainCurrentRoundResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentTimeStampRequest represents a request message
type MessageBlockchainCurrentTimeStampRequest struct {
	Message
}

// NewMessageBlockchainCurrentTimeStampRequest creates a request message
func NewMessageBlockchainCurrentTimeStampRequest() *MessageBlockchainCurrentTimeStampRequest {
	message := &MessageBlockchainCurrentTimeStampRequest{}
	message.Kind = BlockchainCurrentTimeStampRequest

	return message
}

// MessageBlockchainCurrentTimeStampResponse represents a response message
type MessageBlockchainCurrentTimeStampResponse struct {
	Message
	Result uint64
}

// NewMessageBlockchainCurrentTimeStampResponse creates a response message
func NewMessageBlockchainCurrentTimeStampResponse(result uint64) *MessageBlockchainCurrentTimeStampResponse {
	message := &MessageBlockchainCurrentTimeStampResponse{}
	message.Kind = BlockchainCurrentTimeStampResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentRandomSeedRequest represents a request message
type MessageBlockchainCurrentRandomSeedRequest struct {
	Message
}

// NewMessageBlockchainCurrentRandomSeedRequest creates a request message
func NewMessageBlockchainCurrentRandomSeedRequest() *MessageBlockchainCurrentRandomSeedRequest {
	message := &MessageBlockchainCurrentRandomSeedRequest{}
	message.Kind = BlockchainCurrentRandomSeedRequest

	return message
}

// MessageBlockchainCurrentRandomSeedResponse represents a response message
type MessageBlockchainCurrentRandomSeedResponse struct {
	Message
	Result []byte
}

// NewMessageBlockchainCurrentRandomSeedResponse creates a response message
func NewMessageBlockchainCurrentRandomSeedResponse(result []byte) *MessageBlockchainCurrentRandomSeedResponse {
	message := &MessageBlockchainCurrentRandomSeedResponse{}
	message.Kind = BlockchainCurrentRandomSeedResponse
	message.Result = result
	return message
}

// MessageBlockchainCurrentEpochRequest represents a request message
type MessageBlockchainCurrentEpochRequest struct {
	Message
}

// NewMessageBlockchainCurrentEpochRequest creates a request message
func NewMessageBlockchainCurrentEpochRequest() *MessageBlockchainCurrentEpochRequest {
	message := &MessageBlockchainCurrentEpochRequest{}
	message.Kind = BlockchainCurrentEpochRequest

	return message
}

// MessageBlockchainCurrentEpochResponse represents a response message
type MessageBlockchainCurrentEpochResponse struct {
	Message
	Result uint32
}

// NewMessageBlockchainCurrentEpochResponse creates a response message
func NewMessageBlockchainCurrentEpochResponse(result uint32) *MessageBlockchainCurrentEpochResponse {
	message := &MessageBlockchainCurrentEpochResponse{}
	message.Kind = BlockchainCurrentEpochResponse
	message.Result = result
	return message
}

// MessageBlockchainProcessBuiltinFunctionRequest represents a request message
type MessageBlockchainProcessBuiltinFunctionRequest struct {
	Message
	CallInput vmcommon.ContractCallInput
}

// NewMessageBlockchainProcessBuiltinFunctionRequest creates a request message
func NewMessageBlockchainProcessBuiltinFunctionRequest(callInput vmcommon.ContractCallInput) *MessageBlockchainProcessBuiltinFunctionRequest {
	message := &MessageBlockchainProcessBuiltinFunctionRequest{}
	message.Kind = BlockchainProcessBuiltinFunctionRequest
	message.CallInput = callInput

	return message
}

// MessageBlockchainProcessBuiltinFunctionResponse represents a response message
type MessageBlockchainProcessBuiltinFunctionResponse struct {
	Message
	Value       *big.Int
	GasConsumed uint64
}

// NewMessageBlockchainProcessBuiltinFunctionResponse creates a response message
func NewMessageBlockchainProcessBuiltinFunctionResponse(value *big.Int, gasConsumed uint64, err error) *MessageBlockchainProcessBuiltinFunctionResponse {
	message := &MessageBlockchainProcessBuiltinFunctionResponse{}
	message.Kind = BlockchainProcessBuiltinFunctionResponse
	message.Value = value
	message.GasConsumed = gasConsumed
	message.SetError(err)
	return message
}

// MessageBlockchainGetBuiltinFunctionNamesRequest represents a request message
type MessageBlockchainGetBuiltinFunctionNamesRequest struct {
	Message
}

// NewMessageBlockchainGetBuiltinFunctionNamesRequest creates a request message
func NewMessageBlockchainGetBuiltinFunctionNamesRequest() *MessageBlockchainGetBuiltinFunctionNamesRequest {
	message := &MessageBlockchainGetBuiltinFunctionNamesRequest{}
	message.Kind = BlockchainGetBuiltinFunctionNamesRequest

	return message
}

// MessageBlockchainGetBuiltinFunctionNamesResponse represents a response message
type MessageBlockchainGetBuiltinFunctionNamesResponse struct {
	Message
	FunctionNames vmcommon.FunctionNames
}

// NewMessageBlockchainGetBuiltinFunctionNamesResponse creates a response message
func NewMessageBlockchainGetBuiltinFunctionNamesResponse(functionNames vmcommon.FunctionNames) *MessageBlockchainGetBuiltinFunctionNamesResponse {
	message := &MessageBlockchainGetBuiltinFunctionNamesResponse{}
	message.Kind = BlockchainGetBuiltinFunctionNamesResponse
	message.FunctionNames = functionNames

	return message
}
