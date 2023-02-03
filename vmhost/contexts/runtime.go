package contexts

import (
	"bytes"
	"errors"
	"fmt"
	builtinMath "math"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var logRuntime = logger.GetOrCreate("vm/runtime")

var _ vmhost.RuntimeContext = (*runtimeContext)(nil)

const warmCacheSize = 100

// WarmInstancesEnabled controls the usage of warm instances
const WarmInstancesEnabled = true

type runtimeContext struct {
	host                 vmhost.VMHost
	vmInput              *vmcommon.ContractCallInput
	codeAddress          []byte
	codeSize             uint64
	callFunction         string
	vmType               []byte
	readOnly             bool
	verifyCode           bool
	maxInstanceStackSize uint64

	vmExecutor executor.Executor

	iTracker *instanceTracker

	stateStack []*runtimeContext

	validator *wasmValidator
	errors    vmhost.WrappableError
	hasher    vmhost.HashComputer
}

// NewRuntimeContext creates a new runtimeContext
func NewRuntimeContext(
	host vmhost.VMHost,
	vmType []byte,
	builtInFuncContainer vmcommon.BuiltInFunctionContainer,
	vmExecutor executor.Executor,
	hasher vmhost.HashComputer,
) (*runtimeContext, error) {

	if check.IfNil(host) {
		return nil, vmhost.ErrNilVMHost
	}
	if check.IfNil(vmExecutor) {
		return nil, vmhost.ErrNilExecutor
	}
	if len(vmType) == 0 {
		return nil, vmhost.ErrNilVMType
	}
	if check.IfNil(builtInFuncContainer) {
		return nil, vmhost.ErrNilBuiltInFunctionsContainer
	}
	if check.IfNil(hasher) {
		return nil, vmhost.ErrNilHasher
	}

	scAPINames := vmExecutor.FunctionNames()

	context := &runtimeContext{
		host:       host,
		vmType:     vmType,
		stateStack: make([]*runtimeContext, 0),
		validator:  newWASMValidator(scAPINames, builtInFuncContainer),
		hasher:     hasher,
		errors:     nil,
	}

	iTracker, err := NewInstanceTracker()
	if err != nil {
		return nil, err
	}
	context.iTracker = iTracker

	context.vmExecutor = vmExecutor
	context.InitState()

	return context, nil
}

// InitState initializes all the contexts fields with default data.
func (context *runtimeContext) InitState() {
	context.vmInput = &vmcommon.ContractCallInput{}
	context.codeAddress = make([]byte, 0)
	context.callFunction = ""
	context.verifyCode = false
	context.readOnly = false
	context.iTracker.InitState()
	context.errors = nil

	logRuntime.Trace("init state")
}

// ClearWarmInstanceCache clears all elements from warm instance cache
func (context *runtimeContext) ClearWarmInstanceCache() {
	context.iTracker.ClearWarmInstanceCache()
	context.iTracker.UnsetInstance()
}

// GetVMExecutor yields the configured contract executor.
func (context *runtimeContext) GetVMExecutor() executor.Executor {
	return context.vmExecutor
}

// ReplaceVMExecutor force-replaces the current executor.
// Do not use in production code, it is only suited for tests! The executor should normally be immutable, once initialized.
func (context *runtimeContext) ReplaceVMExecutor(vmExecutor executor.Executor) {
	context.vmExecutor = vmExecutor
}

// GetInstanceTracker returns the internal instance tracker
func (context *runtimeContext) GetInstanceTracker() vmhost.InstanceTracker {
	return context.iTracker
}

// StartWasmerInstance creates a new wasmer instance if the maxInstanceStackSize has not been reached.
func (context *runtimeContext) StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error {
	context.iTracker.UnsetInstance()

	if context.GetInstanceStackSize() >= context.maxInstanceStackSize {
		logRuntime.Trace("create instance", "error", vmhost.ErrMaxInstancesReached)
		return vmhost.ErrMaxInstancesReached
	}

	var codeHash []byte
	if newCode {
		codeHash = context.hasher.Compute(string(contract))
	} else {
		blockchain := context.host.Blockchain()
		codeHash = blockchain.GetCodeHash(context.codeAddress)
	}

	context.iTracker.SetCodeHash(codeHash)

	defer func() {
		context.iTracker.LogCounts()
		logRuntime.Trace("code was new", "new", newCode)
	}()

	warmInstanceUsed := context.useWarmInstanceIfExists(gasLimit, newCode)
	if warmInstanceUsed {
		return nil
	}

	compiledCodeUsed := context.makeInstanceFromCompiledCode(gasLimit, newCode)
	if compiledCodeUsed {
		return nil
	}

	return context.makeInstanceFromContractByteCode(contract, gasLimit, newCode)
}

func (context *runtimeContext) makeInstanceFromCompiledCode(gasLimit uint64, newCode bool) bool {
	codeHash := context.iTracker.CodeHash()
	if newCode || len(codeHash) == 0 {
		return false
	}

	blockchain := context.host.Blockchain()
	found, compiledCode := blockchain.GetCompiledCode(codeHash)
	if !found {
		logRuntime.Trace("instance creation", "code", "cached compilation", "error", "compiled code was not found")
		return false
	}

	gasSchedule := context.host.Metering().GasSchedule()
	options := executor.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		MaxMemoryGrow:      uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrow),
		MaxMemoryGrowDelta: uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrowDelta),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := context.vmExecutor.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err != nil {
		logRuntime.Error("instance creation", "from", "cached compilation", "error", err)
		return false
	}

	context.iTracker.SetNewInstance(newInstance, Precompiled)
	context.verifyCode = false

	context.saveWarmInstance()
	logRuntime.Trace("start instance", "from", "cached compilation",
		"id", context.iTracker.Instance().ID(),
		"codeHash", context.iTracker.codeHash,
	)
	return true
}

