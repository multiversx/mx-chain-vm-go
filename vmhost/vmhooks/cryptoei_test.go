package vmhooks

import (
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	mock2 "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_Sha256(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)
	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)

	data := []byte("test data")
	hash := sha256.Sum256(data)
	crypto.Result = hash[:]

	ret := hooks.Sha256(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedSha256(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	data := []byte("test data")
	hash := sha256.Sum256(data)
	managedType.On("GetBytes", mock.Anything).Return(data, nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	crypto.Result = hash[:]
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.ManagedSha256(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_Keccak256(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	data := []byte("test data")
	hash := sha256.Sum256(data) // just a placeholder
	crypto.Result = hash[:]

	ret := hooks.Keccak256(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedKeccak256(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	data := []byte("test data")
	hash := sha256.Sum256(data) // just a placeholder
	managedType.On("GetBytes", mock.Anything).Return(data, nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	crypto.Result = hash[:]
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.ManagedKeccak256(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_Ripemd160(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	data := []byte("test data")
	hash := sha256.Sum256(data) // just a placeholder
	crypto.Result = hash[:]

	ret := hooks.Ripemd160(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedRipemd160(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	data := []byte("test data")
	hash := sha256.Sum256(data) // just a placeholder
	managedType.On("GetBytes", mock.Anything).Return(data, nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	crypto.Result = hash[:]
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.ManagedRipemd160(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_VerifyBLS(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	crypto.Err = nil

	ret := hooks.VerifyBLS(0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedVerifyBLS(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	crypto.Err = nil

	ret := hooks.ManagedVerifyBLS(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_VerifyEd25519(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	crypto.Err = nil

	ret := hooks.VerifyEd25519(0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedVerifyEd25519(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	crypto.Err = nil

	ret := hooks.ManagedVerifyEd25519(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_VerifyCustomSecp256k1(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	crypto.Err = nil

	instance := host.Runtime().GetInstance().(*mockery.MockInstance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return([]byte{0x30, 0x0}, nil)

	ret := hooks.VerifyCustomSecp256k1(0, secp256k1CompressedPublicKeyLength, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedVerifyCustomSecp256k1(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	crypto.Err = nil

	ret := hooks.ManagedVerifyCustomSecp256k1(0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_EncodeSecp256k1DerSignature(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)

	crypto.Result = []byte("signature")

	ret := hooks.EncodeSecp256k1DerSignature(0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedEncodeSecp256k1DerSignature(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	crypto := &mock2.CryptoHookMock{}
	host.On("Crypto").Return(crypto)
	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	crypto.Result = []byte("signature")
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.ManagedEncodeSecp256k1DerSignature(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_CreateEC(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	managedType.On("PutEllipticCurve", mock.Anything).Return(int32(1))
	instance := host.Runtime().GetInstance().(*mockery.MockInstance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return([]byte("p256"), nil)

	ret := hooks.CreateEC(0, 4)
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_ManagedCreateEC(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	managedType.On("GetBytes", mock.Anything).Return([]byte("p256"), nil)
	managedType.On("PutEllipticCurve", mock.Anything).Return(int32(1))

	ret := hooks.ManagedCreateEC(0)
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_EllipticCurve(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	ec := elliptic.P256()
	managedType.On("GetEllipticCurve", mock.Anything).Return(ec, nil)
	managedType.On("GetTwoBigInt", mock.Anything, mock.Anything).Return(big.NewInt(0), big.NewInt(0), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(0), nil)
	managedType.On("Get100xCurveGasCostMultiplier", mock.Anything).Return(int32(100))
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything).Return(nil)

	hooks.AddEC(0, 0, 0, 0, 0, 0, 0)
	hooks.DoubleEC(0, 0, 0, 0, 0)
	hooks.IsOnCurveEC(0, 0, 0)
	hooks.ScalarBaseMultEC(0, 0, 0, 0, 0)
	hooks.ScalarMultEC(0, 0, 0, 0, 0, 0, 0)
	hooks.MarshalEC(0, 0, 0, 0)
	hooks.MarshalCompressedEC(0, 0, 0, 0)
	hooks.UnmarshalEC(0, 0, 0, 0, 0)
	hooks.UnmarshalCompressedEC(0, 0, 0, 0, 0)
	hooks.GenerateKeyEC(0, 0, 0, 0)
	hooks.GetCurveLengthEC(0)
	hooks.GetPrivKeyByteLengthEC(0)
	hooks.EllipticCurveGetValues(0, 0, 0, 0, 0, 0)
}
