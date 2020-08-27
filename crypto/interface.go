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
	BLSVerify(key []byte,  msg []byte, sig []byte) error
}

type Ed25519 interface {
	Ed25519Verify(key []byte,  msg []byte, sig []byte) error
}

type Secp256k1 interface {
	Secp256k1Verify(key []byte,  msg []byte, sig []byte) error
}


// VMCrypto will provide the interface to the main crypto functionalities of the vm
type VMCrypto interface {
	Hasher
	Ed25519
	BLS
	Secp256k1
}
