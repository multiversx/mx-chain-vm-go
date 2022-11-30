package host

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/contexts"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/cryptoapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/crypto/factory"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
	"github.com/ElrondNetwork/elrond-go-core/core/atomic"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var log = logger.GetOrCreate("arwen/host")
var logGasTrace = logger.GetOrCreate("gasTrace")

// MaximumWasmerInstanceCount represents the maximum number of Wasmer instances that can be active at the same time
var MaximumWasmerInstanceCount = uint64(10)

var _ arwen.VMHost = (*vmHost)(nil)

const minExecutionTimeout = time.Second
const internalVMErrors = "internalVMErrors"

// vmHost implements HostContext interface.
type vmHost struct {
	cryptoHook       crypto.VMCrypto
	mutExecution     sync.RWMutex
	closingInstance  bool
	executionTimeout time.Duration

	ethInput []byte

	blockchainContext   arwen.BlockchainContext
	runtimeContext      arwen.RuntimeContext
	outputContext       arwen.OutputContext
	meteringContext     arwen.MeteringContext
	storageContext      arwen.StorageContext
	managedTypesContext arwen.ManagedTypesContext

	gasSchedule          config.GasScheduleMap
	scAPIMethods         *wasmer.Imports
	builtInFuncContainer vmcommon.BuiltInFunctionContainer
	esdtTransferParser   vmcommon.ESDTTransferParser

	multiESDTTransferAsyncCallBackEnableEpoch uint32
	flagMultiESDTTransferAsyncCallBack        atomic.Flag

	fixOOGReturnCodeEnableEpoch uint32
	flagFixOOGReturnCode        atomic.Flag

	removeNonUpdatedStorageEnableEpoch uint32
	flagRemoveNonUpdatedStorage        atomic.Flag

	createNFTThroughExecByCallerEnableEpoch uint32
	flagCreateNFTThroughExecByCaller        atomic.Flag

	fixFailExecutionOnErrorEnableEpoch uint32
	flagFixFailExecutionOnError        atomic.Flag

	useDifferentGasCostForReadingCachedStorageEpoch uint32
	flagUseDifferentGasCostForCachedStorage         atomic.Flag

	disableExecByCallerEnableEpoch uint32
	flagDisableExecByCaller        atomic.Flag

	refactorContextEnableEpoch uint32
	flagRefactorContext        atomic.Flag

	fixAsnycCallArgumentsEnableEpoch uint32
	flagFixAsyncCallArguments        atomic.Flag

	checkExecuteReadOnlyEnableEpoch     uint32
	flagCheckExecuteReadOnlyEnableEpoch atomic.Flag

	activationEpochMap map[uint32]struct{}
}

