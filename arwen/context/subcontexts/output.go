package subcontexts

import (
	"math/big"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	context "github.com/ElrondNetwork/arwen-wasm-vm/arwen/context"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type LogTopicsData struct {
	topics [][]byte
	data   []byte
}

type Output struct {
	host           arwen.VMContext
	outputAccounts map[string]*vmcommon.OutputAccount
	logs           map[string]LogTopicsData
	storageUpdate  map[string](map[string][]byte)
	returnData     [][]byte
	returnCode     vmcommon.ReturnCode
	selfDestruct   map[string][]byte
	refund         uint64
	stateStack     []*Output
}

func NewOutputSubcontext(host arwen.VMContext) (*Output, error) {
	output := &Output{
		outputAccounts: make(map[string]*vmcommon.OutputAccount),
		host:           host,
		stateStack:     make([]*Output, 0),
	}

	return output, nil
}

func (output *Output) InitState() {
	storage.storageUpdate = make(map[string]map[string][]byte, 0)
	host.logs = make(map[string]logTopicsData, 0)
}

func (output *Output) PushState() {
	newState := &Output{
		logs:           output.logs,
		storageUpdate:  output.storageUpdate,
		outputAccounts: output.outputAccounts,
		returnData:     output.returnData,
		returnCode:     output.returnCode,
		selfDestruct:   output.selfDestruct,
		refund:         output.refund,
	}

	output.stateStack = append(output.stateStack, newState)
}

func (output *Output) PopState() error {
	stateStackLen := len(output.stateStack)
	if stateStackLen < 1 {
		return context.StateStackUnderflow
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

func (output *Output) ReturnData() [][]byte {
	return output.returnData
}

func (output *Output) ReturnCode() vmcommon.ReturnCode {
	return output.returnCode
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
	output.returnData = append(output.returnData, data)
}

// adapt vm output and all saved data from sc run into VM Output
func (output *Output) CreateVMOutput(result []byte) *vmcommon.VMOutput {
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
	if len(result) > 0 {
		vmOutput.ReturnData = append(vmOutput.ReturnData, result)
	}

	vmOutput.GasRemaining = output.host.Metering().GasLeft()
	vmOutput.GasRefund = big.NewInt(0).SetUint64(output.refund)
	vmOutput.ReturnCode = output.returnCode

	return vmOutput
}

