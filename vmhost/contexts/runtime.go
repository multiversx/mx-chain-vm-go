package contexts

import (
	"bytes"
	"errors"
	"fmt"
	builtinMath "math"
	"math/big"
	"unsafe"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/math"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
)

var logRuntime = logger.GetOrCreate("vm/runtime")

var _ vmhost.RuntimeContext = (*runtimeContext)(nil)

const warmCacheSize = 100

// WarmInstancesEnabled controls the usage of warm instances
const WarmInstancesEnabled = true

// HashComputer provides hash computation
type HashComputer interface {
	Compute(string) []byte
	Size() int
	IsInterfaceNil() bool
}

type runtimeContext struct {
	host               vmhost.VMHost
	vmInput            *vmcommon.ContractCallInput
	codeAddress        []byte
	codeSize           uint64
	callFunction       string
	vmType             []byte
	readOnly           bool
	verifyCode         bool
	maxWasmerInstances uint64

	iTracker *instanceTracker

	stateStack []*runtimeContext

	asyncCallInfo    *vmhost.AsyncCallInfo
	asyncContextInfo *vmhost.AsyncContextInfo

	validator       *wasmValidator
	instanceBuilder vmhost.InstanceBuilder
	hasher          vmhost.HashComputer

	errors vmhost.WrappableError
}

