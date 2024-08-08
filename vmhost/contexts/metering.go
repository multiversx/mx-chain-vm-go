package contexts

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var logMetering = logger.GetOrCreate("vm/metering")

type meteringContext struct {
	host               vmhost.VMHost
	stateStack         []*meteringContext
	gasSchedule        *config.GasCost
	blockGasLimit      uint64
	initialGasProvided uint64
	initialCost        uint64
	gasForExecution    uint64
	gasUsedByAccounts  map[string]uint64
	restoreGasEnabled  bool

	gasTracer       vmhost.GasTracing
	traceGasEnabled bool
}

// NewMeteringContext creates a new meteringContext
func NewMeteringContext(
	host vmhost.VMHost,
	gasMap config.GasScheduleMap,
	blockGasLimit uint64,
) (*meteringContext, error) {
	if check.IfNil(host) {
		return nil, vmhost.ErrNilVMHost
	}

	gasSchedule, err := config.CreateGasConfig(gasMap)
	if err != nil {
		return nil, err
	}

	context := &meteringContext{
		host:              host,
		stateStack:        make([]*meteringContext, 0),
		gasSchedule:       gasSchedule,
		blockGasLimit:     blockGasLimit,
		gasUsedByAccounts: make(map[string]uint64),
		restoreGasEnabled: true,
	}

	context.InitState()

	return context, nil
}

// InitState resets the internal state of the MeteringContext
func (context *meteringContext) InitState() {
	context.gasUsedByAccounts = make(map[string]uint64)
	context.initialGasProvided = 0
	context.initialCost = 0
	context.gasForExecution = 0
	context.gasUsedByAccounts = make(map[string]uint64)
	context.restoreGasEnabled = true

	var newGasTracer vmhost.GasTracing
	if context.traceGasEnabled {
		newGasTracer = NewEnabledGasTracer()
	} else {
		newGasTracer = NewDisabledGasTracer()
	}
	context.gasTracer = newGasTracer
}

// InitStateFromContractCallInput initializes the internal state of the
// MeteringContext using values taken from the provided ContractCallInput
func (context *meteringContext) InitStateFromContractCallInput(input *vmcommon.VMInput) {
	context.InitState()
	context.unlockGasIfAsyncCallback(input)
	context.initialGasProvided = input.GasProvided
	context.gasForExecution = input.GasProvided
}

// PushState pushes the current state of the MeteringContext on its internal state stack
func (context *meteringContext) PushState() {
	newState := &meteringContext{
		initialGasProvided: context.initialGasProvided,
		initialCost:        context.initialCost,
		gasForExecution:    context.gasForExecution,
		gasUsedByAccounts:  context.cloneGasUsedByAccounts(),
		restoreGasEnabled:  context.restoreGasEnabled,
	}

	context.stateStack = append(context.stateStack, newState)
}

// PopSetActiveState pops the state at the top of the internal state stack, and
// sets it as the current state
func (context *meteringContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.initialGasProvided = prevState.initialGasProvided
	context.initialCost = prevState.initialCost
	context.gasForExecution = prevState.gasForExecution
	context.gasUsedByAccounts = prevState.gasUsedByAccounts
	context.restoreGasEnabled = prevState.restoreGasEnabled
}

// PopDiscard pops the state at the top of the internal state stack, and discards it
func (context *meteringContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// PopMergeActiveState pops the state at the top of the internal stack and
// merges it into the active state
func (context *meteringContext) PopMergeActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.initialGasProvided = prevState.initialGasProvided
	context.initialCost = prevState.initialCost
	context.gasForExecution = prevState.gasForExecution
	context.restoreGasEnabled = prevState.restoreGasEnabled

	context.addToGasUsedByAccounts(prevState.gasUsedByAccounts)
}

func (context *meteringContext) cloneGasUsedByAccounts() map[string]uint64 {
	clone := make(map[string]uint64, len(context.gasUsedByAccounts))

	for address, gasUsed := range context.gasUsedByAccounts {
		clone[address] = gasUsed
	}

	return clone
}

func (context *meteringContext) addToGasUsedByAccounts(gasUsed map[string]uint64) {
	for address, gas := range gasUsed {
		context.gasUsedByAccounts[address] += gas
	}
}

