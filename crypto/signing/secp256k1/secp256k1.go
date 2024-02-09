package secp256k1

import (
	cryptoEcsda "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"github.com/multiversx/mx-chain-vm-go/crypto/hashing"
	"github.com/multiversx/mx-chain-vm-go/crypto/signing"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// MessageHashType defines the type of hash algorithm
type MessageHashType uint8

// constants that define the available hash algorithms
const (
	ECDSAPlainMsg MessageHashType = iota
	ECDSASha256
	ECDSADoubleSha256
	ECDSAKeccak256
	ECDSARipemd160
)

const (
	// fieldSize is the curve domain size.
	fieldSize  = 32
	pubKeySize = fieldSize + 1

	name = "secp256r1"
)

// p256Order returns the curve order for the secp256r1 curve
// NOTE: this is specific to the secp256r1/P256 curve,
// and not taken from the domain params for the key itself
// (which would be a more generic approach for all EC).
var p256Order = elliptic.P256().Params().N

// p256HalfOrder returns half the curve order
// a bit shift of 1 to the right (Rsh) is equivalent
// to division by 2, only faster.
var p256HalfOrder = new(big.Int).Rsh(p256Order, 1)

// signatureR1 holds the r and s values of an ECDSA signature.
type signatureR1 struct {
	R, S *big.Int
}

type secp256k1 struct {
	secp256r1 elliptic.Curve
}

// NewSecp256k1 returns the component able to verify Secp256 signatures
func NewSecp256k1() (*secp256k1, error) {
	secp256r1 := elliptic.P256()

	expected := (secp256r1.Params().BitSize + 7) / 8
	if expected != fieldSize {
		return nil, errors.New("wrong secp256r1 curve")
	}
	return &secp256k1{
		secp256r1: secp256r1,
	}, nil
}

// VerifySecp256k1 checks a secp256k1 signature provided in the DER encoding format.
// The hash type used over the message can also be configured using @param hashType
func (sec *secp256k1) VerifySecp256k1(key, msg, sig []byte, hashType uint8) error {
	pubKey, err := btcec.ParsePubKey(key)
	if err != nil {
		return err
	}

	messageHash, err := sec.hashMessage(msg, hashType)
	if err != nil {
		return err
	}

	signature, err := ecdsa.ParseSignature(sig)
	if err != nil {
		return err
	}

	verified := signature.Verify(messageHash, pubKey)
	if !verified {
		return signing.ErrInvalidSignature
	}

	return nil
}

// EncodeSecp256k1DERSignature creates a DER encoding of a signature provided with r and s.
// Useful when having the plain params - like in the case of ecrecover
//
//	from ethereum
func (sec *secp256k1) EncodeSecp256k1DERSignature(r, s []byte) []byte {
	rScalar := &btcec.ModNScalar{}
	rScalar.SetByteSlice(r)

	sScalar := &btcec.ModNScalar{}
	sScalar.SetByteSlice(s)

	sig := ecdsa.NewSignature(rScalar, sScalar)

	return sig.Serialize()
}

func (sec *secp256k1) hashMessage(msg []byte, hashType uint8) ([]byte, error) {
	hasher := hashing.NewHasher()

	var err error
	var hashedMsg []byte
	switch MessageHashType(hashType) {
	case ECDSASha256:
		hashedMsg, err = hasher.Sha256(msg)
	case ECDSADoubleSha256:
		hashedMsg = chainhash.DoubleHashB(msg)
	case ECDSAKeccak256:
		hashedMsg, err = hasher.Keccak256(msg)
	case ECDSARipemd160:
		hashedMsg, err = hasher.Ripemd160(msg)
	case ECDSAPlainMsg:
		hashedMsg = msg
	default:
		return nil, signing.ErrHasherNotSupported
	}

	if err != nil {
		return nil, err
	}

	return hashedMsg, nil
}

func (sec *secp256k1) VerifySecp256r1(key []byte, msg []byte, sig []byte) error {
	if len(sig) != 64 {
		return errors.New("invalid signature length")
	}

	s := signatureFromBytes(sig)
	if !IsSNormalized(s.S) {
		return errors.New("signature not normalized")
	}

	h := sha256.Sum256(msg)

	verified := cryptoEcsda.Verify(&cpk, h[:], s.R, s.S)
	if !verified {
		return errors.New("signature verification failed")
	}

	return nil
}

func (sec *secp256k1) unmarshalPubKey(key []byte) (cryptoEcsda.PublicKey, error) {
	cpk := cryptoEcsda.PublicKey{Curve: sec.secp256r1}
	cpk.X, cpk.Y = elliptic.UnmarshalCompressed(sec.secp256r1, key)
	return cpk, nil
}

// signatureFromBytes function roughly copied from secp256k1_nocgo.go
// Read Signature struct from R || S. Caller needs to ensure that
// len(sigStr) == 64.
func signatureFromBytes(sigStr []byte) *signatureR1 {
	return &signatureR1{
		R: new(big.Int).SetBytes(sigStr[:32]),
		S: new(big.Int).SetBytes(sigStr[32:64]),
	}
}

// IsSNormalized returns true for the integer sigS if sigS falls in
// lower half of the curve order
func IsSNormalized(sigS *big.Int) bool {
	return sigS.Cmp(p256HalfOrder) != 1
}

// NormalizeS will invert the s value if not already in the lower half
// of curve order value
func NormalizeS(sigS *big.Int) *big.Int {
	if IsSNormalized(sigS) {
		return sigS
	}

	return new(big.Int).Sub(p256Order, sigS)
}
