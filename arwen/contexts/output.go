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
	outputState *outputState
	stateStack  []*outputContext
}

func NewOutputContext(host arwen.VMHost) (*outputContext, error) {
	context := &outputContext{
		host:       host,
		stateStack: make([]*outputContext, 0),
	}

	context.InitState()

	return context, nil
}

func (context *outputContext) InitState() {
	context.outputState = newOutputState()
}

func (context *outputContext) PushState() {
	newState := &outputContext{
		logs:           context.logs,
		storageUpdate:  context.storageUpdate,
		outputAccounts: context.outputAccounts,
		returnData:     context.returnData,
		returnCode:     context.returnCode,
		selfDestruct:   context.selfDestruct,
		returnMessage:  context.returnMessage,
		refund:         context.refund,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *outputContext) PopState() error {
	stateStackLen := len(context.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	for key, log := range prevState.logs {
		context.logs[key] = log
	}
	for account, updates := range prevState.storageUpdate {
		if _, ok := context.storageUpdate[account]; !ok {
			context.storageUpdate[account] = updates
			continue
		}

		for key, value := range updates {
			context.storageUpdate[account][key] = value
		}
	}

	context.outputAccounts = prevState.outputAccounts
	context.returnData = append(context.returnData, prevState.returnData...)
	context.returnCode = prevState.returnCode
	context.returnMessage = prevState.returnMessage

	context.refund += prevState.refund

	return nil
}

func (context *outputContext) GetOutputAccounts() map[string]*vmcommon.OutputAccount {
	return context.outputAccounts
}

func (context *outputContext) GetStorageUpdates() map[string](map[string][]byte) {
	return context.storageUpdate
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

func (context *outputContext) SelfDestruct(addr []byte, beneficiary []byte) {
	panic("not implemented")
}

func (context *outputContext) Finish(data []byte) {
	if len(data) > 0 {
		context.outputState.ReturnData = append(context.outputState.ReturnData, data)
	}
}

func (context *outputContext) FinishValue(value wasmer.Value) {
	if !value.IsVoid() {
		convertedResult := arwen.ConvertReturnValue(value)
		valueBytes := convertedResult.Bytes()

		context.Finish(valueBytes)
	}
}

func (context *outputContext) WriteLog(addr []byte, topics [][]byte, data []byte) {
	if context.host.Runtime().ReadOnly() {
		return
	}

	strAdr := string(addr)

	if _, ok := context.logs[strAdr]; !ok {
		context.logs[strAdr] = logTopicsData{
			topics: make([][]byte, 0),
			data:   make([]byte, 0),
		}
	}

	currLogs := context.logs[strAdr]
	for i := 0; i < len(topics); i++ {
		currLogs.topics = append(currLogs.topics, topics[i])
	}
	currLogs.data = append(currLogs.data, data...)

	context.logs[strAdr] = currLogs
}

// Transfer handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (context *outputContext) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) {
	senderAcc, ok := context.outputState.OutputAccounts[string(sender)]
	if !ok {
		senderAcc = &outputAccount{
			Address:      sender,
			BalanceDelta: big.NewInt(0),
		}
		context.outputState.OutputAccounts[string(senderAcc.Address)] = senderAcc
	}

	destAcc, ok := context.outputState.OutputAccounts[string(destination)]
	if !ok {
		destAcc = &outputAccount{
			Address:      destination,
			BalanceDelta: big.NewInt(0),
		}
		context.outputState.OutputAccounts[string(destAcc.Address)] = destAcc
	}

	senderAcc.BalanceDelta = big.NewInt(0).Sub(senderAcc.BalanceDelta, value)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
	destAcc.Data = append(destAcc.Data, input...)
	destAcc.GasLimit = gasLimit
}

func (context *outputContext) AddTxValueToAccount(address []byte, value *big.Int) {
	destAcc, ok := context.outputState.OutputAccounts[string(address)]
	if !ok {
		destAcc = &outputAccount{
			Address:      address,
			BalanceDelta: big.NewInt(0),
		}
		context.outputState.OutputAccounts[string(destAcc.Address)] = destAcc
	}

	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

// adapt vm output and all saved data from sc run into VM Output
func (context *outputContext) CreateVMOutput(result wasmer.Value) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{}
	// save storage updates
	outAccs := make(map[string]*vmcommon.OutputAccount, 0)
	for addr, updates := range context.storageUpdate {
		if _, ok := outAccs[addr]; !ok {
			outAccs[addr] = &vmcommon.OutputAccount{Address: []byte(addr)}
		}

		for key, value := range updates {
			storageUpdate := &vmcommon.StorageUpdate{
				Offset: []byte(key),
				Data:   value,
			}

			outAccs[addr].StorageUpdates = append(outAccs[addr].StorageUpdates, storageUpdate)
		}
	}

	// add balances
	for addr, outAcc := range context.outputAccounts {
		if _, ok := outAccs[addr]; !ok {
			outAccs[addr] = &vmcommon.OutputAccount{}
		}

		outAccs[addr].Address = outAcc.Address
		outAccs[addr].BalanceDelta = outAcc.BalanceDelta

		if len(outAcc.Code) > 0 {
			outAccs[addr].Code = outAcc.Code
		}
		if outAcc.Nonce > 0 {
			outAccs[addr].Nonce = outAcc.Nonce
		}
		if len(outAcc.Data) > 0 {
			outAccs[addr].Data = outAcc.Data
		}

		outAccs[addr].GasLimit = outAcc.GasLimit
	}

	// save to the output finally
	for _, outAcc := range outAccs {
		vmOutput.OutputAccounts = append(vmOutput.OutputAccounts, outAcc)
	}

	// save logs
	for addr, value := range context.logs {
		logEntry := &vmcommon.LogEntry{
			Address: []byte(addr),
			Data:    value.data,
			Topics:  value.topics,
		}

		vmOutput.Logs = append(vmOutput.Logs, logEntry)
	}

	if len(context.returnData) > 0 {
		vmOutput.ReturnData = append(vmOutput.ReturnData, context.returnData...)
	}

	convertedResult := arwen.ConvertReturnValue(result)
	resultBytes := convertedResult.Bytes()
	if len(resultBytes) > 0 {
		vmOutput.ReturnData = append(vmOutput.ReturnData, resultBytes)
	}

	vmOutput.GasRemaining = context.host.Metering().GasLeft()
	vmOutput.GasRefund = big.NewInt(0).SetUint64(context.refund)
	vmOutput.ReturnCode = context.returnCode
	vmOutput.ReturnMessage = context.returnMessage

	return vmOutput
}

func (context *outputContext) DeployCode(address []byte, code []byte) {
	newSCAcc, ok := context.outputAccounts[string(address)]
	if !ok {
		context.outputAccounts[string(address)] = &vmcommon.OutputAccount{
			Address:        address,
			Nonce:          0,
			BalanceDelta:   big.NewInt(0),
			StorageUpdates: nil,
			Code:           code,
		}
	} else {
		newSCAcc.Code = code
	}
}

func (context *outputContext) CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{GasRemaining: 0, GasRefund: big.NewInt(0)}
	vmOutput.ReturnCode = errCode
	vmOutput.ReturnMessage = message
	return vmOutput
}
