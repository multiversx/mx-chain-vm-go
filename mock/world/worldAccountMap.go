package worldmock

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// AccountMap is a map from address to Account, also implementing the
// AccountsAdapter interface
type AccountMap map[string]*Account

// NewAccountMap creates a new AccountMap instance.
func NewAccountMap() AccountMap {
	return make(AccountMap)
}

// CreateAccount instantiates an empty account for the given address.
func (am AccountMap) CreateAccount(address []byte) *Account {
	newAccount := &Account{
		Exists:          true,
		Address:         make([]byte, len(address)),
		Nonce:           0,
		Balance:         big.NewInt(0),
		BalanceDelta:    big.NewInt(0),
		Storage:         make(map[string][]byte),
		Code:            nil,
		OwnerAddress:    nil,
		ShardID:         0,
		IsSmartContract: false,
		DeveloperReward: big.NewInt(0),
	}
	copy(newAccount.Address, address)
	am.PutAccount(newAccount)

	return newAccount
}

// CreateSmartContractAccount instantiates an account for a smart contract with
// the given address and WASM bytecode.
func (am AccountMap) CreateSmartContractAccount(owner []byte, address []byte, code []byte) *Account {
	newAccount := am.CreateAccount(address)
	newAccount.Code = code
	newAccount.IsSmartContract = true
	newAccount.OwnerAddress = owner
	newAccount.CodeMetadata = []byte{0, vmcommon.MetadataPayable}

	return newAccount
}

// PutAccount inserts account based on address.
func (am AccountMap) PutAccount(account *Account) {
	am[string(account.Address)] = account
}

// PutAccounts inserts multiple accounts based on address.
func (am AccountMap) PutAccounts(accounts []*Account) {
	for _, account := range accounts {
		am.PutAccount(account)
	}
}

// GetAccount retrieves account based on address
func (am AccountMap) GetAccount(address []byte) *Account {
	return am[string(address)]
}

// DeleteAccount removes account based on address
func (am AccountMap) DeleteAccount(address []byte) {
	delete(am, string(address))
}

func (am AccountMap) Clone() AccountMap {
	clone := make(AccountMap, len(am))
	for address, account := range am {
		clone[address] = account.Clone()
	}

	return clone
}

func (am AccountMap) LoadAccountStorageFrom(otherAM AccountMap) error {
	for address, account := range am {
		otherAccount, otherExists := otherAM[address]
		if !otherExists {
			return fmt.Errorf(
				"account %s could not be loaded from AccountMap",
				hex.EncodeToString([]byte(address)))
		}
		account.Storage = otherAccount.Storage
	}

	return nil
}