func (context *runtimeContext) makeInstanceFromContractByteCode(contract []byte, gasLimit uint64, newCode bool) error {
	gasSchedule := context.host.Metering().GasSchedule()
	options := executor.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		MaxMemoryGrow:      uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrow),
		MaxMemoryGrowDelta: uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrowDelta),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := context.vmExecutor.NewInstanceWithOptions(contract, options)
	if err != nil {
		context.iTracker.UnsetInstance()
		logRuntime.Trace("instance creation", "from", "bytecode", "error", err)
		return err
	}

	context.iTracker.SetNewInstance(newInstance, Bytecode)

	if newCode || len(context.iTracker.CodeHash()) == 0 {
		codeHash := context.hasher.Compute(string(contract))
		context.iTracker.SetCodeHash(codeHash)
	}

	if newCode {
		err = context.VerifyContractCode()
		if err != nil {
			context.iTracker.ForceCleanInstance(true)
			logRuntime.Trace("instance creation", "from", "bytecode", "error", err)
			return err
		}
	}

	logRuntime.Trace("start instance",
		"from", "bytecode",
		"id", context.iTracker.Instance().ID(),
		"codeHash", context.iTracker.CodeHash(),
	)
	context.saveCompiledCode()

	return nil
}

func (context *runtimeContext) useWarmInstanceIfExists(gasLimit uint64, newCode bool) bool {
	if !WarmInstancesEnabled {
		return false
	}

	codeHash := context.iTracker.CodeHash()
	if newCode || len(codeHash) == 0 {
		return false
	}

	if context.isContractOrCodeHashOnTheStack() {
		return false
	}

	ok := context.iTracker.UseWarmInstance(codeHash, newCode)
	if !ok {
		return false
	}

	context.SetPointsUsed(0)
	context.iTracker.Instance().SetGasLimit(gasLimit)
	context.SetRuntimeBreakpointValue(vmhost.BreakpointNone)
	context.verifyCode = false
	logRuntime.Trace("start instance", "from", "warm", "id", context.iTracker.Instance().ID())
	return true
}

