package main

import (
	"fmt"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ContractCommand is
type ContractCommand struct {
	Tag         string
	CreateInput *vmcommon.ContractCreateInput
	CallInput   *vmcommon.ContractCallInput
}

func (command *ContractCommand) String() string {
	return fmt.Sprintf("Command [%s]", command.Tag)
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
	Error  error
}
