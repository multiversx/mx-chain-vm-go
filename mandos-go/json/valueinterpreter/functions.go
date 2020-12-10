package mandosvalueinterpreter

import (
	"golang.org/x/crypto/sha3"
)

// Keccak256 cryptographic function
// TODO: externalize the same way as the file resolver
func keccak256(data []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	result := hash.Sum(nil)
	return result, nil
}

// Generates a 32-byte address based on the input.
func address(data []byte) ([]byte, error) {
	if len(data) > 32 {
		return data[:32], nil
	}
	var result [32]byte
	i := 0
	for ; i < len(data); i++ {
		result[i] = data[i]
	}
	for ; i < 32; i++ {
		result[i] = byte('_')
	}
	return result[:], nil
}
