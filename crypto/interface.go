package crypto

type Hasher interface {
	// Sha256 cryptographic function
	Sha256(data []byte) ([]byte, error)

	// Keccak256 cryptographic function
	Keccak256(data []byte) ([]byte, error)

	// Ripemd160 cryptographic function
	Ripemd160(data []byte) ([]byte, error)
}

type BLS interface {
	VerifyBLS(key []byte, msg []byte, sig []byte) error
}

type Ed25519 interface {
	VerifyEd25519(key []byte, msg []byte, sig []byte) error
}

type Secp256k1 interface {
	VerifySecp256k1(key []byte, msg []byte, sig []byte) error
}

// VMCrypto will provide the interface to the main crypto functionalities of the vm
type VMCrypto interface {
	Hasher
	Ed25519
	BLS
	Secp256k1
}
