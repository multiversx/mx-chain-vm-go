package arwenmandos

import (
	"bytes"
	"encoding/hex"
	"fmt"

	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/test/test-util/mandos/json/model"
)

func checkAccounts(
	checkAccounts *mj.CheckAccounts,
	world *worldhook.BlockchainHookMock,
) error {

	if !checkAccounts.OtherAccountsAllowed {
		for worldAcctAddr := range world.AcctMap {
			postAcctMatch := mj.FindCheckAccount(checkAccounts.Accounts, []byte(worldAcctAddr))
			if postAcctMatch == nil {
				return fmt.Errorf("unexpected account address: %s", hex.EncodeToString([]byte(worldAcctAddr)))
			}
		}
	}

	for _, expectedAcct := range checkAccounts.Accounts {
		matchingAcct, isMatch := world.AcctMap[string(expectedAcct.Address.Value)]
		if !isMatch {
			return fmt.Errorf("account %s expected but not found after running test",
				hex.EncodeToString(expectedAcct.Address.Value))
		}

		if !bytes.Equal(matchingAcct.Address, expectedAcct.Address.Value) {
			return fmt.Errorf("bad account address %s", hex.EncodeToString(matchingAcct.Address))
		}

		if !expectedAcct.Nonce.Check(matchingAcct.Nonce) {
			return fmt.Errorf("bad account nonce. Account: %s. Want: %s. Have: %d",
				hex.EncodeToString(matchingAcct.Address), expectedAcct.Nonce.Original, matchingAcct.Nonce)
		}

		if !expectedAcct.Balance.Check(matchingAcct.Balance) {
			return fmt.Errorf("bad account balance. Account: %s. Want: %s. Have: %s",
				hex.EncodeToString(matchingAcct.Address), expectedAcct.Balance.Original, bigIntPretty(matchingAcct.Balance))
		}

		if !expectedAcct.Code.Check(matchingAcct.Code) {
			return fmt.Errorf("bad account code. Account: %s. Want: [%s]. Have: [%s]",
				hex.EncodeToString(matchingAcct.Address), expectedAcct.Code.Original, string(matchingAcct.Code))
		}

		if !expectedAcct.AsyncCallData.Check([]byte(matchingAcct.AsyncCallData)) {
			return fmt.Errorf("bad async call data. Account: %s. Want: [%s]. Have: [%s]",
				hex.EncodeToString(matchingAcct.Address), expectedAcct.AsyncCallData.Original, matchingAcct.AsyncCallData)
		}

		err := checkAccountStorage(expectedAcct, matchingAcct)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkAccountStorage(expectedAcct *mj.CheckAccount, matchingAcct *worldhook.Account) error {
	if expectedAcct.IgnoreStorage {
		return nil
	}

	expectedStorage := make(map[string][]byte)
	for _, stkvp := range expectedAcct.CheckStorage {
		expectedStorage[string(stkvp.Key.Value)] = stkvp.Value.Value
	}

	allKeys := make(map[string]bool)
	for k := range expectedStorage {
		allKeys[k] = true
	}
	for k := range matchingAcct.Storage {
		allKeys[k] = true
	}
	storageError := ""
	for k := range allKeys {
		want, _ := expectedStorage[k]
		have := matchingAcct.StorageValue(k)
		if !bytes.Equal(want, have) {
			storageError += fmt.Sprintf(
				"\n  for key %s: Want: %s. Have: %s",
				byteArrayPretty([]byte(k)), byteArrayPretty(want), byteArrayPretty(have))
		}
	}
	if len(storageError) > 0 {
		return fmt.Errorf("wrong account storage for account 0x%s:%s",
			expectedAcct.Address.Original, storageError)
	}
	return nil
}
