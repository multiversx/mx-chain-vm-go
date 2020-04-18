package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func checkAccounts(
	expectedAccounts []*ij.Account,
	world *worldhook.BlockchainHookMock,
) error {

	for worldAcctAddr := range world.AcctMap {
		postAcctMatch := ij.FindAccount(expectedAccounts, []byte(worldAcctAddr))
		if postAcctMatch == nil {
			return fmt.Errorf("unexpected account address: %s", hex.EncodeToString([]byte(worldAcctAddr)))
		}
	}

	for _, postAcctFromTest := range expectedAccounts {
		postAcct := convertAccount(postAcctFromTest)
		matchingAcct, isMatch := world.AcctMap[string(postAcct.Address)]
		if !isMatch {
			return fmt.Errorf("account %s expected but not found after running test",
				hex.EncodeToString(postAcct.Address))
		}

		if !bytes.Equal(matchingAcct.Address, postAcct.Address) {
			return fmt.Errorf("bad account address %s", hex.EncodeToString(matchingAcct.Address))
		}

		if matchingAcct.Nonce != postAcct.Nonce {
			return fmt.Errorf("bad account nonce. Account: %s. Want: %d. Have: %d",
				hex.EncodeToString(matchingAcct.Address), postAcct.Nonce, matchingAcct.Nonce)
		}

		if matchingAcct.Balance.Cmp(postAcct.Balance) != 0 {
			return fmt.Errorf("bad account balance. Account: %s. Want: %s. Have: %s",
				hex.EncodeToString(matchingAcct.Address), bigIntPretty(postAcct.Balance), bigIntPretty(matchingAcct.Balance))
		}

		if !bytes.Equal(matchingAcct.Code, postAcct.Code) {
			return fmt.Errorf("bad account code. Account: %s. Want: [%s]. Have: [%s]",
				hex.EncodeToString(matchingAcct.Address), postAcct.Code, matchingAcct.Code)
		}

		if matchingAcct.AsyncCallData != postAcct.AsyncCallData {
			return fmt.Errorf("bad async call data. Account: %s. Want: [%s]. Have: [%s]",
				hex.EncodeToString(matchingAcct.Address), postAcct.AsyncCallData, matchingAcct.AsyncCallData)
		}

		// compare storages
		allKeys := make(map[string]bool)
		for k := range postAcct.Storage {
			allKeys[k] = true
		}
		for k := range matchingAcct.Storage {
			allKeys[k] = true
		}
		storageError := ""
		for k := range allKeys {
			want := postAcct.StorageValue(k)
			have := matchingAcct.StorageValue(k)
			if !bytes.Equal(want, have) {
				storageError += fmt.Sprintf(
					"\n  for key %s: Want: %s. Have: %s",
					hex.EncodeToString([]byte(k)), byteArrayPretty(want), byteArrayPretty(have))
			}
		}
		if len(storageError) > 0 {
			return fmt.Errorf("wrong account storage for account 0x%s:%s",
				hex.EncodeToString(postAcct.Address), storageError)
		}
	}

	return nil
}
