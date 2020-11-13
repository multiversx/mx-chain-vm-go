package arwenpart

import (
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/sha3"

	vmcommon "github.com/ElrondNetwork/elrond-go/core/vm-common"
	"golang.org/x/crypto/ripemd160"
)

var _ vmcommon.CryptoHook = (*CryptoHookGateway)(nil)

// CryptoHookGateway is a copy of the CryptoHook implementation from the node
// TODO: Remove this implementation and reference ElrondNetwork/common/crypto when it becomes available
type CryptoHookGateway struct {
}

// NewCryptoHookGateway creates a new crypto hook gateway
func NewCryptoHookGateway() *CryptoHookGateway {
	return &CryptoHookGateway{}
}

// Sha256 returns a sha 256 hash of the input string. Should return in hex format
func (hook *CryptoHookGateway) Sha256(data []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	result := hash.Sum(nil)
	return result, nil
}

// Keccak256 returns a keccak 256 hash of the input string. Should return in hex format
func (hook *CryptoHookGateway) Keccak256(data []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	result := hash.Sum(nil)
	return result, nil
}

// Ripemd160 is a legacy hash and should not be used for new applications
func (hook *CryptoHookGateway) Ripemd160(data []byte) ([]byte, error) {
	hash := ripemd160.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	result := hash.Sum(nil)
	return result, nil
}

// Ecrecover calculates the corresponding Ethereum address for the public key which created the given signature
// https://ewasm.readthedocs.io/en/mkdocs/system_contracts/
func (hook *CryptoHookGateway) Ecrecover(_ []byte, _ []byte, _ []byte, _ []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}
