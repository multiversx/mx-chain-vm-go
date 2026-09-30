package hostCore

import (
	"math"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	hostmock "github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

func TestNewVMHost(t *testing.T) {
	blockchainHook := &contextmock.BlockchainHookStub{}
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	epochNotifier := &hostmock.EpochNotifierStub{}
	vmType := []byte("vmType")
	esdtTransferParser, err := parsers.NewESDTTransferParser(&mock.MarshalizerMock{})
	require.Nil(t, err)

	makeHostParameters := func() *vmhost.VMHostParameters {
		epochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		}

		return &vmhost.VMHostParameters{
			VMType:                    vmType,
			ESDTTransferParser:        esdtTransferParser,
			BuiltInFuncContainer:      bfc,
			EpochNotifier:             epochNotifier,
			EnableEpochsHandler:       epochsHandler,
			Hasher:                    &mock.HasherMock{},
			MapOpcodeAddressIsAllowed: map[string]map[string]struct{}{},
		}
	}

	t.Run("NilBlockchainHook", func(t *testing.T) {
		host, err := NewVMHost(nil, makeHostParameters())
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilBlockChainHook)
	})
	t.Run("NilHostParameters", func(t *testing.T) {
		host, err := NewVMHost(blockchainHook, nil)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilHostParameters)
	})
	t.Run("NilESDTTransferParser", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.ESDTTransferParser = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilESDTTransferParser)
	})
	t.Run("NilBuiltInFunctionsContainer", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.BuiltInFuncContainer = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilBuiltInFunctionsContainer)
	})
	t.Run("NilEpochNotifier", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EpochNotifier = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilEpochNotifier)
	})
	t.Run("NilEnableEpochsHandler", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EnableEpochsHandler = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilEnableEpochsHandler)
	})
	t.Run("InvalidEnableEpochsHandler", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EnableEpochsHandler = &mock.EnableEpochsHandlerStub{
			IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
				return false
			},
		}
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, core.ErrInvalidEnableEpochsHandler)
	})
	t.Run("NilHasher", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.Hasher = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilHasher)
	})
	t.Run("NilVMType", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.VMType = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilVMType)
	})
}

func TestValidateVMInput(t *testing.T) {
	vmInput := &vmcommon.VMInput{
		GasProvided: 0,
	}

	vmInput.GasProvided = math.MaxUint64
	err := validateVMInput(vmInput)
	require.ErrorIs(t, err, vmhost.ErrInvalidGasProvided)

	vmInput.GasProvided = math.MaxInt64
	err = validateVMInput(vmInput)
	require.Nil(t, err)
}
