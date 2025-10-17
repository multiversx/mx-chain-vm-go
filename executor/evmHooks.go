package executor

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"math/big"
)

// EVMHooks contains all VM functions that can be called by the executor during SC execution.
type EVMHooks interface {
	EvmBlockchainHooks
	EvmMeteringHooks
	EvmOutputHooks
	EvmRuntimeHooks
	EvmStorageHooks
	EvmExecutionHooks
}

type EvmBlockchainHooks interface {
	ChainID() *big.Int
	Random() *common.Hash
	GetHash(number uint64) common.Hash
	BlockNumber() *big.Int
	Time() uint64
	GetSelfBalance() *uint256.Int
	GetBalance(address common.Address) *uint256.Int
	GetCodeHash(address common.Address) common.Hash
	GetCode(address common.Address) []byte
	GetCodeSize(address common.Address) int
	SaveAliasAddress() error
}

type EvmMeteringHooks interface {
	GasLeft() uint64
	UseGas(opCode string, gas uint64) bool
	BlockGasLimit() uint64
}

type EvmOutputHooks interface {
	Finish(returnData []byte)
	FinishCreate(returnData []byte)
	TransferBalance(destination common.Address, value *uint256.Int) error
	SelfDestruct(destination common.Address)
	AddLog(log *types.Log)
}

type EvmRuntimeHooks interface {
	FailExecution(err error)
	ReadOnly() bool
	Origin() common.Address
	GasPrice() *big.Int
	CallerAddress() common.Address
	ContractAddress() common.Address
	CallValue() *uint256.Int
	Arguments() [][]byte
	CodeHash() common.Hash
}

type EvmStorageHooks interface {
	GetState(key common.Hash) common.Hash
	SetState(key common.Hash, value common.Hash)
}

type EvmExecutionHooks interface {
	IsSmartContractAddress(address common.Address) bool
	Create(code []byte, gas uint64, value *uint256.Int) ([]byte, common.Address, error)
	Create2(code []byte, gas uint64, value *uint256.Int, salt *uint256.Int) ([]byte, common.Address, error)
	Call(address common.Address, value *uint256.Int, input []byte, gas uint64) ([]byte, error)
	StaticCall(address common.Address, input []byte, gas uint64) ([]byte, error)
	DelegateCall(address common.Address, input []byte, gas uint64) ([]byte, error)
	CallCode(address common.Address, value *uint256.Int, input []byte, gas uint64) ([]byte, error)
}
