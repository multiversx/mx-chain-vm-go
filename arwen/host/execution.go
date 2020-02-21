package host

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput) {
	host.ClearStateStack()
	host.InitState()

	blockchain := host.Blockchain()
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	var returnCode vmcommon.ReturnCode
	var err error
	defer func() {
		if err != nil {
			var message string
			if err == arwen.ErrSignalError {
				message = output.ReturnMessage()
			} else {
				message = err.Error()
			}
			vmOutput = output.CreateVMOutputInCaseOfError(returnCode, message)
		}
	}()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		returnCode = vmcommon.ExecutionFailed
		return vmOutput
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)
	output.AddTxValueToAccount(address, input.CallValue)

	err = metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		returnCode = vmcommon.OutOfGas
		return vmOutput
	}

	vmInput := runtime.GetVMInput()
	err = runtime.CreateWasmerInstance(input.ContractCode, vmInput.GasProvided)
	if err != nil {
		returnCode = vmcommon.ContractInvalid
		return vmOutput
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		returnCode = vmcommon.ContractInvalid
		return vmOutput
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextId(idContext)
	defer func() {
		runtime.CleanInstance()
		arwen.RemoveHostContext(idContext)
	}()

	result, err := host.callInitFunction()
	if err != nil {
		returnCode = vmcommon.FunctionWrongSignature
		return vmOutput
	}

	output.DeployCode(address, input.ContractCode)
	output.FinishValue(result)
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

	var returnCode vmcommon.ReturnCode
	var err error
	defer func() {
		if err != nil {
			var message string
			if err == arwen.ErrSignalError {
				message = output.ReturnMessage()
			} else {
				message = err.Error()
			}
			vmOutput = output.CreateVMOutputInCaseOfError(returnCode, message)
		}
	}()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)

	contract, err := blockchain.GetCode(runtime.GetSCAddress())
	if err != nil {
		returnCode = vmcommon.ContractInvalid
		return vmOutput
	}

	vmInput := runtime.GetVMInput()

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		returnCode = vmcommon.OutOfGas
		return vmOutput
	}

	err = runtime.CreateWasmerInstance(contract, vmInput.GasProvided)
	if err != nil {
		returnCode = vmcommon.ContractInvalid
		return vmOutput
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextId(idContext)
	defer func() {
		runtime.CleanInstance()
		arwen.RemoveHostContext(idContext)
	}()

	result, returnCode, err := host.callSCMethod()

	if err != nil {
		return vmOutput
	}

	metering.UnlockGasIfAsyncStep()

	output.FinishValue(result)
	vmOutput = output.GetVMOutput()

	return vmOutput
}

func (host *vmHost) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.PushState()
	defer host.PopState()

	host.InitState()

	host.Runtime().InitStateFromContractCallInput(input)
	err := host.execute(input)
	if err != nil {
		return nil, err
	}

	vmOutput := host.Output().GetVMOutput()

	return vmOutput, nil
}

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

	if runtime.ReadOnly() {
		return nil, arwen.ErrInvalidCallOnReadOnlyMode
	}

	runtime.PushState()

	defer func() {
		runtime.PopState()
	}()

	runtime.SetVMInput(&input.VMInput)
	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return nil, err
	}

	output.Transfer(address, input.CallerAddr, 0, input.CallValue, nil)
	blockchain.IncreaseNonce(input.CallerAddr)
	runtime.SetSCAddress(address)

	totalGasConsumed := input.GasProvided
	defer func() {
		metering.UseGas(totalGasConsumed)
	}()

	err = metering.DeductInitialGasForIndirectDeployment(input)
	if err != nil {
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	defer func() {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
	}()

	gasForDeployment := runtime.GetVMInput().GasProvided
	err = runtime.CreateWasmerInstance(input.ContractCode, gasForDeployment)
	if err != nil {
		return nil, err
	}

	err = runtime.VerifyContractCode()
	if err != nil {
		return nil, err
	}

	runtime.SetInstanceContextId(idContext)

	result, err := host.callInitFunction()
	if err != nil {
		return nil, err
	}

	output.DeployCode(address, input.ContractCode)
	output.FinishValue(result)

	totalGasConsumed = input.GasProvided - gasForDeployment - runtime.GetPointsUsed()

	return address, nil
}

func (host *vmHost) execute(input *vmcommon.ContractCallInput) error {
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	contract, err := host.Blockchain().GetCode(runtime.GetSCAddress())
	if err != nil {
		return err
	}

	totalGasConsumed := input.GasProvided

	defer func() {
		metering.UseGas(totalGasConsumed)
	}()

	err = metering.DeductInitialGasForExecution(contract)
	if err != nil {
		return err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	defer func() {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
	}()

	gasForExecution := runtime.GetVMInput().GasProvided
	err = runtime.CreateWasmerInstance(contract, gasForExecution)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return err
	}

	runtime.SetInstanceContextId(idContext)

	if host.isInitFunctionBeingCalled() {
		return arwen.ErrInitFuncCalledInRun
	}

	// TODO replace with callSCMethod()?
	exports := runtime.GetInstanceExports()
	functionName := runtime.Function()
	function, ok := exports[functionName]
	if !ok {
		return arwen.ErrFuncNotFound
	}

	result, err := function()
	if err != nil {
		return arwen.ErrFunctionRunError
	}

	if output.ReturnCode() != vmcommon.Ok {
		return arwen.ErrReturnCodeNotOk
	}

	convertedResult := arwen.ConvertReturnValue(result)
	output.Finish(convertedResult)

	metering.UnlockGasIfAsyncStep()

	totalGasConsumed = input.GasProvided - gasForExecution - runtime.GetPointsUsed()

	return nil
}

func (host *vmHost) EthereumCallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmHost) callInitFunction() (wasmer.Value, error) {
	init := host.Runtime().GetInitFunction()
	if init != nil {
		result, err := init()
		if err != nil {
			return wasmer.Void(), err
		}
		return result, nil
	}
	return wasmer.Void(), nil
}

func (host *vmHost) callSCMethod() (wasmer.Value, vmcommon.ReturnCode, error) {
	if host.isInitFunctionBeingCalled() {
		return wasmer.Void(), vmcommon.UserError, arwen.ErrInitFuncCalledInRun
	}

	runtime := host.Runtime()
	output := host.Output()

	function, err := runtime.GetFunctionToCall()
	if err != nil {
		return wasmer.Void(), vmcommon.FunctionNotFound, err
	}

	result, err := function()
	if err != nil {
		breakpointValue := runtime.GetRuntimeBreakpointValue()
		if breakpointValue != arwen.BreakpointNone {
			err = host.handleBreakpoint(breakpointValue, result)
		}
	}

	if err != nil {
		result = wasmer.Void()
	}

	return result, output.ReturnCode(), err
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