// GetSCCode returns the SC code of the current SC.
func (context *runtimeContext) GetSCCode() ([]byte, error) {
	blockchain := context.host.Blockchain()

	code, err := blockchain.GetCode(context.codeAddress)
	if err != nil {
		return nil, err
	}

	context.codeSize = uint64(len(code))
	return code, nil
}

// GetSCCodeSize returns the cached size of the current SC code.
func (context *runtimeContext) GetSCCodeSize() uint64 {
	return context.codeSize
}

func (context *runtimeContext) saveCompiledCode() {
	compiledCode, err := context.iTracker.Instance().Cache()
	if err != nil {
		logRuntime.Error("getCompiledCode from instance", "error", err)
		return
	}

	codeHash := context.iTracker.CodeHash()
	blockchain := context.host.Blockchain()
	blockchain.SaveCompiledCode(codeHash, compiledCode)
	logRuntime.Trace("save compiled code", "codeHash", codeHash)

	found, _ := blockchain.GetCompiledCode(codeHash)
	if !found {
		logRuntime.Trace("save compiled code silent fail, code hash not found")
	}

	context.saveWarmInstance()
}

func (context *runtimeContext) saveWarmInstance() {
	if !WarmInstancesEnabled {
		return
	}

	codeHash := context.iTracker.CodeHash()
	if context.iTracker.IsCodeHashOnTheStack(codeHash) {
		return
	}

	context.iTracker.SaveAsWarmInstance()
}

// MustVerifyNextContractCode sets the verifyCode field to true
func (context *runtimeContext) MustVerifyNextContractCode() {
	context.verifyCode = true
}

// SetMaxInstanceStackSize sets the maximum number of allowed Wasmer instances on
// the instance stack, for recursivity.
func (context *runtimeContext) SetMaxInstanceStackSize(maxInstances uint64) {
	context.maxInstanceStackSize = maxInstances
}

// InitStateFromContractCallInput initializes the state of the runtime context
// (and the async context) from the provided ContractCallInput.
func (context *runtimeContext) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	context.SetVMInput(input)
	context.codeAddress = input.RecipientAddr
	context.callFunction = input.Function

	logRuntime.Trace("init state from call input",
		"caller", input.CallerAddr,
		"contract", input.RecipientAddr,
		"func", input.Function,
		"args", input.Arguments,
		"gas provided", input.GasProvided)
}

// SetCustomCallFunction sets a custom function to be called next, instead of
// the one specified by the current ContractCallInput.
func (context *runtimeContext) SetCustomCallFunction(callFunction string) {
	context.callFunction = callFunction
	logRuntime.Trace("set custom call function", "function", callFunction)
}

// PushState appends the current runtime state to the state stack; this
// includes the currently running Wasmer instance.
func (context *runtimeContext) PushState() {
	newState := &runtimeContext{
		codeAddress:  context.codeAddress,
		callFunction: context.callFunction,
		readOnly:     context.readOnly,
	}
	newState.SetVMInput(context.vmInput)

	context.stateStack = append(context.stateStack, newState)

	// Also preserve the currently running Wasmer instance at the top of the
	// instance stack; when the corresponding call to popInstance() is made, a
	// check is made to ensure that the running instance will not be cleaned
	// while still required for execution.
	context.pushInstance()
}

// PopSetActiveState pops the state at the top of the state stack and sets it as the current state.
func (context *runtimeContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.popInstance()

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.SetVMInput(prevState.vmInput)
	context.codeAddress = prevState.codeAddress
	context.callFunction = prevState.callFunction
	context.readOnly = prevState.readOnly
}