// NewArwenVM creates a new Arwen vmHost
func NewArwenVM(
	blockChainHook vmcommon.BlockchainHook,
	hostParameters *arwen.VMHostParameters,
) (arwen.VMHost, error) {

	if check.IfNil(blockChainHook) {
		return nil, arwen.ErrNilBlockChainHook
	}
	if hostParameters == nil {
		return nil, arwen.ErrNilHostParameters
	}
	if check.IfNil(hostParameters.ESDTTransferParser) {
		return nil, arwen.ErrNilESDTTransferParser
	}
	if check.IfNil(hostParameters.BuiltInFuncContainer) {
		return nil, arwen.ErrNilBuiltInFunctionsContainer
	}
	if check.IfNil(hostParameters.EpochNotifier) {
		return nil, arwen.ErrNilEpochNotifier
	}

	cryptoHook := factory.NewVMCrypto()
	host := &vmHost{
		cryptoHook:           cryptoHook,
		meteringContext:      nil,
		runtimeContext:       nil,
		blockchainContext:    nil,
		storageContext:       nil,
		managedTypesContext:  nil,
		gasSchedule:          hostParameters.GasSchedule,
		scAPIMethods:         nil,
		builtInFuncContainer: hostParameters.BuiltInFuncContainer,
		esdtTransferParser:   hostParameters.ESDTTransferParser,
		executionTimeout:     minExecutionTimeout,
		multiESDTTransferAsyncCallBackEnableEpoch:       hostParameters.MultiESDTTransferAsyncCallBackEnableEpoch,
		fixOOGReturnCodeEnableEpoch:                     hostParameters.FixOOGReturnCodeEnableEpoch,
		removeNonUpdatedStorageEnableEpoch:              hostParameters.RemoveNonUpdatedStorageEnableEpoch,
		createNFTThroughExecByCallerEnableEpoch:         hostParameters.CreateNFTThroughExecByCallerEnableEpoch,
		fixFailExecutionOnErrorEnableEpoch:              hostParameters.FixFailExecutionOnErrorEnableEpoch,
		useDifferentGasCostForReadingCachedStorageEpoch: hostParameters.UseDifferentGasCostForReadingCachedStorageEpoch,
		disableExecByCallerEnableEpoch:                  hostParameters.DisableExecByCallerEnableEpoch,
		refactorContextEnableEpoch:                      hostParameters.RefactorContextEnableEpoch,
		fixAsnycCallArgumentsEnableEpoch:                hostParameters.ManagedCryptoAPIEnableEpoch,
		checkExecuteReadOnlyEnableEpoch:                 hostParameters.CheckExecuteReadOnlyEnableEpoch,
	}

	host.activationEpochMap = createActivationMap(hostParameters)

	newExecutionTimeout := time.Duration(hostParameters.TimeOutForSCExecutionInMilliseconds) * time.Millisecond
	if newExecutionTimeout > minExecutionTimeout {
		host.executionTimeout = newExecutionTimeout
	}
	imports, err := elrondapi.ElrondEIImports()
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.BigIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.BigFloatImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.SmallIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.ManagedEIImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.ManagedBufferImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = cryptoapi.CryptoImports(imports)
	if err != nil {
		return nil, err
	}

	err = wasmer.SetImports(imports)
	if err != nil {
		return nil, err
	}

	host.scAPIMethods = imports

	host.blockchainContext, err = contexts.NewBlockchainContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeContext, err = contexts.NewRuntimeContext(
		host,
		hostParameters.VMType,
		host.builtInFuncContainer,
		hostParameters.EpochNotifier,
		hostParameters.UseDifferentGasCostForReadingCachedStorageEpoch,
		hostParameters.ManagedCryptoAPIEnableEpoch,
		hostParameters.RuntimeStoreLimitEnableEpoch,
	)
	if err != nil {
		return nil, err
	}

	host.meteringContext, err = contexts.NewMeteringContext(host, hostParameters.GasSchedule, hostParameters.BlockGasLimit)
	if err != nil {
		return nil, err
	}

	host.outputContext, err = contexts.NewOutputContext(host)
	if err != nil {
		return nil, err
	}

	host.storageContext, err = contexts.NewStorageContext(
		host,
		blockChainHook,
		hostParameters.EpochNotifier,
		hostParameters.ElrondProtectedKeyPrefix,
		hostParameters.UseDifferentGasCostForReadingCachedStorageEpoch,
	)
	if err != nil {
		return nil, err
	}

	host.managedTypesContext, err = contexts.NewManagedTypesContext(host)
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(host.gasSchedule)
	if err != nil {
		return nil, err
	}

	host.runtimeContext.SetMaxInstanceCount(MaximumWasmerInstanceCount)

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)
	wasmer.SetRkyvSerializationEnabled(true)

	if hostParameters.WasmerSIGSEGVPassthrough {
		wasmer.SetSIGSEGVPassthrough()
	}

	host.initContexts()
	hostParameters.EpochNotifier.RegisterNotifyHandler(host)

	return host, nil
}