// UpdateGasStateOnSuccess performs final gas accounting after a successful execution.
func (context *meteringContext) UpdateGasStateOnSuccess(vmOutput *vmcommon.VMOutput) error {
	logMetering.Trace("UpdateGasStateOnSuccess")
	context.updateSCGasUsed()
	err := context.setGasUsedToOutputAccounts(vmOutput)
	if err != nil {
		return err
	}

	err = context.checkGas(vmOutput)
	if err != nil {
		return err
	}

	logMetering.Trace("UpdateGasStateOnSuccess", "vmOutput.GasRemaining", vmOutput.GasRemaining)
	logMetering.Trace("UpdateGasStateOnSuccess", "instance gas left", context.GasLeft())

	return nil
}

// UpdateGasStateOnFailure performs final gas accounting after a failed execution.
func (context *meteringContext) UpdateGasStateOnFailure(_ *vmcommon.VMOutput) {
	logMetering.Trace("UpdateGasStateOnFailure")
	runtime := context.host.Runtime()
	output := context.host.Output()

	account, _ := output.GetOutputAccount(runtime.GetContextAddress())
	account.GasUsed = math.AddUint64(account.GasUsed, context.GetGasProvided())
	logMetering.Trace("UpdateGasStateOnFailure", "gas used", account.GasUsed)
	logMetering.Trace("UpdateGasStateOnFailure", "instance gas left", context.GasLeft())
}

func (context *meteringContext) updateSCGasUsed() {
	runtime := context.host.Runtime()
	output := context.host.Output()

	currentAccountAddress := runtime.GetContextAddress()
	currentContractAccount, _ := output.GetOutputAccount(currentAccountAddress)
	outputAccounts := context.host.Output().GetOutputAccounts()

	gasTransferredByCurrentAccount := context.getGasTransferredByAccount(currentContractAccount)
	gasUsedByOthers := context.getGasUsedByAllOtherAccounts(outputAccounts)

	gasUsed := context.GasSpentByContract()
	gasUsed = math.SubUint64(gasUsed, gasTransferredByCurrentAccount)
	gasUsed = math.SubUint64(gasUsed, gasUsedByOthers)

	context.gasUsedByAccounts[string(currentAccountAddress)] = gasUsed
}

// TrackGasUsedByOutOfVMFunction computes the gas used by a builtin function
// execution or a function executed via blockchain on another VM  and consumes
// it on the current contract instance.
func (context *meteringContext) TrackGasUsedByOutOfVMFunction(
	builtinInput *vmcommon.ContractCallInput,
	builtinOutput *vmcommon.VMOutput,
	postBuiltinInput *vmcommon.ContractCallInput,
) {

	gasUsed := math.SubUint64(builtinInput.GasProvided, builtinOutput.GasRemaining)

	// If the builtin function indicated that there's a follow-up SC execution
	// after itself, then it has reserved gas for the SC in postBuiltinInput.
	// This gas must not be tracked as if it was used by the builtin function
	// (i.e. used on the instance of the caller).
	if postBuiltinInput != nil {
		gasUsed = math.SubUint64(gasUsed, postBuiltinInput.GasProvided)
	}

	context.useGas(gasUsed)
	logMetering.Trace("gas used by builtin function", "gas", gasUsed)
}

func (context *meteringContext) checkGas(vmOutput *vmcommon.VMOutput) error {
	logMetering.Trace("check gas")
	gasUsed := context.getCurrentTotalUsedGas()
	totalGas := math.AddUint64(gasUsed, vmOutput.GasRemaining)
	gasProvided := context.GetGasProvided()

	context.PrintState()
	if totalGas != gasProvided {
		logOutput.Error("gas usage mismatch", "total gas", totalGas, "gas provided", gasProvided)
		return vmhost.ErrInputAndOutputGasDoesNotMatch
	}

	return nil
}

func (context *meteringContext) getCurrentTotalUsedGas() uint64 {
	outputAccounts := context.host.Output().GetOutputAccounts()

	gasUsed := uint64(0)
	for _, outputAccount := range outputAccounts {
		gasTransferred := context.getGasTransferredByAccount(outputAccount)
		gasUsed = math.AddUint64(gasUsed, outputAccount.GasUsed)
		gasUsed = math.AddUint64(gasUsed, gasTransferred)
	}

	return gasUsed
}

func (context *meteringContext) getGasUsedByAllOtherAccounts(outputAccounts map[string]*vmcommon.OutputAccount) uint64 {
	gasUsedAndTransferred := uint64(0)
	currentAccountAddress := string(context.host.Runtime().GetContextAddress())
	for address, account := range outputAccounts {
		gasTransferred := context.getGasTransferredByAccount(account)

		gasUsed := uint64(0)
		if address != currentAccountAddress {
			gasUsed = context.gasUsedByAccounts[address]
		}

		gasUsedAndTransferred = math.AddUint64(gasUsedAndTransferred, gasUsed)
		gasUsedAndTransferred = math.AddUint64(gasUsedAndTransferred, gasTransferred)
	}

	return gasUsedAndTransferred
}