// PopDiscard removes the latest entry from the state stack
func (context *runtimeContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.iTracker.PopSetActiveState()

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack discards the entire state state stack and initializes it anew.
func (context *runtimeContext) ClearStateStack() {
	context.stateStack = make([]*runtimeContext, 0)
	context.iTracker.ClearStateStack()
}

// pushInstance pushes the current Wasmer instance on the instance stack (separate from the state stack).
func (context *runtimeContext) pushInstance() {
	context.iTracker.PushState()
}

// popInstance removes the latest entry from the wasmer instance stack and sets it
// as the current wasmer instance
func (context *runtimeContext) popInstance() {
	context.iTracker.PopSetActiveState()
}

// GetInstanceStackSize returns the number of the currently running Wasmer instances.
func (context *runtimeContext) GetInstanceStackSize() uint64 {
	return context.iTracker.StackSize()
}

// GetVMType returns the vm type for the current context.
func (context *runtimeContext) GetVMType() []byte {
	return context.vmType
}

// GetVMInput returns the vm input for the current context.
func (context *runtimeContext) GetVMInput() *vmcommon.ContractCallInput {
	return context.vmInput
}

func copyESDTTransfer(esdtTransfer *vmcommon.ESDTTransfer) *vmcommon.ESDTTransfer {
	newESDTTransfer := &vmcommon.ESDTTransfer{
		ESDTValue:      big.NewInt(0).Set(esdtTransfer.ESDTValue),
		ESDTTokenType:  esdtTransfer.ESDTTokenType,
		ESDTTokenNonce: esdtTransfer.ESDTTokenNonce,
		ESDTTokenName:  make([]byte, len(esdtTransfer.ESDTTokenName)),
	}
	copy(newESDTTransfer.ESDTTokenName, esdtTransfer.ESDTTokenName)
	return newESDTTransfer
}

// SetVMInput sets the given vm input as the current context vm input.
func (context *runtimeContext) SetVMInput(vmInput *vmcommon.ContractCallInput) {
	if vmInput == nil {
		context.vmInput = vmInput
		return
	}

	internalVMInput := vmcommon.VMInput{
		CallType:             vmInput.CallType,
		GasPrice:             vmInput.GasPrice,
		GasProvided:          vmInput.GasProvided,
		GasLocked:            vmInput.GasLocked,
		CallValue:            big.NewInt(0),
		ReturnCallAfterError: vmInput.ReturnCallAfterError,
	}
	context.vmInput = &vmcommon.ContractCallInput{
		VMInput:       internalVMInput,
		RecipientAddr: vmInput.RecipientAddr,
		Function:      vmInput.Function,
	}

	if vmInput.CallValue != nil {
		context.vmInput.CallValue.Set(vmInput.CallValue)
	}

	if len(vmInput.CallerAddr) > 0 {
		context.vmInput.CallerAddr = make([]byte, len(vmInput.CallerAddr))
		copy(context.vmInput.CallerAddr, vmInput.CallerAddr)
	}

	context.vmInput.ESDTTransfers = make([]*vmcommon.ESDTTransfer, len(vmInput.ESDTTransfers))

	if len(vmInput.ESDTTransfers) > 0 {
		for i, esdtTransfer := range vmInput.ESDTTransfers {
			context.vmInput.ESDTTransfers[i] = copyESDTTransfer(esdtTransfer)
		}
	}

	if len(vmInput.OriginalTxHash) > 0 {
		context.vmInput.OriginalTxHash = make([]byte, len(vmInput.OriginalTxHash))
		copy(context.vmInput.OriginalTxHash, vmInput.OriginalTxHash)
	}

	if len(vmInput.CurrentTxHash) > 0 {
		context.vmInput.CurrentTxHash = make([]byte, len(vmInput.CurrentTxHash))
		copy(context.vmInput.CurrentTxHash, vmInput.CurrentTxHash)
	}

	if len(vmInput.PrevTxHash) > 0 {
		context.vmInput.PrevTxHash = make([]byte, len(vmInput.PrevTxHash))
		copy(context.vmInput.PrevTxHash, vmInput.PrevTxHash)
	}

	if len(vmInput.Arguments) > 0 {
		context.vmInput.Arguments = make([][]byte, len(vmInput.Arguments))
		for i, arg := range vmInput.Arguments {
			context.vmInput.Arguments[i] = make([]byte, len(arg))
			copy(context.vmInput.Arguments[i], arg)
		}
	}
}

// GetContextAddress returns the SC address from the current context.
func (context *runtimeContext) GetContextAddress() []byte {
	return context.vmInput.RecipientAddr
}

// SetCodeAddress sets the given address as the scAddress for the current context.
func (context *runtimeContext) SetCodeAddress(scAddress []byte) {
	context.codeAddress = scAddress
}

// GetCurrentTxHash returns the hash of the current transaction, as specified by the current VMInput.
func (context *runtimeContext) GetCurrentTxHash() []byte {
	return context.vmInput.CurrentTxHash
}

// GetOriginalTxHash returns the hash of the original transaction, in the case of async calls, as specified by the current VMInput.
func (context *runtimeContext) GetOriginalTxHash() []byte {
	return context.vmInput.OriginalTxHash
}

// GetPrevTxHash returns the hash of the previous transaction, in the case of async calls, as specified by the current VMInput.
func (context *runtimeContext) GetPrevTxHash() []byte {
	return context.vmInput.PrevTxHash
}

// FunctionName returns the name of the contract function to be called next
func (context *runtimeContext) FunctionName() string {
	return context.callFunction
}

// Arguments returns the binary arguments that will be passed to the contract to be executed, as specified by the current VMInput.
func (context *runtimeContext) Arguments() [][]byte {
	return context.vmInput.Arguments
}

// ExtractCodeUpgradeFromArgs extracts the code and code metadata from the
// current VMInput.Arguments, assuming a contract code upgrade has been requested.
func (context *runtimeContext) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	const numMinUpgradeArguments = 2

	arguments := context.vmInput.Arguments
	if len(arguments) < numMinUpgradeArguments {
		return nil, nil, vmhost.ErrInvalidUpgradeArguments
	}

	code := arguments[0]
	codeMetadata := arguments[1]
	context.vmInput.Arguments = context.vmInput.Arguments[numMinUpgradeArguments:]
	return code, codeMetadata, nil
}

