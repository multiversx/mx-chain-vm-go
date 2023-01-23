package host

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-v1_4-go/math"
	arwen "github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/contexts"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) *vmcommon.VMOutput {
	host.InitState()
	defer func() {
		errs := host.GetRuntimeErrors()
		if errs != nil {
			log.Trace("doRunSmartContractCreate full error list", "error", errs)
		}
	}()

	_, blockchain, metering, output, runtime, storage := host.GetContexts()

	var vmOutput *vmcommon.VMOutput
	defer func() {
		if vmOutput == nil || vmOutput.ReturnCode == vmcommon.ExecutionFailed {
			runtime.CleanInstance()
		}
	}()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		vmOutput = output.CreateVMOutputInCaseOfError(err)
		return vmOutput
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput:       input.VMInput,
		RecipientAddr: address,
	}
	runtime.SetVMInput(contractCallInput)
	runtime.SetCodeAddress(address)
	metering.InitStateFromContractCallInput(&input.VMInput)

	output.AddTxValueToAccount(address, input.CallValue)
	storage.SetAddress(runtime.GetContextAddress())

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err = host.performCodeDeployment(codeDeployInput)
	if err != nil {
		log.Trace("doRunSmartContractCreate", "error", err)
		vmOutput = output.CreateVMOutputInCaseOfError(err)
		return vmOutput
	}

	log.Trace("doRunSmartContractCreate",
		"retCode", vmOutput.ReturnCode,
		"message", vmOutput.ReturnMessage,
		"data", vmOutput.ReturnData)

	return vmOutput
}

func (host *vmHost) performCodeDeployment(input arwen.CodeDeployInput) (*vmcommon.VMOutput, error) {
	log.Trace("performCodeDeployment", "address", input.ContractAddress, "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	_, _, metering, output, runtime, _ := host.GetContexts()

	err := metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return nil, err
	}

	runtime.MustVerifyNextContractCode()

	err = runtime.StartWasmerInstance(input.ContractCode, metering.GetGasForExecution(), true)
	if err != nil {
		log.Trace("performCodeDeployment/StartWasmerInstance", "err", err)
		return nil, arwen.ErrContractInvalid
	}

	defer func() {
		if !contexts.WarmInstancesEnabled {
			runtime.CleanInstance()
		}
	}()

	err = host.callInitFunction()
	if err != nil {
		return nil, err
	}

	output.DeployCode(input)
	if host.enableEpochsHandler.IsRemoveNonUpdatedStorageFlagEnabled() {
		output.RemoveNonUpdatedStorage()
	}

	vmOutput := output.GetVMOutput()

	return vmOutput, nil
}

// doRunSmartContractUpgrade upgrades a contract directly
func (host *vmHost) doRunSmartContractUpgrade(input *vmcommon.ContractCallInput) *vmcommon.VMOutput {
	host.InitState()
	defer func() {
		errs := host.GetRuntimeErrors()
		if errs != nil {
			log.Trace("doRunSmartContractUpgrade full error list", "error", errs)
		}
	}()

	_, _, metering, output, runtime, storage := host.GetContexts()

	var vmOutput *vmcommon.VMOutput
	defer func() {
		if vmOutput == nil || vmOutput.ReturnCode == vmcommon.ExecutionFailed {
			runtime.CleanInstance()
		}
	}()

	runtime.InitStateFromContractCallInput(input)
	metering.InitStateFromContractCallInput(&input.VMInput)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetContextAddress())

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		vmOutput = output.CreateVMOutputInCaseOfError(arwen.ErrInvalidUpgradeArguments)
		return vmOutput
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err = host.performCodeDeployment(codeDeployInput)
	if err != nil {
		log.Trace("doRunSmartContractUpgrade", "error", err)
		vmOutput = output.CreateVMOutputInCaseOfError(err)
		return vmOutput
	}

	return vmOutput
}