func createActivationMap(hostParameters *arwen.VMHostParameters) map[uint32]struct{} {
	activationMap := make(map[uint32]struct{})

	activationMap[hostParameters.CheckExecuteReadOnlyEnableEpoch] = struct{}{}
	activationMap[hostParameters.DisableExecByCallerEnableEpoch] = struct{}{}
	activationMap[hostParameters.RefactorContextEnableEpoch] = struct{}{}
	activationMap[hostParameters.FixFailExecutionOnErrorEnableEpoch] = struct{}{}
	activationMap[hostParameters.ManagedCryptoAPIEnableEpoch] = struct{}{}
	activationMap[hostParameters.CreateNFTThroughExecByCallerEnableEpoch] = struct{}{}
	activationMap[hostParameters.FixOOGReturnCodeEnableEpoch] = struct{}{}
	activationMap[hostParameters.MultiESDTTransferAsyncCallBackEnableEpoch] = struct{}{}
	activationMap[hostParameters.RemoveNonUpdatedStorageEnableEpoch] = struct{}{}
	activationMap[hostParameters.UseDifferentGasCostForReadingCachedStorageEpoch] = struct{}{}
	activationMap[hostParameters.RuntimeStoreLimitEnableEpoch] = struct{}{}

	return activationMap
}

// GetVersion returns the Arwen version string
func (host *vmHost) GetVersion() string {
	return arwen.ArwenVersion
}

// Crypto returns the VMCrypto instance of the host
func (host *vmHost) Crypto() crypto.VMCrypto {
	return host.cryptoHook
}

// Blockchain returns the BlockchainContext instance of the host
func (host *vmHost) Blockchain() arwen.BlockchainContext {
	return host.blockchainContext
}

// Runtime returns the RuntimeContext instance of the host
func (host *vmHost) Runtime() arwen.RuntimeContext {
	return host.runtimeContext
}

// Output returns the OutputContext instance of the host
func (host *vmHost) Output() arwen.OutputContext {
	return host.outputContext
}

// Metering returns the MeteringContext instance of the host
func (host *vmHost) Metering() arwen.MeteringContext {
	return host.meteringContext
}

// Storage returns the StorageContext instance of the host
func (host *vmHost) Storage() arwen.StorageContext {
	return host.storageContext
}

// ManagedTypes returns the ManagedTypeContext instance of the host
func (host *vmHost) ManagedTypes() arwen.ManagedTypesContext {
	return host.managedTypesContext
}

// GetContexts returns the main contexts of the host
func (host *vmHost) GetContexts() (
	arwen.ManagedTypesContext,
	arwen.BlockchainContext,
	arwen.MeteringContext,
	arwen.OutputContext,
	arwen.RuntimeContext,
	arwen.StorageContext,
) {
	return host.managedTypesContext,
		host.blockchainContext,
		host.meteringContext,
		host.outputContext,
		host.runtimeContext,
		host.storageContext
}

// InitState resets the contexts of the host and reconfigures its flags
func (host *vmHost) InitState() {
	host.initContexts()
}

func (host *vmHost) close() {
	host.runtimeContext.ClearWarmInstanceCache()
}

// Close will close all underlying processes
func (host *vmHost) Close() error {
	host.mutExecution.Lock()
	host.close()
	host.closingInstance = true
	host.mutExecution.Unlock()

	return nil
}

// Reset is a function which closes the VM and resets the closingInstance variable
func (host *vmHost) Reset() {
	host.mutExecution.Lock()
	host.close()
	// keep closingInstance flag to false
	host.mutExecution.Unlock()
}

func (host *vmHost) initContexts() {
	host.ClearContextStateStack()
	host.managedTypesContext.InitState()
	host.outputContext.InitState()
	host.meteringContext.InitState()
	host.runtimeContext.InitState()
	host.storageContext.InitState()
	host.ethInput = nil
}

// ClearContextStateStack cleans the state stacks of all the contexts of the host
func (host *vmHost) ClearContextStateStack() {
	host.managedTypesContext.ClearStateStack()
	host.outputContext.ClearStateStack()
	host.meteringContext.ClearStateStack()
	host.runtimeContext.ClearStateStack()
	host.storageContext.ClearStateStack()
}

