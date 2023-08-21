package vmhost

import "github.com/multiversx/mx-chain-core-go/core"

const (
	MultiESDTTransferFixOnCallBackFlag core.EnableEpochFlag = "MultiESDTTransferFixOnCallBackFlag"
	RemoveNonUpdatedStorageFlag        core.EnableEpochFlag = "RemoveNonUpdatedStorageFlag"
	CreateNFTThroughExecByCallerFlag   core.EnableEpochFlag = "CreateNFTThroughExecByCallerFlag"
	StorageAPICostOptimizationFlag     core.EnableEpochFlag = "StorageAPICostOptimizationFlag"
	CheckExecuteOnReadOnlyFlag         core.EnableEpochFlag = "CheckExecuteOnReadOnlyFlag"
	FailExecutionOnEveryAPIErrorFlag   core.EnableEpochFlag = "FailExecutionOnEveryAPIErrorFlag"
	ManagedCryptoAPIsFlag              core.EnableEpochFlag = "ManagedCryptoAPIsFlag"
	DisableExecByCallerFlag            core.EnableEpochFlag = "DisableExecByCallerFlag"
	RefactorContextFlag                core.EnableEpochFlag = "RefactorContextFlag"
	RuntimeMemStoreLimitFlag           core.EnableEpochFlag = "RuntimeMemStoreLimitFlag"
	RuntimeCodeSizeFixFlag             core.EnableEpochFlag = "RuntimeCodeSizeFixFlag"
	FixOOGReturnCodeFlag               core.EnableEpochFlag = "FixOOGReturnCodeFlag"
)
