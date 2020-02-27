package common

import (
	"fmt"
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ContractRequest is
type ContractRequest struct {
	Tag         string
	CreateInput *vmcommon.ContractCreateInput
	CallInput   *vmcommon.ContractCallInput
}

func (request *ContractRequest) String() string {
	return fmt.Sprintf("ContractRequest [%s]", request.Tag)
}

// HookCallRequestOrContractResponse is
type HookCallRequestOrContractResponse struct {
	Type             string
	Tag              string
	Hook             string
	Function         string
	Arguments        []interface{}
	VMOutput         *vmcommon.VMOutput
	ErrorMessage     string
	HasCriticalError bool
}

// NewHookCallRequest creates
func NewHookCallRequest(hook string, function string, arguments ...interface{}) *HookCallRequestOrContractResponse {
	return &HookCallRequestOrContractResponse{
		Type:      "HookCallRequest",
		Hook:      hook,
		Function:  function,
		Arguments: arguments,
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
	return fmt.Errorf(message.ErrorMessage)
}

func (message *HookCallRequestOrContractResponse) String() string {
	return fmt.Sprintf("[%s][%s]", message.Type, message.ErrorMessage)
}

// HookCallResponse is
type HookCallResponse struct {
	Tag          string
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

func (response *HookCallResponse) String() string {
	return fmt.Sprintf("[%s][%s]", response.Bytes1, response.ErrorMessage)
}
