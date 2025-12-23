package contexts

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

type vmInputFactory struct{}

// NewVMInputFactory creates new vm input factory
func NewVMInputFactory() *vmInputFactory {
	return &vmInputFactory{}
}

// CreateContractCallInput create contract call input
func (f *vmInputFactory) CreateContractCallInput(
	internalVMInput vmcommon.VMInput,
	vmInput vmcommon.ContractCallInputHandler,
) vmcommon.ContractCallInputHandler {
	return &vmcommon.ContractCallInput{
		VMInput:       internalVMInput,
		RecipientAddr: vmInput.GetRecipientAddr(),
		Function:      vmInput.GetFunction(),
	}
}

// IsInterfaceNil returns true if underlying object is nil
func (f *vmInputFactory) IsInterfaceNil() bool {
	return f == nil
}
