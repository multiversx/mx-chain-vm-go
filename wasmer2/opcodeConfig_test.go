package wasmer2

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

func testOpcodeCosts(t *testing.T) *executor.WASMOpcodeCost {
	gasCostConfig, err := config.CreateGasConfig(config.MakeGasMapForTests())
	require.Nil(t, err)
	return gasCostConfig.WASMOpcodeCost
}

func TestWasmer2Executor_SetOpcodeConfig(t *testing.T) {
	wasmerExecutor, err := CreateExecutor()
	require.Nil(t, err)

	costs := testOpcodeCosts(t)
	for _, version := range []executor.OpcodeVersion{executor.OpcodeVersionV1, executor.OpcodeVersionV2} {
		require.Nil(t, wasmerExecutor.SetOpcodeConfig(version, costs))
		require.NotNil(t, wasmerExecutor.opcodeCost)
	}
}

func TestWasmer2ExecutorFactory_CreateExecutor_WithOpcodeCosts(t *testing.T) {
	factory := ExecutorFactory()

	for _, version := range []executor.OpcodeVersion{executor.OpcodeVersionV1, executor.OpcodeVersionV2} {
		exec, err := factory.CreateExecutor(executor.ExecutorFactoryArgs{
			OpcodeVersion: version,
			OpcodeCosts:   testOpcodeCosts(t),
		})
		require.Nil(t, err)
		require.NotNil(t, exec)
	}
}
