package contexts

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var log = logger.GetOrCreate("arwen/runtime")

var _ arwen.RuntimeContext = (*runtimeContext)(nil)

type runtimeContext struct {
	host         arwen.VMHost
	instance     *wasmer.Instance
	vmInput      *vmcommon.VMInput
	scAddress    []byte
	codeSize     uint64
	callFunction string
	vmType       []byte
	readOnly     bool

	verifyCode bool

	stateStack    []*runtimeContext
	instanceStack []*wasmer.Instance

	maxWasmerInstances uint64

	defaultAsyncCall *arwen.AsyncCall
	asyncContext     *arwen.AsyncContext

	validator *WASMValidator

	useWarmInstance     bool
	warmInstanceAddress []byte
	warmInstance        *wasmer.Instance
}

// NewRuntimeContext creates a new runtimeContext
func NewRuntimeContext(host arwen.VMHost, vmType []byte, useWarmInstance bool) (*runtimeContext, error) {
	scAPINames := host.GetAPIMethods().Names()
	protocolBuiltinFunctions := host.GetProtocolBuiltinFunctions()

	context := &runtimeContext{
		host:                host,
		vmType:              vmType,
		stateStack:          make([]*runtimeContext, 0),
		instanceStack:       make([]*wasmer.Instance, 0),
		validator:           NewWASMValidator(scAPINames, protocolBuiltinFunctions),
		useWarmInstance:     useWarmInstance,
		warmInstanceAddress: nil,
		warmInstance:        nil,
	}

	context.InitState()

	return context, nil
}

func (context *runtimeContext) InitState() {
	context.vmInput = &vmcommon.VMInput{}
	context.scAddress = make([]byte, 0)
	context.callFunction = ""
	context.verifyCode = false
	context.readOnly = false
	context.defaultAsyncCall = nil
	context.asyncContext = arwen.NewAsyncContext()
}

func (context *runtimeContext) setWarmInstanceWhenNeeded(gasLimit uint64) bool {
	scAddress := context.GetSCAddress()
	useWarm := context.useWarmInstance && context.warmInstanceAddress != nil && bytes.Equal(scAddress, context.warmInstanceAddress)
	if scAddress != nil && useWarm {
		log.Trace("Reusing the warm Wasmer instance")

		context.instance = context.warmInstance
		context.SetPointsUsed(0)
		context.instance.SetGasLimit(gasLimit)

		context.SetRuntimeBreakpointValue(arwen.BreakpointNone)
		return true
	}

	return false
}

func (context *runtimeContext) makeInstanceFromCompiledCode(codeHash []byte, gasLimit uint64, newCode bool) bool {
	if !context.host.IsAheadOfTimeCompileEnabled() {
		return false
	}

	if newCode || len(codeHash) == 0 {
		return false
	}

	blockchain := context.host.Blockchain()
	found, compiledCode := blockchain.GetCompiledCode(codeHash)
	if !found {
		log.Debug("compiled code was not found")
		return false
	}

	gasSchedule := context.host.Metering().GasSchedule()
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err != nil {
		log.Warn("NewInstanceFromCompiledCodeWithOptions", "error", err)
		return false
	}

	context.instance = newInstance

	idContext := arwen.AddHostContext(context.host)
	context.instance.SetContextData(idContext)
	context.verifyCode = false

	return true
}

func (context *runtimeContext) StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error {
	if context.RunningInstancesCount() >= context.maxWasmerInstances {
		context.instance = nil
		return arwen.ErrMaxInstancesReached
	}

	warmInstanceUsed := context.setWarmInstanceWhenNeeded(gasLimit)
	if warmInstanceUsed {
		return nil
	}

	blockchain := context.host.Blockchain()
	codeHash := blockchain.GetCodeHash(context.GetSCAddress())
	compiledCodeUsed := context.makeInstanceFromCompiledCode(codeHash, gasLimit, newCode)
	if compiledCodeUsed {
		return nil
	}

	return context.makeInstanceFromContractByteCode(contract, codeHash, gasLimit, newCode)
}

func (context *runtimeContext) makeInstanceFromContractByteCode(contract []byte, codeHash []byte, gasLimit uint64, newCode bool) error {
	log.Trace("Creating a new Wasmer instance")

	gasSchedule := context.host.Metering().GasSchedule()
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	newInstance, err := wasmer.NewInstanceWithOptions(contract, options)
	if err != nil {
		context.instance = nil
		return err
	}

	context.instance = newInstance

	if newCode || len(codeHash) == 0 {
		codeHash, err = context.host.Crypto().Sha256(contract)
		if err != nil {
			context.CleanWasmerInstance()
			return err
		}
	}

	context.saveCompiledCode(codeHash)

	idContext := arwen.AddHostContext(context.host)
	context.instance.SetContextData(idContext)

	if newCode {
		err = context.VerifyContractCode()
		if err != nil {
			context.CleanWasmerInstance()
			return err
		}
	}

	if context.useWarmInstance {
		context.warmInstanceAddress = context.GetSCAddress()
		context.warmInstance = context.instance
	}

	return nil
}

