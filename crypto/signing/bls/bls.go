package bls

import (
	"github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl/singlesig"
)

type bls struct {
	keyGenerator crypto.KeyGenerator
	signer       crypto.SingleSigner
}

// NewBLS returns the component able to verify BLS signatures
func NewBLS() *bls {
	b := &bls{}
	suite := mcl.NewSuiteBLS12()
	b.keyGenerator = signing.NewKeyGenerator(suite)
	b.signer = singlesig.NewBlsSigner()

	return b
}

// VerifyBLS verifies a BLS signatures
func (b *bls) VerifyBLS(key []byte, msg []byte, sig []byte) error {
	publicKey, err := b.keyGenerator.PublicKeyFromByteArray(key)
	if err != nil {
		return err
	}

	return b.signer.Verify(publicKey, msg, sig)
}
