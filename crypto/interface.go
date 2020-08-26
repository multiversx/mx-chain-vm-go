package crypto

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/hashing"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/signing/bls"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/signing/ed25519"
)

// VMCrypto will provide the interface to the main crypto functionalities of the vm
type VMCrypto interface {
	hashing.Hasher
	ed25519.Ed25519
	bls.BLS
}
