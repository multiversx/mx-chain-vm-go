package evmhooks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

func (context *EVMHooksImpl) Finish(returnData []byte) {
	context.GetOutputContext().ClearReturnData()
	if returnData != nil {
		context.GetOutputContext().Finish(returnData)
	}
}

func (context *EVMHooksImpl) FinishCreate(returnData []byte) {
	context.saveCompiledCode(returnData)
	context.Finish(returnData)
}

func (context *EVMHooksImpl) TransferBalance(destination common.Address, value *uint256.Int) error {
	sender := context.ContractMvxAddress()
	return context.GetOutputContext().TransferValueOnly(context.toMVXAddress(destination), sender, value.ToBig(), true)
}

func (context *EVMHooksImpl) SelfDestruct(destination common.Address) {
	err := context.TransferBalance(destination, context.GetSelfBalance())
	if !context.WithFault(err) {
		contract := context.ContractMvxAddress()
		context.GetOutputContext().DeleteAccount(contract)
	}
}

func (context *EVMHooksImpl) AddLog(log *types.Log) {
	topics := make([][]byte, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = topic.Bytes()
	}
	contract := context.ContractMvxAddress()
	context.GetOutputContext().WriteLog(contract, topics, [][]byte{log.Data})
}

func (context *EVMHooksImpl) saveCompiledCode(contract []byte) {
	runtime, output := context.GetRuntimeContext(), context.GetOutputContext()
	runtime.SetTrackerCode(contract)
	runtime.SaveCompiledCode()
	output.ChangeAccountCode(context.ContractMvxAddress(), contract)
}
