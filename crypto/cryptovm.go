package crypto

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/hashing"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/signing/bls"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/signing/ed25519"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/signing/secp256k1"
)

// NewCryptoVm returns a composite struct containing VMCrypto functionality implementations
func NewCryptoVm() VMCrypto {
	return struct {
		Hasher
		Ed25519
		BLS
		Secp256k1
	}{
		Hasher:    hashing.NewHasher(),
		Ed25519:   ed25519.NewEd25519Signer(),
		BLS:       bls.NewBLS(),
		Secp256k1: secp256k1.NewSecp256k1(),
	}
}
