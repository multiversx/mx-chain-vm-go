package hashing

import (
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

type hasher struct {

}

// NewHasher returns a new hasher instance implementing wrappers over different hash functions
func NewHasher() Hasher {
	return &hasher{}
}

// Sha256 returns a sha 256 hash of the input string. Should return in hex format
func (h *hasher) Sha256(data []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	result := hash.Sum(nil)
	return result, nil
}

// Keccak256 returns a keccak 256 hash of the input string. Should return in hex format
func (h *hasher) Keccak256(data []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	result := hash.Sum(nil)
	return result, nil
}

// Ripemd160 is a legacy hash and should not be used for new applications
func (h *hasher) Ripemd160(data []byte) ([]byte, error) {
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
func (h *hasher) Ecrecover(_, _, _, _ []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}


