package host

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/contexts"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-go/core"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, blockchain, _, output, runtime, storage := host.GetContexts()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)

	output.AddTxValueToAccount(address, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err := host.performCodeDeployment(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}
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

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(input.ContractCode, vmInput.GasProvided, true)
	if err != nil {
		log.Debug("performCodeDeployment/StartWasmerInstance", "err", err)
		return nil, arwen.ErrContractInvalid
	}

	err = host.callInitFunction()
	if err != nil {
		return nil, err
	}

	output.DeployCode(input)
	vmOutput := output.GetVMOutput()
	runtime.CleanWasmerInstance()
	return vmOutput, nil
}

// doRunSmartContractUpgrade upgrades a contract directly
func (host *vmHost) doRunSmartContractUpgrade(input *vmcommon.ContractCallInput) *vmcommon.VMOutput {
	host.InitState()
	defer host.Clean()

	_, _, _, output, runtime, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrInvalidUpgradeArguments)
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	vmOutput, err := host.performCodeDeployment(codeDeployInput)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}
	return vmOutput
}

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput) {
	host.InitState()
	defer host.Clean()

	_, _, metering, output, runtime, storage := host.GetContexts()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	contract, err := runtime.GetSCCode()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractNotFound)
	}

	metering.UnlockGasIfAsyncCallback()
	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrNotEnoughGas)
	}

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(contract, vmInput.GasProvided, false)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(arwen.ErrContractInvalid)
	}

	err = host.callSCMethod()
	if err != nil {
		return output.CreateVMOutputInCaseOfError(err)
	}

	vmOutput = output.GetVMOutput()
	runtime.CleanWasmerInstance()
	return
}

func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, asyncInfo *arwen.AsyncContextInfo, err error) {
	log.Trace("ExecuteOnDestContext", "function", input.Function)

	bigInt, _, _, output, runtime, storage := host.GetContexts()

	bigInt.PushState()
	bigInt.InitState()

	output.PushState()
	output.CensorVMOutput()
	output.ResetGas()

	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	storage.PushState()
	storage.SetAddress(runtime.GetSCAddress())

	gasUsed := uint64(0)
	defer func() {
		vmOutput = host.finishExecuteOnDestContext(gasUsed, err)
	}()

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err = output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue)
	if err != nil {
		return
	}

	gasUsed, err = host.execute(input)
	if err != nil {
		return
	}

	asyncInfo = runtime.GetAsyncContextInfo()
	_, err = host.processAsyncInfo(asyncInfo)
	return
}

func computeGasUsedByCurrentSC(
	gasUsed uint64,
	output arwen.OutputContext,
	executeErr error,
) (uint64, error) {
	if executeErr != nil {
		return 0, executeErr
	}

	vmOutput := output.GetVMOutput()
	if vmOutput.ReturnCode != vmcommon.Ok || gasUsed == 0 {
		return 0, nil
	}

	for _, outAcc := range vmOutput.OutputAccounts {
		accumulatedGasLimit := uint64(0)
		for _, outTransfer := range outAcc.OutputTransfers {
			accumulatedGasLimit += outTransfer.GasLimit
		}

		if gasUsed < outAcc.GasUsed+accumulatedGasLimit {
			return 0, arwen.ErrGasUsageError
		}

		gasUsed -= outAcc.GasUsed
		gasUsed -= accumulatedGasLimit
	}

	return gasUsed, nil
}

func (host *vmHost) finishExecuteOnDestContext(gasUsed uint64, executeErr error) *vmcommon.VMOutput {
	bigInt, _, _, output, runtime, storage := host.GetContexts()

	// Extract the VMOutput produced by the execution in isolation, before
	// restoring the contexts. This needs to be done before popping any state
	// stacks.
	gasUsedBySC, err := computeGasUsedByCurrentSC(gasUsed, output, executeErr)
	if err != nil {
		// Execution failed: restore contexts as if the execution didn't happen,
		// but first create a vmOutput to capture the error.
		vmOutput := output.CreateVMOutputInCaseOfError(err)

		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()
		storage.PopSetActiveState()

		return vmOutput
	}

	vmOutput := output.GetVMOutput()

	// Restore the previous context states, except Output, which will be merged
	// into the initial state (VMOutput), but only if it the child execution
	// returned vmcommon.Ok.
	bigInt.PopSetActiveState()
	runtime.PopSetActiveState()
	storage.PopSetActiveState()

	if vmOutput.ReturnCode == vmcommon.Ok {
		output.PopMergeActiveState()
		scAddress := string(runtime.GetSCAddress())
		accumulateGasUsedByContract(vmOutput, scAddress, gasUsedBySC)
	} else {
		output.PopSetActiveState()
	}

	return vmOutput
}

