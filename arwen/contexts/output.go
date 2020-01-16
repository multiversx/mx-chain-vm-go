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
	stateStack  []*outputState
}

func NewOutputContext(host arwen.VMHost) (*outputContext, error) {
	context := &outputContext{
		host:       host,
		stateStack: make([]*outputState, 0),
	}

	context.InitState()

	return context, nil
}

func (context *outputContext) InitState() {
	context.outputState = newOutputState()
}

func (context *outputContext) PushState() {
	newState := newOutputState()
  newState.update(context.outputState)
	context.stateStack = append(context.stateStack, newState)
}

func (context *outputContext) PopState() error {
	stateStackLen := len(context.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

  prevState.update(context.outputState)
  context.outputState = newOutputState()
  context.outputState.update(prevState)

	return nil
}

func (context *outputContext) GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate  {
  account, ok := context.outputState.OutputAccounts[string(address)]
  if !ok {
    account = newOutputAccount(address)
    context.outputState.OutputAccounts[string(address)] = account
  }

  return account.StorageUpdates
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

  logs := context.outputState.Logs

	strAdr := string(addr)

	if _, ok := logs[strAdr]; !ok {
		logs[strAdr] = &vmcommon.LogEntry{
      Address: addr,
			Topics: make([][]byte, 0),
			Data:   make([]byte, 0),
		}
	}

	currLogs := logs[strAdr]
	for i := 0; i < len(topics); i++ {
		topics = append(currLogs.Topics, topics[i])
	}
	data = append(currLogs.Data, data...)

	logs[strAdr] = currLogs
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
  vmOutput := context.outputState.ToVMOutput()
	return vmOutput
}

func (context *outputContext) DeployCode(address []byte, code []byte) {
	newSCAcc, ok := context.outputState.OutputAccounts[string(address)]
	if !ok {
    newSCAcc = newOutputAccount(address)
		context.outputState.OutputAccounts[string(address)] = newSCAcc
  } 
  newSCAcc.Code = code
}

func (context *outputContext) CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{GasRemaining: 0, GasRefund: big.NewInt(0)}
	vmOutput.ReturnCode = errCode
	vmOutput.ReturnMessage = message
	return vmOutput
}
