package host

import (
	"testing"

	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen/mock"
	worldmock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/world"
	"github.com/stretchr/testify/require"
)

func TestNewArwenVM(t *testing.T) {
	blockchainHook := worldmock.NewMockWorld()
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	epochNotifier := &mock.EpochNotifierStub{}
	epochsHandler := &mock.EnableEpochsHandlerStub{}
	vmType := []byte("vmType")
	esdtTransferParser, err := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	require.Nil(t, err)

	makeHostParameters := func() *arwen.VMHostParameters {
		return &arwen.VMHostParameters{
			VMType:               vmType,
			ESDTTransferParser:   esdtTransferParser,
			BuiltInFuncContainer: bfc,
			EpochNotifier:        epochNotifier,
			EnableEpochsHandler:  epochsHandler,
			Hasher:               worldmock.DefaultHasher,
		}
	}

	t.Run("NilBlockchainHook", func(t *testing.T) {
		host, err := NewArwenVM(nil, makeHostParameters())
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilBlockChainHook)
	})
	t.Run("NilHostParameters", func(t *testing.T) {
		host, err := NewArwenVM(blockchainHook, nil)
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilHostParameters)
	})
	t.Run("NilESDTTransferParser", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.ESDTTransferParser = nil
		host, err := NewArwenVM(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilESDTTransferParser)
	})
	t.Run("NilBuiltInFunctionsContainer", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.BuiltInFuncContainer = nil
		host, err := NewArwenVM(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilBuiltInFunctionsContainer)
	})
	t.Run("NilEpochNotifier", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EpochNotifier = nil
		host, err := NewArwenVM(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilEpochNotifier)
	})
	t.Run("NilEnableEpochsHandler", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EnableEpochsHandler = nil
		host, err := NewArwenVM(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilEnableEpochsHandler)
	})
	t.Run("NilHasher", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.Hasher = nil
		host, err := NewArwenVM(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilHasher)
	})
	t.Run("NilVMType", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.VMType = nil
		host, err := NewArwenVM(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, arwen.ErrNilVMType)
	})
}
