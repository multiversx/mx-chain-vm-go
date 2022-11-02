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

func (ens *EpochNotifierStub) IsGlobalMintBurnFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsESDTTransferRoleFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsBuiltInFunctionsFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsCheckCorrectTokenIDForTransferRoleFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsMultiESDTTransferFixOnCallBackFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsFixOOGReturnCodeFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsRemoveNonUpdatedStorageFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsCreateNFTThroughExecByCallerFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsStorageAPICostOptimizationFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsFailExecutionOnEveryAPIErrorFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsManagedCryptoAPIsFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsSCDeployFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsAheadOfTimeGasUsageFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsRepairCallbackFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsDisableExecByCallerFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsRefactorContextFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsCheckFunctionArgumentFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsCheckExecuteOnReadOnlyFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsFixAsyncCallbackCheckFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsSaveToSystemAccountFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsCheckFrozenCollectionFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsSendAlwaysFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsValueLengthCheckFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsCheckTransferFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsTransferToMetaFlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) IsESDTNFTImprovementV1FlagEnabled() bool {
	return true
}

func (ens *EpochNotifierStub) MultiESDTTransferAsyncCallBackEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) FixOOGReturnCodeEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) RemoveNonUpdatedStorageEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) CreateNFTThroughExecByCallerEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) FixFailExecutionOnErrorEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) ManagedCryptoAPIEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) DisableExecByCallerEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) RefactorContextEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) CheckExecuteReadOnlyEnableEpoch() uint32 {
	return 0
}

func (ens *EpochNotifierStub) StorageAPICostOptimizationEnableEpoch() uint32 {
	return 0
}

// IsInterfaceNil -
func (ens *EpochNotifierStub) IsInterfaceNil() bool {
	return ens == nil
}
