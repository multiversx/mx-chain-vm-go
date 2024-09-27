package secp256

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEthereumSig(t *testing.T) {
	t.Parallel()

	msg, _ := hex.DecodeString("ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	r, _ := hex.DecodeString("90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e54998")
	s, _ := hex.DecodeString("4a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc93")
	key, _ := hex.DecodeString("04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")

	verifier, _ := NewSecp256()
	sig := verifier.EncodeSecp256k1DERSignature(r, s)
	err := verifier.VerifySecp256k1(key, msg, sig, byte(ECDSAPlainMsg))

	assert.Nil(t, err)
}

func TestWrongSizeForRSShouldNotPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r != nil {
			assert.Fail(t, fmt.Sprintf("should have not panicked %v", r))
		}
	}()

	r, _ := hex.DecodeString("90f2")                                                                 // too short
	s, _ := hex.DecodeString("4a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc939393") // too long

	verifier, _ := NewSecp256()
	sig := verifier.EncodeSecp256k1DERSignature(r, s)

	assert.NotEmpty(t, sig)
}

func TestBitcoinSig(t *testing.T) {
	t.Parallel()

	pubKey, _ := hex.DecodeString("04d2e670a19c6d753d1a6d8b20bd045df8a08fb162cf508956c31268c6d81ffdabab65528eefbb8057aa85d597258a3fbd481a24633bc9b47a9aa045c91371de52")
	msg, _ := hex.DecodeString("01020304")
	r, _ := hex.DecodeString("fef45d2892953aa5bbcdb057b5e98b208f1617a7498af7eb765574e29b5d9c2c")
	s, _ := hex.DecodeString("d47563f52aac6b04b55de236b7c515eb9311757db01e02cff079c3ca6efb063f")

	verifier, _ := NewSecp256()
	sig := verifier.EncodeSecp256k1DERSignature(r, s)
	err := verifier.VerifySecp256k1(pubKey, msg, sig, byte(ECDSADoubleSha256))

	assert.Nil(t, err)
}

func TestEthereumSig2(t *testing.T) {
	t.Parallel()

	msg, _ := hex.DecodeString("616161")
	key, _ := hex.DecodeString("044338845e8308b819bf33a43dc7f47713f92d8d377dfde399831e9d8da23446be32cef60a7c923332ab06c768242d11017a6bcf419c17b8b184fc19ea603b07d6")
	sig, _ := hex.DecodeString("3046022100da0db89620513df9a90cf8c97edf227e07182d1c91b3cab55a472122d639daee022100d5b9cf4a02274cf5b606df7b4fa73bff1190f54e0c6ef8cd362e63dc1dbecce1")
	verifier, _ := NewSecp256()
	err := verifier.VerifySecp256k1(key, msg, sig, byte(ECDSASha256))

	assert.Nil(t, err)
}

func TestSecp256_VerifySecp256r1(t *testing.T) {
	t.Parallel()

	msg := make([]byte, 100)
	_, _ = rand.Read(msg)
	verifier, _ := NewSecp256()

	sk, _ := GenPrivateKeyR1()

	key := elliptic.MarshalCompressed(elliptic.P256(), sk.PublicKey.X, sk.PublicKey.Y)
	sig, _ := SignMessage(msg, sk)

	err := verifier.VerifySecp256r1(nil, nil, nil)
	assert.NotNil(t, err)

	err = verifier.VerifySecp256r1(key, nil, nil)
	assert.NotNil(t, err)

	err = verifier.VerifySecp256r1(key, nil, sig)
	assert.NotNil(t, err)

	err = verifier.VerifySecp256r1(key, msg, sig)
	assert.Nil(t, err)

	msg[0] += 1
	err = verifier.VerifySecp256r1(key, msg, sig)
	assert.NotNil(t, err)
}

func GenPrivateKeyR1() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}

	return key, nil
}

func SignMessage(msg []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	digest := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, digest[:])
	if err != nil {
		return nil, err
	}

	normS := NormalizeS(s)
	return signatureRaw(r, normS), nil
}

// will serialize signature to R || S.
// R, S are padded to 32 bytes respectively.
// code roughly copied from secp256k1_nocgo.go
func signatureRaw(r *big.Int, s *big.Int) []byte {
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sigBytes := make([]byte, 64)

	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}