func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (asyncInfo *arwen.AsyncContextInfo, err error) {
	log.Trace("ExecuteOnSameContext", "function", input.Function)

	if host.IsBuiltinFunctionName(input.Function) {
		return nil, arwen.ErrBuiltinCallOnSameContextDisallowed
	}

	bigInt, _, _, output, runtime, _ := host.GetContexts()

	// Back up the states of the contexts (except Storage, which isn't affected
	// by ExecuteOnSameContext())
	bigInt.PushState()
	output.PushState()
	runtime.PushState()

	output.ResetGas()
	runtime.InitStateFromContractCallInput(input)

	gasUsed := uint64(0)
	defer func() {
		host.finishExecuteOnSameContext(gasUsed, err)
	}()

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err = output.TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue)
	if err != nil {
		return
	}

	gasUsed, err = host.execute(input)
	if err != nil {
		return
	}

	asyncInfo = runtime.GetAsyncContextInfo()

	return
}

func (host *vmHost) finishExecuteOnSameContext(gasUsed uint64, executeErr error) {
	bigInt, _, _, output, runtime, _ := host.GetContexts()

	gasUsedBySC, err := computeGasUsedByCurrentSC(gasUsed, output, executeErr)
	if output.ReturnCode() != vmcommon.Ok || err != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return
	}

	scAddress := string(runtime.GetSCAddress())
	// Execution successful: discard the backups made at the beginning and
	// resume from the new state.
	bigInt.PopDiscard()
	output.PopDiscard()
	runtime.PopSetActiveState()

	vmOutput := output.GetVMOutput()
	accumulateGasUsedByContract(vmOutput, scAddress, gasUsedBySC)
}

func accumulateGasUsedByContract(vmOutput *vmcommon.VMOutput, scAddress string, gasUsed uint64) {
	if _, ok := vmOutput.OutputAccounts[scAddress]; !ok {
		vmOutput.OutputAccounts[scAddress] = contexts.NewVMOutputAccount([]byte(scAddress))
	}
	vmOutput.OutputAccounts[scAddress].GasUsed += gasUsed
}

func (host *vmHost) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

func (host *vmHost) isBuiltinFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return host.IsBuiltinFunctionName(functionName)
}

func (host *vmHost) IsBuiltinFunctionName(functionName string) bool {
	_, ok := host.protocolBuiltinFunctions[functionName]
	return ok
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
	vmOutput, _, err := host.ExecuteOnDestContext(initCallInput)
	if err != nil {
		return
	}

	metering.RestoreGas(vmOutput.GasRemaining)
	blockchain.IncreaseNonce(input.CallerAddr)

	return
}

func (host *vmHost) checkUpgradePermission(vmInput *vmcommon.ContractCallInput) error {
	contract, err := host.blockChainHook.GetUserAccount(vmInput.RecipientAddr)
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

// executeUpgrade upgrades a contract indirectly (from another contract)
func (host *vmHost) executeUpgrade(input *vmcommon.ContractCallInput) (uint64, error) {
	_, _, metering, output, runtime, _ := host.GetContexts()

	initialGasProvided := input.GasProvided
	err := host.checkUpgradePermission(input)
	if err != nil {
		return 0, err
	}

	code, codeMetadata, err := runtime.ExtractCodeUpgradeFromArgs()
	if err != nil {
		return 0, arwen.ErrInvalidUpgradeArguments
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
		CodeDeployerAddress:  input.CallerAddr,
	}

	metering.UnlockGasIfAsyncCallback()
	err = metering.DeductInitialGasForDirectDeployment(codeDeployInput)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return 0, err
	}

	runtime.PushInstance()
	runtime.MustVerifyNextContractCode()

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(codeDeployInput.ContractCode, vmInput.GasProvided, true)
	if err != nil {
		log.Debug("performCodeDeployment/StartWasmerInstance", "err", err)
		runtime.PopInstance()
		return 0, arwen.ErrContractInvalid
	}

	err = host.callInitFunction()
	if err != nil {
		runtime.PopInstance()
		return 0, err
	}

	output.DeployCode(codeDeployInput)
	if output.ReturnCode() != vmcommon.Ok {
		runtime.PopInstance()
		return 0, arwen.ErrReturnCodeNotOk
	}

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	metering.RestoreGas(gasToRestoreToCaller)

	return initialGasProvided - gasToRestoreToCaller, nil
}