// GetAPIMethods returns the EEI as a set of imports for Wasmer
func (host *vmHost) GetAPIMethods() *wasmer.Imports {
	return host.scAPIMethods
}

// GasScheduleChange applies a new gas schedule to the host
func (host *vmHost) GasScheduleChange(newGasSchedule config.GasScheduleMap) {
	host.mutExecution.Lock()
	defer host.mutExecution.Unlock()

	host.gasSchedule = newGasSchedule
	gasCostConfig, err := config.CreateGasConfig(newGasSchedule)
	if err != nil {
		log.Error("cannot apply new gas config", "err", err)
		return
	}

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.meteringContext.SetGasSchedule(newGasSchedule)
	host.runtimeContext.ClearWarmInstanceCache()
}

// GetGasScheduleMap returns the currently stored gas schedule
func (host *vmHost) GetGasScheduleMap() config.GasScheduleMap {
	return host.gasSchedule
}

// RunSmartContractCreate executes the deployment of a new contract
func (host *vmHost) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	if host.closingInstance {
		return nil, arwen.ErrVMIsClosing
	}

	host.setGasTracerEnabledIfLogIsTrace()
	ctx, cancel := context.WithTimeout(context.Background(), host.executionTimeout)
	defer cancel()

	log.Trace("RunSmartContractCreate begin",
		"len(code)", len(input.ContractCode),
		"metadata", input.ContractCodeMetadata,
		"gasProvided", input.GasProvided,
		"gasLocked", input.GasLocked)

	done := make(chan struct{})
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Error("VM execution panicked", "error", r, "stack", "\n"+string(debug.Stack()))
				err = arwen.ErrExecutionPanicked
			}
			host.Runtime().CleanInstance()
			close(done)
		}()

		vmOutput = host.doRunSmartContractCreate(input)
		logsFromErrors := host.createLogEntryFromErrors(input.CallerAddr, input.CallerAddr, "_init")
		if logsFromErrors != nil {
			vmOutput.Logs = append(vmOutput.Logs, logsFromErrors)
		}

		log.Trace("RunSmartContractCreate end",
			"returnCode", vmOutput.ReturnCode,
			"returnMessage", vmOutput.ReturnMessage,
			"gasRemaining", vmOutput.GasRemaining)
		host.logFromGasTracer("init")
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		host.Runtime().FailExecution(arwen.ErrExecutionFailedWithTimeout)
		<-done
		err = arwen.ErrExecutionFailedWithTimeout
	}

	return
}

// RunSmartContractCall executes the call of an existing contract
func (host *vmHost) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	if host.closingInstance {
		return nil, arwen.ErrVMIsClosing
	}

	host.setGasTracerEnabledIfLogIsTrace()
	ctx, cancel := context.WithTimeout(context.Background(), host.executionTimeout)
	defer cancel()

	log.Trace("RunSmartContractCall begin",
		"function", input.Function,
		"gasProvided", input.GasProvided,
		"gasLocked", input.GasLocked)

	done := make(chan struct{})
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Error("VM execution panicked", "error", r, "stack", "\n"+string(debug.Stack()))
				err = arwen.ErrExecutionPanicked
			}

			host.Runtime().CleanInstance()
			close(done)
		}()

		isUpgrade := input.Function == arwen.UpgradeFunctionName
		if isUpgrade {
			vmOutput = host.doRunSmartContractUpgrade(input)
		} else {
			vmOutput = host.doRunSmartContractCall(input)
		}

		logsFromErrors := host.createLogEntryFromErrors(input.CallerAddr, input.RecipientAddr, input.Function)
		if logsFromErrors != nil {
			vmOutput.Logs = append(vmOutput.Logs, logsFromErrors)
		}

		log.Trace("RunSmartContractCall end",
			"function", input.Function,
			"returnCode", vmOutput.ReturnCode,
			"returnMessage", vmOutput.ReturnMessage,
			"gasRemaining", vmOutput.GasRemaining)
		host.logFromGasTracer(input.Function)
	}()

	select {
	case <-done:
		// Normal termination.
		return
	case <-ctx.Done():
		// Terminated due to timeout. The VM sets the `ExecutionFailed` breakpoint
		// in Wasmer. Also, the VM must wait for Wasmer to reach the end of a WASM
		// basic block in order to close the WASM instance cleanly. This is done by
		// reading the `done` channel once more, awaiting the call to `close(done)`
		// from above.
		host.Runtime().FailExecution(arwen.ErrExecutionFailedWithTimeout)
		<-done
		err = arwen.ErrExecutionFailedWithTimeout
	}

	return
}