func (host *vmHost) checkGasForGetCode(input *vmcommon.ContractCallInput, metering arwen.MeteringContext) error {
	getCodeBaseCost := metering.GasSchedule().BaseOperationCost.GetCode
	if input.GasProvided < getCodeBaseCost {
		return arwen.ErrNotEnoughGas
	}

	return nil
}

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) *vmcommon.VMOutput {
	host.InitState()
	defer func() {
		errs := host.GetRuntimeErrors()
		if errs != nil {
			log.Trace(fmt.Sprintf("doRunSmartContractCall full error list for %s", input.Function), "error", errs)
		}
	}()

	_, _, metering, output, runtime, storage := host.GetContexts()

	var vmOutput *vmcommon.VMOutput
	defer func() {
		if vmOutput == nil || vmOutput.ReturnCode == vmcommon.ExecutionFailed {
			host.Runtime().CleanInstance()
		}
	}()

	runtime.InitStateFromContractCallInput(input)
	metering.InitStateFromContractCallInput(&input.VMInput)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetContextAddress())

	err := host.checkGasForGetCode(input, metering)
	if err != nil {
		log.Trace("doRunSmartContractCall get code", "error", arwen.ErrNotEnoughGas)
		vmOutput = output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
		return vmOutput
	}

	contract, err := runtime.GetSCCode()
	if err != nil {
		log.Trace("doRunSmartContractCall get code", "error", arwen.ErrContractNotFound)
		vmOutput = output.CreateVMOutputInCaseOfError(arwen.ErrContractNotFound)
		return vmOutput
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		log.Trace("doRunSmartContractCall initial gas", "error", arwen.ErrNotEnoughGas)
		vmOutput = output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
		return vmOutput
	}

	err = runtime.StartWasmerInstance(contract, metering.GetGasForExecution(), false)
	if err != nil {
		vmOutput = output.CreateVMOutputInCaseOfError(arwen.ErrContractInvalid)
		return vmOutput
	}

	defer func() {
		if !contexts.WarmInstancesEnabled {
			runtime.CleanInstance()
		}
	}()

	err = host.callSCMethod()
	if err != nil {
		log.Trace("doRunSmartContractCall", "error", err)

		vmOutput = output.CreateVMOutputInCaseOfError(err)
		return vmOutput
	}

	if host.enableEpochsHandler.IsRemoveNonUpdatedStorageFlagEnabled() {
		output.RemoveNonUpdatedStorage()
	}

	vmOutput = output.GetVMOutput()

	log.Trace("doRunSmartContractCall finished",
		"retCode", vmOutput.ReturnCode,
		"message", vmOutput.ReturnMessage,
		"data", vmOutput.ReturnData)

	return vmOutput
}

func copyTxHashesFromContext(runtime arwen.RuntimeContext, input *vmcommon.ContractCallInput) {
	currentVMInput := runtime.GetVMInput()
	if len(currentVMInput.OriginalTxHash) > 0 {
		input.OriginalTxHash = currentVMInput.OriginalTxHash
	}
	if len(currentVMInput.CurrentTxHash) > 0 {
		input.CurrentTxHash = currentVMInput.CurrentTxHash
	}
	if len(currentVMInput.PrevTxHash) > 0 {
		input.PrevTxHash = currentVMInput.PrevTxHash
	}

}

// ExecuteOnDestContext pushes each context to the corresponding stack
// and initializes new contexts for executing the contract call with the given input
func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, asyncInfo *arwen.AsyncContextInfo, err error) {
	log.Trace("ExecuteOnDestContext", "caller", input.CallerAddr, "dest", input.RecipientAddr, "function", input.Function)

	scExecutionInput := input

	blockchain := host.Blockchain()
	blockchain.PushState()

	if host.IsBuiltinFunctionName(input.Function) {
		scExecutionInput, vmOutput, err = host.handleBuiltinFunctionCall(input)
		if err != nil {
			blockchain.PopSetActiveState()
			host.Runtime().AddError(err, input.Function)
			vmOutput = host.Output().CreateVMOutputInCaseOfError(err)
			return
		}
	}

	if scExecutionInput != nil {
		vmOutput, asyncInfo, err = host.executeOnDestContextNoBuiltinFunction(scExecutionInput)
	}

	if err != nil {
		blockchain.PopSetActiveState()
	} else {
		blockchain.PopDiscard()
	}

	return
}

