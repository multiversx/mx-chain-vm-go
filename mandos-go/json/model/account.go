package mandosjsonmodel

import "bytes"

// Account is a json object representing an account.
type Account struct {
	Address       JSONBytesFromString
	Comment       string
	Nonce         JSONUint64
	Balance       JSONBigInt
	Storage       []*StorageKeyValuePair
	Code          JSONBytesFromString
	AsyncCallData string
}

// StorageKeyValuePair is a json key value pair in the storage map.
type StorageKeyValuePair struct {
	Key   JSONBytesFromString
	Value JSONBytesFromTree
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
