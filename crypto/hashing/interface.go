package hashing

type Hasher interface {
	// Sha256 cryptographic function
	Sha256(data []byte) ([]byte, error)

	// Keccak256 cryptographic function
	Keccak256(data []byte) ([]byte, error)

	// Ripemd160 cryptographic function
	Ripemd160(data []byte) ([]byte, error)
}