func (host *vmHost) handleBuiltinFunctionCall(input *vmcommon.ContractCallInput) (*vmcommon.ContractCallInput, *vmcommon.VMOutput, error) {
	output := host.Output()
	postBuiltinInput, builtinOutput, err := host.callBuiltinFunction(input)
	if err != nil {
		log.Trace("ExecuteOnDestContext builtin function", "error", err)
		return nil, nil, err
	}

	output.AddToActiveState(builtinOutput)

	return postBuiltinInput, builtinOutput, nil
}

func (host *vmHost) executeOnDestContextNoBuiltinFunction(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, asyncInfo *arwen.AsyncContextInfo, err error) {
	managedTypes, _, metering, output, runtime, storage := host.GetContexts()
	managedTypes.PushState()
	managedTypes.InitState()

	output.PushState()
	output.CensorVMOutput()

	copyTxHashesFromContext(runtime, input)
	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	metering.PushState()
	metering.InitStateFromContractCallInput(&input.VMInput)

	storage.PushState()
	storage.SetAddress(runtime.GetContextAddress())

	defer func() {
		vmOutput = host.finishExecuteOnDestContext(err)

		if err == nil && vmOutput.ReturnCode != vmcommon.Ok {
			err = arwen.ErrExecutionFailed
		}
	}()

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	if input.CallType != vm.AsynchronousCallBack || input.CallValue.Cmp(arwen.Zero) == 0 {
		err = output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue, false)
		if err != nil {
			log.Trace("ExecuteOnDestContext transfer", "error", err)
			return
		}
	}

	err = host.execute(input)
	if err != nil {
		log.Trace("ExecuteOnDestContext execution", "error", err)
		return
	}

	asyncInfo = runtime.GetAsyncContextInfo()
	_, err = host.processAsyncInfo(asyncInfo)
	return
}

func (host *vmHost) finishExecuteOnDestContext(executeErr error) *vmcommon.VMOutput {
	managedTypes, _, metering, output, runtime, storage := host.GetContexts()

	var vmOutput *vmcommon.VMOutput
	if executeErr != nil {
		// Execution failed: restore contexts as if the execution didn't happen,
		// but first create a vmOutput to capture the error.
		vmOutput = output.CreateVMOutputInCaseOfError(executeErr)
	} else {
		// Retrieve the VMOutput before popping the Runtime state and the previous
		// instance, to ensure accurate GasRemaining
		vmOutput = output.GetVMOutput()
	}

	gasSpentByChildContract := metering.GasSpentByContract()

	// Restore the previous context states
	managedTypes.PopSetActiveState()
	storage.PopSetActiveState()

	if vmOutput.ReturnCode == vmcommon.Ok {
		metering.PopMergeActiveState()
		output.PopMergeActiveState()
	} else {
		metering.PopSetActiveState()
		output.PopSetActiveState()
	}

	// Return to the caller context completely
	runtime.PopSetActiveState()

	// Restore remaining gas to the caller Wasmer instance
	metering.RestoreGas(vmOutput.GasRemaining)

	log.Trace("ExecuteOnDestContext finished", "gas spent", gasSpentByChildContract)

	return vmOutput
}