// NewRuntimeContext creates a new runtimeContext
func NewRuntimeContext(
	host vmhost.VMHost,
	vmType []byte,
	builtInFuncContainer vmcommon.BuiltInFunctionContainer,
	hasher vmhost.HashComputer,
) (*runtimeContext, error) {

	if check.IfNil(host) {
		return nil, vmhost.ErrNilHost
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

	scAPINames := host.GetAPIMethods().Names()

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

	context.instanceBuilder = &WasmerInstanceBuilder{}
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
	context.asyncCallInfo = nil
	context.asyncContextInfo = &vmhost.AsyncContextInfo{
		AsyncContextMap: make(map[string]*vmhost.AsyncContext),
	}
	context.iTracker.InitState()
	context.errors = nil

	logRuntime.Trace("init state")
}

// ClearWarmInstanceCache clears all elements from warm instance cache
func (context *runtimeContext) ClearWarmInstanceCache() {
	context.iTracker.ClearWarmInstanceCache()
	context.iTracker.UnsetInstance()
}

// ReplaceInstanceBuilder replaces the instance builder, allowing the creation
// of mocked Wasmer instances; this is used for tests only
func (context *runtimeContext) ReplaceInstanceBuilder(builder vmhost.InstanceBuilder) {
	context.instanceBuilder = builder
}

// StartWasmerInstance creates a new wasmer instance if the maxWasmerInstances has not been reached.
func (context *runtimeContext) StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error {
	context.iTracker.UnsetInstance()

	if context.RunningInstancesCount() >= context.maxWasmerInstances {
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

	context.iTracker.SetCodeSize(uint64(len(contract)))
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
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		MaxMemoryGrow:      uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrow),
		MaxMemoryGrowDelta: uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrowDelta),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := context.instanceBuilder.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err != nil {
		logRuntime.Error("instance creation", "code", "cached compilation", "error", err)
		return false
	}

	context.iTracker.SetNewInstance(newInstance, Precompiled)

	hostReference := uintptr(unsafe.Pointer(&context.host))
	context.iTracker.Instance().SetContextData(hostReference)
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
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		MaxMemoryGrow:      uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrow),
		MaxMemoryGrowDelta: uint64(gasSchedule.WASMOpcodeCost.MaxMemoryGrowDelta),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := context.instanceBuilder.NewInstanceWithOptions(contract, options)
	if err != nil {
		context.iTracker.UnsetInstance()
		logRuntime.Trace("instance creation", "code", "bytecode", "error", err)
		return err
	}

	context.iTracker.SetNewInstance(newInstance, Bytecode)

	if newCode || len(context.iTracker.CodeHash()) == 0 {
		codeHash := context.hasher.Compute(string(contract))
		context.iTracker.SetCodeHash(codeHash)
	}

	hostReference := uintptr(unsafe.Pointer(&context.host))
	context.iTracker.Instance().SetContextData(hostReference)

	if newCode {
		err = context.VerifyContractCode()
		if err != nil {
			context.iTracker.ForceCleanInstance(true)
			logRuntime.Trace("instance creation", "code", "bytecode", "error", err)
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

	hostReference := uintptr(unsafe.Pointer(&context.host))
	context.iTracker.Instance().SetContextData(hostReference)
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
	context.iTracker.SetCodeSize(uint64(len(code)))
	return code, nil
}

// GetSCCodeSize returns the size of the current SC code.
func (context *runtimeContext) GetSCCodeSize() uint64 {
	if context.host.EnableEpochsHandler().IsRuntimeCodeSizeFixEnabled() {
		return context.iTracker.GetCodeSize()
	}
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

// SetMaxInstanceCount sets the maxWasmerInstances field to the given value
func (context *runtimeContext) SetMaxInstanceCount(maxInstances uint64) {
	context.maxWasmerInstances = maxInstances
}

// InitStateFromContractCallInput initializes the runtime context state with the values from the given input
func (context *runtimeContext) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	context.SetVMInput(input)
	context.codeAddress = input.RecipientAddr
	context.callFunction = input.Function
	// Reset async map for initial state
	context.asyncContextInfo = &vmhost.AsyncContextInfo{
		CallerAddr:      input.CallerAddr,
		AsyncContextMap: make(map[string]*vmhost.AsyncContext),
	}

	logRuntime.Trace("init state from call input",
		"caller", input.CallerAddr,
		"contract", input.RecipientAddr,
		"func", input.Function,
		"args", input.Arguments)
}

// SetCustomCallFunction sets the given string as the callFunction field.
func (context *runtimeContext) SetCustomCallFunction(callFunction string) {
	context.callFunction = callFunction
	logRuntime.Trace("set custom call function", "function", callFunction)
}

// PushState appends the current runtime state to the state stack; this
// includes the currently running Wasmer instance.
func (context *runtimeContext) PushState() {
	newState := &runtimeContext{
		codeAddress:      context.codeAddress,
		callFunction:     context.callFunction,
		readOnly:         context.readOnly,
		asyncCallInfo:    context.asyncCallInfo,
		asyncContextInfo: context.asyncContextInfo,
	}
	newState.SetVMInput(context.vmInput)

	context.stateStack = append(context.stateStack, newState)

	// Also preserve the currently running Wasmer instance at the top of the
	// instance stack; when the corresponding call to popInstance() is made, a
	// check is made to ensure that the running instance will not be cleaned
	// while still required for execution.
	context.pushInstance()
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current
// runtime context state.
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
	context.asyncCallInfo = prevState.asyncCallInfo
	context.asyncContextInfo = prevState.asyncContextInfo
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

// ClearStateStack reinitializes the state stack.
func (context *runtimeContext) ClearStateStack() {
	context.stateStack = make([]*runtimeContext, 0)
	context.iTracker.ClearStateStack()
}

// pushInstance appends the current wasmer instance to the instance stack.
func (context *runtimeContext) pushInstance() {
	context.iTracker.PushState()
}

// popInstance removes the latest entry from the wasmer instance stack and sets it
// as the current wasmer instance
func (context *runtimeContext) popInstance() {
	context.iTracker.PopSetActiveState()
}

// RunningInstancesCount returns the length of the instance stack.
func (context *runtimeContext) RunningInstancesCount() uint64 {
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

// GetCurrentTxHash returns the txHash from the vmInput of the current context.
func (context *runtimeContext) GetCurrentTxHash() []byte {
	return context.vmInput.CurrentTxHash
}

// GetOriginalTxHash returns the originalTxHash from the vmInput of the current context.
func (context *runtimeContext) GetOriginalTxHash() []byte {
	return context.vmInput.OriginalTxHash
}

// Function returns the callFunction for the current context.
func (context *runtimeContext) Function() string {
	return context.callFunction
}

// Arguments returns the arguments from the vmInput of the current context.
func (context *runtimeContext) Arguments() [][]byte {
	return context.vmInput.Arguments
}

// ExtractCodeUpgradeFromArgs extracts the arguments needed for a code upgrade from the vmInput.
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

// FailExecution sets the returnMessage, returnCode and runtimeBreakpoint according to the given error.
func (context *runtimeContext) FailExecution(err error) {
	context.host.Output().SetReturnCode(vmcommon.ExecutionFailed)

	var message string
	breakpoint := vmhost.BreakpointExecutionFailed

	if err != nil {
		message = err.Error()
		context.AddError(err)
		if errors.Is(err, vmhost.ErrNotEnoughGas) && context.host.FixOOGReturnCodeEnabled() {
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

// SignalUserError sets the returnMessage, returnCode and runtimeBreakpoint according an user error.
func (context *runtimeContext) SignalUserError(message string) {
	context.host.Output().SetReturnCode(vmcommon.UserError)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(vmhost.BreakpointSignalError)
	context.AddError(errors.New(message))
	logRuntime.Trace("user error signalled", "message", message)
}

// SetRuntimeBreakpointValue sets the given value as a breakpoint value.
func (context *runtimeContext) SetRuntimeBreakpointValue(value vmhost.BreakpointValue) {
	context.iTracker.Instance().SetBreakpointValue(uint64(value))
	logRuntime.Trace("runtime breakpoint set", "breakpoint", value)
}

// GetRuntimeBreakpointValue returns the breakpoint value for the current wasmer instance.
func (context *runtimeContext) GetRuntimeBreakpointValue() vmhost.BreakpointValue {
	return vmhost.BreakpointValue(context.iTracker.Instance().GetBreakpointValue())
}

// VerifyContractCode checks the current wasmer instance for enough memory and for correct functions.
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

	enableEpochsHandler := context.host.EnableEpochsHandler()
	if !enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled() {
		err = context.checkBackwardCompatibility()
		if err != nil {
			logRuntime.Trace("verify contract code", "error", err)
			return err
		}
	}

	if !enableEpochsHandler.IsManagedCryptoAPIsFlagEnabled() {
		err = context.checkIfContainsNewManagedCryptoAPI()
		if err != nil {
			logRuntime.Trace("verify contract code", "error", err)
			return err
		}
	}

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

func (context *runtimeContext) checkBackwardCompatibility() error {
	if context.iTracker.Instance().IsFunctionImported("mBufferSetByteSlice") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("getESDTLocalRoles") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("validateTokenIdentifier") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedSha256") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedKeccak256") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("mBufferStorageLoadFromAddress") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("cleanReturnData") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("deleteFromReturnData") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("completedTxEvent") {
		return vmhost.ErrContractInvalid
	}

	return nil
}

func (context *runtimeContext) checkIfContainsNewManagedCryptoAPI() error {
	if context.iTracker.Instance().IsFunctionImported("managedIsESDTFrozen") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedIsESDTPaused") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedIsESDTLimitedTransfer") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedBufferToHex") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigIntToString") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedRipemd160") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedVerifyBLS") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedVerifyEd25519") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedVerifySecp256k1") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedVerifyCustomSecp256k1") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedEncodeSecp256k1DerSignature") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedScalarBaseMultEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedScalarMultEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedMarshalEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedUnmarshalEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedMarshalCompressedEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedUnmarshalCompressedEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedGenerateKeyEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("managedCreateEC") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("mBufferToBigFloat") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("mBufferFromBigFloat") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatNewFromParts") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatNewFromFrac") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatNewFromSci") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatAdd") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatSub") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatMul") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatDiv") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatAbs") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatCmp") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatSign") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatClone") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatSqrt") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatPow") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatFloor") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatCeil") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatTruncate") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatIsInt") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatSetInt64") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatSetBigInt") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatGetConstPi") {
		return vmhost.ErrContractInvalid
	}
	if context.iTracker.Instance().IsFunctionImported("bigFloatGetConstE") {
		return vmhost.ErrContractInvalid
	}

	return nil
}

