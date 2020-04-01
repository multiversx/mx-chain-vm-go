package host

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput) {
	host.ClearStateStack()
	host.InitState()

	blockchain := host.Blockchain()
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()
	storage := host.Storage()

	var err error
	defer func() {
		if err != nil {
			var message string
			if err == arwen.ErrSignalError {
				message = output.ReturnMessage()
			} else {
				message = err.Error()
			}
			vmOutput = output.CreateVMOutputInCaseOfError(output.ReturnCode(), message)
		}
	}()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		output.SetReturnCode(vmcommon.ExecutionFailed)
		return vmOutput
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)
	output.AddTxValueToAccount(address, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	err = metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return vmOutput
	}

	vmInput := runtime.GetVMInput()
	err = runtime.CreateWasmerInstance(input.ContractCode, vmInput.GasProvided)
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return vmOutput
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return vmOutput
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextID(idContext)
	defer func() {
		runtime.CleanInstance()
		arwen.RemoveHostContext(idContext)
	}()

	err = host.callInitFunction()
	if err != nil {
		output.SetReturnCode(vmcommon.FunctionWrongSignature)
		return vmOutput
	}

	output.DeployCode(address, input.ContractCode)
	vmOutput = output.GetVMOutput()

	return vmOutput
}

func (host *vmHost) doRunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput) {
	host.ClearStateStack()
	host.InitState()

	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	blockchain := host.Blockchain()
	storage := host.Storage()

	var err error
	defer func() {
		if err != nil {
			var message string
			if err == arwen.ErrSignalError {
				message = output.ReturnMessage()
			} else {
				message = err.Error()
			}
			vmOutput = output.CreateVMOutputInCaseOfError(output.ReturnCode(), message)
		}
	}()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())

	contract, err := blockchain.GetCode(runtime.GetSCAddress())
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return vmOutput
	}

	vmInput := runtime.GetVMInput()

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		output.SetReturnCode(vmcommon.OutOfGas)
		return vmOutput
	}

	err = runtime.CreateWasmerInstance(contract, vmInput.GasProvided)
	if err != nil {
		output.SetReturnCode(vmcommon.ContractInvalid)
		return vmOutput
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextID(idContext)
	defer func() {
		runtime.CleanInstance()
		arwen.RemoveHostContext(idContext)
	}()

	returnCode, err := host.callSCMethod()
	if err != nil {
		output.SetReturnCode(returnCode)
		return vmOutput
	}

	metering.UnlockGasIfAsyncStep()

	vmOutput = output.GetVMOutput()

	return vmOutput
}

/*
	initialize:
		bigint	push, InitState
		output	push, InitState (clone top of stack with censoring)
		runtime	push, InitStateFromContractCallInput
		storage	push, SetAddress

	success:
		bigint	popSetActive
		output	popMergeActive
		runtime	popSetActive
		storage	popSetActive

	fail:
		bigint	popSetActive
		output	popSetActive
		runtime	popSetActive
		storage	popSetActive
*/
func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.PushState()
	defer host.PopState()

	host.InitState()

	host.Runtime().InitStateFromContractCallInput(input)
	host.Storage().SetAddress(host.Runtime().GetSCAddress())

	err := host.execute(input)
	if err != nil {
		return nil, err
	}

	vmOutput := host.Output().GetVMOutput()

	return vmOutput, nil
}

/*
	initialize:
		bigint	push
		output	push
		runtime push, InitStateFromContractCallInput
		storage -

	success:
		bigint	popDiscard
		output	popDiscard
		runtime	popSetActive
		storage	-

	fail:
		bigint	popSetActive
		output	popSetActive
		runtime	popSetActive
		storage	-
*/
func (host *vmHost) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	runtime := host.Runtime()
	runtime.PushState()
	defer runtime.PopState()

	runtime.InitStateFromContractCallInput(input)
	err := host.execute(input)

	return err
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
		runtime.PopState()
		return nil, err
	}

	err = output.Transfer(address, input.CallerAddr, 0, input.CallValue, nil)
	if err != nil {
		runtime.PopState()
		return nil, err
	}

	blockchain.IncreaseNonce(input.CallerAddr)
	runtime.SetSCAddress(address)

	err = metering.DeductInitialGasForIndirectDeployment(input)
	if err != nil {
		runtime.PopState()
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	gasForDeployment := runtime.GetVMInput().GasProvided
	err = runtime.CreateWasmerInstance(input.ContractCode, gasForDeployment)
	if err != nil {
		runtime.PopInstance()
		runtime.PopState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		runtime.PopInstance()
		runtime.PopState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	runtime.SetInstanceContextID(idContext)

	err = host.callInitFunction()
	if err != nil {
		runtime.PopInstance()
		runtime.PopState()
		arwen.RemoveHostContext(idContext)
		return nil, err
	}

	output.DeployCode(address, input.ContractCode)

	gasToRestoreToCaller := metering.GasLeft()

	runtime.PopInstance()
	runtime.PopState()
	arwen.RemoveHostContext(idContext)

	metering.RestoreGas(gasToRestoreToCaller)
	return address, nil
}

func (host *vmHost) execute(input *vmcommon.ContractCallInput) error {
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	// Use all gas initially. In case of successful execution, the unused gas
	// will be restored.
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

	function, err := runtime.GetFunctionToCall()
	if err != nil {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return err
	}

	_, err = function()
	if err != nil {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
		return arwen.ErrFunctionRunError
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

func (host *vmHost) callSCMethod() (vmcommon.ReturnCode, error) {
	runtime := host.Runtime()

	err := host.verifyAllowedFunctionCall()
	if err != nil {
		return vmcommon.UserError, err
	}

	function, err := runtime.GetFunctionToCall()
	if err != nil {
		return vmcommon.FunctionNotFound, err
	}

	_, err = function()
	if err != nil {
		breakpointValue := runtime.GetRuntimeBreakpointValue()
		if breakpointValue != arwen.BreakpointNone {
			err = host.handleBreakpoint(breakpointValue)
		}
	}

	if err != nil {
		var returnCode vmcommon.ReturnCode
		switch err {
		case arwen.ErrSignalError:
			returnCode = vmcommon.UserError
		case arwen.ErrNotEnoughGas:
			returnCode = vmcommon.OutOfGas
		default:
			returnCode = vmcommon.ExecutionFailed
		}

		return returnCode, err
	}

	return vmcommon.Ok, nil
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
