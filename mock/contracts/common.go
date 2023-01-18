package contracts

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/wasm-vm/arwen"
	"github.com/multiversx/wasm-vm/arwen/vmhooks"
)

// DirectCallGasTestConfig is configuration for direct call tests
type DirectCallGasTestConfig struct {
	GasUsedByParent      uint64
	GasUsedByChild       uint64
	GasProvided          uint64
	GasProvidedToChild   uint64
	ParentBalance        int64
	ChildBalance         int64
	ESDTTokensToTransfer uint64
}

// TransferAndExecuteTestConfig is configuration for transfer and execute tests
type TransferAndExecuteTestConfig struct {
	DirectCallGasTestConfig
	TransferFromParentToChild int64
	GasTransferToChild        uint64
}

// AsyncCallBaseTestConfig is base configuration for async call tests
type AsyncCallBaseTestConfig struct {
	GasProvided       uint64
	GasUsedByParent   uint64
	GasUsedByChild    uint64
	GasUsedByCallback uint64
	GasLockCost       uint64

	TransferFromParentToChild int64

	ParentBalance int64
	ChildBalance  int64
}

// AsyncCallTestConfig is configuration for async call tests
type AsyncCallTestConfig struct {
	AsyncCallBaseTestConfig
	TransferToThirdParty         int64
	TransferToVault              int64
	ESDTTokensToTransfer         uint64
	CallbackESDTTokensToTransfer uint64
}

// AsyncBuiltInCallTestConfig is configuration for async call tests of builtin functions
type AsyncBuiltInCallTestConfig struct {
	AsyncCallBaseTestConfig
	TransferFromChildToParent int64
}

// AsyncCallRecursiveTestConfig is configuration for recursive async call tests
type AsyncCallRecursiveTestConfig struct {
	AsyncCallBaseTestConfig
	RecursiveChildCalls int
}

// AsyncCallMultiChildTestConfig is configuration for recursivemultiple children async call tests
type AsyncCallMultiChildTestConfig struct {
	AsyncCallBaseTestConfig
	ChildCalls int
}

// GasTestConfig interface for gas tests configs
type GasTestConfig interface {
	GetGasUsedByChild() uint64
}

// GetGasUsedByChild - getter for GasUsedByChild
func (config AsyncCallTestConfig) GetGasUsedByChild() uint64 {
	return config.GasUsedByChild
}

// GetGasUsedByChild - getter for GasUsedByChild
func (config DirectCallGasTestConfig) GetGasUsedByChild() uint64 {
	return config.GasUsedByChild
}

// ExecuteOnSameContextInMockContracts - calls the corresponding method in elrond api
func ExecuteOnSameContextInMockContracts(host arwen.VMHost, input *vmcommon.ContractCallInput) int32 {
	return vmhooks.ExecuteOnSameContextWithTypedArgs(host, int64(input.GasProvided), input.CallValue, []byte(input.Function), input.RecipientAddr, input.Arguments)
}

// ExecuteOnDestContextInMockContracts - calls the corresponding method in elrond api
func ExecuteOnDestContextInMockContracts(host arwen.VMHost, input *vmcommon.ContractCallInput) int32 {
	return vmhooks.ExecuteOnDestContextWithTypedArgs(host, int64(input.GasProvided), input.CallValue, []byte(input.Function), input.RecipientAddr, input.Arguments)
}