// FailExecution informs Wasmer to immediately stop the execution of the contract
// with BreakpointExecutionFailed and sets the corresponding VMOutput fields accordingly
// FailExecution sets the returnMessage, returnCode and runtimeBreakpoint according to the given error.
func (context *runtimeContext) FailExecution(err error) {
	context.host.Output().SetReturnCode(vmcommon.ExecutionFailed)

	var message string
	breakpoint := vmhost.BreakpointExecutionFailed

	if err != nil {
		message = err.Error()
		context.AddError(err)
		if errors.Is(err, vmhost.ErrNotEnoughGas) {
			breakpoint = vmhost.BreakpointOutOfGas
		}
	} else {
		message = "execution failed"
		context.AddError(errors.New(message))
	}

	context.host.Output().SetReturnMessage(message)
	if !check.IfNil(context.iTracker.Instance()) {
		context.SetRuntimeBreakpointValue(breakpoint)
	}

	traceMessage := message
	if err != nil {
		traceMessage = err.Error()
	}
	logRuntime.Trace("execution failed", "message", traceMessage)
}

// SignalUserError informs Wasmer to immediately stop the execution of the contract
// with BreakpointSignalError and sets the corresponding VMOutput fields accordingly
func (context *runtimeContext) SignalUserError(message string) {
	context.host.Output().SetReturnCode(vmcommon.UserError)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(vmhost.BreakpointSignalError)
	context.AddError(errors.New(message))
	logRuntime.Trace("user error signalled", "message", message)
}

// SetRuntimeBreakpointValue sets the specified runtime breakpoint in Wasmer,
// immediately stopping the contract execution.
func (context *runtimeContext) SetRuntimeBreakpointValue(value vmhost.BreakpointValue) {
	context.iTracker.Instance().SetBreakpointValue(uint64(value))
	logRuntime.Trace("runtime breakpoint set", "breakpoint", value)
}

// GetRuntimeBreakpointValue retrieves the value of the breakpoint that has
// stopped the execution of the contract.
func (context *runtimeContext) GetRuntimeBreakpointValue() vmhost.BreakpointValue {
	return vmhost.BreakpointValue(context.iTracker.Instance().GetBreakpointValue())
}

