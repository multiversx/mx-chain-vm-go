package crypto

// Hasher defines the functionality of a component able to generate hashes
type Hasher interface {
	Sha256(data []byte) ([]byte, error)
	Keccak256(data []byte) ([]byte, error)
	Ripemd160(data []byte) ([]byte, error)
}

// BLS defines the functionality of a component able to verify BLS signatures
type BLS interface {
	VerifyBLS(key []byte, msg []byte, sig []byte) error
	VerifySignatureShare(publicKey []byte, message []byte, sig []byte) error
	VerifyAggregatedSig(pubKeysSigners [][]byte, message []byte, aggSig []byte) error
}

// Ed25519 defines the functionality of a component able to verify Ed25519 signatures
type Ed25519 interface {
	VerifyEd25519(key []byte, msg []byte, sig []byte) error
}

// Secp256 defines the functionality of a component able to verify and encode Secp256 signatures
type Secp256 interface {
	VerifySecp256k1(key []byte, msg []byte, sig []byte, hashType uint8) error
	EncodeSecp256k1DERSignature(r, s []byte) []byte
	VerifySecp256r1(key []byte, msg []byte, sig []byte) error
}

// VMCrypto will provide the interface to the main crypto functionalities of the vm
type VMCrypto interface {
	Hasher
	Ed25519
	BLS
	Secp256
}