func (host *vmHost) createLogEntryFromErrors(sndAddress, rcvAddress []byte, function string) *vmcommon.LogEntry {
	formattedErrors := host.runtimeContext.GetAllErrors()
	if formattedErrors == nil {
		return nil
	}

	logFromError := &vmcommon.LogEntry{
		Identifier: []byte(internalVMErrors),
		Address:    sndAddress,
		Topics:     [][]byte{rcvAddress, []byte(function)},
		Data:       []byte(formattedErrors.Error()),
	}

	return logFromError
}

// AreInSameShard returns true if the provided addresses are part of the same shard
func (host *vmHost) AreInSameShard(leftAddress []byte, rightAddress []byte) bool {
	blockchain := host.Blockchain()
	leftShard := blockchain.GetShardOfAddress(leftAddress)
	rightShard := blockchain.GetShardOfAddress(rightAddress)

	return leftShard == rightShard
}

// IsInterfaceNil returns true if there is no value under the interface
func (host *vmHost) IsInterfaceNil() bool {
	return host == nil
}

// SetRuntimeContext sets the runtimeContext for this host, used in tests
func (host *vmHost) SetRuntimeContext(runtime arwen.RuntimeContext) {
	host.runtimeContext = runtime
}

// GetRuntimeErrors obtains the cumultated error object after running the SC
func (host *vmHost) GetRuntimeErrors() error {
	if host.runtimeContext != nil {
		return host.runtimeContext.GetAllErrors()
	}
	return nil
}

// SetBuiltInFunctionsContainer sets the built in function container - only for testing
func (host *vmHost) SetBuiltInFunctionsContainer(builtInFuncs vmcommon.BuiltInFunctionContainer) {
	if check.IfNil(builtInFuncs) {
		return
	}
	host.builtInFuncContainer = builtInFuncs
}

