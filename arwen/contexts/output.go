package contexts

import (
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.OutputContext = (*outputContext)(nil)

type outputContext struct {
	host        arwen.VMHost
	outputState *vmcommon.VMOutput
	stateStack  []*vmcommon.VMOutput
	codeUpdates map[string]struct{}
}

// NewOutputContext creates a new outputContext
func NewOutputContext(host arwen.VMHost) (*outputContext, error) {
	context := &outputContext{
		host:       host,
		stateStack: make([]*vmcommon.VMOutput, 0),
	}

	context.InitState()

	return context, nil
}

func (context *outputContext) InitState() {
	context.outputState = newVMOutput()
	context.codeUpdates = make(map[string]struct{})
}

func newVMOutput() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}
}

func NewVMOutputAccount(address []byte) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        nil,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

func (context *outputContext) PushState() {
	newState := newVMOutput()
	mergeVMOutputs(newState, context.outputState)
	context.stateStack = append(context.stateStack, newState)
}

func (context *outputContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.outputState = prevState
}

func (context *outputContext) PopMergeActiveState() {
	stateStackLen := len(context.stateStack)

	// Merge the current state into the head of the stateStack,
	// then pop the head of the stateStack into the current state.
	// Doing this allows the VM to execute a SmartContract into a context on top
	// of an existing context (a previous SC) without allowing access to it, but
	// later merging the output of the two SCs in chronological order.
	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	mergeVMOutputs(prevState, context.outputState)
	context.outputState = newVMOutput()
	mergeVMOutputs(context.outputState, prevState)
}

func (context *outputContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	context.stateStack = context.stateStack[:stateStackLen-1]
}

func (context *outputContext) ClearStateStack() {
	context.stateStack = make([]*vmcommon.VMOutput, 0)
}

// CensorVMOutput will cause the next executed SC to appear isolated, as if
// nothing was executed before. Required for ExecuteOnDestContext().
// StorageUpdates are not deleted from context.outputState.OutputAccounts,
// preserving the storage cache.
func (context *outputContext) CensorVMOutput() {
	context.outputState.ReturnData = make([][]byte, 0)
	context.outputState.ReturnCode = vmcommon.Ok
	context.outputState.ReturnMessage = ""
	context.outputState.GasRemaining = 0
	context.outputState.GasRefund = big.NewInt(0)
	context.outputState.Logs = make([]*vmcommon.LogEntry, 0)
}

// ResetGas will set to 0 all gas used from output accounts, in order
// to properly calculate the actual used gas of one smart contract when called in sync
func (context *outputContext) ResetGas() {
	for _, outAcc := range context.outputState.OutputAccounts {
		outAcc.GasUsed = 0
		for _, outTransfer := range outAcc.OutputTransfers {
			outTransfer.GasLimit = 0
		}
	}
}

func (context *outputContext) GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool) {
	accountIsNew := false
	account, ok := context.outputState.OutputAccounts[string(address)]
	if !ok {
		account = NewVMOutputAccount(address)
		context.outputState.OutputAccounts[string(address)] = account
		accountIsNew = true
	}

	return account, accountIsNew
}

func (context *outputContext) DeleteOutputAccount(address []byte) {
	delete(context.outputState.OutputAccounts, string(address))
	delete(context.codeUpdates, string(address))
}

func (context *outputContext) GetRefund() uint64 {
	return uint64(context.outputState.GasRefund.Int64())
}

func (context *outputContext) SetRefund(refund uint64) {
	context.outputState.GasRefund = big.NewInt(int64(refund))
}

func (context *outputContext) ReturnData() [][]byte {
	return context.outputState.ReturnData
}

func (context *outputContext) ReturnCode() vmcommon.ReturnCode {
	return context.outputState.ReturnCode
}

func (context *outputContext) SetReturnCode(returnCode vmcommon.ReturnCode) {
	context.outputState.ReturnCode = returnCode
}

func (context *outputContext) ReturnMessage() string {
	return context.outputState.ReturnMessage
}

func (context *outputContext) SetReturnMessage(returnMessage string) {
	context.outputState.ReturnMessage = returnMessage
}

func (context *outputContext) ClearReturnData() {
	context.outputState.ReturnData = make([][]byte, 0)
}

func (context *outputContext) SelfDestruct(_ []byte, _ []byte) {
}

func (context *outputContext) Finish(data []byte) {
	context.outputState.ReturnData = append(context.outputState.ReturnData, data)
}

func (context *outputContext) WriteLog(address []byte, topics [][]byte, data []byte) {
	if context.host.Runtime().ReadOnly() {
		return
	}

	newLogEntry := &vmcommon.LogEntry{
		Address: address,
		Data:    data,
	}

	if len(topics) == 0 {
		context.outputState.Logs = append(context.outputState.Logs, newLogEntry)
		return
	}

	newLogEntry.Identifier = topics[0]
	newLogEntry.Topics = topics[1:]

	context.outputState.Logs = append(context.outputState.Logs, newLogEntry)
}

