package factory

import (
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/crypto/hashing"
	"github.com/multiversx/mx-chain-vm-go/crypto/signing/bls"
	"github.com/multiversx/mx-chain-vm-go/crypto/signing/ed25519"
	"github.com/multiversx/mx-chain-vm-go/crypto/signing/secp256k1"
)

// NewVMCrypto returns a composite struct containing VMCrypto functionality implementations
func NewVMCrypto() (crypto.VMCrypto, error) {
	blsVerifier, err := bls.NewBLS()
	if err != nil {
		return nil, err
	}

	return struct {
		crypto.Hasher
		crypto.Ed25519
		crypto.BLS
		crypto.Secp256
	}{
		Hasher:  hashing.NewHasher(),
		Ed25519: ed25519.NewEd25519Signer(),
		BLS:     blsVerifier,
		Secp256: secp256k1.NewSecp256k1(),
	}, nil
}
