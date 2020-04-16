package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func convertAccount(testAcct *ij.Account) *worldhook.Account {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		if stkvp.Value == nil {
			panic("why?")
		}
		key := string(stkvp.Key)
		storage[key] = stkvp.Value
	}

	if len(testAcct.Address) != 32 {
		panic("bad test: account address should be 32 bytes long")
	}

	return &worldhook.Account{
		Exists:        true,
		Address:       testAcct.Address,
		Nonce:         testAcct.Nonce.Uint64(),
		Balance:       big.NewInt(0).Set(testAcct.Balance),
		Storage:       storage,
		Code:          []byte(testAcct.Code),
		AsyncCallData: testAcct.AsyncCallData,
	}
}

func convertLogToTestFormat(outputLog *vmi.LogEntry) *ij.LogEntry {
	testLog := ij.LogEntry{
		Address:    outputLog.Address,
		Identifier: outputLog.Identifier,
		Data:       outputLog.Data,
		Topics:     outputLog.Topics,
	}

	return &testLog
}

func bigIntPretty(i *big.Int) string {
	return fmt.Sprintf("0x%x (%d)", i, i)
}

func byteArrayPretty(b []byte) string {
	if len(b) == 0 {
		return "[]"
	}
	asInt := big.NewInt(0).SetBytes(b)
	return fmt.Sprintf("0x%s (%d)", hex.EncodeToString(b), asInt)
}