// VerifyContractCode performs validation on the WASM bytecode (declaration of memory and legal functions).
func (context *runtimeContext) VerifyContractCode() error {
	if !context.verifyCode {
		return nil
	}

	context.verifyCode = false

	err := context.validator.verifyMemoryDeclaration(context.iTracker.Instance())
	if err != nil {
		logRuntime.Trace("verify contract code", "error", err)
		return err
	}

	err = context.validator.verifyFunctions(context.iTracker.Instance())
	if err != nil {
		logRuntime.Trace("verify contract code", "error", err)
		return err
	}

	err = context.validator.verifyProtectedFunctions(context.iTracker.Instance())
	if err != nil {
		logRuntime.Trace("verify contract code", "error", err)
		return err
	}

	enableEpochsHandler := context.host.EnableEpochsHandler()
	if enableEpochsHandler.IsManagedCryptoAPIsFlagEnabled() {
		err = context.validator.verifyProtectedFunctions(context.iTracker.Instance())
		if err != nil {
			logRuntime.Trace("verify contract code", "error", err)
			return err
		}
	}

	logRuntime.Trace("verified contract code")

	return nil
}

// BaseOpsErrorShouldFailExecution returns true
func (context *runtimeContext) BaseOpsErrorShouldFailExecution() bool {
	return true
}

// SyncExecAPIErrorShouldFailExecution specifies whether an error in the
// EEI functions for synchronous execution should abort contract execution.
func (context *runtimeContext) SyncExecAPIErrorShouldFailExecution() bool {
	return true
}

// BigIntAPIErrorShouldFailExecution specifies whether an error in the EEI
// functions for BigInt operations should abort contract execution.
func (context *runtimeContext) BigIntAPIErrorShouldFailExecution() bool {
	return true
}

// BigFloatAPIErrorShouldFailExecution returns true
func (context *runtimeContext) BigFloatAPIErrorShouldFailExecution() bool {
	return true
}

// CryptoAPIErrorShouldFailExecution specifies whether an error in the EEI
// functions for crypto operations should abort contract execution.
func (context *runtimeContext) CryptoAPIErrorShouldFailExecution() bool {
	return true
}

// ManagedBufferAPIErrorShouldFailExecution returns true
func (context *runtimeContext) ManagedBufferAPIErrorShouldFailExecution() bool {
	return true
}

// GetPointsUsed returns the gas amount spent by the currently running Wasmer instance.
func (context *runtimeContext) GetPointsUsed() uint64 {
	if check.IfNil(context.iTracker.Instance()) {
		return 0
	}
	return context.iTracker.Instance().GetPointsUsed()
}

// SetPointsUsed directly sets the gas amount already spent by the currently running Wasmer instance.
func (context *runtimeContext) SetPointsUsed(gasPoints uint64) {
	if gasPoints > builtinMath.MaxInt64 {
		gasPoints = builtinMath.MaxInt64
	}
	context.iTracker.Instance().SetPointsUsed(gasPoints)
}

// ReadOnly verifies whether the read-only execution flag is set.
func (context *runtimeContext) ReadOnly() bool {
	return context.readOnly
}

// SetReadOnly sets the read-only execution flag.
func (context *runtimeContext) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

// GetInstance returns the current wasmer instance
func (context *runtimeContext) GetInstance() executor.Instance {
	return context.iTracker.Instance()
}

// CleanInstance cleans the current instance
func (context *runtimeContext) CleanInstance() {
	context.iTracker.ForceCleanInstance(false)
}

// isContractOrCodeHashOnTheStack iterates over the state stack to find whether the
// provided SC address is already in execution, below the current instance.
func (context *runtimeContext) isContractOrCodeHashOnTheStack() bool {
	if context.isScAddressOnTheStack(context.codeAddress) {
		return true
	}
	return context.iTracker.IsCodeHashOnTheStack(context.iTracker.CodeHash())
}