// BaseOpsErrorShouldFailExecution returns true
func (context *runtimeContext) BaseOpsErrorShouldFailExecution() bool {
	return true
}

// SyncExecAPIErrorShouldFailExecution returns true
func (context *runtimeContext) SyncExecAPIErrorShouldFailExecution() bool {
	return true
}

// BigIntAPIErrorShouldFailExecution returns true
func (context *runtimeContext) BigIntAPIErrorShouldFailExecution() bool {
	return true
}

// BigFloatAPIErrorShouldFailExecution returns true
func (context *runtimeContext) BigFloatAPIErrorShouldFailExecution() bool {
	return true
}

// CryptoAPIErrorShouldFailExecution returns true
func (context *runtimeContext) CryptoAPIErrorShouldFailExecution() bool {
	return true
}

// ManagedBufferAPIErrorShouldFailExecution returns true
func (context *runtimeContext) ManagedBufferAPIErrorShouldFailExecution() bool {
	return true
}

// GetPointsUsed returns the gas points used by the current wasmer instance.
func (context *runtimeContext) GetPointsUsed() uint64 {
	if context.iTracker.Instance() == nil {
		return 0
	}
	return context.iTracker.Instance().GetPointsUsed()
}

// SetPointsUsed sets the given gas points as the gas points used by the current wasmer instance.
func (context *runtimeContext) SetPointsUsed(gasPoints uint64) {
	if gasPoints > builtinMath.MaxInt64 {
		gasPoints = builtinMath.MaxInt64
	}
	context.iTracker.Instance().SetPointsUsed(gasPoints)
}

