package subcontexts

import (
	"math/big"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
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
}

func NewOutputSubcontext(host arwen.VMContext) (*Output, error) {
	output := &Output{
		outputAccounts: make(map[string]*vmcommon.OutputAccount),
		host:           host,
	}

	return output, nil
}

func (output *Output) CreateStateCopy() *Output {
	return &Output{
		logs:           output.logs,
		storageUpdate:  output.storageUpdate,
		outputAccounts: output.outputAccounts,
		returnData:     output.returnData,
		returnCode:     output.returnCode,
		selfDestruct:   output.selfDestruct,
		refund:         output.refund,
	}
}

func (output *Output) LoadFromStateCopy(otherOutput *Output) {
	for key, log := range otherOutput.logs {
		output.logs[key] = log
	}
	for key, storageUpdate := range otherOutput.storageUpdate {
		if _, ok := output.storageUpdate[key]; !ok {
			output.storageUpdate[key] = storageUpdate
			continue
		}

		for internKey, internStore := range storageUpdate {
			output.storageUpdate[key][internKey] = internStore
		}
	}

	output.outputAccounts = otherOutput.outputAccounts
	output.returnData = append(output.returnData, otherOutput.returnData...)
	output.returnCode = otherOutput.returnCode

	for key, selfDestruct := range otherOutput.selfDestruct {
		output.selfDestruct[key] = selfDestruct
	}

	output.refund += otherOutput.refund
}

func (output *Output) GetOutputAccounts() map[string]*vmcommon.OutputAccount {
	return output.outputAccounts
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
	panic("not implemented")
}