// ExecuteOnSameContext executes the contract call with the given input
// on the same runtime context. Some other contexts are backed up.
func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (asyncInfo *arwen.AsyncContextInfo, err error) {
	log.Trace("ExecuteOnSameContext", "function", input.Function)

	if host.IsBuiltinFunctionName(input.Function) {
		return nil, arwen.ErrBuiltinCallOnSameContextDisallowed
	}

	managedTypes, blockchain, metering, output, runtime, _ := host.GetContexts()

	// Back up the states of the contexts (except Storage, which isn't affected
	// by ExecuteOnSameContext())
	managedTypes.PushState()
	managedTypes.InitState()
	output.PushState()

	librarySCAddress := make([]byte, len(input.RecipientAddr))
	copy(librarySCAddress, input.RecipientAddr)

	if host.enableEpochsHandler.IsRefactorContextFlagEnabled() {
		input.RecipientAddr = input.CallerAddr
	}

	copyTxHashesFromContext(runtime, input)
	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)
	runtime.SetCodeAddress(librarySCAddress)

	metering.PushState()
	metering.InitStateFromContractCallInput(&input.VMInput)

	blockchain.PushState()

	defer func() {
		runtime.AddError(err, input.Function)
		host.finishExecuteOnSameContext(err)
	}()

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err = output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue, false)
	if err != nil {
		return
	}

	err = host.execute(input)
	if err != nil {
		return
	}

	asyncInfo = runtime.GetAsyncContextInfo()
	return
}

func (host *vmHost) finishExecuteOnSameContext(executeErr error) {
	managedTypes, blockchain, metering, output, runtime, _ := host.GetContexts()

	if output.ReturnCode() != vmcommon.Ok || executeErr != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		managedTypes.PopSetActiveState()
		metering.PopSetActiveState()
		output.PopSetActiveState()
		blockchain.PopSetActiveState()
		runtime.PopSetActiveState()
		return
	}

	// Execution successful; retrieve the VMOutput before popping the Runtime
	// state and the previous instance, to ensure accurate GasRemaining and
	// GasUsed for all accounts.
	vmOutput := output.GetVMOutput()

	metering.PopMergeActiveState()
	output.PopDiscard()
	blockchain.PopDiscard()
	managedTypes.PopSetActiveState()
	runtime.PopSetActiveState()
	// Restore remaining gas to the caller (parent) Wasmer instance
	metering.RestoreGas(vmOutput.GasRemaining)
}

func (host *vmHost) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

//nolint:all
func (host *vmHost) isBuiltinFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return host.IsBuiltinFunctionName(functionName)
}

// IsBuiltinFunctionName returns true if the given function name is the same as any protocol builtin function
func (host *vmHost) IsBuiltinFunctionName(functionName string) bool {
	function, err := host.builtInFuncContainer.Get(functionName)
	if err != nil {
		return false
	}

	return function.IsActive()
}

// CreateNewContract creates a new contract indirectly (from another Smart Contract)
func (host *vmHost) CreateNewContract(input *vmcommon.ContractCreateInput) (newContractAddress []byte, err error) {
	newContractAddress = nil
	err = nil

	defer func() {
		if err != nil {
			newContractAddress = nil
		}
	}()

	_, blockchain, metering, output, runtime, _ := host.GetContexts()

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      nil,
		CodeDeployerAddress:  input.CallerAddr,
	}
	err = metering.DeductInitialGasForIndirectDeployment(codeDeployInput)
	if err != nil {
		return
	}

	if runtime.ReadOnly() {
		err = arwen.ErrInvalidCallOnReadOnlyMode
		return
	}

	newContractAddress, err = blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return
	}

	if blockchain.AccountExists(newContractAddress) {
		err = arwen.ErrDeploymentOverExistingAccount
		return
	}

	codeDeployInput.ContractAddress = newContractAddress
	output.DeployCode(codeDeployInput)

	defer func() {
		if err != nil {
			output.DeleteOutputAccount(newContractAddress)
		}
	}()

	runtime.MustVerifyNextContractCode()

	initCallInput := &vmcommon.ContractCallInput{
		RecipientAddr:     newContractAddress,
		Function:          arwen.InitFunctionName,
		AllowInitFunction: true,
		VMInput:           input.VMInput,
	}
	_, _, err = host.ExecuteOnDestContext(initCallInput)
	if err != nil {
		return
	}

	blockchain.IncreaseNonce(input.CallerAddr)

	return
}