// ReadOnly returns true if the current context is readOnly
func (context *runtimeContext) ReadOnly() bool {
	return context.readOnly
}

// SetReadOnly sets the readOnly field of the context to the given value.
func (context *runtimeContext) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

// GetInstance returns the current wasmer instance
func (context *runtimeContext) GetInstance() wasmer.InstanceHandler {
	return context.iTracker.Instance()
}

// GetWarmInstance retrieves an instance from the warm cache
func (context *runtimeContext) GetWarmInstance(codeHash []byte) (wasmer.InstanceHandler, bool) {
	return context.iTracker.GetWarmInstance(codeHash)
}

// GetInstanceExports returns the current wasmer instance exports.
func (context *runtimeContext) GetInstanceExports() wasmer.ExportsMap {
	return context.iTracker.Instance().GetExports()
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

// CallFunction calls the specified instance function
func (context *runtimeContext) CallFunction(funcName string) error {
	instance := context.iTracker.Instance()
	if !instance.HasFunction(funcName) {
		return vmhost.ErrFuncNotFound
	}

	_, err := instance.CallFunction(funcName)

	return err
}

// GetFunctionToCall returns the function to call from the wasmer instance exports.
func (context *runtimeContext) GetFunctionToCall() (string, error) {
	instance := context.iTracker.Instance()
	logRuntime.Trace("get function to call", "function", context.callFunction)
	if instance.HasFunction(context.callFunction) {
		return context.callFunction, nil
	}

	if context.callFunction == vmhost.CallbackFunctionName {
		// TODO rewrite this condition, until the AsyncContext is merged
		logRuntime.Trace("get function to call", "error", vmhost.ErrNilCallbackFunction)
		return "", vmhost.ErrNilCallbackFunction
	}

	return "", vmhost.ErrFuncNotFound
}

// GetInitFunction returns the init function from the current wasmer instance exports.
func (context *runtimeContext) GetInitFunction() string {
	instance := context.iTracker.Instance()
	if instance.HasFunction(vmhost.InitFunctionName) {
		return vmhost.InitFunctionName
	}

	return ""
}

// ExecuteAsyncCall locks the necessary gas and sets the async call info and a runtime breakpoint value.
func (context *runtimeContext) ExecuteAsyncCall(address []byte, data []byte, value []byte) error {
	if context.ReadOnly() && context.host.CheckExecuteReadOnly() {
		return vmhost.ErrInvalidCallOnReadOnlyMode
	}
	metering := context.host.Metering()
	err := metering.UseGasForAsyncStep()
	if err != nil {
		return err
	}

	gasToLock := uint64(0)
	if context.HasCallbackMethod() {
		gasToLock = metering.ComputeGasLockedForAsync()
		logRuntime.Trace("ExecuteAsyncCall", "gasToLock", gasToLock)

		err = metering.UseGasBounded(gasToLock)
		if err != nil {
			logRuntime.Trace("ExecuteAsyncCall: cannot lock gas", "err", err)
			return err
		}
	}

	context.SetAsyncCallInfo(&vmhost.AsyncCallInfo{
		Destination: address,
		Data:        data,
		GasLimit:    metering.GasLeft(),
		GasLocked:   gasToLock,
		ValueBytes:  value,
	})
	context.SetRuntimeBreakpointValue(vmhost.BreakpointAsyncCall)

	logRuntime.Trace("prepare async call",
		"caller", context.GetContextAddress(),
		"dest", address,
		"value", big.NewInt(0).SetBytes(value),
		"data", data)
	return nil
}

// SetAsyncCallInfo sets the given data as the async call info for the current context.
func (context *runtimeContext) SetAsyncCallInfo(asyncCallInfo *vmhost.AsyncCallInfo) {
	context.asyncCallInfo = asyncCallInfo
}

// AddAsyncContextCall adds the given async call to the asyncContextMap at the given identifier.
func (context *runtimeContext) AddAsyncContextCall(contextIdentifier []byte, asyncCall *vmhost.AsyncGeneratedCall) error {
	_, ok := context.asyncContextInfo.AsyncContextMap[string(contextIdentifier)]
	currentContextMap := context.asyncContextInfo.AsyncContextMap
	if !ok {
		currentContextMap[string(contextIdentifier)] = &vmhost.AsyncContext{
			AsyncCalls: make([]*vmhost.AsyncGeneratedCall, 0),
		}
	}

	currentContextMap[string(contextIdentifier)].AsyncCalls =
		append(currentContextMap[string(contextIdentifier)].AsyncCalls, asyncCall)

	return nil
}

// GetAsyncContextInfo returns the async context info for the current context.
func (context *runtimeContext) GetAsyncContextInfo() *vmhost.AsyncContextInfo {
	return context.asyncContextInfo
}

// GetAsyncContext returns the async context mapped to the given context identifier.
func (context *runtimeContext) GetAsyncContext(contextIdentifier []byte) (*vmhost.AsyncContext, error) {
	asyncContext, ok := context.asyncContextInfo.AsyncContextMap[string(contextIdentifier)]
	if !ok {
		return nil, vmhost.ErrAsyncContextDoesNotExist
	}

	return asyncContext, nil
}

// GetAsyncCallInfo returns the async call info for the current context.
func (context *runtimeContext) GetAsyncCallInfo() *vmhost.AsyncCallInfo {
	return context.asyncCallInfo
}

// HasCallbackMethod returns true if the current wasmer instance exports has a callback method.
func (context *runtimeContext) HasCallbackMethod() bool {
	_, ok := context.iTracker.Instance().GetExports()[vmhost.CallbackFunctionName]
	return ok
}

// IsFunctionImported returns true if the WASM module imports the specified function.
func (context *runtimeContext) IsFunctionImported(name string) bool {
	return context.iTracker.Instance().IsFunctionImported(name)
}

// MemLoad returns the contents from the given offset of the WASM memory.
func (context *runtimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memory := context.iTracker.Instance().GetInstanceCtxMemory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := math.AddInt32(offset, length)

	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > memoryLength
	isRequestedEndTooLarge := uint32(requestedEnd) > memoryLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("mem load: %w", vmhost.ErrBadBounds)
	}
	if isLengthNegative {
		return nil, fmt.Errorf("mem load: %w", vmhost.ErrNegativeLength)
	}

	result := make([]byte, length)
	if isRequestedEndTooLarge {
		copy(result, memoryView[offset:])
	} else {
		copy(result, memoryView[offset:requestedEnd])
	}

	return result, nil
}

