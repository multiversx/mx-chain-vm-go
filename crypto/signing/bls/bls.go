package bls

import (
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/mcl"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/mcl/singlesig"
)

type bls struct {
}

func NewBLS() *bls {
	return &bls{}
}

func (b *bls) BLSVerify(key []byte,  msg []byte, sig []byte) error {
	suite := mcl.NewSuiteBLS12()
	keyGenerator := signing.NewKeyGenerator(suite)

	publicKey, err := keyGenerator.PublicKeyFromByteArray(key)
	if err != nil {
		return err
	}

	signer := singlesig.NewBlsSigner()

	return signer.Verify(publicKey, msg, sig)
}