// TransferValueOnly will transfer the big.int value and checks if it is possible
func (context *outputContext) TransferValueOnly(destination []byte, sender []byte, value *big.Int) error {
	if value.Cmp(arwen.Zero) < 0 {
		return arwen.ErrTransferNegativeValue
	}

	if !context.hasSufficientBalance(sender, value) {
		return arwen.ErrTransferInsufficientFunds
	}

	payable, err := context.host.Blockchain().IsPayable(destination)
	if err != nil {
		return err
	}

	hasValue := value.Cmp(big.NewInt(0)) == 1
	if !payable && hasValue {
		return arwen.ErrAccountNotPayable
	}

	senderAcc, _ := context.GetOutputAccount(sender)
	destAcc, _ := context.GetOutputAccount(destination)

	senderAcc.BalanceDelta = big.NewInt(0).Sub(senderAcc.BalanceDelta, value)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)

	return nil
}

// Transfer handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (context *outputContext) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte, callType vmcommon.CallType) error {
	err := context.TransferValueOnly(destination, sender, value)
	if err != nil {
		return err
	}

	destAcc, _ := context.GetOutputAccount(destination)
	outputTransfer := vmcommon.OutputTransfer{
		Value:    big.NewInt(0).Set(value),
		GasLimit: gasLimit,
		Data:     input,
		CallType: callType,
	}
	destAcc.OutputTransfers = append(destAcc.OutputTransfers, outputTransfer)

	return nil
}

func (context *outputContext) hasSufficientBalance(address []byte, value *big.Int) bool {
	senderBalance := context.host.Blockchain().GetBalanceBigInt(address)
	return senderBalance.Cmp(value) >= 0
}

