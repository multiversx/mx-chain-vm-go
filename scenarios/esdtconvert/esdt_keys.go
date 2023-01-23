package esdtconvert

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

// esdtTokenKeyPrefix is the prefix of storage keys belonging to ESDT tokens.
var esdtTokenKeyPrefix = []byte(core.ProtectedKeyPrefix + core.ESDTKeyIdentifier)

// esdtRoleKeyPrefix is the prefix of storage keys belonging to ESDT roles.
var esdtRoleKeyPrefix = []byte(core.ProtectedKeyPrefix + core.ESDTRoleIdentifier + core.ESDTKeyIdentifier)

// esdtNonceKeyPrefix is the prefix of storage keys belonging to ESDT nonces.
var esdtNonceKeyPrefix = []byte(core.ProtectedKeyPrefix + core.ESDTNFTLatestNonceIdentifier)

// esdtDataMarshalizer is the global marshalizer to be used for encoding/decoding ESDT data
var esdtDataMarshalizer = &marshal.GogoProtoMarshalizer{}

// errNegativeValue signals that a negative value has been detected and it is not allowed
var errNegativeValue = errors.New("negative value")

// makeTokenKey creates the storage key corresponding to the given tokenName.
func makeTokenKey(tokenName []byte, nonce uint64) []byte {
	nonceBytes := big.NewInt(0).SetUint64(nonce).Bytes()
	tokenKey := append(esdtTokenKeyPrefix, tokenName...)
	tokenKey = append(tokenKey, nonceBytes...)
	return tokenKey
}

// makeTokenRolesKey creates the storage key corresponding to the roles for the
// given tokenName.
func makeTokenRolesKey(tokenName []byte) []byte {
	tokenRolesKey := append(esdtRoleKeyPrefix, tokenName...)
	return tokenRolesKey
}

// makeLastNonceKey creates the storage key corresponding to the last nonce of
// the given tokenName.
func makeLastNonceKey(tokenName []byte) []byte {
	tokenNonceKey := append(esdtNonceKeyPrefix, tokenName...)
	return tokenNonceKey
}

// isTokenKey returns true if the given storage key belongs to an ESDT token.
func isTokenKey(key []byte) bool {
	return bytes.HasPrefix(key, esdtTokenKeyPrefix)
}

// isRoleKey returns true if the given storage key belongs to an ESDT role.
func isRoleKey(key []byte) bool {
	return bytes.HasPrefix(key, esdtRoleKeyPrefix)
}

// isNonceKey returns true if the given storage key belongs to an ESDT nonce.
func isNonceKey(key []byte) bool {
	return bytes.HasPrefix(key, esdtNonceKeyPrefix)
}

// getTokenNameFromKey extracts the token name from the given storage key; it
// does not check whether the key is indeed a token key or not.
func getTokenNameFromKey(key []byte) []byte {
	return key[len(esdtTokenKeyPrefix):]
}
