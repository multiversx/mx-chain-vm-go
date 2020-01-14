package context

import (
	"fmt"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/context/subcontexts"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmContext) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	host.InitState()

	blockchain := host.Blockchain()
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return nil, err
	}

	runtime.SetVMInput(&input.VMInput)
	runtime.SetSCAddress(address)
	output.AddTxValueToAccount(address, input.CallValue)

	gasForDeployment, err := metering.DeductInitialGasForDirectDeployment(input)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(vmcommon.OutOfGas, err.Error()), nil
	}
	runtime.GetVMInput().GasProvided = gasForDeployment

	err = runtime.CreateWasmerInstance(input.ContractCode)

	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		return output.CreateVMOutputInCaseOfError(vmcommon.ContractInvalid, err.Error()), nil
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextId(idContext)
	defer func() {
		runtime.CleanInstance()
		arwen.RemoveHostContext(idContext)
	}()

	result, err := host.callInitFunction()
	if err != nil {
		fmt.Println("arwen.callInitFunction() error:", err.Error())
		return output.CreateVMOutputInCaseOfError(vmcommon.FunctionWrongSignature, err.Error()), nil
	}

	output.DeployCode(address, input.ContractCode)
	vmOutput := output.CreateVMOutput(result)

	return vmOutput, err
}

func (host *vmContext) doRunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.InitState()
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()
	blockchain := host.Blockchain()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)

	contract, err := blockchain.GetCode(runtime.GetSCAddress())
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return output.CreateVMOutputInCaseOfError(vmcommon.ContractInvalid, err.Error()), nil
	}

	gasForExecution, err := metering.DeductInitialGasForExecution(input, contract)
	if err != nil {
		return output.CreateVMOutputInCaseOfError(vmcommon.OutOfGas, err.Error()), nil
	}
	runtime.GetVMInput().GasProvided = gasForExecution

	err = runtime.CreateWasmerInstance(contract)
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return output.CreateVMOutputInCaseOfError(vmcommon.ContractInvalid, err.Error()), nil
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextId(idContext)
	defer func() {
		runtime.CleanInstance()
		arwen.RemoveHostContext(idContext)
	}()

	result, returnCode, err := host.callSCMethod()
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return output.CreateVMOutputInCaseOfError(returnCode, err.Error()), nil
	}

	if returnCode != vmcommon.Ok {
		return output.CreateVMOutputInCaseOfError(returnCode, output.ReturnMessage()), nil
	}

	vmOutput := output.CreateVMOutput(result)

	return vmOutput, nil
}

func (host *vmContext) ExecuteOnDestContext(input *vmcommon.ContractCallInput) error {
	host.PushState()

	var err error
	defer func() {
		popErr := host.PopState()
		if popErr != nil {
			err = popErr
		}
	}()

	host.InitState()

	host.Runtime().InitStateFromContractCallInput(input)
	err = host.execute(input)

	return err
}

func (host *vmContext) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	runtime := host.Runtime()
	runtime.PushState()

	var err error
	defer func() {
		popErr := runtime.PopState()
		if popErr != nil {
			err = popErr
		}
	}()

	runtime.InitStateFromContractCallInput(input)
	err = host.execute(input)

	return err
}

func (host *vmContext) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

func (host *vmContext) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()
	metering := host.Metering()
	output := host.Output()

	if runtime.ReadOnly() {
		return nil, ErrInvalidCallOnReadOnlyMode
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

	gasForDeployment, err := metering.DeductInitialGasForIndirectDeployment(input)
	if err != nil {
		output.Transfer(input.CallerAddr, address, 0, input.CallValue, nil)
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	defer func() {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
	}()

	err = runtime.CreateWasmerInstanceWithGasLimit(input.ContractCode, gasForDeployment)
	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		// TODO use all gas here?
		output.Transfer(input.CallerAddr, address, 0, input.CallValue, nil)
		return nil, err
	}

	runtime.SetInstanceContextId(idContext)

	result, err := host.callInitFunction()
	if err != nil {
		output.Transfer(input.CallerAddr, address, 0, input.CallValue, nil)
		return nil, err
	}

	output.DeployCode(address, input.ContractCode)
	output.FinishValue(result)

	totalGasConsumed = input.GasProvided - gasForDeployment - runtime.GetPointsUsed()

	return address, nil
}

func (host *vmContext) execute(input *vmcommon.ContractCallInput) error {
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

	gasForExecution, err := metering.DeductInitialGasForExecution(input, contract)
	if err != nil {
		return err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	defer func() {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
	}()

	err = runtime.CreateWasmerInstanceWithGasLimit(contract, gasForExecution)
	if err != nil {
		metering.UseGas(input.GasProvided)
		return err
	}

	runtime.SetInstanceContextId(idContext)

	if host.isInitFunctionBeingCalled() {
		return ErrInitFuncCalledInRun
	}

	// TODO replace with callSCMethod()?
	exports := runtime.GetInstanceExports()
	functionName := runtime.Function()
	function, ok := exports[functionName]
	if !ok {
		return subcontexts.ErrFuncNotFound
	}

	result, err := function()
	if err != nil {
		return ErrFunctionRunError
	}

	if output.ReturnCode() != vmcommon.Ok {
		return ErrReturnCodeNotOk
	}

	convertedResult := arwen.ConvertReturnValue(result)
	output.Finish(convertedResult.Bytes())

	totalGasConsumed = input.GasProvided - gasForExecution - runtime.GetPointsUsed()

	return nil
}

func (host *vmContext) EthereumCallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmContext) callInitFunction() (wasmer.Value, error) {
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

func (host *vmContext) callSCMethod() (wasmer.Value, vmcommon.ReturnCode, error) {
	if host.isInitFunctionBeingCalled() {
		return wasmer.Void(), vmcommon.UserError, ErrInitFuncCalledInRun
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
			err = host.handleBreakpoint(breakpointValue, result, err)
		}
	}

	if err != nil {
		strError, _ := wasmer.GetLastError()
		fmt.Println("wasmer Error", strError)
		return wasmer.Void(), vmcommon.ExecutionFailed, err
	}

	return result, output.ReturnCode(), nil
}

// The first four bytes is the method selector. The rest of the input data are method arguments in chunks of 32 bytes.
// The method selector is the kecccak256 hash of the method signature.
func (host *vmContext) createETHCallInput() []byte {
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
