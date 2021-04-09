package arwenmandos

import (
	"bytes"
	"encoding/hex"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/builtInFunctions"
)

func (ae *ArwenTestExecutor) ExecuteCheckStateStep(step *mj.CheckStateStep) error {
	if len(step.Comment) > 0 {
		log.Trace("CheckStateStep", "comment", step.Comment)
	}

	return ae.checkAccounts(step.CheckAccounts)
}

func (ae *ArwenTestExecutor) checkAccounts(checkAccounts *mj.CheckAccounts) error {
	if !checkAccounts.OtherAccountsAllowed {
		for worldAcctAddr := range ae.World.AcctMap {
			postAcctMatch := mj.FindCheckAccount(checkAccounts.Accounts, []byte(worldAcctAddr))
			if postAcctMatch == nil {
				return fmt.Errorf("unexpected account address: %s", hex.EncodeToString([]byte(worldAcctAddr)))
			}
		}
	}

	for _, expectedAcct := range checkAccounts.Accounts {
		matchingAcct, isMatch := ae.World.AcctMap[string(expectedAcct.Address.Value)]
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

func checkAccountStorage(expectedAcct *mj.CheckAccount, matchingAcct *worldmock.Account) error {
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
		if !bytes.Equal(want, have) && !worldmock.IsESDTKey([]byte(k)) {
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

func checkAccountESDT(expectedAcct *mj.CheckAccount, matchingAcct *worldmock.Account) error {
	if expectedAcct.IgnoreESDT {
		return nil
	}

	accountAddress := expectedAcct.Address.Original
	expectedTokens := getExpectedTokens(expectedAcct)
	accountTokens, err := matchingAcct.GetAllTokenData()
	if err != nil {
		return err
	}

	err = detectUnexpectedTokens(expectedTokens, accountTokens)
	if err != nil {
		return fmt.Errorf("mismatch for account %s: %w", accountAddress, err)
	}

	err = detectMissingTokens(expectedTokens, accountTokens)
	if err != nil {
		return fmt.Errorf("mismatch for account %s: %w", accountAddress, err)
	}

	errors := checkTokensState(expectedTokens, accountTokens)
	errorString := makeErrorString(errors)
	if len(errorString) > 0 {
		return fmt.Errorf("mismatch for account %s: %s", accountAddress, errorString)
	}

	return nil
}

func checkTokensState(
	expectedTokens map[string]*mj.CheckESDTData,
	accountTokens map[string]*esdt.ESDigitalToken,
) []error {
	errors := make([]error, 0)
	for tokenName := range accountTokens {
		expectedTokenData := expectedTokens[tokenName]
		accountTokenData := accountTokens[tokenName]
		err := checkTokenState(tokenName, expectedTokenData, accountTokenData)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func getExpectedTokens(expectedAcct *mj.CheckAccount) map[string]*mj.CheckESDTData {
	expectedTokens := make(map[string]*mj.CheckESDTData)
	for _, expectedTokenData := range expectedAcct.CheckESDTData {
		tokenName := expectedTokenData.TokenIdentifier.Value
		tokenNonce := expectedTokenData.Nonce.Value
		tokenKeyStr := string(worldmock.MakeTokenKey(tokenName, tokenNonce))

		expectedTokens[tokenKeyStr] = expectedTokenData
	}

	return expectedTokens
}

func detectUnexpectedTokens(
	expectedTokens map[string]*mj.CheckESDTData,
	accountTokens map[string]*esdt.ESDigitalToken,
) error {
	for tokenName, accountTokenData := range accountTokens {
		_, isExpected := expectedTokens[tokenName]
		if !isExpected && accountTokenData.Value.Sign() > 0 {
			return fmt.Errorf("unexpected ESDT token `%s` found on the account", tokenName)
		}
	}

	return nil
}

func detectMissingTokens(
	expectedTokens map[string]*mj.CheckESDTData,
	accountTokens map[string]*esdt.ESDigitalToken,
) error {
	for tokenName, expectedTokenData := range expectedTokens {
		_, isFound := accountTokens[tokenName]
		if !isFound && expectedTokenData.Value.Value.Sign() > 0 {
			return fmt.Errorf("missing ESDT token %ss", tokenName)
		}
	}

	return nil
}

func checkTokenState(tokenName string, expectedTokenData *mj.CheckESDTData, accountTokenData *esdt.ESDigitalToken) error {
	if expectedTokenData == nil {
		if accountTokenData.Value.Sign() != 0 {
			return fmt.Errorf("bad ESDT balance. Token %s: Want: %d. Have: %d",
				tokenName, 0, accountTokenData.Value)
		}

		return nil
	}

	if !expectedTokenData.Value.Check(accountTokenData.Value) {
		return fmt.Errorf("bad ESDT balance. Token %s: Want: %d. Have: %d",
			tokenName, expectedTokenData.Value.Value, accountTokenData.Value)
	}

	metadataFromBytes := builtInFunctions.ESDTUserMetadataFromBytes(accountTokenData.Properties)
	if !expectedTokenData.Frozen.CheckBool(metadataFromBytes.Frozen) {
		return fmt.Errorf("bad ESDT frozen flag. Token %s: Want: %t. Have: %t",
			tokenName, expectedTokenData.Frozen.Value > 0, metadataFromBytes.Frozen)
	}

	return nil
}

func makeErrorString(errors []error) string {
	errorString := ""
	for i, err := range errors {
		errorString += err.Error()
		if i < len(errors)-1 {
			errorString += "\n"
		}
	}
	return errorString
}
