package host

import (
	"errors"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput) {
	host.ClearContextStateStack()
	host.InitState()

	blockchain := host.Blockchain()
	runtime := host.Runtime()
	output := host.Output()
	storage := host.Storage()

	var err error
	defer func() {
		vmOutput = host.onExitDirectCreateOrCall(err, vmOutput)
	}()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		output.SetReturnCode(vmcommon.ExecutionFailed)
		return
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)

	output.AddTxValueToAccount(address, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
	}

	vmOutput, err = host.performCodeDeploy(codeDeployInput)
	return
}

func (host *vmHost) onExitDirectCreateOrCall(err error, vmOutput *vmcommon.VMOutput) *vmcommon.VMOutput {
	host.Runtime().CleanInstance()
	arwen.RemoveAllHostContexts()

	return host.overrideVMOutputIfError(err, vmOutput)
}

func (host *vmHost) overrideVMOutputIfError(err error, vmOutput *vmcommon.VMOutput) *vmcommon.VMOutput {
	if err == nil {
		return vmOutput
	}

	output := host.Output()

	var message string
	if err == arwen.ErrSignalError {
		message = output.ReturnMessage()
	} else {
		message = err.Error()
	}

	return output.CreateVMOutputInCaseOfError(output.ReturnCode(), message)
}

func (host *vmHost) performCodeDeploy(input arwen.CodeDeployInput) (*vmcommon.VMOutput, error) {
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	err := metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return nil, err
	}

	vmInput := runtime.GetVMInput()
	err = runtime.CreateWasmerInstance(input.ContractCode, vmInput.GasProvided)
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return nil, err
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextID(idContext)

	err = host.callInitFunction()
	returnCode := host.resolveReturnCodeFromError(err)
	output.SetReturnCode(returnCode)
	if err != nil {
		return nil, err
	}

	output.DeployCode(input)
	vmOutput := output.GetVMOutput()
	return vmOutput, nil
}

func (host *vmHost) doRunSmartContractUpgrade(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput) {
	host.ClearContextStateStack()
	host.InitState()

	runtime := host.Runtime()
	output := host.Output()
	storage := host.Storage()

	var err error
	defer func() {
		vmOutput = host.onExitDirectCreateOrCall(err, vmOutput)
	}()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	code, codeMetadata, err := runtime.GetCodeUpgradeFromArgs()
	if err != nil {
		output.SetReturnCode(vmcommon.UpgradeFailed)
		return
	}

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
		ContractAddress:      input.RecipientAddr,
	}

	vmOutput, err = host.performCodeDeploy(codeDeployInput)
	return
}

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput) {
	host.ClearContextStateStack()
	host.InitState()

	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	blockchain := host.Blockchain()
	storage := host.Storage()

	var err error
	defer func() {
		vmOutput = host.onExitDirectCreateOrCall(err, vmOutput)
	}()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	contract, err := blockchain.GetCode(runtime.GetSCAddress())
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return
	}

	vmInput := runtime.GetVMInput()
	err = runtime.CreateWasmerInstance(contract, vmInput.GasProvided)
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextID(idContext)

	err = host.callSCMethod()
	returnCode := host.resolveReturnCodeFromError(err)
	output.SetReturnCode(returnCode)
	if err != nil {
		return
	}

	metering.UnlockGasIfAsyncStep()

	vmOutput = output.GetVMOutput()
	return
}

func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	bigInt := host.BigInt()
	output := host.Output()
	runtime := host.Runtime()
	storage := host.Storage()

	bigInt.PushState()
	bigInt.InitState()

	output.PushState()
	output.CensorVMOutput()

	runtime.PushState()
	runtime.InitStateFromContractCallInput(input)

	storage.PushState()
	storage.SetAddress(host.Runtime().GetSCAddress())

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err := output.Transfer(input.RecipientAddr, input.CallerAddr, 0, input.CallValue, nil)
	if err != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return nil, err
	}

	err = host.execute(input)

	// Extract the VMOutput produced by the execution in isolation, before
	// restoring the contexts. This needs to be done before popping any state
	// stacks.
	vmOutput := host.Output().GetVMOutput()

	if err != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return nil, err
	}

	// Execution successful: restore the previous context states, except Output,
	// which will merge the current state (VMOutput) with the initial state.
	bigInt.PopSetActiveState()
	output.PopMergeActiveState()
	runtime.PopSetActiveState()
	storage.PopSetActiveState()

	return vmOutput, nil
}

func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	bigInt := host.BigInt()
	output := host.Output()
	runtime := host.Runtime()

	// Back up the states of the contexts (except Storage, which isn't affected
	// by ExecuteOnSameContext())
	bigInt.PushState()
	output.PushState()
	runtime.PushState()

	runtime.InitStateFromContractCallInput(input)

	// Perform a value transfer to the called SC. If the execution fails, this
	// transfer will not persist.
	err := output.Transfer(input.RecipientAddr, input.CallerAddr, 0, input.CallValue, nil)
	if err != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return err
	}

	err = host.execute(input)
	if err != nil {
		// Execution failed: restore contexts as if the execution didn't happen.
		bigInt.PopSetActiveState()
		output.PopSetActiveState()
		runtime.PopSetActiveState()

		return err
	}

	// Execution successful: discard the backups made at the beginning and
	// resume from the new state.
	bigInt.PopDiscard()
	output.PopDiscard()
	runtime.PopSetActiveState()

	return nil
}