func (context *runtimeContext) GetSCCode() ([]byte, error) {
	blockchain := context.host.Blockchain()
	code, err := blockchain.GetCode(context.scAddress)
	if err != nil {
		return nil, err
	}

	context.codeSize = uint64(len(code))
	return code, nil
}

func (context *runtimeContext) GetSCCodeSize() uint64 {
	return context.codeSize
}

func (context *runtimeContext) saveCompiledCode(codeHash []byte) {
	compiledCode, err := context.instance.Cache()
	if err != nil {
		log.Error("getCompiledCode from instance", "error", err)
	}

	blockchain := context.host.Blockchain()
	blockchain.SaveCompiledCode(codeHash, compiledCode)
}

func (context *runtimeContext) IsWarmInstance() bool {
	if context.instance != nil && context.instance == context.warmInstance {
		return true
	}

	return false
}

func (context *runtimeContext) ResetWarmInstance() {
	if context.instance == nil {
		return
	}

	arwen.RemoveHostContext(*context.instance.Data)
	context.instance.Clean()

	context.instance = nil
	context.warmInstanceAddress = nil
	context.warmInstance = nil
}

func (context *runtimeContext) MustVerifyNextContractCode() {
	context.verifyCode = true
}

func (context *runtimeContext) SetMaxInstanceCount(maxInstances uint64) {
	context.maxWasmerInstances = maxInstances
}

func (context *runtimeContext) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	context.SetVMInput(&input.VMInput)
	context.scAddress = input.RecipientAddr
	context.callFunction = input.Function

	context.asyncContext = arwen.NewAsyncContext()
	context.asyncContext.CallerAddr = input.CallerAddr
}

func (context *runtimeContext) SetCustomCallFunction(callFunction string) {
	context.callFunction = callFunction
}

func (context *runtimeContext) PushState() {
	newState := &runtimeContext{
		vmInput:          context.vmInput,
		scAddress:        context.scAddress,
		callFunction:     context.callFunction,
		readOnly:         context.readOnly,
		defaultAsyncCall: context.defaultAsyncCall,
		asyncContext:     context.asyncContext,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *runtimeContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.vmInput = prevState.vmInput
	context.scAddress = prevState.scAddress
	context.callFunction = prevState.callFunction
	context.readOnly = prevState.readOnly
	context.defaultAsyncCall = prevState.defaultAsyncCall
	context.asyncContext = prevState.asyncContext
}

func (context *runtimeContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	context.stateStack = context.stateStack[:stateStackLen-1]
}

func (context *runtimeContext) ClearStateStack() {
	context.stateStack = make([]*runtimeContext, 0)
}

func (context *runtimeContext) PushInstance() {
	context.instanceStack = append(context.instanceStack, context.instance)
}

func (context *runtimeContext) PopInstance() {
	instanceStackLen := len(context.instanceStack)
	prevInstance := context.instanceStack[instanceStackLen-1]
	context.instanceStack = context.instanceStack[:instanceStackLen-1]

	context.CleanWasmerInstance()
	context.instance = prevInstance
}

func (context *runtimeContext) RunningInstancesCount() uint64 {
	return uint64(len(context.instanceStack))
}

func (context *runtimeContext) ClearInstanceStack() {
	for _, instance := range context.instanceStack {
		instance.Clean()
	}
	context.instanceStack = make([]*wasmer.Instance, 0)
}

func (context *runtimeContext) GetVMType() []byte {
	return context.vmType
}

func (context *runtimeContext) GetVMInput() *vmcommon.VMInput {
	return context.vmInput
}

func (context *runtimeContext) SetVMInput(vmInput *vmcommon.VMInput) {
	if !context.host.IsArwenV2Enabled() || vmInput == nil {
		context.vmInput = vmInput
		return
	}

	context.vmInput = &vmcommon.VMInput{
		CallType:      vmInput.CallType,
		GasPrice:      vmInput.GasPrice,
		GasProvided:   vmInput.GasProvided,
		GasLocked:     vmInput.GasLocked,
		CallValue:     big.NewInt(0),
		ESDTValue:     big.NewInt(0),
		ESDTTokenName: nil,
	}

	if vmInput.CallValue != nil {
		context.vmInput.CallValue.Set(vmInput.CallValue)
	}

	if len(vmInput.CallerAddr) > 0 {
		context.vmInput.CallerAddr = make([]byte, len(vmInput.CallerAddr))
		copy(context.vmInput.CallerAddr, vmInput.CallerAddr)
	}

	if vmInput.ESDTValue != nil {
		context.vmInput.ESDTValue.Set(vmInput.ESDTValue)
	}

	if len(vmInput.ESDTTokenName) > 0 {
		context.vmInput.ESDTTokenName = make([]byte, len(vmInput.ESDTTokenName))
		copy(context.vmInput.ESDTTokenName, vmInput.ESDTTokenName)
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

func (context *runtimeContext) GetSCAddress() []byte {
	return context.scAddress
}

func (context *runtimeContext) SetSCAddress(scAddress []byte) {
	context.scAddress = scAddress
}

func (context *runtimeContext) GetCurrentTxHash() []byte {
	return context.vmInput.CurrentTxHash
}

func (context *runtimeContext) GetOriginalTxHash() []byte {
	return context.vmInput.OriginalTxHash
}

func (context *runtimeContext) GetPrevTxHash() []byte {
	return context.vmInput.PrevTxHash
}

func (context *runtimeContext) Function() string {
	return context.callFunction
}

func (context *runtimeContext) Arguments() [][]byte {
	return context.vmInput.Arguments
}

func (context *runtimeContext) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	const numMinUpgradeArguments = 2

	arguments := context.vmInput.Arguments
	if len(arguments) < numMinUpgradeArguments {
		return nil, nil, arwen.ErrInvalidUpgradeArguments
	}

	code := arguments[0]
	codeMetadata := arguments[1]
	context.vmInput.Arguments = context.vmInput.Arguments[numMinUpgradeArguments:]
	return code, codeMetadata, nil
}

func (context *runtimeContext) FailExecution(err error) {
	context.host.Output().SetReturnCode(vmcommon.ExecutionFailed)

	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = "execution failed"
	}

	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointExecutionFailed)
}

