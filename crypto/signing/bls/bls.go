package bls

import (
	"github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl"
	mclMultiSig "github.com/multiversx/mx-chain-crypto-go/signing/mcl/multisig"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl/singlesig"
	"github.com/multiversx/mx-chain-crypto-go/signing/multisig"
)

type bls struct {
	keyGenerator crypto.KeyGenerator
	signer       crypto.SingleSigner

	multiSigner crypto.MultiSigner
}

// NewBLS returns the component able to verify BLS signatures
func NewBLS() (*bls, error) {
	b := &bls{}
	suite := mcl.NewSuiteBLS12()
	b.keyGenerator = signing.NewKeyGenerator(suite)
	b.signer = singlesig.NewBlsSigner()

	var err error
	b.multiSigner, err = multisig.NewBLSMultisig(&mclMultiSig.BlsMultiSignerKOSK{}, b.keyGenerator)

	return b, err
}

// VerifyBLS verifies a BLS signatures
func (b *bls) VerifyBLS(key []byte, msg []byte, sig []byte) error {
	publicKey, err := b.keyGenerator.PublicKeyFromByteArray(key)
	if err != nil {
		return err
	}

	return b.signer.Verify(publicKey, msg, sig)
}

// VerifySignatureShare verifies signature share of BLS MultiSig
func (b *bls) VerifySignatureShare(publicKey []byte, message []byte, sig []byte) error {
	return b.multiSigner.VerifySignatureShare(publicKey, message, sig)
}

// VerifyAggregatedSig verifies aggregated signature of BLS MultiSig
func (b *bls) VerifyAggregatedSig(pubKeysSigners [][]byte, message []byte, aggSig []byte) error {
	return b.multiSigner.VerifyAggregatedSig(pubKeysSigners, message, aggSig)
}
