package evmhooks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"math/big"
)

func (context *EVMHooksImpl) FailExecution(err error) {
	context.GetRuntimeContext().FailExecution(err)
}

func (context *EVMHooksImpl) ReadOnly() bool {
	return context.GetRuntimeContext().ReadOnly()
}

func (context *EVMHooksImpl) Origin() common.Address {
	originalCallerAddress := context.GetRuntimeContext().GetOriginalCallerAddress()
	if len(originalCallerAddress) > 0 {
		return context.toEVMAddress(originalCallerAddress)
	}
	return common.Address{}
}

func (context *EVMHooksImpl) GasPrice() *big.Int {
	return new(big.Int).SetUint64(context.GetRuntimeContext().GetVMInput().GasPrice)
}

func (context *EVMHooksImpl) CallerMvxAddress() []byte {
	return context.GetRuntimeContext().GetVMInput().CallerAddr
}

func (context *EVMHooksImpl) CallerAddress() common.Address {
	return context.toEVMAddress(context.CallerMvxAddress())
}

func (context *EVMHooksImpl) CallValue() *uint256.Int {
	return uint256.MustFromBig(context.GetRuntimeContext().GetVMInput().CallValue)
}

func (context *EVMHooksImpl) ContractMvxAddress() []byte {
	return context.GetRuntimeContext().GetContextAddress()
}

func (context *EVMHooksImpl) ContractAddress() common.Address {
	return context.toEVMAddress(context.ContractMvxAddress())
}

func (context *EVMHooksImpl) ContractAliasAddress() common.Address {
	address := context.GetRuntimeContext().GetVMInput().RecipientAliasAddr
	if len(address) > 0 {
		return common.BytesToAddress(address)
	}
	return common.Address{}
}

func (context *EVMHooksImpl) Arguments() [][]byte {
	return context.GetRuntimeContext().Arguments()
}

func (context *EVMHooksImpl) CodeHash() common.Hash {
	return common.BytesToHash(context.GetRuntimeContext().GetSCCodeHash())
}