func (host *vmHost) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

func (host *vmHost) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()
	metering := host.Metering()
	output := host.Output()

	// Use all gas initially. In case of successful deployment, the unused gas
	// will be restored.
	initialGasProvided := input.GasProvided
	metering.UseGas(initialGasProvided)

	if runtime.ReadOnly() {
		return nil, arwen.ErrInvalidCallOnReadOnlyMode
	}

	runtime.PushState()

	runtime.SetVMInput(&input.VMInput)
	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		runtime.PopSetActiveState()
		return nil, err
	}

	err = output.Transfer(address, input.CallerAddr, 0, input.CallValue, nil)
	if err != nil {
		runtime.PopSetActiveState()
		return nil, err
	}

	blockchain.IncreaseNonce(input.CallerAddr)
	runtime.SetSCAddress(address)

	codeDeployInput := arwen.CodeDeployInput{
		ContractCode:         input.ContractCode,
		ContractCodeMetadata: input.ContractCodeMetadata,
		ContractAddress:      address,
	}

	err = metering.DeductInitialGasForIndirectDeployment(codeDeployInput)
	if err != nil {
		runtime.PopSetActiveState()
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	gasForDeployment := runtime.GetVMInput().GasProvided
	err = runtime.CreateWasmerInstance(input.ContractCode, gasForDeployment)
	if err != nil {
		runtime.PopInstance()
		runtime.PopSetActiveState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		runtime.PopInstance()
		runtime.PopSetActiveState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	runtime.SetInstanceContextID(idContext)

	err = host.callInitFunction()
	if err != nil {
		runtime.PopInstance()
		runtime.PopSetActiveState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	output.DeployCode(codeDeployInput)

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	runtime.PopSetActiveState()
	arwen.RemoveHostContext(idContext)

	metering.RestoreGas(gasToRestoreToCaller)
	return address, nil
}

// TODO: Add support for indirect smart contract upgrades.
func (host *vmHost) execute(input *vmcommon.ContractCallInput) error {
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	// Use all gas initially, on the Wasmer instance of the caller
	// (runtime.PushInstance() is called later). In case of successful execution,
	// the unused gas will be restored.
	initialGasProvided := input.GasProvided
	metering.UseGas(initialGasProvided)

	if host.isInitFunctionBeingCalled() {
		return arwen.ErrInitFuncCalledInRun
	}

	contract, err := host.Blockchain().GetCode(runtime.GetSCAddress())
	if err != nil {
		return err
	}

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return err
	}

	idContext := arwen.AddHostContext(host)

	runtime.PushInstance()

	gasForExecution := runtime.GetVMInput().GasProvided
	err = runtime.CreateWasmerInstance(contract, gasForExecution)
	if err != nil {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return err
	}

	runtime.SetInstanceContextID(idContext)

	err = host.callSCMethodIndirect()
	returnCode := host.resolveReturnCodeFromError(err)
	output.SetReturnCode(returnCode)
	if err != nil {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return err
	}

	if output.ReturnCode() != vmcommon.Ok {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return arwen.ErrReturnCodeNotOk
	}

	metering.UnlockGasIfAsyncStep()

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	metering.RestoreGas(gasToRestoreToCaller)
	arwen.RemoveHostContext(idContext)

	return nil
}

func (host *vmHost) callSCMethodIndirect() error {
	function, err := host.Runtime().GetFunctionToCall()
	if err != nil {
		return err
	}

	_, err = function()
	if err != nil {
		return arwen.ErrFunctionRunError
	}

	return nil
}

func (host *vmHost) EthereumCallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmHost) callInitFunction() error {
	init := host.Runtime().GetInitFunction()
	if init != nil {
		_, err := init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (host *vmHost) callSCMethod() error {
	runtime := host.Runtime()

	err := host.verifyAllowedFunctionCall()
	if err != nil {
		return err
	}

	function, err := runtime.GetFunctionToCall()
	if err != nil {
		return err
	}

	_, err = function()
	if err != nil {
		breakpointValue := runtime.GetRuntimeBreakpointValue()
		if breakpointValue != arwen.BreakpointNone {
			err = host.handleBreakpoint(breakpointValue)
		}
	}

	return err
}

func (host *vmHost) resolveReturnCodeFromError(err error) vmcommon.ReturnCode {
	if err != nil {
		if errors.Is(err, arwen.ErrSignalError) {
			return vmcommon.UserError
		}
		if err == arwen.ErrFuncNotFound {
			return vmcommon.FunctionNotFound
		}
		if err == arwen.ErrFunctionNonvoidSignature {
			return vmcommon.FunctionWrongSignature
		}
		if errors.Is(err, arwen.ErrInvalidFunction) {
			return vmcommon.UserError
		}
		if errors.Is(err, arwen.ErrNotEnoughGas) {
			return vmcommon.OutOfGas
		}

		return vmcommon.ExecutionFailed
	}

	return vmcommon.Ok
}

func (host *vmHost) verifyAllowedFunctionCall() error {
	runtime := host.Runtime()
	functionName := runtime.Function()

	isInit := functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
	if isInit {
		return arwen.ErrInitFuncCalledInRun
	}

	isCallBack := functionName == arwen.CallBackFunctionName
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
