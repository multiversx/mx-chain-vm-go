package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type logTopicsData struct {
	topics [][]byte
	data   []byte
}

type outputContext struct {
	host        arwen.VMHost
	outputState *vmcommon.VMOutput
	stateStack  []*vmcommon.VMOutput
}

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

func newVMOutputAccount(address []byte) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

func (context *outputContext) PushState() {
	newState := newVMOutput()
	mergeVMOutputs(newState, context.outputState)
	context.stateStack = append(context.stateStack, newState)
}

func (context *outputContext) PopState() {
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

func (context *outputContext) ClearStateStack() {
	context.stateStack = make([]*vmcommon.VMOutput, 0)
}

func (context *outputContext) GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool) {
	accountIsNew := false
	account, ok := context.outputState.OutputAccounts[string(address)]
	if !ok {
		account = newVMOutputAccount(address)
		context.outputState.OutputAccounts[string(address)] = account
		accountIsNew = true
	}

	return account, accountIsNew
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

func (context *outputContext) SelfDestruct(address []byte, beneficiary []byte) {
	panic("not implemented")
}

func (context *outputContext) Finish(data []byte) {
	if len(data) > 0 {
		context.outputState.ReturnData = append(context.outputState.ReturnData, data)
	}
}

func (context *outputContext) FinishValue(value wasmer.Value) {
	if !value.IsVoid() {
		valueBytes := arwen.ConvertReturnValue(value)
		context.Finish(valueBytes)
	}
}

func (context *outputContext) WriteLog(address []byte, topics [][]byte, data []byte) {
	if context.host.Runtime().ReadOnly() {
		return
	}

	newLogEntry := &vmcommon.LogEntry{
		Address: address,
		Topics:  topics,
		Data:    data,
	}
	context.outputState.Logs = append(context.outputState.Logs, newLogEntry)
}

// Transfer handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (context *outputContext) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) {
	senderAcc, _ := context.GetOutputAccount(sender)
	destAcc, _ := context.GetOutputAccount(destination)

	senderAcc.BalanceDelta = big.NewInt(0).Sub(senderAcc.BalanceDelta, value)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
	destAcc.Data = append(destAcc.Data, input...)
	destAcc.GasLimit = gasLimit
}

func (context *outputContext) AddTxValueToAccount(address []byte, value *big.Int) {
	destAcc, _ := context.GetOutputAccount(address)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

// GetVMOutput updates the current VMOutput and returns it
func (context *outputContext) GetVMOutput() *vmcommon.VMOutput {
	context.outputState.GasRemaining = context.host.Metering().GasLeft()
	return context.outputState
}

func (context *outputContext) DeployCode(address []byte, code []byte) {
	newSCAcc, _ := context.GetOutputAccount(address)
	newSCAcc.Code = code
}

func (context *outputContext) CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		GasRemaining:  0,
		GasRefund:     big.NewInt(0),
		ReturnCode:    errCode,
		ReturnMessage: message,
	}
}

func mergeVMOutputs(leftOutput *vmcommon.VMOutput, rightOutput *vmcommon.VMOutput) {
	if leftOutput.OutputAccounts == nil {
		leftOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
	}
	for address, rightAccount := range rightOutput.OutputAccounts {
		leftAccount, ok := leftOutput.OutputAccounts[address]
		if !ok {
			leftAccount = &vmcommon.OutputAccount{}
			leftOutput.OutputAccounts[address] = leftAccount
		}
		mergeOutputAccounts(leftAccount, rightAccount)
	}

	// TODO merge DeletedAccounts and TouchedAccounts as well?

	leftOutput.Logs = append(leftOutput.Logs, rightOutput.Logs...)
	leftOutput.ReturnData = append(leftOutput.ReturnData, rightOutput.ReturnData...)
	leftOutput.GasRemaining = rightOutput.GasRemaining
	leftOutput.GasRefund = rightOutput.GasRefund
	leftOutput.ReturnCode = rightOutput.ReturnCode
	leftOutput.ReturnMessage = rightOutput.ReturnMessage
}

func mergeOutputAccounts(
	leftAccount *vmcommon.OutputAccount,
	rightAccount *vmcommon.OutputAccount,
) {
	leftAccount.Address = rightAccount.Address
	leftAccount.GasLimit = rightAccount.GasLimit
	mergeStorageUpdates(leftAccount, rightAccount)

	if rightAccount.Balance != nil {
		leftAccount.Balance = rightAccount.Balance
	}
	if leftAccount.BalanceDelta == nil {
		leftAccount.BalanceDelta = big.NewInt(0)
	}
	if rightAccount.BalanceDelta != nil {
		leftAccount.BalanceDelta = big.NewInt(0).Add(leftAccount.BalanceDelta, rightAccount.BalanceDelta)
	}
	if len(rightAccount.Code) > 0 {
		leftAccount.Code = rightAccount.Code
	}
	if len(rightAccount.Data) > 0 {
		leftAccount.Data = rightAccount.Data
	}
	if rightAccount.Nonce > leftAccount.Nonce {
		leftAccount.Nonce = rightAccount.Nonce
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
