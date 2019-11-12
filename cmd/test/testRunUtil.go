package main

import (
	"fmt"
	"math/big"

	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// for nicer error messages
func resultAsString(result []*big.Int) string {
	str := "["
	for i, res := range result {
		str += fmt.Sprintf("0x%x", res)
		if i < len(result)-1 {
			str += ", "
		}
	}
	return str + "]"
}

func convertAccount(testAcct *ij.Account) *worldhook.Account {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		if stkvp.Value == nil {
			panic("why?")
		}
		key := string(stkvp.Key)
		storage[key] = stkvp.Value
	}

	return &worldhook.Account{
		Exists:  true,
		Address: testAcct.Address,
		Nonce:   testAcct.Nonce.Uint64(),
		Balance: big.NewInt(0).Set(testAcct.Balance),
		Storage: storage,
		Code:    []byte(testAcct.Code),
	}
}

func convertLogToTestFormat(outputLog *vmi.LogEntry) *ij.LogEntry {
	testLog := ij.LogEntry{
		Address: outputLog.Address,
		Topics:  outputLog.Topics,
		Data:    outputLog.Data,
	}
	return &testLog
}

func convertBlockHeader(testBlh *ij.BlockHeader) *vmi.SCCallHeader {
	return &vmi.SCCallHeader{
		Beneficiary: testBlh.Beneficiary,
		Number:      testBlh.Number,
		GasLimit:    testBlh.GasLimit,
		Timestamp:   testBlh.UnixTimestamp,
	}
}

var zero = big.NewInt(0)

func zeroIfNil(i *big.Int) *big.Int {
	if i == nil {
		return zero
	}
	return i
}
