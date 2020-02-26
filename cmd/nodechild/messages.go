package main

import (
	"fmt"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// Response is
type Response struct {
	ErrorMessage     string
	HasCriticalError bool
}

// HasError returns
func (response *Response) HasError() bool {
	return response.ErrorMessage != ""
}

// GetError returns
func (response *Response) GetError() error {
	return fmt.Errorf(response.ErrorMessage)
}

func (response *Response) String() string {
	return fmt.Sprintf("[%s][%t]", response.ErrorMessage, response.HasCriticalError)
}

// ContractRequest is
type ContractRequest struct {
	Tag         string
	CreateInput *vmcommon.ContractCreateInput
	CallInput   *vmcommon.ContractCallInput
}

func (request *ContractRequest) String() string {
	return fmt.Sprintf("ContractRequest [%s]", request.Tag)
}

// ContractResponse is
type ContractResponse struct {
	Tag      string
	VMOutput *vmcommon.VMOutput
	Response
}

func (response *ContractResponse) String() string {
	return fmt.Sprintf("ContractResponse [%s] [%v]", response.Tag, response.Response)
}

// HookCallRequest is
type HookCallRequest struct {
	Tag       string
	Hook      string
	Function  string
	Arguments []interface{}
}

// HookCallResponse is
type HookCallResponse struct {
	Tag    string
	Result []interface{}
	Response
}
