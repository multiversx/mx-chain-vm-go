package crypto

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/hashing"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/signing/bls"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/signing/ed25519"
)

// NewCryptoVm returns a composite struct containing VMCrypto functionality implementations
func NewCryptoVm() VMCrypto {
	return struct {
		hashing.Hasher
		ed25519.Ed25519
		bls.BLS
	}{
		hashing.NewHasher(),
		ed25519.NewEd25519Signer(),
		bls.NewBLS(),
	}
}