func (host *vmHost) checkUpgradePermission(vmInput *vmcommon.ContractCallInput) error {
	contract, err := host.Blockchain().GetUserAccount(vmInput.RecipientAddr)
	if err != nil {
		return err
	}
	if check.IfNilReflect(contract) {
		return arwen.ErrNilContract
	}

	codeMetadata := vmcommon.CodeMetadataFromBytes(contract.GetCodeMetadata())
	isUpgradeable := codeMetadata.Upgradeable
	callerAddress := vmInput.CallerAddr
	ownerAddress := contract.GetOwnerAddress()
	isCallerOwner := bytes.Equal(callerAddress, ownerAddress)

	if isUpgradeable && isCallerOwner {
		return nil
	}

	return arwen.ErrUpgradeNotAllowed
}

// executeUpgrade upgrades a contract indirectly (from another contract). This
// function follows the convention of executeSmartContractCall().
func (host *vmHost) executeUpgrade(input *vmcommon.ContractCallInput) error {
	_, _, metering, output, runtime, _ := host.GetContexts()

	err := host.checkUpgradePermission(input)
	if err != nil {
		return err
	}

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		return arwen.ErrInvalidUpgradeArguments
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	err = metering.DeductInitialGasForDirectDeployment(codeDeployInput)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return err
	}

	runtime.MustVerifyNextContractCode()

	err = runtime.StartWasmerInstance(codeDeployInput.ContractCode, metering.GetGasForExecution(), true)
	if err != nil {
		log.Trace("performCodeDeployment/StartWasmerInstance", "err", err)
		return arwen.ErrContractInvalid
	}

	err = host.callInitFunction()
	if err != nil {
		return err
	}

	output.DeployCode(codeDeployInput)
	if output.ReturnCode() != vmcommon.Ok {
		return arwen.ErrReturnCodeNotOk
	}

	return nil
}

// execute executes an indirect call to a smart contract, assuming there is an
// already-running Wasmer instance of another contract that has requested the
// indirect call. This method creates a new Wasmer instance and pushes the
// previous one onto the Runtime instance stack, but it will not pop the
// previous instance back - that remains the responsibility of the calling
// code. Also, this method does not restore the gas remaining after the
// indirect call, it does not push the states of any Host Context onto their
// respective stacks, nor does it pop any state stack. Handling the state
// stacks and the remaining gas are responsibilities of the calling code, which
// must push and pop as required, before and after calling this method, and
// handle the remaining gas. These principles also apply to indirect contract
// upgrading (via host.executeUpgrade(), which also does not pop the previous
// instance from the Runtime instance stack, nor does it restore the remaining
// gas).
func (host *vmHost) execute(input *vmcommon.ContractCallInput) error {
	_, _, metering, output, runtime, _ := host.GetContexts()

	if host.isInitFunctionBeingCalled() && !input.AllowInitFunction {
		return arwen.ErrInitFuncCalledInRun
	}

	// Use all gas initially, on the Wasmer instance of the caller. In case of
	// successful execution, the unused gas will be restored.
	metering.UseGas(input.GasProvided)

	isUpgrade := input.Function == arwen.UpgradeFunctionName
	if isUpgrade {
		return host.executeUpgrade(input)
	}

	contract, err := runtime.GetSCCode()
	if err != nil {
		return err
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return err
	}

	// Replace the current Wasmer instance of the Runtime with a new one; this
	// assumes that the instance was preserved on the Runtime instance stack
	// before calling executeSmartContractCall().
	err = runtime.StartWasmerInstance(contract, metering.GetGasForExecution(), false)
	if err != nil {
		return err
	}

	err = host.callSCMethodIndirect()
	if err != nil {
		return err
	}

	if output.ReturnCode() != vmcommon.Ok {
		return arwen.ErrReturnCodeNotOk
	}

	return nil
}

func (host *vmHost) callSCMethodIndirect() error {
	function, err := host.Runtime().GetFunctionToCall()
	if err != nil {
		if errors.Is(err, arwen.ErrNilCallbackFunction) {
			return nil
		}
		return err
	}

	err = host.Runtime().CallFunction(function)
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	return err
}

