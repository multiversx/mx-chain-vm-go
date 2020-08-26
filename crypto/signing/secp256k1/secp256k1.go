package secp256k1

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type secp256k1 struct {
}

func NewSecp256k1() *secp256k1 {
	return &secp256k1{}
}

func (s *secp256k1) VerifySecp256k1(key []byte,  msg []byte, sig []byte) error {
	pubKey, err := btcec.ParsePubKey(key, btcec.S256())
	if err != nil {
		return err
	}

	signature, err := btcec.ParseSignature(sig, btcec.S256())
	if err != nil {
		return err
	}

	messageHash := chainhash.DoubleHashB(msg)
	verified := signature.Verify(messageHash, pubKey)

	if !verified {
		return crypto.ErrInvalidSignature
	}

	return nil
}

func (s *secp256k1) Ecrecover(hash []byte, recoveryID []byte, r []byte, s []byte) ([]byte, error) {

}
