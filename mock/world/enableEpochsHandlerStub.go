package worldmock

import vmcommon "github.com/ElrondNetwork/elrond-vm-common"

var _ vmcommon.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	CheckExecuteReadOnlyEnableEpochField                 uint32
	CreateNFTThroughExecByCallerEnableEpochField         uint32
	DisableExecByCallerEnableEpochField                  uint32
	FixFailExecutionOnErrorEnableEpochField              uint32
	FixOOGReturnCodeEnableEpochField                     uint32
	IsAheadOfTimeGasUsageFlagEnabledField                bool
	IsBuiltInFunctionsFlagEnabledField                   bool
	IsCheckCorrectTokenIDForTransferRoleFlagEnabledField bool
	IsCheckExecuteOnReadOnlyFlagEnabledField             bool
	IsCheckFrozenCollectionFlagEnabledField              bool
	IsCheckFunctionArgumentFlagEnabledField              bool
	IsCheckTransferFlagEnabledField                      bool
	IsCreateNFTThroughExecByCallerFlagEnabledField       bool
	IsDisableExecByCallerFlagEnabledField                bool
	IsESDTNFTImprovementV1FlagEnabledField               bool
	IsESDTTransferRoleFlagEnabledField                   bool
	IsFailExecutionOnEveryAPIErrorFlagEnabledField       bool
	IsFixAsyncCallbackCheckFlagEnabledField              bool
	IsFixOOGReturnCodeFlagEnabledField                   bool
	IsFixOldTokenLiquidityEnabledField                   bool
	IsGlobalMintBurnFlagEnabledField                     bool
	IsManagedCryptoAPIsFlagEnabledField                  bool
	IsMultiESDTTransferFixOnCallBackFlagEnabledField     bool
	IsRefactorContextFlagEnabledField                    bool
	IsRemoveNonUpdatedStorageFlagEnabledField            bool
	IsRepairCallbackFlagEnabledField                     bool
	IsSCDeployFlagEnabledField                           bool
	IsSaveToSystemAccountFlagEnabledField                bool
	IsSendAlwaysFlagEnabledField                         bool
	IsStorageAPICostOptimizationFlagEnabledField         bool
	IsTransferToMetaFlagEnabledField                     bool
	IsValueLengthCheckFlagEnabledField                   bool
	ManagedCryptoAPIEnableEpochField                     uint32
	MultiESDTTransferAsyncCallBackEnableEpochField       uint32
	RefactorContextEnableEpochField                      uint32
	RemoveNonUpdatedStorageEnableEpochField              uint32
	StorageAPICostOptimizationEnableEpochField           uint32
	IsRuntimeMemStoreLimitEnabledField                   bool
	IsMaxBlockchainHookCountersFlagEnabledField          bool
	IsWipeSingleNFTLiquidityDecreaseEnabledField         bool
}

// IsGlobalMintBurnFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsGlobalMintBurnFlagEnabled() bool {
	return stub.IsGlobalMintBurnFlagEnabledField
}

// IsESDTTransferRoleFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsESDTTransferRoleFlagEnabled() bool {
	return stub.IsESDTTransferRoleFlagEnabledField
}

// IsBuiltInFunctionsFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsBuiltInFunctionsFlagEnabled() bool {
	return stub.IsBuiltInFunctionsFlagEnabledField
}

// IsCheckCorrectTokenIDForTransferRoleFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsCheckCorrectTokenIDForTransferRoleFlagEnabled() bool {
	return stub.IsCheckCorrectTokenIDForTransferRoleFlagEnabledField
}

// IsMultiESDTTransferFixOnCallBackFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsMultiESDTTransferFixOnCallBackFlagEnabled() bool {
	return stub.IsMultiESDTTransferFixOnCallBackFlagEnabledField
}

// IsFixOOGReturnCodeFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsFixOOGReturnCodeFlagEnabled() bool {
	return stub.IsFixOOGReturnCodeFlagEnabledField
}

// IsRemoveNonUpdatedStorageFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsRemoveNonUpdatedStorageFlagEnabled() bool {
	return stub.IsRemoveNonUpdatedStorageFlagEnabledField
}

// IsCreateNFTThroughExecByCallerFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsCreateNFTThroughExecByCallerFlagEnabled() bool {
	return stub.IsCreateNFTThroughExecByCallerFlagEnabledField
}

// IsStorageAPICostOptimizationFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsStorageAPICostOptimizationFlagEnabled() bool {
	return stub.IsStorageAPICostOptimizationFlagEnabledField
}

// IsFailExecutionOnEveryAPIErrorFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsFailExecutionOnEveryAPIErrorFlagEnabled() bool {
	return stub.IsFailExecutionOnEveryAPIErrorFlagEnabledField
}

// IsManagedCryptoAPIsFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsManagedCryptoAPIsFlagEnabled() bool {
	return stub.IsManagedCryptoAPIsFlagEnabledField
}

// IsSCDeployFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsSCDeployFlagEnabled() bool {
	return stub.IsSCDeployFlagEnabledField
}

// IsAheadOfTimeGasUsageFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsAheadOfTimeGasUsageFlagEnabled() bool {
	return stub.IsAheadOfTimeGasUsageFlagEnabledField
}

// IsRepairCallbackFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsRepairCallbackFlagEnabled() bool {
	return stub.IsRepairCallbackFlagEnabledField
}

// IsDisableExecByCallerFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsDisableExecByCallerFlagEnabled() bool {
	return stub.IsDisableExecByCallerFlagEnabledField
}

// IsRefactorContextFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsRefactorContextFlagEnabled() bool {
	return stub.IsRefactorContextFlagEnabledField
}

// IsCheckFunctionArgumentFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsCheckFunctionArgumentFlagEnabled() bool {
	return stub.IsCheckFunctionArgumentFlagEnabledField
}

// IsCheckExecuteOnReadOnlyFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsCheckExecuteOnReadOnlyFlagEnabled() bool {
	return stub.IsCheckExecuteOnReadOnlyFlagEnabledField
}

// IsFixAsyncCallbackCheckFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsFixAsyncCallbackCheckFlagEnabled() bool {
	return stub.IsFixAsyncCallbackCheckFlagEnabledField
}

// IsSaveToSystemAccountFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsSaveToSystemAccountFlagEnabled() bool {
	return stub.IsSaveToSystemAccountFlagEnabledField
}

// IsCheckFrozenCollectionFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsCheckFrozenCollectionFlagEnabled() bool {
	return stub.IsCheckFrozenCollectionFlagEnabledField
}

// IsSendAlwaysFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsSendAlwaysFlagEnabled() bool {
	return stub.IsSendAlwaysFlagEnabledField
}

// IsValueLengthCheckFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsValueLengthCheckFlagEnabled() bool {
	return stub.IsValueLengthCheckFlagEnabledField
}

// IsCheckTransferFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsCheckTransferFlagEnabled() bool {
	return stub.IsCheckTransferFlagEnabledField
}

// IsTransferToMetaFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsTransferToMetaFlagEnabled() bool {
	return stub.IsTransferToMetaFlagEnabledField
}

// IsESDTNFTImprovementV1FlagEnabled -
func (stub *EnableEpochsHandlerStub) IsESDTNFTImprovementV1FlagEnabled() bool {
	return stub.IsESDTNFTImprovementV1FlagEnabledField
}

// IsFixOldTokenLiquidityEnabled -
func (stub *EnableEpochsHandlerStub) IsFixOldTokenLiquidityEnabled() bool {
	return stub.IsFixOldTokenLiquidityEnabledField
}

// MultiESDTTransferAsyncCallBackEnableEpoch -
func (stub *EnableEpochsHandlerStub) MultiESDTTransferAsyncCallBackEnableEpoch() uint32 {
	return stub.MultiESDTTransferAsyncCallBackEnableEpochField
}

// FixOOGReturnCodeEnableEpoch -
func (stub *EnableEpochsHandlerStub) FixOOGReturnCodeEnableEpoch() uint32 {
	return stub.FixOOGReturnCodeEnableEpochField
}

// RemoveNonUpdatedStorageEnableEpoch -
func (stub *EnableEpochsHandlerStub) RemoveNonUpdatedStorageEnableEpoch() uint32 {
	return stub.RemoveNonUpdatedStorageEnableEpochField
}

