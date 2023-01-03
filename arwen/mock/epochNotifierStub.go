package mock

import (
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// EpochNotifierStub -
type EpochNotifierStub struct {
	CurrentEpochCalled          func() uint32
	RegisterNotifyHandlerCalled func(handler vmcommon.EpochSubscriberHandler)
}

// RegisterNotifyHandler -
func (ens *EpochNotifierStub) RegisterNotifyHandler(handler vmcommon.EpochSubscriberHandler) {
	if ens.RegisterNotifyHandlerCalled != nil {
		ens.RegisterNotifyHandlerCalled(handler)
	} else {
		if !check.IfNil(handler) {
			handler.EpochConfirmed(0, 0)
		}
	}
}

// IsGlobalMintBurnFlagEnabled -
func (ens *EpochNotifierStub) IsGlobalMintBurnFlagEnabled() bool {
	return true
}

// IsESDTTransferRoleFlagEnabled -
func (ens *EpochNotifierStub) IsESDTTransferRoleFlagEnabled() bool {
	return true
}

// IsBuiltInFunctionsFlagEnabled -
func (ens *EpochNotifierStub) IsBuiltInFunctionsFlagEnabled() bool {
	return true
}

// IsCheckCorrectTokenIDForTransferRoleFlagEnabled -
func (ens *EpochNotifierStub) IsCheckCorrectTokenIDForTransferRoleFlagEnabled() bool {
	return true
}

// IsMultiESDTTransferFixOnCallBackFlagEnabled -
func (ens *EpochNotifierStub) IsMultiESDTTransferFixOnCallBackFlagEnabled() bool {
	return true
}

// IsFixOOGReturnCodeFlagEnabled -
func (ens *EpochNotifierStub) IsFixOOGReturnCodeFlagEnabled() bool {
	return true
}

// IsRemoveNonUpdatedStorageFlagEnabled -
func (ens *EpochNotifierStub) IsRemoveNonUpdatedStorageFlagEnabled() bool {
	return true
}

// IsCreateNFTThroughExecByCallerFlagEnabled -
func (ens *EpochNotifierStub) IsCreateNFTThroughExecByCallerFlagEnabled() bool {
	return true
}

// IsStorageAPICostOptimizationFlagEnabled -
func (ens *EpochNotifierStub) IsStorageAPICostOptimizationFlagEnabled() bool {
	return true
}

// IsFailExecutionOnEveryAPIErrorFlagEnabled -
func (ens *EpochNotifierStub) IsFailExecutionOnEveryAPIErrorFlagEnabled() bool {
	return true
}

// IsManagedCryptoAPIsFlagEnabled -
func (ens *EpochNotifierStub) IsManagedCryptoAPIsFlagEnabled() bool {
	return true
}

// IsSCDeployFlagEnabled -
func (ens *EpochNotifierStub) IsSCDeployFlagEnabled() bool {
	return true
}

// IsAheadOfTimeGasUsageFlagEnabled -
func (ens *EpochNotifierStub) IsAheadOfTimeGasUsageFlagEnabled() bool {
	return true
}

// IsRepairCallbackFlagEnabled -
func (ens *EpochNotifierStub) IsRepairCallbackFlagEnabled() bool {
	return true
}

// IsDisableExecByCallerFlagEnabled -
func (ens *EpochNotifierStub) IsDisableExecByCallerFlagEnabled() bool {
	return true
}

// IsRefactorContextFlagEnabled -
func (ens *EpochNotifierStub) IsRefactorContextFlagEnabled() bool {
	return true
}

// IsCheckFunctionArgumentFlagEnabled -
func (ens *EpochNotifierStub) IsCheckFunctionArgumentFlagEnabled() bool {
	return true
}

// IsCheckExecuteOnReadOnlyFlagEnabled -
func (ens *EpochNotifierStub) IsCheckExecuteOnReadOnlyFlagEnabled() bool {
	return true
}

// IsFixAsyncCallbackCheckFlagEnabled -
func (ens *EpochNotifierStub) IsFixAsyncCallbackCheckFlagEnabled() bool {
	return true
}

// IsSaveToSystemAccountFlagEnabled -
func (ens *EpochNotifierStub) IsSaveToSystemAccountFlagEnabled() bool {
	return true
}

// IsCheckFrozenCollectionFlagEnabled -
func (ens *EpochNotifierStub) IsCheckFrozenCollectionFlagEnabled() bool {
	return true
}

// IsSendAlwaysFlagEnabled -
func (ens *EpochNotifierStub) IsSendAlwaysFlagEnabled() bool {
	return true
}

// IsValueLengthCheckFlagEnabled -
func (ens *EpochNotifierStub) IsValueLengthCheckFlagEnabled() bool {
	return true
}

// IsCheckTransferFlagEnabled -
func (ens *EpochNotifierStub) IsCheckTransferFlagEnabled() bool {
	return true
}

// IsTransferToMetaFlagEnabled -
func (ens *EpochNotifierStub) IsTransferToMetaFlagEnabled() bool {
	return true
}

// IsESDTNFTImprovementV1FlagEnabled -
func (ens *EpochNotifierStub) IsESDTNFTImprovementV1FlagEnabled() bool {
	return true
}

// MultiESDTTransferAsyncCallBackEnableEpoch -
func (ens *EpochNotifierStub) MultiESDTTransferAsyncCallBackEnableEpoch() uint32 {
	return 0
}

// FixOOGReturnCodeEnableEpoch -
func (ens *EpochNotifierStub) FixOOGReturnCodeEnableEpoch() uint32 {
	return 0
}

// RemoveNonUpdatedStorageEnableEpoch -
func (ens *EpochNotifierStub) RemoveNonUpdatedStorageEnableEpoch() uint32 {
	return 0
}

// CreateNFTThroughExecByCallerEnableEpoch -
func (ens *EpochNotifierStub) CreateNFTThroughExecByCallerEnableEpoch() uint32 {
	return 0
}

// FixFailExecutionOnErrorEnableEpoch -
func (ens *EpochNotifierStub) FixFailExecutionOnErrorEnableEpoch() uint32 {
	return 0
}

// ManagedCryptoAPIEnableEpoch -
func (ens *EpochNotifierStub) ManagedCryptoAPIEnableEpoch() uint32 {
	return 0
}

// DisableExecByCallerEnableEpoch -
func (ens *EpochNotifierStub) DisableExecByCallerEnableEpoch() uint32 {
	return 0
}

// RefactorContextEnableEpoch -
func (ens *EpochNotifierStub) RefactorContextEnableEpoch() uint32 {
	return 0
}

// CheckExecuteReadOnlyEnableEpoch -
func (ens *EpochNotifierStub) CheckExecuteReadOnlyEnableEpoch() uint32 {
	return 0
}

// StorageAPICostOptimizationEnableEpoch -
func (ens *EpochNotifierStub) StorageAPICostOptimizationEnableEpoch() uint32 {
	return 0
}

// IsInterfaceNil -
func (ens *EpochNotifierStub) IsInterfaceNil() bool {
	return ens == nil
}
