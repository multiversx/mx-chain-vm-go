package mockery

import (
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/stretchr/testify/mock"
)

type MockCryptoContext struct {
	mock.Mock
}

func (m *MockCryptoContext) Sha256(data []byte) ([]byte, error) {
	args := m.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptoContext) Keccak256(data []byte) ([]byte, error) {
	args := m.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptoContext) Ripemd160(data []byte) ([]byte, error) {
	args := m.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptoContext) VerifyBLS(key []byte, msg []byte, sig []byte) error {
	args := m.Called(key, msg, sig)
	return args.Error(0)
}

func (m *MockCryptoContext) VerifySignatureShare(publicKey []byte, message []byte, sig []byte) error {
	args := m.Called(publicKey, message, sig)
	return args.Error(0)
}

func (m *MockCryptoContext) VerifyAggregatedSig(pubKeysSigners [][]byte, message []byte, aggSig []byte) error {
	args := m.Called(pubKeysSigners, message, aggSig)
	return args.Error(0)
}

func (m *MockCryptoContext) VerifyEd25519(key []byte, msg []byte, sig []byte) error {
	args := m.Called(key, msg, sig)
	return args.Error(0)
}

func (m *MockCryptoContext) VerifySecp256k1(key []byte, msg []byte, sig []byte, hashType uint8) error {
	args := m.Called(key, msg, sig, hashType)
	return args.Error(0)
}

func (m *MockCryptoContext) EncodeSecp256k1DERSignature(r, s []byte) []byte {
	args := m.Called(r, s)
	return args.Get(0).([]byte)
}

func (m *MockCryptoContext) VerifySecp256r1(key []byte, msg []byte, sig []byte) error {
	args := m.Called(key, msg, sig)
	return args.Error(0)
}

var _ crypto.VMCrypto = &MockCryptoContext{}