func (context *runtimeContext) isScAddressOnTheStack(scAddress []byte) bool {
	for _, state := range context.stateStack {
		if bytes.Equal(scAddress, state.codeAddress) {
			return true
		}
	}
	return false
}

// CountSameContractInstancesOnStack returns the number of times the given contract
// address appears in the state stack.
func (context *runtimeContext) CountSameContractInstancesOnStack(address []byte) uint64 {
	count := uint64(0)
	for _, state := range context.stateStack {
		if bytes.Equal(address, state.vmInput.RecipientAddr) {
			count++
		}
	}

	return count
}

// FunctionNameChecked returns the function name, after checking that it exists in the contract.
func (context *runtimeContext) FunctionNameChecked() (string, error) {
	functionName := context.FunctionName()
	if context.HasFunction(functionName) {
		return functionName, nil
	}

	// If the requested function is missing from the contract exports, but is
	// named like vmhost.CallbackFunctionName, then a different error is returned
	// to indicate that, not just a missing function.
	if context.callFunction == vmhost.CallbackFunctionName {
		logRuntime.Trace("missing function " + vmhost.CallbackFunctionName)
		return "", vmhost.ErrNilCallbackFunction
	}

	return "", executor.ErrFuncNotFound
}

// CallSCFunction will execute the function with given name from the loaded contract.
func (context *runtimeContext) CallSCFunction(functionName string) error {
	return context.iTracker.Instance().CallFunction(functionName)
}

// IsFunctionImported returns true if the WASM module imports the specified function.
func (context *runtimeContext) IsFunctionImported(name string) bool {
	return context.iTracker.Instance().IsFunctionImported(name)
}

// AddError adds an error to the global error list on runtime context
func (context *runtimeContext) AddError(err error, otherInfo ...string) {
	if err == nil {
		return
	}
	if context.errors == nil {
		context.errors = vmhost.WrapError(err, otherInfo...)
		return
	}
	context.errors = context.errors.WrapWithError(err, otherInfo...)
}

// GetAllErrors returns all the errors stored on the RuntimeContext
func (context *runtimeContext) GetAllErrors() error {
	return context.errors
}

// EndExecution performs final steps after execution ends
func (context *runtimeContext) EndExecution() {
	context.iTracker.UnsetInstance()
}

// ValidateInstances checks the state of the instances after execution
func (context *runtimeContext) ValidateInstances() error {
	if !WarmInstancesEnabled {
		return nil
	}

	err := context.iTracker.CheckInstances()
	if err != nil {
		return err
	}

	err = context.checkNumRunningInstances()
	if err != nil {
		return err
	}

	return nil
}

func (context *runtimeContext) checkNumRunningInstances() error {
	_, cold := context.iTracker.NumRunningInstances()
	if cold > 0 {
		return fmt.Errorf("potentially leaked cold instances")
	}

	return nil
}

// ValidateCallbackName verifies whether the provided function name may be used as AsyncCall callback
func (context *runtimeContext) ValidateCallbackName(callbackName string) error {
	err := context.validator.verifyValidFunctionName(callbackName)
	if err != nil {
		return vmhost.ErrInvalidFunctionName
	}
	if callbackName == vmhost.InitFunctionName {
		return vmhost.ErrInvalidFunctionName
	}
	if context.host.IsBuiltinFunctionName(callbackName) {
		return vmhost.ErrCannotUseBuiltinAsCallback
	}
	if !context.HasFunction(callbackName) {
		return executor.ErrFuncNotFound
	}

	return nil
}

// HasFunction checks if loaded contract has a function (endpoint) with given name.
func (context *runtimeContext) HasFunction(functionName string) bool {
	return context.iTracker.Instance().HasFunction(functionName)
}

// EpochConfirmed is called whenever a new epoch is confirmed
func (context *runtimeContext) EpochConfirmed(_ uint32, _ uint64) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (context *runtimeContext) IsInterfaceNil() bool {
	return context == nil
}