func (context *meteringContext) getGasTransferredByAccount(account *vmcommon.OutputAccount) uint64 {
	gasUsed := uint64(0)
	for _, outputTransfer := range account.OutputTransfers {
		gasUsed = math.AddUint64(gasUsed, outputTransfer.GasLimit)
		gasUsed = math.AddUint64(gasUsed, outputTransfer.GasLocked)
	}

	return gasUsed
}

func (context *meteringContext) setGasUsedToOutputAccounts(vmOutput *vmcommon.VMOutput) error {
	for address, account := range vmOutput.OutputAccounts {
		account.GasUsed = context.gasUsedByAccounts[address]
	}

	for address := range context.gasUsedByAccounts {
		_, exists := vmOutput.OutputAccounts[address]
		if !exists {
			return fmt.Errorf("expected OutputAccount has used gas but is missing")
		}
	}

	return nil
}

// ClearStateStack reinitializes the internal state stack to an empty stack
func (context *meteringContext) ClearStateStack() {
	context.stateStack = make([]*meteringContext, 0)
	context.gasTracer = nil
}

// unlockGasIfAsyncCallback unlocks the locked gas if the call type is async callback
func (context *meteringContext) unlockGasIfAsyncCallback(input *vmcommon.VMInput) {
	if input.CallType != vm.AsynchronousCallBack {
		return
	}

	gasProvided := math.AddUint64(input.GasProvided, input.GasLocked)

	context.gasForExecution = gasProvided
	input.GasProvided = gasProvided
	input.GasLocked = 0
}

// GasSchedule returns the current gas schedule
func (context *meteringContext) GasSchedule() *config.GasCost {
	return context.gasSchedule
}

// SetGasSchedule sets the gas schedule to the given gas map
func (context *meteringContext) SetGasSchedule(gasMap config.GasScheduleMap) {
	gasSchedule, err := config.CreateGasConfig(gasMap)
	if err != nil {
		logMetering.Error("SetGasSchedule createGasConfig", "error", err)
		return
	}
	context.gasSchedule = gasSchedule
}

// useGas consumes the specified amount of gas on the currently running Wasmer instance.
func (context *meteringContext) useGas(gas uint64) {
	gasUsed := math.AddUint64(context.host.Runtime().GetPointsUsed(), gas)
	context.host.Runtime().SetPointsUsed(gasUsed)
	logMetering.Trace("used gas", "gas", gas)
}

// useAndTraceGas sets in the runtime context the given gas as gas used and adds to current trace
func (context *meteringContext) useAndTraceGas(gas uint64) {
	context.useGas(gas)
	context.traceGas(gas)
}

// useGasAndAddTracedGas sets in the runtime context the given gas as gas used and adds to current trace
func (context *meteringContext) useGasAndAddTracedGas(functionName string, gas uint64) {
	context.useGas(gas)
	context.addToGasTrace(functionName, gas)
}

// UseGasBoundedAndAddTracedGas sets in the runtime context the given gas as gas used and adds to current trace
func (context *meteringContext) UseGasBoundedAndAddTracedGas(functionName string, gas uint64) error {
	err := context.UseGasBounded(gas)
	if err != nil {
		return err
	}

	context.addToGasTrace(functionName, gas)
	return nil
}

// GetGasTrace returns the gasTrace map
func (context *meteringContext) GetGasTrace() map[string]map[string][]uint64 {
	return context.gasTracer.GetGasTrace()
}

// RestoreGas deducts the specified amount of gas from the gas currently spent on the running Wasmer instance.
func (context *meteringContext) RestoreGas(gas uint64) {
	if !context.restoreGasEnabled {
		logMetering.Trace("restore gas disabled", "gas not restored", gas)
		return
	}
	gasUsed := context.host.Runtime().GetPointsUsed()
	if gas <= gasUsed {
		gasUsed = math.SubUint64(gasUsed, gas)
		context.host.Runtime().SetPointsUsed(gasUsed)
	}
	logMetering.Trace("restored gas", "gas", gas)
}

// DisableRestoreGas disables the RestoreGas() method; needed for gas management of async calls.
func (context *meteringContext) DisableRestoreGas() {
	context.restoreGasEnabled = false
}