func (context *runtimeContext) SignalUserError(message string) {
	context.host.Output().SetReturnCode(vmcommon.UserError)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(arwen.BreakpointSignalError)
}

func (context *runtimeContext) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	context.instance.SetBreakpointValue(uint64(value))
}

func (context *runtimeContext) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(context.instance.GetBreakpointValue())
}

func (context *runtimeContext) VerifyContractCode() error {
	if !context.verifyCode {
		return nil
	}

	context.verifyCode = false

	err := context.validator.verifyMemoryDeclaration(context.instance)
	if err != nil {
		return err
	}

	err = context.validator.verifyFunctions(context.instance)
	if err != nil {
		return err
	}

	return nil
}

func (context *runtimeContext) ElrondAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) ElrondSyncExecAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) BigIntAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) CryptoAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) GetPointsUsed() uint64 {
	return context.instance.GetPointsUsed()
}

func (context *runtimeContext) SetPointsUsed(gasPoints uint64) {
	context.instance.SetPointsUsed(gasPoints)
}

func (context *runtimeContext) ReadOnly() bool {
	return context.readOnly
}

func (context *runtimeContext) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

func (context *runtimeContext) GetInstanceExports() wasmer.ExportsMap {
	return context.instance.Exports
}

func (context *runtimeContext) CleanWasmerInstance() {
	if context.instance == nil || context.IsWarmInstance() {
		return
	}

	arwen.RemoveHostContext(*context.instance.Data)
	context.instance.Clean()
	context.instance = nil
}

func (context *runtimeContext) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	exports := context.instance.Exports
	if function, ok := exports[context.callFunction]; ok {
		return function, nil
	}

	if context.callFunction == arwen.CallbackFunctionName {
		return nil, arwen.ErrNilCallbackFunction
	}

	return nil, arwen.ErrFuncNotFound
}

func (context *runtimeContext) GetInitFunction() wasmer.ExportedFunctionCallback {
	exports := context.instance.Exports
	if init, ok := exports[arwen.InitFunctionName]; ok {
		return init
	}

	return nil
}

func (context *runtimeContext) SetDefaultAsyncCall(asyncCall *arwen.AsyncCall) {
	context.defaultAsyncCall = asyncCall
}

func (context *runtimeContext) ExecuteAsyncCall(address []byte, data []byte, value []byte) error {
	metering := context.host.Metering()
	err := metering.UseGasForAsyncStep()
	if err != nil {
		return err
	}

	gasToLock := uint64(0)
	shouldLockGas := context.HasCallbackMethod() || !context.host.IsDynamicGasLockingEnabled()
	if shouldLockGas {
		gasToLock = metering.ComputeGasLockedForAsync()
		err = metering.UseGasBounded(gasToLock)
		if err != nil {
			return err
		}
	}

	context.SetDefaultAsyncCall(&arwen.AsyncCall{
		Destination: address,
		Data:        data,
		GasLimit:    metering.GasLeft(),
		GasLocked:   gasToLock,
		ValueBytes:  value,
	})
	context.SetRuntimeBreakpointValue(arwen.BreakpointAsyncCall)

	return nil
}