// ExecuteESDTTransfer calls the process built in function with the given transfer for ESDT/ESDTNFT if nonce > 0
// there are no NFTs with nonce == 0, it will call multi transfer if multiple tokens are sent
func (host *vmHost) ExecuteESDTTransfer(destination []byte, sender []byte, transfers []*vmcommon.ESDTTransfer, callType vm.CallType) (*vmcommon.VMOutput, uint64, error) {
	if len(transfers) == 0 {
		return nil, 0, arwen.ErrFailedTransfer
	}

	if host.Runtime().ReadOnly() && host.CheckExecuteReadOnly() {
		return nil, 0, arwen.ErrInvalidCallOnReadOnlyMode
	}

	_, _, metering, _, runtime, _ := host.GetContexts()

	esdtTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    callType,
			GasPrice:    runtime.GetVMInput().GasPrice,
			GasProvided: metering.GasLeft(),
			GasLocked:   0,
		},
		RecipientAddr:     destination,
		Function:          core.BuiltInFunctionESDTTransfer,
		AllowInitFunction: false,
	}

	if len(transfers) == 1 {
		if transfers[0].ESDTTokenNonce > 0 {
			esdtTransferInput.Function = core.BuiltInFunctionESDTNFTTransfer
			esdtTransferInput.RecipientAddr = esdtTransferInput.CallerAddr
			nonceAsBytes := big.NewInt(0).SetUint64(transfers[0].ESDTTokenNonce).Bytes()
			esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, transfers[0].ESDTTokenName, nonceAsBytes, transfers[0].ESDTValue.Bytes(), destination)
		} else {
			esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, transfers[0].ESDTTokenName, transfers[0].ESDTValue.Bytes())
		}
	} else {
		esdtTransferInput.Function = core.BuiltInFunctionMultiESDTNFTTransfer
		esdtTransferInput.RecipientAddr = esdtTransferInput.CallerAddr
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, destination, big.NewInt(int64(len(transfers))).Bytes())
		for _, transfer := range transfers {
			nonceAsBytes := big.NewInt(0).SetUint64(transfer.ESDTTokenNonce).Bytes()
			esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, transfer.ESDTTokenName, nonceAsBytes, transfer.ESDTValue.Bytes())
		}
	}

	vmOutput, err := host.Blockchain().ProcessBuiltInFunction(esdtTransferInput)
	log.Trace("ESDT transfer", "sender", sender, "dest", destination)
	for _, transfer := range transfers {
		log.Trace("ESDT transfer", "token", transfer.ESDTTokenName, "nonce", transfer.ESDTTokenNonce, "value", transfer.ESDTValue)
	}
	if err != nil {
		log.Trace("ESDT transfer", "error", err)
		return vmOutput, esdtTransferInput.GasProvided, err
	}
	if vmOutput.ReturnCode != vmcommon.Ok {
		log.Trace("ESDT transfer", "error", err, "retcode", vmOutput.ReturnCode, "message", vmOutput.ReturnMessage)
		return vmOutput, esdtTransferInput.GasProvided, arwen.ErrExecutionFailed
	}

	gasConsumed := math.SubUint64(esdtTransferInput.GasProvided, vmOutput.GasRemaining)
	for _, outAcc := range vmOutput.OutputAccounts {
		for _, transfer := range outAcc.OutputTransfers {
			gasConsumed = math.SubUint64(gasConsumed, transfer.GasLimit)
		}
	}
	if callType != vm.AsynchronousCallBack {
		if metering.GasLeft() < gasConsumed {
			log.Trace("ESDT transfer", "error", arwen.ErrNotEnoughGas)
			return vmOutput, esdtTransferInput.GasProvided, arwen.ErrNotEnoughGas
		}
		metering.UseGas(gasConsumed)
	}

	return vmOutput, gasConsumed, nil
}

