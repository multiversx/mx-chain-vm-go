package common

import (
	"fmt"
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// Message is
type Message interface {
	GetNonce() uint32
	SetNonce(nonce uint32)
}

// ContractRequest is
type ContractRequest struct {
	Nonce       uint32
	Action      string
	CreateInput *vmcommon.ContractCreateInput
	CallInput   *vmcommon.ContractCallInput
}

// GetNonce gets
func (request *ContractRequest) GetNonce() uint32 {
	return request.Nonce
}

// SetNonce sets
func (request *ContractRequest) SetNonce(nonce uint32) {
	request.Nonce = nonce
}

func (request *ContractRequest) String() string {
	return fmt.Sprintf("ContractRequest [%s]", request.Action)
}

// HookCallRequestOrContractResponse is
type HookCallRequestOrContractResponse struct {
	Type             string
	Nonce            uint32
	Hook             string
	Function         string
	Bytes1           []byte
	Bytes2           []byte
	Uint64_1         uint64
	VMOutput         *vmcommon.VMOutput
	ErrorMessage     string
	HasCriticalError bool
}

// NewHookCallRequest creates
func NewHookCallRequest(hook string, function string) *HookCallRequestOrContractResponse {
	return &HookCallRequestOrContractResponse{
		Type:     "HookCallRequest",
		Hook:     hook,
		Function: function,
	}
}

// NewContractResponse creates
func NewContractResponse(vmOutput *vmcommon.VMOutput, err error) *HookCallRequestOrContractResponse {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}

	return &HookCallRequestOrContractResponse{
		Type:         "ContractResponse",
		VMOutput:     vmOutput,
		ErrorMessage: errorMessage,
	}
}

// NewCriticalError creates
func NewCriticalError(errorMessage string) *HookCallRequestOrContractResponse {
	return &HookCallRequestOrContractResponse{
		ErrorMessage:     errorMessage,
		HasCriticalError: true,
	}
}

// IsHookCallRequest gets
func (message *HookCallRequestOrContractResponse) IsHookCallRequest() bool {
	return message.Type == "HookCallRequest"
}

// IsContractResponse gets
func (message *HookCallRequestOrContractResponse) IsContractResponse() bool {
	return message.Type == "ContractResponse"
}

// IsCriticalError returns
func (message *HookCallRequestOrContractResponse) IsCriticalError() bool {
	return message.HasCriticalError
}

// HasError returns
func (message *HookCallRequestOrContractResponse) HasError() bool {
	return message.ErrorMessage != ""
}

// GetError returns
func (message *HookCallRequestOrContractResponse) GetError() error {
	if message.ErrorMessage == "" {
		return nil
	}

	return fmt.Errorf(message.ErrorMessage)
}

// GetNonce gets
func (message *HookCallRequestOrContractResponse) GetNonce() uint32 {
	return message.Nonce
}

// SetNonce sets
func (message *HookCallRequestOrContractResponse) SetNonce(nonce uint32) {
	message.Nonce = nonce
}

func (message *HookCallRequestOrContractResponse) String() string {
	return fmt.Sprintf("[%s][%s]", message.Type, message.ErrorMessage)
}

// HookCallResponse is
type HookCallResponse struct {
	Nonce        uint32
	ErrorMessage string
	Bool1        bool
	Bytes1       []byte
	Bytes2       []byte
	BigInt1      *big.Int
	Uint64_1     uint64
	Uint32_1     uint32
}

// HasError returns
func (response *HookCallResponse) HasError() bool {
	return response.ErrorMessage != ""
}

// GetError returns
func (response *HookCallResponse) GetError() error {
	return fmt.Errorf(response.ErrorMessage)
}

// GetNonce gets
func (response *HookCallResponse) GetNonce() uint32 {
	return response.Nonce
}

// SetNonce sets
func (response *HookCallResponse) SetNonce(nonce uint32) {
	response.Nonce = nonce
}

func (response *HookCallResponse) String() string {
	return fmt.Sprintf("[%s]", response.ErrorMessage)
}
