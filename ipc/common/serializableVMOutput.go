package common

import (
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type SerializableVMOutput struct {
	ReturnData              [][]byte
	ReturnCode              vmcommon.ReturnCode
	ReturnMessage           string
	GasRemaining            uint64
	GasRefund               *big.Int
	CorrectedOutputAccounts []*SerializableOutputAccount
	DeletedAccounts         [][]byte
	TouchedAccounts         [][]byte
	Logs                    []*vmcommon.LogEntry
}

func NewSerializableVMOutput(vmOutput *vmcommon.VMOutput) *SerializableVMOutput {
	if vmOutput == nil {
		return &SerializableVMOutput{}
	}

	o := &SerializableVMOutput{
		ReturnData:              vmOutput.ReturnData,
		ReturnCode:              vmOutput.ReturnCode,
		ReturnMessage:           vmOutput.ReturnMessage,
		GasRemaining:            vmOutput.GasRemaining,
		GasRefund:               vmOutput.GasRefund,
		CorrectedOutputAccounts: make([]*SerializableOutputAccount, 0, len(vmOutput.OutputAccounts)),
		DeletedAccounts:         vmOutput.DeletedAccounts,
		TouchedAccounts:         vmOutput.TouchedAccounts,
		Logs:                    vmOutput.Logs,
	}

	for _, account := range vmOutput.OutputAccounts {
		o.CorrectedOutputAccounts = append(o.CorrectedOutputAccounts, NewSerializableOutputAccount(account))
	}

	return o
}

func (o *SerializableVMOutput) ConvertToVMOutput() *vmcommon.VMOutput {
	accountsMap := make(map[string]*vmcommon.OutputAccount)

	for _, item := range o.CorrectedOutputAccounts {
		accountsMap[string(item.Address)] = item.ConvertToOutputAccount()
	}

	return &vmcommon.VMOutput{
		ReturnData:      o.ReturnData,
		ReturnCode:      o.ReturnCode,
		ReturnMessage:   o.ReturnMessage,
		GasRemaining:    o.GasRemaining,
		GasRefund:       o.GasRefund,
		OutputAccounts:  accountsMap,
		DeletedAccounts: o.DeletedAccounts,
		TouchedAccounts: o.TouchedAccounts,
		Logs:            o.Logs,
	}
}

type SerializableOutputAccount struct {
	Address        []byte
	Nonce          uint64
	Balance        *big.Int
	BalanceDelta   *big.Int
	StorageUpdates []*vmcommon.StorageUpdate
	Code           []byte
	CodeMetadata   []byte
	Data           [][]byte
	GasLimit       uint64
	CallType       vmcommon.CallType
}

func NewSerializableOutputAccount(account *vmcommon.OutputAccount) *SerializableOutputAccount {
	a := &SerializableOutputAccount{
		Address:        account.Address,
		Nonce:          account.Nonce,
		Balance:        account.Balance,
		BalanceDelta:   account.BalanceDelta,
		StorageUpdates: make([]*vmcommon.StorageUpdate, 0, len(account.StorageUpdates)),
		Code:           account.Code,
		CodeMetadata:   account.CodeMetadata,
		Data:           account.Data,
		GasLimit:       account.GasLimit,
		CallType:       account.CallType,
	}

	for _, storageUpdate := range account.StorageUpdates {
		a.StorageUpdates = append(a.StorageUpdates, storageUpdate)
	}

	return a
}

func (a *SerializableOutputAccount) ConvertToOutputAccount() *vmcommon.OutputAccount {
	updatesMap := make(map[string]*vmcommon.StorageUpdate)

	for _, item := range a.StorageUpdates {
		updatesMap[string(item.Offset)] = item
	}

	return &vmcommon.OutputAccount{
		Address:        a.Address,
		Nonce:          a.Nonce,
		Balance:        a.Balance,
		BalanceDelta:   a.BalanceDelta,
		StorageUpdates: updatesMap,
		Code:           a.Code,
		CodeMetadata:   a.CodeMetadata,
		Data:           a.Data,
		GasLimit:       a.GasLimit,
		CallType:       a.CallType,
	}
}