func (host *vmHost) callBuiltinFunction(input *vmcommon.ContractCallInput) (*vmcommon.ContractCallInput, *vmcommon.VMOutput, error) {
	metering := host.Metering()

	if host.Runtime().ReadOnly() && host.CheckExecuteReadOnly() {
		return nil, nil, arwen.ErrInvalidCallOnReadOnlyMode
	}

	vmOutput, err := host.Blockchain().ProcessBuiltInFunction(input)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return nil, nil, err
	}

	newVMInput, err := host.isSCExecutionAfterBuiltInFunc(input, vmOutput)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return nil, nil, err
	}

	if newVMInput != nil {
		for _, outAcc := range vmOutput.OutputAccounts {
			outAcc.OutputTransfers = make([]vmcommon.OutputTransfer, 0)
		}
	}

	metering.TrackGasUsedByBuiltinFunction(input, vmOutput, newVMInput)

	host.addESDTTransferToVMOutputSCIntraShardCall(input, vmOutput)

	return newVMInput, vmOutput, nil
}

// add output transfer of esdt transfer when sc calling another sc intra shard to log the transfer information
func (host *vmHost) addESDTTransferToVMOutputSCIntraShardCall(
	input *vmcommon.ContractCallInput,
	output *vmcommon.VMOutput,
) {
	if output.ReturnCode != vmcommon.Ok {
		return
	}

	parsedTransfer, err := host.esdtTransferParser.ParseESDTTransfers(input.CallerAddr, input.RecipientAddr, input.Function, input.Arguments)
	if err != nil {
		return
	}

	if !host.AreInSameShard(input.CallerAddr, parsedTransfer.RcvAddr) {
		return
	}

	addOutputTransferToVMOutput(input.Function, input.Arguments, input.CallerAddr, parsedTransfer.RcvAddr, input.CallType, output)
}

func addOutputTransferToVMOutput(
	function string,
	arguments [][]byte,
	sender []byte,
	recipient []byte,
	callType vm.CallType,
	vmOutput *vmcommon.VMOutput,
) {
	esdtTransferTxData := function
	for _, arg := range arguments {
		esdtTransferTxData += "@" + hex.EncodeToString(arg)
	}
	outTransfer := vmcommon.OutputTransfer{
		Value:         big.NewInt(0),
		Data:          []byte(esdtTransferTxData),
		CallType:      callType,
		SenderAddress: sender,
	}

	if len(vmOutput.OutputAccounts) == 0 {
		vmOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
	}
	outAcc, ok := vmOutput.OutputAccounts[string(recipient)]
	if !ok {
		outAcc = &vmcommon.OutputAccount{
			Address:         recipient,
			OutputTransfers: make([]vmcommon.OutputTransfer, 0),
		}
	}
	outAcc.OutputTransfers = append(outAcc.OutputTransfers, outTransfer)
	vmOutput.OutputAccounts[string(recipient)] = outAcc
}

func (host *vmHost) checkFinalGasAfterExit() error {
	totalUsedPoints := host.Runtime().GetPointsUsed()
	if totalUsedPoints > host.Metering().GetGasForExecution() {
		log.Trace("checkFinalGasAfterExit", "failed")
		return arwen.ErrNotEnoughGas
	}

	log.Trace("checkFinalGasAfterExit", "ok")
	return nil
}

func (host *vmHost) callInitFunction() error {
	runtime := host.Runtime()
	init := runtime.GetInitFunction()
	if init == "" {
		return nil
	}

	err := runtime.CallFunction(init)
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	if err == nil {
		err = host.checkFinalGasAfterExit()
	}

	return err
}