// EnableRestoreGas enables the RestoreGas() method; needed for gas management of async calls.
func (context *meteringContext) EnableRestoreGas() {
	context.restoreGasEnabled = true
}

// FreeGas refunds the specified amount of gas to the caller.
func (context *meteringContext) FreeGas(gas uint64) {
	refund := math.AddUint64(context.host.Output().GetRefund(), gas)
	context.host.Output().SetRefund(refund)
}

// GasLeft computes the amount of gas left on the currently running Wasmer instance.
func (context *meteringContext) GasLeft() uint64 {
	gasProvided := context.gasForExecution
	gasUsed := context.host.Runtime().GetPointsUsed()

	if gasProvided < gasUsed {
		return 0
	}

	return gasProvided - gasUsed
}

// GasSpentByContract calculates the entire gas consumption of the contract,
// without any gas forwarding.
func (context *meteringContext) GasSpentByContract() uint64 {
	runtime := context.host.Runtime()
	executionGasUsed := runtime.GetPointsUsed()

	gasSpent := math.AddUint64(context.initialCost, executionGasUsed)
	return gasSpent
}

// GasUsedForExecution returns the actual gas used for execution for the contract which needs to be restored
func (context *meteringContext) GasUsedForExecution() uint64 {
	gasUsed := context.GasSpentByContract()
	gasUsed = math.SubUint64(gasUsed, context.initialCost)
	return gasUsed
}

// GetGasForExecution returns the gas left after the deduction of the initial gas from the provided gas
func (context *meteringContext) GetGasForExecution() uint64 {
	return context.gasForExecution
}

// GetGasProvided returns the fully provided gas for the sc execution
func (context *meteringContext) GetGasProvided() uint64 {
	return context.initialGasProvided
}

// GetSCPrepareInitialCost return the initial prepare cost for the sc execution
func (context *meteringContext) GetSCPrepareInitialCost() uint64 {
	return context.initialCost
}

