// Copyright 2016 The go-ethereum Authors
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

type StateDB interface {
	FailExecution(err error)

	GetSelfBalance() *uint256.Int
	GetBalance(address common.Address) *uint256.Int
	TransferBalance(destination common.Address, value *uint256.Int) error

	GetCodeHash(address common.Address) common.Hash
	GetCode(address common.Address) []byte
	GetCodeSize(address common.Address) int

	GetState(key common.Hash) common.Hash
	SetState(key common.Hash, value common.Hash)

	GetTransientState(address common.Address, key common.Hash) common.Hash
	SetTransientState(address common.Address, key common.Hash, value common.Hash)

	SelfDestruct(destination common.Address)

	AddLog(log *types.Log)

	GasLeft() uint64
	UseGas(opCode string, gas uint64) bool

	IsSmartContractAddress(address common.Address) bool
	Create(code []byte, gas uint64, value *uint256.Int) ([]byte, common.Address, error)
	Create2(code []byte, gas uint64, value *uint256.Int, salt *uint256.Int) ([]byte, common.Address, error)
	Call(address common.Address, value *uint256.Int, input []byte, gas uint64) ([]byte, error)
	StaticCall(address common.Address, input []byte, gas uint64) ([]byte, error)
	CallCode(address common.Address, value *uint256.Int, input []byte, gas uint64) ([]byte, error)
	DelegateCall(address common.Address, input []byte, gas uint64) ([]byte, error)
}