// CreateNFTThroughExecByCallerEnableEpoch -
func (stub *EnableEpochsHandlerStub) CreateNFTThroughExecByCallerEnableEpoch() uint32 {
	return stub.CreateNFTThroughExecByCallerEnableEpochField
}

// FixFailExecutionOnErrorEnableEpoch -
func (stub *EnableEpochsHandlerStub) FixFailExecutionOnErrorEnableEpoch() uint32 {
	return stub.FixFailExecutionOnErrorEnableEpochField
}

// ManagedCryptoAPIEnableEpoch -
func (stub *EnableEpochsHandlerStub) ManagedCryptoAPIEnableEpoch() uint32 {
	return stub.ManagedCryptoAPIEnableEpochField
}

// DisableExecByCallerEnableEpoch -
func (stub *EnableEpochsHandlerStub) DisableExecByCallerEnableEpoch() uint32 {
	return stub.DisableExecByCallerEnableEpochField
}

// RefactorContextEnableEpoch -
func (stub *EnableEpochsHandlerStub) RefactorContextEnableEpoch() uint32 {
	return stub.RefactorContextEnableEpochField
}

// CheckExecuteReadOnlyEnableEpoch -
func (stub *EnableEpochsHandlerStub) CheckExecuteReadOnlyEnableEpoch() uint32 {
	return stub.CheckExecuteReadOnlyEnableEpochField
}

// StorageAPICostOptimizationEnableEpoch -
func (stub *EnableEpochsHandlerStub) StorageAPICostOptimizationEnableEpoch() uint32 {
	return stub.StorageAPICostOptimizationEnableEpochField
}

// IsRuntimeMemStoreLimitEnabled -
func (stub *EnableEpochsHandlerStub) IsRuntimeMemStoreLimitEnabled() bool {
	return stub.IsRuntimeMemStoreLimitEnabledField
}

// IsMaxBlockchainHookCountersFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsMaxBlockchainHookCountersFlagEnabled() bool {
	return stub.IsMaxBlockchainHookCountersFlagEnabledField
}

// IsWipeSingleNFTLiquidityDecreaseEnabled -
func (stub *EnableEpochsHandlerStub) IsWipeSingleNFTLiquidityDecreaseEnabled() bool {
	return stub.IsWipeSingleNFTLiquidityDecreaseEnabledField
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}

// EnableEpochsHandlerStubAllFlags creates a new EnableEpochsHandlerStub with all flags enabled
func EnableEpochsHandlerStubAllFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField:         true,
		IsMultiESDTTransferFixOnCallBackFlagEnabledField:     true,
		IsFixOOGReturnCodeFlagEnabledField:                   true,
		IsRemoveNonUpdatedStorageFlagEnabledField:            true,
		IsCreateNFTThroughExecByCallerFlagEnabledField:       true,
		IsManagedCryptoAPIsFlagEnabledField:                  true,
		IsFailExecutionOnEveryAPIErrorFlagEnabledField:       true,
		IsRefactorContextFlagEnabledField:                    true,
		IsCheckCorrectTokenIDForTransferRoleFlagEnabledField: true,
		IsDisableExecByCallerFlagEnabledField:                true,
		IsESDTTransferRoleFlagEnabledField:                   true,
		IsGlobalMintBurnFlagEnabledField:                     true,
		IsTransferToMetaFlagEnabledField:                     true,
		IsCheckFrozenCollectionFlagEnabledField:              true,
		IsFixAsyncCallbackCheckFlagEnabledField:              true,
		IsESDTNFTImprovementV1FlagEnabledField:               true,
		IsSaveToSystemAccountFlagEnabledField:                true,
		IsValueLengthCheckFlagEnabledField:                   true,
		IsSCDeployFlagEnabledField:                           true,
		IsRepairCallbackFlagEnabledField:                     true,
		IsAheadOfTimeGasUsageFlagEnabledField:                true,
		IsCheckFunctionArgumentFlagEnabledField:              true,
		IsCheckExecuteOnReadOnlyFlagEnabledField:             true,
		IsFixOldTokenLiquidityEnabledField:                   true,
	}
}

// EnableEpochsHandlerStubNoFlags creates a new EnableEpochsHandlerStub with all flags disabled
func EnableEpochsHandlerStubNoFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{}
}
