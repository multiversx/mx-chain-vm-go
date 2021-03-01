package mandosjsonmodel

import "bytes"

// Account is a json object representing an account.
type Account struct {
	Address         JSONBytesFromString
	Shard           JSONUint64
	IsSmartContract bool
	Comment         string
	Nonce           JSONUint64
	Balance         JSONBigInt
	Storage         []*StorageKeyValuePair
	Code            JSONBytesFromString
	AsyncCallData   string
	ESDTData        []*ESDTData
}

// StorageKeyValuePair is a json key value pair in the storage map.
type StorageKeyValuePair struct {
	Key   JSONBytesFromString
	Value JSONBytesFromTree
}

// ESDTData models an account holding an ESDT token
type ESDTData struct {
	TokenName JSONBytesFromString
	Balance   JSONBigInt
	Frozen    JSONUint64
}

// CheckAccount is a json object representing checks for an account.
type CheckAccount struct {
	Address       JSONBytesFromString
	Comment       string
	Nonce         JSONCheckUint64
	Balance       JSONCheckBigInt
	IgnoreStorage bool
	CheckStorage  []*StorageKeyValuePair
	Code          JSONCheckBytes
	AsyncCallData JSONCheckBytes
	ESDTData      []*CheckESDTData
}

// CheckESDTData checks the ESDT tokens held by an account
type CheckESDTData struct {
	TokenName JSONBytesFromString
	Balance   JSONCheckBigInt
	Frozen    JSONCheckUint64
}

// CheckAccounts encodes rules to check mock accounts.
type CheckAccounts struct {
	OtherAccountsAllowed bool
	Accounts             []*CheckAccount
}

// FindAccount searches an account list by address.
func FindAccount(accounts []*Account, address []byte) *Account {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address.Value, address) {
			return acct
		}
	}
	return nil
}

// FindCheckAccount searches a check account list by address.
func FindCheckAccount(accounts []*CheckAccount, address []byte) *CheckAccount {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address.Value, address) {
			return acct
		}
	}
	return nil
}