// BoundGasLimit returns the maximum between the provided amount and the gas
// left on the currently running Wasmer instance.
func (context *meteringContext) BoundGasLimit(value int64) uint64 {
	gasLeft := context.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

// UseGasForAsyncStep consumes the AsyncCallStep gas cost on the currently
// running Wasmer instance
func (context *meteringContext) UseGasForAsyncStep() error {
	gasSchedule := context.GasSchedule().BaseOpsAPICost
	gasToDeduct := gasSchedule.AsyncCallStep
	return context.UseGasBounded(gasToDeduct)
}

// UseGasBounded consumes the specified amount of gas on the currently running
// Wasmer instance, but returns an error if there is not enough gas left.
func (context *meteringContext) UseGasBounded(gasToUse uint64) error {
	if context.GasLeft() < gasToUse {
		return vmhost.ErrNotEnoughGas
	}
	context.useGas(gasToUse)
	context.traceGas(gasToUse)
	return nil
}

// ComputeExtraGasLockedForAsync calculates the minimum amount of gas to lock for async callbacks
func (context *meteringContext) ComputeExtraGasLockedForAsync() uint64 {
	baseGasSchedule := context.GasSchedule().BaseOperationCost
	apiGasSchedule := context.GasSchedule().BaseOpsAPICost
	codeSize := context.host.Runtime().GetSCCodeSize()
	costPerByte := baseGasSchedule.AoTPreparePerByte

	// Exact amount of gas required to compile this SC again, to execute the callback
	compilationGasLock := math.MulUint64(codeSize, costPerByte)

	// Minimum amount required to execute the callback
	executionGasLock := math.AddUint64(apiGasSchedule.AsyncCallStep, apiGasSchedule.AsyncCallbackGasLock)
	gasLockedForAsync := math.AddUint64(compilationGasLock, executionGasLock)

	return gasLockedForAsync
}

// GetGasLocked returns the locked gas
func (context *meteringContext) GetGasLocked() uint64 {
	input := context.host.Runtime().GetVMInput()
	return input.GasLocked
}

// BlockGasLimit returns the maximum amount of gas allowed to be consumed in a block.
func (context *meteringContext) BlockGasLimit() uint64 {
	return context.blockGasLimit
}

// DeductInitialGasForExecution deducts gas for compilation and locks gas if the execution is an asynchronous call
func (context *meteringContext) DeductInitialGasForExecution(contract []byte) error {
	costPerByte := context.gasSchedule.BaseOperationCost.AoTPreparePerByte
	baseCost := context.gasSchedule.BaseOperationCost.GetCode
	err := context.deductInitialGas(contract, baseCost, costPerByte)
	if err != nil {
		return err
	}

	return nil
}

// DeductInitialGasForDirectDeployment deducts gas for the deployment of a contract initiated by a Transaction
func (context *meteringContext) DeductInitialGasForDirectDeployment(input vmhost.CodeDeployInput) error {
	return context.deductInitialGas(
		input.ContractCode,
		context.gasSchedule.BaseOpsAPICost.CreateContract,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
}

// DeductInitialGasForIndirectDeployment deducts gas for the deployment of a contract initiated by another SmartContract
func (context *meteringContext) DeductInitialGasForIndirectDeployment(input vmhost.CodeDeployInput) error {
	return context.deductInitialGas(
		input.ContractCode,
		0,
		context.gasSchedule.BaseOperationCost.CompilePerByte,
	)
}

func (context *meteringContext) deductInitialGas(
	code []byte,
	baseCost uint64,
	costPerByte uint64,
) error {
	input := context.host.Runtime().GetVMInput()
	codeLength := uint64(len(code))
	codeCost := math.MulUint64(codeLength, costPerByte)
	initialCost := math.AddUint64(baseCost, codeCost)

	if initialCost > input.GasProvided {
		return vmhost.ErrNotEnoughGas
	}

	context.initialCost = initialCost
	context.gasForExecution = input.GasProvided - initialCost
	return nil
}

// SetGasTracing enables/disables gas tracing
func (context *meteringContext) SetGasTracing(enableGasTracing bool) {
	context.traceGasEnabled = enableGasTracing
	if context.traceGasEnabled {
		context.gasTracer = NewEnabledGasTracer()
	} else {
		context.gasTracer = NewDisabledGasTracer()
	}
}

// StartGasTracing sets initial trace for the upcoming gas usage.
func (context *meteringContext) StartGasTracing(functionName string) {
	if context.traceGasEnabled {
		scAddress := context.getSCAddress()
		if len(scAddress) != 0 {
			context.gasTracer.BeginTrace(scAddress, functionName)
		}
	}
}

func (context *meteringContext) traceGas(usedGas uint64) {
	context.gasTracer.AddToCurrentTrace(usedGas)
}

func (context *meteringContext) addToGasTrace(functionName string, usedGas uint64) {
	scAddress := context.getSCAddress()
	context.gasTracer.AddTracedGas(scAddress, functionName, usedGas)
}

func (context *meteringContext) getSCAddress() string {
	return string(context.host.Runtime().GetContextAddress())
}

// PrintState dumps the internal state of the meteringContext to the TRACE output
func (context *meteringContext) PrintState() {
	sc := context.host.Runtime().GetContextAddress()
	functionName := context.host.Runtime().FunctionName()
	scAccount, _ := context.host.Output().GetOutputAccount(sc)
	outputAccounts := context.host.Output().GetOutputAccounts()
	gasSpent := context.GasSpentByContract()
	gasTransferred := context.getGasTransferredByAccount(scAccount)
	gasUsedByOthers := context.getGasUsedByAllOtherAccounts(outputAccounts)
	gasUsed := gasSpent
	gasUsed = math.SubUint64(gasUsed, gasTransferred)
	gasUsed = math.SubUint64(gasUsed, gasUsedByOthers)
	logMetering.Trace("metering state", "┌----------            sc", string(sc))
	logMetering.Trace("              ", "|                function", functionName)
	logMetering.Trace("              ", "|        initial provided", context.initialGasProvided)
	logMetering.Trace("              ", "|            initial cost", context.initialCost)
	logMetering.Trace("              ", "|            gas for exec", context.gasForExecution)
	logMetering.Trace("              ", "|            instance gas", context.host.Runtime().GetPointsUsed())
	logMetering.Trace("              ", "|                gas left", context.GasLeft())
	logMetering.Trace("              ", "|         gas spent by sc", gasSpent)
	logMetering.Trace("              ", "|         gas transferred", gasTransferred)
	logMetering.Trace("              ", "|      gas used by others", gasUsedByOthers)
	logMetering.Trace("              ", "| adjusted gas used by sc", gasUsed)
	for key, gas := range context.gasUsedByAccounts {
		logMetering.Trace("              ", "| gas per acct", gas, "key", key)
	}
	logMetering.Trace("              ", "└ stack size", len(context.stateStack))
}
