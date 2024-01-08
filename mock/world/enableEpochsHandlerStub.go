package worldmock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	IsFlagDefinedCalled        func(flag core.EnableEpochFlag) bool
	IsFlagEnabledCalled        func(flag core.EnableEpochFlag) bool
	IsFlagEnabledInEpochCalled func(flag core.EnableEpochFlag, epoch uint32) bool
	GetActivationEpochCalled   func(flag core.EnableEpochFlag) uint32
}

// IsFlagDefined -
func (stub *EnableEpochsHandlerStub) IsFlagDefined(flag core.EnableEpochFlag) bool {
	if stub.IsFlagDefinedCalled != nil {
		return stub.IsFlagDefinedCalled(flag)
	}
	return true
}

// IsFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsFlagEnabled(flag core.EnableEpochFlag) bool {
	if stub.IsFlagEnabledCalled != nil {
		return stub.IsFlagEnabledCalled(flag)
	}
	return false
}

// IsFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsFlagEnabledInEpoch(flag core.EnableEpochFlag, epoch uint32) bool {
	if stub.IsFlagEnabledInEpochCalled != nil {
		return stub.IsFlagEnabledInEpochCalled(flag, epoch)
	}
	return false
}

// GetActivationEpoch -
func (stub *EnableEpochsHandlerStub) GetActivationEpoch(flag core.EnableEpochFlag) uint32 {
	if stub.GetActivationEpochCalled != nil {
		return stub.GetActivationEpochCalled(flag)
	}
	return 0
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}

// EnableEpochsHandlerStubAllFlags creates a new EnableEpochsHandlerStub with all flags enabled
func EnableEpochsHandlerStubAllFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{
		IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
			return flag == vmhost.StorageAPICostOptimizationFlag ||
				flag == vmhost.FixOOGReturnCodeFlag ||
				flag == vmhost.RemoveNonUpdatedStorageFlag ||
				flag == vmhost.CreateNFTThroughExecByCallerFlag ||
				flag == vmhost.ManagedCryptoAPIsFlag ||
				flag == vmhost.FailExecutionOnEveryAPIErrorFlag ||
				flag == vmhost.RefactorContextFlag ||
				flag == vmhost.DisableExecByCallerFlag ||
				flag == vmhost.CheckExecuteOnReadOnlyFlag ||
				flag == builtInFunctions.CheckCorrectTokenIDForTransferRoleFlag ||
				flag == builtInFunctions.ESDTTransferRoleFlag ||
				flag == builtInFunctions.GlobalMintBurnFlag ||
				flag == builtInFunctions.CheckFrozenCollectionFlag ||
				flag == builtInFunctions.FixAsyncCallbackCheckFlag ||
				flag == builtInFunctions.ESDTNFTImprovementV1Flag ||
				flag == builtInFunctions.SaveToSystemAccountFlag ||
				flag == builtInFunctions.ValueLengthCheckFlag ||
				flag == builtInFunctions.CheckTransferFlag ||
				flag == builtInFunctions.CheckFunctionArgumentFlag ||
				flag == builtInFunctions.FixOldTokenLiquidityFlag ||
				flag == builtInFunctions.AlwaysSaveTokenMetaDataFlag ||
				flag == builtInFunctions.SetGuardianFlag ||
				flag == builtInFunctions.WipeSingleNFTLiquidityDecreaseFlag ||
				flag == builtInFunctions.ConsistentTokensValuesLengthCheckFlag ||
				flag == builtInFunctions.ChangeUsernameFlag ||
				flag == builtInFunctions.AutoBalanceDataTriesFlag ||
				flag == builtInFunctions.ScToScLogEventFlag
		},
	}
}

// EnableEpochsHandlerStubNoFlags creates a new EnableEpochsHandlerStub with all flags disabled
func EnableEpochsHandlerStubNoFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{}
}