// EpochConfirmed is called whenever a new epoch is confirmed
func (host *vmHost) EpochConfirmed(epoch uint32, _ uint64) {
	host.flagMultiESDTTransferAsyncCallBack.SetValue(epoch >= host.multiESDTTransferAsyncCallBackEnableEpoch)
	log.Debug("Arwen VM: multi esdt transfer on async callback intra shard", "enabled", host.flagMultiESDTTransferAsyncCallBack.IsSet())

	host.flagFixOOGReturnCode.SetValue(epoch >= host.fixOOGReturnCodeEnableEpoch)
	log.Debug("Arwen VM: fix OutOfGas ReturnCode", "enabled", host.flagFixOOGReturnCode.IsSet())

	host.flagRemoveNonUpdatedStorage.SetValue(epoch >= host.removeNonUpdatedStorageEnableEpoch)
	log.Debug("Arwen VM: remove non updated storage", "enabled", host.flagRemoveNonUpdatedStorage.IsSet())

	host.flagCreateNFTThroughExecByCaller.SetValue(epoch >= host.createNFTThroughExecByCallerEnableEpoch)
	log.Debug("Arwen VM: create NFT through exec by caller", "enabled", host.flagCreateNFTThroughExecByCaller.IsSet())

	host.flagFixFailExecutionOnError.SetValue(epoch >= host.fixFailExecutionOnErrorEnableEpoch)
	log.Debug("Arwen VM: fix fail execution on error", "enabled", host.flagFixFailExecutionOnError.IsSet())

	host.flagUseDifferentGasCostForCachedStorage.SetValue(epoch >= host.useDifferentGasCostForReadingCachedStorageEpoch)
	log.Debug("Arwen VM: use different gas costs when reading cached storage", "enabled", host.flagUseDifferentGasCostForCachedStorage.IsSet())

	host.flagDisableExecByCaller.SetValue(epoch >= host.disableExecByCallerEnableEpoch)
	log.Debug("Arwen VM: disable execute by caller endpoints", "enabled", host.flagDisableExecByCaller.IsSet())

	host.flagRefactorContext.SetValue(epoch >= host.refactorContextEnableEpoch)
	log.Debug("Arwen VM: refactor context", "enabled", host.flagRefactorContext.IsSet())

	host.flagFixAsyncCallArguments.SetValue(epoch >= host.fixAsnycCallArgumentsEnableEpoch)
	log.Debug("Arwen VM: fix asnyccall arguments", "enabled", host.flagFixAsyncCallArguments.IsSet())

	host.flagCheckExecuteReadOnlyEnableEpoch.SetValue(epoch >= host.checkExecuteReadOnlyEnableEpoch)
	log.Debug("Arwen VM: check execute read only mode", "enabled", host.flagCheckExecuteReadOnlyEnableEpoch.IsSet())

	_, ok := host.activationEpochMap[epoch]
	if ok {
		host.Runtime().ClearWarmInstanceCache()
		host.Blockchain().ClearCompiledCodes()
	}
}

// FixOOGReturnCodeEnabled returns true if the corresponding flag is set
func (host *vmHost) FixOOGReturnCodeEnabled() bool {
	return host.flagFixOOGReturnCode.IsSet()
}

// FixFailExecutionEnabled returns true if the corresponding flag is set
func (host *vmHost) FixFailExecutionEnabled() bool {
	return host.flagFixFailExecutionOnError.IsSet()
}

// CreateNFTOnExecByCallerEnabled returns true if the corresponding flag is set
func (host *vmHost) CreateNFTOnExecByCallerEnabled() bool {
	return host.flagCreateNFTThroughExecByCaller.IsSet()
}

// DisableExecByCaller returns true if the corresponding flag is set
func (host *vmHost) DisableExecByCaller() bool {
	return host.flagDisableExecByCaller.IsSet()
}

// CheckExecuteReadOnly returns true if the corresponding flag is set
func (host *vmHost) CheckExecuteReadOnly() bool {
	return host.flagCheckExecuteReadOnlyEnableEpoch.IsSet()
}

func (host *vmHost) setGasTracerEnabledIfLogIsTrace() {
	host.Metering().SetGasTracing(false)
	if logGasTrace.GetLevel() == logger.LogTrace {
		host.Metering().SetGasTracing(true)
	}
}

func (host *vmHost) logFromGasTracer(functionName string) {
	if logGasTrace.GetLevel() == logger.LogTrace {
		scGasTrace := host.meteringContext.GetGasTrace()
		totalGasUsedByAPIs := 0
		for scAddress, gasTrace := range scGasTrace {
			logGasTrace.Trace("Gas Trace for", "SC Address", scAddress, "function", functionName)
			for apiName, value := range gasTrace {
				totalGasUsed := uint64(0)
				for _, usedGas := range value {
					totalGasUsed += usedGas
				}
				logGasTrace.Trace("Gas Trace for", "apiName", apiName, "totalGasUsed", totalGasUsed, "numberOfCalls", len(value))
				totalGasUsedByAPIs += int(totalGasUsed)
			}
			logGasTrace.Trace("Gas Trace for", "TotalGasUsedByAPIs", totalGasUsedByAPIs)
		}
	}
}
