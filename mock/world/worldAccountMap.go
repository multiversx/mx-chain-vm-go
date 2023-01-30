package worldmock

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/crypto/hashing"
)

// AccountMap is a map from address to Account, also implementing the
// AccountsAdapter interface
type AccountMap map[string]*Account

// NewAccountMap creates a new AccountMap instance.
func NewAccountMap() AccountMap {
	return make(AccountMap)
}

// CreateAccount instantiates an empty account for the given address.
func (am AccountMap) CreateAccount(address []byte, world *MockWorld) *Account {
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
		MockWorld:       world,
	}
	copy(newAccount.Address, address)
	am.PutAccount(newAccount)

	return newAccount
}

// CreateSmartContractAccount instantiates an account for a smart contract with
// the given address and WASM bytecode.
func (am AccountMap) CreateSmartContractAccount(owner []byte, address []byte, code []byte, world *MockWorld) *Account {
	return am.CreateSmartContractAccountWithCodeHash(owner, address, code, code, world)
}

// CreateSmartContractAccountWithCodeHash instantiates an account for a smart contract with
// the given address and WASM bytecode.
func (am AccountMap) CreateSmartContractAccountWithCodeHash(owner []byte, address []byte, code []byte, codeHash []byte, world *MockWorld) *Account {
	newAccount := am.CreateAccount(address, world)
	newAccount.Code = code
	if codeHash == nil {
		codeHash = code
	}
	newAccount.CodeHash = codeHash
	newAccount.IsSmartContract = true
	newAccount.OwnerAddress = owner

	metadata := &vmcommon.CodeMetadata{
		Payable: true,
	}

	newAccount.SetCodeAndMetadata(code, metadata)

	return newAccount
}

// PutAccount inserts account based on address.
func (am AccountMap) PutAccount(account *Account) {
	if account.Code != nil && account.CodeHash == nil {
		hash, _ := hashing.NewHasher().Sha256(account.Code)
		account.CodeHash = hash
	}
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

// Clone creates a deep clone of the entire AccountMap.
func (am AccountMap) Clone() AccountMap {
	clone := make(AccountMap, len(am))
	for address, account := range am {
		clone[address] = account.Clone()
	}

	return clone
}

// LoadAccountStorageFrom reassigns the storage of the accounts to the storage
// of the accounts found in otherAM; it only does a reference change, not a deep copy.
func (am AccountMap) LoadAccountStorageFrom(otherAM AccountMap) error {
	for address, account := range am {
		otherAccount, otherExists := otherAM[address]
		if !otherExists {
			if bytes.Equal([]byte(address), vmcommon.SystemAccountAddress) {
				continue
			}

			return fmt.Errorf(
				"account %s could not be loaded from AccountMap",
				hex.EncodeToString([]byte(address)))
		}
		account.Storage = otherAccount.Storage
	}

	return nil
}