func (context *outputContext) AddTxValueToAccount(address []byte, value *big.Int) {
	destAcc, _ := context.GetOutputAccount(address)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

// GetVMOutput updates the current VMOutput and returns it
func (context *outputContext) GetVMOutput() *vmcommon.VMOutput {
	if context.outputState.ReturnCode == vmcommon.Ok {
		context.outputState.GasRemaining = context.host.Metering().GasLeft()
	} else {
		context.outputState.GasRemaining = context.host.Metering().GetGasLockedForAsyncStep()
	}

	context.removeNonUpdatedCode(context.outputState)

	return context.outputState
}

func (context *outputContext) DeployCode(input arwen.CodeDeployInput) {
	newSCAccount, _ := context.GetOutputAccount(input.ContractAddress)
	newSCAccount.Code = input.ContractCode
	newSCAccount.CodeMetadata = input.ContractCodeMetadata
	newSCAccount.CodeDeployerAddress = input.CodeDeployerAddress

	var empty struct{}
	context.codeUpdates[string(input.ContractAddress)] = empty
}

func (context *outputContext) CreateVMOutputInCaseOfError(err error) *vmcommon.VMOutput {
	metering := context.host.Metering()
	var message string

	if err == arwen.ErrSignalError {
		message = context.ReturnMessage()
	} else {
		if len(context.outputState.ReturnMessage) > 0 {
			// another return message was already set previously
			message = context.outputState.ReturnMessage
		} else {
			message = err.Error()
		}
	}

	returnCode := context.resolveReturnCodeFromError(err)

	return &vmcommon.VMOutput{
		GasRemaining:  metering.GetGasLockedForAsyncStep(),
		GasRefund:     big.NewInt(0),
		ReturnCode:    returnCode,
		ReturnMessage: message,
	}
}

func (context *outputContext) removeNonUpdatedCode(vmOutput *vmcommon.VMOutput) {
	for address, account := range vmOutput.OutputAccounts {
		_, ok := context.codeUpdates[address]
		if !ok {
			account.Code = nil
			account.CodeMetadata = nil
			account.CodeDeployerAddress = nil
		}
	}
}

func (context *outputContext) resolveReturnCodeFromError(err error) vmcommon.ReturnCode {
	if err == nil {
		return vmcommon.Ok
	}

	log.Info("resolveReturnCodeFromError", "error", err)
	if errors.Is(err, arwen.ErrSignalError) {
		return vmcommon.UserError
	}
	if errors.Is(err, arwen.ErrFuncNotFound) {
		return vmcommon.FunctionNotFound
	}
	if errors.Is(err, arwen.ErrFunctionNonvoidSignature) {
		return vmcommon.FunctionWrongSignature
	}
	if errors.Is(err, arwen.ErrInvalidFunction) {
		return vmcommon.UserError
	}
	if errors.Is(err, arwen.ErrNotEnoughGas) {
		log.Info("this is out of gas resolveReturnCodeFromError")
		return vmcommon.OutOfGas
	}
	if errors.Is(err, arwen.ErrContractNotFound) {
		return vmcommon.ContractNotFound
	}
	if errors.Is(err, arwen.ErrContractInvalid) {
		return vmcommon.ContractInvalid
	}
	if errors.Is(err, arwen.ErrUpgradeFailed) {
		return vmcommon.UpgradeFailed
	}
	if errors.Is(err, arwen.ErrTransferInsufficientFunds) {
		return vmcommon.OutOfFunds
	}

	return vmcommon.ExecutionFailed
}

func (context *outputContext) AddToActiveState(rightOutput *vmcommon.VMOutput) {
	rightOutput.GasRemaining = 0
	if rightOutput.GasRefund != nil {
		rightOutput.GasRefund.Add(rightOutput.GasRefund, context.outputState.GasRefund)
	}

	for _, rightAccount := range rightOutput.OutputAccounts {
		leftAccount, ok := context.outputState.OutputAccounts[string(rightAccount.Address)]
		if !ok {
			continue
		}

		if rightAccount.BalanceDelta != nil {
			rightAccount.BalanceDelta.Add(rightAccount.BalanceDelta, leftAccount.BalanceDelta)
		}
		if len(rightAccount.OutputTransfers) > 0 {
			leftAccount.OutputTransfers = append(leftAccount.OutputTransfers, rightAccount.OutputTransfers...)
		}
	}

	mergeVMOutputs(context.outputState, rightOutput)
}

func mergeVMOutputs(leftOutput *vmcommon.VMOutput, rightOutput *vmcommon.VMOutput) {
	if leftOutput.OutputAccounts == nil {
		leftOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
	}

	for _, rightAccount := range rightOutput.OutputAccounts {
		leftAccount, ok := leftOutput.OutputAccounts[string(rightAccount.Address)]
		if !ok {
			leftAccount = &vmcommon.OutputAccount{}
			leftOutput.OutputAccounts[string(rightAccount.Address)] = leftAccount
		}
		mergeOutputAccounts(leftAccount, rightAccount)
	}

	leftOutput.Logs = append(leftOutput.Logs, rightOutput.Logs...)
	leftOutput.ReturnData = append(leftOutput.ReturnData, rightOutput.ReturnData...)
	leftOutput.GasRemaining = rightOutput.GasRemaining
	leftOutput.GasRefund = rightOutput.GasRefund
	if leftOutput.GasRefund == nil {
		leftOutput.GasRefund = big.NewInt(0)
	}

	leftOutput.ReturnCode = rightOutput.ReturnCode
	leftOutput.ReturnMessage = rightOutput.ReturnMessage
}

func mergeOutputAccounts(
	leftAccount *vmcommon.OutputAccount,
	rightAccount *vmcommon.OutputAccount,
) {
	if len(rightAccount.Address) != 0 {
		leftAccount.Address = rightAccount.Address
	}

	mergeStorageUpdates(leftAccount, rightAccount)

	if rightAccount.Balance != nil {
		leftAccount.Balance = rightAccount.Balance
	}
	if leftAccount.BalanceDelta == nil {
		leftAccount.BalanceDelta = big.NewInt(0)
	}
	if rightAccount.BalanceDelta != nil {
		leftAccount.BalanceDelta = rightAccount.BalanceDelta
	}
	if len(rightAccount.Code) > 0 {
		leftAccount.Code = rightAccount.Code
	}
	if len(rightAccount.CodeMetadata) > 0 {
		leftAccount.CodeMetadata = rightAccount.CodeMetadata
	}
	if rightAccount.Nonce > leftAccount.Nonce {
		leftAccount.Nonce = rightAccount.Nonce
	}

	lenLeftOutTransfers := len(leftAccount.OutputTransfers)
	lenRightOutTransfers := len(rightAccount.OutputTransfers)
	if lenRightOutTransfers > lenLeftOutTransfers {
		leftAccount.OutputTransfers = append(leftAccount.OutputTransfers, rightAccount.OutputTransfers[lenLeftOutTransfers:]...)
	}

	leftAccount.GasUsed = rightAccount.GasUsed

	if rightAccount.CodeDeployerAddress != nil {
		leftAccount.CodeDeployerAddress = rightAccount.CodeDeployerAddress
	}
}

func mergeStorageUpdates(
	leftAccount *vmcommon.OutputAccount,
	rightAccount *vmcommon.OutputAccount,
) {
	if leftAccount.StorageUpdates == nil {
		leftAccount.StorageUpdates = make(map[string]*vmcommon.StorageUpdate)
	}
	for key, update := range rightAccount.StorageUpdates {
		leftAccount.StorageUpdates[key] = update
	}
}