func (host *vmHost) executeSmartContractCall(
	input *vmcommon.ContractCallInput,
	metering arwen.MeteringContext,
	runtime arwen.RuntimeContext,
	output arwen.OutputContext,
	withInitialGasDeduct bool,
) (uint64, error) {
	// Use all gas initially, on the Wasmer instance of the caller
	// (runtime.PushInstance() is called later). In case of successful execution,
	// the unused gas will be restored.
	initialGasProvided := input.GasProvided
	metering.UseGas(initialGasProvided)

	if host.isInitFunctionBeingCalled() && !input.AllowInitFunction {
		return 0, arwen.ErrInitFuncCalledInRun
	}

	isUpgrade := input.Function == arwen.UpgradeFunctionName
	if isUpgrade {
		return host.executeUpgrade(input)
	}

	contract, err := runtime.GetSCCode()
	if err != nil {
		return 0, err
	}

	metering.UnlockGasIfAsyncCallback()
	if withInitialGasDeduct {
		err = metering.DeductInitialGasForExecution(contract)
		if err != nil {
			return 0, err
		}
	}

	runtime.PushInstance()

	gasForExecution := runtime.GetVMInput().GasProvided
	err = runtime.StartWasmerInstance(contract, gasForExecution, false)
	if err != nil {
		runtime.PopInstance()
		return 0, err
	}

	err = host.callSCMethodIndirect()
	if err != nil {
		runtime.PopInstance()
		return 0, err
	}

	if output.ReturnCode() != vmcommon.Ok {
		runtime.PopInstance()
		return 0, arwen.ErrReturnCodeNotOk
	}

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	metering.RestoreGas(gasToRestoreToCaller)

	return initialGasProvided - gasToRestoreToCaller, nil
}

func (host *vmHost) execute(input *vmcommon.ContractCallInput) (uint64, error) {
	_, _, metering, output, runtime, storage := host.GetContexts()

	if host.isBuiltinFunctionBeingCalled() {
		err := metering.DeductGasIfAsyncStep()
		if err != nil {
			return 0, err
		}

		newVMInput, err := host.callBuiltinFunction(input)
		if err != nil {
			return 0, err
		}

		if newVMInput != nil {
			runtime.InitStateFromContractCallInput(newVMInput)
			storage.SetAddress(runtime.GetSCAddress())
			return host.executeSmartContractCall(newVMInput, metering, runtime, output, false)
		}

		return 0, nil
	}

	return host.executeSmartContractCall(input, metering, runtime, output, true)
}

func (host *vmHost) callSCMethodIndirect() error {
	function, err := host.Runtime().GetFunctionToCall()
	if err != nil {
		if errors.Is(err, arwen.ErrNilCallbackFunction) {
			return nil
		}
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	return err
}

func (host *vmHost) callBuiltinFunction(input *vmcommon.ContractCallInput) (*vmcommon.ContractCallInput, error) {
	_, _, metering, output, _, _ := host.GetContexts()

	vmOutput, err := host.blockChainHook.ProcessBuiltInFunction(input)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return nil, err
	}

	gasConsumed := input.GasProvided - vmOutput.GasRemaining
	if vmOutput.GasRemaining < input.GasProvided {
		metering.UseGas(gasConsumed)
	}

	newVMInput, err := isSCExecutionAfterBuiltInFunc(input, vmOutput)
	if err != nil {
		return nil, err
	}

	output.AddToActiveState(vmOutput)

	return newVMInput, nil
}

func (host *vmHost) EthereumCallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmHost) callInitFunction() error {
	runtime := host.Runtime()
	init := runtime.GetInitFunction()
	if init == nil {
		return nil
	}

	_, err := init()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}

	return err
}

