package host

import (
	"math/big"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (host *vmHost) handleAsyncCallBreakpoint(result wasmer.Value, argError error) error {
	runtime := host.Runtime()
	output := host.Output()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	// TODO also determine whether caller and callee are in the same Shard, by
	// account addresses - this makes the empty SC code an error
	syncCall, err := host.canExecuteSynchronously()
	if err != nil {
		return err
	}
	if !syncCall {
		return argError
	}

	senderVMOutput := output.GetVMOutput(result)
	intermediaryVMOutput := senderVMOutput

	// Start calling the destination SC, synchronously.
	destinationCallInput, err := host.createDestinationContractCallInput()
	if err != nil {
		return err
	}

	destinationVMOutput, err := host.executeOnNewContextAndGetVMOutput(destinationCallInput)
	// TODO pass error to the SC callback (append as argument)
	// TODO consume remaining gas
	if err != nil {
		return err
	}
	intermediaryVMOutput = mergeTwoVMOutputs(intermediaryVMOutput, destinationVMOutput)

	callbackCallInput, err := host.createCallbackContractCallInput(destinationVMOutput)
	// TODO handle this error properly
	if err != nil {
		return err
	}

	callbackVMOutput, err := host.executeOnNewContextAndGetVMOutput(callbackCallInput)
	// TODO handle this error properly
	if err != nil {
		return err
	}
	finalVMOutput := mergeTwoVMOutputs(intermediaryVMOutput, callbackVMOutput)

	return finalVMOutput, nil
}

func (host *vmHost) canExecuteSynchronously() (bool, error) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()
	asyncCallInfo := runtime.GetAsyncCallInfo()
	dest := asyncCallInfo.Destination
	calledSCCode, err := blockchain.GetCode(dest)

	return len(calledSCCode) != 0, err
}

func (host *vmHost) createDestinationContractCallInput() (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()
	asyncCallInfo := runtime.GetAsyncCallInfo()

	argParser := runtime.ArgParser()
	err := argParser.ParseData(string(asyncCallInfo.Data))
	if err != nil {
		return nil, err
	}

	function, err := argParser.GetFunction()
	if err != nil {
		return nil, err
	}

	arguments, err := argParser.GetArguments()
	if err != nil {
		return nil, err
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   arguments,
			CallValue:   big.NewInt(0).SetBytes(asyncCallInfo.ValueBytes),
			GasPrice:    0,
			GasProvided: asyncCallInfo.GasLimit,
		},
		RecipientAddr: asyncCallInfo.Destination,
		Function:      function,
	}

	return contractCallInput, nil
}

func (host *vmHost) createCallbackContractCallInput(destinationVMOutput *vmcommon.VMOutput) (*vmcommon.ContractCallInput, error) {
	metering := host.Metering()
	runtime := host.Runtime()

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
	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(dataLength)
	gasLimit -= gasToUse

	sender := runtime.GetAsyncCallInfo().Destination
	dest := runtime.GetSCAddress()

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

func (host *vmHost) executeOnNewContextAndGetVMOutput(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
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
	if err != nil {
		return nil, err
	}

	vmOutput := host.Output().CreateVMOutput(wasmer.Void())
	return vmOutput, err
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
