package context

import (
	"math/big"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

func (host *vmContext) handleAsyncCallBreakpoint(result wasmer.Value, err error) (*vmcommon.VMOutput, error) {
	convertedResult := arwen.ConvertReturnValue(result)
	senderVMOutput := host.createVMOutput(convertedResult.Bytes())
	intermediaryVMOutput := senderVMOutput

	// If SC code is not found, it means this is either a cross-shard call or a wrong call.
	dest := host.asyncCallDest
	calledSCCode, err := host.GetCode(dest)
	if err != nil || len(calledSCCode) == 0 {
		// TODO detect same Shard call - this makes the empty calledSCCode an error
		// if intraShard {
		//   vmOutputWithError := createVMOutputInCaseOfBreakpointError(err)
		//   return mergeTwoVMOutputs(intermediaryVMOutput, vmOutputWithError)
		// }
		return intermediaryVMOutput, nil
	}

	// Start calling the destination SC, synchronously.
	destinationCallInput, err := host.createDestinationContractCallInput()
	if err != nil {
		vmOutputWithError := host.createVMOutputInCaseOfBreakpointError(err)
		return mergeTwoVMOutputs(intermediaryVMOutput, vmOutputWithError), nil
	}

	destinationVMOutput, err := host.executeOnNewContextAndGetVMOutput(destinationCallInput)
	// TODO pass error to the SC callback (append as argument)
	if err != nil {
		vmOutputWithError := host.createVMOutputInCaseOfBreakpointError(err)
		return mergeTwoVMOutputs(intermediaryVMOutput, vmOutputWithError), nil
	}

	intermediaryVMOutput = mergeTwoVMOutputs(intermediaryVMOutput, destinationVMOutput)

	callbackCallInput, err := host.createCallbackContractCallInput(destinationVMOutput)
	if err != nil {
		vmOutputWithError := host.createVMOutputInCaseOfBreakpointError(err)
		return mergeTwoVMOutputs(intermediaryVMOutput, vmOutputWithError), nil
	}

	callbackVMOutput, err := host.executeOnNewContextAndGetVMOutput(callbackCallInput)
	if err != nil {
		vmOutputWithError := host.createVMOutputInCaseOfBreakpointError(err)
		return mergeTwoVMOutputs(intermediaryVMOutput, vmOutputWithError), nil
	}
	finalVMOutput := mergeTwoVMOutputs(intermediaryVMOutput, callbackVMOutput)

	return finalVMOutput, nil
}

func (host *vmContext) createDestinationContractCallInput() (*vmcommon.ContractCallInput, error) {
	sender := host.GetSCAddress()
	dest := host.asyncCallDest
	valueBytes := host.asyncCallValueBytes
	data := host.asyncCallData
	gasLimit := host.asyncCallGasLimit

	err := host.argParser.ParseData(string(data))
	if err != nil {
		return nil, err
	}

	function, err := host.argParser.GetFunction()
	if err != nil {
		return nil, err
	}

	arguments, err := host.argParser.GetArguments()
	if err != nil {
		return nil, err
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0).SetBytes(valueBytes),
			GasPrice:    0,
			GasProvided: gasLimit,
		},
		RecipientAddr: dest,
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmContext) createCallbackContractCallInput(destinationVMOutput *vmcommon.VMOutput) (*vmcommon.ContractCallInput, error) {
	arguments := destinationVMOutput.ReturnData
	gasLimit := destinationVMOutput.GasRemaining
	function := "callBack"

	// Calculate what length would the Data field have, were it of the
	// form "callback@arg1@arg4...
	dataLength := len(function) + 1
	for i, element := range arguments {
		if len(element) == 0 {
			continue
		}
		if i != 0 && dataLength > 0 {
			dataLength += 1
		}
		dataLength += len(element)
	}

	// TODO define gas cost specific to async calls - asyncCall should cost more
	// than ExecuteOnDestContext, especially cross-shard async calls
	gasToUse := host.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	gasToUse += host.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(dataLength)
	gasLimit -= gasToUse

	sender := host.asyncCallDest
	dest := host.GetSCAddress()

	// Return to the sender SC, calling its callback() method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0),
			GasPrice:    0,
			GasProvided: gasLimit,
		},
		RecipientAddr: dest,
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmContext) executeOnNewContextAndGetVMOutput(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	currVmInput := host.vmInput
	currScAddress := host.scAddress
	currCallFunction := host.callFunction

	currContext := host.copyToNewContext()

	defer func() {
		// Restore the original host context
		host.copyFromContext(currContext)
		host.vmInput = currVmInput
		host.scAddress = currScAddress
		host.callFunction = currCallFunction
	}()

	host.initInternalValues()

	host.vmInput = input.VMInput
	host.scAddress = input.RecipientAddr
	host.callFunction = input.Function

	// TODO handle out-of-gas
	err := host.execute(input)

	var vmOutput *vmcommon.VMOutput
	if host.reachedBreakpoint(err) {
		vmOutput, err = host.handleBreakpoint(wasmer.I32(-1), err)
	}

	if err != nil {
		if err == ErrUnhandledRuntimeBreakpoint {
			return host.CreateVMOutputInCaseOfErrorWithMessage(vmcommon.ExecutionFailed, ErrUnhandledRuntimeBreakpoint.Error()), nil
		}

		strError, _ := wasmer.GetLastError()
		return host.CreateVMOutputInCaseOfErrorWithMessage(vmcommon.ExecutionFailed, strError), nil
	}

	// TODO this will override the VMOutput created by breakpoint handlers, if
	// the handlers didn't change the value of host.returnCode to vmcommon.Ok
	if host.returnCode != vmcommon.Ok {
		return host.createVMOutputInCaseOfError(host.returnCode), nil
	}

	vmOutput = host.createVMOutput(make([]byte, 0))
	return vmOutput, nil
}

