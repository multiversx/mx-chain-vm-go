package common

import (
	"fmt"
)

// MessageKind is
type MessageKind uint32

const (
	FirstKind MessageKind = iota
	Stop
	ContractDeployRequest
	ContractCallRequest
	ContractResponse
	BlockchainAccountExistsRequest
	BlockchainAccountExistsResponse
	BlockchainNewAddressRequest
	BlockchainNewAddressResponse
	BlockchainGetBalanceRequest
	BlockchainGetBalanceResponse
	BlockchainGetNonceRequest
	BlockchainGetNonceResponse
	BlockchainGetStorageDataRequest
	BlockchainGetStorageDataResponse
	BlockchainIsCodeEmptyRequest
	BlockchainIsCodeEmptyResponse
	BlockchainGetCodeRequest
	BlockchainGetCodeResponse
	BlockchainGetBlockhashRequest
	BlockchainGetBlockhashResponse
	BlockchainLastNonceRequest
	BlockchainLastNonceResponse
	BlockchainLastRoundRequest
	BlockchainLastRoundResponse
	BlockchainLastTimeStampRequest
	BlockchainLastTimeStampResponse
	BlockchainLastRandomSeedRequest
	BlockchainLastRandomSeedResponse
	BlockchainLastEpochRequest
	BlockchainLastEpochResponse
	BlockchainGetStateRootHashRequest
	BlockchainGetStateRootHashResponse
	BlockchainCurrentNonceRequest
	BlockchainCurrentNonceResponse
	BlockchainCurrentRoundRequest
	BlockchainCurrentRoundResponse
	BlockchainCurrentTimeStampRequest
	BlockchainCurrentTimeStampResponse
	BlockchainCurrentRandomSeedRequest
	BlockchainCurrentRandomSeedResponse
	BlockchainCurrentEpochRequest
	BlockchainCurrentEpochResponse
	LastKind
)

// MessageHandler is
type MessageHandler interface {
	GetNonce() uint32
	SetNonce(nonce uint32)
	GetKind() MessageKind
	SetKind(kind MessageKind)
	GetError() error
	SetError(err error)
}

// Message is
type Message struct {
	DialogueNonce uint32
	Kind          MessageKind
	ErrorMessage  string
}

// GetNonce gets
func (message *Message) GetNonce() uint32 {
	return message.DialogueNonce
}

// SetNonce sets
func (message *Message) SetNonce(nonce uint32) {
	message.DialogueNonce = nonce
}

// GetKind gets
func (message *Message) GetKind() MessageKind {
	return message.Kind
}

// SetKind sets
func (message *Message) SetKind(kind MessageKind) {
	message.Kind = kind
}

// GetError gets
func (message *Message) GetError() error {
	if message.ErrorMessage == "" {
		return nil
	}

	return fmt.Errorf(message.ErrorMessage)
}

// SetError sets
func (message *Message) SetError(err error) {
	if err != nil {
		message.ErrorMessage = err.Error()
	}
}

func (message *Message) String() string {
	return fmt.Sprintf("[kind=%d nonce=%d err=%s]", message.Kind, message.DialogueNonce, message.ErrorMessage)
}

// MessageStop is
type MessageStop struct {
	Message
}

// NewMessageStop creates a message
func NewMessageStop() *MessageStop {
	message := &MessageStop{}
	message.Kind = Stop
	return message
}

// MessageCallback is a callback
type MessageCallback func(MessageHandler) (MessageHandler, error)

func noopHandler(message MessageHandler) (MessageHandler, error) {
	panic("Noop handler called")
}

// CreateHandlerSlots creates
func CreateHandlerSlots() []MessageCallback {
	slots := make([]MessageCallback, LastKind)
	for i := 0; i < len(slots); i++ {
		slots[i] = noopHandler
	}

	return slots
}

// IsHookCallRequest returns
func IsHookCallRequest(message MessageHandler) bool {
	kind := message.GetKind()
	return kind >= BlockchainAccountExistsRequest && kind <= BlockchainCurrentEpochResponse
}

// IsContractResponse returns
func IsContractResponse(message MessageHandler) bool {
	return message.GetKind() == ContractResponse
}