func (host *vmHost) callSCMethod() error {
	runtime := host.Runtime()

	err := host.verifyAllowedFunctionCall()
	if err != nil {
		return err
	}

	callType := runtime.GetVMInput().CallType
	function, err := host.getFunctionByCallType(callType)
	if err != nil {
		if callType == vmcommon.AsynchronousCallBack && errors.Is(err, arwen.ErrNilCallbackFunction) {
			return host.processCallbackStack()
		}
		return err
	}

	_, err = function()
	if err != nil {
		err = host.handleBreakpointIfAny(err)
	}
	if err != nil {
		return err
	}

	switch callType {
	case vmcommon.AsynchronousCall:
		pendingMap, paiErr := host.processAsyncInfo(runtime.GetAsyncContextInfo())
		if paiErr != nil {
			return paiErr
		}
		if len(pendingMap.AsyncContextMap) == 0 {
			err = host.sendCallbackToCurrentCaller()
		}
		break
	case vmcommon.AsynchronousCallBack:
		err = host.processCallbackStack()
		break
	default:
		_, err = host.processAsyncInfo(runtime.GetAsyncContextInfo())
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
	isInAsyncCallBack := runtime.GetVMInput().CallType == vmcommon.AsynchronousCallBack
	if isCallBack && !isInAsyncCallBack {
		return arwen.ErrCallBackFuncCalledInRun
	}

	return nil
}

// The first four bytes is the method selector. The rest of the input data are method arguments in chunks of 32 bytes.
// The method selector is the kecccak256 hash of the method signature.
func (host *vmHost) createETHCallInput() []byte {
	newInput := make([]byte, 0)

	function := host.Runtime().Function()
	if len(function) > 0 {
		hashOfFunction, err := host.cryptoHook.Keccak256([]byte(function))
		if err != nil {
			return nil
		}

		newInput = append(newInput, hashOfFunction[0:4]...)
	}

	for _, arg := range host.Runtime().Arguments() {
		paddedArg := make([]byte, arwen.ArgumentLenEth)
		copy(paddedArg[arwen.ArgumentLenEth-len(arg):], arg)
		newInput = append(newInput, paddedArg...)
	}

	return newInput
}

func isSCExecutionAfterBuiltInFunc(
	vmInput *vmcommon.ContractCallInput,
	vmOutput *vmcommon.VMOutput,
) (*vmcommon.ContractCallInput, error) {
	if vmOutput.ReturnCode != vmcommon.Ok {
		return nil, nil
	}

	if !core.IsSmartContractAddress(vmInput.RecipientAddr) {
		return nil, nil
	}

	outAcc, ok := vmOutput.OutputAccounts[string(vmInput.RecipientAddr)]
	if !ok {
		return nil, nil
	}
	if len(outAcc.OutputTransfers) != 1 {
		return nil, nil
	}

	callType := vmInput.CallType
	txData := prependCallbackToTxDataIfAsyncCall(outAcc.OutputTransfers[0].Data, callType)

	argParser := parsers.NewCallArgsParser()
	function, arguments, err := argParser.ParseData(txData)
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
			GasProvided:    vmOutput.GasRemaining,
			OriginalTxHash: vmInput.OriginalTxHash,
			CurrentTxHash:  vmInput.CurrentTxHash,
		},
		RecipientAddr:     vmInput.RecipientAddr,
		Function:          function,
		AllowInitFunction: false,
	}

	fillWithESDTValue(vmInput, newVMInput)

	return newVMInput, nil
}

func fillWithESDTValue(fullVMInput *vmcommon.ContractCallInput, newVMInput *vmcommon.ContractCallInput) {
	if fullVMInput.Function != core.BuiltInFunctionESDTTransfer {
		return
	}

	newVMInput.ESDTTokenName = fullVMInput.Arguments[0]
	newVMInput.ESDTValue = big.NewInt(0).SetBytes(fullVMInput.Arguments[1])
}

func prependCallbackToTxDataIfAsyncCall(txData []byte, callType vmcommon.CallType) string {
	if callType == vmcommon.AsynchronousCallBack {
		return string(append([]byte(arwen.CallbackFunctionName), txData...))
	}

	return string(txData)
}
