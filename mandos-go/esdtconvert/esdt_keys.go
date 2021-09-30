package esdtconvert

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
)

// ErrNegativeValue signals that a negative value has been detected and it is not allowed
var ErrNegativeValue = errors.New("negative value")

// ESDTTokenKeyPrefix is the prefix of storage keys belonging to ESDT tokens.
var ESDTTokenKeyPrefix = []byte(core.ElrondProtectedKeyPrefix + core.ESDTKeyIdentifier)

// ESDTRoleKeyPrefix is the prefix of storage keys belonging to ESDT roles.
var ESDTRoleKeyPrefix = []byte(core.ElrondProtectedKeyPrefix + core.ESDTRoleIdentifier + core.ESDTKeyIdentifier)

// ESDTNonceKeyPrefix is the prefix of storage keys belonging to ESDT nonces.
var ESDTNonceKeyPrefix = []byte(core.ElrondProtectedKeyPrefix + core.ESDTNFTLatestNonceIdentifier)

// esdtDataMarshalizer is the global marshalizer to be used for encoding/decoding ESDT data
var esdtDataMarshalizer = &marshal.GogoProtoMarshalizer{}

// MakeTokenKey creates the storage key corresponding to the given tokenName.
func MakeTokenKey(tokenName []byte, nonce uint64) []byte {
	nonceBytes := big.NewInt(0).SetUint64(nonce).Bytes()
	tokenKey := append(ESDTTokenKeyPrefix, tokenName...)
	tokenKey = append(tokenKey, nonceBytes...)
	return tokenKey
}

// MakeTokenRolesKey creates the storage key corresponding to the roles for the
// given tokenName.
func MakeTokenRolesKey(tokenName []byte) []byte {
	tokenRolesKey := append(ESDTRoleKeyPrefix, tokenName...)
	return tokenRolesKey
}

// MakeLastNonceKey creates the storage key corresponding to the last nonce of
// the given tokenName.
func MakeLastNonceKey(tokenName []byte) []byte {
	tokenNonceKey := append(ESDTNonceKeyPrefix, tokenName...)
	return tokenNonceKey
}

// IsESDTKey returns true if the given storage key is ESDT-related
func IsESDTKey(key []byte) bool {
	return IsTokenKey(key) || IsRoleKey(key) || IsNonceKey(key)
}

// IsTokenKey returns true if the given storage key belongs to an ESDT token.
func IsTokenKey(key []byte) bool {
	return bytes.HasPrefix(key, ESDTTokenKeyPrefix)
}

// IsRoleKey returns true if the given storage key belongs to an ESDT role.
func IsRoleKey(key []byte) bool {
	return bytes.HasPrefix(key, ESDTRoleKeyPrefix)
}

// IsNonceKey returns true if the given storage key belongs to an ESDT nonce.
func IsNonceKey(key []byte) bool {
	return bytes.HasPrefix(key, ESDTNonceKeyPrefix)
}

// GetTokenNameFromKey extracts the token name from the given storage key; it
// does not check whether the key is indeed a token key or not.
func GetTokenNameFromKey(key []byte) []byte {
	return key[len(ESDTTokenKeyPrefix):]
}
