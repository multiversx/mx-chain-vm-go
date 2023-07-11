package worldmock

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	GetCurrentEpochField                           uint32
	IsFixOOGReturnCodeFlagEnabledField             bool
	IsRemoveNonUpdatedStorageFlagEnabledField      bool
	IsCreateNFTThroughExecByCallerFlagEnabledField bool
	IsStorageAPICostOptimizationFlagEnabledField   bool
	IsFailExecutionOnEveryAPIErrorFlagEnabledField bool
	IsManagedCryptoAPIsFlagEnabledField            bool
	IsDisableExecByCallerFlagEnabledField          bool
	IsRefactorContextFlagEnabledField              bool
	IsCheckExecuteOnReadOnlyFlagEnabledField       bool
}

// GetCurrentEpoch -
func (stub *EnableEpochsHandlerStub) GetCurrentEpoch() uint32 {
	return stub.GetCurrentEpochField
}

// IsStorageAPICostOptimizationFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsStorageAPICostOptimizationFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsStorageAPICostOptimizationFlagEnabledField
}

// IsManagedCryptoAPIsFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsManagedCryptoAPIsFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsManagedCryptoAPIsFlagEnabledField
}

// IsRemoveNonUpdatedStorageFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsRemoveNonUpdatedStorageFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsRemoveNonUpdatedStorageFlagEnabledField
}

// IsRefactorContextFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsRefactorContextFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsRefactorContextFlagEnabledField
}

// IsFailExecutionOnEveryAPIErrorFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsFailExecutionOnEveryAPIErrorFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsFailExecutionOnEveryAPIErrorFlagEnabledField
}

// IsFixOOGReturnCodeFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsFixOOGReturnCodeFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsFixOOGReturnCodeFlagEnabledField
}

// IsCreateNFTThroughExecByCallerFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsCreateNFTThroughExecByCallerFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsCreateNFTThroughExecByCallerFlagEnabledField
}

// IsDisableExecByCallerFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsDisableExecByCallerFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsDisableExecByCallerFlagEnabledField
}

// IsCheckExecuteOnReadOnlyFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsCheckExecuteOnReadOnlyFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsCheckExecuteOnReadOnlyFlagEnabledField
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}

// EnableEpochsHandlerStubAllFlags creates a new EnableEpochsHandlerStub with all flags enabled
func EnableEpochsHandlerStubAllFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField:   true,
		IsFixOOGReturnCodeFlagEnabledField:             true,
		IsRemoveNonUpdatedStorageFlagEnabledField:      true,
		IsCreateNFTThroughExecByCallerFlagEnabledField: true,
		IsManagedCryptoAPIsFlagEnabledField:            true,
		IsFailExecutionOnEveryAPIErrorFlagEnabledField: true,
		IsRefactorContextFlagEnabledField:              true,
		IsDisableExecByCallerFlagEnabledField:          true,
		IsCheckExecuteOnReadOnlyFlagEnabledField:       true,
	}
}

// EnableEpochsHandlerStubNoFlags creates a new EnableEpochsHandlerStub with all flags disabled
func EnableEpochsHandlerStubNoFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{}
}
