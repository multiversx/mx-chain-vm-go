// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
)

type (
	// GetHashFunc returns the n'th block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

func (evm *EVM) precompile(addr common.Address) (PrecompiledContract, bool) {
	var precompiles map[common.Address]PrecompiledContract
	switch {
	default:
		precompiles = PrecompiledContractsCancun
	}
	p, ok := precompiles[addr]
	return p, ok
}

// BlockContext provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type BlockContext struct {
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        uint64         // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
	BaseFee     *big.Int       // Provides information for BASEFEE (0 if vm runs with NoBaseFee flag and 0 gas price)
	BlobBaseFee *big.Int       // Provides information for BLOBBASEFEE (0 if vm runs with NoBaseFee flag and 0 blob gas price)
	Random      *common.Hash   // Provides information for PREVRANDAO
}

// TxContext provides the EVM with information about a transaction.
// All fields can change between transactions.
type TxContext struct {
	// Message information
	Origin     common.Address // Provides information for ORIGIN
	GasPrice   *big.Int       // Provides information for GASPRICE (and is used to zero the basefee if NoBaseFee is set)
	BlobHashes []common.Hash  // Provides information for BLOBHASH
	BlobFeeCap *big.Int       // Is used to zero the blobbasefee if NoBaseFee is set
}

// EVM is the Ethereum Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The EVM should never be reused and is not thread safe.
type EVM struct {
	// Context provides auxiliary blockchain related information
	Context BlockContext
	TxContext
	// StateDB gives access to the underlying state
	StateDB StateDB
	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// GasConfig contains the gas costs
	GasConfig *GasConfig
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	interpreter *EVMInterpreter
	// abort is used to abort the EVM calling operations
	abort atomic.Bool
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* and later
	// applied in opCall*.
	callGasTemp uint64
}

// NewEVM returns a new EVM. The returned EVM is not thread safe and should
// only ever be used *once*.
func NewEVM(blockCtx BlockContext, txCtx TxContext, statedb StateDB, chainConfig *params.ChainConfig, gasConfig *GasConfig, instructionSet *JumpTable) *EVM {
	evm := &EVM{
		Context:     blockCtx,
		TxContext:   txCtx,
		StateDB:     statedb,
		chainConfig: chainConfig,
		GasConfig:   gasConfig,
	}
	evm.interpreter = NewEVMInterpreter(evm, instructionSet)
	return evm
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM) Cancel() {
	evm.abort.Store(true)
}

// Cancelled returns true if Cancel has been called
func (evm *EVM) Cancelled() bool {
	return evm.abort.Load()
}

// Interpreter returns the current interpreter
func (evm *EVM) Interpreter() *EVMInterpreter {
	return evm.interpreter
}

func (evm *EVM) Call(addr common.Address, input []byte, gas uint64, value *uint256.Int) (ret []byte, err error) {
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, err = evm.RunPrecompiledAndConsumeGas(p, input, gas)
	} else {
		if evm.StateDB.IsSmartContractAddress(addr) {
			ret, err = evm.StateDB.Call(addr, value, input, gas)
		} else {
			ret, err = nil, evm.StateDB.TransferBalance(addr, value)
		}
	}
	return ret, err
}

func (evm *EVM) CallCode(addr common.Address, input []byte, gas uint64, value *uint256.Int) (ret []byte, err error) {
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, err = evm.RunPrecompiledAndConsumeGas(p, input, gas)
	} else {
		ret, err = evm.StateDB.CallCode(addr, value, input, gas)
	}
	return ret, err
}

func (evm *EVM) DelegateCall(addr common.Address, input []byte, gas uint64) (ret []byte, err error) {
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, err = evm.RunPrecompiledAndConsumeGas(p, input, gas)
	} else {
		ret, err = evm.StateDB.DelegateCall(addr, input, gas)
	}
	return ret, err
}

func (evm *EVM) StaticCall(addr common.Address, input []byte, gas uint64) (ret []byte, err error) {
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, err = evm.RunPrecompiledAndConsumeGas(p, input, gas)
	} else {
		ret, err = evm.StateDB.StaticCall(addr, input, gas)
	}
	return ret, err
}

func (evm *EVM) RunPrecompiledAndConsumeGas(p PrecompiledContract, input []byte, suppliedGas uint64) ([]byte, error) {
	ret, remainingGas, err := RunPrecompiledContract(evm, p, input, suppliedGas)

	usedGas := suppliedGas - remainingGas
	precompileIdentifier := fmt.Sprintf("%T", p)
	if !evm.StateDB.UseGas(precompileIdentifier, usedGas) {
		err = ErrOutOfGas
		evm.StateDB.FailExecution(err)
	}

	return ret, err
}
