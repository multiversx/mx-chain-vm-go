package mandosjsonwrite

import (
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func accountsToOJ(accounts []*mj.Account) oj.OJsonObject {
	acctsOJ := oj.NewMap()
	for _, account := range accounts {
		acctOJ := oj.NewMap()
		if len(account.Comment) > 0 {
			acctOJ.Put("comment", stringToOJ(account.Comment))
		}
		acctOJ.Put("nonce", uint64ToOJ(account.Nonce))
		acctOJ.Put("balance", bigIntToOJ(account.Balance))
		appendESDTToOJ(account.ESDTData, acctOJ)
		storageOJ := oj.NewMap()
		for _, st := range account.Storage {
			storageOJ.Put(bytesFromStringToString(st.Key), bytesFromTreeToOJ(st.Value))
		}
		acctOJ.Put("storage", storageOJ)
		acctOJ.Put("code", bytesFromStringToOJ(account.Code))
		if len(account.AsyncCallData) > 0 {
			acctOJ.Put("asyncCallData", stringToOJ(account.AsyncCallData))
		}

		acctsOJ.Put(bytesFromStringToString(account.Address), acctOJ)
	}

	return acctsOJ
}

func checkAccountsToOJ(checkAccounts *mj.CheckAccounts) oj.OJsonObject {
	acctsOJ := oj.NewMap()
	for _, checkAccount := range checkAccounts.Accounts {
		acctOJ := oj.NewMap()
		if len(checkAccount.Comment) > 0 {
			acctOJ.Put("comment", stringToOJ(checkAccount.Comment))
		}
		if !checkAccount.Nonce.IsUnspecified() {
			acctOJ.Put("nonce", checkUint64ToOJ(checkAccount.Nonce))
		}
		if !checkAccount.Balance.IsUnspecified() {
			acctOJ.Put("balance", checkBigIntToOJ(checkAccount.Balance))
		}
		storageOJ := oj.NewMap()
		for _, st := range checkAccount.CheckStorage {
			storageOJ.Put(bytesFromStringToString(st.Key), bytesFromTreeToOJ(st.Value))
		}
		if checkAccount.IgnoreStorage {
			acctOJ.Put("storage", stringToOJ("*"))
		} else {
			acctOJ.Put("storage", storageOJ)
		}
		if !checkAccount.Code.IsUnspecified() {
			acctOJ.Put("code", checkBytesToOJ(checkAccount.Code))
		}
		if !checkAccount.AsyncCallData.IsUnspecified() {
			acctOJ.Put("asyncCallData", checkBytesToOJ(checkAccount.AsyncCallData))
		}

		acctsOJ.Put(bytesFromStringToString(checkAccount.Address), acctOJ)
	}

	if checkAccounts.OtherAccountsAllowed {
		acctsOJ.Put("+", stringToOJ(""))
	}

	return acctsOJ
}
