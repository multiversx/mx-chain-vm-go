package hostCore

import (
	"errors"
	"math"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

func TestNewVMHost(t *testing.T) {
	blockchainHook := worldmock.NewMockWorld()
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	epochNotifier := &mock.EpochNotifierStub{}
	epochsHandler := &worldmock.EnableEpochsHandlerStub{}
	vmType := []byte("vmType")
	esdtTransferParser, err := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	require.Nil(t, err)

	makeHostParameters := func() *vmhost.VMHostParameters {
		return &vmhost.VMHostParameters{
			VMType:                    vmType,
			ESDTTransferParser:        esdtTransferParser,
			BuiltInFuncContainer:      bfc,
			EpochNotifier:             epochNotifier,
			EnableEpochsHandler:       epochsHandler,
			Hasher:                    worldmock.DefaultHasher,
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
		hostParameters.EnableEpochsHandler = &worldmock.EnableEpochsHandlerStub{
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

type opcodeConfigExecutorStub struct {
	executor.Executor
	setOpcodeConfigErr error
	version            executor.OpcodeVersion
	calls              int
}

func (e *opcodeConfigExecutorStub) SetOpcodeConfig(version executor.OpcodeVersion, _ *executor.WASMOpcodeCost) error {
	e.calls++
	if e.setOpcodeConfigErr != nil {
		return e.setOpcodeConfigErr
	}
	e.version = version
	return nil
}

type opcodeConfigRuntimeStub struct {
	vmhost.RuntimeContext
	exec          executor.Executor
	warmCacheWipe int
}

func (r *opcodeConfigRuntimeStub) GetVMExecutor() executor.Executor { return r.exec }
func (r *opcodeConfigRuntimeStub) ClearWarmInstanceCache()          { r.warmCacheWipe++ }

type opcodeConfigMeteringStub struct {
	vmhost.MeteringContext
	schedule config.GasScheduleMap
}

func (m *opcodeConfigMeteringStub) SetGasSchedule(schedule config.GasScheduleMap) {
	m.schedule = schedule
}

type opcodeConfigBlockchainStub struct {
	vmhost.BlockchainContext
	codesCleared int
}

func (b *opcodeConfigBlockchainStub) ClearCompiledCodes() { b.codesCleared++ }

func newOpcodeConfigTestHost(v2Enabled *bool) (*vmHost, *opcodeConfigExecutorStub, *opcodeConfigRuntimeStub, *opcodeConfigMeteringStub, *opcodeConfigBlockchainStub) {
	exec := &opcodeConfigExecutorStub{}
	runtime := &opcodeConfigRuntimeStub{exec: exec}
	metering := &opcodeConfigMeteringStub{}
	blockchain := &opcodeConfigBlockchainStub{}
	host := &vmHost{
		runtimeContext:    runtime,
		meteringContext:   metering,
		blockchainContext: blockchain,
		enableEpochsHandler: &worldmock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return flag == vmhost.OpcodeV2Flag && *v2Enabled
			},
		},
		activationEpochMap: map[uint32]struct{}{},
	}
	return host, exec, runtime, metering, blockchain
}

func TestVMHost_GetOpcodeVersionForCurrentEpoch(t *testing.T) {
	v2 := false
	host, _, _, _, _ := newOpcodeConfigTestHost(&v2)

	require.Equal(t, executor.OpcodeVersionV1, host.getOpcodeVersionForCurrentEpoch())
	v2 = true
	require.Equal(t, executor.OpcodeVersionV2, host.getOpcodeVersionForCurrentEpoch())
}

func TestVMHost_ApplyExecutionConfig(t *testing.T) {
	v2 := false

	t.Run("success clears caches", func(t *testing.T) {
		host, exec, runtime, metering, blockchain := newOpcodeConfigTestHost(&v2)
		schedule := config.MakeGasMapForTests()

		err := host.applyExecutionConfig(executor.OpcodeVersionV2, schedule)
		require.Nil(t, err)
		require.Equal(t, executor.OpcodeVersionV2, exec.version)
		require.Equal(t, executor.OpcodeVersionV2, host.opcodeVersion)
		require.Equal(t, schedule, host.gasSchedule)
		require.Equal(t, schedule, metering.schedule)
		require.Equal(t, 1, runtime.warmCacheWipe)
		require.Equal(t, 1, blockchain.codesCleared)
	})
	t.Run("invalid gas schedule", func(t *testing.T) {
		host, exec, runtime, _, blockchain := newOpcodeConfigTestHost(&v2)

		err := host.applyExecutionConfig(executor.OpcodeVersionV2, config.GasScheduleMap{})
		require.NotNil(t, err)
		require.Equal(t, 0, exec.calls)
		require.Equal(t, 0, runtime.warmCacheWipe)
		require.Equal(t, 0, blockchain.codesCleared)
	})
	t.Run("executor error leaves state untouched", func(t *testing.T) {
		host, exec, runtime, metering, blockchain := newOpcodeConfigTestHost(&v2)
		exec.setOpcodeConfigErr = errors.New("set opcode config failed")

		err := host.applyExecutionConfig(executor.OpcodeVersionV2, config.MakeGasMapForTests())
		require.Equal(t, exec.setOpcodeConfigErr, err)
		require.Equal(t, executor.OpcodeVersionV1, host.opcodeVersion)
		require.Nil(t, host.gasSchedule)
		require.Nil(t, metering.schedule)
		require.Equal(t, 0, runtime.warmCacheWipe)
		require.Equal(t, 0, blockchain.codesCleared)
	})
}

func TestVMHost_ReconfigureExecutorOpcodeVersion(t *testing.T) {
	v2 := false
	host, exec, runtime, _, blockchain := newOpcodeConfigTestHost(&v2)
	host.gasSchedule = config.MakeGasMapForTests()

	// unchanged version: nothing happens
	host.reconfigureExecutorOpcodeVersion()
	require.Equal(t, 0, exec.calls)

	// flag flips: executor reconfigured and caches wiped
	v2 = true
	host.reconfigureExecutorOpcodeVersion()
	require.Equal(t, 1, exec.calls)
	require.Equal(t, executor.OpcodeVersionV2, host.opcodeVersion)
	require.Equal(t, 1, runtime.warmCacheWipe)
	require.Equal(t, 1, blockchain.codesCleared)

	// executor failure is only logged, the old version is kept
	v2 = false
	exec.setOpcodeConfigErr = errors.New("boom")
	host.reconfigureExecutorOpcodeVersion()
	require.Equal(t, executor.OpcodeVersionV2, host.opcodeVersion)
}

func TestVMHost_GasScheduleChange(t *testing.T) {
	v2 := true
	host, exec, _, metering, _ := newOpcodeConfigTestHost(&v2)
	schedule := config.MakeGasMapForTests()

	host.GasScheduleChange(schedule)
	require.Equal(t, executor.OpcodeVersionV2, exec.version)
	require.Equal(t, schedule, metering.schedule)

	// error path is only logged
	exec.setOpcodeConfigErr = errors.New("boom")
	host.GasScheduleChange(schedule)
	require.Equal(t, 2, exec.calls)
}

func TestVMHost_EpochConfirmedReconfiguresExecutor(t *testing.T) {
	v2 := false
	host, exec, _, _, _ := newOpcodeConfigTestHost(&v2)
	host.gasSchedule = config.MakeGasMapForTests()

	v2 = true
	host.EpochConfirmed(5, 0)
	require.Equal(t, executor.OpcodeVersionV2, exec.version)
}