// MemLoadMultiple returns multiple byte slices loaded from the WASM memory, starting at the given offset and having the provided lengths.
func (context *runtimeContext) MemLoadMultiple(offset int32, lengths []int32) ([][]byte, error) {
	if len(lengths) == 0 {
		return [][]byte{}, nil
	}

	results := make([][]byte, len(lengths))

	for i, length := range lengths {
		result, err := context.MemLoad(offset, length)
		if err != nil {
			return nil, err
		}

		results[i] = result
		offset += length
	}

	return results, nil
}

// MemStore stores the given data in the WASM memory at the given offset.
func (context *runtimeContext) MemStore(offset int32, data []byte) error {
	dataLength := int32(len(data))
	if dataLength == 0 {
		return nil
	}

	memory := context.iTracker.Instance().GetInstanceCtxMemory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := math.AddInt32(offset, dataLength)

	isOffsetTooSmall := offset < 0
	if isOffsetTooSmall {
		return vmhost.ErrBadLowerBounds
	}

	isNewPageNecessary := uint32(requestedEnd) > memoryLength
	epochsHandler := context.host.EnableEpochsHandler()

	if isNewPageNecessary {
		if epochsHandler.IsRuntimeMemStoreLimitEnabled() {
			return vmhost.ErrBadUpperBounds
		}

		err := memory.Grow(1)
		if err != nil {
			return err
		}

		memoryView = memory.Data()
		memoryLength = memory.Length()
	}

	isRequestedEndTooLarge := uint32(requestedEnd) > memoryLength
	if isRequestedEndTooLarge {
		return vmhost.ErrBadUpperBounds
	}

	copy(memoryView[offset:requestedEnd], data)
	return nil
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

// IsInterfaceNil returns true if there is no value under the interface
func (context *runtimeContext) IsInterfaceNil() bool {
	return context == nil
}
