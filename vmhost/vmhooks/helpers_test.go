package vmhooks

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockVMHost is a mock for the VMHost interface
type mockVMHost struct {
	mock.Mock
	crypto            vmhost.VMCrypto
	managedType       vmhost.ManagedTypesContext
	metering          vmhost.MeteringContext
	blockchain        vmhost.BlockchainContext
	enableEpochsH     vmhost.EnableEpochsHandler
	runtime           vmhost.RuntimeContext
}

func (m *mockVMHost) Crypto() vmhost.VMCrypto                            { return m.crypto }
func (m *mockVMHost) ManagedTypes() vmhost.ManagedTypesContext           { return m.managedType }
func (m *mockVMHost) Metering() vmhost.MeteringContext                    { return m.metering }
func (m *mockVMHost) Blockchain() vmhost.BlockchainContext                { return m.blockchain }
func (m *mockVMHost) EnableEpochsHandler() vmhost.EnableEpochsHandler    { return m.enableEpochsH }
func (m *mockVMHost) Runtime() vmhost.RuntimeContext                     { return m.runtime }
func (m *mockVMHost) FailExecution(err error)                            { m.Called(err) }
func (m *mockVMHost) GetGasSchedule() vmhost.GasSchedule                   { return nil }
func (m *mockVMHost) AreInSameShard(address1, address2 []byte) bool        { return true }
func (m *mockVMHost) IsBuiltinFunctionName(name string) bool               { return false }
func (m *mockVMHost) GetTxContext() vmhost.TxContext                       { return nil }
func (m *mockVMHost) GetLogEntries() []*vmhost.LogEntry                    { return nil }
func (m *mockVMHost) CompleteLogEntriesWithCallType(output *vmcommon.VMOutput, callType string) {}


// mockCryptoHook is a mock for the VMCrypto interface
type mockCryptoHook struct {
	mock.Mock
}

func (m *mockCryptoHook) Sha256(data []byte) ([]byte, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
func (m *mockCryptoHook) Keccak256(p []byte) ([]byte, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCryptoHook) Ripemd160(p []byte) ([]byte, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
// ... other crypto functions

// mockManagedTypesContext is a mock for the ManagedTypesContext interface
type mockManagedTypesContext struct {
	mock.Mock
}

func (m *mockManagedTypesContext) GetBytes(handle int32) ([]byte, error) {
	args := m.Called(handle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockManagedTypesContext) ConsumeGasForBytes(data []byte) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *mockManagedTypesContext) SetBytes(handle int32, data []byte) {
	m.Called(handle, data)
}

// mockMetering is a mock for the MeteringContext interface
type mockMetering struct {
	mock.Mock
}

func (m *mockMetering) UseGasBoundedAndAddTracedGas(name string, gas uint64) error {
	args := m.Called(name, gas)
	return args.Error(0)
}
func (m *mockMetering) GasLeft() uint64 {
	return 0
}
func (m *mockMetering) UseGas(gas uint64) error {
	return nil
}

// mockEnableEpochsHandler is a mock for the EnableEpochsHandler interface
type mockEnableEpochsHandler struct {
	mock.Mock
}

func (m *mockEnableEpochsHandler) IsFlagEnabled(flag vmhost.EpochFlag) bool {
	return true
}

func TestVMHooksImpl_managedHash(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		cryptoHook := &mockCryptoHook{}
		managedType := &mockManagedTypesContext{}
		metering := &mockMetering{}
		enableEpochsH := &mockEnableEpochsHandler{}
		host := &mockVMHost{
			crypto:        cryptoHook,
			managedType:   managedType,
			metering:      metering,
			enableEpochsH: enableEpochsH,
		}
		context := &VMHooksImpl{host: host}

		inputBytes := []byte("input")
		outputBytes := []byte("output")

		managedType.On("GetBytes", int32(1)).Return(inputBytes, nil)
		managedType.On("ConsumeGasForBytes", inputBytes).Return(nil)
		managedType.On("SetBytes", int32(2), outputBytes).Return()
		cryptoHook.On("Sha256", inputBytes).Return(outputBytes, nil)
		metering.On("UseGasBoundedAndAddTracedGas", "sha256", uint64(100)).Return(nil)

		result := context.managedHash(1, 2, "sha256", 100, cryptoHook.Sha256, vmhost.ErrSha256Hash)

		assert.Equal(t, int32(0), result)
		managedType.AssertExpectations(t)
		cryptoHook.AssertExpectations(t)
		metering.AssertExpectations(t)
	})

	t.Run("should fail on get bytes", func(t *testing.T) {
		t.Parallel()

		cryptoHook := &mockCryptoHook{}
		managedType := &mockManagedTypesContext{}
		metering := &mockMetering{}
		enableEpochsH := &mockEnableEpochsHandler{}
		host := &mockVMHost{
			crypto:        cryptoHook,
			managedType:   managedType,
			metering:      metering,
			enableEpochsH: enableEpochsH,
		}
		context := &VMHooksImpl{host: host}

		err := errors.New("err")
		managedType.On("GetBytes", int32(1)).Return(nil, err)
		metering.On("UseGasBoundedAndAddTracedGas", "sha256", uint64(100)).Return(nil)
		host.On("FailExecution", err).Return()

		result := context.managedHash(1, 2, "sha256", 100, cryptoHook.Sha256, vmhost.ErrSha256Hash)

		assert.Equal(t, int32(1), result)
		host.AssertExpectations(t)
	})
}

func TestVMHooksImpl_getSignatureOperands(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		managedType := &mockManagedTypesContext{}
		host := &mockVMHost{
			managedType: managedType,
		}
		context := &VMHooksImpl{host: host}

		keyBytes := []byte("key")
		msgBytes := []byte("msg")
		sigBytes := []byte("sig")

		managedType.On("GetBytes", int32(1)).Return(keyBytes, nil)
		managedType.On("ConsumeGasForBytes", keyBytes).Return(nil)
		managedType.On("GetBytes", int32(2)).Return(msgBytes, nil)
		managedType.On("ConsumeGasForBytes", msgBytes).Return(nil)
		managedType.On("GetBytes", int32(3)).Return(sigBytes, nil)
		managedType.On("ConsumeGasForBytes", sigBytes).Return(nil)

		k, m, s, err := context.getSignatureOperands(1, 2, 3)

		assert.Nil(t, err)
		assert.Equal(t, keyBytes, k)
		assert.Equal(t, msgBytes, m)
		assert.Equal(t, sigBytes, s)
		managedType.AssertExpectations(t)
	})
}
