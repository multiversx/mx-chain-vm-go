package ed25519

import (
	libed25519 "crypto/ed25519"

	"github.com/ElrondNetwork/wasm-vm/crypto/signing"
)

type ed25519 struct {
}

func NewEd25519Signer() *ed25519 {
	return &ed25519{}
}

func (e *ed25519) VerifyEd25519(key []byte, msg []byte, sig []byte) error {
	if len(key) != libed25519.PublicKeySize {
		return signing.ErrInvalidPublicKey
	}

	isValidSig := libed25519.Verify(key, msg, sig)
	if !isValidSig {
		return signing.ErrInvalidSignature
	}

	return nil
}
