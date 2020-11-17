package mockhookcrypto

import (
	"fmt"

	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// KryptoHookMock is a krypto hook implementation that we use for VM tests
type KryptoHookMock int

// KryptoHookMockInstance is a krypto hook mock singleton
const KryptoHookMockInstance KryptoHookMock = 0

// Sha256 cryptographic function
func (KryptoHookMock) Sha256(data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	result := hash.Sum(nil)
	return result, nil
}

// Keccak256 cryptographic function
func (KryptoHookMock) Keccak256(data []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	result := hash.Sum(nil)
	return result, nil
}

// Ripemd160 cryptographic function
func (KryptoHookMock) Ripemd160(data []byte) ([]byte, error) {
	hash := ripemd160.New()
	hash.Write(data)
	result := hash.Sum(nil)
	return result, nil
}

// Ecrecover calculates the corresponding Ethereum address for the public key which created the given signature
func (KryptoHookMock) Ecrecover(_ []byte, _ []byte, _ []byte, _ []byte) ([]byte, error) {
	fmt.Println(">>>>> Ecrecover")
	return []byte{}, nil
}