func (host *vmHost) callSCMethod() error {
	runtime := host.Runtime()

	log.Trace("call SC method")

	// TODO host.verifyAllowedFunctionCall() performs some checks, but then the
	// function itself is changed by host.getFunctionByCallType(). Order must be
	// reversed, and `getFunctionByCallType()` must be decomposed into smaller functions.

	err := host.verifyAllowedFunctionCall()
	if err != nil {
		log.Trace("call SC method failed", "error", err, "src", "verifyAllowedFunctionCall")
		return err
	}

	callType := runtime.GetVMInput().CallType
	function, err := host.getFunctionByCallType(callType)
	if err != nil {
		if callType == vm.AsynchronousCallBack && errors.Is(err, arwen.ErrNilCallbackFunction) {
			err = host.processCallbackStack()
			if err != nil {
				log.Trace("call SC method failed", "error", err, "src", "processCallbackStack")
			}

			return err
		}
		log.Trace("call SC method failed", "error", err, "src", "getFunctionByCallType")
		return err
	}

	err = runtime.CallFunction(function)
	if err != nil {
		err = host.handleBreakpointIfAny(err)
		log.Trace("breakpoint detected and handled", "err", err)
	}
	if err == nil {
		err = host.checkFinalGasAfterExit()
	}
	if err != nil {
		log.Trace("call SC method failed", "error", err, "src", "sc function")
		return err
	}

	switch callType {
	case vm.AsynchronousCall:
		_, paiErr := host.processAsyncInfo(runtime.GetAsyncContextInfo())
		if paiErr != nil {
			log.Trace("call SC method failed", "error", paiErr)
			return paiErr
		}
		err = host.sendCallbackToCurrentCaller()
	case vm.AsynchronousCallBack:
		err = host.processCallbackStack()
	default:
		_, err = host.processAsyncInfo(runtime.GetAsyncContextInfo())
	}

	if err != nil {
		log.Trace("call SC method failed", "error", err, "src", "async post-process")
	}

	return err
}

func (host *vmHost) verifyAllowedFunctionCall() error {
	runtime := host.Runtime()
	functionName := runtime.Function()

	isInit := functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
	if isInit {
		return arwen.ErrInitFuncCalledInRun
	}

	isCallBack := functionName == arwen.CallbackFunctionName
	isInAsyncCallBack := runtime.GetVMInput().CallType == vm.AsynchronousCallBack
	if isCallBack && !isInAsyncCallBack {
		return arwen.ErrCallBackFuncCalledInRun
	}

	return nil
}

func (host *vmHost) isSCExecutionAfterBuiltInFunc(
	vmInput *vmcommon.ContractCallInput,
	vmOutput *vmcommon.VMOutput,
) (*vmcommon.ContractCallInput, error) {
	if vmOutput.ReturnCode != vmcommon.Ok {
		return nil, nil
	}
	if vmInput.ReturnCallAfterError && vmInput.CallType != vm.AsynchronousCallBack {
		return nil, nil
	}

	parsedTransfer, err := host.esdtTransferParser.ParseESDTTransfers(vmInput.CallerAddr, vmInput.RecipientAddr, vmInput.Function, vmInput.Arguments)
	if err != nil {
		return nil, nil
	}

	if !host.AreInSameShard(vmInput.CallerAddr, parsedTransfer.RcvAddr) {
		return nil, nil
	}
	if !host.Blockchain().IsSmartContract(parsedTransfer.RcvAddr) {
		return nil, nil
	}

	outAcc, ok := vmOutput.OutputAccounts[string(parsedTransfer.RcvAddr)]
	if !ok {
		return nil, nil
	}
	if len(outAcc.OutputTransfers) != 1 {
		return nil, nil
	}

	callType := vmInput.CallType
	scCallOutTransfer := outAcc.OutputTransfers[0]

	argParser := parsers.NewCallArgsParser()
	function, arguments, err := argParser.ParseData(string(scCallOutTransfer.Data))
	if err != nil {
		return nil, err
	}

	newVMInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     vmInput.CallerAddr,
			Arguments:      arguments,
			CallValue:      big.NewInt(0),
			CallType:       callType,
			GasPrice:       vmInput.GasPrice,
			GasProvided:    scCallOutTransfer.GasLimit,
			GasLocked:      scCallOutTransfer.GasLocked,
			OriginalTxHash: vmInput.OriginalTxHash,
			CurrentTxHash:  vmInput.CurrentTxHash,
		},
		RecipientAddr:     parsedTransfer.RcvAddr,
		Function:          function,
		AllowInitFunction: false,
	}

	newVMInput.ESDTTransfers = parsedTransfer.ESDTTransfers

	return newVMInput, nil
}
