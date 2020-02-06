package mock

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type OutputContextMock struct {
	outputState *vmcommon.VMOutput
	stateStack  []*vmcommon.VMOutput
}

func NewOutputContextMock() *OutputContextMock {
	context := &OutputContextMock{
		stateStack: make([]*vmcommon.VMOutput, 0),
	}

	context.InitState()

	return context
}

func (context *OutputContextMock) InitState() {
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

func (context *OutputContextMock) PushState() {
	newState := newVMOutput()
	mergeVMOutputs(newState, context.outputState)
	context.stateStack = append(context.stateStack, newState)
}

func (context *OutputContextMock) PopState() error {
	stateStackLen := len(context.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

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

	return nil
}

func (context *OutputContextMock) GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool) {
	accountIsNew := false
	account, ok := context.outputState.OutputAccounts[string(address)]
	if !ok {
		account = newVMOutputAccount(address)
		context.outputState.OutputAccounts[string(address)] = account
		accountIsNew = true
	}

	return account, accountIsNew
}

func (context *OutputContextMock) GetRefund() uint64 {
	return uint64(context.outputState.GasRefund.Int64())
}

func (context *OutputContextMock) SetRefund(refund uint64) {
	context.outputState.GasRefund = big.NewInt(int64(refund))
}

func (context *OutputContextMock) ReturnData() [][]byte {
	return context.outputState.ReturnData
}

func (context *OutputContextMock) ReturnCode() vmcommon.ReturnCode {
	return context.outputState.ReturnCode
}

func (context *OutputContextMock) SetReturnCode(returnCode vmcommon.ReturnCode) {
	context.outputState.ReturnCode = returnCode
}

func (context *OutputContextMock) ReturnMessage() string {
	return context.outputState.ReturnMessage
}

func (context *OutputContextMock) SetReturnMessage(returnMessage string) {
	context.outputState.ReturnMessage = returnMessage
}

func (context *OutputContextMock) ClearReturnData() {
	context.outputState.ReturnData = make([][]byte, 0)
}

func (context *OutputContextMock) SelfDestruct(_ []byte, _ []byte) {
	panic("not implemented")
}

func (context *OutputContextMock) Finish(data []byte) {
	if len(data) > 0 {
		context.outputState.ReturnData = append(context.outputState.ReturnData, data)
	}
}

func (context *OutputContextMock) FinishValue(value wasmer.Value) {
	if !value.IsVoid() {
		convertedResult := arwen.ConvertReturnValue(value)
		valueBytes := convertedResult.Bytes()

		context.Finish(valueBytes)
	}
}

func (context *OutputContextMock) WriteLog(address []byte, topics [][]byte, data []byte) {
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
func (context *OutputContextMock) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) {
	senderAcc, _ := context.GetOutputAccount(sender)
	destAcc, _ := context.GetOutputAccount(destination)

	senderAcc.BalanceDelta = big.NewInt(0).Sub(senderAcc.BalanceDelta, value)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
	destAcc.Data = append(destAcc.Data, input...)
	destAcc.GasLimit = gasLimit
}

func (context *OutputContextMock) AddTxValueToAccount(address []byte, value *big.Int) {
	destAcc, _ := context.GetOutputAccount(address)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

// adapt vm output and all saved data from sc run into VM Output
func (context *OutputContextMock) GetVMOutput(_ wasmer.Value) *vmcommon.VMOutput {
	return context.outputState
}

func (context *OutputContextMock) DeployCode(address []byte, code []byte) {
	newSCAcc, _ := context.GetOutputAccount(address)
	newSCAcc.Code = code
}

func (context *OutputContextMock) CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{GasRemaining: 0, GasRefund: big.NewInt(0)}
	vmOutput.ReturnCode = errCode
	vmOutput.ReturnMessage = message
	return vmOutput
}

func mergeVMOutputs(leftOutput *vmcommon.VMOutput, rightOutput *vmcommon.VMOutput) {
	for address, rightAccount := range rightOutput.OutputAccounts {
		leftAccount, ok := leftOutput.OutputAccounts[address]
		if !ok {
			leftAccount = &vmcommon.OutputAccount{}
			leftOutput.OutputAccounts[address] = leftAccount
		}
		mergeOutputAccounts(leftAccount, rightAccount)
	}

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
