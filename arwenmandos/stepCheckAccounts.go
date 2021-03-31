package arwenmandos

import (
	"bytes"
	"encoding/hex"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
)

func checkAccounts(
	checkAccounts *mj.CheckAccounts,
	world *worldhook.MockWorld,
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

		// currently ignoring asyncCallData that is unspecified in the json
		if !expectedAcct.AsyncCallData.IsUnspecified() &&
			!expectedAcct.AsyncCallData.Check([]byte(matchingAcct.AsyncCallData)) {
			return fmt.Errorf("bad async call data. Account: %s. Want: [%s]. Have: [%s]",
				hex.EncodeToString(matchingAcct.Address), expectedAcct.AsyncCallData.Original, matchingAcct.AsyncCallData)
		}

		err := checkAccountStorage(expectedAcct, matchingAcct)
		if err != nil {
			return err
		}

		err = checkAccountESDT(expectedAcct, matchingAcct)
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
		want := expectedStorage[k]
		have := matchingAcct.StorageValue(k)
		if !bytes.Equal(want, have) {
			storageError += fmt.Sprintf(
				"\n  for key %s: Want: %s. Have: %s",
				byteArrayPretty([]byte(k)), byteArrayPretty(want), byteArrayPretty(have))
		}
	}
	if len(storageError) > 0 {
		return fmt.Errorf("wrong account storage for account \"%s\":%s",
			expectedAcct.Address.Original, storageError)
	}
	return nil
}

func checkAccountESDT(expectedAcct *mj.CheckAccount, matchingAcct *worldhook.Account) error {
	// check for unexpected tokens
	expectedTokenNames := make(map[string]bool)
	for _, expectedTokenData := range expectedAcct.ESDTData {
		tokenNameStr := string(expectedTokenData.TokenName.Value)
		expectedTokenNames[tokenNameStr] = true
	}
	for tokenName, tokenData := range matchingAcct.ESDTData {
		if tokenData.Balance.Sign() > 0 && !expectedTokenNames[tokenName] {
			return fmt.Errorf("unexpected ESDT token %s for account %s", tokenName, expectedAcct.Address.Original)
		}
	}

	esdtError := ""
	for _, expectedTokenData := range expectedAcct.ESDTData {
		tokenNameStr := string(expectedTokenData.TokenName.Value)
		have := matchingAcct.ESDTData[tokenNameStr]
		if !expectedTokenData.Balance.Check(have.Balance) {
			esdtError += fmt.Sprintf(
				"\n  bad ESDT balance. Token %s: Want: %d. Have: %d",
				tokenNameStr, expectedTokenData.Balance.Value, have.Balance)
		}

		if !expectedTokenData.Frozen.CheckBool(have.Frozen) {
			esdtError += fmt.Sprintf(
				"\n  bad ESDT frozen flag. Token %s: Want: %t. Have: %t",
				tokenNameStr, expectedTokenData.Frozen.Value > 0, have.Frozen)
		}
	}
	if len(esdtError) > 0 {
		return fmt.Errorf("wrong ESDT token state for account \"%s\":%s",
			expectedAcct.Address.Original, esdtError)
	}
	return nil
}
