package worldmock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	IsFlagDefinedCalled               func(flag core.EnableEpochFlag) bool
	IsFlagEnabledInCurrentEpochCalled func(flag core.EnableEpochFlag) bool
}

// IsFlagDefined -
func (stub *EnableEpochsHandlerStub) IsFlagDefined(flag core.EnableEpochFlag) bool {
	if stub.IsFlagDefinedCalled != nil {
		return stub.IsFlagDefinedCalled(flag)
	}
	return true
}

// IsFlagEnabledInCurrentEpoch -
func (stub *EnableEpochsHandlerStub) IsFlagEnabledInCurrentEpoch(flag core.EnableEpochFlag) bool {
	if stub.IsFlagEnabledInCurrentEpochCalled != nil {
		return stub.IsFlagEnabledInCurrentEpochCalled(flag)
	}
	return false
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}

// EnableEpochsHandlerStubAllFlags creates a new EnableEpochsHandlerStub with all flags enabled
func EnableEpochsHandlerStubAllFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{
		IsFlagEnabledInCurrentEpochCalled: func(flag core.EnableEpochFlag) bool {
			return flag == core.StorageAPICostOptimizationFlag ||
				flag == core.FixOOGReturnCodeFlag ||
				flag == core.RemoveNonUpdatedStorageFlag ||
				flag == core.CreateNFTThroughExecByCallerFlag ||
				flag == core.ManagedCryptoAPIsFlag ||
				flag == core.FailExecutionOnEveryAPIErrorFlag ||
				flag == core.RefactorContextFlag ||
				flag == core.DisableExecByCallerFlag ||
				flag == core.CheckExecuteOnReadOnlyFlag
		},
	}
}

// EnableEpochsHandlerStubNoFlags creates a new EnableEpochsHandlerStub with all flags disabled
func EnableEpochsHandlerStubNoFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{}
}