func mergeVMOutputs(sender *vmcommon.VMOutput, destination *vmcommon.VMOutput, callback *vmcommon.VMOutput) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{}
	vmOutput = mergeTwoVMOutputs(vmOutput, sender)
	vmOutput = mergeTwoVMOutputs(vmOutput, destination)
	vmOutput = mergeTwoVMOutputs(vmOutput, callback)
	return vmOutput
}

func mergeTwoVMOutputs(leftVMOutput *vmcommon.VMOutput, rightVMOutput *vmcommon.VMOutput) *vmcommon.VMOutput {
	mergedVMOutput := &vmcommon.VMOutput{}

	leftOutputAccounts := convertAccountsSliceToAccountsMap(leftVMOutput.OutputAccounts)
	rightOutputAccounts := convertAccountsSliceToAccountsMap(rightVMOutput.OutputAccounts)
	mergedOutputAccounts := mergeOutputAccountMaps(leftOutputAccounts, rightOutputAccounts)
	mergedVMOutput.OutputAccounts = convertAccountsMapToAccountsSlice(mergedOutputAccounts)

	mergedVMOutput.Logs = append(mergedVMOutput.Logs, leftVMOutput.Logs...)
	mergedVMOutput.Logs = append(mergedVMOutput.Logs, rightVMOutput.Logs...)

	mergedVMOutput.ReturnData = append(mergedVMOutput.ReturnData, leftVMOutput.ReturnData...)
	mergedVMOutput.ReturnData = append(mergedVMOutput.ReturnData, rightVMOutput.ReturnData...)

	// TODO merge DeletedAccounts and TouchedAccounts as well?

	mergedVMOutput.GasRemaining = rightVMOutput.GasRemaining
	mergedVMOutput.GasRefund = rightVMOutput.GasRefund
	mergedVMOutput.ReturnCode = rightVMOutput.ReturnCode
	mergedVMOutput.ReturnMessage = rightVMOutput.ReturnMessage

	return mergedVMOutput
}

func convertAccountsSliceToAccountsMap(outputAccountsSlice []*vmcommon.OutputAccount) map[string]*vmcommon.OutputAccount {
	outputAccountsMap := make(map[string]*vmcommon.OutputAccount)

	for _, account := range outputAccountsSlice {
		address := string(account.Address)
		outputAccountsMap[address] = account
	}

	return outputAccountsMap
}

func convertAccountsMapToAccountsSlice(outputAccountsMap map[string]*vmcommon.OutputAccount) []*vmcommon.OutputAccount {
	outputAccountsSlice := make([]*vmcommon.OutputAccount, len(outputAccountsMap))
	i := 0
	for _, account := range outputAccountsMap {
		outputAccountsSlice[i] = account
		i++
	}

	return outputAccountsSlice
}

func mergeOutputAccountMaps(leftMap map[string]*vmcommon.OutputAccount, rightMap map[string]*vmcommon.OutputAccount) map[string]*vmcommon.OutputAccount {
	mergedAccountsMap := make(map[string]*vmcommon.OutputAccount)

	for addr, account := range leftMap {
		mergedAccountsMap[addr] = account
	}

	for addr, account := range rightMap {
		if _, ok := mergedAccountsMap[addr]; !ok {
			mergedAccountsMap[addr] = account
		} else {
			mergedAccountsMap[addr] = mergeOutputAccounts(mergedAccountsMap[addr], account)
		}
	}

	return mergedAccountsMap
}

func mergeOutputAccounts(leftAccount *vmcommon.OutputAccount, rightAccount *vmcommon.OutputAccount) *vmcommon.OutputAccount {
	// TODO Discuss merging each of the fields of two OutputAccount instances
	mergedAccount := &vmcommon.OutputAccount{}

	mergedAccount.Address = leftAccount.Address

	leftDelta := leftAccount.BalanceDelta
	rightDelta := rightAccount.BalanceDelta
	if leftDelta == nil {
		leftDelta = big.NewInt(0)
	}
	if rightDelta == nil {
		rightDelta = big.NewInt(0)
	}
	mergedAccount.BalanceDelta = big.NewInt(0).Add(leftDelta, rightDelta)

	if leftAccount.Nonce > 0 {
		mergedAccount.Nonce = leftAccount.Nonce
	}

	if rightAccount.Nonce > mergedAccount.Nonce {
		mergedAccount.Nonce = rightAccount.Nonce
	}

	mergedAccount.StorageUpdates = append(mergedAccount.StorageUpdates, leftAccount.StorageUpdates...)
	mergedAccount.StorageUpdates = append(mergedAccount.StorageUpdates, rightAccount.StorageUpdates...)

	mergedAccount.Code = rightAccount.Code
	mergedAccount.Data = rightAccount.Data

	mergedAccount.GasLimit = rightAccount.GasLimit

	return mergedAccount
}
