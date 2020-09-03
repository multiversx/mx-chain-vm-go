package common

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// MessageContractDeployRequest is a deploy request message (from the Node)
type MessageContractDeployRequest struct {
	Message
	CreateInput *vmcommon.ContractCreateInput
}

// NewMessageContractDeployRequest creates a MessageContractDeployRequest
func NewMessageContractDeployRequest(input *vmcommon.ContractCreateInput) *MessageContractDeployRequest {
	message := &MessageContractDeployRequest{}
	message.Kind = ContractDeployRequest
	message.CreateInput = input
	return message
}

// MessageContractCallRequest is a call request message (from the Node)
type MessageContractCallRequest struct {
	Message
	CallInput *vmcommon.ContractCallInput
}

// NewMessageContractCallRequest creates a MessageContractCallRequest
func NewMessageContractCallRequest(input *vmcommon.ContractCallInput) *MessageContractCallRequest {
	message := &MessageContractCallRequest{}
	message.Kind = ContractCallRequest
	message.CallInput = input
	return message
}

// MessageContractResponse is a contract response message (from Arwen)
type MessageContractResponse struct {
	Message
	SerializableVMOutput *SerializableVMOutput
}

// NewMessageContractResponse creates a MessageContractResponse
func NewMessageContractResponse(vmOutput *vmcommon.VMOutput, err error) *MessageContractResponse {
	message := &MessageContractResponse{}
	message.Kind = ContractResponse
	message.SerializableVMOutput = NewSerializableVMOutput(vmOutput)
	message.SetError(err)
	return message
}

// MessageVersionRequest is a version request message (from the Node)
type MessageVersionRequest struct {
	Message
}

// NewMessageVersionRequest creates a MessageVersionRequest
func NewMessageVersionRequest() *MessageVersionRequest {
	message := &MessageVersionRequest{}
	message.Kind = VersionRequest
	return message
}

// MessageVersionResponse is a version response message (from Arwen)
type MessageVersionResponse struct {
	Message
	Version string
}

// NewMessageVersionResponse creates a MessageVersionResponse
func NewMessageVersionResponse(version string) *MessageVersionResponse {
	message := &MessageVersionResponse{}
	message.Kind = VersionResponse
	message.Version = version
	return message
}