func (context *runtimeContext) CreateAndAddAsyncCall(
	groupID []byte,
	address []byte,
	data []byte,
	value []byte,
	successCallback []byte,
	errorCallback []byte,
	gas uint64,
) error {
	metering := context.host.Metering()

	err := metering.UseGasForAsyncStep()
	if err != nil {
		return err
	}

	var shouldLockGas bool

	if !context.host.IsDynamicGasLockingEnabled() {
		// Legacy mode: static gas locking, always enabled
		shouldLockGas = true
	} else {
		// Dynamic mode: lock only if callBack() exists
		shouldLockGas = context.HasCallbackMethod()
	}

	gasToLock := uint64(0)
	if shouldLockGas {
		gasToLock = metering.ComputeGasLockedForAsync()
		err = metering.UseGasBounded(gasToLock)
		if err != nil {
			return err
		}
	}

	return context.AddAsyncCall(groupID, &arwen.AsyncCall{
		Status:          arwen.AsyncCallPending,
		Destination:     address,
		Data:            data,
		ValueBytes:      value,
		SuccessCallback: string(successCallback),
		ErrorCallback:   string(errorCallback),
		ProvidedGas:     gas,
		GasLocked:       gasToLock,
	})
}

func (context *runtimeContext) AddAsyncCall(groupIDBytes []byte, asyncCall *arwen.AsyncCall) error {
	groupID := string(groupIDBytes)
	if context.host.IsBuiltinFunctionName(asyncCall.SuccessCallback) {
		return arwen.ErrCannotUseBuiltinAsCallback
	}
	if context.host.IsBuiltinFunctionName(asyncCall.ErrorCallback) {
		return arwen.ErrCannotUseBuiltinAsCallback
	}

	asyncCallGroup, ok := context.asyncContext.GetAsyncCallGroup(groupID)
	if !ok {
		asyncCallGroup = arwen.NewAsyncCallGroup(groupID)
		context.asyncContext.AddAsyncGroup(asyncCallGroup)
	}

	asyncCallGroup.AddAsyncCall(asyncCall)

	return nil
}

func (context *runtimeContext) GetAsyncContext() *arwen.AsyncContext {
	return context.asyncContext
}

func (context *runtimeContext) GetAsyncCallGroup(groupID []byte) (*arwen.AsyncCallGroup, error) {
	asyncCallGroup, ok := context.asyncContext.GetAsyncCallGroup(string(groupID))
	if !ok {
		return nil, arwen.ErrAsyncCallGroupDoesNotExist
	}

	return asyncCallGroup, nil
}

func (context *runtimeContext) GetDefaultAsyncCall() *arwen.AsyncCall {
	return context.defaultAsyncCall
}

func (context *runtimeContext) HasCallbackMethod() bool {
	_, ok := context.instance.Exports[arwen.CallbackFunctionName]
	return ok
}

func (context *runtimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memory := context.instance.InstanceCtx.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + length)
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > memoryLength
	isRequestedEndTooLarge := requestedEnd > memoryLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("mem load: %w", arwen.ErrBadBounds)
	}
	if isLengthNegative {
		return nil, fmt.Errorf("mem load: %w", arwen.ErrNegativeLength)
	}

	result := make([]byte, length)
	if isRequestedEndTooLarge {
		copy(result, memoryView[offset:])
	} else {
		copy(result, memoryView[offset:requestedEnd])
	}

	return result, nil
}

func (context *runtimeContext) MemStore(offset int32, data []byte) error {
	dataLength := int32(len(data))
	if dataLength == 0 {
		return nil
	}

	memory := context.instance.InstanceCtx.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + dataLength)
	isOffsetTooSmall := offset < 0
	isNewPageNecessary := requestedEnd > memoryLength

	if isOffsetTooSmall {
		return arwen.ErrBadLowerBounds
	}
	if isNewPageNecessary {
		err := memory.Grow(1)
		if err != nil {
			return err
		}

		memoryView = memory.Data()
		memoryLength = memory.Length()
	}

	isRequestedEndTooLarge := requestedEnd > memoryLength
	if isRequestedEndTooLarge {
		return arwen.ErrBadUpperBounds
	}

	copy(memoryView[offset:requestedEnd], data)
	return nil
}
