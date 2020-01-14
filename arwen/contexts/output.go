package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type LogTopicsData struct {
	topics [][]byte
	data   []byte
}

type Output struct {
	host           arwen.VMHost
	outputAccounts map[string]*vmcommon.OutputAccount
	logs           map[string]LogTopicsData
	storageUpdate  map[string](map[string][]byte)
	returnData     [][]byte
	returnCode     vmcommon.ReturnCode
	returnMessage  string
	selfDestruct   map[string][]byte
	refund         uint64
	stateStack     []*Output
}

func NewOutputContext(host arwen.VMHost) (*Output, error) {
	output := &Output{
		host:       host,
		stateStack: make([]*Output, 0),
	}

	output.InitState()

	return output, nil
}

func (output *Output) InitState() {
	output.outputAccounts = make(map[string]*vmcommon.OutputAccount, 0)
	output.logs = make(map[string]LogTopicsData, 0)
	output.storageUpdate = make(map[string]map[string][]byte, 0)
	output.selfDestruct = make(map[string][]byte)
	output.returnData = nil
	output.returnCode = vmcommon.Ok
	output.returnMessage = ""
	output.refund = 0
}

func (output *Output) PushState() {
	newState := &Output{
		logs:           output.logs,
		storageUpdate:  output.storageUpdate,
		outputAccounts: output.outputAccounts,
		returnData:     output.returnData,
		returnCode:     output.returnCode,
		selfDestruct:   output.selfDestruct,
		returnMessage:  output.returnMessage,
		refund:         output.refund,
	}

	output.stateStack = append(output.stateStack, newState)
}

func (output *Output) PopState() error {
	stateStackLen := len(output.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

	prevState := output.stateStack[stateStackLen-1]
	output.stateStack = output.stateStack[:stateStackLen-1]

	for key, log := range prevState.logs {
		output.logs[key] = log
	}
	for key, storageUpdate := range prevState.storageUpdate {
		if _, ok := output.storageUpdate[key]; !ok {
			output.storageUpdate[key] = storageUpdate
			continue
		}

		for internKey, internStore := range storageUpdate {
			output.storageUpdate[key][internKey] = internStore
		}
	}

	output.outputAccounts = prevState.outputAccounts
	output.returnData = append(output.returnData, prevState.returnData...)
	output.returnCode = prevState.returnCode
	output.returnMessage = prevState.returnMessage

	for key, selfDestruct := range prevState.selfDestruct {
		output.selfDestruct[key] = selfDestruct
	}

	output.refund += prevState.refund

	return nil
}

func (output *Output) GetOutputAccounts() map[string]*vmcommon.OutputAccount {
	return output.outputAccounts
}

func (output *Output) GetStorageUpdates() map[string](map[string][]byte) {
	return output.storageUpdate
}

func (output *Output) GetRefund() uint64 {
	return output.refund
}

func (output *Output) SetRefund(refund uint64) {
	output.refund = refund
}

func (output *Output) ReturnData() [][]byte {
	return output.returnData
}

func (output *Output) ReturnCode() vmcommon.ReturnCode {
	return output.returnCode
}

func (output *Output) SetReturnCode(returnCode vmcommon.ReturnCode) {
	output.returnCode = returnCode
}

func (output *Output) ReturnMessage() string {
	return output.returnMessage
}

func (output *Output) SetReturnMessage(returnMessage string) {
	output.returnMessage = returnMessage
}

func (output *Output) ClearReturnData() {
	output.returnData = make([][]byte, 0)
}

func (output *Output) SelfDestruct(addr []byte, beneficiary []byte) {
	if output.host.Runtime().ReadOnly() {
		return
	}

	output.selfDestruct[string(addr)] = beneficiary
}

func (output *Output) Finish(data []byte) {
	if len(data) > 0 {
		output.returnData = append(output.returnData, data)
	}
}

func (output *Output) FinishValue(value wasmer.Value) {
	if !value.IsVoid() {
		convertedResult := arwen.ConvertReturnValue(value)
		valueBytes := convertedResult.Bytes()

		output.Finish(valueBytes)
	}
}

func (output *Output) WriteLog(addr []byte, topics [][]byte, data []byte) {
	if output.host.Runtime().ReadOnly() {
		return
	}

	strAdr := string(addr)

	if _, ok := output.logs[strAdr]; !ok {
		output.logs[strAdr] = LogTopicsData{
			topics: make([][]byte, 0),
			data:   make([]byte, 0),
		}
	}

	currLogs := output.logs[strAdr]
	for i := 0; i < len(topics); i++ {
		currLogs.topics = append(currLogs.topics, topics[i])
	}
	currLogs.data = append(currLogs.data, data...)

	output.logs[strAdr] = currLogs
}

// Transfer handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (output *Output) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) {
	senderAcc, ok := output.outputAccounts[string(sender)]
	if !ok {
		senderAcc = &vmcommon.OutputAccount{
			Address:      sender,
			BalanceDelta: big.NewInt(0),
		}
		output.outputAccounts[string(senderAcc.Address)] = senderAcc
	}

	destAcc, ok := output.outputAccounts[string(destination)]
	if !ok {
		destAcc = &vmcommon.OutputAccount{
			Address:      destination,
			BalanceDelta: big.NewInt(0),
		}
		output.outputAccounts[string(destAcc.Address)] = destAcc
	}

	senderAcc.BalanceDelta = big.NewInt(0).Sub(senderAcc.BalanceDelta, value)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
	destAcc.Data = append(destAcc.Data, input...)
	destAcc.GasLimit = gasLimit
}

func (output *Output) AddTxValueToAccount(address []byte, value *big.Int) {
	destAcc, ok := output.outputAccounts[string(address)]
	if !ok {
		destAcc = &vmcommon.OutputAccount{
			Address:      address,
			BalanceDelta: big.NewInt(0),
		}
		output.outputAccounts[string(destAcc.Address)] = destAcc
	}

	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

// adapt vm output and all saved data from sc run into VM Output
func (output *Output) CreateVMOutput(result wasmer.Value) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{}
	// save storage updates
	outAccs := make(map[string]*vmcommon.OutputAccount, 0)
	for addr, updates := range output.storageUpdate {
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
	for addr, outAcc := range output.outputAccounts {
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
	for addr, value := range output.logs {
		logEntry := &vmcommon.LogEntry{
			Address: []byte(addr),
			Data:    value.data,
			Topics:  value.topics,
		}

		vmOutput.Logs = append(vmOutput.Logs, logEntry)
	}

	if len(output.returnData) > 0 {
		vmOutput.ReturnData = append(vmOutput.ReturnData, output.returnData...)
	}

	convertedResult := arwen.ConvertReturnValue(result)
	resultBytes := convertedResult.Bytes()
	if len(resultBytes) > 0 {
		vmOutput.ReturnData = append(vmOutput.ReturnData, resultBytes)
	}

	vmOutput.GasRemaining = output.host.Metering().GasLeft()
	vmOutput.GasRefund = big.NewInt(0).SetUint64(output.refund)
	vmOutput.ReturnCode = output.returnCode
	vmOutput.ReturnMessage = output.returnMessage

	return vmOutput
}

func (output *Output) DeployCode(address []byte, code []byte) {
	newSCAcc, ok := output.outputAccounts[string(address)]
	if !ok {
		output.outputAccounts[string(address)] = &vmcommon.OutputAccount{
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

func (output *Output) CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{GasRemaining: 0, GasRefund: big.NewInt(0)}
	vmOutput.ReturnCode = errCode
	vmOutput.ReturnMessage = message
	return vmOutput
}
